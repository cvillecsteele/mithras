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
// # CORE FUNCTIONS: INSTANCE
//

package instance

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.instances.scan](#scan)
// > * [aws.instances.create](#create)
// > * [aws.instances.delete](#delete)
// > * [aws.instances.describe](#describe)
//
// This API allows resource handlers to manage EC2 instances.
//
// ## AWS.INSTANCES.SCAN
// <a name="scan"></a>
// `aws.instances.scan(region);`
//
// Returns a list of instances.
//
// Example:
//
// ```
//
//  var instances =  aws.instances.scan("us-east-1");
//
// ```
//
// ## AWS.INSTANCES.CREATE
// <a name="create"></a>
// `aws.instances.create(region, config);`
//
// Create one or more instances.
//
// Example:
//
// ```
//
//  var inst =  aws.instances.create(
//    "us-east-1",
//    {
//      ImageId:        "ami-60b6c60a"
//      MaxCount:       1
//      MinCount:       1
//      DisableApiTermination: false
//      EbsOptimized:          false
//      IamInstanceProfile: {
//        Name: iamProfileName
//      }
//      InstanceInitiatedShutdownBehavior: "terminate"
//      InstanceType:                      "t2.small"
//      KeyName:                           "my-key"
//      Monitoring: {
//        Enabled: true
//      }
//      NetworkInterfaces: [
//        {
//          AssociatePublicIpAddress: true
//          DeleteOnTermination:      true
//          DeviceIndex:              0
//          Groups:                  [ "sg-1234" ]
//          SubnetId:                "subnet-abcd"
//        }
//      ]
//    });
//
// ```
//
// ## AWS.INSTANCES.DELETE
// <a name="delete"></a>
// `aws.instances.delete(region, instance_id);`
//
// Delete an instance
//
// Example:
//
// ```
//
//  aws.instances.delete("us-east-1", "i-abcd");
//
// ```
//
// ## AWS.INSTANCES.DESCRIBE
// <a name="describe"></a>
// `aws.instances.describe(region, inst_id);`
//
// Get info from AWS about an instance.
//
// Example:
//
// ```
//
//  var i = aws.instances.describe("us-east-1", "i-abcd");
//
// ```
//

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "instance"

func describe(region string, id string) *ec2.Instance {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	}
	resp, err := svc.DescribeInstances(params)

	if err != nil {
		return nil
	}
	if len(resp.Reservations) > 0 && len(resp.Reservations[0].Instances) > 0 {
		return resp.Reservations[0].Instances[0]
	}
	return nil
}

func create(region string, params *ec2.RunInstancesInput, verbose bool) []*ec2.Instance {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.RunInstances(params)
	if err != nil {
		log.Fatalf("Error running instances: %s", err)
	}
	if len(resp.Instances) < 1 {
		log.Fatal("Zero reservations returned")
	}

	instances := []*ec2.Instance{}
	for _, i := range resp.Instances {
		id := *i.InstanceId

		// Wait for it.
		for i := 0; i < 10; i++ {
			target := describe(region, id)
			if target != nil && *target.InstanceId == id && *target.State.Name == "running" {
				break
			}
			time.Sleep(time.Second * 10)
		}

		// Describe it.
		instances = append(instances, describe(region, id))
	}

	return instances
}

func delete(region string, id string, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	}
	_, err := svc.TerminateInstances(params)

	if err != nil {
		log.Fatalf("Error terminating instances: %s", err)
	}

	// Wait for it.
	for i := 0; i < 10; i++ {
		target := describe(region, id)
		if target != nil && *target.InstanceId == id && *target.State.Name == "terminated" {
			break
		}
		time.Sleep(time.Second * 10)
	}

}

func scan(rt *otto.Otto, region string) otto.Value {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	instances := []ec2.Instance{}
	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			instances = append(instances, *i)
		}
	}
	return mcore.Sanitize(rt, instances)
}

func init() {
	mcore.RegisterInit(func(context *mcore.Context) {
		rt := context.Runtime

		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			rt.Object(`aws = {}`)
		}
		o1, _ := rt.Object(`aws.instances = {}`)
		o1.Set("scan", func(region string) otto.Value {
			return scan(rt, region)
		})
		o1.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input ec2.RunInstancesInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for instance create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall instance create json: %s", err)
			}

			region := call.Argument(0).String()
			verbose := mcore.IsVerbose(rt)

			f := mcore.Sanitizer(rt)
			return f(create(region, &input, verbose))
		})
		o1.Set("delete", func(call otto.FunctionCall) otto.Value {
			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			delete(region, id, verbose)
			return otto.Value{}
		})
		o1.Set("describe", func(call otto.FunctionCall) otto.Value {
			f := mcore.Sanitizer(rt)
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			return f(describe(region, id))
		})
	})
}
