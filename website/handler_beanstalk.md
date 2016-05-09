 
 
 # Beanstalk
 
 Beanstalk is a resource handler for dealing with AWS Elastic Beanstalk.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"beantalk"`
 
 Usage:
 
 `var beanstalk = require("beanstalk").init();`
 
  ## Example Resource
 
 ```javascript

 var beanstalk = {
      name: "beanstalk"
      module: "beanstalk"
      params: {
          ensure: "present"
          app: {
             ApplicationName: appName
             Description:     "Description"
          }
          version: {
             ApplicationName:       appName
             VersionLabel:          "VersionLabel"
             AutoCreateApplication: false
             Description:           "Description"
             Process:               true
             SourceBundle: {
                 S3Bucket: sourceBucket
                 S3Key:    "api/v0.1.1-1-g5b33818.zip"
             }
          }
          config: {
             ApplicationName: appName
             TemplateName:    "ConfigurationTemplateName"
             Description:     "Description"
             SolutionStackName: "64bit Amazon Linux 2016.03 v2.1.0 running Go 1.4"
             // See: https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/command-options-general.html
             // OptionSettings: [
             //          {
             //              Namespace:    "aws:elasticbeanstalk:application"
             //              OptionName:   "Application Healthcheck URL"
             //              Value:        "/hc"
             //          }
             // ]
          }
          environment: {
             ApplicationName: appName
             CNAMEPrefix:     cname
             Description:     "test"
             EnvironmentName: "test"
             Tags: [
                 {
                     Key:   "TagKey"
                     Value: "TagValue"
                 }
             ]
             SolutionStackName: "64bit Amazon Linux 2016.03 v2.1.0 running Go 1.4"
             VersionLabel: "VersionLabel"
          }
          terminate: {
             EnvironmentName:    "test"
             ForceTerminate:     true
             TerminateResources: true
          }
      }
 }

 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"` and ...
 
 ### `app`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticbeanstalk.html#type-CreateApplicationInput)

 Specifies the application to be created/deleted.

 ### `version`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticbeanstalk.html#type-CreateApplicationVersionInput)

 Specifies properties for an application version to be created/deleted.

 ### `environment`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticbeanstalk.html#type-CreateEnvironmentInput)

 Specifies properties for a beanstalk application environment to be created/deleted.

 ### `config`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticbeanstalk.html#type-CreateApplicationVersionInput)

 Specifies properties for a beanstalk application configuration to be created/deleted.

 ### `terminate`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/elasticbeanstalk.html#type-TerminateEnvironmentInput)

 Specifies properties for termination of an environment.

 ### `on_find`

 * Required: true
 * Allowed Values: A function taking two parameters: `catalog` and `resource`

 If defined in the resource's `params` object, the `on_find`
 function provides a way for a matching resource to be identified
 using a user-defined way.  The function is called with the current
 `catalog`, as well as the `resource` object itself.  The function
 can look through the catalog, find a matching object using whatever
 logic you want, and return it.  If the function returns `undefined`
 or a n empty Javascript array, (`[]`), the function is indicating
 that no matching resource was found in the `catalog`.
 

