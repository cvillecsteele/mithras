function run() {

    // Filter regions
    mithras.activeRegions = function (catalog) { return ["us-east-1"]; };

    // Talk to AWS
    catalog = mithras.run();

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

    // We tag (and find) our instance based on this tag
    var instanceNameTag = "instance"

    //////////////////////////////////////////////////////////////////////
    // Resource Definitions
    //////////////////////////////////////////////////////////////////////

    // This will launch an instance into your default (classic) VPC
    var rInstance = {
    	name: "instance"
    	module: "instance"
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
	    tags: {
		Name: instanceNameTag
	    }
	    instance: {
		ImageId:                           "ami-22111148"
		MaxCount:                          1
		MinCount:                          1
		DisableApiTermination:             false
		EbsOptimized:                      false
		InstanceInitiatedShutdownBehavior: "terminate"
		InstanceType:                      "t1.micro"
		KeyName:                           "cr"
		Monitoring: {
		    Enabled: false
		}
	    } // instance
	} // params
    };

    mithras.apply(catalog, [ rInstance ], reverse);

    return true;
}
