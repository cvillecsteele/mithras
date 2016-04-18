 
 
 # s3
 
 S3 is resource handler for working with AWS's S3.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"s3"`
 
 Usage:
 
 `var s3 = require("s3").init();`
 
  ## Example Resource
 
 ```javascript
 var rElbDnsEntry = {
     name: "elbDnsEntry"
     module: "s3"
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
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/s3.html#type-ResourceRecordSet)


