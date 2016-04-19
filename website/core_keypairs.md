 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.keypairs.scan](#vscan)
 > * [aws.keypairs.create](#vcreate)
 > * [aws.keypairs.delete](#vdelete)
 > * [aws.keypairs.describe](#vdescribe)

 This API allows resource handlers to manage AWS SSH keypairs.

 ## AWS.KEYPAIRS.SCAN
 <a name="vscan"></a>
 `aws.keypairs.scan(region);`

 Returns a list of keypairs.

 Example:

 ```

  var keypairs =  aws.keypairs.scan("us-east-1");

 ```

 ## AWS.KEYPAIRS.CREATE
 <a name="vcreate"></a>
 `aws.keypairs.create(region, config);`

 Create a keypair.

 Example:

 ```

  var keypair =  aws.keypairs.create(
    "us-east-1",
    {
		  KeyName: "my-keypair"
    });

 ```

 ## AWS.KEYPAIRS.DELETE
 <a name="vdelete"></a>
 `aws.keypairs.delete(region, keypair_id);`

 Delete a keypair.

 Example:

 ```

  aws.keypairs.delete("us-east-1", "my-keypair");

 ```

 ## AWS.KEYPAIRS.DESCRIBE
 <a name="vdescribe"></a>
 `aws.keypairs.describe(region, keypair_id);`

 Get info from AWS about a KEYPAIR.

 Example:

 ```

  var keypair = aws.keypairs.describe("us-east-1", "my-keypair");

 ```


