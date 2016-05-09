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
// # Beanstalk
// 
// Beanstalk is a resource handler for dealing with AWS Elastic Beanstalk.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"beantalk"`
// 
// Usage:
// 
// `var beanstalk = require("beanstalk").init();`
// 
//  ## Example Resource
// 
// ```javascript
//
// var beanstalk = {
//      name: "beanstalk"
//      module: "beanstalk"
//      params: {
//          ensure: "present"
//          app: {
//             ApplicationName: appName
//             Description:     "Description"
//          }
//          version: {
//             ApplicationName:       appName
//             VersionLabel:          "VersionLabel"
//             AutoCreateApplication: false
//             Description:           "Description"
//             Process:               true
//             SourceBundle: {
//                 S3Bucket: sourceBucket
//                 S3Key:    "api/v0.1.1-1-g5b33818.zip"
//             }
//          }
//          config: {
//             ApplicationName: appName
//             TemplateName:    "ConfigurationTemplateName"
//             Description:     "Description"
//             SolutionStackName: "64bit Amazon Linux 2016.03 v2.1.0 running Go 1.4"
//             // See: https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/command-options-general.html
//             // OptionSettings: [
//             //          {
//             //              Namespace:    "aws:elasticbeanstalk:application"
//             //              OptionName:   "Application Healthcheck URL"
//             //              Value:        "/hc"
//             //          }
//             // ]
//          }
//          environment: {
//             ApplicationName: appName
//             CNAMEPrefix:     cname
//             Description:     "test"
//             EnvironmentName: "test"
//             Tags: [
//                 {
//                     Key:   "TagKey"
//                     Value: "TagValue"
//                 }
//             ]
//             SolutionStackName: "64bit Amazon Linux 2016.03 v2.1.0 running Go 1.4"
//             VersionLabel: "VersionLabel"
//          }
//          terminate: {
//             EnvironmentName:    "test"
//             ForceTerminate:     true
//             TerminateResources: true
//          }
//      }
// }
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
// If `"present"` and ...
// 
// ### `app`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticbeanstalk.html#type-CreateApplicationInput)
//
// Specifies the application to be created/deleted.
//
// ### `version`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticbeanstalk.html#type-CreateApplicationVersionInput)
//
// Specifies properties for an application version to be created/deleted.
//
// ### `environment`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticbeanstalk.html#type-CreateEnvironmentInput)
//
// Specifies properties for a beanstalk application environment to be created/deleted.
//
// ### `config`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticbeanstalk.html#type-CreateApplicationVersionInput)
//
// Specifies properties for a beanstalk application configuration to be created/deleted.
//
// ### `terminate`
//
// * Required: true
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticbeanstalk.html#type-TerminateEnvironmentInput)
//
// Specifies properties for termination of an environment.
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
        moduleName: "beanstalk"
        findInCatalog: function(catalog, resource) {
            if (resource.module != handler.moduleName) {
                return [null, false];
            }
            if (typeof(resource.params.on_find) === 'function') {
                result = resource.params.on_find(catalog, resource);
                if (!_.isObject(result)) {
                    return;
                }
                return result;
            }
            // Look in beanstalkApps, beanstalkVersions, beanstalkEnvironment and beanstalkConfigs
            var params = resource.params;
            var appName;

            var app;
            if (params.app) {
                appName = params.app.ApplicationName;
                app = _.find(catalog.beanstalkApps, function(a) { 
                    return (a.ApplicationName === appName);
                });
            }
            var version;
            if (params.version) {
                appName = params.version.ApplicationName;
                version = _.find(catalog.beanstalkVersions, function(v) { 
                    return ((v.ApplicationName === appName) &&
                            (v.VersionLabel === params.version.VersionLabel));
                });
            }
            var env;
            if (params.environment) {
                appName = params.environment.ApplicationName;
                env = _.find(catalog.beanstalkEnvironments, function(e) { 
                    return ((e.ApplicationName === appName) &&
                            (e.EnvironmentName === params.environment.EnvironmentName));
                });
            }
            var config;
            if (params.config) {
                appName = params.config.ApplicationName;
                config = _.find(catalog.beanstalkConfigs, function(c) { 
                    console.log(JSON.stringify(c, null, 2));
                    console.log("c", c, ((c.ApplicationName === appName) &&
                                         (c.TemplateName === params.config.TemplateName)),
                                c.ApplicationName,
                                c.TemplateName);
                    return ((c.ApplicationName === appName) &&
                            (c.TemplateName === params.config.TemplateName));
                });
                console.log("config", config);
            }
            var found;
            if (app || version || env || config) {
                found = {
                    app: app
                    version: version
                    config: config
                    env: env
                };
            }
            console.log("config", config);
            return found;
        }
        handle: function(catalog, resources, resource) {
            if (resource.module != handler.moduleName) {
                return [null, false];
            }

            var params = resource.params;
            var ensure = params.ensure;
            var app = params.app;
            var env = params.environment;
            var config = params.config;
            var version = params.version;
            var found = resource._target;

            // Sanity
            if (!(app || env || config || version)) {
                console.log("Invalid beanstalk params")
                os.exit(3);
            }

            switch(ensure) {
            case "absent":
                if (app) {
                    if (found && found.app) {
                        if (mithras.verbose) {
                            log(sprintf("Deleting app '%s'", app.ApplicationName));
                        }
                        aws.beanstalk.apps.delete(params.region, app.ApplicationName);
                        catalog.beanstalkApps = 
                            _.reject(catalog.beanstalkApps,
                                     function(s) { 
                                         return (s.ApplicationName == 
                                                 app.ApplicationName);
                                     });
                    } else {
                        if (mithras.verbose) {
                            log(sprintf("No action taken for app."));
                        }
                    }
                } 
                if (version) {
                    if (found && found.version) {
                        if (mithras.verbose) {
                            log(sprintf("Deleting app version '%s'", 
                                        version.VersionLabel));
                        }
                        aws.beanstalk.versions.delete(params.region, 
                                                      version.ApplicationName, 
                                                      version.VersionLabel);
                        catalog.beanstalkVersions = 
                            _.reject(catalog.beanstalkVersion,
                                     function(v) { 
                                         return ((v.ApplicationName == 
                                                  version.ApplicationName) &&
                                                 (v.VersionLabel == 
                                                  version.VersionLabel))
                                     });
                    } else {
                        if (mithras.verbose) {
                            log(sprintf("No action taken for app version."));
                        }
                    }
                }
                if (config) {
                    if (found && found.config) {
                        if (mithras.verbose) {
                            log(sprintf("Deleting app config '%s'", 
                                        config.TemplateName));
                        }
                        aws.beanstalk.configs.delete(params.region, 
                                                     config.ApplicationName, 
                                                     config.ConfigLabel);
                        catalog.beanstalkConfigs = 
                            _.reject(catalog.beanstalkConfig,
                                     function(c) { 
                                         return ((c.ApplicationName == 
                                                  config.ApplicationName) &&
                                                 (c.TemplateName == 
                                                  config.TemplateName))
                                     });
                    } else {
                        if (mithras.verbose) {
                            log(sprintf("No action taken for app config."));
                        }
                    }
                }
                if (env) {
                    if (found && found.env) {
                        if (mithras.verbose) {
                            log(sprintf("Deleting app env '%s'", env.EnvironmentName));
                        }
                        aws.beanstalk.environments.terminate(params.region, 
                                                             env.EnvironmentName);
                        catalog.beanstalkEnvironments = 
                            _.reject(catalog.beanstalkEnvironments,
                                     function(e) { 
                                         return ((e.ApplicationName == 
                                                  env.ApplicationName) &&
                                                 (e.TemplateName == env.TemplateName))
                                     });
                    } else {
                        if (mithras.verbose) {
                            log(sprintf("No action taken for app env."));
                        }
                    }
                }
               break;
            case "present":
                if (!found) {
                    found = {};
                }
                if (app) {
                    if (!found.app) {
                        if (mithras.verbose) {
                            log(sprintf("Creating app '%s'", app.ApplicationName));
                        }
                        found.app = aws.beanstalk.apps.create(params.region, app);
                        catalog.beanstalkApps.push(found.app);
                    } else {
                        if (mithras.verbose) {
                            log(sprintf("No action taken for app."));
                        }
                    }
                } 
                if (version) {
                    if (!found.version) {
                        if (mithras.verbose) {
                            log(sprintf("Creating app version '%s'", 
                                        version.VersionLabel));
                        }
                        found.version = aws.beanstalk.versions.create(params.region, 
                                                                      version);
                        catalog.beanstalkVersions.push(found.version);
                    } else {
                        if (mithras.verbose) {
                            log(sprintf("No action taken for app version."));
                        }
                    }
                }
                if (config) {
                    if (!found.config) {
                        if (mithras.verbose) {
                            log(sprintf("Creating app config '%s'", 
                                        config.TemplateName));
                        }
                        found.config = aws.beanstalk.configs.create(params.region, 
                                                                    config);
                        catalog.beanstalkConfigs.push(found.configs);
                    } else {
                        if (mithras.verbose) {
                            log(sprintf("No action taken for app config."));
                        }
                    }
                }
                if (env) {
                    if (!found.env) {
                        if (mithras.verbose) {
                            log(sprintf("Creating app env '%s'", env.EnvironmentName));
                        }
                        found.env = aws.beanstalk.environments.create(params.region, 
                                                                      env);
                        catalog.beanstalkEnvironments.push(found.env);
                    } else {
                        if (mithras.verbose) {
                            log(sprintf("No action taken for app env."));
                        }
                    }
                }
                return [found, true];
            }
            return [null, true];
        }
        preflight: function(catalog, resources, resource) {
            if (resource.module != handler.moduleName) {
                return [null, false];
            }
            var s = handler.findInCatalog(catalog, resource);
            return [s, true];
        }
    };
    
    handler.init = function () {
        mithras.modules.preflight.register(handler.moduleName, handler.preflight);
        mithras.modules.handlers.register(handler.moduleName, handler.handle);
        return handler;
    };
    
    return handler;
});
