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
