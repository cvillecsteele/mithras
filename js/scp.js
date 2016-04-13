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
		} else if (updatedParams.ensure === 'absent') {
		    log("Ensure: absent; skipping.")
		} else if (updatedParams.ensure === 'present') {
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
    };
		   

    handler.init = function () {
	mithras.modules.handlers.register(handler.moduleNames[0], handler.handle);
	return handler;
    };
    
    return handler;
});
