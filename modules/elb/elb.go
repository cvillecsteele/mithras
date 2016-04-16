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
// # CORE FUNCTIONS: ELB
//

package elb

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.elbs.create](#create)
// > * [aws.elbs.delete](#delete)
// > * [aws.elbs.describe](#describe)
// > * [aws.elbs.scan](#scan)
// > * [aws.elbs.register](#register)
// > * [aws.elbs.deRegister](#deregister)
// > * [aws.elbs.setHealth](#health)
// > * [aws.elbs.setAttrs](#attrs)
//
// This API allows exposes functions to manage AWS elastic load balancers.
//
// ## AWS.ELBS.CREATE
// <a name="create"></a>
// `aws.elbs.create(region, config);`
//
// Create an ELB
//
// Example:
//
// ```
//
//  var lb = aws.elbs.create("us-east-1", {
// Listeners: [
//  {
//    InstancePort:     80
//    LoadBalancerPort: 80
//    Protocol:         "http"
//    InstanceProtocol: "http"
//  },
// ]
// LoadBalancerName: "test-lb"
// SecurityGroups: [
//     "sg-xyz"
// ]
// Subnets: [
//     "subnet-123"
//     "subnet-456"
// ]
// Tags: [
//     {
//       Key:   "foo"
//       Value: "bar"
//     },
// ]
// });
//
// ```
//
// ## AWS.ELBS.SETATTRS
// <a name="attrs"></a>
// `aws.elbs.setAttrs(region, lbName, config);`
//
// Set ELB attributes
//
// Example:
//
// ```
//
//  aws.elbs.setAttrs("us-east-1", "test-lb"{
//  LoadBalancerAttributes: {
//      AccessLog: {
//        Enabled:        false
//        EmitInterval:   60
//        S3BucketName:   "my-loadbalancer-logs"
//        S3BucketPrefix: "test-app"
//      }
//      ConnectionDraining: {
//        Enabled: true
//        Timeout: 300
//      }
//      ConnectionSettings: {
//        IdleTimeout: 30
//      }
//      CrossZoneLoadBalancing: {
//        Enabled: true
//      }
//    }
//    LoadBalancerName: "test-lb"
//  }
// });
//
// ```
//
// ## AWS.ELBS.SETHEALTH
// <a name="health"></a>
// `aws.elbs.setHealth(region, lbName, config);`
//
// Set ELB health check.
//
// Example:
//
// ```
//
//  aws.elbs.setHealth("us-east-1", "test-lb", {
//  HealthCheck: {
//      HealthyThreshold:   2
//      Interval:           30
//      Target:             "HTTP:80/hc"
//      Timeout:            5
//      UnhealthyThreshold: 3
//  }
//  LoadBalancerName: "test-lb"
// });
//
// ```
//
// ## AWS.ELBS.DELETE
// <a name="delete"></a>
// `aws.elbs.delete(region, lbName);`
//
// Delete an ELB.
//
// Example:
//
// ```
//
//  aws.elbs.delete("us-east-1", "test-lb");
//
// ```
//
// ## AWS.ELBS.DESCRIBE
// <a name="describe"></a>
// `aws.elbs.describe(region, lbName);`
//
// Get info about an ELB.
//
// Example:
//
// ```
//
//  var elb = aws.elbs.describe("us-east-1", "test-lb");
//
// ```
//
// ## AWS.ELBS.SCAN
// <a name="scan"></a>
// `aws.elbs.scan(region);`
//
// Get info about ELBs.
//
// Example:
//
// ```
//
//  var elbs = aws.elbs.scan("us-east-1");
//
// ```
//
// ## AWS.ELBS.REGISTER
// <a name="register"></a>
// `aws.elbs.register(region, lbName, instance);`
//
// Register an EC2 instance with an ELB.
//
// Example:
//
// ```
//
// aws.elbs.register("us-east-1", "test-lb",
// 		  {
// 		      "AmiLaunchIndex": 0,
// 		      "Architecture": "x86_64",
// 		      "BlockDeviceMappings": [
// 			  {
// 			      "DeviceName": "/dev/xvda",
// 			      "Ebs": {
// 				  "AttachTime": "2016-03-31T19:17:44Z",
// 				  "DeleteOnTermination": true,
// 				  "Status": "attached",
// 				  "VolumeId": "vol-3d1ab09f"
// 			      }
// 			  }
// 		      ],
// 		      "ClientToken": "",
// 		      "EbsOptimized": false,
// 		      "Hypervisor": "xen",
// 		      "IamInstanceProfile": {
// 			  "Arn": "arn:aws:iam::286536233385:instance-profile/cr-webserver",
// 			  "Id": "AIPAIOOUMBIRCV2QCTIYA"
// 		      },
// 		      "ImageId": "ami-60b6c60a",
// 		      "InstanceId": "i-824ad119",
// 		      "InstanceLifecycle": null,
// 		      "InstanceType": "t2.small",
// 		      "KernelId": null,
// 		      "KeyName": "cr",
// 		      "LaunchTime": "2016-03-31T19:17:43Z",
// 		      "Monitoring": {
// 			  "State": "enabled"
// 		      },
// 		      "NetworkInterfaces": [
// 			  {
// 			      "Association": {
// 				  "IpOwnerId": "amazon",
// 				  "PublicDnsName": "ec2-52-90-244-101.compute-1.amazonaws.com",
// 				  "PublicIp": "52.90.244.101"
// 			      },
// 			      "Attachment": {
// 				  "AttachTime": "2016-03-31T19:17:43Z",
// 				  "AttachmentId": "eni-attach-4d1c8ab1",
// 				  "DeleteOnTermination": true,
// 				  "DeviceIndex": 0,
// 				  "Status": "attached"
// 			      },
// 			      "Description": "",
// 			      "Groups": [
// 				  {
// 				      "GroupId": "sg-57c9812f",
// 				      "GroupName": "webserver"
// 				  }
// 			      ],
// 			      "MacAddress": "0e:6b:fe:c1:cb:45",
// 			      "NetworkInterfaceId": "eni-a6b644f6",
// 			      "OwnerId": "286536233385",
// 			      "PrivateDnsName": "ip-172-33-1-178.ec2.internal",
// 			      "PrivateIpAddress": "172.33.1.178",
// 			      "PrivateIpAddresses": [
// 				  {
// 				      "Association": {
// 					  "IpOwnerId": "amazon",
// 					  "PublicDnsName": "ec2-52-90-244-101.compute-1.amazonaws.com",
// 					  "PublicIp": "52.90.244.101"
// 				      },
// 				      "Primary": true,
// 				      "PrivateDnsName": "ip-172-33-1-178.ec2.internal",
// 				      "PrivateIpAddress": "172.33.1.178"
// 				  }
// 			      ],
// 			      "SourceDestCheck": true,
// 			      "Status": "in-use",
// 			      "SubnetId": "subnet-d76ae58f",
// 			      "VpcId": "vpc-b88883dc"
// 			  }
// 		      ],
// 		      "Placement": {
// 			  "Affinity": null,
// 			  "AvailabilityZone": "us-east-1d",
// 			  "GroupName": "",
// 			  "HostId": null,
// 			  "Tenancy": "default"
// 		      },
// 		      "Platform": null,
// 		      "PrivateDnsName": "ip-172-33-1-178.ec2.internal",
// 		      "PrivateIpAddress": "172.33.1.178",
// 		      "ProductCodes": null,
// 		      "PublicDnsName": "ec2-52-90-244-101.compute-1.amazonaws.com",
// 		      "PublicIpAddress": "52.90.244.101",
// 		      "RamdiskId": null,
// 		      "RootDeviceName": "/dev/xvda",
// 		      "RootDeviceType": "ebs",
// 		      "SecurityGroups": [
// 			  {
// 			      "GroupId": "sg-57c9812f",
// 			      "GroupName": "webserver"
// 			  }
// 		      ],
// 		      "SourceDestCheck": true,
// 		      "SpotInstanceRequestId": null,
// 		      "SriovNetSupport": null,
// 		      "State": {
// 			  "Code": 16,
// 			  "Name": "running"
// 		      },
// 		      "StateReason": null,
// 		      "StateTransitionReason": "",
// 		      "SubnetId": "subnet-d76ae58f",
// 		      "Tags": [
// 			  {
// 			      "Key": "Name",
// 			      "Value": "webserver"
// 			  }
// 		      ],
// 		      "VirtualizationType": "hvm",
// 		      "VpcId": "vpc-b88883dc",
// 		      "uname": "Linux ip-172-33-1-178 4.1.10-17.31.amzn1.x86_64 #1 SMP Sat Oct 24 01:31:37 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux\r\n"
// 		  });
//
// ```
//
// ## AWS.ELBS.DEREGISTER
// <a name="deregister"></a>
// `aws.elbs.deregister(region, lbName, instance);`
//
// Deregister an EC2 instance with an ELB.
//
// Example:
//
// ```
//
// aws.elbs.deregister("us-east-1", "test-lb",
// 		  {
// 		      "AmiLaunchIndex": 0,
// 		      "Architecture": "x86_64",
// 		      "BlockDeviceMappings": [
// 			  {
// 			      "DeviceName": "/dev/xvda",
// 			      "Ebs": {
// 				  "AttachTime": "2016-03-31T19:17:44Z",
// 				  "DeleteOnTermination": true,
// 				  "Status": "attached",
// 				  "VolumeId": "vol-3d1ab09f"
// 			      }
// 			  }
// 		      ],
// 		      "ClientToken": "",
// 		      "EbsOptimized": false,
// 		      "Hypervisor": "xen",
// 		      "IamInstanceProfile": {
// 			  "Arn": "arn:aws:iam::286536233385:instance-profile/cr-webserver",
// 			  "Id": "AIPAIOOUMBIRCV2QCTIYA"
// 		      },
// 		      "ImageId": "ami-60b6c60a",
// 		      "InstanceId": "i-824ad119",
// 		      "InstanceLifecycle": null,
// 		      "InstanceType": "t2.small",
// 		      "KernelId": null,
// 		      "KeyName": "cr",
// 		      "LaunchTime": "2016-03-31T19:17:43Z",
// 		      "Monitoring": {
// 			  "State": "enabled"
// 		      },
// 		      "NetworkInterfaces": [
// 			  {
// 			      "Association": {
// 				  "IpOwnerId": "amazon",
// 				  "PublicDnsName": "ec2-52-90-244-101.compute-1.amazonaws.com",
// 				  "PublicIp": "52.90.244.101"
// 			      },
// 			      "Attachment": {
// 				  "AttachTime": "2016-03-31T19:17:43Z",
// 				  "AttachmentId": "eni-attach-4d1c8ab1",
// 				  "DeleteOnTermination": true,
// 				  "DeviceIndex": 0,
// 				  "Status": "attached"
// 			      },
// 			      "Description": "",
// 			      "Groups": [
// 				  {
// 				      "GroupId": "sg-57c9812f",
// 				      "GroupName": "webserver"
// 				  }
// 			      ],
// 			      "MacAddress": "0e:6b:fe:c1:cb:45",
// 			      "NetworkInterfaceId": "eni-a6b644f6",
// 			      "OwnerId": "286536233385",
// 			      "PrivateDnsName": "ip-172-33-1-178.ec2.internal",
// 			      "PrivateIpAddress": "172.33.1.178",
// 			      "PrivateIpAddresses": [
// 				  {
// 				      "Association": {
// 					  "IpOwnerId": "amazon",
// 					  "PublicDnsName": "ec2-52-90-244-101.compute-1.amazonaws.com",
// 					  "PublicIp": "52.90.244.101"
// 				      },
// 				      "Primary": true,
// 				      "PrivateDnsName": "ip-172-33-1-178.ec2.internal",
// 				      "PrivateIpAddress": "172.33.1.178"
// 				  }
// 			      ],
// 			      "SourceDestCheck": true,
// 			      "Status": "in-use",
// 			      "SubnetId": "subnet-d76ae58f",
// 			      "VpcId": "vpc-b88883dc"
// 			  }
// 		      ],
// 		      "Placement": {
// 			  "Affinity": null,
// 			  "AvailabilityZone": "us-east-1d",
// 			  "GroupName": "",
// 			  "HostId": null,
// 			  "Tenancy": "default"
// 		      },
// 		      "Platform": null,
// 		      "PrivateDnsName": "ip-172-33-1-178.ec2.internal",
// 		      "PrivateIpAddress": "172.33.1.178",
// 		      "ProductCodes": null,
// 		      "PublicDnsName": "ec2-52-90-244-101.compute-1.amazonaws.com",
// 		      "PublicIpAddress": "52.90.244.101",
// 		      "RamdiskId": null,
// 		      "RootDeviceName": "/dev/xvda",
// 		      "RootDeviceType": "ebs",
// 		      "SecurityGroups": [
// 			  {
// 			      "GroupId": "sg-57c9812f",
// 			      "GroupName": "webserver"
// 			  }
// 		      ],
// 		      "SourceDestCheck": true,
// 		      "SpotInstanceRequestId": null,
// 		      "SriovNetSupport": null,
// 		      "State": {
// 			  "Code": 16,
// 			  "Name": "running"
// 		      },
// 		      "StateReason": null,
// 		      "StateTransitionReason": "",
// 		      "SubnetId": "subnet-d76ae58f",
// 		      "Tags": [
// 			  {
// 			      "Key": "Name",
// 			      "Value": "webserver"
// 			  }
// 		      ],
// 		      "VirtualizationType": "hvm",
// 		      "VpcId": "vpc-b88883dc",
// 		      "uname": "Linux ip-172-33-1-178 4.1.10-17.31.amzn1.x86_64 #1 SMP Sat Oct 24 01:31:37 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux\r\n"
// 		  });
//
// ```
//
import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "elb"

func setHealth(region string, lbName string, input elb.ConfigureHealthCheckInput, verbose bool) {
	svc := elb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.ConfigureHealthCheck(&input)
	if err != nil {
		log.Fatalf("Can't configure elb health check: %s", err)
	}
}

func setAttrs(region string, lbName string, input elb.ModifyLoadBalancerAttributesInput, verbose bool) {
	svc := elb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.ModifyLoadBalancerAttributes(&input)
	if err != nil {
		log.Fatalf("Can't set elb attributes: %s", err)
	}
}

func describe(region string, id string) *elb.LoadBalancerDescription {
	svc := elb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{
			aws.String(id),
		},
	}
	resp, err := svc.DescribeLoadBalancers(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if "LoadBalancerNotFound" == awsErr.Code() {
				return nil
			}
		}
		log.Fatalf("Error describing elb: %s", err)
	}

	return resp.LoadBalancerDescriptions[0]
}

func create(region string, params *elb.CreateLoadBalancerInput, verbose bool) *elb.LoadBalancerDescription {
	svc := elb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.CreateLoadBalancer(params)
	if err != nil {
		log.Fatal(err.Error())
	}
	id := *params.LoadBalancerName

	// Wait for it.
	var target *elb.LoadBalancerDescription
	for i := 0; i < 10; i++ {
		target = describe(region, id)
		if target != nil && *target.LoadBalancerName == id {
			break
		}
		time.Sleep(time.Second * 10)
	}

	return target
}

func delete(region string, id string, verbose bool) {
	svc := elb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	i := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(id),
	}
	_, err := svc.DeleteLoadBalancer(i)

	if err != nil {
		log.Fatal(err.Error())
	}

	// Wait for it.
	for i := 0; i < 10; i++ {
		target := describe(region, id)
		if target == nil {
			return
		}
		time.Sleep(time.Second * 10)
	}

	log.Fatal("Timeout waiting for elb deletion.")
}

func register(region string, lbName string, instances []*elb.Instance, verbose bool) {
	svc := elb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &elb.RegisterInstancesWithLoadBalancerInput{
		Instances:        instances,
		LoadBalancerName: aws.String(lbName),
	}

	if _, err := svc.RegisterInstancesWithLoadBalancer(params); err != nil {
		log.Fatalf("Error adding instances to elb '%s': %s", lbName, err)
	}
}

func deRegister(region string, lbName string, instances []*elb.Instance, verbose bool) {
	svc := elb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &elb.DeregisterInstancesFromLoadBalancerInput{
		Instances:        instances,
		LoadBalancerName: aws.String(lbName),
	}

	if _, err := svc.DeregisterInstancesFromLoadBalancer(params); err != nil {
		log.Fatalf("Error removing instances from elb '%s': %s", lbName, err)
	}
}

func scan(rt *otto.Otto, region string) otto.Value {
	svc := elb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.DescribeLoadBalancers(nil)
	if err != nil {
		panic(err)
	}

	// shove instances into jsland
	lbs := []elb.LoadBalancerDescription{}
	for _, i := range resp.LoadBalancerDescriptions {
		lbs = append(lbs, *i)
	}
	return mcore.Sanitize(rt, lbs)
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			rt.Object(`aws = {}`)
		}
		o1, _ := rt.Object(`aws.elbs = {}`)
		o1.Set("scan", func(region string) otto.Value {
			return scan(rt, region)
		})
		o1.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input elb.CreateLoadBalancerInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for ELB create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall ELB create json: %s", err)
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
		o1.Set("register", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input []*elb.Instance
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(2))
			if err != nil {
				log.Fatalf("Can't create json for ELB register input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall ELB register json: %s", err)
			}

			region := call.Argument(0).String()
			lbName := call.Argument(1).String()
			verbose := mcore.IsVerbose(rt)
			register(region, lbName, input, verbose)
			return otto.Value{}
		})
		o1.Set("deRegister", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input []*elb.Instance
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(2))
			if err != nil {
				log.Fatalf("Can't create json for ELB deregister input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall ELB deregister json: %s", err)
			}

			region := call.Argument(0).String()
			lbName := call.Argument(1).String()
			verbose := mcore.IsVerbose(rt)
			deRegister(region, lbName, input, verbose)
			return otto.Value{}
		})
		o1.Set("setHealth", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input elb.ConfigureHealthCheckInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(2))
			if err != nil {
				log.Fatalf("Can't create json for ELB setHealth input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall ELB setHealth json: %s", err)
			}
			region := call.Argument(0).String()
			lbName := call.Argument(1).String()
			verbose := mcore.IsVerbose(rt)
			setHealth(region, lbName, input, verbose)
			return otto.Value{}
		})
		o1.Set("setAttrs", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input elb.ModifyLoadBalancerAttributesInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(2))
			if err != nil {
				log.Fatalf("Can't create json for ELB setAttrs input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall ELB setAttrs json: %s", err)
			}
			region := call.Argument(0).String()
			lbName := call.Argument(1).String()
			verbose := mcore.IsVerbose(rt)
			setAttrs(region, lbName, input, verbose)
			return otto.Value{}
		})
	})
}
