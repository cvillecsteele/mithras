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
	    if (typeof(p.hosts) != "object") {
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
		
		var port = p.port || 22;
		var timeout = p.timeout || 120;
		var ok = network.check(host.PublicIpAddress, port, timeout);
		
		if (ok) {
		    if (mithras.verbose) {
			log(sprintf("Success."));
		    }
		    target[host.PublicIpAddress] = ok;
		    target[host.InstanceId] = ok;
		} else {
		    log(sprintf("Network error remote system '%s', port %d",
				host.PublicIpAddress, 
				port));
		    os.exit(2);
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
