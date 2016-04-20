# WALKTHROUGH, PART 2: VPC & Instance Configuration

Use this document to get up and working quickly and easily with
Mithras.

* [Part One](guide1.html): An EC2 instance
* [Part Two](guide2.html): VPC & Configuring our instance
* [Part Three](guide3.html): A complete application stack
* [Part Four](guide4.html): A dynamically-built script

## Part Two: VPC & Configuring our EC2 Instance

This part of the guide demonstrates a somewhat more involved AWS
configuration, including a VPC, two subnets, a security group, an IAM
role, and an instance with some configuration done to it.

You'll find the script in `$MITHRASHOME/example/intermediate.js` You
should copy it to your Mithras sandbox:

    cp $MITHRASHOME/example/intermediate.js ~/projects/mysite

You'll see use of caching in this example, to reduce the overhead of
querying AWS repeatedly for resources.  (Since it can take awhile.)

We'll use explicit dependencies to order the resources so they can be
built and torn down succesfully.  This example will demonstrate
runtime configuration using `mithras.watch`.

It also uses `mithras.bootstrap`, which prepares an instance for
complex configuration operations by moving a copy of Mithras itself
onto the target instance.

Before you get going, make sure that you've [installed](usage.html)
Mithras first.  Also, double check that your AWS credentials are set
up correctly and that if you've run the first part of this script,
you've torn it down. (See [Part One](guide1.html).)

Then:

    cp -r $MITHRASHOME/example ~/mysite

Fire up your favorite editor and load
`~/mysite/example/intermediate.js` to follow along.

### Caching

The first part of this example sets up a cache.  When Mithras queries
AWS, the resulting information is saved to disk in a `"cache"` folder,
and subsequently used as the catalog:

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

### Our VPC

Resource definitions follow a brief set of variable declarations in
`intermediate.js`.  The first definition is for our VPC:

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

This sets up a VPC with a `CidrBlock` of `"172.33.0.0/16"`, and tags
it.  The default mechanism used by the `"vpc"` handler is to identify
VPCs by their `CidrBlock` value.  If, for example, you would like to
identify VPCs based on tags, consult the VPC handler
[documentation](handler_vpc.html) for more information on the
`on_find` property.

### Two Subnets

We define two subnet resources:

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

You'll notice that for some of the properties in these resource
definitions, a call to `mithras.watch` is used instead of a value.
You can find the documentation for `mithras.watch`
[here](handler_mithras.html#watch).  This function allows a resource
to examine another resource for information gathered at runtime.

In this case, the id of the VPC is not known until runtime, when AWS
has been queried and a set of existing VPCs is known.  Or perhaps this
VPC will be created during this run of Mithras.  In either case, the
VPC resource defined above will have a `_target` property set to the
VPC it matches, at runtime.  Using `mithras.watch`, we reach into the
VPC resource and extract the `VpcId` property from the matching VPC.

Use `mithras.watch` when one resource depends on information stored in
the `_target` property of another resource.

Here, the first part of the first argument to `mithras.watch` is
`VPC`.  That's the value of the `name` property in the `rVpc` object
defined above.  `Mithras.watch` uses the
[object-path](https://github.com/mariocasciaro/object-path) library to
extract property values from other resources you've defined.

Also, note that our subnets have dependencies on the VPC resource,
which is declared in the `dependsOn` property.

### A Security Group

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

Here, we define a simple "firewall" that our instances will use.  It
allows incoming connections to port `22` and `80`, and outgoing
connections to anywhere.

### IAM Instance Profile

Commonly, instances need to interact with the AWS API in some way.  A
best practice on AWS is to use instance profiles.  Here, we give our
instances permissions for any operation on S3 (just as an example):

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

See the [docs](handler_iam.html) for the `"iamProfile"` handler and
for the [IAM Core Functions](core_iam.html) for more information.

### Instance Resource

The next interesting resource is our instance:

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

We want our instance to have a security group and an IAM instance
profile, so this resource depends on those.

### Bootstrapping

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

Mithras uses copies of itself on remote instances to perform complex
configuration tasks.  Handlers like `"packager"` and `"service"`
require it.  Before we go any farther, we call `mithras.bootstrap` to
obtain a standard set of resources for setting up an instance for
complex ops using Mithras.

### The Punchline: Package Update!

The final resource tells Mithras to run `"yum update"` on our instance:

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

This is the simplest possible per-instance configuration.  See the API
<a href="api.html">Reference</a> to get an idea about what you can do
with Mithras.

### Applying Resources to the Catalog

Finally we are ready to apply these resources to the current catalog
on AWS.  For convenience, we group the resources into to containing
resources:

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

Then we apply:

    catalog = mithras.apply(catalog, [ rStack, rWSTier ], reverse);

### Caching Again

`Mithras.apply` returns an updated copy of the catalog, which includes
any changes it makes as a result of running your resources.  We take
that updated catalog and store it on disk for 5 minutes:

    // Cache it for 5 mintues.
    Cache.put("catalog", catalog, (5 * 60));

### Running the script

Run the script:

    mithras -v run -f example/intermediate.js

Watch Mithras build your stack, run your instance, and configure it.






