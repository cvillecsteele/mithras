package secgroup

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
var ModuleName = "secgroup"

func describe(region string, id string) *ec2.SecurityGroup {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String(id),
		},
	}
	resp, err := svc.DescribeSecurityGroups(params)

	if err != nil {
		return nil
	}

	return resp.SecurityGroups[0]
}

func authorizeIngress(region string, params *ec2.AuthorizeSecurityGroupIngressInput, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	if verbose {
		log.Printf("  ### Creating security group authorization for ingress")
	}

	_, err := svc.AuthorizeSecurityGroupIngress(params)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func authorizeEgress(region string, params *ec2.AuthorizeSecurityGroupEgressInput, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	if verbose {
		log.Printf("  ### Creating security group authorization for egress")
	}

	_, err := svc.AuthorizeSecurityGroupEgress(params)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func create(region string, params *ec2.CreateSecurityGroupInput, verbose bool) *ec2.SecurityGroup {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	if verbose {
		log.Printf("  ### Creating security group '%s'", *params.GroupName)
	}

	resp, err := svc.CreateSecurityGroup(params)
	if err != nil {
		log.Fatal(err.Error())
	}
	id := *resp.GroupId

	// Wait for it.
	for i := 0; i < 10; i++ {
		target := describe(region, id)
		if target != nil && *target.GroupId == id {
			break
		}
		time.Sleep(time.Second * 10)
	}

	// Describe it.
	target := describe(region, id)

	return target
}

func delete(region string, id string, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	if verbose {
		log.Printf("  ### Deleting security group '%s'", id)
	}

	i := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(id),
	}
	_, err := svc.DeleteSecurityGroup(i)

	if err != nil {
		log.Fatal(err.Error())
	}

	// TODO: Wait for it.
}

func scan(rt *otto.Otto, region string) otto.Value {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.DescribeSecurityGroups(nil)
	if err != nil {
		panic(err)
	}

	// shove instances into jsland
	sgs := []ec2.SecurityGroup{}
	for _, i := range resp.SecurityGroups {
		sgs = append(sgs, *i)
	}
	return mcore.Sanitize(rt, sgs)
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			rt.Object(`aws = {}`)
		}
		o1, _ := rt.Object(`aws.securityGroups = {}`)
		o1.Set("scan", func(region string) otto.Value {
			return scan(rt, region)
		})
		o1.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input ec2.CreateSecurityGroupInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for secgroup create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall secgroup create json: %s", err)
			}

			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()

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
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			f := mcore.Sanitizer(rt)
			return f(describe(region, id))
		})
		o1.Set("authorizeIngress", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input ec2.AuthorizeSecurityGroupIngressInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for secgroup create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall secgroup ingress json: %s", err)
			}

			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()
			authorizeIngress(region, &input, verbose)
			return otto.Value{}
		})
		o1.Set("authorizeEgress", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input ec2.AuthorizeSecurityGroupEgressInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for secgroup create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall secgroup egress json: %s", err)
			}

			verbose := mcore.IsVerbose(rt)
			region := call.Argument(0).String()
			authorizeEgress(region, &input, verbose)
			return otto.Value{}
		})
	})
}
