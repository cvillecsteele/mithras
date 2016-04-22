 


 # CORE FUNCTIONS: SNS


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.sns.topics.scan](#tscan)
 > * [aws.sns.topics.create](#tcreate)
 > * [aws.sns.topics.delete](#tdelete)
 > * [aws.sns.topics.describe](#tdescribe)

 This API allows resource handlers to manage SNS.

 ## AWS.SNS.TOPICS.SCAN
 <a name="tscan"></a>
 `aws.sns.topics.scan(region);`

 Returns a list of sns topics.

 Example:

 ```

  var topics = aws.sns.topics.scan("us-east-1");

 ```

 ## AWS.SNS.TOPICS.CREATE
 <a name="tcreate"></a>
 `aws.sns.create(region, config);`

 Create a SNS topic.

 Example:

 ```

  var sns =  aws.sns.topics.create(
    "us-east-1",
    {
		  Name: "mytopic"
    });

 ```

 ## AWS.SNS.TOPICS.DELETE
 <a name="tdelete"></a>
 `aws.sns.topics.delete(region, name);`

 Delete an SNS topic.

 Example:

 ```

  aws.sns.topics.delete("us-east-1", "arn:aws:sns:us-east-1:286536233385:Test");

 ```

 ## AWS.SNS.TOPICS.DESCRIBE
 <a name="tdescribe"></a>
 `aws.sns.topics.describe(region, sns_id);`

 Get info from AWS about an SNS topic.

 Example:

 ```

  var topic = aws.sns.topics.describe("us-east-1", "arn:aws:sns:us-east-1:286536233385:Test");

 ```

 ## AWS.SNS.SUBS.SCAN
 <a name="sscan"></a>
 `aws.sns.subs.scan(region);`

 Returns a list of sns subs.

 Example:

 ```

  var subs = aws.sns.subs.scan("us-east-1");

 ```

 ## AWS.SNS.SUBS.CREATE
 <a name="screate"></a>
 `aws.sns.create(region, config);`

 Create a SNS sub.

 Example:

 ```

  var sns =  aws.sns.subs.create(
    "us-east-1",
    {
		  Protocol: "email"
		  TopicArn: mithras.watch(rTopic.name+"._target.topic")
		  Endpoint: "colin@mithras.io"
    });

 ```

 ## AWS.SNS.SUBS.DELETE
 <a name="sdelete"></a>
 `aws.sns.subs.delete(region, sub_id);`

 Delete an SNS sub.

 Example:

 ```

  aws.sns.subs.delete("us-east-1", "arn:aws:sns:us-east-1:286536233385:Test:317abc61-7374-4d94-9947-19b1c26e119d");

 ```

 ## AWS.SNS.SUBS.DESCRIBE
 <a name="sdescribe"></a>
 `aws.sns.subs.describe(region, sub_id);`

 Get info from AWS about an SNS sub.

 Example:

 ```

  aws.sns.subs.describe("us-east-1", "arn:aws:sns:us-east-1:286536233385:Test:317abc61-7374-4d94-9947-19b1c26e119d");

 ```


