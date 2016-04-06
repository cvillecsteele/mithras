package iam

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "iam"

func createProfile(region string, name string, verbose bool) *iam.InstanceProfile {
	svc := iam.New(session.New(), aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	}

	if verbose {
		log.Printf("  ### Creating IAM instance profile '%s'", name)
	}

	resp, err := svc.CreateInstanceProfile(params)
	if err != nil {
		log.Fatalf("Error creating iam instance profile: %s", err)
	}
	id := *resp.InstanceProfile.InstanceProfileId

	// Wait for it.
	for i := 0; i < 10; i++ {
		p := describeProfile(region, id)
		if p != nil {
			break
		}
		time.Sleep(time.Second * 10)
	}

	return describeProfile(region, id)
}

func deleteProfile(region string, id string, verbose bool) {
	svc := iam.New(session.New(), aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	if verbose {
		log.Printf("  ### Deleting IAM profile '%s'\n", id)
	}

	params := &iam.DeleteInstanceProfileInput{
		InstanceProfileName: aws.String(id),
	}
	_, err := svc.DeleteInstanceProfile(params)

	if err != nil {
		log.Fatal(err.Error())
	}

}

func describeProfile(region string, id string) *iam.InstanceProfile {
	svc := iam.New(session.New(), aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(id),
	}
	resp, err := svc.GetInstanceProfile(params)

	if err != nil {
		return nil
	}

	return resp.InstanceProfile
}

func scanProfiles(region string) []*iam.InstanceProfile {
	svc := iam.New(session.New(), aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &iam.ListInstanceProfilesInput{}
	resp, err := svc.ListInstanceProfiles(params)

	if err != nil {
		panic(err)
	}

	return resp.InstanceProfiles
}

func createRole(region string, name string, trustPolicy string, verbose bool) *iam.Role {
	svc := iam.New(session.New(), aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(trustPolicy),
		RoleName:                 aws.String(name),
		// Path:                     aws.String("pathType"),
	}

	if verbose {
		log.Printf("  ### Creating IAM role '%s'", name)
	}

	resp, err := svc.CreateRole(params)
	if err != nil {
		log.Fatalf("Error creating IAM role: %s", err)
	}
	id := *resp.Role.RoleId

	// Wait for it.
	for i := 0; i < 10; i++ {
		p := describeRole(region, id)
		if p != nil {
			break
		}
		time.Sleep(time.Second * 10)
	}

	return describeRole(region, id)
}

func deleteRole(region string, id string, verbose bool) {
	svc := iam.New(session.New(), aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	if verbose {
		log.Printf("  ### Deleting IAM role '%s'\n", id)
	}

	params := &iam.DeleteRoleInput{
		RoleName: aws.String(id),
	}
	_, err := svc.DeleteRole(params)

	if err != nil {
		log.Fatal(err.Error())
	}

}

func describeRole(region string, id string) *iam.Role {
	svc := iam.New(session.New(), aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &iam.GetRoleInput{
		RoleName: aws.String(id),
	}
	resp, err := svc.GetRole(params)

	if err != nil {
		return nil
	}

	return resp.Role
}

func scanRoles(region string) []*iam.Role {
	svc := iam.New(session.New(), aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &iam.ListRolesInput{}
	resp, err := svc.ListRoles(params)

	if err != nil {
		panic(err)
	}

	return resp.Roles
}

func putRolePolicy(region string, roleName string, policyName string, policy string) {
	svc := iam.New(session.New(), aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(policy),
		PolicyName:     aws.String(policyName),
		RoleName:       aws.String(roleName),
	}
	_, err := svc.PutRolePolicy(params)

	if err != nil {
		panic(err)
	}

}

func addRoleToProfile(region string, profileName string, roleName string) {
	svc := iam.New(session.New(), aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
		RoleName:            aws.String(roleName),
	}
	_, err := svc.AddRoleToInstanceProfile(params)

	if err != nil {
		panic(err)
	}
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		var o1 *otto.Object
		var awsObj *otto.Object
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			awsObj, _ = rt.Object(`aws = {}`)
		} else {
			awsObj = a.Object()
		}

		if b, err := awsObj.Get("iam"); err != nil || b.IsUndefined() {
			o1, _ = rt.Object(`aws.iam = {}`)
		} else {
			o1 = b.Object()
		}

		// Profiles
		var o2 *otto.Object
		if b, err := o1.Get("profiles"); err != nil || b.IsUndefined() {
			o2, _ = rt.Object(`aws.iam.profiles = {}`)
		} else {
			o2 = b.Object()
		}
		o2.Set("scan", func(region string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(scanProfiles(region))
		})
		o2.Set("delete", func(region, name string) otto.Value {
			verbose := mcore.IsVerbose(rt)
			deleteProfile(region, name, verbose)
			return otto.Value{}
		})
		o2.Set("create", func(region, name string) otto.Value {
			verbose := mcore.IsVerbose(rt)
			f := mcore.Sanitizer(rt)
			return f(createProfile(region, name, verbose))
		})
		o2.Set("describe", func(region, id string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(describeProfile(region, id))
		})

		var o3 *otto.Object
		if b, err := o1.Get("roles"); err != nil || b.IsUndefined() {
			o3, _ = rt.Object(`aws.iam.roles = {}`)
		} else {
			o3 = b.Object()
		}
		o3.Set("scan", func(region string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(scanRoles(region))
		})
		o3.Set("delete", func(region, id string) otto.Value {
			verbose := mcore.IsVerbose(rt)
			deleteRole(region, id, verbose)
			return otto.Value{}
		})
		o3.Set("create", func(region, name, trust string) otto.Value {
			verbose := mcore.IsVerbose(rt)
			f := mcore.Sanitizer(rt)
			return f(createRole(region, name, trust, verbose))
		})
		o3.Set("describe", func(region, name string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(describeRole(region, name))
		})
		o3.Set("putRolePolicy", func(region, roleName string, policyName string, policy string) otto.Value {
			putRolePolicy(region, roleName, policyName, policy)
			return otto.Value{}
		})
		o3.Set("addRoleToProfile", func(region, profileName string, roleName string) otto.Value {
			addRoleToProfile(region, profileName, roleName)
			return otto.Value{}
		})
		o3.Set("ec2TrustPolicy", `{
      "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "",
          "Effect": "Allow",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          },
          "Action": "sts:AssumeRole"
        }
      ]
    }`)

	})
}
