 
 
 # SQS
 
 SQS is a resource handler for dealing with AWS SQS resources.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"sqs"`
 
 Usage:
 
 `var sqs = require("sqs").init();`
 
  ## Example Resource
 
 ```javascript
 var rQueue = {
     name: "sqsQueue"
     module: "sqs"
     params: {
        region: defaultRegion
        ensure: ensure
        queue: {
          QueueName: "myqueue"
          Attributes: [
            "Key": "value"
          ]
        }
     }
 };
 var rSub = {
     name: "sqsSub"
     module: "sqs"
     dependsOn: [rTopic.name]
     params: {
         region: defaultRegion
         ensure: ensure
         sub: {
           Protocol: "..."
           TopicArn: "..."
           Endpoint: "..."
         }
     }
 };
 var rPub = {
     name: "sqsPub"
     module: "sqs"
     dependsOn: [rTopic.name]
     params: {
         region: defaultRegion
         ensure: ensure
         sub: {...}
     }
 };
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"` and the sqs queue `params.queue.QueueName` does not
 exist, it is created.  If `"absent"`, and it exists, it is removed.
 
 If `"present"` and the the `params.message` property is set, a message
 is published to the queue.  This is NOT an idempotent operation.
 
 ### `queue`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/sqs.html#type-CreateQueueInput)

 Parameters for queue creation.

 ### `message`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/sqs.html#type-SendMessageInput)

 Parameters for publishing a message to a queue.

 ### `on_find`

 * Required: true
 * Allowed Values: A function taking two parameters: `catalog` and `resource`

 If defined in the resource's `params` object, the `on_find`
 function provides a way for a matching resource to be identified
 using a user-defined way.  The function is called with the current
 `catalog`, as well as the `resource` object itself.  The function
 can look through the catalog, find a matching object using whatever
 logic you want, and return it.  If the function returns `undefined`
 or a n empty Javascript array, (`[]`), the function is indicating
 that no matching resource was found in the `catalog`.
 

