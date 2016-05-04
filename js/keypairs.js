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
// # keypairs
// 
// Keypairs is resource handler for manipulating AWS SSH keypairs.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"keypairs"`
// 
// Usage:
// 
// `var keypairs = require("keypairs").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rKey = {
//      name: "key"
//      module: "keypairs"
//      skip: (ensure === 'absent') // Don't delete keys
//      params: {
//          region: "us-east-1"
//          ensure: "present"
//          key: {
//              KeyName: "my-fancy-key"
//          }
//          savePath: os.expandEnv("$HOME/.ssh/" + keyName + ".pem")
//      }
// };
// ```
// 
// ## Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "absent", "present"
//
// If `"present"`, the handler will ensure the keypair exists, and it
// not, it will be created.  If `"absent"`, the keypair is removed.
// 
// ### `key`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-CreateKeyPairInput)
//
// Specifies parameters for keypair creation.
//
// ### `savepath`
//
// * Required: true
// * Allowed Values: A valid path for saving the pemfile when a keypair is created
//
// When the handler creates a new keypair, the contents of the key are saved to this path.
//
// ### `on_find`
//
// * Required: true
// * Allowed Values: A function taking two parameters: `catalog` and `resource`
//
// If defined in the resource's `params` object, the `on_find`
// function provides a way for a matching resource to be identified
// using a user-defined way.  The function is called with the current
// `catalog`, as well as the `resource` object itself.  The function
// can look through the catalog, find a matching object using whatever
// logic you want, and return it.  If the function returns `undefined`
// or a n empty Javascript array, (`[]`), the function is indicating
// that no matching resource was found in the `catalog`.
// 
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
        module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
        moduleNames: ["keypairs"]
        findInCatalog: function(catalog, resource) {
            if (typeof(resource.params.on_find) === 'function') {
		result = resource.params.on_find(catalog, resource);
		if (!result || 
		    (Array.isArray(result) && result.length == 0)) {
		    return;
		}
		return result;
	    }
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
                os.exit(3);
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
                var err = fs.write(params.savePath, raw, 0400);
		if (err) {
                    log(sprintf("Error writing key '%s' to '%s': %s", 
				params.key.KeyName, 
				params.savePath,
				err));
		    os.exit(3);
		}

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
        preflight: function(catalog, resources, resource) {
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
