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
            test('fs.create(), fs.close()', function(){
                assert.doesNotThrow(function() {
                    var path = filepath.join(dir, "fs.create");
                    var created = fs.create(path);
                    assert(created[1] == null);
                    fs.close(created[0]);
                    var cmd = sprintf("test -e %s && echo ok", path);
                    var result = exec.run(cmd);
                    assert(result[2] === true);
                    assert(result[3] === 0);
                    assert(result[0].trim() === 'ok');
                }, "threw an exception");
            });
            test('fs.write()', function(){
                assert.doesNotThrow(function() {
                    var path = filepath.join(dir, "fs.write");
                    var created = fs.write(path, "foo", 0644);
                    assert(created == null);
                    var cmd = sprintf("cat %s", path);
                    var result = exec.run(cmd);
                    assert(result[2] === true);
                    assert(result[3] === 0);
                    assert(result[0].trim() === 'foo');
                    cmd = sprintf("ls -l %s", path);
                    result = exec.run(cmd);
                    assert(result[0].spltest(/\s+/)[0] === "-rw-r--r--");
                }, "threw an exception");
            });
            test('fs.read()', function(){
                assert.doesNotThrow(function() {
                    var path = filepath.join(dir, "fs.read");
                    var cmd = sprintf("echo bar > %s", path);
                    var result = exec.run(cmd);
                    assert(fs.read(path)[0] === "bar\n");
                }, "threw an exception");
            });
            test('fs.copy()', function(){
                assert.doesNotThrow(function() {
                    var path = filepath.join(dir, "fs.copy");
                    var path2 = filepath.join(dir, "fs.copy.copy");
                    var cmd = sprintf("echo bar > %s", path);
                    var result = exec.run(cmd);
                    fs.copy(path, path2, 0644);
                    assert(fs.read(path2)[0] === "bar\n");
                }, "threw an exception");
            });
            test('fs.getwd()', function(){
                assert.doesNotThrow(function() {
                    var cwd = fs.getwd()[0];
                    var cmd = sprintf("pwd");
                    var result = exec.run(cmd);
                    assert(result[0].trim() === cwd);
                }, "threw an exception");
            });
            test('fs.chdir()', function(){
                assert.doesNotThrow(function() {
                    var cwd = fs.getwd()[0];
                    fs.chdir("/");
                    assert(fs.getwd()[0] === "/");
                    fs.chdir(cwd);
                }, "threw an exception");
            });
            test('fs.rename()', function(){
                assert.doesNotThrow(function() {
                    var path = filepath.join(dir, "fs.rename");
                    var created = fs.create(path);
                    fs.rename(path, path + ".new")
                    var cmd = sprintf("test -e %s.new && echo ok", path);
                    var result = exec.run(cmd);
                    assert(result[2] === true);
                    assert(result[3] === 0);
                    assert(result[0].trim() === 'ok');
                }, "threw an exception");
            });
            test('fs.symlink()', function(){
                assert.doesNotThrow(function() {
                    var path = filepath.join(dir, "fs.symlink");
                    var created = fs.create(path);
                    fs.symlink(path, path + ".new")
                    var cmd = sprintf("test -e %s.new && echo ok", path);
                    var result = exec.run(cmd);
                    assert(result[2] === true);
                    assert(result[3] === 0);
                    assert(result[0].trim() === 'ok');
                }, "threw an exception");
            });
            test('fs.link()', function(){
                assert.doesNotThrow(function() {
                    var path = filepath.join(dir, "fs.link");
                    var created = fs.create(path);
                    fs.link(path, path + ".new")
                    var cmd = sprintf("test -e %s.new && echo ok", path);
                    var result = exec.run(cmd);
                    assert(result[2] === true);
                    assert(result[3] === 0);
                    assert(result[0].trim() === 'ok');
                }, "threw an exception");
            });
            test('fs.dir()', function(){
                assert.doesNotThrow(function() {
                    var path = filepath.join(dir, "fs.dir");
                    var created = fs.create(path);
                    var contents = fs.dir(dir)[0];
                    assert(contents.length == 1);
                    assert(contents[0] === "fs.dir");
                }, "threw an exception");
            });
            test('fs.stat()', function(){
                assert.doesNotThrow(function() {
                    var path = filepath.join(dir, "fs.stat");
                    var created = fs.create(path);
                    var stat = fs.stat(path);
                    assert(stat.Name === "fs.stat");
                    assert(stat.Size === 0);
                    assert(stat.Perm === 420);
                    assert(stat.IsRegular === true);
                    assert(stat.IsDir === false);
                    assert(stat.Error === null);
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
