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
// # VPC Stack
// 
// Usage:
// 
//     mithras -v run -f example/intermediate.js
// 
// In this example, we set up a "starter" VPC stack, including:
// 
// * VPC
// * Subnets
// * Security Group
// * IAM Instance Role
// * Keypair
// * Instances
// * Instance Bootstrapping
// * Instance Configuration
// 
function run() {

    // Filter regions
    mithras.activeRegions = function (catalog) { return ["us-east-1"]; };

    // Set up caching
    var Cache = (new (require("cache").Cache)).init();

    // Look for cached catalog and if none found get one.
    var catalog = Cache.get("catalog");
    if (!catalog) {
        catalog = mithras.run();
    } else {
        log0("Using cached catalog.");
    }
  

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
    var defaultZone = "us-east-1d";
    var altZone = "us-east-1b";
    var iamProfileName = "test-webserver";
    var iamRoleName = "test-webserver-iam-role";
    var keyName = "mithras";
    var webserverTagName = "webserver";

    ///////////////////////////////////////////////////////////////////////////
    // Resource Definitions
    ///////////////////////////////////////////////////////////////////////////

    var rVpc = {
        name: "VPC"
        module: "vpc"
        params: {
            region: defaultRegion
            ensure: ensure

            vpc: {
                CidrBlock:       "172.33.0.0/16"
            }
            gateway: true
            tags: {
                Name: "my-vpc"
            }
        }
    };

    var rSubnetA = {
        name: "subnetA"
        module: "subnet"
        dependsOn: [rVpc.name]
        params: {
            region: defaultRegion
            ensure: ensure

            subnet: {
                CidrBlock:        "172.33.1.0/24"
                VpcId:            mithras.watch("VPC._target.VpcId")
                AvailabilityZone: defaultZone
            }
            tags: {
                Name: "primary-subnet"
            }
            routes: [
                {
                    DestinationCidrBlock: "0.0.0.0/0"
                    GatewayId:            mithras.watch("VPC._target.VpcId", mithras.findGWByVpcId)
                }
            ]
        }
    };

    var rSubnetB = _.extend({}, rSubnetA, {
        name: "subnetB"
    });
    rSubnetB.params = _.extend({}, rSubnetA.params, {
        subnet: {
            CidrBlock: "172.33.2.0/24"
            VpcId:            mithras.watch("VPC._target.VpcId")
            AvailabilityZone: altZone
        }
    });

    var rwsSG = {
        name: "webserverSG"
        module: "secgroup"
        dependsOn: [rVpc.name]
        params: {
            region: defaultRegion
            ensure: ensure

            secgroup: {
                Description: "Webserver security group"
                GroupName:   "webserver"
                VpcId:       mithras.watch("VPC._target.VpcId")
            }
            tags: {
                Name: "webserver"
            }
            ingress: {
                IpPermissions: [
                    {
                        FromPort:   22
                        IpProtocol: "tcp"
                        IpRanges: [ {CidrIp: "0.0.0.0/0"} ]
                        ToPort: 22
                    },
                    {
                        FromPort:   80
                        IpProtocol: "tcp"
                        IpRanges: [ {CidrIp: "0.0.0.0/0"} ]
                        ToPort: 80
                    }
                ]
            }
            egress: {
                IpPermissions: [
                    {
                        FromPort:   0
                        IpProtocol: "tcp"
                        IpRanges: [ {CidrIp: "0.0.0.0/0"} ]
                        ToPort:     65535
                    }
                ]
            }

        }
    };

    var rIAM = {
        name: "IAM"
        module: "iamProfile"
        dependsOn: [rVpc.name]
        params: {
            region: defaultRegion
            ensure: ensure
            profile: {
                InstanceProfileName: iamProfileName
            }
            role: {
                RoleName: iamRoleName
                AssumeRolePolicyDocument: aws.iam.roles.ec2TrustPolicy
            }
            policies: {
                "s3_full_access": {
                    "Version": "2012-10-17",
                    "Statement": [
                        {
                            "Effect": "Allow",
                            "Action": "s3:*",
                            "Resource": "*"
                        }
                    ]
                },
            }
        }
    }
    
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

    var rWebServer = {
        name: "webserver"
        module: "instance"
        dependsOn: [rwsSG.name, rSubnetA.name, rIAM.name, rKey.name]
        params: {
            region: defaultRegion
            ensure: ensure
            on_find: function(catalog) {
                var matches = _.filter(catalog.instances, function (i) {
                    if (i.State.Name != "running") {
                        return false;
                    }
                    return (_.where(i.Tags, {"Key": "Name", 
                                             "Value": webserverTagName}).length > 0);
                });
                return matches;
            }
            tags: {
                Name: webserverTagName
            }
            instance: {
                ImageId:        "ami-60b6c60a"
                MaxCount:       1
                MinCount:       1
                DisableApiTermination: false
                EbsOptimized:          false
                IamInstanceProfile: {
                    Name: iamProfileName
                }
                InstanceInitiatedShutdownBehavior: "terminate"
                InstanceType:                      "t2.small"
                KeyName:                           keyName
                Monitoring: {
                    Enabled: true
                }
                NetworkInterfaces: [
                    {
                        AssociatePublicIpAddress: true
                        DeleteOnTermination:      true
                        DeviceIndex:              0
                        Groups:                  [ mithras.watch("webserverSG._target.GroupId") ]
                        SubnetId:                mithras.watch("subnetA._target.SubnetId")
                    }
                ]
            } // instance
        } // params
    };

    var rBootstrap = new mithras.bootstrap({
        name: "bootstrap"
        dependsOn: [rWebServer.name]
        params: {
            ensure: ensure
            become: true
            becomeMethod: "sudo"
            becomeUser: "root"
            hosts: mithras.watch(rWebServer.name+"._target")
        }
    });

    var rUpdatePkgs = {
        name: "updatePackages"
        module: "packager"
        skip: (ensure === 'absent')
        dependsOn: [rBootstrap.name]
        params: {
            ensure: "latest"
            name: ""
            become: true
            becomeMethod: "sudo"
            becomeUser: "root"
            hosts: mithras.watch(rWebServer.name+"._target")
        }
    };

    var rStack = {
        name: "stack"
        includes: [
            rVpc, 
            rSubnetA, 
            rSubnetB, 
            rwsSG, 
            rIAM
        ]
    }

    var rWSTier = {
        name: "wsTier"
        includes: [
            rKey,
            rWebServer, 
            rBootstrap,
            rUpdatePkgs,
        ]
    }

    catalog = mithras.apply(catalog, [ rStack, rWSTier ], reverse);

    // Cache it for 5 mintues.
    Cache.put("catalog", catalog, (5 * 60));

    return true;
}
