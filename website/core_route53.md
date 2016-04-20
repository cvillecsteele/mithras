 


 # CORE FUNCTIONS: ROUTE53


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.route53.zones.scan](#zscan)

 This API allows resource handlers to manipulate DNS records in Route53.

 ## AWS.ROUTE53.ZONES.SCAN
 <a name="zscan"></a>
 `aws.route53.zones.scan(region);`

 Query zones in Route53.

 Example:

 ```

 var zones = aws.route53.zones.scan("us-east-1");

 ```

 ## AWS.ROUTE53.RRS.SCAN
 <a name="scan"></a>
 `aws.route53.rrs.scan(region);`

 Query resource records sets (RRSs) in Route53.

 Example:

 ```

 var rrs = aws.route53.rrs.scan("us-east-1");

 ```

 ## AWS.ROUTE53.RRS.CREATE
 <a name="create"></a>
 `aws.route53.rrs.create(region, zoneId, config);`

 Create dns records.

 Example:

 ```

 var rrs = aws.route53.rrs.create("us-east-1", "Z111111...",
 {
 		Name:         "mithras.io."
 		Type:         "A"
 		AliasTarget: {
 		    DNSName:              "s3-website-us-east-1.amazonaws.com"
 		    EvaluateTargetHealth: false
 		    HostedZoneId:         "Z3AQBSTGFYJSTF"
 		}
 });

 ```

 ## AWS.ROUTE53.RRS.DELETE
 <a name="delete"></a>
 `aws.route53.rrs.delete(region, zoneId, config);`

 Delete dns records.

 Example:

 ```

 var rrs = aws.route53.rrs.delete("us-east-1", "Z111111...",
 {
 		Name:         "mithras.io."
 		Type:         "A"
 		AliasTarget: {
 		    DNSName:              "s3-website-us-east-1.amazonaws.com"
 		    EvaluateTargetHealth: false
 		    HostedZoneId:         "Z3AQBSTGFYJSTF"
 		}
 });

 ```

 ## AWS.ROUTE53.RRS.DESCRIBE
 <a name="describe"></a>
 `aws.route53.rrs.describe(region, zoneId, name, type);`

 Get info about dns records.

 Example:

 ```

 var rrs = aws.route53.rrs.describe("us-east-1", "Z111111...", "mithras.io." "A");

 ```


