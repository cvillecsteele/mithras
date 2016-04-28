
## Object Format: DB


    {
      "AllocatedStorage": 10,
      "AutoMinorVersionUpgrade": true,
      "AvailabilityZone": "us-east-1d",
      "BackupRetentionPeriod": 1,
      "CACertificateIdentifier": "rds-ca-2015",
      "CharacterSetName": null,
      "CopyTagsToSnapshot": false,
      "DBClusterIdentifier": null,
      "DBInstanceClass": "db.m1.small",
      "DBInstanceIdentifier": "cr-staging",
      "DBInstanceStatus": "available",
      "DBName": "cr",
      "DBParameterGroups": [
        {
          "DBParameterGroupName": "default.postgres9.4",
          "ParameterApplyStatus": "in-sync"
        }
      ],
      "DBSecurityGroups": null,
      "DBSubnetGroup": {
        "DBSubnetGroupDescription": "postgres cache subnet group",
        "DBSubnetGroupName": "cr-postgres-subnet-group-staging",
        "SubnetGroupStatus": "Complete",
        "Subnets": [
          {
            "SubnetAvailabilityZone": {
              "Name": "us-east-1b"
            },
            "SubnetIdentifier": "subnet-bf6cedc9",
            "SubnetStatus": "Active"
          },
          {
            "SubnetAvailabilityZone": {
              "Name": "us-east-1d"
            },
            "SubnetIdentifier": "subnet-e88d56b0",
            "SubnetStatus": "Active"
          }
        ],
        "VpcId": "vpc-f6732c92"
      },
      "DbInstancePort": 0,
      "DbiResourceId": "db-D22VWSG5EQGTKCLC74VVHMZ75Y",
      "Endpoint": {
        "Address": "cr-staging.cpxl92pgez22.us-east-1.rds.amazonaws.com",
        "HostedZoneId": null,
        "Port": 5432
      },
      "Engine": "postgres",
      "EngineVersion": "9.4.5",
      "EnhancedMonitoringResourceArn": null,
      "InstanceCreateTime": "2016-02-18T16:33:41.001Z",
      "Iops": null,
      "KmsKeyId": null,
      "LatestRestorableTime": "2016-04-28T19:02:20Z",
      "LicenseModel": "postgresql-license",
      "MasterUsername": "cr_user",
      "MonitoringInterval": 0,
      "MonitoringRoleArn": null,
      "MultiAZ": false,
      "OptionGroupMemberships": [
        {
          "OptionGroupName": "default:postgres-9-4",
          "Status": "in-sync"
        }
      ],
      "PendingModifiedValues": {
        "AllocatedStorage": null,
        "BackupRetentionPeriod": null,
        "CACertificateIdentifier": null,
        "DBInstanceClass": null,
        "DBInstanceIdentifier": null,
        "EngineVersion": null,
        "Iops": null,
        "MasterUserPassword": null,
        "MultiAZ": null,
        "Port": null,
        "StorageType": null
      },
      "PreferredBackupWindow": "06:44-07:14",
      "PreferredMaintenanceWindow": "thu:03:41-thu:04:11",
      "PubliclyAccessible": false,
      "ReadReplicaDBInstanceIdentifiers": null,
      "ReadReplicaSourceDBInstanceIdentifier": null,
      "SecondaryAvailabilityZone": null,
      "StatusInfos": null,
      "StorageEncrypted": false,
      "StorageType": "standard",
      "TdeCredentialArn": null,
      "VpcSecurityGroups": [
        {
          "Status": "active",
          "VpcSecurityGroupId": "sg-02b1727a"
        }
      ]
    }
