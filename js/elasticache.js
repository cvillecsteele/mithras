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
// # Elasticache
// 
// Elasticache is resource handler for managing AWS caches.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"elasticache"`
// 
// Usage:
// 
// `var elasticache = require("elasticache").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rCache = {
//      name: "redis"
//      module: "elasticache"
//      dependsOn: [otherResource.name]
//      params: {
//          ensure: ensure
//          region: "us-east-1"
//          wait: true
//          subnetGroup: {
//              CacheSubnetGroupDescription: "Redis Subnet Group"
//              CacheSubnetGroupName:        "redis-subnet-group"
//              SubnetIds: [
//                  "subnet-123",
//                  "subnet-456"
//              ]
//          }
//          cache: {
//              CacheClusterId:          "test-redis"
//              AutoMinorVersionUpgrade: true
//              CacheNodeType:           "cache.t2.small"
//              CacheSubnetGroupName:    "redis-subnet-group"
//              Engine:                  "redis"
//              NumCacheNodes:           1
//              SecurityGroupIds:        []
//              Tags: [
//                  {
//                      Key:   "Name"
//                      Value: "test-cluster"
//                  },
//              ]
//          }
//          delete: {
//              CacheClusterId:          "test-redis"
//          }
//
//      }
// };
// ```
// 
// ## Copy Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "present" or "absent"
//
// If `"present"`, the cache specified by `cache` will be created, and
// if `"absent"`, the cache will be removed using the `delete` property.
// 
// ### `region`
//
// * Required: true
// * Allowed Values: string, any valid AWS region; eg "us-east-1"
//
// The region for calls to the AWS API.
// 
// ### `wait`
//
// * Required: false
// * Allowed Values: true or false
//
// If `true`, delay execution until the cache has been created in AWS.
// 
// ### `subnetGroup`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticache.html#type-CreateCacheSubnetGroupInput)
//
// If set, a subnet group will be created for your cache.
// 
// ### `cache`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticache.html#type-CreateCacheInput)
//
// Parameters for cache creation.
// 
// ### `delete`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticache.html#type-DeleteCacheInput)
//
// Parameters for cache deletion.
// 
// ### `on_find`
//
// * Required: false
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
                os.exit(3);
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
        findInCatalog: function(catalog, resource, id) {
            if (typeof(resource.params.on_find) === 'function') {
		result = resource.params.on_find(catalog, resource);
		if (!result || 
		    (Array.isArray(result) && result.length == 0)) {
		    return;
		}
		return result;
	    }
            return _.find(catalog.caches, function(inst) { 
                return inst.CacheClusterId === id;
            });
        }
        preflight: function(catalog, resources, resource) {
            if (!_.find(handler.moduleNames, function(m) { 
                return resource.module === m; 
            })) {
                return [null, false];
            }
            var cache = handler.findInCatalog(catalog, 
					      resource,
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
