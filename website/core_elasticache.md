


# CORE FUNCTIONS: ELASTICACHE




This package exports several entry points into the JS environment,
including:

> * [aws.elasticache.create](#create)
> * [aws.elasticache.delete](#delete)
> * [aws.elasticache.describe](#describe)
> * [aws.elasticache.scan](#scan)
> * [aws.elasticache.subnetGroups.create](#sgcreate)
> * [aws.elasticache.subnetGroups.delete](#sgdelete)
> * [aws.elasticache.subnetGroups.describe](#sgdescribe)

This API allows exposes functions to manage AWS cache clusters and
their associated subnet groups.

## AWS.ELASTICACHE.CREATE
<a name="create"></a>
`aws.elasticache.create(region, config);`

Create a cache cluster.

Example:

```

var cache = aws.elasticache.create("us-east-1", {
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
});

```

## AWS.ELASTICACHE.DELETE
<a name="delete"></a>
`aws.elasticache.delete(region, cache_id);`

Delete a cache cluster.

Example:

```

aws.elasticache.delete("us-east-1", "test-redis");

```

## AWS.ELASTICACHE.DESCRIBE
<a name="describe"></a>
`aws.elasticache.describe(region, cache_id);`

Get information about a cache cluster.

Example:

```

var cache = aws.elasticache.describe("us-east-1", "test-redis");

```

## AWS.ELASTICACHE.SCAN
<a name="scan"></a>
`aws.elasticache.scan(region, cache_id);`

Get information about all cache clusters.

Example:

```

var caches = aws.elasticache.scan("us-east-1");

```

## AWS.ELASTICACHE.SUBNETGROUPS.CREATE
<a name="sgcreate"></a>
`aws.elasticache.subnetGroups.create(region, config);`

Create an elasticache subnet group.

Example:

```

var cache = aws.elasticache.create("us-east-1", {
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
});

```

## AWS.ELASTICACHE.SUBNETGROUPS.DELETE
<a name="sgdelete"></a>
`aws.elasticache.subnetGroups.delete(region, cache_id);`

Delete a cache subnet group.

Example:

```

aws.elasticache.subnetGroups.delete("us-east-1", "redis-subnet-group");

```

## AWS.ELASTICACHE.SUBNETGROUPS.DESCRIBE
<a name="sgdescribe"></a>
`aws.elasticache.subnetGroups.describe(region, cache_id);`

Get information about a cache cluster.

Example:

```

var group = aws.elasticache.subnetGroups.describe("us-east-1", "redis-subnet-group");

```


