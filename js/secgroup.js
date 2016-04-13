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
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleName: "secgroup"
	findInCatalog: function(catalog, vpcId, groupName) {
	    return _.find(catalog.securityGroups, function(s) { 
		return (s.GroupName === groupName && s.VpcId == vpcId);
	    });
	}
	handle: function(catalog, resources, resource) {
	    if (resource.module != handler.moduleName) {
		return [null, false];
	    }

	    // Sanity
	    if (!resource.params.secgroup) {
		console.log("Invalid secgroup params")
		os.exit(1);
	    }

	    var ensure = resource.params.ensure;
	    var params = resource.params;
	    var sg = resource._target;

	    switch(ensure) {
	    case "absent":
		if (!sg) {
		    if (mithras.verbose) {
			log(sprintf("Security group not found, no action taken."));
		    }
		    break;
		}
		aws.securityGroups.delete(params.region, sg.GroupId);
		catalog.securityGroups = _.reject(catalog.securityGroups,
						  function(s) { 
						      return s.GroupId == sg.GroupId;
						  });
		break;
	    case "present":
		if (sg) {
		    if (mithras.verbose) {
			log(sprintf("Security group found, no action taken."));
		    }
		    break;
		}
		// create sg
		if (mithras.verbose) {
		    log(sprintf("Creating security group '%s'", 
				params.secgroup.GroupName));
		}
		var sg = aws.securityGroups.create(params.region, params.secgroup);

		// do authorizations
		if (params.ingress) {
		    params.ingress.GroupId = sg.GroupId;
		    if (mithras.verbose) {
			log(sprintf("Creating security group ingress authorization"));
		    }
		    aws.securityGroups.authorizeIngress(params.region, params.ingress);
		}
		if (params.egress) {
		    params.egress.GroupId = sg.GroupId;
		    if (mithras.verbose) {
			log(sprintf("Creating security group egress authorization"));
		    }
		    aws.securityGroups.authorizeEgress(params.region, params.egress);
		}
		
		// Tag it
		if (params.tags) {
		    if (mithras.verbose) {
			log(sprintf("Tagging security group."));
		    }
		    aws.tags.create(params.region, sg.GroupId, params.tags);
		}
		
		// Reload it to get tags
		var sg = aws.securityGroups.describe(params.region, sg.GroupId);
		
		// add to catalog
		catalog.securityGroups.push(sg);
		
		// return it
		return [sg, true];
	    }
	    return [null, true];
	}
	preflight: function(catalog, resource) {
	    if (resource.module != handler.moduleName) {
		return [null, false];
	    }
	    var params = resource.params;
	    var sg = params.secgroup;
	    var s = handler.findInCatalog(catalog, sg.VpcId, sg.GroupName);
	    if (s) {
		return [s, true];
	    }
	    return [null, true];
	}
    };
    
    handler.init = function () {
	mithras.modules.preflight.register(handler.moduleName, handler.preflight);
	mithras.modules.handlers.register(handler.moduleName, handler.handle);
	return handler;
    };
    
    return handler;
});
