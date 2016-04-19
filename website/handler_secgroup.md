 
 
 # secgroup
 
 Secgroup is a resource handler for dealing with AWS security groups.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"secgroup"`
 
 Usage:
 
 `var secgroup = require("secgroup").init();`
 
  ## Example Resource
 
 ```javascript
 var sgName = "simple-sg";
 var security = {
      name: "webserverSG"
      module: "secgroup"
      params: {
          region: defaultRegion
          ensure: ensure

          secgroup: {
              Description: "Webserver security group"
              GroupName:   sgName
          }
          tags: {
              Name: "webserver"
          }
          ingress: {
              IpPermissions: [
                  {
                      FromPort:   22
                      IpProtocol: "tcp"
                      IpRanges: [ {CidrIp: "0.0.0.0/0"} ]
                      ToPort: 22
                  },
                  {
                      FromPort:   80
                      IpProtocol: "tcp"
                      IpRanges: [ {CidrIp: "0.0.0.0/0"} ]
                      ToPort: 80
                  }
              ]
          }
      }
 };
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"` and the security group named
 `params.secgroup.GroupName` does not exist, it is created.  If
 `"absent"`, and the group exists, it is removed.
 
 ### `secgroup`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-CreateSecurityGroupInput)

 This file is copied to remote hosts.

 ### `tags`

 * Required: false
 * Allowed Values: a map of tags to be set on a created security group

 For tagging.

 ### `ingress`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-AuthorizeSecurityGroupIngressInput)

 Set inbound rules for the SG

 ### `egress`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-AuthorizeSecurityGroupEgressInput)

 Set outbound rules for the SG


