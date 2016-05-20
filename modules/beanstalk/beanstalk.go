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
//
// # CORE FUNCTIONS: BEANSTALK
//

package beanstalk

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.beanstalk.apps.create](#app_create)
// > * [aws.beanstalk.apps.delete](#app_delete)
// > * [aws.beanstalk.apps.describe](#app_describe)
// > * [aws.beanstalk.apps.scan](#app_scan)
//
// > * [aws.beanstalk.versions.create](#version_create)
// > * [aws.beanstalk.versions.delete](#version_delete)
// > * [aws.beanstalk.versions.describe](#version_describe)
// > * [aws.beanstalk.versions.scan](#version_scan)
//
// > * [aws.beanstalk.environments.create](#environment_create)
// > * [aws.beanstalk.environments.terminate](#environment_terminate)
// > * [aws.beanstalk.environments.describe](#environment_describe)
// > * [aws.beanstalk.environments.scan](#environment_scan)
//
// > * [aws.beanstalk.configs.create](#config_create)
// > * [aws.beanstalk.configs.terminate](#config_terminate)
// > * [aws.beanstalk.configs.describe](#config_describe)
// > * [aws.beanstalk.configs.scan](#config_scan)
//
// > * [aws.beanstalk.storage.create](#storage_create)
//
// > * [aws.beanstalk.check](#check)
// > * [aws.beanstalk.swapCNAMEs](#swap)
//
// This API exposes functions to manage AWS Elastic Beanstalk.
//
// ## AWS.BEANSTALK.APPS.CREATE
// <a name="app_create"></a>
// `aws.beanstalk.apps.create(region, config);`
//
// Create an application.
//
// Example:
//
// ```
//
// var app = aws.beanstalk.apps.create(defaultRegion,
//            {
//              ApplicationName: "appName"
//              Description:     "Description"
//            });
//
//
// ```
//
// ## AWS.BEANSTALK.APPS.DELETE
// <a name="app_delete"></a>
// `aws.beanstalk.apps.delete(region, appName, force);`
//
// Delete an application.
//
// Example:
//
// ```
//
// aws.beanstalk.apps.delete("us-east-1", "appName", true);
//
// ```
//
//
// ## AWS.BEANSTALK.APPS.DESCRIBE
// <a name="app_describe"></a>
// `aws.beanstalk.apps.describe(region, appName);`
//
// Get information about a beanstalk app.
//
// Example:
//
// ```
//
//  var app = aws.beanstalk.apps.describe("us-east-1", "appName");
//
// ```
//
// ## AWS.BEANSTALK.APPS.SCAN
// <a name="app_scan"></a>
// `aws.beanstalk.apps.scan(region, cache_id);`
//
// Get information about all beanstalk apps.
//
// Example:
//
// ```
//
//  var apps = aws.beanstalk.apps.scan("us-east-1");
//
// ```
//
// ## AWS.BEANSTALK.VERSIONS.CREATE
// <a name="version_create"></a>
// `aws.beanstalk.versions.create(region, config);`
//
// Create an application version.
//
// Example:
//
// ```
//
// var ver = aws.beanstalk.versions.create(defaultRegion,
//                            {
//                                ApplicationName:       "appName"
//                                VersionLabel:          "VersionLabel"
//                                AutoCreateApplication: false
//                                Description:           "Description"
//                                Process:               true
//                                SourceBundle: {
//                                    S3Bucket: "sourceBucket"
//                                    S3Key:    "api/v0.1.1-1-g5b33818.zip"
//                                }
//                            });
//
//
// ```
//
// ## AWS.BEANSTALK.VERSIONS.DELETE
// <a name="version_delete"></a>
// `aws.beanstalk.versions.delete(region, config);`
//
// Delete an application version.
//
// Example:
//
// ```
//
// aws.beanstalk.versions.delete("us-east-1",
// {
//   ApplicationName:    "appName"
//   VersionLabel:       "VersionLabel"
//   DeleteSourceBundle: true
// });
//
// ```
//
//
// ## AWS.BEANSTALK.VERSIONS.DESCRIBE
// <a name="version_describe"></a>
// `aws.beanstalk.versions.describe(region, appName, versionLabel);`
//
// Get information about a beanstalk application version.
//
// Example:
//
// ```
//
//  var ver = aws.beanstalk.versions.describe("us-east-1", "appName", "VersionLabel");
//
// ```
//
// ## AWS.BEANSTALK.VERSIONS.SCAN
// <a name="version_scan"></a>
// `aws.beanstalk.versions.scan(region);`
//
// Get information about all beanstalk versions.
//
// Example:
//
// ```
//
//  var versions = aws.beanstalk.versions.scan("us-east-1");
//
// ```
//
// ## AWS.BEANSTALK.ENVIRONMENTS.CREATE
// <a name="environment_create"></a>
// `aws.beanstalk.environments.create(region, config);`
//
// Create an application environment.
//
// Example:
//
// ```
//
// var env = aws.beanstalk.environments.create(defaultRegion,
// {
//     ApplicationName: appName
//     CNAMEPrefix:     cname
//     Description:     "test"
//     EnvironmentName: "test"
//     Tags: [
//         {
//             Key:   "TagKey"
//             Value: "TagValue"
//         }
//     ]
//     SolutionStackName: "64bit Amazon Linux 2016.03 v2.1.0 running Go 1.4"
//     VersionLabel: "VersionLabel"
// });
//
//
// ```
//
// ## AWS.BEANSTALK.ENVIRONMENTS.TERMINATE
// <a name="environment_terminate"></a>
// `aws.beanstalk.environments.terminate(region, config);`
//
// Terminate an application environment.
//
// Example:
//
// ```
//
// aws.beanstalk.environments.terminate("us-east-1",
//                                   {
//                                       EnvironmentName:    "test"
//                                       ForceTerminate:     true
//                                       TerminateResources: true
//                                   });
//
// ```
//
//
// ## AWS.BEANSTALK.ENVIRONMENTS.DESCRIBE
// <a name="environment_describe"></a>
// `aws.beanstalk.environments.describe(region, appName, envName);`
//
// Get information about a beanstalk application environment.
//
// Example:
//
// ```
//
//  var ver = aws.beanstalk.environments.describe("us-east-1", "appName", "EnvironmentLabel");
//
// ```
//
// ## AWS.BEANSTALK.ENVIRONMENTS.SCAN
// <a name="environment_scan"></a>
// `aws.beanstalk.environments.scan(region);`
//
// Get information about all beanstalk environments.
//
// Example:
//
// ```
//
//  var apps = aws.beanstalk.environments.scan("us-east-1");
//
// ```
//
// ## AWS.BEANSTALK.CONFIGS.CREATE
// <a name="config_create"></a>
// `aws.beanstalk.configs.create(region, config);`
//
// Create an application config.
//
// Example:
//
// ```
//
// var c = aws.beanstalk.configs.create(defaultRegion,
//                           {
//                               ApplicationName: appName
//                               TemplateName:    "ConfigurationTemplateName"
//                               Description:     "Description"
//                               SolutionStackName: "64bit Amazon Linux 2016.03 v2.1.0 running Go 1.4"
//                               // See: https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/command-options-general.html
//                               // OptionSettings: [
//                               //          {
//                               //              Namespace:    "aws:elasticbeanstalk:application"
//                               //              OptionName:   "Application Healthcheck URL"
//                               //              Value:        "/hc"
//                               //          }
//                               // ]
//                           });
//
//
// ```
//
// ## AWS.BEANSTALK.CONFIGS.DELETE
// <a name="config_delete"></a>
// `aws.beanstalk.configs.delete(region, appName, templateName);`
//
// Delete an application config.
//
// Example:
//
// ```
//
// aws.beanstalk.configs.delete("us-east-1", "appName", "ConfigurationTemplateName")
//
// ```
//
//
// ## AWS.BEANSTALK.CONFIGS.DESCRIBE
// <a name="config_describe"></a>
// `aws.beanstalk.configs.describe(region, appName, templateName);`
//
// Get information about a beanstalk application config.
//
// Example:
//
// ```
//
//  var cfg = aws.beanstalk.configs.describe("us-east-1", "appName", "ConfigurationTemplateName");
//
// ```
//
// ## AWS.BEANSTALK.CONFIGS.SCAN
// <a name="config_scan"></a>
// `aws.beanstalk.configs.scan(region);`
//
// Get information about all beanstalk configs.
//
// Example:
//
// ```
//
//  var apps = aws.beanstalk.configs.scan("us-east-1");
//
// ```
//
// ## AWS.BEANSTALK.STORAGE.CREATE
// <a name="storate_create"></a>
// `aws.beanstalk.storage.create(region);`
//
// Create a storage location for logs.
//
// Example:
//
// ```
//
//  aws.beanstalk.storage.create("us-east-1");
//
// ```
//
// ## AWS.BEANSTALK.SWAPCNAMES
// <a name="swapCNAMEs"></a>
// `aws.beanstalk.swapCNAMEs(region, config);`
//
// Exchange CNAMES to enable red/blue pattern.
//
// Example:
//
// ```
//
//  aws.beanstalk.swapCNAMEs("us-east-1",
// {
// 	DestinationEnvironmentId:   "Environment1Id"
// 	DestinationEnvironmentName: "Environment1Name"
// 	SourceEnvironmentId:        "Environment2Id"
// 	SourceEnvironmentName:      "Environment2Name"
// });
//
// ```
//
// ## AWS.BEANSTALK.CHECK
// <a name="check"></a>
// `aws.beanstalk.check(region, cname);`
//
// Checks to see if a beanstalk CNAME is available.
//
// Example:
//
// ```
//
//  var avail = aws.beanstalk.check("us-east-1", "my-beanstalk");
//
// ```
//

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	eb "github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "elasticbeanstalk"

//////////////////////////////////////////////////////////////////////
// APPS
//////////////////////////////////////////////////////////////////////

func createApp(region string, params *eb.CreateApplicationInput) *eb.ApplicationDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.CreateApplication(params)
	if err != nil {
		log.Fatalf("Error creating application '%s': %s",
			*params.ApplicationName,
			err)
	}
	name := *resp.Application.ApplicationName

	// Wait for it.
	avail := false
	for i := 0; i < 100; i++ {
		target := describeApp(region, name)
		if target != nil && *target.ApplicationName == name {
			avail = true
			break
		}
		time.Sleep(time.Second * 10)
	}

	if !avail {
		log.Fatalf("Error creating application '%s'", name)
	}

	// Describe it.
	return describeApp(region, name)
}

func deleteApp(region string, appName string, force bool) {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &eb.DeleteApplicationInput{
		ApplicationName:     aws.String(appName),
		TerminateEnvByForce: aws.Bool(force),
	}
	_, err := svc.DeleteApplication(params)
	if err != nil {
		log.Fatalf("Error deleting application '%s': %s",
			appName,
			err)
	}

	// Wait for it.
	avail := true
	for i := 0; i < 100; i++ {
		target := describeApp(region, appName)
		if target != nil {
			avail = false
			break
		}
		time.Sleep(time.Second * 10)
	}

	if avail {
		log.Fatalf("Error deleting application '%s'", appName)
	}
}

func scanApps(region string) []*eb.ApplicationDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	in := eb.DescribeApplicationsInput{}

	resp, err := svc.DescribeApplications(&in)

	if err != nil {
		log.Fatalf("Error describing beanstalk application: %s", err)
	}

	return resp.Applications
}

func describeApp(region string, id string) *eb.ApplicationDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	in := eb.DescribeApplicationsInput{
		ApplicationNames: []*string{
			aws.String(id),
		},
	}

	resp, err := svc.DescribeApplications(&in)

	if err != nil {
		log.Fatalf("Error describing beanstalk application: %s", err)
	}

	if resp != nil && len(resp.Applications) > 0 {
		return resp.Applications[0]
	}
	return nil
}

//////////////////////////////////////////////////////////////////////
// VERSIONS
//////////////////////////////////////////////////////////////////////

func deleteVersion(region string, params *eb.DeleteApplicationVersionInput) {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	appName := *params.ApplicationName
	label := *params.VersionLabel

	_, err := svc.DeleteApplicationVersion(params)

	if err != nil {
		log.Fatalf("Error deleting version '%s' v'%s': %s",
			appName,
			label,
			err)
	}

	// Wait for it.
	avail := true
	for i := 0; i < 100; i++ {
		target := describeVersion(region, appName, label)
		if target == nil {
			avail = false
			break
		}
		time.Sleep(time.Second * 10)
	}

	if !avail {
		log.Fatalf("Error deleting version '%s' v'%s': %s",
			appName,
			label,
			err)
	}
}

func scanVersions(region string) []*eb.ApplicationVersionDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	apps := scanApps(region)

	all := []*eb.ApplicationVersionDescription{}
	for _, app := range apps {
		params := &eb.DescribeApplicationVersionsInput{
			ApplicationName: aws.String(*app.ApplicationName),
		}
		resp, err := svc.DescribeApplicationVersions(params)

		if err != nil {
			log.Fatalf("Error describing beanstalk application versions: %s", err)
		}

		if resp != nil {
			for _, v := range resp.ApplicationVersions {
				all = append(all, v)
			}
		}
	}

	return all
}

func createVersion(region string, params *eb.CreateApplicationVersionInput) *eb.ApplicationVersionDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.CreateApplicationVersion(params)
	if err != nil {
		log.Fatalf("Error creating application version '%s' v '%s': %s",
			*params.ApplicationName,
			*params.VersionLabel,
			err)
	}
	name := *params.ApplicationName
	label := *params.VersionLabel

	// Wait for it.
	avail := false
	for i := 0; i < 100; i++ {
		target := describeVersion(region, name, label)
		if target != nil {
			avail = true
			break
		}
		time.Sleep(time.Second * 10)
	}

	if !avail {
		log.Fatalf("Error creating application '%s'", name)
	}

	// Describe it.
	return describeVersion(region, name, label)
}

func describeVersions(region string, appName string) []*eb.ApplicationVersionDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &eb.DescribeApplicationVersionsInput{
		ApplicationName: aws.String(appName),
	}
	resp, err := svc.DescribeApplicationVersions(params)

	if err != nil {
		log.Fatalf("Error describing beanstalk application versions: %s", err)
	}

	if resp != nil {
		return resp.ApplicationVersions
	}
	return []*eb.ApplicationVersionDescription{}
}

func describeVersion(region string, appName string, label string) *eb.ApplicationVersionDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &eb.DescribeApplicationVersionsInput{
		ApplicationName: aws.String(appName),
		VersionLabels: []*string{
			aws.String(label),
		},
	}
	resp, err := svc.DescribeApplicationVersions(params)

	if err != nil {
		log.Fatalf("Error describing beanstalk application versions: %s", err)
	}

	if resp != nil {
		return resp.ApplicationVersions[0]
	}
	return nil
}

//////////////////////////////////////////////////////////////////////
// ENVIRONMENTS
//////////////////////////////////////////////////////////////////////

func rebuildEnvironment(region string, params *eb.RebuildEnvironmentInput) {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.RebuildEnvironment(params)
	if err != nil {
		log.Fatalf("Error rebuilding environment '%s': %s", *params.EnvironmentName, err)
	}
}

func scanEnvironments(region string) []*eb.EnvironmentDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	apps := scanApps(region)

	all := []*eb.EnvironmentDescription{}
	for _, app := range apps {
		params := &eb.DescribeEnvironmentsInput{
			ApplicationName: aws.String(*app.ApplicationName),
			IncludeDeleted:  aws.Bool(false),
		}
		resp, err := svc.DescribeEnvironments(params)
		if err != nil {
			log.Fatalf("Error describing environment '%s': %s", *app.ApplicationName, err)
		}
		if resp != nil {
			for _, e := range resp.Environments {
				all = append(all, e)
			}
		}
	}

	return all
}

func describeEnvironment(region string, appName string, envName string) *eb.EnvironmentDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &eb.DescribeEnvironmentsInput{
		ApplicationName: aws.String(appName),
		EnvironmentNames: []*string{
			aws.String(envName),
		},
		IncludeDeleted: aws.Bool(false),
	}
	resp, err := svc.DescribeEnvironments(params)

	if err != nil {
		log.Fatalf("Error describing environment '%s': %s", appName, err)
	}

	if resp != nil && len(resp.Environments) > 0 {
		return resp.Environments[0]
	}
	return nil
}

func terminateEnvironment(region string, params *eb.TerminateEnvironmentInput) *eb.EnvironmentDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	resp, err := svc.TerminateEnvironment(params)

	if err != nil {
		log.Fatalf("Error terminating environment '%s': %s",
			*params.EnvironmentName,
			err)
	}
	return resp
}

func createEnvironment(region string, params *eb.CreateEnvironmentInput) *eb.EnvironmentDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.CreateEnvironment(params)

	if err != nil {
		log.Fatalf("Error creating environment '%s': %s",
			*params.ApplicationName,
			err)
	}
	appName := *params.ApplicationName
	envName := *params.EnvironmentName

	// Wait for it.
	avail := false
	for i := 0; i < 100; i++ {
		target := describeEnvironment(region, appName, envName)
		if target != nil {
			avail = true
			break
		}
		time.Sleep(time.Second * 10)
	}

	if !avail {
		log.Fatalf("Error creating environment '%s'", appName)
	}

	// Describe it.
	return describeEnvironment(region, appName, envName)
}

//////////////////////////////////////////////////////////////////////
// CONFIG TEMPLATES
//////////////////////////////////////////////////////////////////////

func scanConfigTemplates(region string) []*eb.ConfigurationSettingsDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	envs := scanEnvironments(region)
	apps := scanApps(region)

	all := []*eb.ConfigurationSettingsDescription{}
	for _, app := range apps {
		for _, env := range envs {
			params := &eb.DescribeConfigurationSettingsInput{
				ApplicationName: aws.String(*app.ApplicationName),
				EnvironmentName: aws.String(*env.EnvironmentName),
			}
			resp, err := svc.DescribeConfigurationSettings(params)
			if err != nil {
				log.Fatalf("Error scanning config template for app '%s': %s", *app.ApplicationName, err)
			}
			if resp != nil {
				for _, c := range resp.ConfigurationSettings {
					if c != nil {
						all = append(all, c)
					}
				}
			}
		}
	}
	return all
}

func describeConfigTemplate(region string, appName string, templateName string) *eb.ConfigurationSettingsDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &eb.DescribeConfigurationSettingsInput{
		ApplicationName: aws.String(appName),
		TemplateName:    aws.String(templateName),
		// Not used?
		// EnvironmentName: aws.String("EnvironmentName"),
	}
	resp, err := svc.DescribeConfigurationSettings(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if "InvalidParameterValue" == awsErr.Code() {
				return nil
			}
		}
		log.Fatalf("Error describing config template '%s': %s", appName, err)
	}

	if resp != nil && len(resp.ConfigurationSettings) > 0 {
		return resp.ConfigurationSettings[0]
	}
	return nil
}

func deleteConfigTemplate(region string, appName string, templateName string) {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &eb.DeleteConfigurationTemplateInput{
		ApplicationName: aws.String(appName),
		TemplateName:    aws.String(templateName),
	}
	_, err := svc.DeleteConfigurationTemplate(params)

	if err != nil {
		log.Fatalf("Error deleting config template '%s': %s",
			appName,
			err)
	}

	// Wait for it.
	avail := true
	for i := 0; i < 100; i++ {
		target := describeConfigTemplate(region, appName, templateName)
		if target == nil {
			avail = false
			break
		}
		time.Sleep(time.Second * 10)
	}

	if avail {
		log.Fatalf("Error deleting config template '%s'", appName)
	}
}

func createConfigTemplate(region string, params *eb.CreateConfigurationTemplateInput) *eb.ConfigurationSettingsDescription {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.CreateConfigurationTemplate(params)

	if err != nil {
		log.Fatalf("Error creating config template '%s': %s",
			*params.ApplicationName,
			err)
	}
	appName := *params.ApplicationName
	templateName := *params.TemplateName

	// Wait for it.
	avail := false
	for i := 0; i < 100; i++ {
		target := describeConfigTemplate(region, appName, templateName)
		if target != nil {
			avail = true
			break
		}
		time.Sleep(time.Second * 10)
	}

	if !avail {
		log.Fatalf("Error creating config template '%s'", appName)
	}

	// Describe it.
	return describeConfigTemplate(region, appName, templateName)
}

//////////////////////////////////////////////////////////////////////
// OTHER
//////////////////////////////////////////////////////////////////////

func swapCNAMEs(region string, params *eb.SwapEnvironmentCNAMEsInput) {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	_, err := svc.SwapEnvironmentCNAMEs(params)
	if err != nil {
		log.Fatalf("Error swapping cnames for '%s': %s",
			*params.DestinationEnvironmentName,
			err)
	}
}

func checkDNS(region string, cname string) *eb.CheckDNSAvailabilityOutput {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	params := &eb.CheckDNSAvailabilityInput{
		CNAMEPrefix: aws.String(cname),
	}
	resp, err := svc.CheckDNSAvailability(params)
	if err != nil {
		log.Fatalf("Error checking DNS availability for '%s': %s",
			cname,
			err)
	}

	return resp
}

func createStorage(region string) string {
	svc := eb.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	var params *eb.CreateStorageLocationInput
	resp, err := svc.CreateStorageLocation(params)

	if err != nil {
		log.Fatalf("Error creating storage location: %s",
			err)
	}

	return *resp.S3Bucket
}

//////////////////////////////////////////////////////////////////////
// INJECT INTO JS-LAND
//////////////////////////////////////////////////////////////////////

func init() {
	mcore.RegisterInit(func(context *mcore.Context) {
		rt := context.Runtime
		var o1 *otto.Object
		var o2 *otto.Object
		var awsObj *otto.Object
		var ebObj *otto.Object
		if a, err := rt.Get("aws"); err != nil || a.IsUndefined() {
			awsObj, _ = rt.Object(`aws = {}`)
		} else {
			awsObj = a.Object()
		}
		if b, err := awsObj.Get("beanstalk"); err != nil || b.IsUndefined() {
			ebObj, _ = rt.Object(`aws.beanstalk = {}`)
		} else {
			ebObj = b.Object()
		}

		// APPS
		if b, err := ebObj.Get("apps"); err != nil || b.IsUndefined() {
			o2, _ = rt.Object(`aws.beanstalk.apps = {}`)
		} else {
			o1 = b.Object()
			v, _ := o1.Get("apps")
			o2 = v.Object()
		}
		o2.Set("scan", func(region string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(scanApps(region))
		})
		o2.Set("describe", func(region, id string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(describeApp(region, id))
		})
		o2.Set("delete", func(region, appName string, force bool) otto.Value {
			deleteApp(region, appName, force)
			return otto.Value{}
		})
		o2.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input eb.CreateApplicationInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for elasticbeanstalk create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall elasticbeanstalk create json: %s", err)
			}

			region := call.Argument(0).String()
			f := mcore.Sanitizer(rt)
			return f(createApp(region, &input))
		})

		// Environments
		if b, err := ebObj.Get("environments"); err != nil || b.IsUndefined() {
			o2, _ = rt.Object(`aws.beanstalk.environments = {}`)
		} else {
			o1 = b.Object()
			v, _ := o1.Get("environments")
			o2 = v.Object()
		}
		o2.Set("scan", func(region string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(scanEnvironments(region))
		})
		o2.Set("describe", func(region, appName, envName string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(describeEnvironment(region, appName, envName))
		})
		o2.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input eb.CreateEnvironmentInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for elasticbeanstalk create input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall elasticbeanstalk create json: %s", err)
			}

			region := call.Argument(0).String()
			f := mcore.Sanitizer(rt)
			return f(createEnvironment(region, &input))
		})
		o2.Set("terminate", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input eb.TerminateEnvironmentInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for elasticbeanstalk terminate env input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall elasticbeanstalk terminate env json: %s", err)
			}

			region := call.Argument(0).String()
			f := mcore.Sanitizer(rt)
			return f(terminateEnvironment(region, &input))
		})

		// Versions
		if b, err := ebObj.Get("versions"); err != nil || b.IsUndefined() {
			o2, _ = rt.Object(`aws.beanstalk.versions = {}`)
		} else {
			o1 = b.Object()
			v, _ := o1.Get("versions")
			o2 = v.Object()
		}
		o2.Set("scan", func(region string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(scanVersions(region))
		})
		o2.Set("describeForApp", func(region, appName string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(describeVersions(region, appName))
		})
		o2.Set("describe", func(region, appName, label string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(describeVersion(region, appName, label))
		})
		o2.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input eb.CreateApplicationVersionInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for elasticbeanstalk create version input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall elasticbeanstalk create version json: %s", err)
			}

			region := call.Argument(0).String()
			f := mcore.Sanitizer(rt)
			return f(createVersion(region, &input))
		})
		o2.Set("delete", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input eb.DeleteApplicationVersionInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for elasticbeanstalk delete version input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall elasticbeanstalk delete version json: %s", err)
			}

			region := call.Argument(0).String()
			deleteVersion(region, &input)
			return otto.Value{}
		})

		// Configs
		if b, err := ebObj.Get("configs"); err != nil || b.IsUndefined() {
			o2, _ = rt.Object(`aws.beanstalk.configs = {}`)
		} else {
			o1 = b.Object()
			v, _ := o1.Get("configs")
			o2 = v.Object()
		}
		o2.Set("scan", func(region string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(scanConfigTemplates(region))
		})
		o2.Set("describe", func(region, appName, configName string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(describeConfigTemplate(region, appName, configName))
		})
		o2.Set("create", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input eb.CreateConfigurationTemplateInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for elasticbeanstalk create config input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall elasticbeanstalk create config json: %s", err)
			}

			region := call.Argument(0).String()
			f := mcore.Sanitizer(rt)
			return f(createConfigTemplate(region, &input))
		})
		o2.Set("delete", func(region, appName, templateName string) otto.Value {
			deleteConfigTemplate(region, appName, templateName)
			return otto.Value{}
		})

		// Other
		if b, err := ebObj.Get("storage"); err != nil || b.IsUndefined() {
			o2, _ = rt.Object(`aws.beanstalk.storage = {}`)
		} else {
			o1 = b.Object()
			v, _ := o1.Get("storage")
			o2 = v.Object()
		}
		ebObj.Set("check", func(region string, cname string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(checkDNS(region, cname))
		})
		ebObj.Set("swapCNAMEs", func(call otto.FunctionCall) otto.Value {
			// Translate params input into a struct
			var input eb.SwapEnvironmentCNAMEsInput
			js := `(function (o) { return JSON.stringify(o); })`
			s, err := rt.Call(js, nil, call.Argument(1))
			if err != nil {
				log.Fatalf("Can't create json for elasticbeanstalk swap input: %s", err)
			}
			err = json.Unmarshal([]byte(s.String()), &input)
			if err != nil {
				log.Fatalf("Can't unmarshall elasticbeanstalk swap json: %s", err)
			}

			region := call.Argument(0).String()
			swapCNAMEs(region, &input)
			return otto.Value{}
		})
		o2.Set("create", func(region string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(createStorage(region))
		})
		if b, err := ebObj.Get("dns"); err != nil || b.IsUndefined() {
			o2, _ = rt.Object(`aws.beanstalk.dns = {}`)
		} else {
			o1 = b.Object()
			v, _ := o1.Get("dns")
			o2 = v.Object()
		}

	})
}
