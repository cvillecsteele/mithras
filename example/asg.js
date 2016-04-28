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
// # Autoscaling Groups Example
// 
// This is a more advanced example, demonstrating how to work with
// Mithras, which is an agentless system, in the context of AWS's
// Autoscaling Groups, which present an interesting context and
// challenge.
// 
// If you're not using ASGs, then don't worry about this example.
// 
// Usage:
// 
//     mithras -v run -f example/asg.js
// 
// This example works with Mithas daemon mode.  After you've set up
// the ASG, using the above script, you'll run:
// 
// 
//     mithras -v daemon start -f example/asg_daemon.js
// 
// 
function run() {

    // Requires
    var sprintf = require("sprintf").sprintf;

    // Filter regions
    mithras.activeRegions = function (catalog) { return ["us-east-1"]; };
    
    // Talk to AWS
    var catalog = mithras.run();

    // Setup, variables, etc.
    var ensure = "present";
    var reverse = false;
    if (mithras.ARGS[0] === "down") { 
        var ensure = "absent";
        var reverse = true;
    }
    var defaultRegion = "us-east-1";
    var defaultZone = "us-east-1d";
    var altZone = "us-east-1b";
    var keyName = "mithras";
    var ami = "ami-22111148";
    var instanceType = "m1.small";
    var sgName = "simple-sg";
    var iamProfileName = "test-asg";
    var iamRoleName = "test-asg-iam-role";
    var asgName = "my-asg";

    // We tag (and find) our instance based on this tag
    var instanceNameTag = "mithras-instance";

    //////////////////////////////////////////////////////////////////////
    // Resource Definitions
    //////////////////////////////////////////////////////////////////////

    var rQueue = {
	name: "sqsQueue"
	module: "sqs"
	params: {
	    region: defaultRegion
	    ensure: ensure
	    queue: {
		QueueName: "myqueue"
	    }
	}
    };
    
    var rIAM = {
        name: "IAM"
        module: "iamProfile"
        params: {
            region: defaultRegion
            ensure: ensure
            profile: {
                InstanceProfileName: iamProfileName
            }
            role: {
                RoleName: iamRoleName
                AssumeRolePolicyDocument: aws.iam.roles.asgTrustPolicy
            }
            policies: {
                "sqs_full_access": {
                    "Version": "2012-10-17",
                    "Statement": [
                        {
                            "Effect": "Allow",
                            "Action": "sqs:*",
                            "Resource": "*"
                        }
                    ]
                }
            }
        }
    };
    
    // Create a keypair for instances
    var rKey = {
        name: "key"
        module: "keypairs"
        skip: (ensure === 'absent') // Don't delete keys
        params: {
            region: defaultRegion
            ensure: ensure
            key: {
                KeyName: keyName
            }
            savePath: os.expandEnv("$HOME/.ssh/" + keyName + ".pem")
        }
    };
    
    var rASG = {
        name: "ASG"
        module: "autoscaling"
        dependsOn: [rIAM.name, rKey.name]
        params: {
            region: "us-east-1"
            ensure: ensure
            group: {
                AutoScalingGroupName: asgName
                MaxSize:              2
                MinSize:              1
                AvailabilityZones:    [ defaultZone ]
                DefaultCooldown:         100
                DesiredCapacity:         1
                HealthCheckGracePeriod:  100
                LaunchConfigurationName: "lcName"
                NewInstancesProtectedFromScaleIn: true
                Tags: [
                    {
                        Key:               "Name"
                        PropagateAtLaunch: true
                        ResourceId:        asgName
                        ResourceType:      "auto-scaling-group"
                        Value:             "test"
                    },
                ]
            }
            hook: {
                AutoScalingGroupName:  asgName
                LifecycleHookName:     "hookName"
                DefaultResult:         "CONTINUE"
                HeartbeatTimeout:      100
                LifecycleTransition:   "autoscaling:EC2_INSTANCE_LAUNCHING"
                NotificationTargetARN: "arn:aws:sqs:us-east-1:286536233385:myqueue"
                RoleARN:               mithras.watch("IAM._target.Roles.0.Arn")
            }
            launchConfig: {
                LaunchConfigurationName:  "lcName"
                EbsOptimized:       false
                ImageId:            ami
		KeyName:            keyName
                InstanceMonitoring: {
                    Enabled: false
                }
                InstanceType:     instanceType
            }
        } // params
    };
    
    catalog = mithras.apply(catalog, [ rQueue, rIAM, rKey, rASG ], reverse);
    
    if (rQueue._target) {
	console.log(sprintf("Queue ARN: %s", rQueue._target));
    }

    return true;
}
