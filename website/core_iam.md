


# CORE FUNCTIONS: IAM




This package exports entry points into the JS environment:

> * [aws.iam.profiles.scan](#scan)
> * [aws.iam.profiles.delete](#delete)
> * [aws.iam.profiles.create](#create)
> * [aws.iam.profiles.describe](#describe)

> * [aws.iam.roles.scan](#scanRole)
> * [aws.iam.roles.delete](#deleteRole)
> * [aws.iam.roles.create](#createRole)
> * [aws.iam.roles.describe](#describeRole)
> * [aws.iam.roles.putRolePolicy](#putPolicy)
> * [aws.iam.roles.deleteRolePolicy](#deletePolicy)
> * [aws.iam.roles.addRoleToProfile](#add)
> * [aws.iam.roles.removeRoleFromProfile](#remove)

This API allows the caller to work with IAM profiles.

## AWS.IAM.PROFILES.SCAN
<a name="scan"></a>
`aws.iam.profiles.scan(region);`

Scan AWS for IAM profiles

Example:

```

var profiles = aws.iam.profiles.scan("us-east-1");

```

## AWS.IAM.PROFILES.DELETE
<a name="delete"></a>
`aws.iam.profiles.delete(region, id);`

Delete a profile.

Example:

```

aws.iam.profiles.delete("us-east-1", "my-profile");

```

## AWS.IAM.PROFILES.CREATE
<a name="create"></a>
`aws.iam.profiles.create(region, id);`

Create a profile.

Example:

```

var profile = aws.iam.profiles.create("us-east-1", "my-profile");

```

## AWS.IAM.PROFILES.DESCRIBE
<a name="describe"></a>
`aws.iam.profiles.describe(region, id);`

Get info about a profile

Example:

```

var p = aws.iam.profiles.describe("us-east-1", "my-profile");

```

## AWS.IAM.ROLES.SCAN
<a name="scanRole"></a>
`aws.iam.roles.scan(region);`

Scan AWS for IAM roles

Example:

```

var roles = aws.iam.roles.scan("us-east-1");

```

## AWS.IAM.ROLES.DELETE
<a name="deleteRole"></a>
`aws.iam.roles.delete(region, id);`

Delete a profile.

Example:

```

aws.iam.roles.delete("us-east-1", "my-role");

```

## AWS.IAM.ROLES.CREATE
<a name="createRole"></a>
`aws.iam.roles.create(region, id, trust);`

Create a role.

Example:

```

var role = aws.iam.roles.create("us-east-1", "my-role", aws.iam.roles.ec2TrustPolicy);

```

## AWS.IAM.ROLES.DESCRIBE
<a name="describeRole"></a>
`aws.iam.roles.describe(region, id);`

Get info about a role

Example:

```

var p = aws.iam.roles.describe("us-east-1", "my-role");

```

## AWS.IAM.ROLES.PUTROLEPOLICY
<a name="putPolicy"></a>
`aws.iam.roles.putRolePolicy(region, id, policyName, policy);`

Attach a policy to a role.

Example:

```

aws.iam.roles.putRolePolicy("us-east-1", "my-role", "s3_full_access",
{
"Version": "2012-10-17",
"Statement": [
{
"Effect": "Allow",
"Action": "s3:*",
"Resource": "*"
}
]
});

```

## AWS.IAM.ROLES.DELETEROLEPOLICY
<a name="deletePolicy"></a>
`aws.iam.roles.deleteRolePolicy(region, id, policyName);`

Delete a policy from a role.

Example:

```

aws.iam.roles.deleteRolePolicy("us-east-1", "my-role", "s3_full_access");

```

## AWS.IAM.ROLES.ADDROLETOPROFILE
<a name="add"></a>
`aws.iam.roles.addRoleToProfile(region, profileName, roleName);`

Add a role to a profile.

Example:

```

aws.iam.roles.addRoleToProfile("us-east-1", "my-profile", "my-role");

```

## AWS.IAM.ROLES.REMOVEROLEFROMPROFILE
<a name="remove"></a>
`aws.iam.roles.removeRoleFromProfile(region, profileName, roleName);`

Remove a role from a profile.

Example:

```

aws.iam.roles.removeRoleFromProfile("us-east-1", "my-profile", "my-role");

```


