 
 
 # ELB
 
 Elb is resource handler for managing AWS elastic load balancers.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"elb"`
 
 Usage:
 
 `var elb = require("elb").init();`
 
  ## Example Resource
 
 ```javascript
 var lbName = "my-lb";
 var rElb = {
      name: "elb"
      module: "elb"
      dependsOn: [otherResource.name]
      on_delete: function(elb) { 
          // Sometimes aws takes a bit to delete an elb, and we can't
          // proceed with deleting until it's GONE.
          this.delay = 30; 
          return true;
      }
      params: {
          region: "us-east-1"
          ensure: "present"
          elb: {
              Listeners: [
                  {
                      InstancePort:     80
                      LoadBalancerPort: 80
                      Protocol:         "http"
                      InstanceProtocol: "http"
                  },
              ]
              LoadBalancerName: lbName
              SecurityGroups: [ "sg-123" ]
              Subnets: [
                  "subnet-abc"
                  "subnet-def"
              ]
              Tags: [
                  {
                      Key:   "foo"
                      Value: "bar"
                  },
              ]
          }
          attributes: {
              LoadBalancerAttributes: {
                  AccessLog: {
                      Enabled:        false
                      EmitInterval:   60
                      S3BucketName:   "my-loadbalancer-logs"
                      S3BucketPrefix: "test-app"
                  }
                  ConnectionDraining: {
                      Enabled: true
                      Timeout: 300
                  }
                  ConnectionSettings: {
                      IdleTimeout: 30
                  }
                  CrossZoneLoadBalancing: {
                      Enabled: true
                  }
              }
              LoadBalancerName: lbName
          }
          health: {
              HealthCheck: {
                  HealthyThreshold:   2
                  Interval:           30
                  Target:             "HTTP:80/hc"
                  Timeout:            5
                  UnhealthyThreshold: 3
              }
              LoadBalancerName: lbName
          }
      }
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
 
 ### `elb`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elb.html#type-CreateLoadBalanerInput)

 Parameters for resource creation.
 
 ### `on_delete`

 * Required: false
 * Allowed Values: `function(elb) { ... }`

 Called after resource deletion.  May be used to modify `wait`
 property (see example above).
 
 ### `attributes`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elb.html#type-ModifyLoadBalancerAttributesInput)

 If specified, used to apply attributes to a created ELB.
 
 ### `health`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elb.html#type-ConfigureHealthCheckInput)

 If specified, used to specify health check params for a created ELB.
 
 ### `on_find`

 * Required: false
 * Allowed Values: A function taking two parameters: `catalog` and `resource`

 If defined in the resource's `params` object, the `on_find`
 function provides a way for a matching resource to be identified
 using a user-defined way.  The function is called with the current
 `catalog`, as well as the `resource` object itself.  The function
 can look through the catalog, find a matching object using whatever
 logic you want, and return it.  If the function returns `undefined`
 or a n empty Javascript array, (`[]`), the function is indicating
 that no matching resource was found in the `catalog`.
 

