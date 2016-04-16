 


 # CORE FUNCTIONS: ELB


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.elbs.create](#create)
 > * [aws.elbs.delete](#delete)
 > * [aws.elbs.describe](#describe)
 > * [aws.elbs.scan](#scan)
 > * [aws.elbs.register](#register)
 > * [aws.elbs.deRegister](#deregister)
 > * [aws.elbs.setHealth](#health)
 > * [aws.elbs.setAttrs](#attrs)

 This API allows exposes functions to manage AWS elastic load balancers.

 ## AWS.ELBS.CREATE
 <a name="create"></a>
 `aws.elbs.create(region, config);`

 Create an ELB

 Example:

 ```

  var lb = aws.elbs.create("us-east-1", {
 Listeners: [
  {
    InstancePort:     80
    LoadBalancerPort: 80
    Protocol:         "http"
    InstanceProtocol: "http"
  },
 ]
 LoadBalancerName: "test-lb"
 SecurityGroups: [
     "sg-xyz"
 ]
 Subnets: [
     "subnet-123"
     "subnet-456"
 ]
 Tags: [
     {
       Key:   "foo"
       Value: "bar"
     },
 ]
 });

 ```

 ## AWS.ELBS.SETATTRS
 <a name="attrs"></a>
 `aws.elbs.setAttrs(region, lbName, config);`

 Set ELB attributes

 Example:

 ```

  aws.elbs.setAttrs("us-east-1", "test-lb"{
  LoadBalancerAttributes: {
      AccessLog: {
        Enabled:        false
        EmitInterval:   60
        S3BucketName:   "my-loadbalancer-logs"
        S3BucketPrefix: "test-app"
      }
      ConnectionDraining: {
        Enabled: true
        Timeout: 300
      }
      ConnectionSettings: {
        IdleTimeout: 30
      }
      CrossZoneLoadBalancing: {
        Enabled: true
      }
    }
    LoadBalancerName: "test-lb"
  }
 });

 ```

 ## AWS.ELBS.SETHEALTH
 <a name="health"></a>
 `aws.elbs.setHealth(region, lbName, config);`

 Set ELB health check.

 Example:

 ```

  aws.elbs.setHealth("us-east-1", "test-lb", {
  HealthCheck: {
      HealthyThreshold:   2
      Interval:           30
      Target:             "HTTP:80/hc"
      Timeout:            5
      UnhealthyThreshold: 3
  }
  LoadBalancerName: "test-lb"
 });

 ```

 ## AWS.ELBS.DELETE
 <a name="delete"></a>
 `aws.elbs.delete(region, lbName);`

 Delete an ELB.

 Example:

 ```

  aws.elbs.delete("us-east-1", "test-lb");

 ```

 ## AWS.ELBS.DESCRIBE
 <a name="describe"></a>
 `aws.elbs.describe(region, lbName);`

 Get info about an ELB.

 Example:

 ```

  var elb = aws.elbs.describe("us-east-1", "test-lb");

 ```

 ## AWS.ELBS.SCAN
 <a name="scan"></a>
 `aws.elbs.scan(region);`

 Get info about ELBs.

 Example:

 ```

  var elbs = aws.elbs.scan("us-east-1");

 ```

 ## AWS.ELBS.REGISTER
 <a name="register"></a>
 `aws.elbs.register(region, lbName, instance);`

 Register an EC2 instance with an ELB.

 Example:

 ```

 aws.elbs.register("us-east-1", "test-lb",
 		  {
 		      "AmiLaunchIndex": 0,
 		      "Architecture": "x86_64",
 		      "BlockDeviceMappings": [
 			  {
 			      "DeviceName": "/dev/xvda",
 			      "Ebs": {
 				  "AttachTime": "2016-03-31T19:17:44Z",
 				  "DeleteOnTermination": true,
 				  "Status": "attached",
 				  "VolumeId": "vol-3d1ab09f"
 			      }
 			  }
 		      ],
 		      "ClientToken": "",
 		      "EbsOptimized": false,
 		      "Hypervisor": "xen",
 		      "IamInstanceProfile": {
 			  "Arn": "arn:aws:iam::286536233385:instance-profile/cr-webserver",
 			  "Id": "AIPAIOOUMBIRCV2QCTIYA"
 		      },
 		      "ImageId": "ami-60b6c60a",
 		      "InstanceId": "i-824ad119",
 		      "InstanceLifecycle": null,
 		      "InstanceType": "t2.small",
 		      "KernelId": null,
 		      "KeyName": "cr",
 		      "LaunchTime": "2016-03-31T19:17:43Z",
 		      "Monitoring": {
 			  "State": "enabled"
 		      },
 		      "NetworkInterfaces": [
 			  {
 			      "Association": {
 				  "IpOwnerId": "amazon",
 				  "PublicDnsName": "ec2-52-90-244-101.compute-1.amazonaws.com",
 				  "PublicIp": "52.90.244.101"
 			      },
 			      "Attachment": {
 				  "AttachTime": "2016-03-31T19:17:43Z",
 				  "AttachmentId": "eni-attach-4d1c8ab1",
 				  "DeleteOnTermination": true,
 				  "DeviceIndex": 0,
 				  "Status": "attached"
 			      },
 			      "Description": "",
 			      "Groups": [
 				  {
 				      "GroupId": "sg-57c9812f",
 				      "GroupName": "webserver"
 				  }
 			      ],
 			      "MacAddress": "0e:6b:fe:c1:cb:45",
 			      "NetworkInterfaceId": "eni-a6b644f6",
 			      "OwnerId": "286536233385",
 			      "PrivateDnsName": "ip-172-33-1-178.ec2.internal",
 			      "PrivateIpAddress": "172.33.1.178",
 			      "PrivateIpAddresses": [
 				  {
 				      "Association": {
 					  "IpOwnerId": "amazon",
 					  "PublicDnsName": "ec2-52-90-244-101.compute-1.amazonaws.com",
 					  "PublicIp": "52.90.244.101"
 				      },
 				      "Primary": true,
 				      "PrivateDnsName": "ip-172-33-1-178.ec2.internal",
 				      "PrivateIpAddress": "172.33.1.178"
 				  }
 			      ],
 			      "SourceDestCheck": true,
 			      "Status": "in-use",
 			      "SubnetId": "subnet-d76ae58f",
 			      "VpcId": "vpc-b88883dc"
 			  }
 		      ],
 		      "Placement": {
 			  "Affinity": null,
 			  "AvailabilityZone": "us-east-1d",
 			  "GroupName": "",
 			  "HostId": null,
 			  "Tenancy": "default"
 		      },
 		      "Platform": null,
 		      "PrivateDnsName": "ip-172-33-1-178.ec2.internal",
 		      "PrivateIpAddress": "172.33.1.178",
 		      "ProductCodes": null,
 		      "PublicDnsName": "ec2-52-90-244-101.compute-1.amazonaws.com",
 		      "PublicIpAddress": "52.90.244.101",
 		      "RamdiskId": null,
 		      "RootDeviceName": "/dev/xvda",
 		      "RootDeviceType": "ebs",
 		      "SecurityGroups": [
 			  {
 			      "GroupId": "sg-57c9812f",
 			      "GroupName": "webserver"
 			  }
 		      ],
 		      "SourceDestCheck": true,
 		      "SpotInstanceRequestId": null,
 		      "SriovNetSupport": null,
 		      "State": {
 			  "Code": 16,
 			  "Name": "running"
 		      },
 		      "StateReason": null,
 		      "StateTransitionReason": "",
 		      "SubnetId": "subnet-d76ae58f",
 		      "Tags": [
 			  {
 			      "Key": "Name",
 			      "Value": "webserver"
 			  }
 		      ],
 		      "VirtualizationType": "hvm",
 		      "VpcId": "vpc-b88883dc",
 		      "uname": "Linux ip-172-33-1-178 4.1.10-17.31.amzn1.x86_64 #1 SMP Sat Oct 24 01:31:37 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux\r\n"
 		  });

 ```

 ## AWS.ELBS.DEREGISTER
 <a name="deregister"></a>
 `aws.elbs.deregister(region, lbName, instance);`

 Deregister an EC2 instance with an ELB.

 Example:

 ```

 aws.elbs.deregister("us-east-1", "test-lb",
 		  {
 		      "AmiLaunchIndex": 0,
 		      "Architecture": "x86_64",
 		      "BlockDeviceMappings": [
 			  {
 			      "DeviceName": "/dev/xvda",
 			      "Ebs": {
 				  "AttachTime": "2016-03-31T19:17:44Z",
 				  "DeleteOnTermination": true,
 				  "Status": "attached",
 				  "VolumeId": "vol-3d1ab09f"
 			      }
 			  }
 		      ],
 		      "ClientToken": "",
 		      "EbsOptimized": false,
 		      "Hypervisor": "xen",
 		      "IamInstanceProfile": {
 			  "Arn": "arn:aws:iam::286536233385:instance-profile/cr-webserver",
 			  "Id": "AIPAIOOUMBIRCV2QCTIYA"
 		      },
 		      "ImageId": "ami-60b6c60a",
 		      "InstanceId": "i-824ad119",
 		      "InstanceLifecycle": null,
 		      "InstanceType": "t2.small",
 		      "KernelId": null,
 		      "KeyName": "cr",
 		      "LaunchTime": "2016-03-31T19:17:43Z",
 		      "Monitoring": {
 			  "State": "enabled"
 		      },
 		      "NetworkInterfaces": [
 			  {
 			      "Association": {
 				  "IpOwnerId": "amazon",
 				  "PublicDnsName": "ec2-52-90-244-101.compute-1.amazonaws.com",
 				  "PublicIp": "52.90.244.101"
 			      },
 			      "Attachment": {
 				  "AttachTime": "2016-03-31T19:17:43Z",
 				  "AttachmentId": "eni-attach-4d1c8ab1",
 				  "DeleteOnTermination": true,
 				  "DeviceIndex": 0,
 				  "Status": "attached"
 			      },
 			      "Description": "",
 			      "Groups": [
 				  {
 				      "GroupId": "sg-57c9812f",
 				      "GroupName": "webserver"
 				  }
 			      ],
 			      "MacAddress": "0e:6b:fe:c1:cb:45",
 			      "NetworkInterfaceId": "eni-a6b644f6",
 			      "OwnerId": "286536233385",
 			      "PrivateDnsName": "ip-172-33-1-178.ec2.internal",
 			      "PrivateIpAddress": "172.33.1.178",
 			      "PrivateIpAddresses": [
 				  {
 				      "Association": {
 					  "IpOwnerId": "amazon",
 					  "PublicDnsName": "ec2-52-90-244-101.compute-1.amazonaws.com",
 					  "PublicIp": "52.90.244.101"
 				      },
 				      "Primary": true,
 				      "PrivateDnsName": "ip-172-33-1-178.ec2.internal",
 				      "PrivateIpAddress": "172.33.1.178"
 				  }
 			      ],
 			      "SourceDestCheck": true,
 			      "Status": "in-use",
 			      "SubnetId": "subnet-d76ae58f",
 			      "VpcId": "vpc-b88883dc"
 			  }
 		      ],
 		      "Placement": {
 			  "Affinity": null,
 			  "AvailabilityZone": "us-east-1d",
 			  "GroupName": "",
 			  "HostId": null,
 			  "Tenancy": "default"
 		      },
 		      "Platform": null,
 		      "PrivateDnsName": "ip-172-33-1-178.ec2.internal",
 		      "PrivateIpAddress": "172.33.1.178",
 		      "ProductCodes": null,
 		      "PublicDnsName": "ec2-52-90-244-101.compute-1.amazonaws.com",
 		      "PublicIpAddress": "52.90.244.101",
 		      "RamdiskId": null,
 		      "RootDeviceName": "/dev/xvda",
 		      "RootDeviceType": "ebs",
 		      "SecurityGroups": [
 			  {
 			      "GroupId": "sg-57c9812f",
 			      "GroupName": "webserver"
 			  }
 		      ],
 		      "SourceDestCheck": true,
 		      "SpotInstanceRequestId": null,
 		      "SriovNetSupport": null,
 		      "State": {
 			  "Code": 16,
 			  "Name": "running"
 		      },
 		      "StateReason": null,
 		      "StateTransitionReason": "",
 		      "SubnetId": "subnet-d76ae58f",
 		      "Tags": [
 			  {
 			      "Key": "Name",
 			      "Value": "webserver"
 			  }
 		      ],
 		      "VirtualizationType": "hvm",
 		      "VpcId": "vpc-b88883dc",
 		      "uname": "Linux ip-172-33-1-178 4.1.10-17.31.amzn1.x86_64 #1 SMP Sat Oct 24 01:31:37 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux\r\n"
 		  });

 ```


