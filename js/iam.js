(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleNames: ["iam", "iamProfile", "iamRole"]
	findProfile: function(catalog, name) {
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
		    role = aws.iam.roles.create(params.region, 
						roleName,
						trust);
		    
		    // add policy to role
		    _.each(params.policies, function(policy, name) {
			aws.iam.roles.putRolePolicy(params.region,
						    roleName,
						    name,
						    JSON.stringify(policy));
		    });

		    // stick the role to the profile
		    aws.iam.roles.addRoleToProfile(params.region,
						   profileName,
						   roleName);
		    
		    // add to catalog
		    catalog.iamProfiles.push(aws.iam.profiles.describe(params.region, 
								       profileName));
		    catalog.iamRoles.push(aws.iam.roles.describe(params.region, 
								 roleName));

		    // return profile
		    return [created, true];
		}
		return [null, true];
		break;
	    case "iam":
		return [null, true];
		break;
	    }
	}
	preflight: function(catalog, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
	    // Profiles
	    if (resource.module === handler.moduleNames[1]) {
		var params = resource.params;
		var profile = params.profile;
		var p = handler.findProfile(catalog, profile.InstanceProfileName);
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
