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

// @public @example
// 
// # A Simple Example
// 
// Start here.  This example script will set up a security group, a
// keypair, and then launch an instance into your default VPC on AWS.
// 
// This example is a good place to begin learning about Mithras.
// 
// Usage:
// 
//     mithras -v run -f example/simple.js
// 
function run() {

    // Filter regions
    mithras.activeRegions = function (catalog) { return ["us-east-1"]; };

    // Talk to AWS
    var catalog = mithras.run();

    // Setup, variables, etc.
    var ensure = "present";
    var reverse = false;
    if (mithras.ARGS[0] === "down") { 
	var ensure = "absent";
	var reverse = true;
    }
    var defaultRegion = "us-east-1";
    var defaultZone = "us-east-1d";
    var altZone = "us-east-1b";
    var keyName = "mithras"
    var ami = "ami-22111148";
    var sgName = "simple-sg";

    // We tag (and find) our instance based on this tag
    var instanceNameTag = "mithras-instance";

    //////////////////////////////////////////////////////////////////////
    // Resource Definitions
    //////////////////////////////////////////////////////////////////////

    // A simple firewall
    var security = {
    	name: "webserverSG"
    	module: "secgroup"
    	params: {
	    region: defaultRegion
    	    ensure: ensure
	    secgroup: {
		Description: "Webserver security group"
		GroupName:   sgName
	    }
	    tags: {
		Name: "webserver"
	    }
	    ingress: {
		IpPermissions: [
		    {
			FromPort:   22
			IpProtocol: "tcp"
			IpRanges: [ {CidrIp: "0.0.0.0/0"} ]
			ToPort: 22
		    },
		    {
			FromPort:   80
			IpProtocol: "tcp"
			IpRanges: [ {CidrIp: "0.0.0.0/0"} ]
			ToPort: 80
		    }
		]
	    }
    	}
    };

    // Define a keypair resource for the instance
    var rKey = {
    	name: "key"
    	module: "keypairs"
	skip: (ensure === 'absent') // Don't delete keys
    	params: {
	    region: defaultRegion
    	    ensure: ensure
	    key: {
		KeyName: keyName
	    }
	    savePath: os.expandEnv("$HOME/.ssh/" + keyName + ".pem")
	}
    };

    // This will launch an instance into your default (classic) VPC
    var rInstance = {
    	name: "instance"
    	module: "instance"
	dependsOn: [rKey.name, security.name]
    	params: {
	    region: defaultRegion
    	    ensure: ensure
	    on_find: function(catalog) {
		var matches = _.filter(catalog.instances, function (i) {
		    if (i.State.Name != "running") {
			return false;
		    }
		    return (_.where(i.Tags, {"Key": "Name", 
					     "Value": instanceNameTag}).length > 0);
		});
		return matches;
	    }
	    instance: {
		ImageId:                           ami
		MaxCount:                          1
		MinCount:                          1
		DisableApiTermination:             false
		EbsOptimized:                      false
		InstanceInitiatedShutdownBehavior: "terminate"
		KeyName:                           keyName
		Monitoring: {
		    Enabled: false
		}
		SecurityGroups: [ sgName ]
	    } // instance
	    tags: {
		Name: instanceNameTag
	    }
	} // params
    };

    catalog = mithras.apply(catalog, [ security, rKey, rInstance ], reverse);

    if (mithras.ARGS[0] === "doc") { 
	var sprintf = require("sprintf").sprintf;
	_.each(Object.keys(catalog), function(key) {
	    key = key.toString();
	    key = key.substring(0, key.length - 1);
	    var file;
	    var title;
	    if (catalog[key+"s"]) {
		if (key === "rr") {
		    file = "website/objects/rrs.md";
		    title = "RRS";
		} else {
		    file = "website/objects/"+key+".md";
		    title = key.toUpperCase();
		}
		var contents = sprintf("\n## Object Format: %s\n\n%s\n", title,
				       ("\n" + JSON.stringify(catalog[key+"s"][0], null, 2)).replace(/\n([^\n])/g, "\n    $1"));
		fs.write(file, contents, 0644);
	    } else {
		console.log(JSON.stringify(key));
	    }
	});
    }

    return true;
}
