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

//
// # CORE FUNCTIONS: RDS
//

package rds

// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.rds.scan](#scan)
// > * [aws.rds.create](#create)
// > * [aws.rds.delete](#delete)
// > * [aws.rds.describe](#describe)

// > * [aws.rds.subnetGroups.create](#gcreate)
// > * [aws.rds.subnetGroups.delete](#gdelete)
// > * [aws.rds.subnetGroups.describe](#gdescribe)
//
// This API allows resource handlers to manage RDS.
//
// ## AWS.RDS.SCAN
// <a name="scan"></a>
// `aws.rds.scan(region);`
//
// Returns a list of RDS clusters.
//
// Example:
//
// ```
//
//  var dbs = aws.rds.scan("us-east-1");
//
// ```
//
// ## AWS.RDS.CREATE
// <a name="create"></a>
// `aws.rds.create(region, config, wait);`
//
// Create an RDS cluster.
//
// Example:
//
// ```
//
//  var db = aws.rds.create("us-east-1",
//   {
//      DBInstanceClass:         "db.m1.small"
//      DBInstanceIdentifier:    "test-rds"
//      Engine:                  "mysql"
//      AllocatedStorage:        10
//      AutoMinorVersionUpgrade: true
//      AvailabilityZone:        "us-east-1"
//      MasterUserPassword:      "test123456789"
//      MasterUsername:          "test"
//      DBSubnetGroupName:       "test-subnet-group"
//      DBName:                  "test"
//      PubliclyAccessible:      false
//      Tags: [
//          {
//             Key:   "foo"
//             Value: "bar"
//          },
//      ]
//   },
//   true);
//
// ```
//
// ## AWS.RDS.DELETE
// <a name="delete"></a>
// `aws.rds.delete(region, config);`
//
// Delete an RDS cluster.
//
// Example:
//
// ```
//
//  var db = aws.rds.delete("us-east-1",
//     {
//		   DBInstanceIdentifier:      "db-xyz",
//		   FinalDBSnapshotIdentifier: "byebye" + Date.now()
//		   SkipFinalSnapshot:         true
//     });
//
// ```
//
// ## AWS.RDS.DESCRIBE
// <a name="describe"></a>
// `aws.rds.describe(region, id);`
//
// Get info about an RDS cluster.
//
// Example:
//
// ```
//
//  var db = aws.rds.describe("us-east-1", "db-xyz");
//
// ```
//
// ## AWS.RDS.SUBNETGROUPS.DESCRIBE
// <a name="gdescribe"></a>
// `aws.rds.subnetGroups.describe(region, id);`
//
// Get info about an RDS subnet group.
//
// Example:
//
// ```
//
//  var group = aws.rds.subnetGroups.describe("us-east-1", "sg-xyz");
//
// ```
//
// ## AWS.RDS.SUBNETGROUPS.CREATE
// <a name="gcreate"></a>
// `aws.rds.subnetGroups.create(region, config);`
//
// Create an RDS subnet group.
//
// Example:
//
// ```
//
//  var group = aws.rds.subnetGroups.create("us-east-1",
// {
// 		DBSubnetGroupDescription: "test subnet group"
// 		DBSubnetGroupName: "test-subnet-group"
// 		SubnetIds: [
//       "subnet-1"
//       "subnet-2"
// 		]
// 		Tags: [
// 		    {
// 			     Key:   "Foo"
// 			     Value: "Bar"
// 		    }
// 		]
// });
//
// ```
//
// ## AWS.RDS.SUBNETGROUPS.DELETE
// <a name="gdelete"></a>
// `aws.rds.subnetGroups.delete(region, id);`
//
// Delete an RDS subnet group.
//
// Example:
//
// ```
//
//  var group = aws.rds.subnetGroups.delete("us-east-1", "sg-xyz");
//
// ```
//
import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "rds"

func describe(region string, id string) *rds.DBInstance {
	svc := rds.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.DescribeDBInstances(nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	for _, i := range resp.DBInstances {
		if *i.DBInstanceIdentifier == id {
			return i
		}
	}

	return nil
}

func describeSubnetGroup(region string, id string) *rds.DBSubnetGroup {
	svc := rds.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.DescribeDBSubnetGroups(nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	for _, i := range resp.DBSubnetGroups {
		if *i.DBSubnetGroupName == id {
			return i
		}
	}

	return nil
}

func createSubnetGroup(region string, params *rds.CreateDBSubnetGroupInput, verbose bool) *rds.DBSubnetGroup {
	svc := rds.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.CreateDBSubnetGroup(params)

	if err != nil {
		log.Fatalf("Error creating DB Subnet Group '%s': %s",
			*params.DBSubnetGroupName,
			err)
	}

	// Wait for it.
	avail := false
	for i := 0; i < 10; i++ {
		target := describeSubnetGroup(region, *params.DBSubnetGroupName)
		if target != nil && *target.DBSubnetGroupName == *params.DBSubnetGroupName &&
			*target.SubnetGroupStatus == "Complete" {
			avail = true
			break
		}
		time.Sleep(time.Second * 10)
	}

	if !avail {
		log.Fatalf("Error creating DB Subnet Group '%s'", *params.DBSubnetGroupName)
	}

	return resp.DBSubnetGroup
}

func create(region string, params *rds.CreateDBInstanceInput, wait bool, verbose bool) *rds.DBInstance {
	svc := rds.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.CreateDBInstance(params)
	if err != nil {
		log.Fatalf("Error creating DB Instance '%s': %s",
			*params.DBInstanceIdentifier,
			err)
	}
	db := *resp.DBInstance
	id := *db.DBInstanceIdentifier

	if wait {
		// Wait for it.
		avail := false
		for i := 0; i < 100; i++ {
			target := describe(region, *db.DBInstanceIdentifier)
			if target != nil && *target.DBInstanceIdentifier == id &&
				*target.DBInstanceStatus == "available" {
				avail = true
				break
			}
			time.Sleep(time.Second * 10)
		}

		if !avail {
			log.Fatalf("Error creating DB Instance '%s'", id)
		}
	}

	// Describe it.
	target := describe(region, id)

	return target
}

func deleteSubnetGroup(region string, id string, verbose bool) {
	svc := rds.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &rds.DeleteDBSubnetGroupInput{
		DBSubnetGroupName: aws.String(id),
	}
	_, err := svc.DeleteDBSubnetGroup(params)

	if err != nil {
		log.Fatal(err.Error())
	}

	// TODO: Wait for it.
}

func delete(region string, params *rds.DeleteDBInstanceInput, verbose bool) {
	svc := rds.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.DeleteDBInstance(params)

	if err != nil {
		log.Fatal(err.Error())
	}

	// Wait for it.
	avail := true
	for i := 0; i < 100; i++ {
		target := describe(region, *params.DBInstanceIdentifier)
		if target == nil {
			avail = false
			break
		}
		time.Sleep(time.Second * 10)
	}

	if avail {
		log.Fatalf("Error deleting DB Instance '%s'", *params.DBInstanceIdentifier)
	}
}

func scan(rt *otto.Otto, region string) otto.Value {
	svc := rds.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.DescribeDBInstances(nil)
	if err != nil {
		panic(err)
	}

	dbs := []rds.DBInstance{}
	// shove instances into jsland
	for _, i := range resp.DBInstances {
		dbs = append(dbs, *i)
	}
	return mcore.Sanitize(rt, dbs)
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

		if b, err := awsObj.Get("rcs"); err != nil || b.IsUndefined() {
			o1, _ = rt.Object(`aws.rds = {}`)
			o2, _ = rt.Object(`aws.rds.subnetGroups = {}`)
		} else {
			o1 = b.Object()
			v, _ := o1.Get("subnetGroups")
			o2 = v.Object()
		}

		o1.Set("scan", func(region string) otto.Value {
			return scan(rt, region)
		})
		o1.Set("describe", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			f := mcore.Sanitizer(rt)
			return f(describe(region, id))
		})
		o1.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input rds.CreateDBInstanceInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for RDS create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall RDS create json: %s", err)
			}

			region := call.Argument(0).String()
			wait, err := call.Argument(2).ToBoolean()
			if err != nil {
				log.Fatalf("Invalid gateway arg to RDS create: %s", err)
			}
			verbose := mcore.IsVerbose(rt)

			f := mcore.Sanitizer(rt)
			return f(create(region, &input, wait, verbose))
		})
		o1.Set("delete", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input rds.DeleteDBInstanceInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for RDS delete input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall RDS delete json: %s", err)
			}

			region := call.Argument(0).String()
			verbose := mcore.IsVerbose(rt)

			delete(region, &input, verbose)
			return otto.Value{}
		})

		o2.Set("describe", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			f := mcore.Sanitizer(rt)
			return f(describeSubnetGroup(region, id))
		})
		o2.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input rds.CreateDBSubnetGroupInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for RDS create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall RDS create json: %s", err)
			}

			region := call.Argument(0).String()
			verbose := mcore.IsVerbose(rt)

			f := mcore.Sanitizer(rt)
			return f(createSubnetGroup(region, &input, verbose))
		})
		o2.Set("delete", func(call otto.FunctionCall) otto.Value {
			region := call.Argument(0).String()
			id := call.Argument(1).String()
			verbose := mcore.IsVerbose(rt)
			deleteSubnetGroup(region, id, verbose)
			return otto.Value{}
		})
	})
}
