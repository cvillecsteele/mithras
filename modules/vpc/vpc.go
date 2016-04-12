//
// # CORE FUNCTIONS: VPC
//

package vpc

// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.vpcs.scan](#vscan)
// > * [aws.vpcs.create](#vcreate)
// > * [aws.vpcs.delete](#vdelete)
// > * [aws.vpcs.describe](#vdescribe)
// > * [aws.gateways.scan](#gscan)
// > * [aws.gateways.create](#gcreate)
// > * [aws.gateways.delete](#gdelete)
// > * [aws.gateways.describe](#gdescribe)
//
// This API allows resource handlers to manage VPCS.
//
// ## AWS.VPCS.SCAN
// <a name="vscan"></a>
// `aws.vpcs.scan(region);`
//
// Returns a list of vpcs.
//
// Example:
//
// ```
//
//  var vpcs =  aws.vpcs.scan("us-east-1");
//
// ```
//
// ## AWS.VPCS.CREATE
// <a name="vcreate"></a>
// `aws.vpcs.create(region, config, gateway);`
//
// Create a VPC.
//
// Example:
//
// ```
//
//  var vpc =  aws.vpcs.create(
//    "us-east-1",
//    {
//		  CidrBlock:       "172.33.0.0/16"
//    },
//    true);
//
// ```
//
// ## AWS.VPCS.DELETE
// <a name="vdelete"></a>
// `aws.vpcs.delete(region, vpc_id);`
//
// Delete a VPC.
//
// Example:
//
// ```
//
//  aws.vpcs.delete("us-east-1", "vpc-abcd");
//
// ```
//
// ## AWS.VPCS.DESCRIBE
// <a name="vdescribe"></a>
// `aws.vpcs.describe(region, vpc_id);`
//
// Get info from AWS about a VPC.
//
// Example:
//
// ```
//
//  var vpc = aws.vpcs.describe("us-east-1", "vpc-abcd");
//
// ```
//
// ## AWS.GATEWAYS.SCAN
// <a name="gscan"></a>
// `aws.gateways.scan(region);`
//
// Returns a list of gateways.
//
// Example:
//
// ```
//
//  var gateways =  mithras.gateways.scan("us-east-1");
//
// ```
//
// ## AWS.GATEWAYS.CREATE
// <a name="gcreate"></a>
// `aws.gateways.create(region);`
//
// Create a gateway.
//
// Example:
//
// ```
//
//  var gateway =  mithras.gateways.create("us-east-1");
//
// ```
//
// ## AWS.GATEWAYS.DELETE
// <a name="gdelete"></a>
// `aws.gateways.delete(region, gateway_id);`
//
// Delete a gateway.
//
// Example:
//
// ```
//
//  mithras.gateways.delete("us-east-1", "gw-abcd");
//
// ```
//
// ## AWS.GATEWAYS.DESCRIBE
// <a name="gdescribe"></a>
// `aws.gateways.describe(region, gateway_id);`
//
// Get info from AWS about a gateway.
//
// Example:
//
// ```
//
//  var gateway = mithras.gateways.describe("us-east-1", "gw-abcd");
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
var ModuleName = "vpc"

func createGateway(region string) ec2.InternetGateway {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	p := &ec2.CreateInternetGatewayInput{}
	resp, err := svc.CreateInternetGateway(p)
	if err != nil {
		log.Fatal(err.Error())
	}
	return *resp.InternetGateway
}

func createVPC(params *ec2.CreateVpcInput, region string, gateway bool, verbose bool) (ec2.Vpc, ec2.InternetGateway) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.CreateVpc(params)
	if err != nil {
		log.Fatalf("Error creating vpc: %s", err)
	}
	id := *resp.Vpc.VpcId

	// Wait for it.
	for i := 0; i < 10; i++ {
		vpc := describeVPC(region, id)
		if vpc.State != nil && *vpc.State == "available" {
			break
		}
		time.Sleep(time.Second * 10)
	}

	modParams := &ec2.ModifyVpcAttributeInput{
		VpcId: aws.String(id),
		EnableDnsHostnames: &ec2.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
	}
	if _, err := svc.ModifyVpcAttribute(modParams); err != nil {
		log.Fatalf("Error modifying vpc '%s': %s", id, err)
	}

	// Create gateway
	var gw ec2.InternetGateway
	if gateway {
		gw = createGateway(region)
		// Now attach it.
		params := &ec2.AttachInternetGatewayInput{
			InternetGatewayId: aws.String(*gw.InternetGatewayId),
			VpcId:             aws.String(id),
		}
		_, err := svc.AttachInternetGateway(params)
		if err != nil {
			log.Fatal(err.Error())
		}
		gw = describeGW(region, *gw.InternetGatewayId)
	}
	return describeVPC(region, id), gw
}

func deleteVPC(region string, id string, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.DeleteVpc(&ec2.DeleteVpcInput{VpcId: aws.String(id)})

	if err != nil {
		log.Fatal(err.Error())
	}

	// TODO: Wait for it.
}

func deleteGW(region string, vpcId string, gwId string, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	if verbose {
		log.Printf("  ### Deleting GW '%s'\n", gwId)
	}

	_, err := svc.DetachInternetGateway(&ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String(gwId),
		VpcId:             aws.String(vpcId),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = svc.DeleteInternetGateway(&ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(gwId),
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}

func describeVPC(region string, vpcId string) ec2.Vpc {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String(vpcId)},
	}
	resp, err := svc.DescribeVpcs(params)

	if err != nil {
		return ec2.Vpc{}
	}

	return *resp.Vpcs[0]
}

func describeGW(region string, id string) ec2.InternetGateway {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{aws.String(id)},
	}
	resp, err := svc.DescribeInternetGateways(params)

	if err != nil {
		return ec2.InternetGateway{}
	}

	return *resp.InternetGateways[0]
}

func scanVPCs(rt *otto.Otto, region string) otto.Value {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.DescribeVpcs(nil)
	if err != nil {
		panic(err)
	}

	// shove instances into jsland
	vpcs := []ec2.Vpc{}
	for _, i := range resp.Vpcs {
		vpcs = append(vpcs, *i)
	}
	return mcore.Sanitize(rt, vpcs)
}

func scanGateways(rt *otto.Otto, region string) otto.Value {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.DescribeInternetGateways(nil)
	if err != nil {
		panic(err)
	}

	// shove instances into jsland
	gws := []ec2.InternetGateway{}
	for _, i := range resp.InternetGateways {
		gws = append(gws, *i)
	}
	return mcore.Sanitize(rt, gws)
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

		if b, err := awsObj.Get("vpcs"); err != nil || b.IsUndefined() {
			o1, _ = rt.Object(`aws.vpcs = {}`)
		} else {
			o1 = b.Object()
		}

		// VPCs
		o1.Set("scan", func(region string) otto.Value {
			return scanVPCs(rt, region)
		})
		o1.Set("delete", func(call otto.FunctionCall) otto.Value {
			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			deleteVPC(region, id, verbose)
			return otto.Value{}
		})
		o1.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input ec2.CreateVpcInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for VPC create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall VPC create json: %s", err)
			}

			region := call.Argument(0).String()
			gateway, err := call.Argument(2).ToBoolean()
			if err != nil {
				log.Fatalf("Invalid gateway arg to VPC create: %s", err)
			}
			verbose := mcore.IsVerbose(rt)

			f := mcore.Sanitizer(rt)
			return f(createVPC(&input, region, gateway, verbose))
		})
		o1.Set("describe", func(call otto.FunctionCall) otto.Value {
			f := mcore.Sanitizer(rt)
			region := call.Argument(0).String()
			vpcId := call.Argument(1).String()
			return f(describeVPC(region, vpcId))
		})

		// Internet Gatways
		o2, _ := rt.Object(`aws.vpcs.gateways = {}`)
		o2.Set("scan", func(region string) otto.Value {
			return scanGateways(rt, region)
		})
		o2.Set("create", func(call otto.FunctionCall) otto.Value {
			f := mcore.Sanitizer(rt)
			region := call.Argument(0).String()
			return f(createGateway(region))
		})
		// TODO: add associate function
		o2.Set("delete", func(call otto.FunctionCall) otto.Value {
			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()
			vpcId := call.Argument(1).String()
			gwId := call.Argument(2).String()
			deleteGW(region, vpcId, gwId, verbose)
			return otto.Value{}
		})
		o2.Set("describe", func(call otto.FunctionCall) otto.Value {
			f := mcore.Sanitizer(rt)
			region := call.Argument(0).String()
			vpcId := call.Argument(1).String()
			return f(describeGW(region, vpcId))
		})
	})
}
