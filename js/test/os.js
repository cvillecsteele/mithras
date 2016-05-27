(function() {

    var Run = function() {
        var sprintf = require('sprintf').sprintf;
        var assert = require('assert');
        var dir = "/tmp/mithrastest";
        suite('os', function() {
            beforeEach(function() {
                fs.removeAll(dir);
                fs.mkdirAll(dir, 0777);
            });
            test('os.getenv()', function(){
                var pwd = exec.run("echo $PWD")[0].trim();
                assert(os.getenv("PWD") === pwd);
            });
            test('os.hostname()', function(){
                var h = exec.run("hostname")[0].trim();
                assert(os.hostname()[0] === h);
            });
            afterEach(function() {
                fs.removeAll(dir);
            });
        });
    }
    
    // Export
    if (typeof(exports) != 'undefined') {
	exports.run = Run;
    }
})();
