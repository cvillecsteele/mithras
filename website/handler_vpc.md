 
 
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


