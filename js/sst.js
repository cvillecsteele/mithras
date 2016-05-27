(function() {
    var sprintf = require("sprintf").sprintf;

    var SST = function SST(config) {
        this.suites = {};
        this.tests = {};
        this.config = config;
    };

    SST.prototype = {
        suite: function(message, f) {
            this.currentSuite = message;
            this.suites[message] = {
                name: message
                setup: f
                tests: {}
                beforeEach: null
                afterEach: null
                before: null
                after: null
            };
        }
        before: function(f) {
            this.suites[this.currentSuite].before = f;
        }
        after: function(f) {
            this.suites[this.currentSuite].after = f;
        }
        test: function(message, f) {
            this.suites[this.currentSuite].tests[message] = {
                func: f
            };
        }
        beforeEach: function(f) {
            this.suites[this.currentSuite].beforeEach = f;
        }
        afterEach: function(f) {
            this.suites[this.currentSuite].afterEach = f;
        }
        run: function(suites) {
            if (suites && !(_.isArray(suites))) {
                suites = [suites];
            }
            _.each(this.suites, function(suite, suiteName) {
                if (suites && suites.length > 0) {
                    if (!_.find(suites, function(s) {
                        return s === suiteName;
                    })) {
                        return;
                    }
                }
                (suite.setup.bind(this))();
                if (suite.before) {
                    (suite.before)();
                }
                _.each(suite.tests, function(t, message) {
                    var f = t.func;
                    if (suite.beforeEach) {
                        (suite.beforeEach)();
                    }
                    t.ok = true;
                    try {
                        (f)();
                    } catch(err) {
                        t.err = err;
                        t.ok = false;
                    };
                    if (suite.afterEach) {
                        (suite.afterEach)();
                    }
                });
                if (suite.after) {
                    (suite.after)();
                }
            });
        }
        report: function(suites) {
            if (suites && !(_.isArray(suites))) {
                suites = [suites];
            }
            _.each(this.suites, function(suite, suiteName) {
                if (suites && suites.length > 0) {
                    if (!_.find(suites, function(s) {
                        return s === suiteName;
                    })) {
                        return;
                    }
                }
                var idx = 1;
                _.each(suite.tests, function(test, testName) {
                    var f = test.func;
                    console.log(sprintf("%s %d - %s",
                                        test.ok ? "ok" : "not ok",
                                        idx,
                                        testName));
                    if (!test.ok) {
                        console.log(sprintf("# %s", JSON.stringify(test.err)));
                    }
                    idx = idx + 1;
                });
            });
        }
    };

    // Export
    if (typeof(exports) != 'undefined') {
	module.exports = SST;
    }
})();
