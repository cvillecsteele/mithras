
## Object Format: AUTOSCALINGGROUP


    {
      "AutoScalingGroupARN": "arn:aws:autoscaling:us-east-1:286536233385:autoScalingGroup:c46af726-3bda-458f-90d9-e030b786a473:autoScalingGroupName/my-asg",
      "AutoScalingGroupName": "my-asg",
      "AvailabilityZones": [
        "us-east-1d"
      ],
      "CreatedTime": "2016-04-28T13:48:52.004Z",
      "DefaultCooldown": 100,
      "DesiredCapacity": 1,
      "EnabledMetrics": null,
      "HealthCheckGracePeriod": 100,
      "HealthCheckType": "EC2",
      "Instances": [
        {
          "AvailabilityZone": "us-east-1d",
          "HealthStatus": "Healthy",
          "InstanceId": "i-98cbd803",
          "LaunchConfigurationName": "lcName",
          "LifecycleState": "InService",
          "ProtectedFromScaleIn": true
        }
      ],
      "LaunchConfigurationName": "lcName",
      "LoadBalancerNames": null,
      "MaxSize": 2,
      "MinSize": 1,
      "NewInstancesProtectedFromScaleIn": true,
      "PlacementGroup": null,
      "Status": null,
      "SuspendedProcesses": null,
      "Tags": [
        {
          "Key": "Name",
          "PropagateAtLaunch": true,
          "ResourceId": "my-asg",
          "ResourceType": "auto-scaling-group",
          "Value": "test"
        }
      ],
      "TerminationPolicies": [
        "Default"
      ],
      "VPCZoneIdentifier": ""
    }
