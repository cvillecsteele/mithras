function run() {

    console.log("RUNNING", filepath.join(mithras.HOME, "js", "test"));

    var sst = new (require('sst'))();
    _.extend(this, sst);
    
    var tests = {};
    filepath.walk(filepath.join(mithras.HOME, "js", "test"), function(path, info, err) {
        if (!info.IsDir) {
            var ext = filepath.ext(path).substring(1);
            if (ext === "js") {
                tests[path] = require(path);
            }
        }
    });

    _.each(tests, function(t, path) {
        t.run();
    });
    
    sst.run();
    sst.report();

    return true;
}
