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
// # CORE FUNCTIONS: S3
//

package s3

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.s3.buckets.delete](#delete)
// > * [aws.s3.buckets.get](#get)
// > * [aws.s3.buckets.describe](#describe)
// > * [aws.s3.buckets.create](#create)
// > * [aws.s3.buckets.website](#website)
// > * [aws.s3.buckets.putACL](#putACL)
//
// > * [aws.s3.objects.delete](#Odelete)
// > * [aws.s3.objects.create](#Ocreate)
// > * [aws.s3.objects.describe](#Odescribe)
//
// This API allows resource handlers to manipulate S3 buckets and objects.
//
// ## AWS.S3.BUCKETS.DESCRIBE
// <a name="describe"></a>
// `aws.s3.buckets.describe(region, bucket-name);`
//
// Get info about the supplied bucket.
//
// Example:
//
// ```
//
// var bucket = aws.s3.buckets.describe("us-east-1", "mithras.io.");
//
// ```
//
// ## AWS.S3.BUCKETS.CREATE
// <a name="create"></a>
// `aws.s3.buckets.create(region, config);`
//
// Create a bucket.
//
// Example:
//
// ```
//
// var bucket = aws.s3.buckets.create("us-east-1",
// {
//                 Bucket: "my-bucket"
//                 ACL:    "public-read"
//       		       LocationConstraint: "us-east-1"
// });
//
// ```
//
// ## AWS.S3.BUCKETS.DELETE
// <a name="delete"></a>
// `aws.s3.buckets.delete(region, bucket-name);`
//
// Delete a bucket
//
// Example:
//
// ```
//
// aws.s3.buckets.delete("us-east-1", "my-bucket");
//
// ```
//
// ## AWS.S3.BUCKETS.GET
// <a name="get"></a>
// `aws.s3.buckets.get(region, bucket, key);`
//
// Get an object in a bucket.
//
// Example:
//
// ```
//
// aws.s3.buckets.get("us-east-1", "my-bucket", "index.html");
//
// ```
//
// ## AWS.S3.BUCKETS.WEBSITE
// <a name="website"></a>
// `aws.s3.buckets.website(region, config);
//
// Set up a bucket to serve a static website.
//
// Example:
//
// ```
//
// aws.s3.buckets.website("us-east-1",
//         website: {
//             Bucket: bucketName
//             WebsiteConfiguration: {
//                 ErrorDocument: {
//                     Key: "error.html"
//                 }
//                 IndexDocument: {
//                     Suffix: "index.html"
//                 }
//              }
//          });
//
// ```
//
// ## AWS.S3.BUCKETS.PUTACL
// <a name="putACL"></a>
// `aws.s3.buckets.putACL(region, config);
//
// Set up bucket access control config.
//
// Example:
//
// ```
//
// aws.s3.buckets.putACL("us-east-1",
// {
// 	   	Bucket: bucketName
// 	   	ACL:    "BucketCannedACL"
// 	   	AccessControlPolicy: {
// 	               Grants: [
// 	   		{
// 	                       Grantee: {
// 	   			Type:         "Type"
// 	   			DisplayName:  "DisplayName"
// 	   			EmailAddress: "EmailAddress"
// 	   			ID:           "ID"
// 	   			URI:          "URI"
// 	                       }
// 	                       Permission: "Permission"
// 	   		}
// 	               ]
// 	               Owner: {
// 	   		DisplayName: "DisplayName"
// 	   		ID:          "ID"
// 	               }
// 	   	}
// 	   	GrantFullControl: "GrantFullControl"
// 	   	GrantRead:        "GrantRead"
// 	   	GrantReadACP:     "GrantReadACP"
// 	   	GrantWrite:       "GrantWrite"
// 	   	GrantWriteACP:    "GrantWriteACP"
// 	});
//
// ```
//
// ## AWS.S3.OBJECTS.DELETE
// <a name="Odelete"></a>
// `aws.s3.objects.delete(region, bucket, key);`
//
// Delete an object in a bucket.
//
// Example:
//
// ```
//
// aws.s3.objects.delete("us-east-1", "my-bucket", "index.html");
//
// ```
//
// ## AWS.S3.OBJECTS.CREATE
// <a name="Ocreate"></a>
// `aws.s3.objects.create(region, config);`
//
// Create an object
//
// Example:
//
// ```
//
// aws.s3.objects.create("us-east-1",
// {
//                         Bucket:             "my-bucket"
//                         Key:                "index.html"
//                         ACL:                "public-read"
//                         Body:               "contents"
//                         ContentType:        type
// });
//
// ```
//
// ## AWS.S3.OBJECTS.DESCRIBE
// <a name="Odescribe"></a>
// `aws.s3.objects.describe(region, bucket, prefix);`
//
// Get object info.
//
// Example:
//
// ```
//
// aws.s3.objects.create("us-east-1", "my-bucket", "index.html");
//
// ```
//
import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"

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

func putACL(region string, params s3.PutBucketAclInput) {
	svc := s3.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.PutBucketAcl(&params)
	if err != nil {
		log.Fatalf("Error putting bucket acl: %s", err)
	}
}

func createBucket(region string, params s3.CreateBucketInput) string {
	svc := s3.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.CreateBucket(&params)
	if err != nil {
		log.Fatalf("Error creating bucket: %s", err)
	}

	return *resp.Location
}

func createObject(region string, params s3.PutObjectInput) s3.PutObjectOutput {
	svc := s3.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.PutObject(&params)
	if err != nil {
		log.Fatalf("Error creating object: %s", err)
	}

	return *resp
}

func deleteBucket(region, id string) {
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

func deleteObject(region, key, bucket string) {
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

func putWebsite(region string, params s3.PutBucketWebsiteInput) {
	svc := s3.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.PutBucketWebsite(&params)

	if err != nil {
		log.Fatalf("Error putting bucket website configuration: %s", err)
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

		// Buckets
		o2.Set("delete", deleteBucket)
		o2.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate target into a struct
			var input s3.CreateBucketInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for S3 createbucket input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall s3 createbucket json: %s", err)
			}
			region := call.Argument(0).String()
			return mcore.Sanitize(rt, createBucket(region, input))
		})
		o2.Set("putACL", func(call otto.FunctionCall) otto.Value {
			// Translate target into a struct
			var input s3.PutBucketAclInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for S3 putACL input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall s3 putACL json: %s", err)
			}
			region := call.Argument(0).String()
			putACL(region, input)
			return otto.Value{}
		})
		o2.Set("get", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			bucket := call.Argument(1).String()
			key := call.Argument(2).String()
			return mcore.Sanitize(rt, getObject(region, bucket, key))
		})
		o2.Set("describe", describeBucket)
		o2.Set("website", func(call otto.FunctionCall) otto.Value {
			// Translate target into a struct
			var input s3.PutBucketWebsiteInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for S3 website input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall s3 website json: %s", err)
			}
			region := call.Argument(0).String()
			putWebsite(region, input)
			return otto.Value{}
		})

		// Objects
		o3.Set("delete", deleteObject)
		o3.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate target into a struct
			var input s3.PutObjectInput
			body, err := call.Argument(1).Object().Get("Body")

			js := `(function (o) { return JSON.stringify(_.omit(o, "Body")); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for S3 putobject input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall s3 putobject json: %s", err)
			}
			if body.IsString() {
				input.Body = bytes.NewReader([]byte(body.String()))
			} else {
				// assume it is an array of JS "number"
				x, _ := body.Export()
				input.Body = bytes.NewReader(x.([]byte))
			}

			region := call.Argument(0).String()

			return mcore.Sanitize(rt, createObject(region, input))
		})
		o3.Set("describe", describeObject)
	})
}
