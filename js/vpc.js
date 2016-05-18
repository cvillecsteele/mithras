// MITHRAS: Javascript configuration management tool for AWS.
// Copyright (C) 2016, Colin Steele
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//   the Free Software Foundation, either version 3 of the License, or
//                  (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//              GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// @public
// 
// # vpc
// 
// Vpc is a resource handler for dealing with AWS security groups.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"vpc"`
// 
// Usage:
// 
// `var vpc = require("vpc").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rVpc = {
//      name: "VPC"
//      module: "vpc"
//      params: {
//          region: defaultRegion
//          ensure: ensure
//          vpc: {
//              CidrBlock:       "172.33.0.0/16"
//          }
//          gateway: true
//          tags: {
//              Name: "my-vpc"
//          }
//      }
// };
// ```
// 
// ## Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "present" or "absent"
//
// If `"present"` and the vpc named
// `params.vpc.GroupName` does not exist, it is created.  If
// `"absent"`, and the vpc exists, it is removed.
// 
// ### `vpc`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-CreateVpcInput)
//
// Parameters for VPC creation.
//
// ### `tags`
//
// * Required: false
// * Allowed Values: a map of tags to be set on a created vpc
//
// For tagging.
//
// ### `gateway`
//
// * Required: false
// * Allowed Values: true or false
//
// If true, an internet gateway is created for the VPC.
//
// ### `on_find`
//
// * Required: true
// * Allowed Values: A function taking two parameters: `catalog` and `resource`
//
// If defined in the resource's `params` object, the `on_find`
// function provides a way for a matching resource to be identified
// using a user-defined way.  The function is called with the current
// `catalog`, as well as the `resource` object itself.  The function
// can look through the catalog, find a matching object using whatever
// logic you want, and return it.  If the function returns `undefined`
// or a n empty Javascript array, (`[]`), the function is indicating
// that no matching resource was found in the `catalog`.
// 
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
        module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
        moduleNames: ["vpc"]
        findVPCInCatalog: function(catalog, resource, cidr) {
            if (typeof(resource.params.on_find) === 'function') {
		result = resource.params.on_find(catalog);
		if (!result || 
		    (Array.isArray(result) && result.length == 0)) {
		    return;
		}
		return result;
	    }
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
            if (!_.find(handler.moduleNames, function(m) { 
                return resource.module === m; 
            })) {
                return [null, false];
            }

            // Sanity
            if (!resource.params.vpc) {
                console.log("Invalid vpc params")
                os.exit(3);
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
        preflight: function(catalog, resources, resource) {
            if (!_.find(handler.moduleNames, function(m) { 
                return resource.module === m; 
            })) {
                return [null, false];
            }
            var vpc = handler.findVPCInCatalog(catalog, 
					       resource,
					       resource.params.vpc.CidrBlock);
            if (vpc) {
                return [vpc, true];
            }
            return [null, true];
        }
    };
    
    handler.init = function () {
        _.each(handler.moduleNames, function(name) {
            mithras.modules.preflight.register(name, handler.preflight);
            mithras.modules.handlers.register(name, handler.handle);
        });
        return handler;
    };
    
    return handler;
});
