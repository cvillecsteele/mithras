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
// # CORE FUNCTIONS: SQS
//

package sqs

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.sqs.scan](#scan)
// > * [aws.sqs.create](#create)
// > * [aws.sqs.delete](#delete)
// > * [aws.sqs.describe](#describe)
// > * [aws.sqs.attributes](#attributes)
// > * [aws.sqs.messages.send](#msend)
// > * [aws.sqs.messages.receive](#mreceive)
// > * [aws.sqs.messages.delete](#mdelete)
//
// This API allows resource handlers to manage SQS.
//
// ## AWS.SQS.SCAN
// <a name="scan"></a>
// `aws.sqs.scan(region);`
//
// Returns a list of sqs queues.
//
// Example:
//
// ```
//
//  var queues = aws.sqs.scan("us-east-1");
//
// ```
//
// ## AWS.SQS.CREATE
// <a name="create"></a>
// `aws.sqs.create(region, config);`
//
// Create a SQS queue.
//
// Example:
//
// ```
//
//  var sqs =  aws.sqs.create(
//    "us-east-1",
//    {
//      QueueName: "myqueue"
//      Attributes: [
//        "Key": "value"
//      ]
//    });
//
// ```
//
// ## AWS.SQS.DELETE
// <a name="delete"></a>
// `aws.sqs.delete(region, url);`
//
// Delete an SQS queue.
//
// Example:
//
// ```
//
//  aws.sqs.delete("us-east-1", "queueUrl");
//
// ```
//
// ## AWS.SQS.MESSAGES.SEND
// <a name="msend"></a>
// `aws.sqs.messages.send(region, input);`
//
// Publish to an SQS queue.
//
// Example:
//
// ```
//
//  aws.sqs.messages.send("us-east-1",
//                  {
//                    MessageBody:  "body"
//                    QueueUrl:     "url"
//                    DelaySeconds: 1
//                    MessageAttributes: {
//                      "Key": {
//                        DataType: "type"
//                        BinaryListValues: [
//                          "PAYLOAD"
//                        ]
//                        BinaryValue: "PAYLOAD"
//                        StringListValues: [
//                          "String"
//                        ]
//                        StringValue: "String"
//                      }
//                    }
//                   });
//
// ```
//
// ## AWS.SQS.DESCRIBE
// <a name="describe"></a>
// `aws.sqs.describe(region, sqs_id);`
//
// Get info from AWS about an SQS queue.
//
// Example:
//
// ```
//
//  var queue = aws.sqs.describe("us-east-1", "queueName");
//
// ```
//
// ## AWS.SQS.ATTRIBUTES
// <a name="attributes"></a>
// `aws.sqs.attributes(region, sqsUrl);`
//
// Get info from AWS about an SQS queue.
//
// Example:
//
// ```
//
//  var queue = aws.sqs.attributes("us-east-1", "queueUrl");
//
// ```
//
// ## AWS.SQS.MESSAGES.RECEIVE
// <a name="mreceive"></a>
// `aws.sqs.messages.receive(region, input);`
//
// Get a message from a queue
//
// Example:
//
// ```
//  var message =
//  aws.sqs.messages.receive("us-east-1",
//                  {
//                    QueueUrl: "url"
//                    AttributeNames: [
//                      "QueueAttributeName"
//                    ]
//                    MaxNumberOfMessages: 1
//                    MessageAttributeNames: [
//                      "MessageAttributeName"
//                    ]
//                    VisibilityTimeout: 1
//                    WaitTimeSeconds:   1
//                  });
//
// ```
//
// ## AWS.SQS.MESSAGES.DELETE
// <a name="mdelete"></a>
// `aws.sqs.messages.delete(region, input);`
//
// Get a message from a queue
//
// Example:
//
// ```
//  aws.sqs.messages.delete("us-east-1",
//  {
//    QueueUrl:      "queueUrl"
//    ReceiptHandle: "123456"
//  });
//
// ```
//

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "sqs"

func createQueue(region string, params *sqs.CreateQueueInput) *string {
	svc := sqs.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.CreateQueue(params)
	if err != nil {
		log.Fatalf("Error creating sqs queue: %s", err)
	}
	id := *params.QueueName

	// Wait for it.
	for i := 0; i < 10; i++ {
		url := describeQueue(region, id)
		if url != nil {
			break
		}
		time.Sleep(time.Second * 10)
	}

	return describeQueue(region, id)
}

func deleteQueue(region string, id string) {
	svc := sqs.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.DeleteQueue(&sqs.DeleteQueueInput{QueueUrl: aws.String(id)})

	if err != nil {
		log.Fatal(err.Error())
	}

	// Wait for it.
	for i := 0; i < 10; i++ {
		url := describeQueue(region, id)
		if url == nil {
			break
		}
		time.Sleep(time.Second * 10)
	}
}

func describeQueue(region string, name string) *string {
	svc := sqs.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &sqs.ListQueuesInput{
		QueueNamePrefix: aws.String(name),
	}
	resp, err := svc.ListQueues(params)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if "400" == awsErr.Code() {
				log.Println(awsErr.Error())
				return nil
			}
		}
		log.Fatal(err)
	}

	if len(resp.QueueUrls) > 0 {
		return resp.QueueUrls[0]
	}
	return nil
}

func attributesQueue(region string, url string) map[string]*string {
	svc := sqs.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(url),
		AttributeNames: []*string{
			aws.String("All"),
		},
	}
	resp, err := svc.GetQueueAttributes(params)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Attributes
}

func scanQueues(region string) []string {
	svc := sqs.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &sqs.ListQueuesInput{}
	resp, err := svc.ListQueues(params)
	if err != nil {
		log.Fatal(err)
	}

	queues := []string{}
	for _, q := range resp.QueueUrls {
		queues = append(queues, *q)
	}
	return queues
}

func sendMessage(region string, params *sqs.SendMessageInput) *sqs.SendMessageOutput {
	svc := sqs.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.SendMessage(params)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func receiveMessage(region string, params *sqs.ReceiveMessageInput) []*sqs.Message {
	svc := sqs.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.ReceiveMessage(params)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Messages
}

func deleteMessage(region string, params *sqs.DeleteMessageInput) {
	svc := sqs.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.DeleteMessage(params)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		var o1 *otto.Object
		var o2 *otto.Object
		var awsObj *otto.Object
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			awsObj, _ = rt.Object(`aws = {}`)
		} else {
			awsObj = a.Object()
		}

		if b, err := awsObj.Get("sqs"); err != nil || b.IsUndefined() {
			o1, _ = rt.Object(`aws.sqs = {}`)
		} else {
			o1 = b.Object()
		}

		if c, err := o1.Get("message"); err != nil || c.IsUndefined() {
			o2, _ = rt.Object(`aws.sqs.messages = {}`)
		} else {
			o2 = c.Object()
		}

		// Queues
		o1.Set("scan", func(region string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(scanQueues(region))
		})
		o1.Set("delete", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			deleteQueue(region, id)
			return otto.Value{}
		})
		o1.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input sqs.CreateQueueInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for SQS queue create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall SQS create queue json: %s", err)
			}

			region := call.Argument(0).String()

			f := mcore.Sanitizer(rt)
			return f(createQueue(region, &input))
		})
		o1.Set("describe", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			sqsId := call.Argument(1).String()
			f := mcore.Sanitizer(rt)
			return f(describeQueue(region, sqsId))
		})
		o1.Set("attributes", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			sqsId := call.Argument(1).String()
			f := mcore.Sanitizer(rt)
			return f(attributesQueue(region, sqsId))
		})

		// Messages
		o2.Set("send", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input sqs.SendMessageInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for SQS publish input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall SQS publish json: %s", err)
			}

			region := call.Argument(0).String()

			f := mcore.Sanitizer(rt)
			return f(sendMessage(region, &input))
		})
		o2.Set("receive", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input sqs.ReceiveMessageInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for SQS receive message input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall SQS receieve message json: %s", err)
			}

			region := call.Argument(0).String()

			f := mcore.Sanitizer(rt)
			return f(receiveMessage(region, &input))
		})
		o2.Set("delete", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input sqs.DeleteMessageInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for SQS delete message input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall SQS delete message json: %s", err)
			}

			region := call.Argument(0).String()

			deleteMessage(region, &input)
			return otto.Value{}
		})
	})
}
