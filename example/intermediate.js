function run() {

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
    var keyName = "mithras"
    var ami = "ami-22111148";

    // We tag (and find) our instance based on this tag
    var instanceNameTag = "mithras-instance"

    //////////////////////////////////////////////////////////////////////
    // Resource Definitions
    //////////////////////////////////////////////////////////////////////

    // Define a keypair resource for the instance
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

    // This will launch an instance into your default (classic) VPC
    var rInstance = {
    	name: "instance"
    	module: "instance"
	dependsOn: [rKey.name]
    	params: {
	    region: defaultRegion
    	    ensure: ensure
	    on_find: function(catalog) {
		var matches = _.filter(catalog.instances, function (i) {
		    if (i.State.Name != "running") {
			return false;
		    }
		    return (_.where(i.Tags, {"Key": "Name", 
					     "Value": instanceNameTag}).length > 0);
		});
		return matches;
	    }
	    instance: {
		ImageId:                           ami
		MaxCount:                          1
		MinCount:                          1
		DisableApiTermination:             false
		EbsOptimized:                      false
		InstanceInitiatedShutdownBehavior: "terminate"
		KeyName:                           keyName
		Monitoring: {
		    Enabled: false
		}
	    } // instance
	    tags: {
		Name: instanceNameTag
	    }
	} // params
    };

    mithras.apply(catalog, [ rKey, rInstance ], reverse);

    return true;
}
