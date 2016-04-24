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

    var rQueue = {
	name: "sqsQueue"
	module: "sqs"
	params: {
            region: defaultRegion
            ensure: ensure
            queue: {
		QueueName:  "Test"
            }
	}
    };

    mithras.apply(catalog, [ rQueue ], reverse);

    return true;
}
