// MITHRAS: Javascript configuration management tool for AWS.
// Copyright (C) 2016, Colin Steele
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//   the Free Software Foundation, either version 3 of the License, or
//                  (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//              GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf,
        vsprintf = require("sprintf.js").vsprintf;

    var handler = {
	moduleName: "copy"
	run: function () {
	    var sprintf = require("sprintf.js").sprintf;

	    var stat = fs.stat(params.dest);
	    if (params.owner) {
		var u = user.lookup(params.owner)[0];
		if (!u) {
		    console.log(sprintf("User '%s' does not exist", params.owner));
		    os.exit(1);
		}
	    }
	    if (stat.Err) {
		// File does not exist...
		if (params.ensure === 'present') {
		    var error;
		    if (params.content) {
			error = fs.write(params.dest, params.content, params.mode);
			if (error) {
			    console.log("Copy error", JSON.stringify(error, null, 2));
			    os.exit(2);
			}
		    } else {
			error = fs.copy(params.src, params.dest, params.mode);
			if (error) {
			    console.log("Copy error", JSON.stringify(error, null, 2));
			    os.exit(3);
			}
		    }

		    var stat = fs.stat(params.dest);
		    if (stat.Err) {
		    	console.log("Copy error", JSON.stringify(stat, null, 2));
			os.exit(4);
		    }

		    if (params.owner) {
			result = user.lookup(params.owner);
			user = result[0];
			error = result[1];
			if (error) {
			    console.log("Copy error user lookupg", 
					JSON.stringify(error, null, 2));
			    os.exit(5);
			}
			error = fs.chown(params.dest, user.Uid);
			if (error) {
			    console.log("Copy error chown", 
					JSON.stringify(error, null, 2));
			    os.exit(5);
			}
		    }

		    console.log("Created");
		    return true;
		}
	    } else {
		// File is there
		switch(params.ensure) {
		case "absent":
		    var error = fs.remove(params.dest);
		    if (error) {
		    	console.log("Remove error", JSON.stringify(error, null, 2));
			os.exit(6);
		    }
		    console.log("Removed");
		    return true;
		    break;
		case "present":
		    if (params.content && params.overwrite) {
			error = fs.write(params.dest, params.content, params.mode);
		    }
		    if (error) {
		    	console.log("Overwrite error", JSON.stringify(error, null, 2));
			os.exit(7);
		    }
		    console.log("Overwritten");
		    return true;
		}
	    }
	    console.log("No action taken");
	    return true;
	}
	handle: function(catalog, resources, resource) {
	    if (resource.module != handler.moduleName) {
		return [null, false];
	    }
	    var js = sprintf("var params = %s;\n", JSON.stringify(_.omit(resource.params, 'hosts')));
	    var js = sprintf("%svar run = %s;\n", js, handler.run.toString());
	    for (var i in resource.params.hosts) {
		var instance = resource.params.hosts[i];
		var result = mithras.remote.mithras(instance, 
						    mithras.sshUserForInstance(resource, instance), 
						    mithras.sshKeyPathForInstance(resource, instance), 
						    js,
						    resource.params.become,
						    resource.params.becomeUser,
						    resource.params.becomeMethod);
		if (result[3] == 0) {
		    log(sprintf("File '%s' %s: %s.", 
				resource.params.dest, 
				resource.params.ensure,
			       result[0].trim()));
		}
	    }
	    return [null, true];
	}
    };
    
    handler.init = function () {
	mithras.modules.handlers.register("copy", handler.handle);
    };
    
    return handler;
});
