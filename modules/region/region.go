package region

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "region"

func scanRegions(rt *otto.Otto, cb otto.Value) {
	// Need a context (this) for Call below
	ctx, _ := rt.Get("mithras")
	svc := ec2.New(session.New(), aws.NewConfig().WithRegion("us-east-1").WithMaxRetries(5))

	resp, err := svc.DescribeRegions(nil)
	if err != nil {
		panic(err)
	}

	for idx, _ := range resp.Regions {
		cb.Call(ctx, *resp.Regions[idx].RegionName)
	}
}

func scan(rt *otto.Otto) otto.Value {
	svc := ec2.New(session.New(), aws.NewConfig().WithRegion("us-east-1").WithMaxRetries(5))

	resp, err := svc.DescribeRegions(nil)
	if err != nil {
		panic(err)
	}

	regions := []ec2.Region{}
	for _, r := range resp.Regions {
		regions = append(regions, *r)
	}
	return mcore.Sanitize(rt, regions)
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		rt.Object(`aws = aws || {}`)
		o2, _ := rt.Object(`aws.regions = {}`)
		o2.Set("scan", func() otto.Value {
			return scan(rt)
		})
	})
}
