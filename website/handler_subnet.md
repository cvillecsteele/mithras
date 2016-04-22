 
 
 # subnet
 
 Subnet is a resource handler for dealing with AWS security groups.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"subnet"`
 
 Usage:
 
 `var subnet = require("subnet").init();`
 
  ## Example Resource
 
 ```javascript
 var rSubnetA = {
     name: "subnetA"
     module: "subnet"
     dependsOn: [rVpc.name]
     params: {
         region: defaultRegion
         ensure: ensure
         subnet: {
             CidrBlock:        "172.33.1.0/24"
             VpcId:            mithras.watch("VPC._target.VpcId")
             AvailabilityZone: defaultZone
         }
         tags: {
             Name: "primary-subnet"
         }
         routes: [
             {
                 DestinationCidrBlock: "0.0.0.0/0"
                 GatewayId:            mithras.watch("VPC._target.VpcId", mithras.findGWByVpcId)
             }
         ]
     }
 };
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"` and the subnet named
 `params.subnet.GroupName` does not exist, it is created.  If
 `"absent"`, and the subnet exists, it is removed.
 
 ### `subnet`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-CreateSubnetInput)

 This file is copied to remote hosts.

 ### `tags`

 * Required: false
 * Allowed Values: a map of tags to be set on a created subnet

 For tagging.

 ### `routes`

 * Required: false
 * Allowed Values: an array of JSON objects corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-CreateRouteInput)

 An list of routes to be created and associated with the subnet.

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
 

