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

// @public
// 
// # File
// 
// File is resource handler for manipulating files.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"file"`
// 
// Usage:
// 
// `var file = require("file").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rFile = {
//   name: "someFile"
//   module: "file"
//   dependsOn: [otherResource.name]
//   params: {
//     dest: "/etc/foo/bar"
//     src: scp://localhost/file.txt
//     ensure: "file"
//     mode: 0644
//     hosts: [<array of ec2 instance objects>]
//   }
// };
// ```
// 
// ## Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "absent", "file", "directory", "link", "hard", "touch"
//
// If `"directory"`, all immediate subdirectories will be created if
// they do not exist. If `"file"`, the file will NOT be created if it
// does not exist, see the copy or template module if you want that
// behavior. If `"link"`, the symbolic link will be created or
// changed. Use `"hard"` for hardlinks. If `"absent"`, directories
// will be recursively deleted, and files or symlinks will be
// unlinked. If `"touch"`, an empty file will be created if the path
// does not exist, while an existing file or directory will receive
// updated file access and modification times (similar to the way
// `touch` works from the command line).
// 
// ### `become`
//
// * Required: false
// * Allowed Values: true or false
//
// If `true`, the copy will attempt to run with escalated privs, as
// specified in the properties `becomeMethod` and `becomeUser`.
// 
// ### `becomeMethod`
//
// * Required: false
// * Allowed Values: "su" or "sudo"
//
// The method of privilege escalation.
// 
// ### `becomeUser`
//
// * Required: false
// * Allowed Values: any string specifying a username suitable for use by `becomeMethod`
//
// ### `hosts`
//
// * Required: true
// * Allowed Values: an array of ec2 instance objects
//
// This property specifies the hosts on which this resource is to be applied.
//
// ### `mode`
//
// * Required: true
// * Allowed Values: octal number specifying a valid permission mask
//
// This property specifies the path to the file/link/directory to be manipulated
//
// ### `owner`
//
// * Required: false
// * Allowed Values: username of the user to which the file will be `chown`'ed
//
// This property specifies the path to the file/link/directory to be manipulated
//
// ### `dest`
//
// * Required: true
// * Allowed Values: a valid path on the target host
//
// This property specifies the path to the file/link/directory to be manipulated
//
// ### `src`
//
// * Required: false
// * Allowed Values: a valid path on the target host
//
// If ensure=`"file"` or `"directory"`, the value of thie property may
// take one of three forms.  If of the form
// `"scp://localhost/foo/bar"`, then the *local* file specified by the
// `src` is SCP'd to the remote host, to the value of `dest`.  If of
// the form `"http://www.someplace.com/foo/bar"`, then from the remote
// instance, an HTTP GET request is performed to the value of `src`,
// and the contents of the response are written to `dest`.
//
// If ensure=`"link"`, specifies the path of the file to link to.
// Will accept absolute, relative and nonexisting paths. Relative
// paths are not expanded.
//
// ### `content`
//
// * Required: false
// * Allowed Values: a string of file contents to be written
//
// If ensure=`"file"`, the value of this property (presumably a
// string) will be written to `dest`.
//
// ### `force`
//
// * Required: false
// * Allowed Values: boolean
//
// If `true`, any `file` will be overwritten, even if it already exists.
//
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleName: "file"

        preflight: function(catalog, resources, resource) {
	    if (resource.module != handler.moduleName) {
		return [null, false];
	    }

            var p = resource.params;
            var ensure = p.ensure;

            // Sanity
            if (!p || !p.dest) {
                console.log(sprintf("Invalid resource params: %s", 
                                    JSON.stringify(p)));
                os.exit(3);
            }

            // Any hosts?
            if (!Array.isArray(p.hosts)) {
                log("No hosts.");
                return [null, true];
            }

            // Loop over hosts
            var target = resource._target = {};
            _.each(p.hosts, function(host) {
                var key = mithras.sshKeyPathForInstance(resource, host);
                var user = mithras.sshUserForInstance(resource, host);

                // update for by-host variance
                _.find(resources, function(r) {
                    return r.name === resource.name;
                })._currentHost = host;
                resource._currentHost = host;
                var updated = mithras.updateResource(resource, 
                                                     catalog, 
                                                     resources,
                                                     resource.name);
                var updatedParams = updated.params;

                cmd = sprintf("test -e '%s' && echo 'found'", updatedParams.dest);
                var result = mithras.remote.shell(host.PublicIpAddress, 
                                                  user, 
                                                  key, 
                                                  "",
                                                  cmd, 
                                                  {});
                var out = result[0] ? result[0].trim() : "";
                var err = result[1] ? result[1].trim() : "";
                if (result[3] == 0) {
                    log(sprintf("File '%s': %s.", 
                                updatedParams.dest,
                                out != "" ? out : "success"));
                    target[host.PublicIpAddress] = out;
                    target[host.InstanceId] = out;
                }
            });
            return [target, true];
        }
   
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

            // Overwrite
            var force = params.force;

	    var stat = fs.stat(dest);
	    var uid;
	    var gid;
	    if (chown) {
		var result = user.lookup(chown);
		var u = result[0];
		var e = result[1];
		if (!u) {
		    // Ugly hack alert.  Apparently cross-compiled go
		    // binaries for linux don't have a robust (read:
		    // working) way to do user lookup.  We fall back
		    // to exec'ing 'id'...
		    var result = exec.run(sprintf("id %s", chown));
		    if (result[3] == 0) {
			out = result[0].trim();
			var match = out.match(/^uid=([0-9]+)/);
			if (match) {
			    uid = parseInt(match[1], 10);
			}
			var match = out.match(/gid=([0-9]+)/);
			if (match) {
			    gid = parseInt(match[1], 10);
			}
		    } else {
			console.log(sprintf("User '%s' does not exist", chown));
			os.exit(3);
		    }
		} else {
		    uid = user.Uid;
		    gid = user.Gid;
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
		if (!present || force) {
                    if (params.content) {
                        error = fs.write(params.dest, params.content, params.mode);
                        if (error) {
                            console.log("File write error", JSON.stringify(error, null, 2));
                            os.exit(3);
                        }
                    } else if (params.src) {
			var url = web.url.parse(params.src);
			var scheme = url.scheme;
			switch (scheme) {
			case "scp":
			    break;
                        case "s3":
                            var bucket = url.host;
                            var path = url.path;
                            aws.s3.objects.writeInto(params.region, bucket, path, params.dest, params.mode);
                            break;
			case "http":
			case "https":
			    web.get(params.src, params.dest, params.mode);
			    break;
			default:
                            error = fs.copy(params.src, params.dest, params.mode);
                            if (error) {
				console.log("File copy error", JSON.stringify(error, null, 2));
				os.exit(3);
                            }
			    break;
			}
		    } else {
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
		    }
		}

		// Check it.
		var stat = fs.stat(dest);
		if (stat.Err) {
		    console.log("Create file error", JSON.stringify(stat, null, 2));
		    os.exit(4);
		}

		// Chown it.
		if (chown) {
		    error = fs.chown(dest, uid, gid);
		    if (error) {
			console.log("File error chown", 
				    JSON.stringify(error, null, 2));
			os.exit(5);
		    }
		}

		// Chmod it.
		if (mode) {
		    error = fs.chmod(dest, mode);
		    if (error) {
			console.log("Chmod file error", JSON.stringify(error, null, 2));
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

		console.log("Success");
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

	scp: function(path, resource, updatedParams, host) {
            var pre = resource._target || {};
            var key = mithras.sshKeyPathForInstance(resource, host);
            var user = mithras.sshUserForInstance(resource, host);

            if (updatedParams.skip == true) {
                log("Skipped.");
            } else if ((updatedParams.ensure != 'absent') &&
                       (updatedParams.ensure != 'present')) {
                console.log(sprintf("Invalid 'file' param property ensure: %s", updatedParams.ensure));
                os.exit(3);
            } else if ((updatedParams.ensure === 'absent')  &&
                       (pre[host.PublicIpAddress] != "found")) {
                log("Ensure: absent; skipping.")
	    } else if (updatedParams.ensure === 'absent') {
                log("Ensure: absent but scp handler does not remove files.")
	    } else if ((updatedParams.ensure === 'present') &&
                       (pre[host.PublicIpAddress] != "found")) {
                var result = mithras.remote.scp(host.PublicIpAddress, 
                                                user, 
                                                key, 
                                                path,
                                                updatedParams.dest);
                
                var out = result[0];
                var err = result[1];
                var ok = result[2];
                var status = result[3];
                if (ok) {
                    if (mithras.verbose) {
                        log(sprintf("SCP to '%s' success.", updatedParams.dest));
                    }
                } else if (status == 255) {
                    log(sprintf("SCP error remote system '%s', dest '%s': %s %s",
                                host.PublicIpAddress, 
                                updatedParams.dest, 
                                err ? err.trim() : "", 
                                out ? out.trim() : ""));
                    os.exit(3);
                } else if (status == 1) {
                    if (mithras.verbose) {
                        log(sprintf("SCP dest '%s' error: %s\n%s", updatedParams.dest, err, out));
                    }
                    os.exit(3);
                } else {
                    if (mithras.verbose) {
                        log(sprintf("SCP dest '%s': status %d; out %s", 
                                    updatedParams.dest, 
                                    status, 
                                    out));
                    }
                }
            } else if (updatedParams.ensure === 'present') {
		log("Ensure: present but destination already exists.")
	    } 
	}

	remote: function(catalog, resources, resource, host) {
	    var updatedParams = resource.params;
            var key = mithras.sshKeyPathForInstance(resource, host);
            var user = mithras.sshUserForInstance(resource, host);
	    var js = sprintf("var run = function() {\n (%s)(%s); };\n", 
			     handler.run.toString(),
			     JSON.stringify(_.omit(updatedParams, 'hosts')));
	    var result = mithras.remote.mithras(host,
						user,
						key,
						js,
						updatedParams.become,
						updatedParams.becomeUser,
						updatedParams.becomeMethod);
	    var out = result[0] ? result[0].trim() : "";
            var err = result[1] ? result[1].trim() : "";
	    if (result[3] == 0) {
		log(sprintf("File '%s' %s: %s.", 
			    updatedParams.dest,
			    updatedParams.ensure,
			    out != "" ? out : "success"));
	    } else {
                log(sprintf("File handler '%s' %s error: %s %s", 
                            updatedParams.dest,
                            updatedParams.ensure,
			    out,
                            err));
		os.exit(3);
	    }
	}

	action: function(catalog, resources, resource, host) {
	    var updatedParams = resource.params;

	    var src = updatedParams.src;
	    var chown = updatedParams.owner;
	    var mode = updatedParams.mode;
	    if (src) {
		var url = web.url.parse(src);
		var scheme = url.scheme;
		if (scheme === "scp") {
		    return function() {
			handler.scp(url.path, resource, updatedParams, host);
			if (chown || mode) {
			    handler.remote(catalog, resources, resource, host);
			}
		    };
		}
	    }
	    return function() {
		handler.remote(catalog, resources, resource, host);
	    };
	}

	handle: function(catalog, resources, resource) {
	    if (resource.module != handler.moduleName) {
		return [null, false];
	    }
	    var params = resource.params;
	    if (Array.isArray(params.hosts)) {
		// Loop over hosts
		_.each(params.hosts, function(host) {

                    // update for by-host variance
                    _.find(resources, function(r) {
			return r.name === resource.name;
                    })._currentHost = host;
                    resource._currentHost = host;
                    var updated = mithras.updateResource(resource, 
							 catalog, 
							 resources,
							 resource.name);
                    if (mithras.verbose) {
			log(sprintf("Host: '%s' (%s)", host.PublicIpAddress, 
                                    host.InstanceId));
                    }
		    if (updated.params.skip == true) {
			if (mithras.verbose) {
			    log("Skipped.");
			}
			return;
		    }
		    var action = handler.action(catalog, resources, updated, host);
		    action();
		});
	    }
	    return [null, true];
	}
    };
    
    handler.init = function () {
        mithras.modules.preflight.register(handler.moduleName, handler.preflight);
	mithras.modules.handlers.register(handler.moduleName, handler.handle);
    };
    
    return handler;
});
