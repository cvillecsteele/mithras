(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var become = require("become").become;
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleNames: ["elasticache"]
	handle: function(catalog, resources, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }

	    var p = resource.params;
	    var ensure = p.ensure;
	    var cache = resource._target;
		    
	    // Sanity
	    if (!p || !p.cache || !p.delete) {
		console.log(sprintf("Invalid elasticache resource '%s'", resource.name));
		os.exit(1);
	    }

	    switch(ensure) {
	    case "absent":

		if (!cache) {
		    if (mithras.verbose) {
			log(sprintf("Elasticache '%s' not found, no action taken.", p.cache.CacheClusterId));
		    }
		    break;
		}
		var id = resource._target.CacheClusterId;
		
		// delete
		if (mithras.verbose) {
		    log(sprintf("Deleting elasticache instance '%s'", id));
		}
		aws.elasticache.delete(p.region, p.delete);

		// remove subnet group
		if (p.subnetGroup &&
		    aws.elasticache.subnetGroups.describe(p.region, 
							  p.subnetGroup.CacheSubnetGroupName)) {
		    if (mithras.verbose) {
			log(sprintf("Deleting elasticache subnet group '%s'", 
				    p.subnetGroup.CacheSubnetGroupName));
		    }
		    aws.elasticache.subnetGroups.delete(p.region, p.subnetGroup.CacheSubnetGroupName);
		}
		
		// remove from catalog
		catalog.caches = _.reject(catalog.caches, 
				       function(x) { 
					   return x.CacheClusterId == id;
				       });
		break;
	    case "present":
		if (cache) {
		    if (mithras.verbose) {
			log(sprintf("Elasticache '%s' found, no action taken.", p.cache.CacheClusterId));
		    }
		    break;
		}
		var id = p.cache.CacheClusterId;

		// create subnet group
		if (p.subnetGroup && 
		    (!aws.elasticache.subnetGroups.describe(p.region, 
						    p.subnetGroup.CacheSubnetGroupName))) {
		    if (mithras.verbose) {
			log(sprintf("Creating elasticache subnet group '%s'", 
				    p.subnetGroup.CacheSubnetGroupName));
		    }
		    aws.elasticache.subnetGroups.create(p.region, p.subnetGroup);
		}
		
		// create instance
		if (mithras.verbose) {
		    log(sprintf("Creating elasticache instance '%s' (WAIT FOR IT...)", id));
		}
		aws.elasticache.create(p.region, p.cache, p.wait);

		// re-describe it
		cache = aws.elasticache.describe(p.region, id);

		// add to catalog
		catalog.caches.push(cache);
		break;
	    }
	    return [null, true];
	}
	findInCatalog: function(catalog, id) {
	    return _.find(catalog.caches, function(inst) { 
		return inst.CacheClusterId === id;
	    });
	}
	preflight: function(catalog, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
	    var cache = handler.findInCatalog(catalog, 
					      resource.params.cache.CacheClusterId);
	    if (cache) {
		return [cache, true];
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
