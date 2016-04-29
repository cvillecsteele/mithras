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
// # CORE FUNCTIONS: ELASTICACHE
//

package elasticache

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.elasticache.create](#create)
// > * [aws.elasticache.delete](#delete)
// > * [aws.elasticache.describe](#describe)
// > * [aws.elasticache.scan](#scan)
// > * [aws.elasticache.subnetGroups.create](#sgcreate)
// > * [aws.elasticache.subnetGroups.delete](#sgdelete)
// > * [aws.elasticache.subnetGroups.describe](#sgdescribe)
//
// This API allows exposes functions to manage AWS cache clusters and
// their associated subnet groups.
//
// ## AWS.ELASTICACHE.CREATE
// <a name="create"></a>
// `aws.elasticache.create(region, config);`
//
// Create a cache cluster.
//
// Example:
//
// ```
//
//  var cache = aws.elasticache.create("us-east-1", {
//     CacheClusterId:          "test-redis"
//     AutoMinorVersionUpgrade: true
//     CacheNodeType:           "cache.t2.small"
//     CacheSubnetGroupName:    "redis-subnet-group"
//     Engine:                  "redis"
//     NumCacheNodes:           1
//     SecurityGroupIds:        []
//     Tags: [
//         {
//            Key:   "Name"
//            Value: "test-cluster"
//         },
//     ]
// });
//
// ```
//
// ## AWS.ELASTICACHE.DELETE
// <a name="delete"></a>
// `aws.elasticache.delete(region, cache_id);`
//
// Delete a cache cluster.
//
// Example:
//
// ```
//
//  aws.elasticache.delete("us-east-1", "test-redis");
//
// ```
//
// ## AWS.ELASTICACHE.DESCRIBE
// <a name="describe"></a>
// `aws.elasticache.describe(region, cache_id);`
//
// Get information about a cache cluster.
//
// Example:
//
// ```
//
//  var cache = aws.elasticache.describe("us-east-1", "test-redis");
//
// ```
//
// ## AWS.ELASTICACHE.SCAN
// <a name="scan"></a>
// `aws.elasticache.scan(region, cache_id);`
//
// Get information about all cache clusters.
//
// Example:
//
// ```
//
//  var caches = aws.elasticache.scan("us-east-1");
//
// ```
//
// ## AWS.ELASTICACHE.SUBNETGROUPS.CREATE
// <a name="sgcreate"></a>
// `aws.elasticache.subnetGroups.create(region, config);`
//
// Create an elasticache subnet group.
//
// Example:
//
// ```
//
//  var cache = aws.elasticache.create("us-east-1", {
//     CacheClusterId:          "test-redis"
//     AutoMinorVersionUpgrade: true
//     CacheNodeType:           "cache.t2.small"
//     CacheSubnetGroupName:    "redis-subnet-group"
//     Engine:                  "redis"
//     NumCacheNodes:           1
//     SecurityGroupIds:        []
//     Tags: [
//         {
//            Key:   "Name"
//            Value: "test-cluster"
//         },
//     ]
// });
//
// ```
//
// ## AWS.ELASTICACHE.SUBNETGROUPS.DELETE
// <a name="sgdelete"></a>
// `aws.elasticache.subnetGroups.delete(region, cache_id);`
//
// Delete a cache subnet group.
//
// Example:
//
// ```
//
//  aws.elasticache.subnetGroups.delete("us-east-1", "redis-subnet-group");
//
// ```
//
// ## AWS.ELASTICACHE.SUBNETGROUPS.DESCRIBE
// <a name="sgdescribe"></a>
// `aws.elasticache.subnetGroups.describe(region, cache_id);`
//
// Get information about a cache cluster.
//
// Example:
//
// ```
//
//  var group = aws.elasticache.subnetGroups.describe("us-east-1", "redis-subnet-group");
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
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "elasticache"

func describe(region string, id string) *elasticache.CacheCluster {
	svc := elasticache.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	in := elasticache.DescribeCacheClustersInput{
		CacheClusterId:    aws.String(id),
		ShowCacheNodeInfo: aws.Bool(true),
	}

	resp, err := svc.DescribeCacheClusters(&in)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if "CacheClusterNotFound" == awsErr.Code() {
				return nil
			}
		}
		log.Fatalf("Error describing cache: %s", err)
	}

	if resp != nil && len(resp.CacheClusters) > 0 {
		return resp.CacheClusters[0]
	}
	return nil
}

func describeSubnetGroup(region string, id string) *elasticache.CacheSubnetGroup {
	svc := elasticache.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &elasticache.DescribeCacheSubnetGroupsInput{
		CacheSubnetGroupName: aws.String(id),
	}
	resp, err := svc.DescribeCacheSubnetGroups(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if "CacheSubnetGroupNotFoundFault" == awsErr.Code() {
				return nil
			}
		}
		log.Fatalf("Error describing subnet group: %s", err)
	}

	if len(resp.CacheSubnetGroups) > 0 {
		return resp.CacheSubnetGroups[0]
	}
	return nil
}

func createSubnetGroup(region string, params *elasticache.CreateCacheSubnetGroupInput, verbose bool) *elasticache.CacheSubnetGroup {
	svc := elasticache.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.CreateCacheSubnetGroup(params)

	if err != nil {
		log.Fatalf("Error creating Cache Subnet Group '%s': %s",
			*params.CacheSubnetGroupName,
			err)
	}

	// Wait for it.
	avail := false
	for i := 0; i < 10; i++ {
		target := describeSubnetGroup(region, *params.CacheSubnetGroupName)
		if target != nil && *target.CacheSubnetGroupName == *params.CacheSubnetGroupName {
			avail = true
			break
		}
		time.Sleep(time.Second * 10)
	}

	if !avail {
		log.Fatalf("Error creating Cache Subnet Group '%s'", *params.CacheSubnetGroupName)
	}

	return resp.CacheSubnetGroup
}

func create(region string, params *elasticache.CreateCacheClusterInput, wait bool, verbose bool) *elasticache.CacheCluster {
	svc := elasticache.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.CreateCacheCluster(params)
	if err != nil {
		log.Fatalf("Error creating Cache '%s': %s",
			*params.CacheClusterId,
			err)
	}
	cache := *resp.CacheCluster
	id := *cache.CacheClusterId

	if wait {
		// Wait for it.
		avail := false
		for i := 0; i < 100; i++ {
			target := describe(region, id)
			if target != nil && *target.CacheClusterId == id &&
				*target.CacheClusterStatus == "available" {
				avail = true
				break
			}
			time.Sleep(time.Second * 10)
		}

		if !avail {
			log.Fatalf("Error creating Cache Instance '%s'", id)
		}
	}

	// Describe it.
	target := describe(region, id)

	return target
}

func deleteSubnetGroup(region string, id string, verbose bool) {
	svc := elasticache.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &elasticache.DeleteCacheSubnetGroupInput{
		CacheSubnetGroupName: aws.String(id),
	}
	_, err := svc.DeleteCacheSubnetGroup(params)

	if err != nil {
		log.Fatalf("Error deleting subnet group: %s", err)
	}

	// TODO: Wait for it.
}

func delete(region string, params *elasticache.DeleteCacheClusterInput, verbose bool) {
	svc := elasticache.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.DeleteCacheCluster(params)

	if err != nil {
		log.Fatal("Error deleting cache: %s", err)
	}

	// Wait for it.
	avail := true
	for i := 0; i < 100; i++ {
		target := describe(region, *params.CacheClusterId)
		if target == nil {
			avail = false
			break
		} else if *target.CacheClusterStatus == "deleted" {
			avail = false
			break
		}
		time.Sleep(time.Second * 10)
	}

	if avail {
		log.Fatalf("Error deleting cache '%s'", *params.CacheClusterId)
	}
}

func scan(rt *otto.Otto, region string) otto.Value {
	svc := elasticache.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.DescribeCacheClusters(nil)
	if err != nil {
		panic(err)
	}

	caches := []elasticache.CacheCluster{}
	// shove instances into jsland
	for _, i := range resp.CacheClusters {
		caches = append(caches, *i)
	}
	return mcore.Sanitize(rt, caches)
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

		if b, err := awsObj.Get("elasticache"); err != nil || b.IsUndefined() {
			o1, _ = rt.Object(`aws.elasticache = {}`)
			o2, _ = rt.Object(`aws.elasticache.subnetGroups = {}`)
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
			var input elasticache.CreateCacheClusterInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for elasticache create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall elasticache create json: %s", err)
			}

			region := call.Argument(0).String()
			wait, err := call.Argument(2).ToBoolean()
			if err != nil {
				log.Fatalf("Invalid gateway arg to elasticache create: %s", err)
			}
			verbose := mcore.IsVerbose(rt)

			f := mcore.Sanitizer(rt)
			return f(create(region, &input, wait, verbose))
		})
		o1.Set("delete", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input elasticache.DeleteCacheClusterInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for elasticache delete input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall elasticache delete json: %s", err)
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
			var input elasticache.CreateCacheSubnetGroupInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for elasticache subnet group create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall elasticache subnet group create json: %s", err)
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
