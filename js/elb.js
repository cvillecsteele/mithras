(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleNames: ["elb", "elbMembership"]
	findInCatalog: function(catalog, elbName) {
	    return _.find(catalog.elbs, function(elb) {
		return elb.LoadBalancerName === elbName;
	    });
	}
	handle: function(catalog, resources, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
		
	    var ensure = resource.params.ensure;
	    var params = resource.params;
	    
	    switch(resource.module) {
	    case "elb":
		// Sanity
		if (!resource.params.elb) {
		    console.log("Invalid elb params")
		    exit(1);
		}
		var elb = resource._target;
		
		switch(ensure) {
		case "absent":
		    if (!elb) {
			log(sprintf("ELB not found, no action taken."));
			break;
		    }
		    
		    // Remove it
		    aws.elbs.delete(params.region, params.elb.LoadBalancerName);

		    break;
		case "present":
		    if (elb) {
			log(sprintf("ELB found, no action taken."));
			break;
		    }
		    
		    // create 
		    created = aws.elbs.create(params.region, params.elb);
		    
		    // Set health
		    aws.elbs.setHealth(params.region, 
				       params.elb.LoadBalancerName, 
				       params.health);
		    
		    // Set attrs
		    aws.elbs.setAttrs(params.region, 
				      params.elb.LoadBalancerName, 
				      params.attributes);
		    
		    // add to catalog
		    catalog.elbs.push(created);
		    
		    // return it
		    return [created, true];
		}
		return [null, true];
		break;
	    case "elbMembership":
		var elb = handler.findInCatalog(catalog, 
						params.membership.LoadBalancerName);

		if (ensure === "present") {
		    if (!elb) { 
			console.log(sprintf("Can't find elb '%s' in catalog", 
					    params.membership.LoadBalancerName));
			os.exit(1);
		    }
		} else if (ensure === "absent") {
		    if (!elb) {
			return [null, true];
		    }
		}

		// intersection of input and LB
		lbInstanceIds = _.pluck(elb.Instances, "InstanceId");
		inInstanceIds = _.pluck(params.membership.Instances, "InstanceId");
		var inBoth = _.intersection(lbInstanceIds, inInstanceIds);
		// What's in the LB minus what's in input
		var inLB = _.difference(lbInstanceIds, inInstanceIds);
		// inputs minus what's in the LB
		var inInput = _.difference(inInstanceIds, lbInstanceIds);

		// build map id -> instance
		var byIds = _.reduce(elb.Instances, function(memo, i) { 
		    memo[i.InstanceId] = i;
		    return memo;
		}, {});
		byIds = _.reduce(params.membership.Instances, function(memo, i) { 
		    memo[i.InstanceId] = i;
		    return memo;
		}, byIds);

		// Turn lists of ids back into lists of instances
		inInput = _.reduce(inInput, function(memo, id) {
		    memo.push(byIds[id]);
		    return memo;
		}, []);
		inLB = _.reduce(inLB, function(memo, id) {
		    memo.push(byIds[id]);
		    return memo;
		}, []);

		switch(ensure) {
		case "absent":
		    if (inInput.length > 0) {
			aws.elbs.deRegister(params.region, 
					    params.membership.LoadBalancerName, 
					    inInput);
		    } else {
			if (mithras.verbose) {
			    log(sprintf("No action taken."));
			}
		    }
		    break;
		case "present":
		    if (inInput.length > 0) {
			aws.elbs.register(params.region, 
					  params.membership.LoadBalancerName, 
					  inInput);
		    } else {
			if (mithras.verbose) {
			    log(sprintf("No action taken."));
			}
		    }
		    break;
		case "converge":
		    if (inInput.length > 0) {
			aws.elbs.register(params.region, 
					  params.membership.LoadBalancerName, 
					  inInput);
		    }
		    if (inLB.length > 0) {
			aws.elbs.deRegister(params.region, 
					    params.membership.LoadBalancerName, 
					    inLB);
		    }
		    break;
		}
		return [null, true];
		break;
	    }
	}
	preflight: function(catalog, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
	    if (resource.module === handler.moduleNames[0]) {
		var params = resource.params;
		var elb = params.elb;
		var t = handler.findInCatalog(catalog, elb.LoadBalancerName);
		if (t) {
		    return [t, true];
		}
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
