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
// # Copy
// 
// Copy is resource handler for creating and copying files.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"copy"`
// 
// Usage:
// 
// `var copy = require("copy").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rFile = {
//   name: "someFile"
//   module: "copy"
//   dependsOn: [otherResource.name]
//   params: {
//     ensure: "present"          // "present" or "absent"
//     become: true               // priv escalation
//     becomeMethod: "sudo"       // escalation method
//     becomeUser: "root"         // desired user
//     dest: "/tmp/foo"           // destination file
//     src: "/etc/hosts"          // source file
//     mode: 0644                 // permissions for destination file
//     hosts: [<array of ec2 instance objects>]
//   }
// };
// ```
// 
// ## Copy Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "present" or "absent"
//
// If `"present"`, the file specified in `dest` will be created,
// either with the contents of the file at `src`, or the file contents
// specified in `content`.  If `"absent"`, the file at `dest` will be
// removed.
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
// ### `dest`
//
// * Required: true
// * Allowed Values: a valid path on the target host
//
// This property specifies the path to the file to be copied into.
//
// ### `src`
//
// * Required: false
// * Allowed Values: a valid path on the target host
//
// One of `src` or `content` must be specified.  If `src` is set, it
// is the path to the file whose contents are to be copied into
// `dest`.
//
// ### `content`
//
// * Required: false
// * Allowed Values: string
//
// One of `src` or `content` must be specified.  If `content` is set, it
// is the path to the file whose contents are to be copied into
// `dest`.
//
// ### `mode`
//
// * Required: true
// * Allowed Values: octal number specifying a valid permission mask
//
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
