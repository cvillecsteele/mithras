 


 # CORE FUNCTIONS: BEANSTALK


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.beanstalk.apps.create](#app_create)
 > * [aws.beanstalk.apps.delete](#app_delete)
 > * [aws.beanstalk.apps.describe](#app_describe)
 > * [aws.beanstalk.apps.scan](#app_scan)

 > * [aws.beanstalk.versions.create](#version_create)
 > * [aws.beanstalk.versions.delete](#version_delete)
 > * [aws.beanstalk.versions.describe](#version_describe)
 > * [aws.beanstalk.versions.scan](#version_scan)

 > * [aws.beanstalk.environments.create](#environment_create)
 > * [aws.beanstalk.environments.terminate](#environment_terminate)
 > * [aws.beanstalk.environments.describe](#environment_describe)
 > * [aws.beanstalk.environments.scan](#environment_scan)

 > * [aws.beanstalk.configs.create](#config_create)
 > * [aws.beanstalk.configs.terminate](#config_terminate)
 > * [aws.beanstalk.configs.describe](#config_describe)
 > * [aws.beanstalk.configs.scan](#config_scan)

 > * [aws.beanstalk.storage.create](#storage_create)

 > * [aws.beanstalk.check](#check)
 > * [aws.beanstalk.swapCNAMEs](#swap)

 This API exposes functions to manage AWS Elastic Beanstalk.

 ## AWS.BEANSTALK.APPS.CREATE
 <a name="app_create"></a>
 `aws.beanstalk.apps.create(region, config);`

 Create an application.

 Example:

 ```

 var app = aws.beanstalk.apps.create(defaultRegion,
            {
              ApplicationName: "appName"
              Description:     "Description"
            });


 ```

 ## AWS.BEANSTALK.APPS.DELETE
 <a name="app_delete"></a>
 `aws.beanstalk.apps.delete(region, appName, force);`

 Delete an application.

 Example:

 ```

 aws.beanstalk.apps.delete("us-east-1", "appName", true);

 ```


 ## AWS.BEANSTALK.APPS.DESCRIBE
 <a name="app_describe"></a>
 `aws.beanstalk.apps.describe(region, appName);`

 Get information about a beanstalk app.

 Example:

 ```

  var app = aws.beanstalk.apps.describe("us-east-1", "appName");

 ```

 ## AWS.BEANSTALK.APPS.SCAN
 <a name="app_scan"></a>
 `aws.beanstalk.apps.scan(region, cache_id);`

 Get information about all beanstalk apps.

 Example:

 ```

  var apps = aws.beanstalk.apps.scan("us-east-1");

 ```

 ## AWS.BEANSTALK.VERSIONS.CREATE
 <a name="version_create"></a>
 `aws.beanstalk.versions.create(region, config);`

 Create an application version.

 Example:

 ```

 var ver = aws.beanstalk.versions.create(defaultRegion,
                            {
                                ApplicationName:       "appName"
                                VersionLabel:          "VersionLabel"
                                AutoCreateApplication: false
                                Description:           "Description"
                                Process:               true
                                SourceBundle: {
                                    S3Bucket: "sourceBucket"
                                    S3Key:    "api/v0.1.1-1-g5b33818.zip"
                                }
                            });


 ```

 ## AWS.BEANSTALK.VERSIONS.DELETE
 <a name="version_delete"></a>
 `aws.beanstalk.versions.delete(region, config);`

 Delete an application version.

 Example:

 ```

 aws.beanstalk.versions.delete("us-east-1",
 {
   ApplicationName:    "appName"
   VersionLabel:       "VersionLabel"
   DeleteSourceBundle: true
 });

 ```


 ## AWS.BEANSTALK.VERSIONS.DESCRIBE
 <a name="version_describe"></a>
 `aws.beanstalk.versions.describe(region, appName, versionLabel);`

 Get information about a beanstalk application version.

 Example:

 ```

  var ver = aws.beanstalk.versions.describe("us-east-1", "appName", "VersionLabel");

 ```

 ## AWS.BEANSTALK.VERSIONS.SCAN
 <a name="version_scan"></a>
 `aws.beanstalk.versions.scan(region);`

 Get information about all beanstalk versions.

 Example:

 ```

  var versions = aws.beanstalk.versions.scan("us-east-1");

 ```

 ## AWS.BEANSTALK.ENVIRONMENTS.CREATE
 <a name="environment_create"></a>
 `aws.beanstalk.environments.create(region, config);`

 Create an application environment.

 Example:

 ```

 var env = aws.beanstalk.environments.create(defaultRegion,
 {
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
 });


 ```

 ## AWS.BEANSTALK.ENVIRONMENTS.TERMINATE
 <a name="environment_terminate"></a>
 `aws.beanstalk.environments.terminate(region, config);`

 Terminate an application environment.

 Example:

 ```

 aws.beanstalk.environments.terminate("us-east-1",
                                   {
                                       EnvironmentName:    "test"
                                       ForceTerminate:     true
                                       TerminateResources: true
                                   });

 ```


 ## AWS.BEANSTALK.ENVIRONMENTS.DESCRIBE
 <a name="environment_describe"></a>
 `aws.beanstalk.environments.describe(region, appName, envName);`

 Get information about a beanstalk application environment.

 Example:

 ```

  var ver = aws.beanstalk.environments.describe("us-east-1", "appName", "EnvironmentLabel");

 ```

 ## AWS.BEANSTALK.ENVIRONMENTS.SCAN
 <a name="environment_scan"></a>
 `aws.beanstalk.environments.scan(region);`

 Get information about all beanstalk environments.

 Example:

 ```

  var apps = aws.beanstalk.environments.scan("us-east-1");

 ```

 ## AWS.BEANSTALK.CONFIGS.CREATE
 <a name="config_create"></a>
 `aws.beanstalk.configs.create(region, config);`

 Create an application config.

 Example:

 ```

 var c = aws.beanstalk.configs.create(defaultRegion,
                           {
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
                           });


 ```

 ## AWS.BEANSTALK.CONFIGS.DELETE
 <a name="config_delete"></a>
 `aws.beanstalk.configs.delete(region, appName, templateName);`

 Delete an application config.

 Example:

 ```

 aws.beanstalk.configs.delete("us-east-1", "appName", "ConfigurationTemplateName")

 ```


 ## AWS.BEANSTALK.CONFIGS.DESCRIBE
 <a name="config_describe"></a>
 `aws.beanstalk.configs.describe(region, appName, templateName);`

 Get information about a beanstalk application config.

 Example:

 ```

  var cfg = aws.beanstalk.configs.describe("us-east-1", "appName", "ConfigurationTemplateName");

 ```

 ## AWS.BEANSTALK.CONFIGS.SCAN
 <a name="config_scan"></a>
 `aws.beanstalk.configs.scan(region);`

 Get information about all beanstalk configs.

 Example:

 ```

  var apps = aws.beanstalk.configs.scan("us-east-1");

 ```

 ## AWS.BEANSTALK.STORAGE.CREATE
 <a name="storate_create"></a>
 `aws.beanstalk.storage.create(region);`

 Create a storage location for logs.

 Example:

 ```

  aws.beanstalk.storage.create("us-east-1");

 ```

 ## AWS.BEANSTALK.SWAPCNAMES
 <a name="swapCNAMEs"></a>
 `aws.beanstalk.swapCNAMEs(region, config);`

 Exchange CNAMES to enable red/blue pattern.

 Example:

 ```

  aws.beanstalk.swapCNAMEs("us-east-1",
 {
 	DestinationEnvironmentId:   "Environment1Id"
 	DestinationEnvironmentName: "Environment1Name"
 	SourceEnvironmentId:        "Environment2Id"
 	SourceEnvironmentName:      "Environment2Name"
 });

 ```

 ## AWS.BEANSTALK.CHECK
 <a name="check"></a>
 `aws.beanstalk.check(region, cname);`

 Checks to see if a beanstalk CNAME is available.

 Example:

 ```

  var avail = aws.beanstalk.check("us-east-1", "my-beanstalk");

 ```


