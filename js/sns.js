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
// # SNS
// 
// SNS is a resource handler for dealing with AWS SNS resources.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"sns"`
// 
// Usage:
// 
// `var sns = require("sns").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rTopic = {
//     name: "snsTopic"
//     module: "sns"
//     params: {
//         region: defaultRegion
//         ensure: ensure
//         topic: {
//             Name:  "my-topic"
//         }
//     }
// };
// var rSub = {
//     name: "snsSub"
//     module: "sns"
//     dependsOn: [rTopic.name]
//     params: {
//         region: defaultRegion
//         ensure: ensure
//         sub: {
//           Protocol: "..."
//           TopicArn: "..."
//           Endpoint: "..."
//         }
//     }
// };
// var rPub = {
//     name: "snsPub"
//     module: "sns"
//     dependsOn: [rTopic.name]
//     params: {
//         region: defaultRegion
//         ensure: ensure
//         sub: {...}
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
// If `"present"` and the sns topic `params.topic.Name` does not
// exist, it is created.  If `"absent"`, and it exists, it is removed.
// 
// If `"present"` and the sns subscription referencing
// `params.topic.Name` does not exist, it is created.  If `"absent"`,
// and it exists, it is removed.
// 
// If `"present"` and the the `params.pub` property is set, a message
// is published to the topic.  This is NOT an idempotent operation.
// 
// ### `topic`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/sns.html#type-CreateTopicInput)
//
// Parameters for topic creation.
//
// ### `sub`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/sns.html#type-SubscribeInput)
//
// Parameters for subscription creation.
//
// ### `pub`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/sns.html#type-PublishInput)
//
// Parameters for publishing a message to a topic.
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
        moduleName: "sns"
        findInCatalog: function(catalog, resource) {
            if (typeof(resource.params.on_find) === 'function') {
		result = resource.params.on_find(catalog, resource);
		if (!result || 
		    (Array.isArray(result) && result.length == 0)) {
		    return;
		}
		return result;
	    }
	    var topic = null;
	    var sub = null;
	    if (resource.params.topic) {
		var re = new RegExp(":"+resource.params.topic.Name+"$");
		topic = _.find(catalog.topics, function(t) { 
                    return re.exec(t);
		});
	    }
	    if (resource.params.sub) {
		sub = _.find(catalog.subs, function(s) { 
		    var match = (s.TopicArn === resource.params.sub.TopicArn &&
				 s.Protocol === resource.params.sub.Protocol &&
				 s.Endpoint === resource.params.sub.Endpoint)
                    return match;
		});
	    }
	    if (topic || sub) {
		return {topic: topic, sub: sub};
	    }
	    return;
        }
        handle: function(catalog, resources, resource) {
            if (resource.module != handler.moduleName) {
                return [null, false];
            }

            // Sanity
            if (!resource.params.topic && !resource.params.sub) {
                console.log("Invalid sns params")
                os.exit(3);
            }

            var ensure = resource.params.ensure;
            var params = resource.params;
	    var topic = resource._target ? resource._target.topic : null;
	    var sub = resource._target ? resource._target.sub : null;

            switch(ensure) {
            case "absent":
		if (topic && params.topic) {
                    if (mithras.verbose) {
			log(sprintf("Deleting topic '%s'", topic));
                    }
                    aws.sns.topics.delete(params.region, topic);
                    catalog.topics = _.reject(catalog.topics,
                                              function(t) { 
						  return t === topic;
                                              });
		}
		if (sub && params.sub) {
                    if (mithras.verbose) {
			log(sprintf("Deleting sub '%s'", sub.SubscriptionArn));
                    }
                    aws.sns.subs.delete(params.region, sub.SubscriptionArn);
                    catalog.subs = _.reject(catalog.subs,
                                            function(s) { 
						return s.SubscriptionArn === sub.SubscriptionArn;
                                            });
		}
		if (!sub && !topic) {
                    if (mithras.verbose) {
			log("No action taken.");
                    }
		}
                break;
            case "present":
		if (params.topic) {
		    if (!topic) {
			if (mithras.verbose) {
			    log(sprintf("Creating topic '%s'", params.topic.Name));
			}
			topic = aws.sns.topics.create(params.region, params.topic);
			catalog.topics.push(topic);
		    } else {
			log(sprintf("Topic '%s' found, no action taken.", params.topic.Name));
		    }
		}
		if (params.pub) {
		    if (mithras.verbose) {
			log(sprintf("Publishing message"))
		    }
		    aws.sns.topics.publish(params.region, params.pub);
		}
		if (params.sub) {
		    if (!sub) {
			if (mithras.verbose) {
			    log(sprintf("Creating sub for '%s'", params.sub.TopicArn));
			}
			sub = aws.sns.subs.create(params.region, params.sub);
			catalog.subs.push(sub);
		    } else {
			log(sprintf("Subscription found, no action taken."));
		    }
		}
                // return it
                return [{topic: topic, sub: sub}, true];
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
