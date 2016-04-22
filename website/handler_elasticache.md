 
 
 # Elasticache
 
 Elasticache is resource handler for managing AWS caches.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"elasticache"`
 
 Usage:
 
 `var elasticache = require("elasticache").init();`
 
  ## Example Resource
 
 ```javascript
 var rCache = {
      name: "redis"
      module: "elasticache"
      dependsOn: [otherResource.name]
      params: {
          ensure: ensure
          region: "us-east-1"
          wait: true
          subnetGroup: {
              CacheSubnetGroupDescription: "Redis Subnet Group"
              CacheSubnetGroupName:        "redis-subnet-group"
              SubnetIds: [
                  "subnet-123",
                  "subnet-456"
              ]
          }
          cache: {
              CacheClusterId:          "test-redis"
              AutoMinorVersionUpgrade: true
              CacheNodeType:           "cache.t2.small"
              CacheSubnetGroupName:    "redis-subnet-group"
              Engine:                  "redis"
              NumCacheNodes:           1
              SecurityGroupIds:        []
              Tags: [
                  {
                      Key:   "Name"
                      Value: "test-cluster"
                  },
              ]
          }
          delete: {
              CacheClusterId:          "test-redis"
          }

      }
 };
 ```
 
 ## Copy Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"`, the cache specified by `cache` will be created, and
 if `"absent"`, the cache will be removed using the `delete` property.
 
 ### `region`

 * Required: true
 * Allowed Values: string, any valid AWS region; eg "us-east-1"

 The region for calls to the AWS API.
 
 ### `wait`

 * Required: false
 * Allowed Values: true or false

 If `true`, delay execution until the cache has been created in AWS.
 
 ### `subnetGroup`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticache.html#type-CreateCacheSubnetGroupInput)

 If set, a subnet group will be created for your cache.
 
 ### `cache`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticache.html#type-CreateCacheInput)

 Parameters for cache creation.
 
 ### `delete`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticache.html#type-DeleteCacheInput)

 Parameters for cache deletion.
 
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
 

