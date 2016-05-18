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
// # packager
// 
// Packager is resource handler for working with packager connections.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"packager"`
// 
// Usage:
// 
// `var packager = require("packager").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rUpdatePkgs = {
//     name: "updatePackages"
//     module: "packager"
//     dependsOn: [otherResource.name]
//     params: {
//         ensure: "latest"
//         name: ""
//         become: true
//         becomeMethod: "sudo"
//         becomeUser: "root"
//         hosts: [ec2 instances objects]
//     }
// };
// ```
// 
// ## Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "present", "absent", "latest"
//
// If `"present"`, the package will be installed if is not already.
// If `"absent"`, the package will be removed if it is present.  If
// `"latest"`, the package will be installed if not present, and
// updated if it is already on the remote system.
// 
// ### `name`
//
// * Required: true
// * Allowed Values: any valid package name; eg "nginx"
//
// If set to the empty string, `""`, this resource handler will omit
// the package name from the command line invocation of `yum`.  Use
// `""` in conjunction with `ensure: "latest"` to upgrade all packages
// on the remote system.
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
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
        module.exports = factory();
    }
})(this, function() {
    
    var become = require("become").become;
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
        moduleNames: ["packager"]
        checkInstalled: function(resource, instance, user, key) {
            var p = resource.params;
            var cmd = sprintf("yum list installed %s", p.name);
            cmd = become(p.become, p.becomeUser, p.becomeMethod, cmd);
            cmd = cmd.split(/[ ]+/);
            result = mithras.remote.wrapper(instance, user, key, cmd, null);
            var out = result[0];
            var err = result[1];
            var ok = result[2];
            var status = result[3];
            var version = "";
            if (ok && status == 0) {
                var lines = out.split("\n");
                for (var i = 0; i < lines.length; i++) {
                    var parts = lines[i].split(/[ \t]/);
                    for (var j = 0; j < parts.length; j++) {
                        if (parts[j].match(/^[0-9].*/)) {
                            version = parts[j];
                        }
                    }
                }
                if (mithras.verbose && version != "") {
                    log(sprintf("Package '%s' found version '%s'",
                                p.name,
                                version));
                    return true;
                }
            } else if (status == 255) {
                log(sprintf("Error communiating with remote system, package '%s': %s",
                            p.name, err.trim()));
                os.exit(3);
            } else if (status == 1 && mithras.verbose) {
                log(sprintf("Package '%s' not found", p.name));
            } else if (mithras.verbose) {
                    log(sprintf("Package '%s': status %d; out %s", 
                                p.name, status, out));
            }
            return false;
        }

        install: function(resource, inst, user, key) {
            var p = resource.params;
            var cmd = "";
            switch (p.ensure) {
            case "present":
                cmd = sprintf("yum -y install %s", p.name);
                break;
            case "absent":
                cmd = sprintf("yum -y remove %s", p.name);
                break;
            case "latest":
                cmd = sprintf("yum -y update %s", p.name);
                break;
            }
            cmd = become(p.become, p.becomeUser, p.becomeMethod, cmd).trim().split(/[ \t]+/);
            var result = mithras.remote.wrapper(inst, user, key, cmd, null);
            var out = result[0].trim();
            var err = result[1].trim();
            var ok = result[2];
            var status = result[3];
            if (ok && status == 0) {
                if (mithras.verbose) {
                    log(sprintf("Package '%s' %s", p.name ? p.name : "ALL", p.ensure));
                }
                return true;
            } else if (status == 255) {
                log(sprintf("Error communiating with remote system, package '%s': %s",
                            p.name, err.trim()));
                os.exit(3);
            } else if (status == 1) {
                if (mithras.verbose) {
                    log(sprintf("Package '%s' error: %s", p.name, err));
                }
                os.exit(3);
            } else {
                if (mithras.verbose) {
                    log(sprintf("Package '%s': status %d; out %s", 
                                p.name, 
                                status, 
                                out));
                }
            }
            return false;
        }

        handle: function(catalog, resources, resource) {
            if (!_.find(handler.moduleNames, function(m) { 
                return resource.module === m; 
            })) {
                return [null, false];
            }
                
            var params = resource.params;

            // Sanity
            if (!params || !(typeof(params.name) === 'string')) {
                console.log("Invalid packager params", JSON.stringify(params, null, 2));
                os.exit(3);
            }
            
            // Loop over hosts
            if (!Array.isArray(params.hosts)) {
                return [null, true];
            }
            _.each(params.hosts, function(host) {
                var key = mithras.sshKeyPathForInstance(resource, host);
                var user = mithras.sshUserForInstance(resource, host);
                if (mithras.verbose) {
                    log(sprintf("Host: '%s' (%s)", host.PublicIpAddress, 
                                host.InstanceId));
                }

                // update for by-host variance
                _.find(resources, function(r) {
                    return r.name === resource.name;
                })._currentHost = host;
                var updated = mithras.updateResource(resource, 
                                                     catalog, 
                                                     resources, 
                                                     resource.name);
                params = updated.params;
                var ensure = params.ensure;
                
                switch(ensure) {
                case "absent":
                    if (handler.checkInstalled(updated, host, user, key)) {
                        handler.install(updated, host, user, key);
                    }
                    break;
                case "present":
                    if (!handler.checkInstalled(updated, host, user, key)) {
                        handler.install(updated, host, user, key);
                    }
                    break;
                case "latest":
                    handler.install(updated, host, user, key);
                    break;
                }
            });
            
            return [null, true];
        }
    }

    handler.init = function () {
        _.each(handler.moduleNames, function(name) {
            mithras.modules.handlers.register(name, handler.handle);
        });
        return handler;
    };
    
    return handler;
});
