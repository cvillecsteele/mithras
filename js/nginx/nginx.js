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
    'use strict';
    
    /*istanbul ignore next:cant test*/
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    } else if (typeof define === 'function' && define.amd) {
	// AMD. Register as an anonymous module.
	define([], factory);
    } else {
	// Browser globals
	root.objectPath = factory();
    }
})(this, function(){
    var sprintf = require("sprintf.js").sprintf;

    var packageFile = function (path) {
	var parts = filepath.split(module.filename);
	return filepath.join(parts[0], path);
    }
    
    var pkg = {
    	name: "nginxPkg"
    	module: "packager"
    	params: {
    	    name: "nginx"
    	}
    };
    
    var baseDir = "/etc/nginx";
    var confDir = "/etc/nginx/conf.d";
    var availDir = "/etc/nginx/sites-available";
    var enabledDir = "/etc/nginx/sites-enabled";
    var dirs = _.map([confDir, availDir, enabledDir],
		     function(dir, idx) { 
			 return {
    			     name: dir
    			     module: "file"
    			     dependsOn: ["nginxPkg"]
    			     params: {
				 dest: dir
				 ensure: "directory"
				 mode: 0777
    			     }
			 };
		     });
    
    var proxy = {
    	name: "nginxProxyConf"
    	module: "copy"
    	dependsOn: [confDir]
    	params: {
	    dest: filepath.join(confDir, "proxy.conf")
	    content: fs.read(packageFile("proxy.conf"))[0]
	    mode: 0644
    	}
    };
    
    var otherConfig = {
    	name: "nginxOtherConfig"
    	module: "copy"
    	dependsOn: [confDir]
    	params: {
	    mode: 0644
    	}
    };
    
    var enabledLink = {
    	name: "nginxEnabledLink"
    	module: "file"
    	dependsOn: [confDir, availDir, enabledDir]
    	params: {
	    mode: 0644
	    src: "set below"
	    dest: "set below"
	    ensure: "link"
    	}
    };
    
    var siteConfig = {
    	name: "nginxSiteConfig"
    	module: "copy"
    	dependsOn: [availDir]
    	params: {
	    mode: 0644
    	}
    };
    
    var base = {
    	name: "nginxBaseConf"
    	module: "copy"
    	dependsOn: ["nginxPkg"]
    	params: {
	    dest: filepath.join(baseDir, "nginx.conf")
	    overwrite: true
	    content: fs.read(packageFile("nginx.conf"))[0]
	    mode: 0644
    	}
    };
    
    var svc = {
    	name: "nginxSvc"
    	module: "service"
    	dependsOn: [
	    "nginxPkg", "nginxBaseConf", "nginxProxyConf"
	].concat(_.map(dirs,
		       function(dir, idx) {
			   return dir.name;
		       }))
    	params: {
    	    name: "nginx"
    	}
    };
    
    var all = {
	name: "nginx"
	includes: [pkg, svc, base, proxy].concat(dirs)
    }
    
    var nginx = function(template, baseConf, configs, sites) {
	for (var i in all.includes) {
	    r = all.includes[i];
	    r.params = _.extend({}, r.params, template.params);
	}
	if (template.params.ensure === 'present') {
	    for (var i in dirs) {
		r = dirs[i];
		r.params = _.extend({}, r.params, {ensure: "directory"});
	    }
	}

	// handle base configuration
	if (baseConf) {
	    base.params.content = baseConf;
	}

	// For each confg, add a resource to base and make svc depend on it
	for (var key in configs) {
	    if (configs.hasOwnProperty(key)) {
		var name = otherConfig.name + key;
		var r = _.extend({}, otherConfig, 
				 {
				     name: name
				     params: _.extend({}, otherConfig.params, template.params)
				 });
		r.params.dest = filepath.join(confDir, key);
		r.params.content = configs[key];
		all.includes.push(r);
		svc.dependsOn.push(name);
	    }
	}
	
	// For each site, add a resource to base and make svc depend on it
	for (var key in sites) {
	    if (sites.hasOwnProperty(key)) {
		var name = sprintf("%s_%s", siteConfig.name, key);
		var r = _.extend({}, siteConfig, 
				 {
				     name: name
				     params: _.extend({}, siteConfig.params, template.params)
				 });
		r.params.dest = filepath.join(availDir, key);
		r.params.content = sites[key];
		all.includes.push(r);
		svc.dependsOn.push(name);

		// symlink from sites-available to sites-enabled
		var name = sprintf("%s_%s", enabledLink.name, key);
		var r = _.extend({}, enabledLink,
				 {
				     name: name
				     params: _.extend({}, enabledLink.params, template.params)
				 });
		r.params.src = filepath.join(availDir, key);
		r.params.dest = filepath.join(enabledDir, key);
		r.params.ensure = (template.params.ensure === 'present') ? "link" : "absent"
		all.includes.push(r);
		svc.dependsOn.push(name);
	    }
	}

	// Add template dependencies to all included resources
	for (var i in all.includes) {
	    r = all.includes[i];
	    r.dependsOn = (r.dependsOn || []).concat(template.dependsOn || []);
	}

	return all;
    };
    
    return nginx;
});
