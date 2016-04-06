(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var become = require("become").become;
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleNames: ["packager"]
	checkInstalled: function(resource, instance, user, key) {
	    var p = resource.params;
	    var cmd = sprintf("yum list installed %s", p.name);
	    cmd = become(p.become, p.becomeUser, p.becomeMethod, cmd);
	    cmd = cmd.split(/[ ]+/);
	    result = mithras.remote.wrapper(instance, user, key, cmd, null);
	    var out = result[0];
	    var err = result[1];
	    var ok = result[2];
	    var status = result[3];
	    var version = "";
	    if (ok && status == 0) {
		var lines = out.split("\n");
		for (var i = 0; i < lines.length; i++) {
		    var parts = lines[i].split(/[ \t]/);
		    for (var j = 0; j < parts.length; j++) {
			if (parts[j].match(/^[0-9].*/)) {
			    version = parts[j];
			}
		    }
		}
		if (mithras.verbose && version != "") {
		    log(sprintf("Package '%s' found version '%s'",
				p.name,
				version));
		    return true;
		}
	    } else if (status == 255) {
		log(sprintf("Error communiating with remote system, package '%s': %s",
			    p.name, err.trim()));
		os.exit(1);
	    } else if (status == 1 && mithras.verbose) {
		log(sprintf("Package '%s' not found", p.name));
	    } else if (mithras.verbose) {
		    log(sprintf("Package '%s': status %d; out %s", 
				p.name, status, out));
	    }
	    return false;
	}

	install: function(resource, inst, user, key) {
	    var p = resource.params;
	    var cmd = "";
	    switch (p.ensure) {
	    case "present":
		cmd = sprintf("yum -y install %s", p.name);
		break;
	    case "absent":
		cmd = sprintf("yum -y remove %s", p.name);
		break;
	    case "latest":
		cmd = sprintf("yum -y update %s", p.name);
		break;
	    }
	    cmd = become(p.become, p.becomeUser, p.becomeMethod, cmd).trim().split(/[ \t]+/);
	    var result = mithras.remote.wrapper(inst, user, key, cmd, null);
	    var out = result[0].trim();
	    var err = result[1].trim();
	    var ok = result[2];
	    var status = result[3];
	    if (ok && status == 0) {
		if (mithras.verbose) {
		    log(sprintf("Package '%s' %s", p.name ? p.name : "ALL", p.ensure));
		}
		return true;
	    } else if (status == 255) {
		log(sprintf("Error communiating with remote system, package '%s': %s",
			    p.name, err.trim()));
		os.exit(2);
	    } else if (status == 1) {
		if (mithras.verbose) {
		    log(sprintf("Package '%s' error: %s", p.name, err));
		}
		os.exit(3);
	    } else {
		if (mithras.verbose) {
		    log(sprintf("Package '%s': status %d; out %s", 
				p.name, 
				status, 
				out));
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
	    if (!params || !(typeof(params.name) === 'string')) {
		console.log("Invalid packager params", JSON.stringify(params, null, 2));
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
		
		switch(ensure) {
		case "absent":
		    if (handler.checkInstalled(updated, host, user, key)) {
			handler.install(updated, host, user, key);
		    }
		    break;
		case "present":
		    if (!handler.checkInstalled(updated, host, user, key)) {
			handler.install(updated, host, user, key);
		    }
		    break;
		case "latest":
		    handler.install(updated, host, user, key);
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
