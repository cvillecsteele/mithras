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
// # Instance
// 
// Instance is resource handler for managing AWS caches.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"instance"`
// 
// Usage:
// 
// `var instance = require("instance").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var webserverTagName: "my-web-server";
// var keyname: "test-key";
// var webServer = {
//      name: "webserver"
//      module: "instance"
//      dependsOn: [resource.name]
//      params: {
//          region: "us-east-1"
//          ensure: "present"
//          on_find: function(catalog) {
//              var matches = _.filter(catalog.instances, function (i) {
//                  if (i.State.Name != "running") {
//                      return false;
//                  }
//                  return (_.where(i.Tags, {"Key": "Name", 
//                                           "Value": webserverTagName}).length > 0);
//              });
//              return matches;
//          }
//          tags: {
//              Name: webserverTagName
//          }
//          instance: {
//              ImageId:        "ami-60b6c60a"
//              MaxCount:       1
//              MinCount:       1
//              DisableApiTermination: false
//              EbsOptimized:          false
//              IamInstanceProfile: {
//                  Name: iamProfileName
//              }
//              InstanceInitiatedShutdownBehavior: "terminate"
//              InstanceType:                      "t2.small"
//              KeyName:                           keyName
//              Monitoring: {
//                  Enabled: true
//              }
//              NetworkInterfaces: [
//                  {
//                      AssociatePublicIpAddress: true
//                      DeleteOnTermination:      true
//                      DeviceIndex:              0
//                      Groups:                  [ "sg-abc" ]
//                      SubnetId:                "subnet-123"
//                  }
//              ]
//          } // instance
//      } // params
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
// If `"present"`, the db specified by `db` will be created, and
// if `"absent"`, it will be removed using the `delete` property.
// 
// ### `region`
//
// * Required: true
// * Allowed Values: string, any valid AWS region; eg "us-east-1"
//
// The region for calls to the AWS API.
// 
// ### `instance`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-RunInstancesInput)
//
// Parameters for resource creation.
// 
// ### `on_find`
//
// * Required: true
// * Allowed Values: `function(catalog) { ... }`
//
// Called to find matching instances in `catalog.instances`.  Returns
// an array of matching EC2 instance objects.
// 
// ### `tags`
//
// * Required: false
// * Allowed Values: A map of tags to be applied to created instances
//
// A map of tags to be applied to created instances
// 
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
        module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
        moduleNames: ["instance"]
        findInCatalog: function(catalog, resources, resource) {
            if (!typeof(resource.params.on_find) === 'function') {
                console.log(sprintf("Instance resource '%s' has no 'on_find' param.", 
                                    resource.name));
                os.exit(3);
            }
	    result = resource.params.on_find(catalog, resources);
	    if (Array.isArray(result) && result.length == 0) {
		return;
	    }
	    return result;
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
                if (mithras.verbose) {
                    log(sprintf("Deleting %d instances", found.length));
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
                    if (mithras.verbose) {
                        log(sprintf("Deleting %d instances", numToDelete));
                    }
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
                    if (mithras.verbose) {
                        log(sprintf("Launching %d instances", numToAdd));
                    }
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
                return [handler.findInCatalog(catalog, resources, resource), true];
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
            var params = resource.params;
            var found = handler.findInCatalog(catalog, resources, resource);
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
