 


 # CORE FUNCTIONS: SQS


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.sqs.scan](#scan)
 > * [aws.sqs.create](#create)
 > * [aws.sqs.delete](#delete)
 > * [aws.sqs.describe](#describe)
 > * [aws.sqs.attributes](#attributes)
 > * [aws.sqs.setAttributes](#setAttributes)
 > * [aws.sqs.messages.send](#msend)
 > * [aws.sqs.messages.receive](#mreceive)
 > * [aws.sqs.messages.delete](#mdelete)

 This API allows resource handlers to manage SQS.

 ## AWS.SQS.SCAN
 <a name="scan"></a>
 `aws.sqs.scan(region);`

 Returns a list of sqs queues.

 Example:

 ```

  var queues = aws.sqs.scan("us-east-1");

 ```

 ## AWS.SQS.CREATE
 <a name="create"></a>
 `aws.sqs.create(region, config);`

 Create a SQS queue.

 Example:

 ```

  var sqs =  aws.sqs.create(
    "us-east-1",
    {
      QueueName: "myqueue"
      Attributes: [
        "Key": "value"
      ]
    });

 ```

 ## AWS.SQS.DELETE
 <a name="delete"></a>
 `aws.sqs.delete(region, url);`

 Delete an SQS queue.

 Example:

 ```

  aws.sqs.delete("us-east-1", "queueUrl");

 ```

 ## AWS.SQS.MESSAGES.SEND
 <a name="msend"></a>
 `aws.sqs.messages.send(region, input);`

 Publish to an SQS queue.

 Example:

 ```

  aws.sqs.messages.send("us-east-1",
                  {
                    MessageBody:  "body"
                    QueueUrl:     "url"
                    DelaySeconds: 1
                    MessageAttributes: {
                      "Key": {
                        DataType: "type"
                        BinaryListValues: [
                          "PAYLOAD"
                        ]
                        BinaryValue: "PAYLOAD"
                        StringListValues: [
                          "String"
                        ]
                        StringValue: "String"
                      }
                    }
                   });

 ```

 ## AWS.SQS.DESCRIBE
 <a name="describe"></a>
 `aws.sqs.describe(region, sqs_id);`

 Get info from AWS about an SQS queue.

 Example:

 ```

  var queue = aws.sqs.describe("us-east-1", "queueName");

 ```

 ## AWS.SQS.ATTRIBUTES
 <a name="attributes"></a>
 `aws.sqs.attributes(region, sqsUrl);`

 Get info from AWS about an SQS queue.

 Example:

 ```

  var queue = aws.sqs.attributes("us-east-1", "queueUrl");

 ```

 ## AWS.SQS.SETATTRIBUTES
 <a name="setAttributes"></a>
 `aws.sqs.setAttributes(region, attrs);`

 Put queue attributes.

 Example:

 ```

 aws.sqs.setAttributes(region, {
     Attributes: {
       "Policy": "..."
     }
     QueueUrl: "..."
 });

 ```

 ## AWS.SQS.MESSAGES.RECEIVE
 <a name="mreceive"></a>
 `aws.sqs.messages.receive(region, input);`

 Get a message from a queue

 Example:

 ```
  var message =
  aws.sqs.messages.receive("us-east-1",
                  {
                    QueueUrl: "url"
                    AttributeNames: [
                      "QueueAttributeName"
                    ]
                    MaxNumberOfMessages: 1
                    MessageAttributeNames: [
                      "MessageAttributeName"
                    ]
                    VisibilityTimeout: 1
                    WaitTimeSeconds:   1
                  });

 ```

 ## AWS.SQS.MESSAGES.DELETE
 <a name="mdelete"></a>
 `aws.sqs.messages.delete(region, input);`

 Get a message from a queue

 Example:

 ```
  aws.sqs.messages.delete("us-east-1",
  {
    QueueUrl:      "queueUrl"
    ReceiptHandle: "123456"
  });

 ```


