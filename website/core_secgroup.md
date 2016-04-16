 


 # CORE FUNCTIONS: SECGROUP


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.securityGroups.scan](#scan)
 > * [aws.securityGroups.create](#create)
 > * [aws.securityGroups.delete](#delete)
 > * [aws.securityGroups.describe](#describe)
 > * [aws.securityGroups.authorizeIngress](#ingress)
 > * [aws.securityGroups.authorizeEgress](#egress)

 This API allows resource handlers to manage secgroups in an AWS VPC.

 ## AWS.SECURITYGROUPS.SCAN
 <a name="scan"></a>
 `aws.securityGroups.scan(region);`

 Returns a list of security groups.

 Example:

 ```

  var secgroups =  aws.securityGroups.scan("us-east-1");

 ```

 ## AWS.SECURITYGROUPS.CREATE
 <a name="create"></a>
 `aws.securityGroups.create(region, config);`

 Create a security group.

 Example:

 ```

  var secgroup =  aws.securityGroups.create(
    "us-east-1",
    {
      Description: "Webserver security group"
      GroupName:   "webserver"
      VpcId:       "vpc-xyz"
    });

 ```

 ## AWS.SECURITYGROUPS.DELETE
 <a name="delete"></a>
 `aws.securityGroups.delete(region, secgroup_id);`

 Delete a security group.

 Example:

 ```

  aws.securityGroups.delete("us-east-1", "sg-abcd");

 ```

 ## AWS.SECURITYGROUPS.DESCRIBE
 <a name="describe"></a>
 `aws.securityGroups.describe(region, secgroup_id);`

 Get info from AWS about a security group.

 Example:

 ```

  var secgroup = aws.securityGroups.describe("us-east-1", "sg-abcd");

 ```

 ## AWS.SECURITYGROUPS.AUTHORIZEINGRESS
 <a name="ingress"></a>
 `aws.securityGroups.authorizeIngress(region, permissions);`

 Authorize ingress routes for a security group.

 Example:

 ```

 aws.securityGroups.authorizeIngress("us-east-1", {
   GroupId: "sg-xyz"
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
 });

 ```

 ## AWS.SECURITYGROUPS.AUTHORIZEEGRESS
 <a name="egress"></a>
 `aws.securityGroups.authorizeEgress(region, permissions);`

 Authorize egress routes for a security group.

 Example:

 ```

 aws.securityGroups.authorizeEgress("us-east-1", {
   GroupId: "sg-xyz"
   IpPermissions: [
     {
       FromPort:   0
       IpProtocol: "tcp"
       IpRanges: [ {CidrIp: "0.0.0.0/0"} ]
       ToPort: 65535
     },
   ]
 });

 ```


