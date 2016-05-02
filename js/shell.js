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

//
// # RESOURDCE HANDLER: SHELL
// 
// @public
// 
// # shell
// 
// Shell is a resource handler for dealing with shells running on instances.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"shell"`
// 
// Usage:
// 
// `var shell = require("shell").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var foo = {
//     name: "dostuff"
//     module: "shell"
//     params: {
//       command: "whoami"
//       hosts: [ec2 instance objects...]
//     }
// };
// ```
// 
// ## Parameter Properties
// 
// ### `command`
//
// * Required: true
// * Allowed Values: a valid shell command for the remote instance
//
// The command in this property is run on remote instances.
// 
// ### `hosts`
//
// * Required: false
// * Allowed Values: an array of ec2 instance objects
//
// The command is executed on these instances.
//
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var become = require("become").become;
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleNames: ["shell"]
	handle: function(catalog, resources, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
		
	    var p = resource.params;
	    var ensure = p.ensure;
		    
	    // Sanity
	    if (!p || !p.command) {
		console.log(sprintf("Invalid shell resource params: %s", 
				    JSON.stringify(p)));
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
		    return;
		}
		
		var cmd = become(updatedParams.become, 
				 updatedParams.becomeUser, 
				 updatedParams.becomeMethod, 
				 updatedParams.command);
		var result = mithras.remote.shell(host.PublicIpAddress, 
						  user, 
						  key, 
						  null, 
						  cmd, 
						  null);
		
		var out = result[0];
		var err = result[1];
		var ok = result[2];
		var status = result[3];
		if (ok && status == 0) {
		    target[host.PublicIpAddress] = out;
		    target[host.InstanceId] = out;
		    if (mithras.verbose) {
			log("Success.");
		    }
		} else if (status == 255) {
		    log(sprintf("SSH error remote system '%s', command '%s': %s %s",
				host.PublicIpAddress, 
				updatedParams.command, 
				err.trim(), 
				out.trim()));
		    os.exit(2);
		} else if (status == 1) {
		    if (mithras.verbose) {
			log(sprintf("Shell '%s' error: %s\n%s", updatedParams.command, err, out));
		    }
		    os.exit(3);
		} else {
		    if (mithras.verbose) {
			log(sprintf("Shell '%s': status %d; out %s", 
				    updatedParams.command, 
				    status, 
				    out));
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
