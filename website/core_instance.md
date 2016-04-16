 


 # CORE FUNCTIONS: INSTANCE


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.instances.scan](#scan)
 > * [aws.instances.create](#create)
 > * [aws.instances.delete](#delete)
 > * [aws.instances.describe](#describe)

 This API allows resource handlers to manage EC2 instances.

 ## AWS.INSTANCES.SCAN
 <a name="scan"></a>
 `aws.instances.scan(region);`

 Returns a list of instances.

 Example:

 ```

  var instances =  aws.instances.scan("us-east-1");

 ```

 ## AWS.INSTANCES.CREATE
 <a name="create"></a>
 `aws.instances.create(region, config);`

 Create one or more instances.

 Example:

 ```

  var inst =  aws.instances.create(
    "us-east-1",
    {
      ImageId:        "ami-60b6c60a"
      MaxCount:       1
      MinCount:       1
      DisableApiTermination: false
      EbsOptimized:          false
      IamInstanceProfile: {
        Name: iamProfileName
      }
      InstanceInitiatedShutdownBehavior: "terminate"
      InstanceType:                      "t2.small"
      KeyName:                           "my-key"
      Monitoring: {
        Enabled: true
      }
      NetworkInterfaces: [
        {
          AssociatePublicIpAddress: true
          DeleteOnTermination:      true
          DeviceIndex:              0
          Groups:                  [ "sg-1234" ]
          SubnetId:                "subnet-abcd"
        }
      ]
    });

 ```

 ## AWS.INSTANCES.DELETE
 <a name="delete"></a>
 `aws.instances.delete(region, instance_id);`

 Delete an instance

 Example:

 ```

  aws.instances.delete("us-east-1", "i-abcd");

 ```

 ## AWS.INSTANCES.DESCRIBE
 <a name="describe"></a>
 `aws.instances.describe(region, inst_id);`

 Get info from AWS about an instance.

 Example:

 ```

  var i = aws.instances.describe("us-east-1", "i-abcd");

 ```


