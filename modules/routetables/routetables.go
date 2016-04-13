package routetables

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
