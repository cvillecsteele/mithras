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
// # Elasticache
// 
// Elasticache is resource handler for managing AWS caches.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"elasticache"`
// 
// Usage:
// 
// `var elasticache = require("elasticache").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rds = {
//   name: "rdsA"
//   module: "rds"
//   dependsOn: [otherResource.name]
//   params: {
//     ensure: ensure
//     region: defaultRegion
//     wait: true
//     subnetGroup: {
//         DBSubnetGroupDescription: "test subnet group"
//         DBSubnetGroupName: "test-subnet-group"
//         SubnetIds: [
//              "subnet-123",
//              "subnet-456"
//         ]
//         Tags: [
//             {
//                 Key:   "Foo"
//                 Value: "Bar"
//             }
//         ]
//     }
//     db: {
//         DBInstanceClass:         "db.m1.small"
//         DBInstanceIdentifier:    "test-rds"
//         Engine:                  "mysql"
//         AllocatedStorage:        10
//         AutoMinorVersionUpgrade: true
//         AvailabilityZone:        defaultZone
//         MasterUserPassword:      "test123456789"
//         MasterUsername:          "test"
//         DBSubnetGroupName:       "test-subnet-group"
//         DBName:                  "test"
//         PubliclyAccessible:      false
//         Tags: [
//             {
//                 Key:   "foo"
//                 Value: "bar"
//             },
//         ]
//     }
//     delete: {
//         DBInstanceIdentifier:      "db-abcd"
//         FinalDBSnapshotIdentifier: "byebye" + Date.now()
//         SkipFinalSnapshot:         true
//     }
//   }
// };
// 
// ```
// 
// ## Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "present" or "absent"
//
// If `"present"`, the db specified by `db` will be created, and
// if `"absent"`, it will be removed using the `delete` property.
// 
// ### `region`
//
// * Required: true
// * Allowed Values: string, any valid AWS region; eg "us-east-1"
//
// The region for calls to the AWS API.
// 
// ### `wait`
//
// * Required: false
// * Allowed Values: true or false
//
// If `true`, delay execution until the db has been created in AWS.
// 
// ### `subnetGroup`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/rds.html#type-CreateDbSubnetGroupInput)
//
// If set, a subnet group will be created for your db.
// 
// ### `db`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/rds.html#type-CreateDBClusterInput)
//
// Parameters for resource creation.
// 
// ### `delete`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/rds.html#type-DeleteDBClusterInput)
//
// Parameters for deletion.
// 
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var become = require("become").become;
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleNames: ["rds"]
	handle: function(catalog, resources, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }

	    var p = resource.params;
	    var ensure = p.ensure;
	    var db = resource._target;
		    
	    // Sanity
	    if (!p || !p.db || !p.delete) {
		console.log(sprintf("Invalid rds resource '%s'", resource.name));
		exit(1);
	    }

	    switch(ensure) {
	    case "absent":

		// remove subnet group
		if (p.subnetGroup &&
		    aws.rds.subnetGroups.describe(p.region, 
						  p.subnetGroup.DBSubnetGroupName)) {
		    if (mithras.verbose) {
			log(sprintf("Deleting rds subnet group '%s'", 
				    p.subnetGroup.DBSubnetGroupName));
		    }
		    aws.rds.subnetGroups.delete(p.region, p.subnetGroup.DBSubnetGroupName);
		}
		
		if (!db) {
		    break;
		}
		var id = resource._target.DBInstanceIdentifier;
		
		// delete
		if (mithras.verbose) {
		    log(sprintf("Deleting rds instance '%s'", id));
		}
		console.log(JSON.stringify(p.delete, null, 2));
		aws.rds.delete(p.region, p.delete);
		
		// remove from catalog
		catalog.dbs = _.reject(catalog.dbs, 
				       function(x) { 
					   return x.DBInstanceIdentifier == id;
				       });
		break;
	    case "present":
		if (db) {
		    break;
		}
		var id = p.db.DBInstanceIdentifier;

		// create subnet group
		if (p.subnetGroup && 
		    (!aws.rds.subnetGroups.describe(p.region, 
						    p.subnetGroup.DBSubnetGroupName))) {
		    if (mithras.verbose) {
			log(sprintf("Creating rds subnet group '%s'", 
				    p.subnetGroup.DBSubnetGroupName));
		    }
		    aws.rds.subnetGroups.create(p.region, p.subnetGroup);
		}
		
		// create instance
		if (mithras.verbose) {
		    log(sprintf("Creating rds instance '%s' (WAIT FOR IT...)", id));
		}
		aws.rds.create(p.region, p.db, p.wait);

		// re-describe it
		db = aws.rds.describe(p.region, id);

		// add to catalog
		catalog.dbs.push(db);
		break;
	    }
	    return [null, true];
	}
	findInCatalog: function(catalog, resource, id) {
            if (typeof(resource.params.on_find) === 'function') {
		result = resource.params.on_find(catalog, resource);
		if (!result || 
		    (Array.isArray(result) && result.length == 0)) {
		    return;
		}
		return result;
	    }
	    return _.find(catalog.dbs, function(inst) { 
		return inst.DBInstanceIdentifier === id;
	    });
	}
	preflight: function(catalog, resources, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
	    var db = handler.findInCatalog(catalog, 
					   resource,
					   resource.params.db.DBInstanceIdentifier);
	    if (db) {
		return [db, true];
	    }
	    return [null, true];
	}
    };
		   

    handler.init = function () {
	mithras.modules.preflight.register(handler.moduleNames[0], handler.preflight);
	mithras.modules.handlers.register(handler.moduleNames[0], handler.handle);
	return handler;
    };
    
    return handler;
});
