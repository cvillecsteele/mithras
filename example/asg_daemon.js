function watcher(input) {
    var messages = aws.sqs.messages.receive("us-east-1",
					    {
						QueueUrl: mithras.ARGS[0]
						MaxNumberOfMessages: 1
						VisibilityTimeout: 10
						WaitTimeSeconds:   1
					    });
    _.each(messages, function(m) {
	var handle = m.ReceiptHandle;
	var id = m.MessageId;
	var message = JSON.parse(m.Body);
	var instance = message.EC2InstanceId;
	var transition = message.LifecycleTransition;
	if (transition === "autoscaling:EC2_INSTANCE_LAUNCHING") {
	    //
	    // Code to run Mithas against the new instance goes here
	    //
	}
	aws.sqs.messages.delete("us-east-1", 
				{
				    QueueUrl:      mithras.ARGS[0]
				    ReceiptHandle: handle
				});
    });
    return;
}

function run() {
    if (!mithras.ARGS[0]) {
	log("Missing QueueUrl argument on command line.");
	os.exit(1);
    }
    log0("Starting daemon.")
    workers.create("watcher", watcher.toString(), 10);
    workers.run("watcher")
    return true;
}

function stop(signal) {
    log0("Daemon terminating.")
    workers.stop("watcher");
    return true;
}
