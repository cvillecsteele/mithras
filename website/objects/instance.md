
## Object Format: INSTANCE


    {
      "AmiLaunchIndex": 0,
      "Architecture": "x86_64",
      "BlockDeviceMappings": [
        {
          "DeviceName": "/dev/xvda",
          "Ebs": {
            "AttachTime": "2016-02-29T18:24:32Z",
            "DeleteOnTermination": true,
            "Status": "attached",
            "VolumeId": "vol-57df9cf4"
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
      "InstanceId": "i-6c7f79f4",
      "InstanceLifecycle": null,
      "InstanceType": "t2.small",
      "KernelId": null,
      "KeyName": "cr",
      "LaunchTime": "2016-02-29T18:24:30Z",
      "Monitoring": {
        "State": "disabled"
      },
      "NetworkInterfaces": [
        {
          "Association": {
            "IpOwnerId": "amazon",
            "PublicDnsName": "ec2-54-210-146-169.compute-1.amazonaws.com",
            "PublicIp": "54.210.146.169"
          },
          "Attachment": {
            "AttachTime": "2016-02-29T18:24:30Z",
            "AttachmentId": "eni-attach-b3c5455c",
            "DeleteOnTermination": true,
            "DeviceIndex": 0,
            "Status": "attached"
          },
          "Description": "",
          "Groups": [
            {
              "GroupId": "sg-15b1726d",
              "GroupName": "cr-webserver-sg-staging"
            }
          ],
          "MacAddress": "0e:e3:1d:63:2f:29",
          "NetworkInterfaceId": "eni-040cc555",
          "OwnerId": "286536233385",
          "PrivateDnsName": "ip-172-22-1-116.ec2.internal",
          "PrivateIpAddress": "172.22.1.116",
          "PrivateIpAddresses": [
            {
              "Association": {
                "IpOwnerId": "amazon",
                "PublicDnsName": "ec2-54-210-146-169.compute-1.amazonaws.com",
                "PublicIp": "54.210.146.169"
              },
              "Primary": true,
              "PrivateDnsName": "ip-172-22-1-116.ec2.internal",
              "PrivateIpAddress": "172.22.1.116"
            }
          ],
          "SourceDestCheck": true,
          "Status": "in-use",
          "SubnetId": "subnet-e88d56b0",
          "VpcId": "vpc-f6732c92"
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
      "PrivateDnsName": "ip-172-22-1-116.ec2.internal",
      "PrivateIpAddress": "172.22.1.116",
      "ProductCodes": null,
      "PublicDnsName": "ec2-54-210-146-169.compute-1.amazonaws.com",
      "PublicIpAddress": "54.210.146.169",
      "RamdiskId": null,
      "RootDeviceName": "/dev/xvda",
      "RootDeviceType": "ebs",
      "SecurityGroups": [
        {
          "GroupId": "sg-15b1726d",
          "GroupName": "cr-webserver-sg-staging"
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
      "SubnetId": "subnet-e88d56b0",
      "Tags": [
        {
          "Key": "Env",
          "Value": "staging"
        },
        {
          "Key": "API-SHA",
          "Value": "0e2cd4f"
        },
        {
          "Key": "Name",
          "Value": "cr-webserver-staging"
        },
        {
          "Key": "CRUI-SHA",
          "Value": "be0408a"
        }
      ],
      "VirtualizationType": "hvm",
      "VpcId": "vpc-f6732c92"
    }
