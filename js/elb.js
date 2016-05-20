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
// # ELB
// 
// Elb is resource handler for managing AWS elastic load balancers.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"elb"`
// 
// Usage:
// 
// `var elb = require("elb").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var lbName = "my-lb";
// var rElb = {
//      name: "elb"
//      module: "elb"
//      dependsOn: [otherResource.name]
//      on_delete: function(elb) { 
//          // Sometimes aws takes a bit to delete an elb, and we can't
//          // proceed with deleting until it's GONE.
//          this.delay = 30; 
//          return true;
//      }
//      params: {
//          region: "us-east-1"
//          ensure: "present"
//          elb: {
//              Listeners: [
//                  {
//                      InstancePort:     80
//                      LoadBalancerPort: 80
//                      Protocol:         "http"
//                      InstanceProtocol: "http"
//                  },
//              ]
//              LoadBalancerName: lbName
//              SecurityGroups: [ "sg-123" ]
//              Subnets: [
//                  "subnet-abc"
//                  "subnet-def"
//              ]
//              Tags: [
//                  {
//                      Key:   "foo"
//                      Value: "bar"
//                  },
//              ]
//          }
//          attributes: {
//              LoadBalancerAttributes: {
//                  AccessLog: {
//                      Enabled:        false
//                      EmitInterval:   60
//                      S3BucketName:   "my-loadbalancer-logs"
//                      S3BucketPrefix: "test-app"
//                  }
//                  ConnectionDraining: {
//                      Enabled: true
//                      Timeout: 300
//                  }
//                  ConnectionSettings: {
//                      IdleTimeout: 30
//                  }
//                  CrossZoneLoadBalancing: {
//                      Enabled: true
//                  }
//              }
//              LoadBalancerName: lbName
//          }
//          health: {
//              HealthCheck: {
//                  HealthyThreshold:   2
//                  Interval:           30
//                  Target:             "HTTP:80/hc"
//                  Timeout:            5
//                  UnhealthyThreshold: 3
//              }
//              LoadBalancerName: lbName
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
// ### `elb`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elb.html#type-CreateLoadBalanerInput)
//
// Parameters for resource creation.
// 
// ### `on_delete`
//
// * Required: false
// * Allowed Values: `function(elb) { ... }`
//
// Called after resource deletion.  May be used to modify `wait`
// property (see example above).
// 
// ### `attributes`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elb.html#type-ModifyLoadBalancerAttributesInput)
//
// If specified, used to apply attributes to a created ELB.
// 
// ### `health`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elb.html#type-ConfigureHealthCheckInput)
//
// If specified, used to specify health check params for a created ELB.
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
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
        moduleNames: ["elb", "elbMembership"]
        findInCatalog: function(catalog, resource, elbName) {
            if (typeof(resource.params.on_find) === 'function') {
		result = resource.params.on_find(catalog, resource);
		if (!result || 
		    (Array.isArray(result) && result.length == 0)) {
		    return;
		}
		return result;
	    }
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
                    if (mithras.verbose) {
                        log(sprintf("Deleting ELB '%s'", 
                                    params.elb.LoadBalancerName));
                    }
                    aws.elbs.delete(params.region, params.elb.LoadBalancerName);

                    break;
                case "present":
                    if (elb) {
                        log(sprintf("ELB found, no action taken."));
                        return [elb, true];
                        break;
                    }
                    
                    // create 
                    if (mithras.verbose) {
                        log(sprintf("Creating elb '%s'", 
                                    params.elb.LoadBalancerName));
                    }
                    created = aws.elbs.create(params.region, params.elb);
                    
                    // Set health
                    if (mithras.verbose) {
                        log(sprintf("Setting health for elb '%s'", 
                                    params.elb.LoadBalancerName));
                    }
                    aws.elbs.setHealth(params.region, 
                                       params.elb.LoadBalancerName, 
                                       params.health);
                    
                    // Set attrs
                    if (mithras.verbose) {
                        log(sprintf("Setting attributes for elb '%s'", 
                                    params.elb.LoadBalancerName));
                    }
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
						resource,
                                                params.membership.LoadBalancerName);

                if (ensure === "present") {
                    if (!elb) { 
                        console.log(sprintf("Can't find elb '%s' in catalog", 
                                            params.membership.LoadBalancerName));
                        os.exit(3);
                    }
                } else if (ensure === "absent") {
                    if (!elb) {
                        return [null, true];
                    }
                }

                // intersection of input and LB
                lbInstanceIds = elb.Instances ? _.pluck(elb.Instances, "InstanceId") : [];
                inInstanceIds = [];
                if (_.isObject(params.membership.Instances) &&
                    !_.isFunction(params.membership.Instances) &&
                    _.isArray(params.membership.Instances)) {
                    inInstanceIds = _.pluck(params.membership.Instances, "InstanceId");
                }

                var inBoth = _.intersection(lbInstanceIds, inInstanceIds);
                // What's in the LB minus what's in input
                var inLB = _.difference(lbInstanceIds, inInstanceIds);
                // inputs minus what's in the LB
                var inInput = _.difference(inInstanceIds, lbInstanceIds);

                // build map id -> instance
                var byIds = _.reduce(elb.Instances ? elb.Instances : [], function(memo, i) { 
                    memo[i.InstanceId] = i;
                    return memo;
                }, {});
                byIds = _.reduce(params.membership.Instances ? params.membership.Instances : [], function(memo, i) { 
                    if (i && i.InstanceId) {
                        memo[i.InstanceId] = i;
                    }
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
                        if (mithras.verbose) {
                            log(sprintf("Deregistering %d instances", 
                                        inInput.length));
                        }
                        console.log(JSON.stringify(inInput, null, 2));
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
                        if (mithras.verbose) {
                            log(sprintf("Registering %d instances", 
                                        inInput.length));
                        }
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
                        if (mithras.verbose) {
                            log(sprintf("Registering %d instances", 
                                        inInput.length));
                        }
                        aws.elbs.register(params.region, 
                                          params.membership.LoadBalancerName, 
                                          inInput);
                    }
                    if (inLB.length > 0) {
                        if (mithras.verbose) {
                            log(sprintf("Deregistering %d instances", 
                                        inLB.length));
                        }
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
        preflight: function(catalog, resources, resource) {
            if (!_.find(handler.moduleNames, function(m) { 
                return resource.module === m; 
            })) {
                return [null, false];
            }
            if (resource.module === handler.moduleNames[0]) {
                var params = resource.params;
                var elb = params.elb;
                var t = handler.findInCatalog(catalog, resource, elb.LoadBalancerName);
                if (t) {
                    return [t, true];
                }
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
