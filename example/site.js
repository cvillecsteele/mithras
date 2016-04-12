function run() {

    // Set up caching
    var Cache = (new (require("cache").Cache)).init();

    // Filter regions
    mithras.activeRegions = function (catalog) { return ["us-east-1"]; };

    // Look for cached catalog and if none found get one.
    var catalog = Cache.get("catalog");
    if (!catalog) {
	catalog = mithras.run();
    } else {
	log0("Using cached catalog.");
    }
  

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
    var lbName = "webserverLB";
    var apiSHA = "HEAD";
    var bucketName = "test-9987x"
    var iamProfileName = "test-webserver"
    var iamRoleName = "test-webserver-iam-role"

    // Resource Definitions

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

    var rRdsA = {
    	name: "rdsA"
    	module: "rds"
	dependsOn: [rVpc.name, rSubnetA.name, rSubnetB.name]
    	params: {
	    ensure: ensure
	    region: defaultRegion
	    wait: true

	    subnetGroup: {
		DBSubnetGroupDescription: "test subnet group"
		DBSubnetGroupName: "test-subnet-group"
		SubnetIds: [
		    mithras.watch("subnetA._target.SubnetId"),
		    mithras.watch("subnetB._target.SubnetId")
		]
		Tags: [
		    {
			Key:   "Foo"
			Value: "Bar"
		    }
		]
	    }

	    db: {
		DBInstanceClass:         "db.m1.small"
		DBInstanceIdentifier:    "test-rds"
		Engine:                  "mysql"
		AllocatedStorage:        10
		AutoMinorVersionUpgrade: true
		AvailabilityZone:        defaultZone
		MasterUserPassword:      "test123456789"
		MasterUsername:          "test"
		DBSubnetGroupName:       "test-subnet-group"
		DBName:                  "test"
		PubliclyAccessible:      false

		// BackupRetentionPeriod:   aws.Int64(1),
		// CharacterSetName:        aws.String("String"),
		// CopyTagsToSnapshot:      true
		// DBClusterIdentifier:     aws.String("String"),
		// DBParameterGroupName:    aws.String("String"),
		// DBSecurityGroups: []*string{
		//     aws.String("String"), // Required
		//     // More values...
		// },
		
		// EngineVersion:      aws.String("String"),
		// Iops:               aws.Int64(1),
		// KmsKeyId:           aws.String("String"),
		// LicenseModel:       aws.String("String"),
		// MonitoringInterval: aws.Int64(1),
		// MonitoringRoleArn:  aws.String("String"),
		// MultiAZ:            aws.Bool(true),
		// OptionGroupName:    aws.String("String"),
		// Port:               aws.Int64(1),
		// PreferredBackupWindow:      aws.String("String"),
		// PreferredMaintenanceWindow: aws.String("String"),
		// StorageEncrypted:           aws.Bool(true),
		// StorageType:                aws.String("String"),

		Tags: [
		    {
			Key:   "foo"
			Value: "bar"
		    },
		]

		// TdeCredentialArn:      aws.String("String"),
		// TdeCredentialPassword: aws.String("String"),
		// VpcSecurityGroupIds: [
		//     aws.String("String"), // Required
		// ]
	    }
	    delete: {
		DBInstanceIdentifier:      mithras.watch("rdsA._target.DBInstanceIdentifier")
		FinalDBSnapshotIdentifier: "byebye" + Date.now()
		SkipFinalSnapshot:         true
	    }
	}
    };

    var rCache = {
	name: "redis"
	module: "elasticache"
	dependsOn: [rVpc.name, rSubnetA.name, rSubnetB.name, rwsSG.name]
	params: {
	    ensure: ensure
	    region: defaultRegion
	    wait: true

	    subnetGroup: {
		CacheSubnetGroupDescription: "Redis Subnet Group"
		CacheSubnetGroupName:        "redis-subnet-group"
		SubnetIds: [
		    mithras.watch("subnetA._target.SubnetId"),
		    mithras.watch("subnetB._target.SubnetId")
		]
	    }
	    cache: {
		CacheClusterId:          "test-redis"
		AutoMinorVersionUpgrade: true
		CacheNodeType:           "cache.t2.small"
		CacheSubnetGroupName:    "redis-subnet-group"
		Engine:                  "redis"
		NumCacheNodes:           1
		SecurityGroupIds:        []
		Tags: [
		    {
			Key:   "Name"
			Value: "test-cluster"
		    },
		]
	    }
	    delete: {
		CacheClusterId:          "test-redis"
	    }

	}
    };

    var rElb = {
    	name: "elb"
    	module: "elb"
	dependsOn: [rVpc.name, rSubnetA.name, rSubnetB.name, rwsSG.name]
	on_delete: function(elb) { 
	    // Sometimes aws takes a bit to delete an elb, and we can't
	    // proceed with deleting until it's GONE.
	    this.delay = 30; 
	    return true;
	}
    	params: {
	    region: defaultRegion
    	    ensure: ensure

	    elb: {
		Listeners: [
		    {
			InstancePort:     80
			LoadBalancerPort: 80
			Protocol:         "http"
			InstanceProtocol: "http"
		    },
		]
		LoadBalancerName: lbName
		SecurityGroups: [
		    mithras.watch("webserverSG._target.GroupId")
		]
		Subnets: [
		    mithras.watch("subnetA._target.SubnetId"),
		    mithras.watch("subnetB._target.SubnetId")
		]
		Tags: [
		    {
			Key:   "foo"
			Value: "bar"
		    },
		]
	    }
	    attributes: {
		LoadBalancerAttributes: {
		    AccessLog: {
			Enabled:        false
			EmitInterval:   60
			S3BucketName:   "my-loadbalancer-logs"
			S3BucketPrefix: "test-app"
		    }
		    ConnectionDraining: {
			Enabled: true
			Timeout: 300
		    }
		    ConnectionSettings: {
			IdleTimeout: 30
		    }
		    CrossZoneLoadBalancing: {
			Enabled: true
		    }
		}
		LoadBalancerName: lbName
	    }
	    health: {
		HealthCheck: {
		    HealthyThreshold:   2
		    Interval:           30
		    Target:             "HTTP:80/hc"
		    Timeout:            5
		    UnhealthyThreshold: 3
		}
		LoadBalancerName: lbName
	    }
	}
    };

    var rElbDnsEntry = {
    	name: "elbDnsEntry"
    	module: "route53"
	dependsOn: [rElb.name]
    	params: {
	    region: defaultRegion
    	    ensure: ensure
	    domain: "mithras.io."
	    resource: {
		Name: "test.mithras.io."
		Type: "A"
		AliasTarget: {
		    DNSName:              mithras.watch("elb._target.DNSName")
		    EvaluateTargetHealth: true
		    HostedZoneId:         mithras.watch("elb._target.CanonicalHostedZoneNameID")
		}
	    }
	} // params
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
    
    var rWebServer = {
    	name: "webserver"
    	module: "instance"
	dependsOn: [rwsSG.name, rSubnetA.name, rIAM.name]
    	params: {
	    region: defaultRegion
    	    ensure: ensure
	    on_find: function(catalog) {
		var matches = _.filter(catalog.instances, function (i) {
		    if (i.State.Name != "running") {
			return false;
		    }
		    return (_.where(i.Tags, {"Key": "Name", "Value": "webserver"}).length > 0);
		});
		return matches;
	    }
	    tags: {
		Name: "webserver"
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
		KeyName:                           "cr"
		Monitoring: {
		    Enabled: true
		}
		// SubnetId: mithras.watch("subnetA._target.SubnetId")
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

    var rELBMembership = {
    	name: "elbmembership"
    	module: "elbMembership"
	dependsOn: [rWebServer.name, rElb.name]
    	params: {
	    region: defaultRegion
    	    ensure: ensure
	    membership: {
		LoadBalancerName: lbName
		Instances: mithras.watch("webserver._target")
	    }
	}
    };

    var rBootstrap = new mithras.bootstrap({
	name: "bootstrap"
	dependsOn: [rWebServer.name]
	params: {
	    ensure: ensure
	    become: true
	    becomeMethod: "sudo"
	    becomeUser: "root"
	    hosts: mithras.watch("webserver._target")
	}
    });

    var template = {dependsOn: [rBootstrap.name]
		    params: {
			ensure: ensure 
     			hosts: mithras.watch("webserver._target")
     			become: true
     			becomeMethod: "sudo"
     			becomeUser: "root"
		    }
		   };
    var nginx = require("nginx")(template, 
				 // base conf content - use default
				 null, 
				 // included configs - none in this case
				 null,
				 // config files for our sites
				 {
				     site1: fs.read("example/site1.conf")[0]
				 });
    nginx.dependsOn = [rBootstrap.name]

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
	    hosts: mithras.watch("webserver._target")
	}
    };

    var rGitPkg = {
	name: "git"
    	module: "packager"
	dependsOn: [rBootstrap.name]
    	params: {
	    ensure: ensure // absent, latest
	    name: "git"
	    become: true
	    becomeMethod: "sudo"
	    becomeUser: "root"
	    hosts: mithras.watch("webserver._target")
	}
    };

    var rFile = {
	name: "someFile"
    	module: "copy"
	dependsOn: [rBootstrap.name]
    	params: {
	    ensure: ensure // present, absent
	    become: true
	    becomeMethod: "sudo"
	    becomeUser: "root"
	    dest: "/tmp/foo"
	    src: "/etc/hosts"
	    content: "example content"
	    // not supported yet:
	    //   owner: "root"
	    //   group: "wheel"
	    mode: 0644
	    hosts: mithras.watch("webserver._target")
	}
    };

    var rRepo = {
	name: "apiRepo"
    	module: "git"
	dependsOn: [rGitPkg.name]
    	params: {
	    ensure: ensure
	    repo: "git@github.com:cvillecsteele/mithras.git"
	    version: apiSHA
	    dest: "api"
	    hosts: mithras.watch("webserver._target")
	}
    };

    var rS3Bucket = {
	name: "s3bucket"
	module: "s3"
	params: {
	    ensure: ensure
	    region: defaultRegion
	    bucket: {
		Bucket: bucketName
		ACL:    "public-read"
	    }
	}
    };

    var rS3Object = {
	name: "s3object"
	module: "s3"
	dependsOn: [rS3Bucket.name, rWebServer.name]
	params: {
	    ensure: ensure
	    region: defaultRegion
	    object: {
		Bucket:             bucketName
		Key:                "foo.txt"
		ACL:                "public-read"
		Body:               "This is a test\n"
		ContentType:        "text/plain"
	    }
	}
    };

    var rS3 = {
	name: "s3"
	includes: [rS3Object, rS3Bucket]
    }
    
    var rWSTier = {
	name: "instance"
	includes: [
	    rWebServer, 
	    rELBMembership, 
	    rBootstrap,
	    rUpdatePkgs,
	    rGitPkg, 
	    nginx, 
	    rRepo, rFile
	]
    }

    var rStack = {
	name: "stack"
	includes: [
	    rVpc, 
	    rSubnetA, 
	    rSubnetB, 
	    rwsSG, 
	    rElb, 
	    rElbDnsEntry,
	    rIAM
	]
    }

    mithras.apply(catalog, 
    		  [
    		      rStack,
    		      // rWSTier,
    		      // rRdsA,
    		      // rCache,
    		      // rS3
    		  ], 
    		  reverse);

    // Cache it for 5 mintues.
    Cache.put("catalog", catalog, (5 * 60));

    return true;
}
