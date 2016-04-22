 


 # CORE FUNCTIONS: SNS


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.sns.scan](#vscan)
 > * [aws.sns.create](#vcreate)
 > * [aws.sns.delete](#vdelete)
 > * [aws.sns.describe](#vdescribe)

 This API allows resource handlers to manage SNS.

 ## AWS.SNS.SCAN
 <a name="vscan"></a>
 `aws.sns.scan(region);`

 Returns a list of sns.

 Example:

 ```

  var sns =  aws.sns.scan("us-east-1");

 ```

 ## AWS.SNS.CREATE
 <a name="vcreate"></a>
 `aws.sns.create(region, config, gateway);`

 Create a SNS.

 Example:

 ```

  var sns =  aws.sns.create(
    "us-east-1",
    {
		  CidrBlock:       "172.33.0.0/16"
    },
    true);

 ```

 ## AWS.SNS.DELETE
 <a name="vdelete"></a>
 `aws.sns.delete(region, sns_id);`

 Delete a SNS.

 Example:

 ```

  aws.sns.delete("us-east-1", "sns-abcd");

 ```

 ## AWS.SNS.DESCRIBE
 <a name="vdescribe"></a>
 `aws.sns.describe(region, sns_id);`

 Get info from AWS about a SNS.

 Example:

 ```

  var sns = aws.sns.describe("us-east-1", "sns-abcd");

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


