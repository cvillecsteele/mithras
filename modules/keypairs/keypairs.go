package keypairs

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "keypairs"

func describe(region string, id string) *ec2.KeyPairInfo {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &ec2.DescribeKeyPairsInput{
		KeyNames: []*string{
			aws.String(id),
		},
	}
	resp, err := svc.DescribeKeyPairs(params)

	if err != nil {
		return nil
	}
	if len(resp.KeyPairs) > 0 {
		return resp.KeyPairs[0]
	}
	return nil
}

func create(region string, name string) string {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &ec2.CreateKeyPairInput{
		KeyName: aws.String(name),
	}
	resp, err := svc.CreateKeyPair(params)
	if err != nil {
		log.Fatalf("Error creating key: %s", err)
	}

	return *resp.KeyMaterial
}

func delete(region string, id string) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &ec2.DeleteKeyPairInput{
		KeyName: aws.String(id),
	}
	_, err := svc.DeleteKeyPair(params)

	if err != nil {
		log.Fatalf("Error deleting key: %s", err)
	}
}

func scan(region string) []*ec2.KeyPairInfo {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &ec2.DescribeKeyPairsInput{}
	resp, err := svc.DescribeKeyPairs(params)

	if err != nil {
		log.Fatalf("Error scanning keys: %s", err)
	}

	return resp.KeyPairs
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			rt.Object(`aws = {}`)
		}
		o1, _ := rt.Object(`aws.keypairs = {}`)
		o1.Set("scan", func(region string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(scan(region))
		})
		o1.Set("create", func(region, id string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(create(region, id))
		})
		o1.Set("delete", func(region, id string) otto.Value {
			delete(region, id)
			return otto.Value{}
		})
		o1.Set("describe", func(region, id string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(describe(region, id))
		})
	})
}
