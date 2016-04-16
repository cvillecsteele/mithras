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
// # CORE FUNCTIONS: ROUTETABLES
//

package routetables

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.routetables.scan](#scan)
// > * [aws.routetables.describeForSubnet](#describeForSubnet)
// > * [aws.routetables.describe](#describe)
// > * [aws.routetables.create](#create)
// > * [aws.routetables.delete](#delete)
// > * [aws.routetables.associate](#associate)
// > * [aws.routetables.disassociate](#disassociate)
// > * [aws.routetables.deleteAssociation](#deleteAssociation)
//
// This API allows resource handlers to manipulate DNS records in Route53.
//
// ## AWS.ROUTETABLES.SCAN
// <a name="scan"></a>
// `aws.routetables.scan(region);`
//
// Query routetables.
//
// Example:
//
// ```
//
// var tables = aws.routetables.scan("us-east-1");
//
// ```
//
// ## AWS.ROUTETABLES.DESCRIBEFORSUBNET
// <a name="describeforsubnet"></a>
// `aws.routetables.describeForSubnet(region, subnet-id);`
//
// Get routetables associated with the supplied subnet.
//
// Example:
//
// ```
//
// var tables = aws.routetables.describeForSubnet("us-east-1", "subnet-abc");
//
// ```
//
// ## AWS.ROUTETABLES.DESCRIBEFORSUBNET
// <a name="describeforsubnet"></a>
// `aws.routetables.describeForSubnet(region, subnet-id);`
//
// Get routetables associated with the supplied subnet.
//
// Example:
//
// ```
//
// var tables = aws.routetables.describeForSubnet("us-east-1", "subnet-abc");
//
// ```
//
// ## AWS.ROUTETABLES.DESCRIBE
// <a name="describe"></a>
// `aws.routetables.describe(region, route-table-id);`
//
// Get info about the supplied route table.
//
// Example:
//
// ```
//
// var table = aws.routetables.describe("us-east-1", "routetable-abc");
//
// ```
//
// ## AWS.ROUTETABLES.CREATE
// <a name="create"></a>
// `aws.routetables.create(region, vpc-id);`
//
// Create a route table in a vpc.
//
// Example:
//
// ```
//
// var table = aws.routetables.create("us-east-1", "vpc-123");
//
// ```
//
// ## AWS.ROUTETABLES.DELETE
// <a name="delete"></a>
// `aws.routetables.delete(region, subnet-id);`
//
// Delete a route table.
//
// Example:
//
// ```
//
// aws.routetables.delete("us-east-1", "routetable-123");
//
// ```
//
// ## AWS.ROUTETABLES.ASSOCIATE
// <a name="associate"></a>
// `aws.routetables.associate(region, subnet-id, route-table-id);`
//
// Assocate a route table with a subnet.
//
// Example:
//
// ```
//
// aws.routetables.associate("us-east-1", "subnet-abc", "routetable-123");
//
// ```
//
// ## AWS.ROUTETABLES.DISASSOCIATE
// <a name="disassociate"></a>
// `aws.routetables.disassociate(region, association-id);`
//
// Disassocate a route table and a subnet.
//
// Example:
//
// ```
//
// aws.routetables.disassociate("us-east-1", "association-xyz");
//
// ```
//
// ## AWS.ROUTETABLES.DELETEASSOCIATION
// <a name="deleteAssociation"></a>
// `aws.routetables.deleteAssociation(region, association-id);`
//
// Delete a route table association.
//
// Example:
//
// ```
//
// aws.routetables.deleteAssociation("us-east-1", "association-xyz");
//
// ```
//
import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "routetables"

func associate(region string, routeTableId string, subnetId string, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	// Associate the subnet with the route table
	_, err := svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(routeTableId),
		SubnetId:     aws.String(subnetId),
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func disassociate(region string, associationId string, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.DisassociateRouteTable(&ec2.DisassociateRouteTableInput{
		AssociationId: aws.String(associationId),
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}

func deleteRouteTable(region string, id string, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.DeleteRouteTable(&ec2.DeleteRouteTableInput{
		RouteTableId: aws.String(id),
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	// Wait for it.
	for i := 0; i < 10; i++ {
		x := describe(region, id)
		if x == nil {
			break
		}
		time.Sleep(time.Second * 10)
	}
}

func describe(region string, id string) []*ec2.RouteTable {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{
			aws.String(id),
		},
	}

	// Get route table info
	rtResp, err := svc.DescribeRouteTables(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if "InvalidRouteTableID.NotFound" == awsErr.Code() {
				return []*ec2.RouteTable{}
			}
		}
		log.Fatalf("Error listing bucket: %s", err)
	}

	return rtResp.RouteTables
}

func DescribeRouteTables(region string) []*ec2.RouteTable {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	// Get route table info
	rtResp, err := svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{})

	if err != nil {
		panic(err)
	}

	return rtResp.RouteTables
}

func scan(rt *otto.Otto, region string) otto.Value {
	tables := []ec2.RouteTable{}
	for _, t := range DescribeRouteTables(region) {
		tables = append(tables, *t)
	}
	return mcore.Sanitize(rt, tables)
}

func describeForSubnet(region string, subnetId string) []*ec2.RouteTable {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	// Get route table info
	rtResp, err := svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("association.subnet-id"),
				Values: []*string{aws.String(subnetId)},
			},
		},
	})

	if err != nil {
		panic(err)
	}

	return rtResp.RouteTables
}

func deleteAssociation(region string, id string, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.DisassociateRouteTable(&ec2.DisassociateRouteTableInput{
		AssociationId: aws.String(id),
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}

func CreateRouteTable(region string, vpcId string, verbose bool) *ec2.RouteTable {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	// Create a route table
	rtResp, err := svc.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(vpcId),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	return rtResp.RouteTable
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		var o1 *otto.Object
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			rt.Object(`aws = {}`)
		} else {
			if b, err := a.Object().Get("routeTables"); err != nil || b.IsUndefined() {
				o1, _ = rt.Object(`aws.routeTables = {}`)
			} else {
				log.Fatal("Logic error: aws.routeTables already set.")
			}
		}

		o1.Set("scan", func(region string) otto.Value {
			return scan(rt, region)
		})
		o1.Set("describeForSubnet", func(call otto.FunctionCall) otto.Value {
			f := mcore.Sanitizer(rt)
			region := call.Argument(0).String()
			subnetId := call.Argument(1).String()
			return f(describeForSubnet(region, subnetId))
		})
		o1.Set("describe", func(call otto.FunctionCall) otto.Value {
			f := mcore.Sanitizer(rt)
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			return f(describe(region, id))
		})
		o1.Set("create", func(call otto.FunctionCall) otto.Value {
			f := mcore.Sanitizer(rt)
			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()
			vpcId := call.Argument(1).String()
			return f(CreateRouteTable(region, vpcId, verbose))
		})
		o1.Set("delete", func(call otto.FunctionCall) otto.Value {
			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			deleteRouteTable(region, id, verbose)
			return otto.Value{}
		})
		o1.Set("associate", func(call otto.FunctionCall) otto.Value {
			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()
			subnetId := call.Argument(1).String()
			rtId := call.Argument(2).String()
			associate(region, rtId, subnetId, verbose)
			return otto.Value{}
		})
		o1.Set("disassociate", func(call otto.FunctionCall) otto.Value {
			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()
			associationId := call.Argument(1).String()
			disassociate(region, associationId, verbose)
			return otto.Value{}
		})
		o1.Set("deleteAssociation", func(call otto.FunctionCall) otto.Value {
			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()
			associationId := call.Argument(1).String()
			deleteAssociation(region, associationId, verbose)
			return otto.Value{}
		})
	})
}
