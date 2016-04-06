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
