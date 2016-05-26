(function() {

    var Run = function() {
        var sprintf = require('sprintf').sprintf;
        var assert = require('assert');
        var dir = "/tmp/mithrastest";
        suite('yourModuleName', function() {
            beforeEach(function() {
                fs.removeAll(dir);
                fs.mkdirAll(dir, 0777);
            });
            it('creates a file fs.create()', function(){
                assert.doesNotThrow(function() {
                    var path = filepath.join(dir, "fs.create");
                    fs.create(path);
                    var cmd = sprintf("test -e %s && echo ok", path);
                    var result = exec.run(cmd);
                    assert(result[2] === true);
                    assert(result[3] === 0);
                    assert(result[0].trim() === 'ok');
                }, "threw an exception");
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
