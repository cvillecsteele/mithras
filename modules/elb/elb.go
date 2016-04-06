package elb

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

	if verbose {
		log.Printf("  ### Modifying health check for elb '%s'\n", lbName)
	}

	_, err := svc.ConfigureHealthCheck(&input)
	if err != nil {
		log.Fatalf("Can't configure elb health check: %s", err)
	}
}

func setAttrs(region string, lbName string, input elb.ModifyLoadBalancerAttributesInput, verbose bool) {
	svc := elb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	if verbose {
		log.Printf("  ### Modifying attributes for elb '%s'\n", lbName)
	}

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

	if verbose {
		log.Printf("  ### Creating ELB '%s'\n", *params.LoadBalancerName)
	}

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

	if verbose {
		log.Printf("  ### Deleting ELB '%s'\n", id)
	}

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

	if verbose {
		log.Printf("  ### Adding %d instances to elb '%s'\n", len(instances), lbName)
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

	if verbose {
		log.Printf("  ### Removing %d instances from elb '%s'\n", len(instances), lbName)
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
