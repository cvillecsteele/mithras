function watcher(input) {
    console.log(Object.keys(mithras));
    var catalog = mithras.run();
    console.log(Object.keys(catalog));
    return;
}

function run() {
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
