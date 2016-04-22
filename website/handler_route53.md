 
 
 # route53
 
 Route53 is resource handler for working with AWS's route53 dns system.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"route53"`
 
 Usage:
 
 `var route53 = require("route53").init();`
 
  ## Example Resource
 
 ```javascript
 var rElbDnsEntry = {
     name: "elbDnsEntry"
     module: "route53"
     dependsOn: [rElb.name]
     params: {
         region: defaultRegion
         ensure: ensure
         domain: "mithras.io."
         resource: {
             Name: "test.mithras.io."
             Type: "A"
             AliasTarget: {
                 DNSName:              mithras.watch("elb._target.DNSName")
                 EvaluateTargetHealth: true
                 HostedZoneId:         mithras.watch("elb._target.CanonicalHostedZoneNameID")
             }
         }
     } // params
 };
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present", "absent"

 If `"present"`, the dns entry will be created if it doesn't already
 exist.  If `"absent"`, the dns entry will be removed if it is
 present.
 
 ### `resource`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/route53.html#type-ResourceRecordSet)

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
 

