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
// # mithras
// 
// Mithras is a core object for building scripts to manage your AWS
// stacks.  The `mithras` object is loaded prior to loading your
// scripts.
// 
// The `mithras` object has the following properties:
//
// > * [ARGS](#args)
// > * [HOMe](#home)
// > * [CWD](#cwd)
// > * [GOPATH](#gopath)
// > * [MODULES](#modules)
// > * [VERSION](#version)
// > * [activeRegions](#active)
// > * [apply](#apply)
// > * [become](#become)
// > * [bootstrap](#bootstrap)
// > * [buildDeps](#buildDeps)
// > * [depGraph](#depGraph)
// > * [doIncludes](#doIncludes)
// > * [findGWByVpcId](#findGWByVpcId)
// > * [modules.handlers.register](#modules.handlers.register)
// > * [modules.handlers.run](#modules.handlers.run)
// > * [modules.preflight.register](#modules.preflight.register)
// > * [modules.preflight.run](#modules.preflight.run)
// > * [objectPath](#objectPath)
// > * [remote](#remote)
// > * [resourceMap](#resourceMap)
// > * [run](#run) 
// > * [sshKeyPathForInstance](#sshKeyPathForInstance)
// > * [sshUserForInstance](#sshUserForInstance)
// > * [traverse](#traverse)
// > * [updateResource](#updateResource)
// > * [verbose](#verbose)
// > * [watch](#watch)
// 
// ## Properties
// 
// ### `remote`
//
// See the documentation for the [remote](core_remote.html) core.
// 
// ### `traverse` <a name="traverse"></a>
//
// See the documentation for the [traverse js](https://github.com/substack/js-traverse) module.
// 
// ### `objectPath` <a name="objectPath"></a>
//
// See the documentation for the [object-path](https://github.com/mariocasciaro/object-path) module.
// 
// ### `depGraph` <a name="depGraph"></a>
//
// See the documentation for the [dep-graph.js](https://github.com/TrevorBurnham/dep-graph) module.
// 
// ### `resourceMap(resources) {...}` <a name="resourceMap"></a>
//
// Helper function.  Returns a map of resources by their names.
// 
// ### `become(command, become, becomeUser, becomeMethod) {...}` <a name="become"></a>
//
// Helper function.  Returns a string with the command wrapped in a privilege escalation.
// 
// ### `modules.handlers.register(name, cb) {...}` <a name="modules.handlers.register"></a>
//
// Register a resource handler function.
// 
// A resource handler function takes three arguments:
// 
// > * `catalog`: the current value of the AWS resources found by `mithras.run()`, possibly modified by other resource execution.
// > * `resources`: the set of resources passed into `mithras.apply()`
// > * `targetResource`: the resource object being evaluated for execution.
// 
// Handler functions return a Javascript array.  The second element is
// a boolean. If `true`, it indicates that the handler function "owns"
// the resource and has handled it.  If so, and the first element of
// the return value is defined, the value of the first element is set
// as the `_target` property of the resource object being evaluated.
// 
// See [Design and Concepts](design.html) for a more detailed
// explanation of how this all works.
// 
// ### `modules.handlers.run(catalog, resources, targetResource, dict) {...}` <a name="modules.handlers.run"></a>
//
// Run a resource through handler functions.
//
// In turn, every registered handler function is called with three
// arguments: `catalog`, `resources` and `targetResource`.  The
// catalog is the current value of the AWS resources gathered by
// `mithras.run()`, which may have been modified by previous resource
// execution.  The `resources` argument is the set of resources that
// was passed into `mithras.apply()`.  The final argument is the
// resource object currently being evaluated.
//
// Handler functions return a Javascript array.  The second element is
// a boolean. If `true`, it indicates that the handler function "owns"
// the resource and has handled it.  If so, and the first element of
// the return value is defined, the value of the first element is set
// as the `_target` property of the resource object being evaluated.
// 
// See [Design and Concepts](design.html) for a more detailed
// explanation of how this all works.
// 
// ### `modules.preflight.register(name, cb) {...}` <a name="modules.preflight.register"></a>
//
// Register a resource handler preflight function.
// 
// A resource handler preflight function takes three arguments:
// 
// > * `catalog`: the current value of the AWS resources found by `mithras.run()`, possibly modified by other resource execution.
// > * `resources`: the set of resources passed into `mithras.apply()`
// > * `targetResource`: the resource object being evaluated for preflight.
// 
// Preflight functions return a Javascript array.  The second element
// is a boolean. If `true`, it indicates that the function "owns" the
// resource and has handled it.  If so, and the first element of the
// return value is defined, the value of the first element is set as
// the `_target` property of the resource object being evaluated.
// 
// See [Design and Concepts](design.html) for a more detailed
// explanation of how this all works.
// 
// ### `modules.preflight.run(catalog, resources, order) {...}` <a name="modules.preflight.run"></a>
//
// Run preflight functions on resources.
// 
// For each resource, first the resource is recursively traversed by
// [`updateResource`](#updateResource), and given an opportunity to
// access runtime values in other resources.
// 
// Next, every registered preflight function is called with three
// arguments: `catalog`, `resources` and `targetResource`.  The
// catalog is the current value of the AWS resources gathered by
// `mithras.run()`, which may have been modified by previous resource
// execution.  The `resources` argument is the set of resources that
// was passed into `mithras.apply()`.  The final argument is the
// resource object currently being evaluated.
//
// Handler functions return a Javascript array.  The second element is
// a boolean. If `true`, it indicates that the handler function "owns"
// the resource and has succesfully preflighted it.  If so, and the
// first element of the return value is defined, the value of the
// first element is set as the `_target` property of the resource
// object being evaluated.
// 
// See [Design and Concepts](design.html) for a more detailed
// explanation of how this all works.
// 
// ### `MODULES` <a name="modules"></a>
//
// A map of loaded core module names to their version strings.
// 
// ### `GOPATH` <a name="gopath"></a>
//
// The current GOPATH when mithras is invoked (if set).
// 
// ### `HOME` <a name="home"></a>
//
// The value of $MITHRASHOME or -m if set on the command line.
// 
// ### `CWD` <a name="cwd"></a>
//
// The current working directory of Mithras.
// 
// ### `ARGS` <a name="args"></a>
//
// When `mithras run` is invoked at the command line, any additional
// non-flag parameters supplied on the command line are passed through
// to the user script in this array.
// 
// ### `verbose` <a name="verbose"></a>
//
// Set to `true` if the `-v` global flag is used to invoke Mithras on
// the command line.  Eg., `mithras -v run ...`
// 
(function(){

    var objectPath = require("object_path.js");
    var sprintf = require("sprintf.js").sprintf;

    var become = function(command, become, becomeUser, becomeMethod) {
        if (become) {
            if (becomeUser) {
                return sprintf("%s -u %s %s",
                               becomeMethod || "sudo",
                               becomeUser,
                               command);
            } else {
                return sprintf("%s -u %s %s",
                               becomeMethod || "sudo",
                               "root",
                               command);
            }
        } else {
            return command;
        }
    }

    var resourceMap = function(resourceList) {
        return _.reduce(resourceList, function(memo, r){ 
            memo[r.name] = r;
            return memo;
        }, {});
    }

    _.extend(mithras, {
        traverse: require("traverse.js")
        objectPath: require("object_path.js")
        depGraph: require("dep-graph").DepGraph
        resourceMap: resourceMap
        become: become

        modules: {
            handlers: {
                funcs: {}
                register: function(name, cb) {
                    mithras.modules.handlers.funcs[name] = cb;
                }
                run: function(catalog, resources, targetResource, dict) {
                    var name = targetResource.name;
                    var handled = false;
                    _.each(mithras.modules.handlers.funcs, function(f, fname) {
                        var result = f(catalog, resources, targetResource);
                        var target = result[0];
                        if (!handled) {
                            handled = result[1];
                            if (handled && target) {
                                var r = _.find(resources, function(r){ 
                                    return r.name === name; 
                                });
                                r._target = target;
                            }
                        }
                    });
                    if (!handled && targetResource.module) {
                        log0(sprintf("Resource '%s': module '%s' not handled.", 
                                     targetResource.name,
                                     targetResource.module));
                        os.exit(3);
                    }
                }
            }
            preflight: {
                funcs: []
                register: function(name, cb) {
                    mithras.modules.preflight.funcs.push(cb);
                }
                run: function(catalog, resources, order) {
                    var dict = resourceMap(resources);
                    _.each(order, function(name) {
                        log0(sprintf("PREFLIGHT: %s", name));

                        var updated = mithras.updateResource(dict[name], 
                                                             catalog, 
                                                             resources, 
                                                             name);
                        dict[name] = updated;

                        _.each(mithras.modules.preflight.funcs, function(f) {
                            var result = f(catalog, resources, updated);
                            var target = result[0];
                            var handled = result[1];
                            if (handled && target) {
                                var r = _.find(resources, function(r){ 
                                    return r.name === name; 
                                });
                                r.params = updated.params;
				log.debug("  --- Setting _target")
                                r._target = target;
                            }
                        });
                    });
                }
            }
        }

        // @public
        // <a name="bootstrap"></a>
        // 
        // ### `bootstrap(template) {...}`
        //
        // Returns a single resource object, which encapsulates all of
        // the resources need to bootstrap an instance for use with
        // Mithras.
        //
        // Example template supplied by caller:
        // ```
        // { 
        //   dependsOn: ["webserver"]
        //   params: {
        //         become: true
        //         becomeMethod: "sudo"
        //         becomUser: "root"
        //         hosts: mithras.watch("webserver._target")
        //  }
        // }
        // ```
        bootstrap: function(template) {
	    var ensure = template.params.ensure
            var ssh = _.extend({}, template, {
                name: "mithrasSshAvailable"
                module: "network"
                skip: (ensure === "absent")
            });
            ssh.params = _.extend({}, template.params, {
                timeout: 120
            });
            
            var uname = _.extend({}, template, {
                name: "mithrasUname"
                module: "shell"
            });
            uname.dependsOn = (uname.dependsOn || []).concat([ssh.name]);
            uname.params = _.extend({}, template.params, {
                command: "uname -s"
            });

            var paramsWithoutBecome = _.omit(template.params, 
					     "become", "becomeUser", "becomeMethod");

            var dir = _.extend({}, template, {
                name: "mithrasDir"
                module: "shell"
            });
            
            var mkdirCmd = function(dir) {
                return sprintf("sh -c 'if test -d %s; then echo exists; else mkdir %s && echo created; fi'", dir, dir);
            }

            var skipper = function(resourceName) {
                var path = resourceName + "._currentHost.PublicIpAddress";
                return mithras.watch(resourceName + "._target",
                                     function(catalog, resources, results) {
                                         var ip = objectPath.get(resources, path);
                                         if (ip) {
					     var found = results[ip] ? results[ip].trim() : "";
                                             return (found === "found");
                                         }
                                     });
            };

            // Mithras dir
            var cmd = mkdirCmd(".mithras");
            dir.dependsOn = (dir.dependsOn || []).concat([ssh.name]);
            dir.params = _.extend({}, 
                                  paramsWithoutBecome,
                                  {
				      skip: skipper("mithrasDir")
                                      command: cmd
                                  });

	    // Binaries
            var binDir = _.extend({}, template, {
                name: "mithrasBinDir"
                module: "shell"
            });
            var cmd = mkdirCmd(".mithras/bin");
            binDir.dependsOn = (binDir.dependsOn || []).concat([dir.name]);
            binDir.params = _.extend({},
                                     paramsWithoutBecome,
                                     {
					 skip: (ensure === "absent") ? ensure : skipper(binDir.name)
                                         command: cmd
                                     });

	    // Scripts
            var scriptsDir = _.extend({}, template, {
                name: "mithrasScriptsDir"
                module: "shell"
            });
            var cmd = mkdirCmd(".mithras/scripts");
            scriptsDir.dependsOn = (scriptsDir.dependsOn || []).concat([dir.name, binDir.name]);
            scriptsDir.params = _.extend({}, 
                                         paramsWithoutBecome, 
                                         {
					     skip: (ensure === "absent") ? true : skipper(scriptsDir.name)
                                             command: cmd
                                         });

            var jsDir = _.extend({}, template, {
                name: "mithrasJsDir"
                module: "file"
            });
            jsDir.dependsOn = (jsDir.dependsOn || []).concat([binDir.name, dir.name]);
            jsDir.params = _.extend({}, template.params, {
                skip: (ensure === "absent") ? true : skipper(jsDir.name)
                src: "scp://localhost"+mithras.JSDIR
                dest: ".mithras/js"
            });

            var osAndArchForInstance = function(catalog, resources, instance) {
                if (!instance) {
                    return;
                }
                var ip = instance.PublicIpAddress;
                var u = objectPath.get(resources, "mithrasUname._target");
                if (!u || !u[ip] || typeof(u[ip]) != "string") {
                    console.log(sprintf("No uname for '%s'", ip));
                    os.exit(3);
                }

                // Instance os
                var theOS = null;
                switch (u[ip].trim()) {
                case "Linux":
                    theOS = "linux";
                    break;
                case "Darwin":
                    theOS = "darwin";
                    break;
                default:
                    console.log(sprintf("Unknown os '%s'", u[ip]));
                    os.exit(3);
                }
                
                // Instance architecture
                var arch = null;
                switch (instance.Architecture) {
                case "x86_64":
                    arch = "amd64";
                    break;
                default:
                    console.log(sprintf("Unknown architecture '%s'", instance.Architecture));
                    os.exit(3);
                }
                
                return [theOS, arch];
            }

	    // Wrapper binary
            var wrapper = _.extend({}, template, {
                name: "mithrasWrapper"
                module: "file"
            });
            wrapper.dependsOn = (wrapper.dependsOn || []).concat([uname.name, binDir.name]);
            wrapper.params = _.extend({}, template.params, {
                dest: ".mithras/bin/wrapper"
                skip: (ensure === "absent") ? true : skipper(wrapper.name)
                src: mithras.watch("mithrasWrapper._currentHost", 
                                   function(catalog, resources, inst) {
                                       var u = objectPath.get(resources, 
							      "mithrasUname._target");
                                       var ip = inst.PublicIpAddress;
                                       if (!u || !u[ip] || typeof(u[ip]) != "string") {
                                           return;
                                       }
                                       var result = osAndArchForInstance(catalog, 
									 resources, 
									 inst);
                                       if (!result) {
                                           return;
                                       }
                                       var theOS = result[0];
                                       var arch = result[1];
                                       return sprintf("scp://localhost/%s/cache/wrapper_%s_%s", 
						      mithras.HOME,
                                                      theOS, 
                                                      arch);
                                   })
            });

	    // Runner binary
            var runner = _.extend({}, template, {
                name: "mithrasRunner"
                module: "file"
            });
            runner.dependsOn = (runner.dependsOn || []).concat([uname.name, binDir.name, wrapper.name]);
            runner.params = _.extend({}, template.params, {
                dest: ".mithras/bin/runner"
                skip: (ensure === "absent") ? true : skipper(runner.name)
                src: mithras.watch("mithrasWrapper._currentHost", 
                                   function(catalog, resources, inst) {
                                       var u = objectPath.get(resources, "mithrasUname._target");
                                       var ip = inst.PublicIpAddress;
                                       if (!u || !u[ip] || typeof(u[ip]) != "string") {
                                           return;
                                       }
                                       var result = osAndArchForInstance(catalog, resources, inst);
                                       if (!result) {
                                           return;
                                       }
                                       var os = result[0];
                                       var arch = result[1];
                                       return sprintf("scp://localhost/%s/cache/runner_%s_%s", 
						      mithras.HOME,
                                                      os, 
                                                      arch);
                                   })
            });

            // Give back the base resource, which includes all the dependencies.
            _.extend(this, template, {
                includes: [dir, binDir, jsDir, scriptsDir, uname, wrapper, runner, ssh]
                dependsOn: (template.dependsOn || []).concat([wrapper.name, 
                                                              runner.name, 
                                                              scriptsDir.name,
                                                              jsDir.name])
            });
        }

        // @public
        // <a name="apply"></a>
	// 
        // ### `apply(catalog, resources, reverse) {...}`
        //
        // The "core" function of Mithras.  Given a `catalog`, an
        // array of resource objects in `resources`, and a boolean
        // (`reverse`), apply the resources to the catalog.
        //
        // First a dependency graph is built.  In forward order, all
        // resources are preflighted.
        //
        // Next, in the desired order (reversed, if `reverse` is
        // `true`), the resources are run through their handlers in
        // dependency order.
        //
        // The catalog, after update by handlers, is returned.
        //
	// See [Design and Concepts](design.html) for a more detailed
	// explanation of how this all works.
	// 
        apply: function(catalog, resources, reverse) {
            // include sub-resources
            var resources = mithras.doIncludes(resources);
            
            // build dep graph
            var fwdDeps = new mithras.depGraph();
            _.each(resources, function(r) {
                fwdDeps.addNode(r.name);
            });
            _.each(resources, function(r) {
                if (r.hasOwnProperty("dependsOn")) {
                    if (typeof(r.dependsOn) === "string") {
                        fwdDeps.addDependency(r.name, r.dependsOn);
                    } else if (Array.isArray(r.dependsOn)) {
                        _.each(r.dependsOn, function(d) {
                            fwdDeps.addDependency(r.name, d);
                        });
                    }
                }
            });
            var revDeps = new mithras.depGraph();
            _.each(resources, function(r) {
                revDeps.addNode(r.name);
            });
            _.each(resources, function(r) {
                if (r.hasOwnProperty("dependsOn")) {
                    if (typeof(r.dependsOn) === "string") {
                        revDeps.addDependency(r.dependsOn, r.name);
                    } else if (Array.isArray(r.dependsOn)) {
                        _.each(r.dependsOn, function(d) {
                            revDeps.addDependency(d, r.name);
                        });
                    }
                }
            });
            
            // Get forward and reverse order
            var fwdOrder = fwdDeps.overallOrder();
            var revOrder = revDeps.overallOrder();
            
            // Preflight in fwd deps order
            mithras.modules.preflight.run(catalog, resources, fwdOrder);
            
            // Build map of resource name to resource
            var dict = resourceMap(resources);

            // Call handlers in specified order
            var order = fwdOrder;
            if (reverse) {
                order = revOrder;
            }
            _.each(order, function(rName) {
                if (dict[rName].skip) {
                    if (mithras.verbose) {
                        log0(sprintf("SKIPPING: %s", rName));
                    }
                } else {
                    log0(sprintf("RESOURCE: %s", rName));

                    // Update the resource
                    var updated = mithras.updateResource(dict[rName], 
							 catalog, 
							 resources, 
							 rName);

                    // Run handlers on updated resource
                    mithras.modules.handlers.run(catalog, resources, updated, dict);
                    
                    // Sleeeeep (possibly)
                    if (updated.delay && updated.delay > 0) {
                        if (mithras.verbose) {
                            log(sprintf("Delay %d seconds", updated.delay));
                        }
                        time.sleep(updated.delay);
                    }
                }
            });

            return catalog;
        }

        // @public
        // 
        // ### `buildDeps(resources, add_node, add_dep) {...}`
        //
        // Helper function.  Loop through resources, calling the
        // `add_node` and `add_dep` functions supplied by the caller
        // to create adependency graph.
        //
        buildDeps: function(resources, add_node, add_dep) {
            for (var key in resources) {
                add_node(resources[key].name, resources[key]);
            }
            for (var key in resources) {
                if (resources[key].hasOwnProperty("dependsOn")) {
                    if (typeof(resources[key].dependsOn) === "string") {
                        add_dep(resources[key].name, resources[key].dependsOn);
                    } else if (Array.isArray(resources[key].dependsOn)) {
                        for (var k in resources[key].dependsOn) {
                            if (resources[key].dependsOn.hasOwnProperty(k)) {
                                add_dep(resources[key].name, resources[key].dependsOn[k]);
                            }
                        }
                    }
                }
            }
        }

        // @public
        // 
        // ### `updateResource(resource, catalog, resources, name) {...}`
        //
        // Given a `resource`, the `catalog` and a list of all
        // `resources`, update the resource allowing it to reach into
        // the catalog and/or other resources to set its properties.
        //
        // Returns a COPY of the resource with updated fields.
        //
        updateResource: function(resource, catalog, resources, name) {
            var dict = resourceMap(resources);
            var re1 = new RegExp("on_");
            var updated = mithras.traverse(_.omit(resource, 
                                                  "includes")).map(function (f) {
                if (!re1.exec(this.key) && (typeof(f) == "function")) {
                    val = f(catalog, dict);
                    if (typeof(val) != "undefined") {
                        this.update(val);
                    }
                }
            });
            updated.includes = resource.includes;
            return updated;
        }

        // @public
        // <a name="doIncludes"></a>
        // 
        // ### `doIncludes(resources) {...}`
        //
        // Recursively descend through supplied `resources`, adding
        // their dependencies via their `includes` property.
        //
        doIncludes: function(resources) {
            mithras.traverse(resources).map(function (f) {
                if (this.key === "includes") {
                    for (var idx in f) {
                        resources.push(f[idx]);
                    }
                }
            });
            return resources;
        }
        
        // @public
        // <a name="sshKeyPathForInstance"></a>
        // 
        // ### `sshKeyPathForInstance(resource, instance) {}`
        //
        // Given a `resource` and an ec2 `instance` object, return the
        // appropriate path to the SSH key for the instance.
        //
        // If the resource has a property named
        // `sshKeyPathForInstance`, it is invoked and its return value
        // used.
        //
        // The default return value is:
        //
        // `"~/.ssh/" + instance.KeyName + ".pem"`
        //
        sshKeyPathForInstance: function(resource, instance) {
            if (resource.params &&
                typeof(resource.params.sshKeyPathForInstance) === 'function') {
                return resource.params.sshKeyPathForInstance(instance);
            } else if (instance && instance.KeyName) {
                return "~/.ssh/" + instance.KeyName + ".pem";
            }
        }

        // @public
        // <a name="sshUserForInstance"></a>
        // 
        // ### `sshUserForInstance(resource, instance) {}`
        //
        // Given a `resource` and an ec2 `instance` object, return the
        // appropriate SSH username for the instance.
        //
        // If the resource has a property named
        // `sshUserForInstance`, it is invoked and its return value
        // used.
        //
        // The default return value is:
        //
        // `"ec2-user"`
        //
        sshUserForInstance: function(resource, instance) {
            if (resource.params &&
                typeof(resource.params.sshUserForInstance) === 'function') {
                return resource.params.sshUserForInstance(instance);
            } else {
                return "ec2-user";
            }
        }

        paramFunction: function(f) {
            return function() { return f; };
        }

        // @public
        // <a name="watch"></a>
        // 
        // ### `watch(path, cb) {...}`
        //
        // A resource may set any property value to a function,
        // instead of a string, array, etc.  When the resource is
        // preflight'ed, that function will be called with two
        // arguments, the current `catalog`, and an array of
        // `resources`.  The parameter function may return `undefined`,
        // and if so, it will remain a function.  If it returns any
        // other value, the value of the property to which it is
        // attached is retplaced with the parameter function's return
        // value.
        //
        // This allows parameter properties to be evaluated at
        // runtime, not just when the resource is defined.  The use
        // case is appropriate to AWS, when a given resource needs the
        // value from some other resource in order to be handled.  For
        // example, instances may be placed into subnets, which are
        // defined on some other resource.  This gives the target
        // resource the ability to reach into resources it depends on
        // and extract values for use in parameters.
        //
        // All of this functionality is wrapped up in a neat little
        // package here, in `mithras.wrapper`.  Using this function,
        // one resource can examine another for its runtime
        // properties.  Here's an example:
        //
        // ```
        // var rSubnetA = {
        //      name: "subnetA"
        //      module: "subnet"
        //      dependsOn: [rVpc.name]
        //      params: {
        //          region: defaultRegion
        //          ensure: ensure
        //          subnet: {
        //              CidrBlock:        "172.33.1.0/24"
        //              VpcId:            mithras.watch("VPC._target.VpcId")
        //              AvailabilityZone: defaultZone
        //          }
        //          tags: {
        //              Name: "primary-subnet"
        //          }
        //          routes: [
        //              {
        //                  DestinationCidrBlock: "0.0.0.0/0"
        //                  GatewayId:            mithras.watch("VPC._target.VpcId", mithras.findGWByVpcId)
        //              }
        //          ]
        //      }
        // };
        // ```
        //
        // This example resource defines a subnet.  Subnets are
        // created in VPCs.  Here, the subnet property `VpcId` uses
        // `mithras.watch` at runtime to obtain the id of the VPC on
        // which it depends, since the id is not known when the
        // resource is defined.
        //
        // Supply two parameters to `watch`: a path (expressed in
        // [object-path](https://github.com/mariocasciaro/object-path)
        // form), and an optional callback.  If the object path is
        // defined, the callback will be invoked `cb(catalog,
        // resources, value_at_path)` The return value of the callback
        // will be returned by the watch function during preflight,
        // and follows the same rules outlined above.
        //
	// See [Design and Concepts](design.html) for a more detailed
	// explanation of how the `_target` property gets set on
	// resource objects.
	// 
        watch: function(path, cb) {
            return function(cat, resources) {
                var p = objectPath.get(resources, path);
                if (typeof(p) != "undefined") {
                    if (typeof(cb) === "function") {
                        return cb(cat, resources, p);
                    }
                    return p;
                }
            }
        }

        // @public
        // <a name="findGWByVpcId"></a>
        // 
        // ### `findGWByVpcId(cat, resources, vpcId) {...}`
        //
        // Given a `vpcId`, look through the `catalog` and find a
        // matching internet gateway.  If one is found, return its
        // `InternetGatewayId` property.
        //
        findGWByVpcId: function (cat, resources, vpcId) {
            var gw = _.find(cat.gateways, function(gw){ 
                for (var i in gw.Attachments) {
                    if (gw.Attachments[i].VpcId === vpcId) {
                        return true;
                    }
                }
                return false;
            });
            if (gw) {
                return gw.InternetGatewayId;
            }
        }

        // @public
        // <a name="run"></a>
        // 
        // ### `run() {...}`
        //
        // Called by user scripts to interrogate AWS and return a
        // `catalog` of resources.
        //
        run: function (targets) {
            if (mithras.verbose) {
                log0(sprintf("--- MITHRAS v %s --- ###", mithras.VERSION));
            }

	    var scanners = {
		autoscalingGroups: aws.autoscaling.groups.scan,
		autoscalingLaunchConfigs: aws.autoscaling.launchConfigs.scan,
		autoscalingHooks: aws.autoscaling.hooks.scan,
		caches: aws.elasticache.scan,
		dbs: aws.rds.scan,
		instances: aws.instances.scan,
		securityGroups: aws.securityGroups.scan,
		vpcs: aws.vpcs.scan,
		gateways: aws.vpcs.gateways.scan,
		subnets: aws.subnets.scan,
		routeTables: aws.routeTables.scan,
		elbs: aws.elbs.scan,
		zones: aws.route53.zones.scan,
		rrs: aws.route53.rrs.scan,
		iamProfiles: aws.iam.profiles.scan,
		iamRoles: aws.iam.roles.scan,
		keypairs: aws.keypairs.scan,
		subs: aws.sns.subs.scan,
		topics: aws.sns.topics.scan,
		queues: aws.sqs.scan,
	    };

	    if (!targets) {
		targets = Object.keys(scanners);
	    }
	    var cat = _.reduce(targets, function(memo, t){ 
		memo[t] = [];
		return memo;
            }, {});

	    targets = _.reduce(targets, function(memo, t){ 
		memo[t] = scanners[t];
		return memo;
            }, {});

	    var regions = aws.regions.scan();
	    cat.regions = regions;
	    _.each(mithras.activeRegions(cat), function(region) {
                if (mithras.verbose) {
                    log(sprintf("Scanning ec2 region: %s", region));
                }
                // This concat nonsense is because scan functions return array-LIKE things, but we want real arrays.
		_.each(targets, function(f, target) {
		    cat[target] = cat[target].concat(f(region));
		});
	    });

            return cat;
        }

        // @public
        // <a name="active"></a>
        // ### `activeRegions(catalog) {...}`
        //
        // Returns an array of AWS regions.  User scripts may replace
        // this function with their own to limit the scope of queries
        // that `mithras.run()` will execute in looking for resources
        // on AWS.
        //
        // If not replaced, `mithras.activeRegions` will return all
        // regions.
        //
        activeRegions: function (catalog) {
            return _.map(catalog.regions, function(r) {
                return r.RegionName;
            });
        }

    }); // extend



    // Load handlers

    var elasticache = require("elasticache").init();
    var elb = require("elb").init();
    var file = require("file").init();
    var git = require("git").init();
    var instance = require("instance").init();
    var packager = require("packager").init();
    var rds = require("rds").init();
    var route53 = require("route53").init();
    var s3 = require("s3").init();
    var service = require("service").init();
    var sg = require("secgroup").init();
    var shell = require("shell").init();
    var subnet = require("subnet").init();
    var vpc = require("vpc").init();
    var iam = require("iam").init();
    var network = require("network").init();
    var keypairs = require("keypairs").init();
    var sns = require("sns").init();
    var sqs = require("sqs").init();
    var autoscaling = require("autoscaling").init();

}());
