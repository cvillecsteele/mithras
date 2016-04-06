(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
        module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
        moduleName: "subnet"
        findSubnetInCatalog: function(catalog, vpcId, cidr) {
            return _.find(catalog.subnets, function(s) { 
                return (s.CidrBlock === cidr && s.VpcId == vpcId);
            });
        }
        handle: function(catalog, resources, resource) {
            if (resource.module != handler.moduleName) {
                return [null, false];
            }

            // Sanity
            if (!resource.params.subnet) {
                console.log("Invalid subnet params")
                os.exit(1);
            }

            var ensure = resource.params.ensure;
            var params = resource.params;
            var cidr = params.subnet.CidrBlock;
            var vpcId = params.subnet.VpcId;
            var subnet = resource._target;

            switch(ensure) {
            case "absent":
                if (!subnet) {
		    if (mithras.verbose) {
			log(sprintf("No action taken."));
		    }
                    break;
                }
                var tables = aws.routeTables.describeForSubnet(params.region, 
							       subnet.SubnetId);
                for (var idx in tables) {
                    var t = tables[idx];
                    
                    // Delete its routes; TODO this seems to be
                    // superfluous and done for us...?

                    // for (var i in t.Routes) {
                    //  var r = t.Routes[i];
                    //  aws.subnets.routes.delete(params.region, 
                    //                            r.DestinationCidrBlock, 
                    //                            t.RouteTableId);
                    // }
                    
                    // Delete its associations
                    for (var i in t.Associations) {
                        var a = t.Associations[i];
                        aws.routeTables.disassociate(params.region, 
                                                     a.RouteTableAssociationId);
                    }

		    // Delete table
                    aws.routeTables.delete(params.region, t.RouteTableId);
                    
                    // Remove table from catalog
                    catalog.routeTables = 
                        _.reject(catalog.routeTables, 
                                 function(rt) { 
                                     return rt.RouteTableId == t.RouteTableId;
                                 });
                }
                
                // Throw away subnet
                aws.subnets.delete(params.region, subnet.SubnetId);
                catalog.subnets = _.reject(catalog.subnets,
                                           function(s) { 
                                               return s.SubnetId == subnet.SubnetId;
                                           });
                break;
            case "present":
                if (subnet) {
		    if (mithras.verbose) {
			log(sprintf("No action taken."));
		    }
                    break;
                }
                // create subnet
                var subnet = aws.subnets.create(params.region, params.subnet);
                
                // Tag it
                aws.tags.create(params.region, subnet.SubnetId, params.tags);
                
                // Reload it to get tags
                var subnet = aws.subnets.describe(params.region, subnet.SubnetId);
                
                // create route table and associate subnet with it
                var rt = aws.routeTables.create(params.region, params.subnet.VpcId);
                aws.routeTables.associate(params.region, subnet.SubnetId, rt.RouteTableId);
                
                // for each route in params, create route using routetableid
                for (var idx in params.routes) {
                    var r = params.routes[idx];
                    r.RouteTableId = rt.RouteTableId;
                    aws.subnets.routes.create(params.region, r);
                }
                
                // add table and subnet to catalog
                catalog.routeTables.push(rt);
                catalog.subnets.push(subnet);
                
                // return it
                return [subnet, true];
            }
            return [null, true];
        }
        preflight: function(catalog, resource) {
            if (resource.module != handler.moduleName) {
                return [null, false];
            }
            var params = resource.params;
            var s = handler.findSubnetInCatalog(catalog, params.subnet.VpcId, params.subnet.CidrBlock);

            if (s) {
                return [s, true];
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
