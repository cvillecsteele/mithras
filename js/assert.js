(function() {

    function AssertionError(message, info) {
        this.message = message;
        this.name = "AssertionError";
        this.info = info;
    }

    var Assert = function(value, message) {
        if (!value) {
            throw new AssertionError(message);
        }
    };

    Assert.ok = Assert;

    Assert.equals = function(actual, expected, message) {
        if (!_.isEqual(actual, expected)) {
            throw new AssertionError(message);
        }
    };

    Assert.doesNotEqual = function(actual, expected, message) {
        if (_.isEqual(actual, expected)) {
            throw new AssertionError(message);
        }
    };

    Assert.throws = function(f, pred, message) {
        try {
            f()
        } catch (e) {
            if (!pred(e)) {
                throw new AssertionError(message);
            }
        }
        throw new AssertionError(message);
    };

    Assert.doesNotThrow = function(f, message) {
        try {
            f()
        } catch (e) {
            throw new AssertionError(message, e);
        }
    };

    // Export
    if (typeof(exports) != 'undefined') {
	module.exports = Assert;
    }
})();
