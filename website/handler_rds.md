 
 
 # Elasticache
 
 Elasticache is resource handler for managing AWS caches.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"elasticache"`
 
 Usage:
 
 `var elasticache = require("elasticache").init();`
 
  ## Example Resource
 
 ```javascript
 var rds = {
   name: "rdsA"
   module: "rds"
   dependsOn: [otherResource.name]
   params: {
     ensure: ensure
     region: defaultRegion
     wait: true
     subnetGroup: {
         DBSubnetGroupDescription: "test subnet group"
         DBSubnetGroupName: "test-subnet-group"
         SubnetIds: [
              "subnet-123",
              "subnet-456"
         ]
         Tags: [
             {
                 Key:   "Foo"
                 Value: "Bar"
             }
         ]
     }
     db: {
         DBInstanceClass:         "db.m1.small"
         DBInstanceIdentifier:    "test-rds"
         Engine:                  "mysql"
         AllocatedStorage:        10
         AutoMinorVersionUpgrade: true
         AvailabilityZone:        defaultZone
         MasterUserPassword:      "test123456789"
         MasterUsername:          "test"
         DBSubnetGroupName:       "test-subnet-group"
         DBName:                  "test"
         PubliclyAccessible:      false
         Tags: [
             {
                 Key:   "foo"
                 Value: "bar"
             },
         ]
     }
     delete: {
         DBInstanceIdentifier:      "db-abcd"
         FinalDBSnapshotIdentifier: "byebye" + Date.now()
         SkipFinalSnapshot:         true
     }
   }
 };
 
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"`, the db specified by `db` will be created, and
 if `"absent"`, it will be removed using the `delete` property.
 
 ### `region`

 * Required: true
 * Allowed Values: string, any valid AWS region; eg "us-east-1"

 The region for calls to the AWS API.
 
 ### `wait`

 * Required: false
 * Allowed Values: true or false

 If `true`, delay execution until the db has been created in AWS.
 
 ### `subnetGroup`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/rds.html#type-CreateDbSubnetGroupInput)

 If set, a subnet group will be created for your db.
 
 ### `db`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/rds.html#type-CreateDBClusterInput)

 Parameters for resource creation.
 
 ### `delete`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/rds.html#type-DeleteDBClusterInput)

 Parameters for deletion.
 

