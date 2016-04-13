(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleNames: ["keypairs"]
	findInCatalog: function(catalog, resource) {
	    return _.find(catalog.keypairs, function(key) { 
		return key.KeyName === resource.params.key.KeyName;
	    });
	}
	handle: function(catalog, resources, resource) {
	    if (!_.find(handler.moduleNames, function(m) { return resource.module === m; })) {
		return [null, false];
	    }
		
	    var ensure = resource.params.ensure;
	    var params = resource.params;
	    
	    if (!params.ensure || !params.key || !params.savePath || !params.region) {
		console.log("Invalid key params")
		os.exit(1);
	    }
	    var found = resource._target;

	    switch(ensure) {
	    case "absent":
		if (!found) {
		    if (mithras.verbose) {
			log(sprintf("Key '%s' not found, no action taken.", params.key.KeyName));
			break;
		    }
		}
		if (mithras.verbose) {
		    log(sprintf("Deleting key '%s'", params.key.KeyName));
		}
		aws.keypairs.delete(params.region, params.key.KeyName);
		catalog.keypairs = 
		    _.reject(catalog.keypairs,
			     function(k) { 
				 return k.KeyName == params.key.KeyName;
			     });
		break;
	    case "present":
		if (found) {
		    if (mithras.verbose) {
			log(sprintf("Key '%s' found, no action taken.", params.key.KeyName));
			break;
		    }
		}

		if (mithras.verbose) {
		    log(sprintf("Creating key '%s'", params.key.KeyName));
		}
		var raw = aws.keypairs.create(params.region, params.key.KeyName);
		if (mithras.verbose) {
		    log(sprintf("Writing key '%s' to '%s'", params.key.KeyName, params.savePath));
		}
		fs.write(params.savePath, raw, 0400);

		var key = aws.keypairs.describe(params.region, params.key.KeyName);
		do {
		    time.sleep(10);
		    key = aws.keypairs.describe(params.region, params.key.KeyName);
		} while (!key);
		catalog.keypairs.push(key);

		// return it
		return [handler.findInCatalog(catalog, resource), true];
		break;
	    }
	    return [null, true];
	}
	preflight: function(catalog, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
	    var found = handler.findInCatalog(catalog, resource);
	    if (found) {
		return [found, true];
	    }
	    return [null, true];
	}
    };
    
    handler.init = function () {
	mithras.modules.preflight.register(handler.moduleNames[0], handler.preflight);
	mithras.modules.handlers.register(handler.moduleNames[0], handler.handle);
	return handler;
    };
    
    return handler;
});
