 


 # CORE FUNCTIONS: VPC


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.vpcs.scan](#vscan)
 > * [aws.vpcs.create](#vcreate)
 > * [aws.vpcs.delete](#vdelete)
 > * [aws.vpcs.describe](#vdescribe)
 > * [aws.gateways.scan](#gscan)
 > * [aws.gateways.create](#gcreate)
 > * [aws.gateways.delete](#gdelete)
 > * [aws.gateways.describe](#gdescribe)

 This API allows resource handlers to manage VPCS.

 ## AWS.VPCS.SCAN
 <a name="vscan"></a>
 `aws.vpcs.scan(region);`

 Returns a list of vpcs.

 Example:

 ```

  var vpcs =  aws.vpcs.scan("us-east-1");

 ```

 ## AWS.VPCS.CREATE
 <a name="vcreate"></a>
 `aws.vpcs.create(region, config, gateway);`

 Create a VPC.

 Example:

 ```

  var vpc =  aws.vpcs.create(
    "us-east-1",
    {
		  CidrBlock:       "172.33.0.0/16"
    },
    true);

 ```

 ## AWS.VPCS.DELETE
 <a name="vdelete"></a>
 `aws.vpcs.delete(region, vpc_id);`

 Delete a VPC.

 Example:

 ```

  aws.vpcs.delete("us-east-1", "vpc-abcd");

 ```

 ## AWS.VPCS.DESCRIBE
 <a name="vdescribe"></a>
 `aws.vpcs.describe(region, vpc_id);`

 Get info from AWS about a VPC.

 Example:

 ```

  var vpc = aws.vpcs.describe("us-east-1", "vpc-abcd");

 ```

 ## AWS.GATEWAYS.SCAN
 <a name="gscan"></a>
 `aws.gateways.scan(region);`

 Returns a list of gateways.

 Example:

 ```

  var gateways =  mithras.gateways.scan("us-east-1");

 ```

 ## AWS.GATEWAYS.CREATE
 <a name="gcreate"></a>
 `aws.gateways.create(region);`

 Create a gateway.

 Example:

 ```

  var gateway =  mithras.gateways.create("us-east-1");

 ```

 ## AWS.GATEWAYS.DELETE
 <a name="gdelete"></a>
 `aws.gateways.delete(region, gateway_id);`

 Delete a gateway.

 Example:

 ```

  mithras.gateways.delete("us-east-1", "gw-abcd");

 ```

 ## AWS.GATEWAYS.DESCRIBE
 <a name="gdescribe"></a>
 `aws.gateways.describe(region, gateway_id);`

 Get info from AWS about a gateway.

 Example:

 ```

  var gateway = mithras.gateways.describe("us-east-1", "gw-abcd");

 ```


