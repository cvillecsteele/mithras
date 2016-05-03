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
// # SQS
// 
// SQS is a resource handler for dealing with AWS SQS resources.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"sqs"`
// 
// Usage:
// 
// `var sqs = require("sqs").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rQueue = {
//     name: "sqsQueue"
//     module: "sqs"
//     params: {
//        region: defaultRegion
//        ensure: ensure
//        queue: {
//          QueueName: "myqueue"
//          Attributes: [
//            "Key": "value"
//          ]
//        }
//     }
// };
// var rPub = {
//     name: "sqsPub"
//     module: "sqs"
//     dependsOn: [rTopic.name]
//     params: {
//         ensure: ensure
//         message: {
//            MessageBody:  "body"
//            QueueUrl:     "url"
//            DelaySeconds: 1
//            MessageAttributes: {
//              "Key": {
//                DataType: "type"
//                BinaryListValues: [
//                  "PAYLOAD"
//                ]
//                BinaryValue: "PAYLOAD"
//                StringListValues: [ "String" ]
//                StringValue: "String"
//              }
//            }
//         }
//     }
// };
// ```
// 
// ## Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "present" or "absent"
//
// If `"present"` and the sqs queue `params.queue.QueueName` does not
// exist, it is created.  If `"absent"`, and it exists, it is removed.
// 
// If `"present"` and the the `params.message` property is set, a message
// is published to the queue.  This is NOT an idempotent operation.
// 
// ### `queue`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/sqs.html#type-CreateQueueInput)
//
// Parameters for queue creation.
//
// ### `message`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/sqs.html#type-SendMessageInput)
//
// Parameters for publishing a message to a queue.
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
        moduleName: "sqs"
        findInCatalog: function(catalog, resource) {
            if (typeof(resource.params.on_find) === 'function') {
		result = resource.params.on_find(catalog, resource);
		if (!result || 
		    (Array.isArray(result) && result.length == 0)) {
		    return;
		}
		return result;
	    }
	    var queue = null;
	    if (resource.params.queue) {
		var re = new RegExp("/"+resource.params.queue.QueueName+"$");
		queue = _.find(catalog.queues, function(t) { 
                    return re.exec(t);
		});
	    }
	    return queue;
        }
        handle: function(catalog, resources, resource) {
            if (resource.module != handler.moduleName) {
                return [null, false];
            }

            // Sanity
            if (!resource.params.queue && !resource.params.queue) {
                console.log("Invalid sqs params")
                os.exit(3);
            }

            var ensure = resource.params.ensure;
            var params = resource.params;
	    var queue = resource._target;

            switch(ensure) {
            case "absent":
		if (queue && params.queue) {
                    if (mithras.verbose) {
			log(sprintf("Deleting queue '%s'", params.queue.QueueName));
                    }
                    aws.sqs.delete(params.region, queue);
                    catalog.queues = _.reject(catalog.queues,
                                              function(t) { 
						  return t === queue;
                                              });
		} else if (!queue) {
                    if (mithras.verbose) {
			log("No action taken.");
                    }
		}
                break;
            case "present":
		if (params.queue) {
		    if (!queue) {
			if (mithras.verbose) {
			    log(sprintf("Creating queue '%s'", params.queue.QueueName));
			}
			queue = aws.sqs.create(params.region, params.queue);
			catalog.queues.push(queue);
		    } else {
			log(sprintf("Queue '%s' found, no action taken.", 
				    params.queue.QueueName));
		    }
		}
		if (params.message) {
		    if (mithras.verbose) {
			log(sprintf("Sending message"))
		    }
		    aws.sqs.messages.send(params.region, params.message);
		}
                // return it
                return [queue, true];
            }
            return [null, true];
        }
        preflight: function(catalog, resources, resource) {
            if (resource.module != handler.moduleName) {
                return [null, false];
            }
            var params = resource.params;
            var s = handler.findInCatalog(catalog, resource);
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
