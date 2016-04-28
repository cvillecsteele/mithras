
## Object Format: ELB


    {
      "AvailabilityZones": [
        "us-east-1b",
        "us-east-1d"
      ],
      "BackendServerDescriptions": null,
      "CanonicalHostedZoneName": "cr-ws-elb-staging-1835873414.us-east-1.elb.amazonaws.com",
      "CanonicalHostedZoneNameID": "Z3DZXE0Q79N41H",
      "CreatedTime": "2016-02-18T16:40:01.9Z",
      "DNSName": "cr-ws-elb-staging-1835873414.us-east-1.elb.amazonaws.com",
      "HealthCheck": {
        "HealthyThreshold": 10,
        "Interval": 30,
        "Target": "HTTP:80/hc",
        "Timeout": 5,
        "UnhealthyThreshold": 2
      },
      "Instances": [
        {
          "InstanceId": "i-6c7f79f4"
        }
      ],
      "ListenerDescriptions": [
        {
          "Listener": {
            "InstancePort": 80,
            "InstanceProtocol": "HTTP",
            "LoadBalancerPort": 443,
            "Protocol": "HTTPS",
            "SSLCertificateId": "arn:aws:acm:us-east-1:286536233385:certificate/7ec23cde-ea82-4740-b44f-91302b48f9f9"
          },
          "PolicyNames": [
            "ELBSecurityPolicy-2015-05"
          ]
        },
        {
          "Listener": {
            "InstancePort": 80,
            "InstanceProtocol": "HTTP",
            "LoadBalancerPort": 80,
            "Protocol": "HTTP",
            "SSLCertificateId": null
          },
          "PolicyNames": null
        }
      ],
      "LoadBalancerName": "cr-ws-elb-staging",
      "Policies": {
        "AppCookieStickinessPolicies": null,
        "LBCookieStickinessPolicies": null,
        "OtherPolicies": [
          "ELBSecurityPolicy-2015-05"
        ]
      },
      "Scheme": "internet-facing",
      "SecurityGroups": [
        "sg-f9b17281"
      ],
      "SourceSecurityGroup": {
        "GroupName": "cr-elb-sg-staging",
        "OwnerAlias": "286536233385"
      },
      "Subnets": [
        "subnet-bf6cedc9",
        "subnet-e88d56b0"
      ],
      "VPCId": "vpc-f6732c92"
    }
