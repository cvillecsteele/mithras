
## Object Format: Beanstalk Configuration

    {
      "ApplicationName": "test",
      "DateCreated": "2016-05-09T04:17:47Z",
      "DateUpdated": "2016-05-09T04:17:47Z",
      "DeploymentStatus": "deployed",
      "Description": "test",
      "EnvironmentName": "test",
      "OptionSettings": [
        {
          "Namespace": "aws:autoscaling:asg",
          "OptionName": "Availability Zones",
          "ResourceName": "AWSEBAutoScalingGroup",
          "Value": "Any"
        },
        {
          "Namespace": "aws:autoscaling:asg",
          "OptionName": "Cooldown",
          "ResourceName": "AWSEBAutoScalingGroup",
          "Value": "360"
        },
        {
          "Namespace": "aws:autoscaling:asg",
          "OptionName": "Custom Availability Zones",
          "ResourceName": "AWSEBAutoScalingGroup",
          "Value": ""
        },
        {
          "Namespace": "aws:autoscaling:asg",
          "OptionName": "MaxSize",
          "ResourceName": "AWSEBAutoScalingGroup",
          "Value": "4"
        },
        {
          "Namespace": "aws:autoscaling:asg",
          "OptionName": "MinSize",
          "ResourceName": "AWSEBAutoScalingGroup",
          "Value": "1"
        },
        {
          "Namespace": "aws:autoscaling:launchconfiguration",
          "OptionName": "BlockDeviceMappings",
          "ResourceName": "AWSEBAutoScalingLaunchConfiguration",
          "Value": null
        },
        {
          "Namespace": "aws:autoscaling:launchconfiguration",
          "OptionName": "EC2KeyName",
          "ResourceName": "AWSEBAutoScalingLaunchConfiguration",
          "Value": null
        },
        {
          "Namespace": "aws:autoscaling:launchconfiguration",
          "OptionName": "IamInstanceProfile",
          "ResourceName": "AWSEBAutoScalingLaunchConfiguration",
          "Value": null
        },
        {
          "Namespace": "aws:autoscaling:launchconfiguration",
          "OptionName": "ImageId",
          "ResourceName": "AWSEBAutoScalingLaunchConfiguration",
          "Value": "ami-7bbbb011"
        },
        {
          "Namespace": "aws:autoscaling:launchconfiguration",
          "OptionName": "InstanceType",
          "ResourceName": null,
          "Value": "t1.micro"
        },
        {
          "Namespace": "aws:autoscaling:launchconfiguration",
          "OptionName": "MonitoringInterval",
          "ResourceName": "AWSEBAutoScalingLaunchConfiguration",
          "Value": "5 minute"
        },
        {
          "Namespace": "aws:autoscaling:launchconfiguration",
          "OptionName": "RootVolumeIOPS",
          "ResourceName": "AWSEBAutoScalingLaunchConfiguration",
          "Value": null
        },
        {
          "Namespace": "aws:autoscaling:launchconfiguration",
          "OptionName": "RootVolumeSize",
          "ResourceName": "AWSEBAutoScalingLaunchConfiguration",
          "Value": null
        },
        {
          "Namespace": "aws:autoscaling:launchconfiguration",
          "OptionName": "RootVolumeType",
          "ResourceName": "AWSEBAutoScalingLaunchConfiguration",
          "Value": null
        },
        {
          "Namespace": "aws:autoscaling:launchconfiguration",
          "OptionName": "SSHSourceRestriction",
          "ResourceName": null,
          "Value": "tcp,22,22,0.0.0.0/0"
        },
        {
          "Namespace": "aws:autoscaling:launchconfiguration",
          "OptionName": "SecurityGroups",
          "ResourceName": "AWSEBAutoScalingLaunchConfiguration",
          "Value": "awseb-e-9zv3vythg2-stack-AWSEBSecurityGroup-UT80P94G03C2"
        },
        {
          "Namespace": "aws:autoscaling:trigger",
          "OptionName": "BreachDuration",
          "ResourceName": "AWSEBCloudwatchAlarmLow",
          "Value": "5"
        },
        {
          "Namespace": "aws:autoscaling:trigger",
          "OptionName": "EvaluationPeriods",
          "ResourceName": "AWSEBCloudwatchAlarmLow",
          "Value": "1"
        },
        {
          "Namespace": "aws:autoscaling:trigger",
          "OptionName": "LowerBreachScaleIncrement",
          "ResourceName": "AWSEBAutoScalingScaleDownPolicy",
          "Value": "-1"
        },
        {
          "Namespace": "aws:autoscaling:trigger",
          "OptionName": "LowerThreshold",
          "ResourceName": "AWSEBCloudwatchAlarmLow",
          "Value": "2000000"
        },
        {
          "Namespace": "aws:autoscaling:trigger",
          "OptionName": "MeasureName",
          "ResourceName": "AWSEBCloudwatchAlarmLow",
          "Value": "NetworkOut"
        },
        {
          "Namespace": "aws:autoscaling:trigger",
          "OptionName": "Period",
          "ResourceName": "AWSEBCloudwatchAlarmLow",
          "Value": "5"
        },
        {
          "Namespace": "aws:autoscaling:trigger",
          "OptionName": "Statistic",
          "ResourceName": "AWSEBCloudwatchAlarmLow",
          "Value": "Average"
        },
        {
          "Namespace": "aws:autoscaling:trigger",
          "OptionName": "Unit",
          "ResourceName": "AWSEBCloudwatchAlarmLow",
          "Value": "Bytes"
        },
        {
          "Namespace": "aws:autoscaling:trigger",
          "OptionName": "UpperBreachScaleIncrement",
          "ResourceName": "AWSEBAutoScalingScaleUpPolicy",
          "Value": "1"
        },
        {
          "Namespace": "aws:autoscaling:trigger",
          "OptionName": "UpperThreshold",
          "ResourceName": "AWSEBCloudwatchAlarmHigh",
          "Value": "6000000"
        },
        {
          "Namespace": "aws:autoscaling:updatepolicy:rollingupdate",
          "OptionName": "MaxBatchSize",
          "ResourceName": "AWSEBAutoScalingGroup",
          "Value": null
        },
        {
          "Namespace": "aws:autoscaling:updatepolicy:rollingupdate",
          "OptionName": "MinInstancesInService",
          "ResourceName": "AWSEBAutoScalingGroup",
          "Value": null
        },
        {
          "Namespace": "aws:autoscaling:updatepolicy:rollingupdate",
          "OptionName": "PauseTime",
          "ResourceName": "AWSEBAutoScalingGroup",
          "Value": null
        },
        {
          "Namespace": "aws:autoscaling:updatepolicy:rollingupdate",
          "OptionName": "RollingUpdateEnabled",
          "ResourceName": "AWSEBAutoScalingGroup",
          "Value": "false"
        },
        {
          "Namespace": "aws:autoscaling:updatepolicy:rollingupdate",
          "OptionName": "RollingUpdateType",
          "ResourceName": "AWSEBAutoScalingGroup",
          "Value": "Time"
        },
        {
          "Namespace": "aws:autoscaling:updatepolicy:rollingupdate",
          "OptionName": "Timeout",
          "ResourceName": "AWSEBAutoScalingGroup",
          "Value": "PT30M"
        },
        {
          "Namespace": "aws:cloudformation:template:parameter",
          "OptionName": "AppSource",
          "ResourceName": null,
          "Value": "https://elasticbeanstalk-samples-us-east-1.s3.amazonaws.com/elasticbeanstalk-golang-sampleapp.zip"
        },
        {
          "Namespace": "aws:cloudformation:template:parameter",
          "OptionName": "EnvironmentVariables",
          "ResourceName": null,
          "Value": "PATH=/bin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:/usr/local/go/bin:/var/app/current,APP_STAGING_DIR=/var/app/staging,APP_DEPLOY_DIR=/var/app/current,GOPATH=/var/app/current"
        },
        {
          "Namespace": "aws:cloudformation:template:parameter",
          "OptionName": "HooksPkgUrl",
          "ResourceName": null,
          "Value": "https://s3.amazonaws.com/elasticbeanstalk-env-resources-us-east-1/stalks/eb_golang_1.0.19.1/lib/hooks.tar.gz"
        },
        {
          "Namespace": "aws:cloudformation:template:parameter",
          "OptionName": "InstancePort",
          "ResourceName": null,
          "Value": "80"
        },
        {
          "Namespace": "aws:cloudformation:template:parameter",
          "OptionName": "StaticFiles",
          "ResourceName": null,
          "Value": null
        },
        {
          "Namespace": "aws:ec2:vpc",
          "OptionName": "AssociatePublicIpAddress",
          "ResourceName": "AWSEBAutoScalingLaunchConfiguration",
          "Value": null
        },
        {
          "Namespace": "aws:ec2:vpc",
          "OptionName": "ELBScheme",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "public"
        },
        {
          "Namespace": "aws:ec2:vpc",
          "OptionName": "ELBSubnets",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": null
        },
        {
          "Namespace": "aws:ec2:vpc",
          "OptionName": "Subnets",
          "ResourceName": "AWSEBAutoScalingGroup",
          "Value": null
        },
        {
          "Namespace": "aws:ec2:vpc",
          "OptionName": "VPCId",
          "ResourceName": "AWSEBSecurityGroup",
          "Value": null
        },
        {
          "Namespace": "aws:elasticbeanstalk:application",
          "OptionName": "Application Healthcheck URL",
          "ResourceName": null,
          "Value": ""
        },
        {
          "Namespace": "aws:elasticbeanstalk:application:environment",
          "OptionName": "APP_DEPLOY_DIR",
          "ResourceName": null,
          "Value": "/var/app/current"
        },
        {
          "Namespace": "aws:elasticbeanstalk:application:environment",
          "OptionName": "APP_STAGING_DIR",
          "ResourceName": null,
          "Value": "/var/app/staging"
        },
        {
          "Namespace": "aws:elasticbeanstalk:application:environment",
          "OptionName": "GOPATH",
          "ResourceName": null,
          "Value": "/var/app/current"
        },
        {
          "Namespace": "aws:elasticbeanstalk:application:environment",
          "OptionName": "PATH",
          "ResourceName": null,
          "Value": "/bin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:/usr/local/go/bin:/var/app/current"
        },
        {
          "Namespace": "aws:elasticbeanstalk:command",
          "OptionName": "BatchSize",
          "ResourceName": null,
          "Value": "100"
        },
        {
          "Namespace": "aws:elasticbeanstalk:command",
          "OptionName": "BatchSizeType",
          "ResourceName": null,
          "Value": "Percentage"
        },
        {
          "Namespace": "aws:elasticbeanstalk:command",
          "OptionName": "IgnoreHealthCheck",
          "ResourceName": null,
          "Value": "false"
        },
        {
          "Namespace": "aws:elasticbeanstalk:command",
          "OptionName": "Timeout",
          "ResourceName": null,
          "Value": "600"
        },
        {
          "Namespace": "aws:elasticbeanstalk:control",
          "OptionName": "DefaultSSHPort",
          "ResourceName": null,
          "Value": "22"
        },
        {
          "Namespace": "aws:elasticbeanstalk:control",
          "OptionName": "LaunchTimeout",
          "ResourceName": null,
          "Value": "0"
        },
        {
          "Namespace": "aws:elasticbeanstalk:control",
          "OptionName": "LaunchType",
          "ResourceName": null,
          "Value": "Migration"
        },
        {
          "Namespace": "aws:elasticbeanstalk:control",
          "OptionName": "RollbackLaunchOnFailure",
          "ResourceName": null,
          "Value": "false"
        },
        {
          "Namespace": "aws:elasticbeanstalk:environment",
          "OptionName": "EnvironmentType",
          "ResourceName": null,
          "Value": "LoadBalanced"
        },
        {
          "Namespace": "aws:elasticbeanstalk:environment",
          "OptionName": "ExternalExtensionsS3Bucket",
          "ResourceName": null,
          "Value": null
        },
        {
          "Namespace": "aws:elasticbeanstalk:environment",
          "OptionName": "ExternalExtensionsS3Key",
          "ResourceName": null,
          "Value": null
        },
        {
          "Namespace": "aws:elasticbeanstalk:environment",
          "OptionName": "ServiceRole",
          "ResourceName": null,
          "Value": null
        },
        {
          "Namespace": "aws:elasticbeanstalk:healthreporting:system",
          "OptionName": "ConfigDocument",
          "ResourceName": null,
          "Value": "{\"Version\":1,\"CloudWatchMetrics\":{\"Instance\":{\"CPUIrq\":null,\"LoadAverage5min\":null,\"ApplicationRequests5xx\":null,\"ApplicationRequests4xx\":null,\"CPUUser\":null,\"LoadAverage1min\":null,\"ApplicationLatencyP50\":null,\"CPUIdle\":null,\"InstanceHealth\":null,\"ApplicationLatencyP95\":null,\"ApplicationLatencyP85\":null,\"RootFilesystemUtil\":null,\"ApplicationLatencyP90\":null,\"CPUSystem\":null,\"ApplicationLatencyP75\":null,\"CPUSoftirq\":null,\"ApplicationLatencyP10\":null,\"ApplicationLatencyP99\":null,\"ApplicationRequestsTotal\":null,\"ApplicationLatencyP99.9\":null,\"ApplicationRequests3xx\":null,\"ApplicationRequests2xx\":null,\"CPUIowait\":null,\"CPUNice\":null},\"Environment\":{\"InstancesSevere\":null,\"InstancesDegraded\":null,\"ApplicationRequests5xx\":null,\"ApplicationRequests4xx\":null,\"ApplicationLatencyP50\":null,\"ApplicationLatencyP95\":null,\"ApplicationLatencyP85\":null,\"InstancesUnknown\":null,\"ApplicationLatencyP90\":null,\"InstancesInfo\":null,\"InstancesPending\":null,\"ApplicationLatencyP75\":null,\"ApplicationLatencyP10\":null,\"ApplicationLatencyP99\":null,\"ApplicationRequestsTotal\":null,\"InstancesNoData\":null,\"ApplicationLatencyP99.9\":null,\"ApplicationRequests3xx\":null,\"ApplicationRequests2xx\":null,\"InstancesOk\":null,\"InstancesWarning\":null}}}"
        },
        {
          "Namespace": "aws:elasticbeanstalk:healthreporting:system",
          "OptionName": "HealthCheckSuccessThreshold",
          "ResourceName": null,
          "Value": "Ok"
        },
        {
          "Namespace": "aws:elasticbeanstalk:healthreporting:system",
          "OptionName": "SystemType",
          "ResourceName": null,
          "Value": "basic"
        },
        {
          "Namespace": "aws:elasticbeanstalk:hostmanager",
          "OptionName": "LogPublicationControl",
          "ResourceName": null,
          "Value": "false"
        },
        {
          "Namespace": "aws:elasticbeanstalk:managedactions",
          "OptionName": "ManagedActionsEnabled",
          "ResourceName": null,
          "Value": "false"
        },
        {
          "Namespace": "aws:elasticbeanstalk:managedactions",
          "OptionName": "PreferredStartTime",
          "ResourceName": null,
          "Value": null
        },
        {
          "Namespace": "aws:elasticbeanstalk:managedactions:platformupdate",
          "OptionName": "InstanceRefreshEnabled",
          "ResourceName": null,
          "Value": "false"
        },
        {
          "Namespace": "aws:elasticbeanstalk:managedactions:platformupdate",
          "OptionName": "UpdateLevel",
          "ResourceName": null,
          "Value": null
        },
        {
          "Namespace": "aws:elasticbeanstalk:monitoring",
          "OptionName": "Automatically Terminate Unhealthy Instances",
          "ResourceName": null,
          "Value": "true"
        },
        {
          "Namespace": "aws:elasticbeanstalk:sns:topics",
          "OptionName": "Notification Endpoint",
          "ResourceName": null,
          "Value": null
        },
        {
          "Namespace": "aws:elasticbeanstalk:sns:topics",
          "OptionName": "Notification Protocol",
          "ResourceName": null,
          "Value": "email"
        },
        {
          "Namespace": "aws:elasticbeanstalk:sns:topics",
          "OptionName": "Notification Topic ARN",
          "ResourceName": null,
          "Value": null
        },
        {
          "Namespace": "aws:elasticbeanstalk:sns:topics",
          "OptionName": "Notification Topic Name",
          "ResourceName": null,
          "Value": null
        },
        {
          "Namespace": "aws:elb:healthcheck",
          "OptionName": "HealthyThreshold",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "3"
        },
        {
          "Namespace": "aws:elb:healthcheck",
          "OptionName": "Interval",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "10"
        },
        {
          "Namespace": "aws:elb:healthcheck",
          "OptionName": "Target",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "TCP:80"
        },
        {
          "Namespace": "aws:elb:healthcheck",
          "OptionName": "Timeout",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "5"
        },
        {
          "Namespace": "aws:elb:healthcheck",
          "OptionName": "UnhealthyThreshold",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "5"
        },
        {
          "Namespace": "aws:elb:listener:80",
          "OptionName": "InstancePort",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "80"
        },
        {
          "Namespace": "aws:elb:listener:80",
          "OptionName": "InstanceProtocol",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "HTTP"
        },
        {
          "Namespace": "aws:elb:listener:80",
          "OptionName": "ListenerEnabled",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "true"
        },
        {
          "Namespace": "aws:elb:listener:80",
          "OptionName": "ListenerProtocol",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "HTTP"
        },
        {
          "Namespace": "aws:elb:listener:80",
          "OptionName": "PolicyNames",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": null
        },
        {
          "Namespace": "aws:elb:listener:80",
          "OptionName": "SSLCertificateId",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": null
        },
        {
          "Namespace": "aws:elb:loadbalancer",
          "OptionName": "CrossZone",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "false"
        },
        {
          "Namespace": "aws:elb:loadbalancer",
          "OptionName": "LoadBalancerHTTPPort",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "80"
        },
        {
          "Namespace": "aws:elb:loadbalancer",
          "OptionName": "LoadBalancerHTTPSPort",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "OFF"
        },
        {
          "Namespace": "aws:elb:loadbalancer",
          "OptionName": "LoadBalancerPortProtocol",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "HTTP"
        },
        {
          "Namespace": "aws:elb:loadbalancer",
          "OptionName": "LoadBalancerSSLPortProtocol",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "HTTPS"
        },
        {
          "Namespace": "aws:elb:loadbalancer",
          "OptionName": "SSLCertificateId",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": null
        },
        {
          "Namespace": "aws:elb:loadbalancer",
          "OptionName": "SecurityGroups",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "{\"Fn::GetAtt\":[\"AWSEBLoadBalancerSecurityGroup\",\"GroupId\"]},{\"Ref\":\"AWSEBLoadBalancerSecurityGroup\"}"
        },
        {
          "Namespace": "aws:elb:policies",
          "OptionName": "ConnectionDrainingEnabled",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "false"
        },
        {
          "Namespace": "aws:elb:policies",
          "OptionName": "ConnectionDrainingTimeout",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "20"
        },
        {
          "Namespace": "aws:elb:policies",
          "OptionName": "ConnectionSettingIdleTimeout",
          "ResourceName": "AWSEBLoadBalancer",
          "Value": "60"
        }
      ],
      "SolutionStackName": "64bit Amazon Linux 2016.03 v2.1.0 running Go 1.4",
      "TemplateName": null
    }
