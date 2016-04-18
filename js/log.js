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
// # log
// 
// Log is resource handler for writing messages to the mithras log.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"log"`
// 
// Usage:
// 
// `var log = require("log").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rLog = {
//      name: "log"
//      module: "log"
//      params: {
//          messsage:  "Hello, world!"
//      }
// };
// ```
// 
// ## Parameter Properties
// 
// ### `message`
//
// * Required: true
// * Allowed Values: any string
//
// The value of this property will be written to the Mithras log.
// 
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var handler = {
	moduleNames: ["log"]
	handle: function(catalog, resources, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
		
	    var p = resource.params;

	    log(p.message);

	    return [null, true];
	}
    };
		   

    handler.init = function () {
	mithras.modules.handlers.register(handler.moduleNames[0], handler.handle);
	return handler;
    };
    
    return handler;
});
