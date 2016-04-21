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
// # secgroup
// 
// Secgroup is a resource handler for dealing with AWS security groups.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"secgroup"`
// 
// Usage:
// 
// `var secgroup = require("secgroup").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var sgName = "simple-sg";
// var security = {
//      name: "webserverSG"
//      module: "secgroup"
//      params: {
//          region: defaultRegion
//          ensure: ensure
//
//          secgroup: {
//              Description: "Webserver security group"
//              GroupName:   sgName
//          }
//          tags: {
//              Name: "webserver"
//          }
//          ingress: {
//              IpPermissions: [
//                  {
//                      FromPort:   22
//                      IpProtocol: "tcp"
//                      IpRanges: [ {CidrIp: "0.0.0.0/0"} ]
//                      ToPort: 22
//                  },
//                  {
//                      FromPort:   80
//                      IpProtocol: "tcp"
//                      IpRanges: [ {CidrIp: "0.0.0.0/0"} ]
//                      ToPort: 80
//                  }
//              ]
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
// If `"present"` and the security group named
// `params.secgroup.GroupName` does not exist, it is created.  If
// `"absent"`, and the group exists, it is removed.
// 
// ### `secgroup`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-CreateSecurityGroupInput)
//
// This file is copied to remote hosts.
//
// ### `tags`
//
// * Required: false
// * Allowed Values: a map of tags to be set on a created security group
//
// For tagging.
//
// ### `ingress`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-AuthorizeSecurityGroupIngressInput)
//
// Set inbound rules for the SG
//
// ### `egress`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-AuthorizeSecurityGroupEgressInput)
//
// Set outbound rules for the SG
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
        moduleName: "secgroup"
        findInCatalog: function(catalog, resource, vpcId, groupName) {
            if (typeof(resource.params.on_find) === 'function') {
		result = resource.params.on_find(catalog, resource);
		if (!result || 
		    (Array.isArray(result) && result.length == 0)) {
		    return;
		}
		return result;
	    }
            return _.find(catalog.securityGroups, function(s) { 
                return (s.GroupName === groupName && s.VpcId == vpcId);
            });
        }
        handle: function(catalog, resources, resource) {
            if (resource.module != handler.moduleName) {
                return [null, false];
            }

            // Sanity
            if (!resource.params.secgroup) {
                console.log("Invalid secgroup params")
                os.exit(1);
            }

            var ensure = resource.params.ensure;
            var params = resource.params;
            var sg = resource._target;

            switch(ensure) {
            case "absent":
                if (!sg) {
                    if (mithras.verbose) {
                        log(sprintf("Security group not found, no action taken."));
                    }
                    break;
                }
                if (mithras.verbose) {
                    log(sprintf("Deleting security group '%s'",
                               params.secgroup.GroupName));
                }
                aws.securityGroups.delete(params.region, sg.GroupId);
                catalog.securityGroups = _.reject(catalog.securityGroups,
                                                  function(s) { 
                                                      return s.GroupId == sg.GroupId;
                                                  });
                break;
            case "present":
                if (sg) {
                    if (mithras.verbose) {
                        log(sprintf("Security group found, no action taken."));
                    }
                    break;
                }
                // create sg
                if (mithras.verbose) {
                    log(sprintf("Creating security group '%s'", 
                                params.secgroup.GroupName));
                }
                var sg = aws.securityGroups.create(params.region, params.secgroup);

                // do authorizations
                if (params.ingress) {
                    params.ingress.GroupId = sg.GroupId;
                    if (mithras.verbose) {
                        log(sprintf("Creating security group ingress authorization"));
                    }
                    aws.securityGroups.authorizeIngress(params.region, params.ingress);
                }
                if (params.egress) {
                    params.egress.GroupId = sg.GroupId;
                    if (mithras.verbose) {
                        log(sprintf("Creating security group egress authorization"));
                    }
                    aws.securityGroups.authorizeEgress(params.region, params.egress);
                }
                
                // Tag it
                if (params.tags) {
                    if (mithras.verbose) {
                        log(sprintf("Tagging security group."));
                    }
                    aws.tags.create(params.region, sg.GroupId, params.tags);
                }
                
                // Reload it to get tags
                var sg = aws.securityGroups.describe(params.region, sg.GroupId);
                
                // add to catalog
                catalog.securityGroups.push(sg);
                
                // return it
                return [sg, true];
            }
            return [null, true];
        }
        preflight: function(catalog, resources, resource) {
            if (resource.module != handler.moduleName) {
                return [null, false];
            }
            var params = resource.params;
            var sg = params.secgroup;
            var s = handler.findInCatalog(catalog, resource, sg.VpcId, sg.GroupName);
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
