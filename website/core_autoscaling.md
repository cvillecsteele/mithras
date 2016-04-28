 


 # CORE FUNCTIONS: AUTOSCALING


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.autoscaling.groups.scan](#gscan)
 > * [aws.autoscaling.groups.create](#gcreate)
 > * [aws.autoscaling.groups.delete](#gdelete)
 > * [aws.autoscaling.groups.describe](#gdescribe)

 > * [aws.autoscaling.hooks.scan](#hscan)
 > * [aws.autoscaling.hooks.create](#hcreate)
 > * [aws.autoscaling.hooks.delete](#hdelete)
 > * [aws.autoscaling.hooks.describe](#hdescribe)
 > * [aws.autoscaling.hooks.complete](#hcomplete)
 > * [aws.autoscaling.hooks.recordHeartbeat](#hrecord)

 > * [aws.autoscaling.launchConfigs.scan](#lscan)
 > * [aws.autoscaling.launchConfigs.create](#lcreate)
 > * [aws.autoscaling.launchConfigs.delete](#ldelete)
 > * [aws.autoscaling.launchConfigs.describe](#ldescribe)

 This API allows resource handlers to manage Autoscaling groups.

 ## AWS.AUTOSCALING.GROUPS.SCAN
 <a name="gscan"></a>
 `aws.autoscaling.groups.scan(region);`

 Returns a list of autoscaling groups.

 Example:

 ```

  var groups = aws.autoscaling.groups.scan("us-east-1");

 ```

 ## AWS.AUTOSCALING.GROUPS.DESCRIBE
 <a name="gdescribe"></a>
 `aws.autoscaling.groups.desribe(region, groupName);`

 Get info about an autoscaling group.

 Example:

 ```

  var group = aws.autoscaling.groups.describe("us-east-1", "groupName");

 ```

 ## AWS.AUTOSCALING.GROUPS.CREATE
 <a name="gcreate"></a>
 `aws.autoscaling.groups.create(region, config);`

 Create an autoscaling group.

 Example:

 ```

  var scaling =  aws.autoscaling.groups.create(
    "us-east-1",
 {
   AutoScalingGroupName: "name"
   MaxSize:              1
   MinSize:              1
   AvailabilityZones: [ "zone" ]
   DefaultCooldown:         100
   DesiredCapacity:         5
   HealthCheckGracePeriod:  100
   HealthCheckType:         "..."
   InstanceId:              "..."
   LaunchConfigurationName: "ResourceName"
   LoadBalancerNames: [ "..." ]
   NewInstancesProtectedFromScaleIn: true
   PlacementGroup:                   "groupName"
   Tags: [
     {
       Key:               "TagKey"
       PropagateAtLaunch: true
       ResourceId:        "id"
       ResourceType:      "type"
       Value:             "TagValue"
     },
   ]
   TerminationPolicies: [ "policyName" ]
   VPCZoneIdentifier: "zoneId"
 });

 ```

 ## AWS.AUTOSCALING.GROUPS.DELETE
 <a name="gdelete"></a>
 `aws.autoscaling.groups.delete(region, groupName);`

 Delete an autoscaling group.

 Example:

 ```

  aws.autoscaling.groups.delete("us-east-1", "groupName");

 ```

 ## AWS.AUTOSCALING.HOOKS.SCAN
 <a name="hscan"></a>
 `aws.autoscaling.hooks.scan(region, groupName);`

 Returns a list of autoscaling hooks.

 Example:

 ```

  var hooks = aws.autoscaling.hooks.scan("us-east-1", "groupName");

 ```

 ## AWS.AUTOSCALING.HOOKS.DESCRIBE
 <a name="hdescribe"></a>
 `aws.autoscaling.hooks.desribe(region, groupName, hookName);`

 Get info about an autoscaling hook.

 Example:

 ```

  var hook = aws.autoscaling.hooks.describe("us-east-1", "groupName", "hookName");

 ```

 ## AWS.AUTOSCALING.HOOKS.CREATE
 <a name="hcreate"></a>
 `aws.autoscaling.hooks.create(region, config);`

 Create an autoscaling hook.

 Example:

 ```

  var hook =  aws.autoscaling.hooks.create(
    "us-east-1",
 {
   AutoScalingGroupName:  "groupName"
   LifecycleHookName:     "hookName"
   DefaultResult:         "CONTINUE"
   HeartbeatTimeout:      100
   LifecycleTransition:   "autoscaling:EC2_INSTANCE_LAUNCHING"
   NotificationMetadata:  "data here"
   NotificationTargetARN: "arn:aws:sns:us-west-2:123456789012:my-sns-topic"
   RoleARN:               "arn:aws:iam::123456789012:role/my-auto-scaling-role"
 });

 ```

 ## AWS.AUTOSCALING.HOOKS.DELETE
 <a name="hdelete"></a>
 `aws.autoscaling.hooks.delete(region, gropuName, hookName);`

 Delete an autoscaling hook.

 Example:

 ```

  aws.autoscaling.hooks.delete("us-east-1", "groupName", "hookName");

 ```

 ## AWS.AUTOSCALING.HOOKS.COMPLETE
 <a name="hcomplete"></a>
 `aws.autoscaling.hooks.complete(region, config);`

 Complete a lifecycle hook.

 Example:

 ```

  aws.autoscaling.hooks.complete("us-east-1",
 {
   AutoScalingGroupName:  "groupName"
   LifecycleActionResult: "COMPLETE"
   LifecycleHookName:     "hookName"
   InstanceId:            "1234"
   LifecycleActionToken:  "LifecycleActionToken"
 });

 ```

 ## AWS.AUTOSCALING.HOOKS.RECORDHEARTBEAT
 <a name="hrecord"></a>
 `aws.autoscaling.hooks.recordHeartbeat(region, config);`

 Complete a lifecycle hook.

 Example:

 ```

  aws.autoscaling.hooks.recordHeartbeat("us-east-1",
 {
   AutoScalingGroupName:  "groupName"
   LifecycleHookName:     "hookName"
   InstanceId:            "1234"
   LifecycleActionToken:  "LifecycleActionToken"
 });

 ```

 ## AWS.AUTOSCALING.LAUNCHCONFIGS.SCAN
 <a name="lscan"></a>
 `aws.autoscaling.launchconfigs.scan(region);`

 Returns a list of autoscaling launchconfigs.

 Example:

 ```

  var launchconfigs = aws.autoscaling.launchConfigs.scan("us-east-1");

 ```

 ## AWS.AUTOSCALING.LAUNCHCONFIGS.DESCRIBE
 <a name="ldescribe"></a>
 `aws.autoscaling.launchconfigs.desribe(region, launchconfigName);`

 Get info about an autoscaling launchconfig.

 Example:

 ```

  var launchconfig = aws.autoscaling.launchconfigs.describe("us-east-1", "launchconfigName");

 ```

 ## AWS.AUTOSCALING.LAUNCHCONFIGS.CREATE
 <a name="lcreate"></a>
 `aws.autoscaling.launchConfigs.create(region, config);`

 Create an autoscaling launchconfig.

 Example:

 ```

  var config =  aws.autoscaling.launchConfigs.create(
    "us-east-1",
 {
     LaunchConfigurationName:  "lcName"
     AssociatePublicIpAddress: true
     BlockDeviceMappings: [
      {
         DeviceName: "string"
         Ebs: {
           DeleteOnTermination: true
           Encrypted:           false
           Iops:                100
           SnapshotId:          "StringMaxLen255"
           VolumeSize:          100
           VolumeType:          "BlockDeviceEbsVolumeType"
        }
        NoDevice:    false
        VirtualName: "StringMaxLen255"
      }
     ]
     ClassicLinkVPCId: "StringMaxLen255"
     ClassicLinkVPCSecurityGroups: [ "StringMaxLen255" ]
     EbsOptimized:       false
     IamInstanceProfile: "StringMaxLen1600"
     ImageId:            "StringMaxLen255"
     InstanceId:         "StringMaxLen19"
     InstanceMonitoring: {
       Enabled: true
     }
     InstanceType:     "StringMaxLen255"
     KernelId:         "StringMaxLen255"
     KeyName:          "StringMaxLen255"
     PlacementTenancy: "StringMaxLen64"
     RamdiskId:        "StringMaxLen255"
     SecurityGroups: [ "sg-1234" ]
     SpotPrice: "SpotPrice"
     UserData:  "StringUserData"
 });

 ```

 ## AWS.AUTOSCALING.LAUNCHCONFIGS.DELETE
 <a name="ldelete"></a>
 `aws.autoscaling.launchConfigs.delete(region, launchconfigName);`

 Delete an autoscaling launchconfig.

 Example:

 ```

  aws.autoscaling.launchConfigs.delete("us-east-1", "launchconfigName");

 ```


