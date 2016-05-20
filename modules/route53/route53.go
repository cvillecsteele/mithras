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
// # CORE FUNCTIONS: ROUTE53
//

package route53

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.route53.zones.scan](#zscan)
//
// This API allows resource handlers to manipulate DNS records in Route53.
//
// ## AWS.ROUTE53.ZONES.SCAN
// <a name="zscan"></a>
// `aws.route53.zones.scan(region);`
//
// Query zones in Route53.
//
// Example:
//
// ```
//
// var zones = aws.route53.zones.scan("us-east-1");
//
// ```
//
// ## AWS.ROUTE53.RRS.SCAN
// <a name="scan"></a>
// `aws.route53.rrs.scan(region);`
//
// Query resource records sets (RRSs) in Route53.
//
// Example:
//
// ```
//
// var rrs = aws.route53.rrs.scan("us-east-1");
//
// ```
//
// ## AWS.ROUTE53.RRS.CREATE
// <a name="create"></a>
// `aws.route53.rrs.create(region, zoneId, config);`
//
// Create dns records.
//
// Example:
//
// ```
//
// var rrs = aws.route53.rrs.create("us-east-1", "Z111111...",
// {
// 		Name:         "mithras.io."
// 		Type:         "A"
// 		AliasTarget: {
// 		    DNSName:              "s3-website-us-east-1.amazonaws.com"
// 		    EvaluateTargetHealth: false
// 		    HostedZoneId:         "Z3AQBSTGFYJSTF"
// 		}
// });
//
// ```
//
// ## AWS.ROUTE53.RRS.DELETE
// <a name="delete"></a>
// `aws.route53.rrs.delete(region, zoneId, config);`
//
// Delete dns records.
//
// Example:
//
// ```
//
// var rrs = aws.route53.rrs.delete("us-east-1", "Z111111...",
// {
// 		Name:         "mithras.io."
// 		Type:         "A"
// 		AliasTarget: {
// 		    DNSName:              "s3-website-us-east-1.amazonaws.com"
// 		    EvaluateTargetHealth: false
// 		    HostedZoneId:         "Z3AQBSTGFYJSTF"
// 		}
// });
//
// ```
//
// ## AWS.ROUTE53.RRS.DESCRIBE
// <a name="describe"></a>
// `aws.route53.rrs.describe(region, zoneId, name, type);`
//
// Get info about dns records.
//
// Example:
//
// ```
//
// var rrs = aws.route53.rrs.describe("us-east-1", "Z111111...", "mithras.io." "A");
//
// ```
//
import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
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
	mcore.RegisterInit(func(context *mcore.Context) {
		rt := context.Runtime

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
