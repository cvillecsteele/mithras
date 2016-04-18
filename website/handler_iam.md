 
 
 # IAM
 
 Iam is resource handler for manipulating AWS IAM resources
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"iam"`
 
 Usage:
 
 `var iam = require("iam").init();`
 
  ## Example Resource
 
 ```javascript
 var iamProfileName = "my-instance-profile";
 var iamRoleName = "test-webserver-iam-role";
 var rIAM = {
      name: "IAM"
      module: "iamProfile"
      dependsOn: [otherResource.name]
      params: {
          region: "us-east-1"
          ensure: "present"
          profile: {
              InstanceProfileName: iamProfileName
          }
          role: {
              RoleName: iamRoleName
              AssumeRolePolicyDocument: aws.iam.roles.ec2TrustPolicy
          }
          policies: {
              "s3_full_access": {
                  "Version": "2012-10-17",
                  "Statement": [
                      {
                          "Effect": "Allow",
                          "Action": "s3:*",
                          "Resource": "*"
                      }
                  ]
              },
          }
      }
 }
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "absent", "present"

 If `"present"`, the handler will ensure the profile exists, and it
 not, it will be created.  If `"absent"`, the profile is removed.
 
 ### `profile`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/iam.html#type-CreateInstanceProfileInput)

 Specifies parameters for profile creation.

 ### `role`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/iam.html#type-CreateRoleInput)

 ### `policies`

 * Required: false
 * Allowed Values: map of policyname => IAM policy document

 Adds an inline policy document that is embedded in the specified role.
 
 When you embed an inline policy in a role, the inline policy is
 used as the role's access (permissions) policy. The role's trust
 policy is created at the same time as the role.  For more
 information about roles, go to [Using Roles to Delegate Permissions
 and Federate
 Identities](http://docs.aws.amazon.com/IAM/latest/UserGuide/roles-toplevel.html).
 
 A role can also have a managed policy attached to it. Refer to
 [Managed Policies and Inline
 Policies](http://docs.aws.amazon.com/IAM/latest/UserGuide/policies-managed-vs-inline.html)
 in the IAM User Guide.
 
 For information about limits on the number of inline policies that
 you can embed with a role, see [Limitations on IAM
 Entities](http://docs.aws.amazon.com/IAM/latest/UserGuide/LimitationsOnEntities.html)
 in the IAM User Guide.


