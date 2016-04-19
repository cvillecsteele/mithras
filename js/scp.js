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
// # scp
// 
// SCP is a resource handler for moving files between systems.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"scp"`
// 
// Usage:
// 
// `var scp = require("scp").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var testFile = {
//     name: "test"
//     module: "scp"
//     dependsOn: [other.name]
//     params: {
//         region: defaultRegion
//         ensure: ensure
//         src: "/etc/hosts"
//         dest: "/tmp/foo"
//         hosts: [array of ec2 objects]
//     }
// }
// ```
// 
// ## Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "present" or "absent"
//
// If `"present"` and the file at `dest` does not exist, the file at
// `src`, locally, is copied to the remote system.  If `"absent"`, and
// the file at `dest` on the remore system exists, it is removed.
// 
// ### `src`
//
// * Required: false
// * Allowed Values: a path to a file on the local system.
//
// This file is copied to remote hosts.
//
// ### `dest`
//
// * Required: false
// * Allowed Values: a path to a file on the remote system.
//
// The file from `src` on the local host is copied to this path on
// remote hosts.
//
// ### `hosts`
//
// * Required: false
// * Allowed Values: an array of ec2 instance objects
//
// This property specifies the hosts on which this resource is to be applied.
//
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
        module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
        moduleNames: ["scp"]
        handle: function(catalog, resources, resource) {
            if (!_.find(handler.moduleNames, function(m) { 
                return resource.module === m; 
            })) {
                return [null, false];
            }
                
            var p = resource.params;
            var ensure = p.ensure;

            // Sanity
            if (!p || !p.src || !p.dest) {
                console.log(sprintf("Invalid scp resource params: %s", 
                                    JSON.stringify(p)));
                os.exit(1);
            }
            
            // Loop over hosts
            if (typeof(p.hosts) != "object") {
                return [null, true];
            }

            var pre = resource._target;
            var target = resource._target = {};

            _.each(p.hosts, function(host) {
                if (mithras.verbose) {
                    log(sprintf("Host: '%s' (%s)", host.PublicIpAddress, 
                                host.InstanceId));
                }
                
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

                if (updatedParams.skip == true) {
                    log("Skipped.");
                } else if ((updatedParams.ensure === 'absent')  &&
                           (pre[host.PublicIpAddress] === "found")) {
                    log("Ensure: absent; skipping.")
                } else if ((updatedParams.ensure === 'present') &&
                           (pre[host.PublicIpAddress] != "found")) {
                    var result = mithras.remote.scp(host.PublicIpAddress, 
                                                    user, 
                                                    key, 
                                                    updatedParams.src,
                                                    updatedParams.dest);
                    
                    var out = result[0];
                    var err = result[1];
                    var ok = result[2];
                    var status = result[3];
                    if (ok) {
                        if (mithras.verbose) {
                            log(sprintf("Copy to '%s' success.", updatedParams.dest));
                        }
                    } else if (status == 255) {
                        log(sprintf("SCP error remote system '%s', dest '%s': %s %s",
                                    host.PublicIpAddress, 
                                    updatedParams.dest, 
                                    err.trim(), 
                                    out.trim()));
                        os.exit(2);
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
                }
            });
            return [null, true];
        }
        preflight: function(catalog, resource) {
            if (!_.find(handler.moduleNames, function(m) { 
                return resource.module === m; 
            })) {
                return [null, false];
            }
                
            var p = resource.params;
            var ensure = p.ensure;

            // Sanity
            if (!p || !p.src || !p.dest) {
                console.log(sprintf("Invalid scp resource params: %s", 
                                    JSON.stringify(p)));
                os.exit(1);
            }

            // Any hosts?
            if (typeof(p.hosts) != "object") {
                log("No hosts.");
                return [null, true];
            }

            // Loop over hosts
            var target = resource._target = {};
            _.each(p.hosts, function(host) {
                if (mithras.verbose) {
                    log(sprintf("Host: '%s' (%s)", host.PublicIpAddress, 
                                host.InstanceId));
                }
                
                var key = mithras.sshKeyPathForInstance(resource, host);
                var user = mithras.sshUserForInstance(resource, host);
                
                if (p.skip == true) {
                    log("Skipped.");
                } else {
                    cmd = sprintf("test -f '%s'&& echo 'found'", p.dest);
                    var result = mithras.remote.shell(host.PublicIpAddress, 
                                                      user, 
                                                      key, 
                                                      "",
                                                      cmd, 
                                                      {});
                    if (result[3] == 0) {
                        var out = result[0].trim();
                        log(sprintf("File '%s' %s: %s.", 
                                    p.dest,
                                    p.ensure,
                                    out != "" ? out : "success"));
                        target[host.PublicIpAddress] = out;
                        target[host.InstanceId] = out;
                    }
                }
            });
            return [target, true];
        }
    };
   

    handler.init = function () {
        mithras.modules.preflight.register(handler.moduleNames[0], handler.preflight);
        mithras.modules.handlers.register(handler.moduleNames[0], handler.handle);
        return handler;
    };
    
    return handler;
});
