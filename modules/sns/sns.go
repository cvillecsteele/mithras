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
// # CORE FUNCTIONS: SNS
//

package sns

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.sns.scan](#vscan)
// > * [aws.sns.create](#vcreate)
// > * [aws.sns.delete](#vdelete)
// > * [aws.sns.describe](#vdescribe)
//
// This API allows resource handlers to manage SNS.
//
// ## AWS.SNS.SCAN
// <a name="vscan"></a>
// `aws.sns.scan(region);`
//
// Returns a list of sns.
//
// Example:
//
// ```
//
//  var sns =  aws.sns.scan("us-east-1");
//
// ```
//
// ## AWS.SNS.CREATE
// <a name="vcreate"></a>
// `aws.sns.create(region, config, gateway);`
//
// Create a SNS.
//
// Example:
//
// ```
//
//  var sns =  aws.sns.create(
//    "us-east-1",
//    {
//		  CidrBlock:       "172.33.0.0/16"
//    },
//    true);
//
// ```
//
// ## AWS.SNS.DELETE
// <a name="vdelete"></a>
// `aws.sns.delete(region, sns_id);`
//
// Delete a SNS.
//
// Example:
//
// ```
//
//  aws.sns.delete("us-east-1", "sns-abcd");
//
// ```
//
// ## AWS.SNS.DESCRIBE
// <a name="vdescribe"></a>
// `aws.sns.describe(region, sns_id);`
//
// Get info from AWS about a SNS.
//
// Example:
//
// ```
//
//  var sns = aws.sns.describe("us-east-1", "sns-abcd");
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
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "sns"

func createTopic(region string, params *sns.CreateTopicInput) string {
	svc := sns.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.CreateTopic(params)
	if err != nil {
		log.Fatalf("Error creating sns topic: %s", err)
	}
	id := *resp.TopicArn

	// Wait for it.
	for i := 0; i < 10; i++ {
		sns := describeTopic(region, id)
		if sns != "" {
			break
		}
		time.Sleep(time.Second * 10)
	}

	return describeTopic(region, id)
}

func deleteTopic(region string, id string) {
	svc := sns.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.DeleteTopic(&sns.DeleteTopicInput{TopicArn: aws.String(id)})

	if err != nil {
		log.Fatal(err.Error())
	}

	// TODO: Wait for it.
}

func describeTopic(region string, snsId string) string {
	for _, topic := range scanTopics(region) {
		if topic == snsId {
			return topic
		}
	}
	return ""
}

func scanTopics(region string) []string {
	svc := sns.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &sns.ListTopicsInput{}
	resp, err := svc.ListTopics(params)
	if err != nil {
		log.Fatal(err)
	}

	topics := []string{}
	for _, t := range resp.Topics {
		topics = append(topics, *t.TopicArn)
	}
	return topics
}

func createSubscription(region string, params *sns.SubscribeInput) *sns.Subscription {
	svc := sns.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.Subscribe(params)
	if err != nil {
		log.Fatalf("Error creating sns subscription: %s", err)
	}
	id := *resp.SubscriptionArn

	return describeSubscription(region, id)
}

func deleteSubscription(region string, id string) {
	svc := sns.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.Unsubscribe(&sns.UnsubscribeInput{SubscriptionArn: aws.String(id)})

	if err != nil {
		log.Fatal(err.Error())
	}

	// TODO: Wait for it.
}

func describeSubscription(region string, subArn string) *sns.Subscription {
	for _, subscription := range scanSubscriptions(region) {
		if *subscription.SubscriptionArn == subArn {
			return subscription
		}
	}
	return nil
}

func scanSubscriptions(region string) []*sns.Subscription {
	svc := sns.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &sns.ListSubscriptionsInput{}
	resp, err := svc.ListSubscriptions(params)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Subscriptions
}

func publish(region string, params *sns.PublishInput) *string {
	svc := sns.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.Publish(params)
	if err != nil {
		log.Fatal(err)
	}

	return resp.MessageId
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		var o1 *otto.Object
		var o2 *otto.Object
		var o3 *otto.Object
		var awsObj *otto.Object
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			awsObj, _ = rt.Object(`aws = {}`)
		} else {
			awsObj = a.Object()
		}

		if b, err := awsObj.Get("sns"); err != nil || b.IsUndefined() {
			o1, _ = rt.Object(`aws.sns = {}`)
		} else {
			o1 = b.Object()
		}

		if c, err := o1.Get("topics"); err != nil || c.IsUndefined() {
			o2, _ = rt.Object(`aws.sns.topics = {}`)
		} else {
			o2 = c.Object()
		}

		if d, err := o1.Get("subs"); err != nil || d.IsUndefined() {
			o3, _ = rt.Object(`aws.sns.subs = {}`)
		} else {
			o3 = d.Object()
		}

		// Topics
		o2.Set("scan", func(region string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(scanTopics(region))
		})
		o2.Set("delete", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			deleteTopic(region, id)
			return otto.Value{}
		})
		o2.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input sns.CreateTopicInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for SNS topic create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall SNS create topic json: %s", err)
			}

			region := call.Argument(0).String()

			f := mcore.Sanitizer(rt)
			return f(createTopic(region, &input))
		})
		o2.Set("describe", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			snsId := call.Argument(1).String()
			f := mcore.Sanitizer(rt)
			return f(describeTopic(region, snsId))
		})

		// Subs
		o3.Set("scan", func(region string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(scanSubscriptions(region))
		})
		o3.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input sns.SubscribeInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for SNS subs create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall SNS create subs json: %s", err)
			}

			region := call.Argument(0).String()

			f := mcore.Sanitizer(rt)
			return f(createSubscription(region, &input))
		})
		o3.Set("delete", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			snsId := call.Argument(1).String()
			deleteSubscription(region, snsId)
			return otto.Value{}
		})
		o3.Set("describe", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			snsId := call.Argument(1).String()
			f := mcore.Sanitizer(rt)
			return f(describeSubscription(region, snsId))
		})
	})
}
