// MITHRAS: Javascript configuration management tool for AWS.
// Copyright (C) 2016, Colin Steele
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//   the Free Software Foundation, either version 3 of the License, or
//                  (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//              GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// @public
// 
// # IAM
// 
// Iam is resource handler for manipulating AWS IAM resources
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"iam"`
// 
// Usage:
// 
// `var iam = require("iam").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var iamProfileName = "my-instance-profile";
// var iamRoleName = "test-webserver-iam-role";
// var rIAM = {
//      name: "IAM"
//      module: "iamProfile"
//      dependsOn: [otherResource.name]
//      params: {
//          region: "us-east-1"
//          ensure: "present"
//          profile: {
//              InstanceProfileName: iamProfileName
//          }
//          role: {
//              RoleName: iamRoleName
//              AssumeRolePolicyDocument: aws.iam.roles.ec2TrustPolicy
//          }
//          policies: {
//              "s3_full_access": {
//                  "Version": "2012-10-17",
//                  "Statement": [
//                      {
//                          "Effect": "Allow",
//                          "Action": "s3:*",
//                          "Resource": "*"
//                      }
//                  ]
//              },
//          }
//      }
// }
// ```
// 
// ## Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "absent", "present"
//
// If `"present"`, the handler will ensure the profile exists, and it
// not, it will be created.  If `"absent"`, the profile is removed.
// 
// ### `profile`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/iam.html#type-CreateInstanceProfileInput)
//
// Specifies parameters for profile creation.
//
// ### `role`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/iam.html#type-CreateRoleInput)
//
// ### `policies`
//
// * Required: false
// * Allowed Values: map of policyname => IAM policy document
//
// Adds an inline policy document that is embedded in the specified role.
// 
// When you embed an inline policy in a role, the inline policy is
// used as the role's access (permissions) policy. The role's trust
// policy is created at the same time as the role.  For more
// information about roles, go to [Using Roles to Delegate Permissions
// and Federate
// Identities](http://docs.aws.amazon.com/IAM/latest/UserGuide/roles-toplevel.html).
// 
// A role can also have a managed policy attached to it. Refer to
// [Managed Policies and Inline
// Policies](http://docs.aws.amazon.com/IAM/latest/UserGuide/policies-managed-vs-inline.html)
// in the IAM User Guide.
// 
// For information about limits on the number of inline policies that
// you can embed with a role, see [Limitations on IAM
// Entities](http://docs.aws.amazon.com/IAM/latest/UserGuide/LimitationsOnEntities.html)
// in the IAM User Guide.
//
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
        module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
        moduleNames: ["iam", "iamProfile", "iamRole"]
        findProfile: function(catalog, resource, name) {
            if (typeof(resource.params.on_find) === 'function') {
		result = resource.params.on_find(catalog, resource);
		if (!result || 
		    (Array.isArray(result) && result.length == 0)) {
		    return;
		}
		return result;
	    }
            return _.find(catalog.iamProfiles, function(p) {
                return p.InstanceProfileName === name;
            });
        }
        findRole: function(catalog, name) {
            return _.find(catalog.iamRoles, function(p) {
                return p.RoleName === name;
            });
        }
        handle: function(catalog, resources, resource) {
            if (!_.find(handler.moduleNames, function(m) { 
                return resource.module === m; 
            })) {
                return [null, false];
            }
                
            var ensure = resource.params.ensure;
            var params = resource.params;
            
            switch(resource.module) {
            case "iamProfile":
                // Sanity
                if (!resource.params.profile) {
                    console.log("Invalid iamProfile params")
                    exit(1);
                }
                var p = resource._target;
                
                switch(ensure) {
                case "absent":
                    if (!p) {
                        if (mithras.verbose) {
                            log(sprintf("IAM Profile not found, no action taken."));
                        }
                        break;
                    }
                    
                    // Role?
                    var roleName = params.role.RoleName;
                    var role = aws.iam.roles.describe(params.region, 
                                                      roleName);
                    if (role) {
                        var profileName = params.profile.InstanceProfileName;

                        // Remove from profile first
                        if (mithras.verbose) {
                            log("Removing role from profile");
                        }
                        aws.iam.roles.removeRoleFromProfile(params.region,
                                                            profileName,
                                                            roleName);

                        // delete policies from role
                        var roleName = params.role.RoleName;
                        _.each(params.policies, function(policy, name) {
                            if (mithras.verbose) {
                                log("Deleting role policies");
                            }
                            aws.iam.roles.deleteRolePolicy(params.region,
                                                           roleName,
                                                           name)
                        });

                        // nuke it
                        if (mithras.verbose) {
                            log(sprintf("Deleting role '%s'", roleName));
                        }
                        aws.iam.roles.delete(params.region, roleName);

                        // nix from catalog
                        catalog.iamRoles = _.reject(catalog.iamRoles, 
                                                    function(x) { 
                                                        return x.RoleName == roleName;
                                                    });

                    } else {
                        log("No role to delete.")
                    }
                    
                    // Remove it
                    var profileName = params.profile.InstanceProfileName;
                    aws.iam.profiles.delete(params.region, profileName);

                    // nix from catalog
                    catalog.iamProfiles = 
                        _.reject(catalog.iamProfiles, 
                                 function(x) { 
                                     return x.InstanceProfileName == profileName;
                                 });
                    
                    break;
                case "present":
                    if (p) {
                        log(sprintf("IAM profile found, no action taken."));
                        break;
                    }
                    
                    // create 
                    var profileName = params.profile.InstanceProfileName;
                    if (mithras.verbose) {
                        log(sprintf("Creating IAM instance profile '%s'", 
                                    profileName));
                    }
                    created = aws.iam.profiles.create(params.region, profileName);
                    
                    // create role
                    var roleName = params.role.RoleName;
                    var trust = params.role.AssumeRolePolicyDocument;
                    if (mithras.verbose) {
                        log(sprintf("Creating IAM role '%s'", 
                                    roleName));
                    }
                    role = aws.iam.roles.create(params.region, 
                                                roleName,
                                                trust);
                    
                    // add policy to role
                    _.each(params.policies, function(policy, name) {
                        if (mithras.verbose) {
                            log(sprintf("Putting policy '%s' to IAM role '%s'", 
                                        name,
                                        roleName));
                        }
                        aws.iam.roles.putRolePolicy(params.region,
                                                    roleName,
                                                    name,
                                                    JSON.stringify(policy));
                    });

                    // stick the role to the profile
                    if (mithras.verbose) {
                        log(sprintf("Adding role '%s' to IAM profile '%s'", 
                                    roleName,
                                    profileName));
                    }
                    aws.iam.roles.addRoleToProfile(params.region,
                                                   profileName,
                                                   roleName);

                    // Wait for association between profile and role
                    var profile = aws.iam.profiles.describe(params.region, profileName);
                    do {
                        time.sleep(10);
                        profile = aws.iam.profiles.describe(params.region, profileName);
                    } while (!profile || !profile.Roles || 
                             (profile.Roles.length == 0) || 
                             (profile.Roles[0].RoleName != roleName))

                    // add to catalog
                    catalog.iamProfiles.push(profile);
                    catalog.iamRoles.push(aws.iam.roles.describe(params.region, 
                                                                 roleName));

                    // return profile
                    return [profile, true];
                }
                return [null, true];
                break;
            case "iam":
                return [null, true];
                break;
            }
        }
        preflight: function(catalog, resources, resource) {
            if (!_.find(handler.moduleNames, function(m) { 
                return resource.module === m; 
            })) {
                return [null, false];
            }
            // Profiles
            if (resource.module === handler.moduleNames[1]) {
                var params = resource.params;
                var profile = params.profile;
                var p = handler.findProfile(catalog, resource, profile.InstanceProfileName);
                if (p) {
                    return [p, true];
                }
            }
            return [null, true];
        }
    };
    
    handler.init = function () {
        mithras.modules.preflight.register(handler.moduleNames[0], handler.preflight);
        mithras.modules.handlers.register(handler.moduleNames[0], handler.handle);
        return handler;
    };
    
    return handler;
});
