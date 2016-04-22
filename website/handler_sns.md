 
 
 # SNS
 
 SNS is a resource handler for dealing with AWS SNS resources.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"sns"`
 
 Usage:
 
 `var sns = require("sns").init();`
 
  ## Example Resource
 
 ```javascript
 var rTopic = {
     name: "snsTopic"
     module: "sns"
     dependsOn: [rVpc.name]
     params: {
         region: defaultRegion
         ensure: ensure
         topic: {
             Name:  "my-topic"
         }
     }
 };
 var rSub = {
     name: "snsSub"
     module: "sns"
     dependsOn: [rVpc.name]
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
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"` and the sns topic `params.topic.Name` does not
 exist, it is created.  If `"absent"`, and it exists, it is removed.
 
 If `"present"` and the sns subscription referencing
 `params.topic.Name` does not exist, it is created.  If `"absent"`,
 and it exists, it is removed.
 
 ### `topic`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/sns.html#type-CreateTopicInput)

 Parameters for topic creation.

 ### `sub`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/sns.html#type-SubscribeInput)

 Parameters for subscription creation.

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
 

