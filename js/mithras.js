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
			os.exit(1);
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
			    var result = f(catalog, updated);
			    var target = result[0];
			    var handled = result[1];
			    if (handled && target) {
				var r = _.find(resources, function(r){ 
				    return r.name === name; 
				});
				r.params = updated.params;
				r._target = target;
			    }
			});
		    });
		}
	    }
	}

	// Example template supplied by caller:
	// { 
	//   dependsOn: ["webserver"]
	//   params: {
	// 	   become: true
	// 	   becomeMethod: "sudo"
	// 	   becomUser: "root"
	// 	   hosts: mithras.watch("webserver._target")
	//  }
	// }
	bootstrap: function(template) {
	    var ssh = _.extend({}, template, {
		name: "mithrasSshAvailable"
    		module: "network"
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

	    var paramsWithoutBecome = _.omit(template.params, "become", "becomeUser", "becomeMethod");

	    var dir = _.extend({}, template, {
		name: "mithrasDir"
    		module: "shell"
	    });
	    
	    var mkdirCmd = function(dir) {
		return sprintf("sh -c 'if test -d %s; then echo exists; else mkdir %s && echo created; fi'", dir, dir);
	    }

	    var skipper = function(resourceName) {
		var path = resourceName + "._currentHost.PublicIpAddress";
		return mithras.watch("mithrasBinDir._target",
				     function(catalog, resources, binDirResults) {
					 var ip = objectPath.get(resources, path);
					 if (ip) {
					     return (binDirResults[ip].trim() === "exists");
					 }
				     });
	    };
	    
	    var cmd = mkdirCmd(".mithras");
	    dir.dependsOn = (dir.dependsOn || []).concat([ssh.name]);
	    dir.params = _.extend({}, 
				  paramsWithoutBecome,
				  {
				      command: cmd
				  });

	    var binDir = _.extend({}, template, {
		name: "mithrasBinDir"
    		module: "shell"
	    });
	    var cmd = mkdirCmd(".mithras/bin");
	    binDir.dependsOn = (binDir.dependsOn || []).concat([dir.name]);
	    binDir.params = _.extend({},
				     paramsWithoutBecome,
				     {
					 command: cmd
				     });

	    var scriptsDir = _.extend({}, template, {
		name: "mithrasScriptsDir"
    		module: "shell"
	    });
	    var cmd = mkdirCmd(".mithras/scripts");
	    scriptsDir.dependsOn = (scriptsDir.dependsOn || []).concat([dir.name, binDir.name]);
	    scriptsDir.params = _.extend({}, 
					 paramsWithoutBecome, 
					 {
					     command: cmd
					 });

	    var jsDir = _.extend({}, template, {
		name: "mithrasJsDir"
    		module: "scp"
	    });
	    jsDir.dependsOn = (jsDir.dependsOn || []).concat([binDir.name, dir.name]);
	    jsDir.params = _.extend({}, template.params, {
		skip: skipper("mithrasJsDir")
		src: "js"
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
		    os.exit(1);
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
		    os.exit(1);
		}
		
		// Instance architecture
		var arch = null;
		switch (instance.Architecture) {
		case "x86_64":
		    arch = "amd64";
		    break;
		default:
		    console.log(sprintf("Unknown architecture '%s'", instance.Architecture));
		    os.exit(1);
		}
		
		return [theOS, arch];
	    }

	    var wrapper = _.extend({}, template, {
		name: "mithrasWrapper"
    		module: "scp"
	    });
	    wrapper.dependsOn = (wrapper.dependsOn || []).concat([uname.name, binDir.name]);
	    wrapper.params = _.extend({}, template.params, {
		dest: ".mithras/bin/wrapper"
		skip: skipper("mithrasWrapper")
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
				       var theOS = result[0];
				       var arch = result[1];
				       return sprintf("cache/wrapper_%s_%s", 
						      theOS, 
						      arch);
				   })
	    });

	    var runner = _.extend({}, template, {
		name: "mithrasRunner"
    		module: "scp"
	    });
	    runner.dependsOn = (runner.dependsOn || []).concat([uname.name, binDir.name]);
	    runner.params = _.extend({}, template.params, {
		dest: ".mithras/bin/runner"
		skip: skipper("mithrasRunner")
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
				       return sprintf("cache/runner_%s_%s", 
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
		    var updated = mithras.updateResource(dict[rName], catalog, resources, rName);
		    
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

	// returns a COPY of the resource with updated fields
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
	
	sshKeyPathForInstance: function(resource, instance) {
	    if (resource.params &&
		typeof(resource.params.sshKeyPathForInstance) === 'function') {
                return resource.params.sshKeyPathForInstance(instance);
            } else {
		return "~/.ssh/" + instance.KeyName + ".pem";
	    }
	}

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

	run: function () {
	    var cat = {
		zones: [],
		iamProfiles: [],
		iamRoles: [],
		keypairs: [],
		rrs: [],
		dbs: [],
		caches: [],
		securityGroups: [],
		regions: [],
		instances: [],
		vpcs: [],
		gateways: [],
		subnets: [],
		routeTables: [],
		elbs: [],
	    }
	    
	    if (mithras.verbose) {
		log(sprintf("MITHRAS v %s", mithras.VERSION));
	    }

	    cat.regions = aws.regions.scan();
	    
	    var regions = mithras.activeRegions(cat);
	    var cnt = regions.length;
	    for (var i = 0; i < cnt; i++) {
		if (mithras.verbose) {
		    log(sprintf("Scanning ec2 region: %s", regions[i]));
		}

		// This concat nonsense is because scan functions return array-LIKE things, but we want real arrays.
		cat.caches         = cat.caches.concat(aws.elasticache.scan(regions[i]));
		cat.dbs            = cat.dbs.concat(aws.rds.scan(regions[i]));
		cat.instances      = cat.instances.concat(aws.instances.scan(regions[i]));
		cat.securityGroups = cat.securityGroups.concat(aws.securityGroups.scan(regions[i]));
		cat.vpcs           = cat.vpcs.concat(aws.vpcs.scan(regions[i]));
		cat.gateways       = cat.gateways.concat(aws.vpcs.gateways.scan(regions[i]));
		cat.subnets        = cat.subnets.concat(aws.subnets.scan(regions[i]));
		cat.routeTables    = cat.routeTables.concat(aws.routeTables.scan(regions[i]));
		cat.elbs           = cat.elbs.concat(aws.elbs.scan(regions[i]));
		cat.zones          = cat.zones.concat(aws.route53.zones.scan(regions[i]));
		cat.rrs            = cat.rrs.concat(aws.route53.rrs.scan(regions[i]));
		cat.iamProfiles    = cat.iamProfiles.concat(aws.iam.profiles.scan(regions[i]));
		cat.iamRoles       = cat.iamRoles.concat(aws.iam.roles.scan(regions[i]));
		cat.keypairs       = cat.keypairs.concat(aws.keypairs.scan(regions[i]));
	    }

	    var cnt = cat.instances.length;
	    for (var i = 0; i < cnt; i++) {
		if (cat.instances[i].PublicIpAddress) {
		    peek(cat.instances[i].PublicIpAddress, 
			 mithras.sshKeyPathForInstance({}, cat.instances[i]), 
			 mithras.sshUserForInstance({}, cat.instances[i]),
			 function(output) {
			     cat.instances[i].uname = output;
			 });
		}
	    }

	    return cat;
	}

	// Override if you want to filter regions
	activeRegions: function (catalog) {
	    return _.map(catalog.regions, function(r) {
		return r.RegionName;
	    });
	}

    }); // extend



    // Load handlers
    var copy = require("copy").init();
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
    var scp = require("scp").init();
    var keypairs = require("keypairs").init();

}());
