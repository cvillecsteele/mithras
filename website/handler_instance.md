 
 
 # Instance
 
 Instance is resource handler for managing AWS caches.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"instance"`
 
 Usage:
 
 `var instance = require("instance").init();`
 
  ## Example Resource
 
 ```javascript
 var webserverTagName: "my-web-server";
 var keyname: "test-key";
 var webServer = {
      name: "webserver"
      module: "instance"
      dependsOn: [resource.name]
      params: {
          region: "us-east-1"
          ensure: "present"
          on_find: function(catalog) {
              var matches = _.filter(catalog.instances, function (i) {
                  if (i.State.Name != "running") {
                      return false;
                  }
                  return (_.where(i.Tags, {"Key": "Name", 
                                           "Value": webserverTagName}).length > 0);
              });
              return matches;
          }
          tags: {
              Name: webserverTagName
          }
          instance: {
              ImageId:        "ami-60b6c60a"
              MaxCount:       1
              MinCount:       1
              DisableApiTermination: false
              EbsOptimized:          false
              IamInstanceProfile: {
                  Name: iamProfileName
              }
              InstanceInitiatedShutdownBehavior: "terminate"
              InstanceType:                      "t2.small"
              KeyName:                           keyName
              Monitoring: {
                  Enabled: true
              }
              NetworkInterfaces: [
                  {
                      AssociatePublicIpAddress: true
                      DeleteOnTermination:      true
                      DeviceIndex:              0
                      Groups:                  [ "sg-abc" ]
                      SubnetId:                "subnet-123"
                  }
              ]
          } // instance
      } // params
 };
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"`, the db specified by `db` will be created, and
 if `"absent"`, it will be removed using the `delete` property.
 
 ### `region`

 * Required: true
 * Allowed Values: string, any valid AWS region; eg "us-east-1"

 The region for calls to the AWS API.
 
 ### `instance`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-RunInstancesInput)

 Parameters for resource creation.
 
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
 
 ### `tags`

 * Required: false
 * Allowed Values: A map of tags to be applied to created instances

 A map of tags to be applied to created instances
 

