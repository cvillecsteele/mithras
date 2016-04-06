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
	findInCatalog: function(catalog, id) {
	    return _.find(catalog.dbs, function(inst) { 
		return inst.DBInstanceIdentifier === id;
	    });
	}
	preflight: function(catalog, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
	    var db = handler.findInCatalog(catalog, 
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
