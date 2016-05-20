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
package keypairs

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.keypairs.scan](#vscan)
// > * [aws.keypairs.create](#vcreate)
// > * [aws.keypairs.delete](#vdelete)
// > * [aws.keypairs.describe](#vdescribe)
//
// This API allows resource handlers to manage AWS SSH keypairs.
//
// ## AWS.KEYPAIRS.SCAN
// <a name="vscan"></a>
// `aws.keypairs.scan(region);`
//
// Returns a list of keypairs.
//
// Example:
//
// ```
//
//  var keypairs =  aws.keypairs.scan("us-east-1");
//
// ```
//
// ## AWS.KEYPAIRS.CREATE
// <a name="vcreate"></a>
// `aws.keypairs.create(region, config);`
//
// Create a keypair.
//
// Example:
//
// ```
//
//  var keypair =  aws.keypairs.create(
//    "us-east-1",
//    {
//		  KeyName: "my-keypair"
//    });
//
// ```
//
// ## AWS.KEYPAIRS.DELETE
// <a name="vdelete"></a>
// `aws.keypairs.delete(region, keypair_id);`
//
// Delete a keypair.
//
// Example:
//
// ```
//
//  aws.keypairs.delete("us-east-1", "my-keypair");
//
// ```
//
// ## AWS.KEYPAIRS.DESCRIBE
// <a name="vdescribe"></a>
// `aws.keypairs.describe(region, keypair_id);`
//
// Get info from AWS about a KEYPAIR.
//
// Example:
//
// ```
//
//  var keypair = aws.keypairs.describe("us-east-1", "my-keypair");
//
// ```
//

import (
	log "github.com/Sirupsen/logrus"

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
	mcore.RegisterInit(func(context *mcore.Context) {
		rt := context.Runtime

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
