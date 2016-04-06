package route53

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "route53"

func allRRSs(region string) []*route53.ResourceRecordSet {
	svc := route53.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.ListHostedZones(nil)
	if err != nil {
		log.Fatal("Can't list zones: %s", err)
	}

	all := []*route53.ResourceRecordSet{}
	for _, hz := range resp.HostedZones {
		id := hz.Id
		params := &route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String(*id),
			MaxItems:     aws.String("100"),
		}
		resp, err := svc.ListResourceRecordSets(params)
		if err != nil {
			log.Fatal("Can't list resource set in zone '%s': %s", id, err)
		}

		all = append(all, resp.ResourceRecordSets...)
	}

	return all
}

func describe(region string, zoneId string, rName string, rType string) *route53.ResourceRecordSet {
	svc := route53.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneId),
		StartRecordName: aws.String(rName),
		StartRecordType: aws.String(rType),
		MaxItems:        aws.String("1"),
	}
	resp, err := svc.ListResourceRecordSets(params)

	if err != nil {
		log.Fatalf("Error describing route53 resources: %s", err)
	}

	if len(resp.ResourceRecordSets) > 0 {
		return resp.ResourceRecordSets[0]
	}
	return nil
}

func modify(region string, rrs *route53.ChangeResourceRecordSetsInput, zoneId string, verbose bool) *route53.ResourceRecordSet {
	svc := route53.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.ChangeResourceRecordSets(rrs)
	if err != nil {
		log.Fatalf("Error creating resource record: %s", err)
	}

	rName := rrs.ChangeBatch.Changes[0].ResourceRecordSet.Name
	rType := rrs.ChangeBatch.Changes[0].ResourceRecordSet.Type

	// Wait for it.
	var target *route53.ResourceRecordSet
	for i := 0; i < 10; i++ {
		target = describe(region, zoneId, *rName, *rType)
		if target != nil && *target.Name == *rName && *target.Type == *rType {
			break
		}
		time.Sleep(time.Second * 10)
	}

	return target
}

func scanZones(rt *otto.Otto, region string) otto.Value {
	svc := route53.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.ListHostedZones(nil)
	if err != nil {
		log.Fatal("Can't list zones: %s", err)
	}

	zones := []route53.HostedZone{}
	for _, z := range resp.HostedZones {
		zones = append(zones, *z)
	}
	return mcore.Sanitize(rt, zones)
}

func scanResources(rt *otto.Otto, region string) otto.Value {
	all := allRRSs(region)
	rrs := []route53.ResourceRecordSet{}
	for _, r := range all {
		rrs = append(rrs, *r)
	}
	return mcore.Sanitize(rt, rrs)
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			rt.Object(`aws = {}`)
		}
		rt.Object(`aws.route53 = {}`)

		o1, _ := rt.Object(`aws.route53.zones = {}`)
		o1.Set("scan", func(region string) otto.Value {
			return scanZones(rt, region)
		})

		o2, _ := rt.Object(`aws.route53.rrs = {}`)
		gimmieChange := func(action string) func(otto.FunctionCall) otto.Value {
			return func(call otto.FunctionCall) otto.Value {
				// Translate params input into a struct
				var input route53.ResourceRecordSet
				js := `(function (o) { return JSON.stringify(o); })`
				s, err := rt.Call(js, nil, call.Argument(2))
				if err != nil {
					log.Fatalf("Can't create json for route53 resource input: %s", err)
				}
				err = json.Unmarshal([]byte(s.String()), &input)
				if err != nil {
					log.Fatalf("Can't unmarshall route53 resource json: %s", err)
				}

				verbose := mcore.IsVerbose(rt)
				region := call.Argument(0).String()
				zoneId := call.Argument(1).String()

				changes := []*route53.Change{}
				change := route53.Change{
					Action:            aws.String(action),
					ResourceRecordSet: &input,
				}
				changes = append(changes, &change)
				changeInput := route53.ChangeResourceRecordSetsInput{
					ChangeBatch: &route53.ChangeBatch{
						Changes: changes,
					},
					HostedZoneId: aws.String(zoneId),
				}

				f := mcore.Sanitizer(rt)
				return f(modify(region, &changeInput, zoneId, verbose))
			}
		}

		o2.Set("scan", func(region string) otto.Value {
			return scanResources(rt, region)
		})
		o2.Set("create", gimmieChange("CREATE"))
		o2.Set("delete", gimmieChange("DELETE"))
		o2.Set("describe", func(region string, zoneId string, rName string, rType string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(describe(region, zoneId, rName, rType))
		})

	})
}
