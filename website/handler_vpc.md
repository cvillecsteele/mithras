 
 
 # vpc
 
 Vpc is a resource handler for dealing with AWS security groups.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"vpc"`
 
 Usage:
 
 `var vpc = require("vpc").init();`
 
  ## Example Resource
 
 ```javascript
 var rVpc = {
      name: "VPC"
      module: "vpc"
      params: {
          region: defaultRegion
          ensure: ensure
          vpc: {
              CidrBlock:       "172.33.0.0/16"
          }
          gateway: true
          tags: {
              Name: "my-vpc"
          }
      }
 };
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"` and the vpc named
 `params.vpc.GroupName` does not exist, it is created.  If
 `"absent"`, and the vpc exists, it is removed.
 
 ### `vpc`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-CreateVpcInput)

 Parameters for VPC creation.

 ### `tags`

 * Required: false
 * Allowed Values: a map of tags to be set on a created vpc

 For tagging.

 ### `gateway`

 * Required: false
 * Allowed Values: true or false

 If true, an internet gateway is created for the VPC.

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
 

