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
// # route53
// 
// Route53 is resource handler for working with AWS's route53 dns system.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"route53"`
// 
// Usage:
// 
// `var route53 = require("route53").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rElbDnsEntry = {
//     name: "elbDnsEntry"
//     module: "route53"
//     dependsOn: [rElb.name]
//     params: {
//         region: defaultRegion
//         ensure: ensure
//         domain: "mithras.io."
//         resource: {
//             Name: "test.mithras.io."
//             Type: "A"
//             AliasTarget: {
//                 DNSName:              mithras.watch("elb._target.DNSName")
//                 EvaluateTargetHealth: true
//                 HostedZoneId:         mithras.watch("elb._target.CanonicalHostedZoneNameID")
//             }
//         }
//     } // params
// };
// ```
// 
// ## Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "present", "absent"
//
// If `"present"`, the dns entry will be created if it doesn't already
// exist.  If `"absent"`, the dns entry will be removed if it is
// present.
// 
// ### `resource`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/route53.html#type-ResourceRecordSet)
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
        moduleName: "route53"
        findInCatalog: function(catalog, resource, set) {
            if (typeof(resource.params.on_find) === 'function') {
		result = resource.params.on_find(catalog, resource);
		if (!result || 
		    (Array.isArray(result) && result.length == 0)) {
		    return;
		}
		return result;
	    }
            var alias = set.AliasTarget;
            var rrs = set.ResourceRecords;
            var name = set.Name;
            var type = set.Type;
            return _.find(catalog.rrs, function(r) { 
                if (r.Name != name || r.Type != type) {
                    return false;
                }
                if (alias && r.AliasTarget) {
                    return true;
                } else if (rrs && r.ResourceRecords.length > 0) {
                    return true;
                }
            });
        }
        zoneForDomain: function(catalog, domain) {
            return _.find(catalog.zones, function(z) { 
                return z.Name === domain;
            });
        }
        handle: function(catalog, resources, resource) {
            if (resource.module != handler.moduleName) {
                return [null, false];
            }

            // Sanity
            if (!resource.params.resource) {
                console.log("Invalid route53 params")
                exit(1);
            }

            var ensure = resource.params.ensure;
            var params = resource.params;
            var r = resource._target;

            switch(ensure) {
            case "absent":
                if (!r) {
                    if (mithras.verbose) {
                        log(sprintf("DNS entry '%s' not found, no action taken.", 
                                    params.resource.Name));
                    }
                    break;
                }
                if (mithras.verbose) {
                    log(sprintf("Deleting DNS entry '%s'", params.resource.Name));
                }
                var zone = handler.zoneForDomain(catalog, params.domain);
                if (!zone) {
                    if (mithras.verbose) {
                        console.log(sprintf("Can't find zone for domain '%s'", 
                                            params.domain));
                    }
                    os.exit(2);
                }
                var zoneId = zone.Id;
                aws.route53.rrs.delete(params.region, zoneId, params.resource);
                break;
            case "present":
                if (r) {
                    if (mithras.verbose) {
                        log(sprintf("DNS entry '%s' found, no action taken.", 
                                    params.resource.Name));
                    }
                    break;
                }
                if (mithras.verbose) {
                    log(sprintf("Creating DNS entry '%s'.", params.resource.Name));
                }

                // create 
                var zone = handler.zoneForDomain(catalog, params.domain);
                if (!zone) {
                    log(sprintf("Can't find zone for domain '%s'", 
                                params.domain));
                    os.exit(2);
                }
                var zoneId = zone.Id;
                created = aws.route53.rrs.create(params.region, zoneId, params.resource);

                // add to catalog
                catalog.rrs.push(created);
                
                // return it
                return [r, true];
            }
            return [null, true];
        }
        preflight: function(catalog, resources, resource) {
            if (resource.module != handler.moduleName) {
                return [null, false];
            }
            var params = resource.params;
            var r = params.resource;
            var t = handler.findInCatalog(catalog, resource, r);
            if (t) {
                return [t, true];
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
