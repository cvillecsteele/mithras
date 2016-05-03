function run() {

    var sprintf = require("sprintf").sprintf;

    function lowercaseFirstLetter(string) {
	return string.charAt(0).toLowerCase() + string.slice(1);
    }

    var publicRE = new RegExp("@public", 'm');
    var modRE = new RegExp("@module", 'm');
    var exRE = new RegExp("@example", 'm');
    var commentRE= new RegExp("^[ \t]*//");

    var re = new RegExp("go|js$");
    var skip = [new RegExp("^vendor"), new RegExp("^website")];
    var core = [];
    var handlers = [];
    var mod = [];
    var example = [];
    var object = [];

    // Deal with object doc first
    var where = filepath.join(mithras.HOME, "website", "objects");
    var result = fs.dir(where);
    if (result[1]) {
	log(sprintf("Can't read dir '%s'", where));
	os.exit(1);
    }
    var files = result[0];
    _.each(files, function(f) {
	var ext = filepath.ext(f);
	if (ext === ".md") {
	    var result = filepath.split(f);
	    var dir = result[0];
	    var file = result[1];
	    object.push(file);
	}
    });
    var n = object.length / 3;
    var contents = "div.container-fluid\n  div.row\n";
    var lists = _.chain(object).groupBy(function(element, index) {
	return Math.floor(index/n);
    }).toArray().value();
    _.each(lists, function(l) {
	contents = contents + "    div.col-md-4\n      ul.list-unstyled\n";
	_.each(l, function(file) {
	    var results = filepath.split(file);
	    var text = results[1].replace(/(.*).md/, '$1');
	    var file = results[1].replace(/(.*).md/, 
					  function (match, p1, offset, string) {
					      return "objects/" + p1
					  });
	    contents = contents + 
		sprintf("        li: a(href='%s.html') %s\n", file, text);
	});
    });
    fs.write("website/objects.jade", contents, 0644);

    
    filepath.walk(".", function(path, info, err) {
        if (info.IsDir || (!re.exec(path))) {
	    return;
	}
	var results = filepath.split(path);
	var dir = results[0];
	var file = results[1];
	for (i in skip) {
	    if (skip[i].exec(dir)) {
		return;
	    }
	}
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
		// Out of a comment section... If we already have one
		// built up, it means it's time to use it.
		if (current != "") {
		    // Only pay attention to public comments
		    if (current.match(publicRE)) {
			// If it's the first public comment, do special stuff
			if (comments.length == 0) {
			    if (current.match(modRE)) {
				var f = "mod_" + file.replace(extRE, ".md");
				path = filepath.join("website", f);
				mod.push(path);
			    } else if (current.match(exRE)) {
				var f = "example_" + file.replace(extRE, ".md");
				path = filepath.join("website", f);
				example.push(path);
			    } else if (ext === ".go") {
				var f = "core_" + file.replace(extRE, ".md");
				path = filepath.join("website", f);
				core.push(path);
			    } else if (ext === ".js") {
				var f = "handler_" + file.replace(extRE, ".md");
				path = filepath.join("website", f);
				handlers.push(path);
			    }
			}
			if (current.match(modRE)) {
			    comments.
				push(current.replace(publicRE, "").
				     replace(modRE, ""));
			} else if (current.match(exRE)) {
			    comments.
				push(current.replace(publicRE, "").
				     replace(exRE, ""));
			} else if (ext === ".go") {
			    comments.push(current.replace(publicRE, ""));
			} else if (ext === ".js") {
			    comments.push(current.replace(publicRE, ""));
			}
		    }
		    current = "";
		}
	    }
	}
	var md = "";
	_.each(comments, function (c) {
	    md = md + c + "\n";
	});
	console.log(sprintf("%32s: %d comments.", file, comments.length));
	if (comments.length > 0) {
	    fs.write(path, md, 0644);	    
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
		sprintf("        li: a(href='%s.html') %s\n", file, text);
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
		sprintf("        li: a(href='%s.html') %s\n", file, text);
	});
    });
    fs.write("website/handlers.jade", contents, 0644);

    // MODULES
    var n = mod.length / 3;
    var contents = "div.container-fluid\n  div.row\n";
    var lists = _.chain(mod).groupBy(function(element, index) {
	return Math.floor(index/n);
    }).toArray().value();
    _.each(lists, function(l) {
	contents = contents + "    div.col-md-4\n      ul.list-unstyled\n";
	_.each(l, function(file) {
	    var results = filepath.split(file);
	    var text = results[1].replace(/mod_(.*).md/, '$1');
	    var file = results[1].replace(/(mod_.*).md/, '$1');
	    contents = contents + 
		sprintf("        li: a(href='%s.html') %s\n", file, text);
	});
    });
    fs.write("website/modules.jade", contents, 0644);

    // EXAMPLES
    var n = example.length / 3;
    var contents = "div.container-fluid\n  div.row\n";
    var lists = _.chain(example).groupBy(function(element, index) {
	return Math.floor(index/n);
    }).toArray().value();
    _.each(lists, function(l) {
	contents = contents + "    div.col-md-4\n      ul.list-unstyled\n";
	_.each(l, function(file) {
	    var results = filepath.split(file);
	    var text = results[1].replace(/example_(.*).md/, '$1');
	    var file = results[1].replace(/(example_.*).md/, '$1');
	    contents = contents + 
		sprintf("        li: a(href='%s.html') %s\n", file, text);
	});
    });
    fs.write("website/examples.jade", contents, 0644);

    // BUILD

    fs.chdir("website");
    console.log(JSON.stringify(exec.run("harp compile", "", {
	PATH: os.getenv("PATH")
	HOME: os.getenv("HOME")
    })));
}

