
## Object Format: AUTOSCALINGHOOK


    {
      "AutoScalingGroupName": "my-asg",
      "DefaultResult": "CONTINUE",
      "GlobalTimeout": 10000,
      "HeartbeatTimeout": 100,
      "LifecycleHookName": "hookName",
      "LifecycleTransition": "autoscaling:EC2_INSTANCE_LAUNCHING",
      "NotificationMetadata": null,
      "NotificationTargetARN": "arn:aws:sqs:us-east-1:286536233385:myqueue",
      "RoleARN": "arn:aws:iam::286536233385:role/test-asg-iam-role"
    }
