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
package s3

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "s3"

func getObject(region string, bucket string, key string) s3.GetObjectOutput {
	svc := s3.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	resp, err := svc.GetObject(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if "NoSuchBucket" == awsErr.Code() {
				return *resp
			}
		}
		log.Fatalf("Error getting object: %s", err)
	}

	return *resp
}

func describeObject(region string, bucket string, prefix string) []*s3.Object {
	svc := s3.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	// TODO: paginate
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
		// Delimiter:    aws.String("Delimiter"),
		// EncodingType: aws.String("EncodingType"),
		// Marker:       aws.String("Marker"),
		// MaxKeys:      aws.Int64(1),
	}
	resp, err := svc.ListObjects(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if "NoSuchBucket" == awsErr.Code() {
				return resp.Contents
			}
		}
		log.Fatalf("Error listing objects: %s", err)
	}

	return resp.Contents
}

func describeBucket(region string, bucket string) []*s3.Bucket {
	svc := s3.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	var params *s3.ListBucketsInput
	resp, err := svc.ListBuckets(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if "404" == awsErr.Code() {
				return []*s3.Bucket{}
			}
		}
		log.Fatalf("Error listing bucket: %s", err)
	}

	if bucket != "*" {
		for _, b := range resp.Buckets {
			if *b.Name == bucket {
				return []*s3.Bucket{b}
			}
		}
	}

	return resp.Buckets
}

func createBucket(params s3.CreateBucketInput, region string, verbose bool) string {
	svc := s3.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.CreateBucket(&params)
	if err != nil {
		log.Fatalf("Error creating bucket: %s", err)
	}

	return *resp.Location
}

func createObject(params s3.PutObjectInput, region string, verbose bool) s3.PutObjectOutput {
	svc := s3.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.PutObject(&params)
	if err != nil {
		log.Fatalf("Error creating object: %s", err)
	}

	return *resp
}

func deleteBucket(id string, region string) {
	svc := s3.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &s3.DeleteBucketInput{
		Bucket: aws.String(id),
	}
	_, err := svc.DeleteBucket(params)

	if err != nil {
		log.Fatalf("Error deleting bucket: %s", err)
	}
}

func deleteObject(bucket string, key string, region string) {
	svc := s3.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		// MFA:          aws.String("MFA"),
		// RequestPayer: aws.String("RequestPayer"),
		// VersionId:    aws.String("ObjectVersionId"),
	}
	_, err := svc.DeleteObject(params)

	if err != nil {
		log.Fatalf("Error deleting bucket: %s", err)
	}
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			rt.Object(`aws = {}`)
		}
		rt.Object(`aws.s3 = {}`)
		o2, _ := rt.Object(`aws.s3.buckets = {}`)
		o3, _ := rt.Object(`aws.s3.objects = {}`)

		// Objects
		o2.Set("delete", deleteBucket)
		o2.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate target into a struct
			var input s3.CreateBucketInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(0))
			if err != nil {
				log.Fatalf("Can't create json for S3 createbucket input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall s3 createbucket json: %s", err)
			}
			region := call.Argument(1).String()
			verbose, _ := call.Argument(2).ToBoolean()
			return mcore.Sanitize(rt, createBucket(input, region, verbose))
		})
		o2.Set("get", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			bucket := call.Argument(1).String()
			key := call.Argument(2).String()
			return mcore.Sanitize(rt, getObject(region, bucket, key))
		})
		o2.Set("describe", describeBucket)

		// Buckets
		o3.Set("delete", deleteObject)
		o3.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate target into a struct
			var input s3.PutObjectInput
			body, err := call.Argument(0).Object().Get("Body")
			js := `(function (o) { return JSON.stringify(_.omit(o, "Body")); })`
			s, err := rt.Call(js, nil, call.Argument(0))
			if err != nil {
				log.Fatalf("Can't create json for S3 putobject input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall s3 putobject json: %s", err)
			}
			input.Body = bytes.NewReader([]byte(body.String()))

			region := call.Argument(1).String()
			verbose, _ := call.Argument(2).ToBoolean()

			return mcore.Sanitize(rt, createObject(input, region, verbose))
		})
		o3.Set("describe", describeObject)
	})
}
