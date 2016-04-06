(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var become = function(become, becomeUser, becomeMethod, command) {
	var cmd = command;
	if (become) {
	    if (becomeUser != "") {
		cmd = becomeMethod + " -u " + becomeUser + " " + command;
	    } else {
		cmd = params.BecomeMethod + " " + command;
	    }
	}
	return cmd;
    }	

    // export
    if (typeof exports !== 'undefined') {
	exports.become = become;
    }

    return module.exports;

});
