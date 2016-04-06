(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleName: "file"
	run: function (params) {
	    var sprintf = require("sprintf.js").sprintf;

	    // If directory, all immediate subdirectories will be
	    // created if they do not exist. If file, the file will
	    // NOT be created if it does not exist, see the copy or
	    // template module if you want that behavior. If link, the
	    // symbolic link will be created or changed. Use hard for
	    // hardlinks. If absent, directories will be recursively
	    // deleted, and files or symlinks will be unlinked. If
	    // touch, an empty file will be created if the path does
	    // not exist, while an existing file or directory will
	    // receive updated file access and modification times
	    // (similar to the way `touch` works from the command
	    // line).
	    var ensure = params.ensure;

	    // name of the user that should own the file/directory, as
	    // would be fed to chown
	    var chown = params.owner;

	    // mode the file or directory should be. For those used to
	    // /usr/bin/chmod remember that modes are actually octal
	    // numbers (like 0644). Leaving off the leading zero will
	    // likely have unexpected results
	    var mode = params.mode;

	    // path of the file to link to (applies only to
	    // ensure=link). Will accept absolute, relative and
	    // nonexisting paths. Relative paths are not expanded.
	    var src = params.src;

	    // Target
	    var dest = params.dest;

	    var stat = fs.stat(dest);
	    if (chown) {
		var u = user.lookup(chown)[0];
		if (!u) {
		    console.log(sprintf("User '%s' does not exist", chown));
		    os.exit(1);
		}
	    }
	    var present = !stat.Err;

	    switch(ensure) {
	    case "directory":
		if (!present) {
		    error = fs.mkdirAll(dest, mode);
		    if (error) {
			console.log("Create directory error", 
				    JSON.stringify(error, null, 2));
			os.exit(3);
		    }
		    // Check it.
		    var stat = fs.stat(dest);
		    if (stat.Err) {
			console.log(sprintf("Error creating directory '%s': %s", 
					    dest,
					    JSON.stringify(stat, null, 2)));
			os.exit(4);
		    }
		    console.log("Created");
		    return true;
		}
		log(sprintf("Found '%s', no action taken", dest));
		return true;
		break;
	    case "file":
	    case "touch":
		if (present) {
		    break;
		}
		result = fs.create(dest);
		f = result[0];
		error = result[1];
		if (error) {
		    console.log("Create file error", 
				JSON.stringify(error, null, 2));
		    os.exit(3);
		}
		error = fs.close(f);
		if (error) {
		    console.log("Create file close error", 
				JSON.stringify(error, null, 2));
		    os.exit(3);
		}

		// Check it.
		var stat = fs.stat(dest);
		if (stat.Err) {
		    console.log("Create file error", JSON.stringify(stat, null, 2));
		    os.exit(4);
		}

		// Chown it.
		if (chown) {
		    result = user.lookup(chown);
		    user = result[0];
		    error = result[1];
		    if (error) {
			console.log("Copy error user lookupg", 
				    JSON.stringify(error, null, 2));
			os.exit(5);
		    }
		    error = fs.chown(dest, user.Uid);
		    if (error) {
			console.log("Copy error chown", 
				    JSON.stringify(error, null, 2));
			os.exit(5);
		    }
		}

		// Chmod it.
		if (mode) {
		    error = user.chown(dest, u.Uid);
		    if (error) {
			console.log("Create file error", JSON.stringify(error, null, 2));
			os.exit(6);
		    }
		}

		if (ensure === "touch") {
		    error = user.chtimes(dest)
		    if (error) {
			console.log("Create chtimes error", 
				    JSON.stringify(error, null, 2));
			os.exit(7);
		    }
		}

		console.log(sprintf("Created '%s'", dest));
		break;
	    case "link":
		error = fs.symlink(src, dest);
		if (error) {
		    if (error.Err == 17) {
			console.log(sprintf("%s already exists", dest));
		    } else {
			console.log("Create symlink error", JSON.stringify(error, null, 2));
			os.exit(8);
		    }
		} else {
		    console.log(sprintf("Created '%s'", dest));
		}
		break;
	    case "hard":
		error = fs.link(dest, src);
		if (error) {
		    console.log("Create hard link error", 
				JSON.stringify(error, null, 2));
		    os.exit(9);
		}
		break;
	    case "absent":
		error = sanitize(fs.removeAll(dest));
		if (error) {
		    console.log(sprintf("Can't remove '%s'", dest),
				JSON.stringify(error, null, 2));
		    os.exit(10);
		}
		break;
	    default:
		console.log(sprintf("Invalid 'ensure': %s", ensure))
		os.exit(11);
	    }
	}
	handle: function(catalog, resources, resource) {
	    if (resource.module != handler.moduleName) {
		return [null, false];
	    }
	    var params = resource.params;
	    if (params.hosts) {
		var js = sprintf("var run = function() {\n (%s)(%s); };\n", 
				 handler.run.toString(),
				 JSON.stringify(_.omit(params, 'hosts')));
		for (var i in resource.params.hosts) {
		    var instance = resource.params.hosts[i];
		    var result = mithras.remote.mithras(instance, 
							mithras.sshUserForInstance(resource, 
										   instance), 
							mithras.sshKeyPathForInstance(resource, 
										      instance), 
							js,
							resource.params.become,
							resource.params.becomeUser,
							resource.params.becomeMethod);
		    if (result[3] == 0) {
			var out = result[0].trim();
			log(sprintf("File '%s' %s: %s.", 
				    resource.params.dest,
				    resource.params.ensure,
				    out != "" ? out : "success"));
		    }
		}
	    } else {
		handler.run(params);
	    }
	    return [null, true];
	}
    };
    
    handler.init = function () {
	mithras.modules.handlers.register("file", handler.handle);
    };
    
    return handler;
});
