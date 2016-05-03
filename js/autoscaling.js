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
// # Autoscaling
// 
// Autoscaling is resource handler for managing AWS autoscaling groups.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"autoscaling"`
// 
// Usage:
// 
// `var autoscaling = require("autoscaling").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var asg = {
//      name: "asg"
//      module: "autoscaling"
//      dependsOn: [resource.name]
//      params: {
//          region: "us-east-1"
//          ensure: "present"
//          group: {
//            AutoScalingGroupName: "name"
//            MaxSize:              2
//            MinSize:              1
//            AvailabilityZones: [ "us-east-1" ]
//            DefaultCooldown:         100
//            DesiredCapacity:         1
//            HealthCheckGracePeriod:  100
//            LaunchConfigurationName: "lcName"
//            NewInstancesProtectedFromScaleIn: true
//            Tags: [
//                {
//                    Key:               "Name"
//                    PropagateAtLaunch: true
//                    ResourceId:        "name"
//                    ResourceType:      "auto-scaling-group"
//                    Value:             "test"
//                },
//            ]
//          }
//          hook: {
// 	      AutoScalingGroupName:  "name"
// 	      LifecycleHookName:     "hookName"
// 	      DefaultResult:         "CONTINUE"
// 	      HeartbeatTimeout:      100
// 	      LifecycleTransition:   "autoscaling:EC2_INSTANCE_LAUNCHING"
// 	      NotificationTargetARN: sqsArn // Like "arn:aws:sqs:us-east-1:286536233385:myqueue"
// 	      RoleARN:               roleArn // Like "arn:aws:iam::286536233385:role/asg"
// 	    }
//          launchConfig: {
//            LaunchConfigurationName:  "lcName"
//            EbsOptimized:       false
//            ImageId:            ami
//            InstanceMonitoring: {
//                Enabled: false
//            }
//            InstanceType:     instanceType
//          }
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
// ### `group`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/scaling.html#type-CreateAutoScalingGroupInput)
//
// Parameters for resource creation.  If present, an autoscaling group is created/deleted.
// 
// ### `launchConfig`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/scaling.html#type-CreateLaunchConfigurationInput)
//
// Parameters for resource creation.  If present, an autoscaling launch configuration is created/deleted.
// 
// ### `hook`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/scaling.html#type-PutLifecycleHookInput)
//
// Parameters for resource creation.  If present, an autoscaling lifecycle hook is created/deleted.
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
        moduleNames: ["autoscaling"]
        findInCatalog: function(catalog, resources, resource) {
            if (typeof(resource.params.on_find) === 'function') {
		result = resource.params.on_find(catalog, resource);
		if (!result || 
		    (Array.isArray(result) && result.length == 0)) {
		    return;
		}
		return result;
	    }
	    var group = null;
	    var config = null;
	    var hook = null;
	    if (resource.params.group) {
		group = _.find(catalog.autoscalingGroups, function(g) { 
                    return g.AutoScalingGroupName === 
			resource.params.group.AutoScalingGroupName;
		});
	    }
	    if (resource.params.hook) {
		hook = _.find(catalog.autoscalingHooks, function(h) { 
                    return h.LifecycleHookName === 
			resource.params.hook.LifecycleHookName;
		});
	    }
	    if (resource.params.launchConfig) {
		config = _.find(catalog.autoscalingLaunchConfigs, function(c) { 
                    return c.LaunchConfigurationName === 
			resource.params.launchConfig.LaunchConfigurationName;
		});
	    }
	    if (group || config || hook) {
		return {group: group, hook: hook, config: config};
	    }
	    return;
        }
        handle: function(catalog, resources, resource) {
            if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
                return [null, false];
            }
            
            var ensure = resource.params.ensure;
            var params = resource.params;
            
            if (!resource.params.group) {
                console.log("Invalid autoscaling resource params")
                os.exit(3);
            }
            var found = resource._target || {};
	    var group = found.group;
	    var hook = found.hook;
	    var config = found.config;
	    var name = group ? group.AutoScalingGroupName : params.group.AutoScalingGroupName;
	    
            if (mithras.verbose && group) {
                log(sprintf("Found ASG '%s'", name));
            }
            
            switch(ensure) {
            case "absent":
		
		// Group
		if (group) {
                    if (mithras.verbose) {
			log(sprintf("Deleting ASG"));
                    }
                    aws.autoscaling.groups.delete(params.region,  name);
                    catalog.autoscalingGroups = 
			_.reject(catalog.autoscalingGroups,
				 function(g) { 
                                     return g.AutoScalingGroupName = name;
				 });
		} else {
                    if (mithras.verbose) {
			log(sprintf("No action taken on group."));
                    }
		}

		Config
		if (config) {
                    if (mithras.verbose) {
			log(sprintf("Deleting ASG launch config"));
                    }
		    var cName = params.launchConfig.LaunchConfigurationName;
                    aws.autoscaling.launchConfigs.delete(params.region, cName);
                    catalog.autoscalingLaunchConfigs = 
			_.reject(catalog.autoscalingLaunchConfigs,
				 function(c) { 
                                     return c.LaunchConfigurationName = cName;
				 });
		} else {
                    if (mithras.verbose) {
			log(sprintf("No action taken on launch config."));
                    }
		}

		// Hook
		if (hook) {
                    if (mithras.verbose) {
			log(sprintf("Deleting ASG lifecycle hook"));
                    }
		    var gName = hook.AutoScalingGroupName;
		    var hName = params.hook.LifecycleHookName;
                    aws.autoscaling.hooks.delete(params.region, gName, hName);
                    catalog.autoscaling.hooks = 
			_.reject(catalog.autoscalingHooks,
				 function(h) { 
                                     return h.LifeCycleHookName = hName;
				 });
		} else {
                    if (mithras.verbose) {
			log(sprintf("No action taken on lifecycle hook."));
                    }
		}
                break;
            case "present":
		// Config
		if (!config) {
		    if (params.launchConfig) {
			var cName = params.launchConfig.LaunchConfigurationName;
			if (mithras.verbose) {
			    log(sprintf("Creating ASG launch config"));
			}
			aws.autoscaling.launchConfigs.create(params.region,
							     params.launchConfig);
			config = aws.autoscaling.launchConfigs.describe(params.region, 
									cName);
			catalog.autoscalingLaunchConfigs.push(config);
		    }
		} else {
                    if (mithras.verbose) {
			log(sprintf("No action taken on launch config."));
                    }
		}

		// Group
		if (!group) {
                    if (mithras.verbose) {
			log(sprintf("Creating ASG '%s'", name));
                    }
                    aws.autoscaling.groups.create(params.region, params.group);
		    group = aws.autoscaling.groups.describe(params.region, name);
                    catalog.autoscalingGroups.push(group);
		} else {
                    if (mithras.verbose) {
			log(sprintf("No action taken on group."));
                    }
		}
		
		// Hook
		if (!hook) {
		    if (params.hook) {
			var hName = params.hook.LifecycleHookName;
			if (mithras.verbose) {
			    log(sprintf("Creating ASG lifecycle hook"));
			}
			aws.autoscaling.hooks.create(params.region, params.hook);
			hook = aws.autoscaling.hooks.describe(params.region, hName);
			catalog.autoscalingHooks.push(config);
		    }
		} else {
                    if (mithras.verbose) {
			log(sprintf("No action taken on lifecycle hook."));
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
