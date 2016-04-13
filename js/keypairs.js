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
	moduleNames: ["keypairs"]
	findInCatalog: function(catalog, resource) {
	    return _.find(catalog.keypairs, function(key) { 
		return key.KeyName === resource.params.key.KeyName;
	    });
	}
	handle: function(catalog, resources, resource) {
	    if (!_.find(handler.moduleNames, function(m) { return resource.module === m; })) {
		return [null, false];
	    }
		
	    var ensure = resource.params.ensure;
	    var params = resource.params;
	    
	    if (!params.ensure || !params.key || !params.savePath || !params.region) {
		console.log("Invalid key params")
		os.exit(1);
	    }
	    var found = resource._target;

	    switch(ensure) {
	    case "absent":
		if (!found) {
		    if (mithras.verbose) {
			log(sprintf("Key '%s' not found, no action taken.", params.key.KeyName));
			break;
		    }
		}
		if (mithras.verbose) {
		    log(sprintf("Deleting key '%s'", params.key.KeyName));
		}
		aws.keypairs.delete(params.region, params.key.KeyName);
		catalog.keypairs = 
		    _.reject(catalog.keypairs,
			     function(k) { 
				 return k.KeyName == params.key.KeyName;
			     });
		break;
	    case "present":
		if (found) {
		    if (mithras.verbose) {
			log(sprintf("Key '%s' found, no action taken.", params.key.KeyName));
			break;
		    }
		}

		if (mithras.verbose) {
		    log(sprintf("Creating key '%s'", params.key.KeyName));
		}
		var raw = aws.keypairs.create(params.region, params.key.KeyName);
		if (mithras.verbose) {
		    log(sprintf("Writing key '%s' to '%s'", params.key.KeyName, params.savePath));
		}
		fs.write(params.savePath, raw, 0400);

		var key = aws.keypairs.describe(params.region, params.key.KeyName);
		do {
		    time.sleep(10);
		    key = aws.keypairs.describe(params.region, params.key.KeyName);
		} while (!key);
		catalog.keypairs.push(key);

		// return it
		return [handler.findInCatalog(catalog, resource), true];
		break;
	    }
	    return [null, true];
	}
	preflight: function(catalog, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
	    var found = handler.findInCatalog(catalog, resource);
	    if (found) {
		return [found, true];
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
