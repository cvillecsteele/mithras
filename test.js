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
    var sgName = "simple-sg";

    // We tag (and find) our instance based on this tag
    var instanceNameTag = "mithras-instance";

    //////////////////////////////////////////////////////////////////////
    // Resource Definitions
    //////////////////////////////////////////////////////////////////////

    var rTopic = {
	name: "snsTopic"
	module: "sns"
	params: {
            region: defaultRegion
            ensure: ensure
            topic: {
		Name:  "Test"
            }
	}
    };
    var rSub = {
	name: "snsSub"
	module: "sns"
	dependsOn: [rTopic.name]
	params: {
            region: defaultRegion
            ensure: ensure
            sub: {
		Protocol: "email"
		TopicArn: mithras.watch(rTopic.name+"._target.topic")
		Endpoint: "cvillecsteele@gmail.com"
            }
	}
    };

    mithras.apply(catalog, [ rTopic, rSub ], reverse);

    return true;
}
