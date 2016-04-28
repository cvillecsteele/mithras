 
 
 # Autoscaling
 
 Autoscaling is resource handler for managing AWS autoscaling groups.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"autoscaling"`
 
 Usage:
 
 `var autoscaling = require("autoscaling").init();`
 
  ## Example Resource
 
 ```javascript
 var asg = {
      name: "asg"
      module: "autoscaling"
      dependsOn: [resource.name]
      params: {
          region: "us-east-1"
          ensure: "present"
          group: {
            AutoScalingGroupName: "name"
            MaxSize:              2
            MinSize:              1
            AvailabilityZones: [ "us-east-1" ]
            DefaultCooldown:         100
            DesiredCapacity:         1
            HealthCheckGracePeriod:  100
            LaunchConfigurationName: "lcName"
            NewInstancesProtectedFromScaleIn: true
            Tags: [
                {
                    Key:               "Name"
                    PropagateAtLaunch: true
                    ResourceId:        "name"
                    ResourceType:      "auto-scaling-group"
                    Value:             "test"
                },
            ]
          }
          hook: {
 	      AutoScalingGroupName:  "name"
 	      LifecycleHookName:     "hookName"
 	      DefaultResult:         "CONTINUE"
 	      HeartbeatTimeout:      100
 	      LifecycleTransition:   "autoscaling:EC2_INSTANCE_LAUNCHING"
 	      NotificationTargetARN: sqsArn // Like "arn:aws:sqs:us-east-1:286536233385:myqueue"
 	      RoleARN:               roleArn // Like "arn:aws:iam::286536233385:role/asg"
 	    }
          launchConfig: {
            LaunchConfigurationName:  "lcName"
            EbsOptimized:       false
            ImageId:            ami
            InstanceMonitoring: {
                Enabled: false
            }
            InstanceType:     instanceType
          }
      } // params
 };
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"`, the db specified by `db` will be created, and
 if `"absent"`, it will be removed using the `delete` property.
 
 ### `region`

 * Required: true
 * Allowed Values: string, any valid AWS region; eg "us-east-1"

 The region for calls to the AWS API.
 
 ### `group`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/scaling.html#type-CreateAutoScalingGroupInput)

 Parameters for resource creation.  If present, an autoscaling group is created/deleted.
 
 ### `launchConfig`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/scaling.html#type-CreateLaunchConfigurationInput)

 Parameters for resource creation.  If present, an autoscaling launch configuration is created/deleted.
 
 ### `hook`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/scaling.html#type-PutLifecycleHookInput)

 Parameters for resource creation.  If present, an autoscaling lifecycle hook is created/deleted.
 
 ### `on_find`

 * Required: true
 * Allowed Values: A function taking two parameters: `catalog` and `resource`

 If defined in the resource's `params` object, the `on_find`
 function provides a way for a matching resource to be identified
 using a user-defined way.  The function is called with the current
 `catalog`, as well as the `resource` object itself.  The function
 can look through the catalog, find a matching object using whatever
 logic you want, and return it.  If the function returns `undefined`
 or a n empty Javascript array, (`[]`), the function is indicating
 that no matching resource was found in the `catalog`.
 
 ### `tags`

 * Required: false
 * Allowed Values: A map of tags to be applied to created instances

 A map of tags to be applied to created instances
 

