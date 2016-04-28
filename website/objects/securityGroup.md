
## Object Format: SECURITYGROUP


    {
      "Description": "default group",
      "GroupId": "sg-c823c2a1",
      "GroupName": "default",
      "IpPermissions": [
        {
          "FromPort": 0,
          "IpProtocol": "tcp",
          "IpRanges": null,
          "PrefixListIds": null,
          "ToPort": 65535,
          "UserIdGroupPairs": [
            {
              "GroupId": "sg-c823c2a1",
              "GroupName": "default",
              "UserId": "286536233385"
            }
          ]
        },
        {
          "FromPort": 0,
          "IpProtocol": "udp",
          "IpRanges": null,
          "PrefixListIds": null,
          "ToPort": 65535,
          "UserIdGroupPairs": [
            {
              "GroupId": "sg-c823c2a1",
              "GroupName": "default",
              "UserId": "286536233385"
            }
          ]
        },
        {
          "FromPort": -1,
          "IpProtocol": "icmp",
          "IpRanges": null,
          "PrefixListIds": null,
          "ToPort": -1,
          "UserIdGroupPairs": [
            {
              "GroupId": "sg-c823c2a1",
              "GroupName": "default",
              "UserId": "286536233385"
            }
          ]
        },
        {
          "FromPort": 22,
          "IpProtocol": "tcp",
          "IpRanges": [
            {
              "CidrIp": "0.0.0.0/0"
            }
          ],
          "PrefixListIds": null,
          "ToPort": 22,
          "UserIdGroupPairs": null
        },
        {
          "FromPort": 80,
          "IpProtocol": "tcp",
          "IpRanges": [
            {
              "CidrIp": "0.0.0.0/0"
            }
          ],
          "PrefixListIds": null,
          "ToPort": 80,
          "UserIdGroupPairs": null
        }
      ],
      "IpPermissionsEgress": null,
      "OwnerId": "286536233385",
      "Tags": null,
      "VpcId": null
    }
