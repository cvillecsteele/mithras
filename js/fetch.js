var semver = require("semver");
var sprintf = require("sprintf").sprintf;

function unpack(dir, file, name) {
    var path = filepath.join(dir, name)
    fs.removeAll(path);
    fs.mkdir(path, 0777);

    var results = exec.run(sprintf("tar -C %s -xvzf %s --strip-components 1", path, file));
    if (results[3] != 0) {
	log(results[1].trim());
	log(results[3]);
	os.exit(3);
    }
}

function install(sourceDir, name) {
    var src = filepath.join(sourceDir, name);
    var dest = filepath.join(mithras.DESTDIR || "./js", name);
    if (mithras.verbose) {
	log0(sprintf("Installing '%s'", dest));
    }
    fs.rename(src, dest);
}

function download(url, tempDir) {
    var re = new RegExp("[^/]+$")
    var tarball = re.exec(url)[0];
    var dest = filepath.join(tempDir, tarball);
    fs.removeAll(dest);
    var results = exec.run(sprintf("curl -o %s %s", dest, url));
    if (results[3] != 0) {
	log(results[1].trim());
	log(results[3]);
	os.exit(3);
    }
    return dest;
}
function getTags(package) {
    var results = web.get("https://registry.npmjs.org/-/package/" + package + "/dist-tags");
    if (results != "") {
	return JSON.parse(results);
    } else {
	log(sprintf("Error getting tags for package '%s'", package));
	os.exit(3);
    }
}
function depTree(package, versionOrTag) {
    var results = web.get("https://registry.npmjs.org/" + package);
    var info;
    if (results != "") {
	info = JSON.parse(results);
    } else {
	log(sprintf("Error getting info for package '%s'", package));
	os.exit(3);
    }

    var tags = getTags(package);
    var targetVersion;
    if (versionOrTag === "") {
	targetVersion = tags.latest;
    } else {
	if (tags[versionOrTag]) {
	    versionOrTag = tags[versionOrTag];
	}
	targetVersion = semver.maxSatisfying(Object.keys(info.versions), versionOrTag);
    }

    // Does it exist?    
    if (!info.versions[targetVersion]) {
	log(sprintf("Invalid version '%s' for package '%s'", targetVersion, package));
	os.exit(3);
    }

    var deps = info.versions[targetVersion].dependencies;
    var dist = info.versions[targetVersion].dist;
    var name = info.versions[targetVersion].name;
    return {
	version: targetVersion, 
	dist: dist
	name: name
	deps: _.mapObject(deps, function(val, key) {
	    return depTree(key, val);
	})
    };
}
function run() {
    var source = mithras.ARGS[0];
    if (!source) {
	log("No source repository specified.");
	os.exit(3);
    }

    var re = new RegExp("([^@]+)(@?([^@]+))?$")
    var match = source.match(re);
    var package = match[1];
    var version = match[3] || "";

    var tree = depTree(package, version);
    mithras.traverse(tree).map(function (url) {
	if (this.key === "tarball") {
	    var name = this.parent.parent.node.name;
	    var tempDir = fs.tempDir();
	    var dest = download(url, tempDir);
	    unpack(tempDir, dest, name);
	    install(tempDir, name);
	}
    });
    
    return true;
}
