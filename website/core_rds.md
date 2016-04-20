 


 # CORE FUNCTIONS: RDS


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.rds.scan](#scan)
 > * [aws.rds.create](#create)
 > * [aws.rds.delete](#delete)
 > * [aws.rds.describe](#describe)

 > * [aws.rds.subnetGroups.create](#gcreate)
 > * [aws.rds.subnetGroups.delete](#gdelete)
 > * [aws.rds.subnetGroups.describe](#gdescribe)

 This API allows resource handlers to manage RDS.

 ## AWS.RDS.SCAN
 <a name="scan"></a>
 `aws.rds.scan(region);`

 Returns a list of RDS clusters.

 Example:

 ```

  var dbs = aws.rds.scan("us-east-1");

 ```

 ## AWS.RDS.CREATE
 <a name="create"></a>
 `aws.rds.create(region, config, wait);`

 Create an RDS cluster.

 Example:

 ```

  var db = aws.rds.create("us-east-1",
   {
      DBInstanceClass:         "db.m1.small"
      DBInstanceIdentifier:    "test-rds"
      Engine:                  "mysql"
      AllocatedStorage:        10
      AutoMinorVersionUpgrade: true
      AvailabilityZone:        "us-east-1"
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
   },
   true);

 ```

 ## AWS.RDS.DELETE
 <a name="delete"></a>
 `aws.rds.delete(region, config);`

 Delete an RDS cluster.

 Example:

 ```

  var db = aws.rds.delete("us-east-1",
     {
		   DBInstanceIdentifier:      "db-xyz",
		   FinalDBSnapshotIdentifier: "byebye" + Date.now()
		   SkipFinalSnapshot:         true
     });

 ```

 ## AWS.RDS.DESCRIBE
 <a name="describe"></a>
 `aws.rds.describe(region, id);`

 Get info about an RDS cluster.

 Example:

 ```

  var db = aws.rds.describe("us-east-1", "db-xyz");

 ```

 ## AWS.RDS.SUBNETGROUPS.DESCRIBE
 <a name="gdescribe"></a>
 `aws.rds.subnetGroups.describe(region, id);`

 Get info about an RDS subnet group.

 Example:

 ```

  var group = aws.rds.subnetGroups.describe("us-east-1", "sg-xyz");

 ```

 ## AWS.RDS.SUBNETGROUPS.CREATE
 <a name="gcreate"></a>
 `aws.rds.subnetGroups.create(region, config);`

 Create an RDS subnet group.

 Example:

 ```

  var group = aws.rds.subnetGroups.create("us-east-1",
 {
 		DBSubnetGroupDescription: "test subnet group"
 		DBSubnetGroupName: "test-subnet-group"
 		SubnetIds: [
       "subnet-1"
       "subnet-2"
 		]
 		Tags: [
 		    {
 			     Key:   "Foo"
 			     Value: "Bar"
 		    }
 		]
 });

 ```

 ## AWS.RDS.SUBNETGROUPS.DELETE
 <a name="gdelete"></a>
 `aws.rds.subnetGroups.delete(region, id);`

 Delete an RDS subnet group.

 Example:

 ```

  var group = aws.rds.subnetGroups.delete("us-east-1", "sg-xyz");

 ```


