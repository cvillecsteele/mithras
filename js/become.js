// MITHRAS: Javascript configuration management tool for AWS.
// Copyright (C) 2016, Colin Steele
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//   the Free Software Foundation, either version 3 of the License, or
//                  (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//              GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
