


# CORE FUNCTIONS: SUBNET




This package exports several entry points into the JS environment,
including:

> * [aws.subnets.scan](#vscan)
> * [aws.subnets.create](#vcreate)
> * [aws.subnets.delete](#vdelete)
> * [aws.subnets.describe](#vdescribe)
> * [aws.subnets.routes.create](#gcreate)
> * [aws.subnets.routes.delete](#gdelete)

This API allows resource handlers to manage subnets in an AWS VPC.

## AWS.SUBNETS.SCAN
<a name="sscan"></a>
`aws.subnets.scan(region);`

Returns a list of subnets.

Example:

```

var subnets =  aws.subnets.scan("us-east-1");

```

## AWS.SUBNETS.CREATE
<a name="screate"></a>
`aws.subnets.create(region, config);`

Create a SUBNET.

Example:

```

var subnet =  aws.subnets.create(
"us-east-1",
{
CidrBlock:        "172.33.1.0/24"
VpcId:            "vpc-abcd",
AvailabilityZone: "some zone"
});

```

## AWS.SUBNETS.DELETE
<a name="sdelete"></a>
`aws.subnets.delete(region, subnet_id);`

Delete a SUBNET.

Example:

```

aws.subnets.delete("us-east-1", "subnet-abcd");

```

## AWS.SUBNETS.DESCRIBE
<a name="sdescribe"></a>
`aws.subnets.describe(region, subnet_id);`

Get info from AWS about a SUBNET.

Example:

```

var subnet = aws.subnets.describe("us-east-1", "subnet-abcd");

```

## AWS.SUBNETS.ROUTES.CREATE
<a name="rcreate"></a>
`aws.subnets.routes.create(region);`

Create a route.

Example:

```

var route =  aws.subnets.routes.create("us-east-1", {
DestinationCidrBlock: "0.0.0.0/0",
GatewayId:            "gw-1234"
});

```

## AWS.SUBNETS.ROUTES.DELETE
<a name="rdelete"></a>
`aws.subnets.routes.delete(region, route_id);`

Delete a route.

Example:

```

aws.subnets.routes.delete("us-east-1", "route-abcd");

```


