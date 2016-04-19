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
// # subnet
// 
// Subnet is a resource handler for dealing with AWS security groups.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"subnet"`
// 
// Usage:
// 
// `var subnet = require("subnet").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rSubnetA = {
//     name: "subnetA"
//     module: "subnet"
//     dependsOn: [rVpc.name]
//     params: {
//         region: defaultRegion
//         ensure: ensure
//         subnet: {
//             CidrBlock:        "172.33.1.0/24"
//             VpcId:            mithras.watch("VPC._target.VpcId")
//             AvailabilityZone: defaultZone
//         }
//         tags: {
//             Name: "primary-subnet"
//         }
//         routes: [
//             {
//                 DestinationCidrBlock: "0.0.0.0/0"
//                 GatewayId:            mithras.watch("VPC._target.VpcId", mithras.findGWByVpcId)
//             }
//         ]
//     }
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
// If `"present"` and the subnet named
// `params.subnet.GroupName` does not exist, it is created.  If
// `"absent"`, and the subnet exists, it is removed.
// 
// ### `subnet`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-CreateSubnetInput)
//
// This file is copied to remote hosts.
//
// ### `tags`
//
// * Required: false
// * Allowed Values: a map of tags to be set on a created subnet
//
// For tagging.
//
// ### `routes`
//
// * Required: false
// * Allowed Values: an array of JSON objects corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-CreateRouteInput)
//
// An list of routes to be created and associated with the subnet.
//
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
                        if (mithras.verbose) {
                            log(sprintf("Deleting route table association '%s'",
                                        a.RouteTableAssociationId));
                        }
                        aws.routeTables.disassociate(params.region, 
                                                     a.RouteTableAssociationId);
                    }

                    // Delete table
                    if (mithras.verbose) {
                        log(sprintf("Deleting route table '%s'",
                                    t.RouteTableId));
                    }
                    aws.routeTables.delete(params.region, t.RouteTableId);
                    
                    // Remove table from catalog
                    catalog.routeTables = 
                        _.reject(catalog.routeTables, 
                                 function(rt) { 
                                     return rt.RouteTableId == t.RouteTableId;
                                 });
                }
                
                // Throw away subnet
                if (mithras.verbose) {
                    log(sprintf("Deleting subnet '%s'", subnet.SubnetId));
                }
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
                if (mithras.verbose) {
                    log(sprintf("Creating subnet with cidr '%s'", params.subnet.CidrBlock));
                }
                var subnet = aws.subnets.create(params.region, params.subnet);
                
                // Tag it
                aws.tags.create(params.region, subnet.SubnetId, params.tags);
                
                // Reload it to get tags
                var subnet = aws.subnets.describe(params.region, subnet.SubnetId);
                
                // create route table and associate subnet with it
                if (mithras.verbose) {
                    log(sprintf("Creating route table for subnet '%s'",
                               subnet.SubnetId));
                }
                var rt = aws.routeTables.create(params.region, params.subnet.VpcId);
                if (mithras.verbose) {
                    log(sprintf("Associating route table with subnet '%s'",
                               subnet.SubnetId));
                }
                aws.routeTables.associate(params.region, 
                                          subnet.SubnetId, 
                                          rt.RouteTableId);
                
                // for each route in params, create route using routetableid
                for (var idx in params.routes) {
                    var r = params.routes[idx];
                    r.RouteTableId = rt.RouteTableId;
                    if (mithras.verbose) {
                        log(sprintf("Creating route."));
                    }
                    if (typeof(r.GatewayId) != 'string') {
                        console.log(JSON.stringify(resource, null, 2));
                        console.log(JSON.stringify(catalog.gateways, null, 2));
                        os.exit(1);
                    }
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
