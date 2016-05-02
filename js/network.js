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
// # network
// 
// Network is resource handler for working with network connections.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"network"`
// 
// Usage:
// 
// `var network = require("network").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var available = {
//   name: "sshAvailable"
//   module: "network"
//   dependsOn: [otherResource.name]
//   params: {
//     timeout: 120
//     port: 22
//     hosts: [<array of ec2 instance objects>]
//   }
// };
// ```
// 
// ## Parameter Properties
// 
// ### `port`
//
// * Required: false
// * Allowed Values: integer
//
// The port on the remote host to attempt to connect to.  Defaults to
// 22.
// 
// ### `timeout`
//
// * Required: false
// * Allowed Values: integer
//
// A number of seconds to attempt to connect to the remote host.
// Defaults to 120.
// 
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleNames: ["network"]
	handle: function(catalog, resources, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
	    
	    var p = resource.params;
	    
	    // Sanity
	    if (!p) {
		console.log("Invalid network resource");
		os.exit(1);
	    }
	    
	    // Loop over hosts
	    if (!Array.isArray(p.hosts)) {
		return [null, true];
	    }
	    var target = resource._target = {};
	    _.each(p.hosts, function(host) {
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
		p = updated.params;
		var ensure = p.ensure;
		
		var port = p.port || 22;
		var timeout = p.timeout || 120;
		var ok = network.check(host.PublicIpAddress, port, timeout);
		
		if (ok) {
		    if (ensure === "present") {
			if (mithras.verbose) {
			    log(sprintf("Success."));
			}
			target[host.PublicIpAddress] = ok;
			target[host.InstanceId] = ok;
		    } else if (ensure === "absent") {
			log("Error: network connection still alive.");
			os.exit(1);
		    }
		} else {
		    if (ensure === "present") {
			log(sprintf("Network error remote system '%s', port %d",
				    host.PublicIpAddress, 
				    port));
			os.exit(2);
		    } else if (ensure === "absent") {
			if (mithras.verbose) {
			    log(sprintf("Success."));
			}
			target[host.PublicIpAddress] = ok;
			target[host.InstanceId] = ok;
		    }
		}

	    });
	    return [target, true];
	}
    };
		   

    handler.init = function () {
	mithras.modules.handlers.register(handler.moduleNames[0], handler.handle);
	return handler;
    };
    
    return handler;
});
