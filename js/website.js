function run() {

    var sprintf = require("sprintf").sprintf;

    var publicRE = new RegExp("@public", 'm');
    var commentRE= new RegExp("^\s*//\s*");
    var re = new RegExp("go|js$");
    var skip = [new RegExp("^vendor"), new RegExp("^website")];
    var core = [];
    var handlers = [];
    filepath.walk(".", function(path, info, err) {
        if (!info.IsDir) {
	    var results = filepath.split(path);
	    var dir = results[0];
	    var file = results[1];
	    for (i in skip) {
		if (skip[i].exec(dir)) {
		    return;
		}
	    }
	    if (re.exec(path)) {
		var results = fs.read(path);
		if (results[1]) {
		    console.log("Error", results[1]);
		    os.exit(1);
		}
		var ext = filepath.ext(file);
		var extRE = new RegExp(ext);
		var contents = results[0].split("\n");
		var comments = [];
		var current = "";
		for (var lineNo in contents) {
		    var line = contents[lineNo];
		    if (line.match(commentRE)) {
			var stripped = line.replace(commentRE, '');
			current = current + stripped + "\n";
		    } else {
			if (current != "") {
			    if (current.match(publicRE)) {
				comments.push(current.replace(publicRE, ""));
			    }
			    current = "";
			}
		    }
		}
		var md = "";
		_.each(comments, function (c) {
		    md = md + c + "\n";
		});
		if (comments.length > 0) {
		    var path;
		    if (ext === ".go") {
			path = filepath.join("website", 
					     "core_" + file.replace(extRE, ".md"));
			core.push(path);
		    } else if (ext === ".js") {
			path = filepath.join("website", 
					     "handler_" + file.replace(extRE, ".md"));
			handlers.push(path);
		    }
		    console.log(sprintf("%32s: Updated", path));
		    fs.write(path, md, 0644);
		} else {
		    console.log(sprintf("%32s: No comments found", file));
		}
	    }
	}
    });
    
    // CORE
    var n = core.length / 3;
    var contents = "div.container-fluid\n  div.row\n";
    var lists = _.chain(core).groupBy(function(element, index) {
	return Math.floor(index/n);
    }).toArray().value();
    _.each(lists, function(l) {
	contents = contents + "    div.col-md-4\n      ul.list-unstyled\n";
	_.each(l, function(file) {
	    var results = filepath.split(file);
	    var text = results[1].replace(/core_(.*).md/, '$1');
	    var file = results[1].replace(/(core_.*).md/, '$1');
	    contents = contents + 
		sprintf("        li: a(href='%s') %s\n", file, text);
	});
    });
    fs.write("website/core.jade", contents, 0644);

    // HANDLERS
    var n = handlers.length / 3;
    var contents = "div.container-fluid\n  div.row\n";
    var lists = _.chain(handlers).groupBy(function(element, index) {
	return Math.floor(index/n);
    }).toArray().value();
    _.each(lists, function(l) {
	contents = contents + "    div.col-md-4\n      ul.list-unstyled\n";
	_.each(l, function(file) {
	    var results = filepath.split(file);
	    var text = results[1].replace(/handler_(.*).md/, '$1');
	    var file = results[1].replace(/(handler_.*).md/, '$1');
	    contents = contents + 
		sprintf("        li: a(href='%s') %s\n", file, text);
	});
    });
    fs.write("website/handlers.jade", contents, 0644);
}

