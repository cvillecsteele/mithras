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
//
// # CORE FUNCTIONS: AUTOSCALING
//

package autoscaling

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.autoscaling.groups.scan](#gscan)
// > * [aws.autoscaling.groups.create](#gcreate)
// > * [aws.autoscaling.groups.delete](#gdelete)
// > * [aws.autoscaling.groups.describe](#gdescribe)
//
// > * [aws.autoscaling.hooks.scan](#hscan)
// > * [aws.autoscaling.hooks.create](#hcreate)
// > * [aws.autoscaling.hooks.delete](#hdelete)
// > * [aws.autoscaling.hooks.describe](#hdescribe)
// > * [aws.autoscaling.hooks.complete](#hcomplete)
// > * [aws.autoscaling.hooks.recordHeartbeat](#hrecord)
//
// > * [aws.autoscaling.launchConfigs.scan](#lscan)
// > * [aws.autoscaling.launchConfigs.create](#lcreate)
// > * [aws.autoscaling.launchConfigs.delete](#ldelete)
// > * [aws.autoscaling.launchConfigs.describe](#ldescribe)
//
// This API allows resource handlers to manage Autoscaling groups.
//
// ## AWS.AUTOSCALING.GROUPS.SCAN
// <a name="gscan"></a>
// `aws.autoscaling.groups.scan(region);`
//
// Returns a list of autoscaling groups.
//
// Example:
//
// ```
//
//  var groups = aws.autoscaling.groups.scan("us-east-1");
//
// ```
//
// ## AWS.AUTOSCALING.GROUPS.DESCRIBE
// <a name="gdescribe"></a>
// `aws.autoscaling.groups.desribe(region, groupName);`
//
// Get info about an autoscaling group.
//
// Example:
//
// ```
//
//  var group = aws.autoscaling.groups.describe("us-east-1", "groupName");
//
// ```
//
// ## AWS.AUTOSCALING.GROUPS.CREATE
// <a name="gcreate"></a>
// `aws.autoscaling.groups.create(region, config);`
//
// Create an autoscaling group.
//
// Example:
//
// ```
//
//  var scaling =  aws.autoscaling.groups.create(
//    "us-east-1",
// {
//   AutoScalingGroupName: "name"
//   MaxSize:              1
//   MinSize:              1
//   AvailabilityZones: [ "zone" ]
//   DefaultCooldown:         100
//   DesiredCapacity:         5
//   HealthCheckGracePeriod:  100
//   HealthCheckType:         "..."
//   InstanceId:              "..."
//   LaunchConfigurationName: "ResourceName"
//   LoadBalancerNames: [ "..." ]
//   NewInstancesProtectedFromScaleIn: true
//   PlacementGroup:                   "groupName"
//   Tags: [
//     {
//       Key:               "TagKey"
//       PropagateAtLaunch: true
//       ResourceId:        "id"
//       ResourceType:      "type"
//       Value:             "TagValue"
//     },
//   ]
//   TerminationPolicies: [ "policyName" ]
//   VPCZoneIdentifier: "zoneId"
// });
//
// ```
//
// ## AWS.AUTOSCALING.GROUPS.DELETE
// <a name="gdelete"></a>
// `aws.autoscaling.groups.delete(region, groupName);`
//
// Delete an autoscaling group.
//
// Example:
//
// ```
//
//  aws.autoscaling.groups.delete("us-east-1", "groupName");
//
// ```
//
// ## AWS.AUTOSCALING.HOOKS.SCAN
// <a name="hscan"></a>
// `aws.autoscaling.hooks.scan(region, groupName);`
//
// Returns a list of autoscaling hooks.
//
// Example:
//
// ```
//
//  var hooks = aws.autoscaling.hooks.scan("us-east-1", "groupName");
//
// ```
//
// ## AWS.AUTOSCALING.HOOKS.DESCRIBE
// <a name="hdescribe"></a>
// `aws.autoscaling.hooks.desribe(region, groupName, hookName);`
//
// Get info about an autoscaling hook.
//
// Example:
//
// ```
//
//  var hook = aws.autoscaling.hooks.describe("us-east-1", "groupName", "hookName");
//
// ```
//
// ## AWS.AUTOSCALING.HOOKS.CREATE
// <a name="hcreate"></a>
// `aws.autoscaling.hooks.create(region, config);`
//
// Create an autoscaling hook.
//
// Example:
//
// ```
//
//  var hook =  aws.autoscaling.hooks.create(
//    "us-east-1",
// {
//   AutoScalingGroupName:  "groupName"
//   LifecycleHookName:     "hookName"
//   DefaultResult:         "CONTINUE"
//   HeartbeatTimeout:      100
//   LifecycleTransition:   "autoscaling:EC2_INSTANCE_LAUNCHING"
//   NotificationMetadata:  "data here"
//   NotificationTargetARN: "arn:aws:sns:us-west-2:123456789012:my-sns-topic"
//   RoleARN:               "arn:aws:iam::123456789012:role/my-auto-scaling-role"
// });
//
// ```
//
// ## AWS.AUTOSCALING.HOOKS.DELETE
// <a name="hdelete"></a>
// `aws.autoscaling.hooks.delete(region, gropuName, hookName);`
//
// Delete an autoscaling hook.
//
// Example:
//
// ```
//
//  aws.autoscaling.hooks.delete("us-east-1", "groupName", "hookName");
//
// ```
//
// ## AWS.AUTOSCALING.HOOKS.COMPLETE
// <a name="hcomplete"></a>
// `aws.autoscaling.hooks.complete(region, config);`
//
// Complete a lifecycle hook.
//
// Example:
//
// ```
//
//  aws.autoscaling.hooks.complete("us-east-1",
// {
//   AutoScalingGroupName:  "groupName"
//   LifecycleActionResult: "COMPLETE"
//   LifecycleHookName:     "hookName"
//   InstanceId:            "1234"
//   LifecycleActionToken:  "LifecycleActionToken"
// });
//
// ```
//
// ## AWS.AUTOSCALING.HOOKS.RECORDHEARTBEAT
// <a name="hrecord"></a>
// `aws.autoscaling.hooks.recordHeartbeat(region, config);`
//
// Complete a lifecycle hook.
//
// Example:
//
// ```
//
//  aws.autoscaling.hooks.recordHeartbeat("us-east-1",
// {
//   AutoScalingGroupName:  "groupName"
//   LifecycleHookName:     "hookName"
//   InstanceId:            "1234"
//   LifecycleActionToken:  "LifecycleActionToken"
// });
//
// ```
//
// ## AWS.AUTOSCALING.LAUNCHCONFIGS.SCAN
// <a name="lscan"></a>
// `aws.autoscaling.launchconfigs.scan(region);`
//
// Returns a list of autoscaling launchconfigs.
//
// Example:
//
// ```
//
//  var launchconfigs = aws.autoscaling.launchConfigs.scan("us-east-1");
//
// ```
//
// ## AWS.AUTOSCALING.LAUNCHCONFIGS.DESCRIBE
// <a name="ldescribe"></a>
// `aws.autoscaling.launchconfigs.desribe(region, launchconfigName);`
//
// Get info about an autoscaling launchconfig.
//
// Example:
//
// ```
//
//  var launchconfig = aws.autoscaling.launchconfigs.describe("us-east-1", "launchconfigName");
//
// ```
//
// ## AWS.AUTOSCALING.LAUNCHCONFIGS.CREATE
// <a name="lcreate"></a>
// `aws.autoscaling.launchConfigs.create(region, config);`
//
// Create an autoscaling launchconfig.
//
// Example:
//
// ```
//
//  var config =  aws.autoscaling.launchConfigs.create(
//    "us-east-1",
// {
//     LaunchConfigurationName:  "lcName"
//     AssociatePublicIpAddress: true
//     BlockDeviceMappings: [
//      {
//         DeviceName: "string"
//         Ebs: {
//           DeleteOnTermination: true
//           Encrypted:           false
//           Iops:                100
//           SnapshotId:          "StringMaxLen255"
//           VolumeSize:          100
//           VolumeType:          "BlockDeviceEbsVolumeType"
//        }
//        NoDevice:    false
//        VirtualName: "StringMaxLen255"
//      }
//     ]
//     ClassicLinkVPCId: "StringMaxLen255"
//     ClassicLinkVPCSecurityGroups: [ "StringMaxLen255" ]
//     EbsOptimized:       false
//     IamInstanceProfile: "StringMaxLen1600"
//     ImageId:            "StringMaxLen255"
//     InstanceId:         "StringMaxLen19"
//     InstanceMonitoring: {
//       Enabled: true
//     }
//     InstanceType:     "StringMaxLen255"
//     KernelId:         "StringMaxLen255"
//     KeyName:          "StringMaxLen255"
//     PlacementTenancy: "StringMaxLen64"
//     RamdiskId:        "StringMaxLen255"
//     SecurityGroups: [ "sg-1234" ]
//     SpotPrice: "SpotPrice"
//     UserData:  "StringUserData"
// });
//
// ```
//
// ## AWS.AUTOSCALING.LAUNCHCONFIGS.DELETE
// <a name="ldelete"></a>
// `aws.autoscaling.launchConfigs.delete(region, launchconfigName);`
//
// Delete an autoscaling launchconfig.
//
// Example:
//
// ```
//
//  aws.autoscaling.launchConfigs.delete("us-east-1", "launchconfigName");
//
// ```
//

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "scaling"

func createLaunchConfiguration(region string, params *autoscaling.CreateLaunchConfigurationInput) {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.CreateLaunchConfiguration(params)
	if err != nil {
		log.Fatalf("Error creating scaling launch config: %s", err)
	}
	name := *params.LaunchConfigurationName

	// Wait for it.
	for i := 0; i < 10; i++ {
		g := describeLaunchConfiguration(region, name)
		if g != nil {
			break
		}
		time.Sleep(time.Second * 10)
	}
}

func deleteLaunchConfiguration(region string, id string) {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &autoscaling.DeleteLaunchConfigurationInput{
		LaunchConfigurationName: aws.String(id),
	}
	_, err := svc.DeleteLaunchConfiguration(params)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Wait for it.
	for i := 0; i < 10; i++ {
		g := describeLaunchConfiguration(region, id)
		if g == nil {
			break
		}
		time.Sleep(time.Second * 10)
	}
}

func describeLaunchConfiguration(region string, name string) *autoscaling.LaunchConfiguration {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: []*string{
			aws.String(name),
		},
		MaxRecords: aws.Int64(1),
	}
	resp, err := svc.DescribeLaunchConfigurations(params)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.LaunchConfigurations) > 0 {
		return resp.LaunchConfigurations[0]
	}
	return nil
}

func scanLaunchConfigurations(region string) []*autoscaling.LaunchConfiguration {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: []*string{},
		MaxRecords:               aws.Int64(100),
	}
	resp, err := svc.DescribeLaunchConfigurations(params)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.LaunchConfigurations) > 0 {
		return resp.LaunchConfigurations
	}
	return []*autoscaling.LaunchConfiguration{}
}

func createLifecycleHook(region string, params *autoscaling.PutLifecycleHookInput) {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.PutLifecycleHook(params)
	if err != nil {
		log.Fatalf("Error creating scaling lifecycle hook: %s", err)
	}
	name := *params.LifecycleHookName

	// Wait for it.
	for i := 0; i < 10; i++ {
		g := describeLifecycleHook(region, *params.AutoScalingGroupName, name)
		if g != nil {
			break
		}
		time.Sleep(time.Second * 10)
	}
}

func deleteLifecycleHook(region string, groupName string, hookName string) {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &autoscaling.DeleteLifecycleHookInput{
		AutoScalingGroupName: aws.String(groupName),
		LifecycleHookName:    aws.String(hookName),
	}
	_, err := svc.DeleteLifecycleHook(params)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Wait for it.
	for i := 0; i < 10; i++ {
		g := describeLifecycleHook(region, groupName, hookName)
		if g == nil {
			break
		}
		time.Sleep(time.Second * 10)
	}
}

func describeLifecycleHook(region string, groupName string, hookName string) *autoscaling.LifecycleHook {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &autoscaling.DescribeLifecycleHooksInput{
		AutoScalingGroupName: aws.String(groupName),
		LifecycleHookNames: []*string{
			aws.String(hookName),
		},
	}
	resp, err := svc.DescribeLifecycleHooks(params)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.LifecycleHooks) > 0 {
		return resp.LifecycleHooks[0]
	}
	return nil
}

func scanLifecycleHooks(region string) []*autoscaling.LifecycleHook {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	groups := scanAutoScalingGroups(region)

	hooks := []*autoscaling.LifecycleHook{}
	for _, g := range groups {
		params := &autoscaling.DescribeLifecycleHooksInput{
			AutoScalingGroupName: aws.String(*g.AutoScalingGroupName),
			LifecycleHookNames:   []*string{},
		}
		resp, err := svc.DescribeLifecycleHooks(params)
		if err != nil {
			log.Fatal(err)
		}
		for _, h := range resp.LifecycleHooks {
			hooks = append(hooks, h)
		}
	}

	return hooks
}

func completeLifecycleAction(region string, params *autoscaling.CompleteLifecycleActionInput) {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.CompleteLifecycleAction(params)
	if err != nil {
		log.Fatal(err)
	}
}

func recordLifecycleActionHeartbeat(region string, params *autoscaling.RecordLifecycleActionHeartbeatInput) {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.RecordLifecycleActionHeartbeat(params)
	if err != nil {
		log.Fatal(err)
	}
}

func createAutoScalingGroup(region string, params *autoscaling.CreateAutoScalingGroupInput) {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.CreateAutoScalingGroup(params)
	if err != nil {
		log.Fatalf("Error creating scaling group: %s", err)
	}
	name := *params.AutoScalingGroupName

	// Wait for it.
	for i := 0; i < 10; i++ {
		g := describeAutoScalingGroup(region, name)
		if g != nil {
			break
		}
		time.Sleep(time.Second * 10)
	}
}

func deleteAutoScalingGroup(region string, id string) {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(id),
		ForceDelete:          aws.Bool(true),
	}
	_, err := svc.DeleteAutoScalingGroup(params)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Wait for it.
	for i := 0; i < 10; i++ {
		g := describeAutoScalingGroup(region, id)
		if g == nil {
			break
		}
		time.Sleep(time.Second * 10)
	}
}

func describeAutoScalingGroup(region string, name string) *autoscaling.Group {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(name),
		},
		MaxRecords: aws.Int64(1),
	}
	resp, err := svc.DescribeAutoScalingGroups(params)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.AutoScalingGroups) > 0 {
		return resp.AutoScalingGroups[0]
	}
	return nil
}

func scanAutoScalingGroups(region string) []*autoscaling.Group {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{},
		MaxRecords:            aws.Int64(100),
	}
	resp, err := svc.DescribeAutoScalingGroups(params)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.AutoScalingGroups) > 0 {
		return resp.AutoScalingGroups
	}
	return []*autoscaling.Group{}
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		var o1 *otto.Object
		var o2 *otto.Object
		var awsObj *otto.Object
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			awsObj, _ = rt.Object(`aws = {}`)
		} else {
			awsObj = a.Object()
		}

		if b, err := awsObj.Get("autoscaling"); err != nil || b.IsUndefined() {
			o1, _ = rt.Object(`aws.autoscaling = {}`)
		} else {
			o1 = b.Object()
		}

		// Launch Configs
		if c, err := o1.Get("launchConfigs"); err != nil || c.IsUndefined() {
			o2, _ = rt.Object(`aws.autoscaling.launchConfigs = {}`)
		} else {
			o2 = c.Object()
		}
		o2.Set("scan", func(call otto.FunctionCall) otto.Value {
			f := mcore.Sanitizer(rt)
			region := call.Argument(0).String()
			return f(scanLaunchConfigurations(region))
		})
		o2.Set("delete", deleteLaunchConfiguration)
		o2.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input autoscaling.CreateLaunchConfigurationInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for autoscaling launch config create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall autoscaling launch config json: %s", err)
			}

			region := call.Argument(0).String()

			createLaunchConfiguration(region, &input)
			return otto.Value{}
		})
		o2.Set("describe", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			f := mcore.Sanitizer(rt)
			return f(describeLaunchConfiguration(region, id))
		})

		// Lifcycle Hooks
		if c, err := o1.Get("hooks"); err != nil || c.IsUndefined() {
			o2, _ = rt.Object(`aws.autoscaling.hooks = {}`)
		} else {
			o2 = c.Object()
		}
		o2.Set("scan", func(call otto.FunctionCall) otto.Value {
			f := mcore.Sanitizer(rt)
			region := call.Argument(0).String()
			return f(scanLifecycleHooks(region))
		})
		o2.Set("delete", deleteLifecycleHook)
		o2.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input autoscaling.PutLifecycleHookInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for autoscaling hook create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall autoscaling hook json: %s", err)
			}

			region := call.Argument(0).String()

			createLifecycleHook(region, &input)
			return otto.Value{}
		})
		o2.Set("describe", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			group := call.Argument(1).String()
			hook := call.Argument(1).String()
			f := mcore.Sanitizer(rt)
			return f(describeLifecycleHook(region, group, hook))
		})
		o2.Set("complete", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input autoscaling.CompleteLifecycleActionInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for autoscaling hook complete input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall autoscaling hook complete json: %s", err)
			}

			region := call.Argument(0).String()

			completeLifecycleAction(region, &input)
			return otto.Value{}
		})
		o2.Set("recordHeartbeat", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input autoscaling.RecordLifecycleActionHeartbeatInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for autoscaling hook heartbeat input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall autoscaling hook heatbeat json: %s", err)
			}

			region := call.Argument(0).String()

			recordLifecycleActionHeartbeat(region, &input)
			return otto.Value{}
		})

		// ASGs
		if c, err := o1.Get("groups"); err != nil || c.IsUndefined() {
			o2, _ = rt.Object(`aws.autoscaling.groups = {}`)
		} else {
			o2 = c.Object()
		}
		o2.Set("scan", func(call otto.FunctionCall) otto.Value {
			f := mcore.Sanitizer(rt)
			region := call.Argument(0).String()
			return f(scanAutoScalingGroups(region))
		})
		o2.Set("delete", deleteAutoScalingGroup)
		o2.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input autoscaling.CreateAutoScalingGroupInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for autoscaling group create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall autoscaling group json: %s", err)
			}

			region := call.Argument(0).String()

			createAutoScalingGroup(region, &input)
			return otto.Value{}
		})
		o2.Set("describe", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			f := mcore.Sanitizer(rt)
			return f(describeAutoScalingGroup(region, id))
		})
	})
}
