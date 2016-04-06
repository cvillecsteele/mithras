(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleNames: ["instance"]
	findInCatalog: function(catalog, resource) {
	    if (!typeof(resource.params.on_find) === 'function') {
		console.log(sprintf("Instance resource '%s' has no 'on_find' param.", 
				    resource.name));
		os.exit(3);
	    }
	    return resource.params.on_find(catalog);
	}
	handle: function(catalog, resources, resource) {
	    if (!_.find(handler.moduleNames, function(m) { return resource.module === m; })) {
		return [null, false];
	    }
		
	    var ensure = resource.params.ensure;
	    var params = resource.params;
	    
	    if (!resource.params.instance) {
		console.log("Invalid instance params")
		os.exit(1);
	    }
	    var found = resource._target || [];
	    var matchingCount = found.length;

	    if (mithras.verbose) {
		log(sprintf("Found %d matching instances, max %d, min %d",
			    matchingCount,
			    params.instance.MaxCount,
			    params.instance.MinCount));
	    }
	    
	    switch(ensure) {
	    case "absent":
		if (found.length == 0 && mithras.verbose) {
		    log(sprintf("No action taken."));
		    break;
		}
		for (var idx in found) {
		    var inst = found[idx];
		    aws.instances.delete(params.region, 
					 inst.InstanceId);
		    catalog.instances = 
			_.reject(catalog.instances,
				 function(i) { 
				     return i.InstanceId == inst.InstanceId;
				 });
		}
		break;
	    case "present":
		// Too many?
		if (matchingCount > params.instance.MaxCount) {
		    var numToDelete = matchingCount - params.instance.MaxCount;
		    for (var idx in _.range(numToDelete)) {
			var inst = found[idx];
			aws.instances.delete(params.region, 
					     inst.InstanceId);
			catalog.instances = 
			    _.reject(catalog.instances,
				     function(i) { 
					 return i.InstanceId == inst.InstanceId;
				     });
		    }
		} else if (matchingCount < params.instance.MinCount) {
		    var numToAdd = params.instance.MaxCount - matchingCount;
		    
		    // create 
		    params.instance.MaxCount = numToAdd;
		    created = aws.instances.create(params.region, params.instance);
		    
		    // Set tags
		    for (var idx in created) {
			aws.tags.create(params.region, 
					created[idx].InstanceId,
					params.tags);
		    }
		    
		    // describe 'em (to get tags) and add to catalog
		    for (var idx in created) {
			var inst = aws.instances.describe(params.region,
							  created[idx].InstanceId);
			catalog.instances.push(inst);
		    }
		} else {
		    if (mithras.verbose) {
			log(sprintf("No action taken."));
		    }
		}
		
		// return 'em
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
	    var params = resource.params;
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
