(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var become = require("become").become;
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleNames: ["service"]
	running: function(resource, instance, user, key) {
	    var p = resource.params;
	    var cmd = sprintf("service %s status", p.name)
	    cmd = become(p.become, p.becomeUser, p.becomeMethod, cmd);
	    result = mithras.remote.shell(instance.PublicIpAddress, 
					  user, 
					  key, 
					  null,
					  cmd,
					  null);
	    var out = result[0].trim();
	    var err = result[1].trim();
	    var ok = result[2];
	    var status = result[3];
	    if (ok && status == 0) {
		if (mithras.verbose && out != "") {
		    log(sprintf("Service '%s' status '%s'",
				p.name,
				out));
		    return [null, true];
		}
	    } else if (status == 255) {
		log(sprintf("Error communicating with remote system '%s', svc '%s': %s",
			    instance.PublicIpAddress, p.name, err.trim()));
		os.exit(1);
	    } else if (status == 1 && mithras.verbose) {
		log(sprintf("Service '%s', ok: %t, status: %d %s %s", 
			    p.name, ok, status, out, err));
	    } else if (status == 3 && mithras.verbose) {
		log(sprintf("Service '%s' is not running.", p.name));
	    } else if (mithras.verbose) {
		log(sprintf("Service '%s': status %d; out %s", 
			    p.name, status, out));
	    }
	    return [null, false];
	}
	
	startStop: function(resource, inst, user, key, action) {
	    var p = resource.params;
	    var cmd = "";
	    cmd = sprintf("service %s %s", p.name, action)
	    cmd = become(p.become, p.becomeUser, p.becomeMethod, cmd);
	    var result = mithras.remote.shell(inst.PublicIpAddress, user, key, null, cmd, null);
	    
	    var out = result[0];
	    var err = result[1];
	    var ok = result[2];
	    var status = result[3];
	    if (ok && status == 0) {
		if (mithras.verbose) {
		    log(sprintf("Service '%s': %s", p.name, out.trim()));
		}
		return true;
	    } else if (status == 255) {
		log(sprintf("SSH error communicating with remote system '%s', service '%s': %s %s",
			    inst.PublicIpAddress, p.name, err.trim(), out.trim()));
		os.exit(2);
	    } else if (status == 1) {
		if (mithras.verbose) {
		    log(sprintf("Service '%s' error: %s\n%s", p.name, err, out.trim()));
		}
		os.exit(3);
	    } else {
		if (mithras.verbose) {
		    log(sprintf("Service '%s': status %d; out %s", 
				p.name, 
				status, 
				out.trim()));
		}
	    }
	    return false;
	}

	handle: function(catalog, resources, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
		
	    var params = resource.params;
		    
	    // Sanity
	    if (!params || !params.name) {
		console.log("Invalid service resource params")
		os.exit(1);
	    }
	    
	    // Loop over hosts
	    if (typeof(params.hosts) != "object") {
		return [null, true];
	    }
	    _.each(params.hosts, function(host) {
		var key = mithras.sshKeyPathForInstance(resource, host);
		var user = mithras.sshUserForInstance(resource, host);
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
		params = updated.params;
		var ensure = params.ensure;
		
		var running = handler.running(updated, host, user, key);
		
		switch(ensure) {
		case "absent":
		    if (running) {
			handler.startStop(updated, host, user, key, "stop");
		    }
		    break;
		case "restart":
		    if (running) {
			handler.startStop(updated, host, user, key, "restart");
		    }
		    break;
		case "present":
		    if (!running) {
			handler.startStop(updated, host, user, key, "start");
		    }
		    break;
		}
	    });
	    
	    return [null, true];
	}
    }

    handler.init = function () {
	mithras.modules.handlers.register(handler.moduleNames[0], handler.handle);
	return handler;
    };
    
    return handler;
});
