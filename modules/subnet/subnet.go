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
// # CORE FUNCTIONS: SUBNET
//

package subnet

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.subnets.scan](#vscan)
// > * [aws.subnets.create](#vcreate)
// > * [aws.subnets.delete](#vdelete)
// > * [aws.subnets.describe](#vdescribe)
// > * [aws.subnets.routes.create](#gcreate)
// > * [aws.subnets.routes.delete](#gdelete)
//
// This API allows resource handlers to manage subnets in an AWS VPC.
//
// ## AWS.SUBNETS.SCAN
// <a name="sscan"></a>
// `aws.subnets.scan(region);`
//
// Returns a list of subnets.
//
// Example:
//
// ```
//
//  var subnets =  aws.subnets.scan("us-east-1");
//
// ```
//
// ## AWS.SUBNETS.CREATE
// <a name="screate"></a>
// `aws.subnets.create(region, config);`
//
// Create a SUBNET.
//
// Example:
//
// ```
//
//  var subnet =  aws.subnets.create(
//    "us-east-1",
//    {
//      CidrBlock:        "172.33.1.0/24"
//      VpcId:            "vpc-abcd",
//      AvailabilityZone: "some zone"
//    });
//
// ```
//
// ## AWS.SUBNETS.DELETE
// <a name="sdelete"></a>
// `aws.subnets.delete(region, subnet_id);`
//
// Delete a SUBNET.
//
// Example:
//
// ```
//
//  aws.subnets.delete("us-east-1", "subnet-abcd");
//
// ```
//
// ## AWS.SUBNETS.DESCRIBE
// <a name="sdescribe"></a>
// `aws.subnets.describe(region, subnet_id);`
//
// Get info from AWS about a SUBNET.
//
// Example:
//
// ```
//
//  var subnet = aws.subnets.describe("us-east-1", "subnet-abcd");
//
// ```
//
// ## AWS.SUBNETS.ROUTES.CREATE
// <a name="rcreate"></a>
// `aws.subnets.routes.create(region);`
//
// Create a route.
//
// Example:
//
// ```
//
//  var route =  aws.subnets.routes.create("us-east-1", {
//		   DestinationCidrBlock: "0.0.0.0/0",
//		   GatewayId:            "gw-1234"
//		 });
//
// ```
//
// ## AWS.SUBNETS.ROUTES.DELETE
// <a name="rdelete"></a>
// `aws.subnets.routes.delete(region, route_id);`
//
// Delete a route.
//
// Example:
//
// ```
//
//  aws.subnets.routes.delete("us-east-1", "route-abcd");
//
// ```
//

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "subnet"

func describeSubnet(region string, id string) *ec2.Subnet {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &ec2.DescribeSubnetsInput{
		SubnetIds: []*string{
			aws.String(id),
		},
	}
	resp, err := svc.DescribeSubnets(params)

	if err != nil {
		return nil
	}

	return resp.Subnets[0]
}

func createSubnet(region string, params *ec2.CreateSubnetInput, verbose bool) *ec2.Subnet {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.CreateSubnet(params)
	if err != nil {
		log.Fatal(err.Error())
	}
	id := *resp.Subnet.SubnetId

	// Wait for it.
	var subnet *ec2.Subnet
	for i := 0; i < 10; i++ {
		subnet = describeSubnet(region, id)
		if subnet != nil && *subnet.State == "available" {
			break
		}
		time.Sleep(time.Second * 10)
	}

	return subnet
}

func createRoute(region string, input ec2.CreateRouteInput, verbose bool) *ec2.CreateRouteOutput {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	output, err := svc.CreateRoute(&input)
	if err != nil {
		log.Fatalf("AWS error creating route: %s", err)
	}
	return output
}

func deleteRoute(region string, cidr string, routeTableId string, verbose bool) *ec2.DeleteRouteOutput {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &ec2.DeleteRouteInput{
		DestinationCidrBlock: aws.String(cidr),
		RouteTableId:         aws.String(routeTableId),
	}
	resp, err := svc.DeleteRoute(params)

	if err != nil {
		log.Fatalf("AWS error deleting route: %s", err)
	}
	return resp
}

func deleteSubnet(region string, id string, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	i := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(id),
	}
	_, err := svc.DeleteSubnet(i)

	if err != nil {
		log.Fatal(err.Error())
	}

	// TODO: Wait for it.
}

func scanSubnets(rt *otto.Otto, region string) otto.Value {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.DescribeSubnets(nil)
	if err != nil {
		panic(err)
	}

	// shove instances into jsland
	subnets := []ec2.Subnet{}
	for _, i := range resp.Subnets {
		subnets = append(subnets, *i)
	}
	return mcore.Sanitize(rt, subnets)
}

func init() {
	// mcore.RegisterHandler(handle)
	mcore.RegisterInit(func(rt *otto.Otto) {
		var o0 *otto.Object
		var o1 *otto.Object
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			rt.Object(`aws = {}`)
		} else {
			if b, err := a.Object().Get("aws.subnets"); err != nil || b.IsUndefined() {
				o0, _ = rt.Object(`aws.subnets = {}`)
				o1, _ = rt.Object(`aws.subnets.routes = {}`)
			} else {
				log.Fatalf("Logic error: aws.subnets already defined")
			}
		}

		o0.Set("scan", func(region string) otto.Value {
			return scanSubnets(rt, region)
		})
		o0.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input ec2.CreateSubnetInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for subnet create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall subnet create json: %s", err)
			}

			region := call.Argument(0).String()
			verbose := mcore.IsVerbose(rt)
			f := mcore.Sanitizer(rt)
			return f(createSubnet(region, &input, verbose))
		})
		o0.Set("delete", func(call otto.FunctionCall) otto.Value {
			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			deleteSubnet(region, id, verbose)
			return otto.Value{}
		})
		o0.Set("describe", func(call otto.FunctionCall) otto.Value {
			f := mcore.Sanitizer(rt)
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			return f(describeSubnet(region, id))
		})

		o1.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input ec2.CreateRouteInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for subnet create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall subnet create json: %s", err)
			}

			region := call.Argument(0).String()
			verbose := mcore.IsVerbose(rt)
			f := mcore.Sanitizer(rt)
			return f(createRoute(region, input, verbose))
		})
		o1.Set("delete", func(call otto.FunctionCall) otto.Value {
			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()
			cidr := call.Argument(1).String()
			routeTableId := call.Argument(2).String()
			deleteRoute(region, cidr, routeTableId, verbose)
			return otto.Value{}
		})

	})
}
