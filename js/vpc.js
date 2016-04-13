(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleName: "vpc"
	findVPCInCatalog: function(catalog, cidr) {
	    return _.find(catalog.vpcs, function(vpc) { 
		return vpc.CidrBlock === cidr;
	    });
	}
	findGWInCatalog: function(catalog, vpcId) {
	    return _.find(catalog.gateways, function(gw) { 
		// skip detached gateways
		if (gw.Attachments) {
		    return gw.Attachments[0].VpcId === vpcId;
		}
	    });
	}
	handle: function(catalog, resources, resource) {
	    if (resource.module != handler.moduleName) {
		return [null, false];
	    }

	    // Sanity
	    if (!resource.params.vpc) {
		console.log("Invalid vpc params")
		os.exit(1);
	    }

	    var params = resource.params;
	    var ensure = resource.params.ensure;
	    var cidr = params.vpc.CidrBlock;
	    var vpc = resource._target;

	    switch(ensure) {
	    case "absent":
		if (vpc) {
		    var vpcId = resource._target.VpcId;

		    // delete gw
		    var gw = handler.findGWInCatalog(catalog, vpcId);
		    if (gw) {
			aws.vpcs.gateways.delete(params.region,
						 vpcId,
						 gw.InternetGatewayId)
		    }

		    // delete vpc
		    if (mithras.verbose) {
			log(sprintf("Deleting VPC '%s'", vpcId));
		    }
		    aws.vpcs.delete(params.region, vpcId)

		    // remove both from catalog
		    if (gw) {
			catalog.gateways = _.reject(catalog.gateways, 
						    function(g) { 
							return g.InternetGatewayId == gw.InternetGatewayId;
						    });
		    }
		    catalog.vpcs = _.reject(catalog.vpcs, 
						function(v) { 
						    return v.VpcId == vpcId;
						});
		} else {
		    log(sprintf("VPC not found, no action taken."));
		}
		break;
	    case "present":
		if (vpc) {
		    log(sprintf("VPC found, no action taken."));
		} else {
		    // create vpc & gw
		    if (mithras.verbose) {
			log(sprintf("Creating VPC with cidr '%s'", 
				    params.vpc.CidrBlock));
		    }

		    var result = aws.vpcs.create(params.region,
						 params.vpc, 
						 params.gateway);
		    var newVPC = result[0];
		    var newGW = result[1];

		    // Tag VPC
		    if (params.tags) {
			aws.tags.create(params.region, newVPC.VpcId, params.tags);
		    }

		    // Reload it to get tags, associations
		    var newVPC = aws.vpcs.describe(params.region, newVPC.VpcId);
		    newGW = aws.vpcs.gateways.describe(params.region, 
						       newGW.InternetGatewayId);
		    do {
			newGW = aws.vpcs.gateways.describe(params.region, 
							   newGW.InternetGatewayId);
			time.sleep(10);
		    } while ((newGW.InternetGatewayId == null) ||
			     (newGW.Attachments == null));
		    

		    // add both to catalog
		    catalog.vpcs.push(newVPC);
		    if (newGW) {
			catalog.gateways.push(newGW);
		    }
		    resource._target = newVPC;
		    return [newVPC, true];
		}
		break;
	    }
	    return [null, true];
	}
	preflight: function(catalog, resource) {
	    if (resource.module != handler.moduleName) {
		return [null, false];
	    }
	    var vpc = handler.findVPCInCatalog(catalog, resource.params.vpc.CidrBlock);
	    if (vpc) {
		return [vpc, true];
	    }
	    return [null, true];
	}
    };
    
    handler.init = function () {
	mithras.modules.preflight.register(handler.moduleName, handler.preflight);
	mithras.modules.handlers.register(handler.moduleName, handler.handle);
	return handler;
    };
    
    return handler;
});
