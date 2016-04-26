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
// # CORE FUNCTIONS: SCALING
//

package scaling

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.scaling.scan](#scan)
// > * [aws.scaling.create](#create)
// > * [aws.scaling.delete](#delete)
// > * [aws.scaling.describe](#describe)
// > * [aws.scaling.messages.send](#msend)
// > * [aws.scaling.messages.receive](#mreceive)
// > * [aws.scaling.messages.delete](#mdelete)
//
// This API allows resource handlers to manage Autoscaling groups.
//
// ## AWS.SCALING.SCAN
// <a name="scan"></a>
// `aws.scaling.scan(region);`
//
// Returns a list of scaling queues.
//
// Example:
//
// ```
//
//  var queues = aws.scaling.scan("us-east-1");
//
// ```
//
// ## AWS.SCALING.CREATE
// <a name="create"></a>
// `aws.scaling.create(region, config);`
//
// Create a SCALING queue.
//
// Example:
//
// ```
//
//  var scaling =  aws.scaling.create(
//    "us-east-1",
//    {
//      QueueName: "myqueue"
//      Attributes: [
//        "Key": "value"
//      ]
//    });
//
// ```
//
// ## AWS.SCALING.DELETE
// <a name="delete"></a>
// `aws.scaling.delete(region, url);`
//
// Delete an SCALING queue.
//
// Example:
//
// ```
//
//  aws.scaling.delete("us-east-1", "queueUrl");
//
// ```
//
// ## AWS.SCALING.MESSAGES.SEND
// <a name="msend"></a>
// `aws.scaling.messages.send(region, input);`
//
// Publish to an SCALING queue.
//
// Example:
//
// ```
//
//  aws.scaling.messages.send("us-east-1",
//                  {
//                    MessageBody:  "body"
//                    QueueUrl:     "url"
//                    DelaySeconds: 1
//                    MessageAttributes: {
//                      "Key": {
//                        DataType: "type"
//                        BinaryListValues: [
//                          "PAYLOAD"
//                        ]
//                        BinaryValue: "PAYLOAD"
//                        StringListValues: [
//                          "String"
//                        ]
//                        StringValue: "String"
//                      }
//                    }
//                   });
//
// ```
//
// ## AWS.SCALING.DESCRIBE
// <a name="describe"></a>
// `aws.scaling.describe(region, scaling_id);`
//
// Get info from AWS about an SCALING queue.
//
// Example:
//
// ```
//
//  var queue = aws.scaling.describe("us-east-1", "queueName");
//
// ```
//
// ## AWS.SCALING.MESSAGES.RECEIVE
// <a name="mreceive"></a>
// `aws.scaling.messages.receive(region, input);`
//
// Get a message from a queue
//
// Example:
//
// ```
//  var message =
//  aws.scaling.messages.receive("us-east-1",
//                  {
//                    QueueUrl: "url"
//                    AttributeNames: [
//                      "QueueAttributeName"
//                    ]
//                    MaxNumberOfMessages: 1
//                    MessageAttributeNames: [
//                      "MessageAttributeName"
//                    ]
//                    VisibilityTimeout: 1
//                    WaitTimeSeconds:   1
//                  });
//
// ```
//
// ## AWS.SCALING.MESSAGES.DELETE
// <a name="mdelete"></a>
// `aws.scaling.messages.delete(region, input);`
//
// Get a message from a queue
//
// Example:
//
// ```
//  aws.scaling.messages.delete("us-east-1",
//  {
//    QueueUrl:      "queueUrl"
//    ReceiptHandle: "123456"
//  });
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

	return resp.LaunchConfigurations
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

func scanLifecycleHooks(region string, groupName string) []*autoscaling.LifecycleHook {
	svc := autoscaling.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &autoscaling.DescribeLifecycleHooksInput{
		AutoScalingGroupName: aws.String(groupName),
		LifecycleHookNames:   []*string{},
	}
	resp, err := svc.DescribeLifecycleHooks(params)
	if err != nil {
		log.Fatal(err)
	}

	return resp.LifecycleHooks
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

	return resp.AutoScalingGroups
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
			group := call.Argument(1).String()
			return f(scanLifecycleHooks(region, group))
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
