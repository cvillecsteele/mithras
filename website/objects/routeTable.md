
## Object Format: ROUTETABLE


    {
      "Associations": [
        {
          "Main": false,
          "RouteTableAssociationId": "rtbassoc-dda3d6ba",
          "RouteTableId": "rtb-d7ca2eb0",
          "SubnetId": "subnet-bf6cedc9"
        }
      ],
      "PropagatingVgws": null,
      "RouteTableId": "rtb-d7ca2eb0",
      "Routes": [
        {
          "DestinationCidrBlock": "172.22.0.0/16",
          "DestinationPrefixListId": null,
          "GatewayId": "local",
          "InstanceId": null,
          "InstanceOwnerId": null,
          "NatGatewayId": null,
          "NetworkInterfaceId": null,
          "Origin": "CreateRouteTable",
          "State": "active",
          "VpcPeeringConnectionId": null
        },
        {
          "DestinationCidrBlock": "0.0.0.0/0",
          "DestinationPrefixListId": null,
          "GatewayId": "igw-463be322",
          "InstanceId": null,
          "InstanceOwnerId": null,
          "NatGatewayId": null,
          "NetworkInterfaceId": null,
          "Origin": "CreateRoute",
          "State": "active",
          "VpcPeeringConnectionId": null
        }
      ],
      "Tags": null,
      "VpcId": "vpc-f6732c92"
    }
