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

// @public @example
// 
// # Beanstalk example
// 
// Usage:
// 
//     mithras -v run -f example/beanstalk.js
// 
// This example shows how to work with Elastic Beanstalk.
// 
function run() {
    
    // Log level
    log.setLevel("debug");

    // Filter regions
    mithras.activeRegions = function (catalog) { return ["us-east-1"]; };

    // Talk to AWS
    var catalog = mithras.run();

    ///////////////////////////////////////////////////////////////////////////
    // Variables
    ///////////////////////////////////////////////////////////////////////////

    var ensure = "present";
    var reverse = false;
    if (mithras.ARGS[0] === "down") { 
        var ensure = "absent";
        var reverse = true;
    }
    var defaultRegion = "us-east-1";
    var appName = "test";
    var cname = "tz-test";
    var sourceBucket = "tz-build-artifacts";

    var beanstalk = {
    	name: "beanstalk"
    	module: "beanstalk"
	params: {
            region: defaultRegion
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
                // See: [here](https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/command-options-general.html)
                // OptionSettings: [
                //          {
                //              Namespace:    "aws:elasticbeanstalk:application"
                //              OptionName:   "Application Healthcheck URL"
                //              Value:        "/hc"
                //          }
                // ]
            }
	    terminate: {
                EnvironmentName:    "test"
                ForceTerminate:     true
                TerminateResources: true
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
                // TemplateName: "ConfigurationTemplateName"
                VersionLabel: "VersionLabel"
            }
	}
    }

    catalog = mithras.apply(catalog, [ beanstalk ], reverse);

    return true;
}
