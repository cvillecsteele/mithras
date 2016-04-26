 


 # CORE FUNCTIONS: SCALING


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.scaling.scan](#scan)
 > * [aws.scaling.create](#create)
 > * [aws.scaling.delete](#delete)
 > * [aws.scaling.describe](#describe)
 > * [aws.scaling.messages.send](#msend)
 > * [aws.scaling.messages.receive](#mreceive)
 > * [aws.scaling.messages.delete](#mdelete)

 This API allows resource handlers to manage Autoscaling groups.

 ## AWS.SCALING.SCAN
 <a name="scan"></a>
 `aws.scaling.scan(region);`

 Returns a list of scaling queues.

 Example:

 ```

  var queues = aws.scaling.scan("us-east-1");

 ```

 ## AWS.SCALING.CREATE
 <a name="create"></a>
 `aws.scaling.create(region, config);`

 Create a SCALING queue.

 Example:

 ```

  var scaling =  aws.scaling.create(
    "us-east-1",
    {
      QueueName: "myqueue"
      Attributes: [
        "Key": "value"
      ]
    });

 ```

 ## AWS.SCALING.DELETE
 <a name="delete"></a>
 `aws.scaling.delete(region, url);`

 Delete an SCALING queue.

 Example:

 ```

  aws.scaling.delete("us-east-1", "queueUrl");

 ```

 ## AWS.SCALING.MESSAGES.SEND
 <a name="msend"></a>
 `aws.scaling.messages.send(region, input);`

 Publish to an SCALING queue.

 Example:

 ```

  aws.scaling.messages.send("us-east-1",
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

 ## AWS.SCALING.DESCRIBE
 <a name="describe"></a>
 `aws.scaling.describe(region, scaling_id);`

 Get info from AWS about an SCALING queue.

 Example:

 ```

  var queue = aws.scaling.describe("us-east-1", "queueName");

 ```

 ## AWS.SCALING.MESSAGES.RECEIVE
 <a name="mreceive"></a>
 `aws.scaling.messages.receive(region, input);`

 Get a message from a queue

 Example:

 ```
  var message =
  aws.scaling.messages.receive("us-east-1",
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

 ## AWS.SCALING.MESSAGES.DELETE
 <a name="mdelete"></a>
 `aws.scaling.messages.delete(region, input);`

 Get a message from a queue

 Example:

 ```
  aws.scaling.messages.delete("us-east-1",
  {
    QueueUrl:      "queueUrl"
    ReceiptHandle: "123456"
  });

 ```


