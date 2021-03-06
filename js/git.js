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

// @public
// 
// # Git
// 
// Git is resource handler for manipulating git repos.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"git"`
// 
// Usage:
// 
// `var git = require("git").init();`
// 
//  ## Example Resource
// 
// ```javascript
// var rGit = {
//   name: "someGit"
//   module: "git"
//   dependsOn: [otherResource.name]
//   params: {
//     ensure: "present"
//     repo: "git@github.com:cvillecsteele/mithras.git"
//     version: "HEAD"
//     dest: "myrepo"
//     hosts: [<array of ec2 instance objects>]
//   }
// };
// ```
// 
// ## Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "absent, "absent"
//
// If `"present"`, the handler will ensure the repo exists, and it
// not, it will be cloned into `dest`.  If `"absent"`, the repo is
// removed from `dest`.
// 
// ### `hosts`
//
// * Required: true
// * Allowed Values: an array of ec2 instance objects
//
// This property specifies the hosts on which this resource is to be applied.
//
// ### `repo`
//
// * Required: true
// * Allowed Values: string specifying git repo
//
// This git repo will be cloned.
//
// ### `dest`
//
// * Required: true
// * Allowed Values: path to destination for cloned repo
//
// This property specifies the path into which the repo will be cloned.
//
// ### `version`
//
// * Required: true
// * Allowed Values: a git SHA or other clone-able commit-ish id
//
// The clone will fetch this version of the repo.
//
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
	module.exports = factory();
    }
})(this, function() {
    
    var become = require("become").become;
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleNames: ["git"]
	check: function(resource, instance, user, key) {
	    var p = resource.params;
	    var cmd = sprintf("test -d '%s' && cd '%s' && git rev-parse --is-inside-work-tree > /dev/null 2>&1 && git rev-parse HEAD", p.dest, p.dest)
	    cmd = become(p.become, p.becomeUser, p.becomeMethod, cmd);
	    result = mithras.remote.shell(instance.PublicIpAddress, 
					  user, 
					  key, 
					  null,
					  cmd,
					  null);
	    var out = result[0].trim();
	    var err = result[1].trim();
	    var ok = result[2];
	    var status = result[3];
	    var version = "";
	    if (ok && status == 0) {
		var lines = out.split("\n");
		for (var i = 0; i < lines.length; i++) {
		    var parts = lines[i].split(/[ \t]/);
		    for (var j = 0; j < parts.length; j++) {
			if (parts[j].match(/^[a-z0-9]+/)) {
			    version = parts[j];
			}
		    }
		}
		if (mithras.verbose && version != "") {
		    log(sprintf("Repo '%s' @ SHA '%s'",
				p.repo,
				version));
		    return version;
		}
	    } else if (status == 255) {
		log(sprintf("Error communicating with remote system '%s', repo '%s', dest '%s': %s",
			    instance.PublicIpAddress, p.repo, p.dest, err.trim()));
		os.exit(3);
	    } else if (status == 1 && mithras.verbose) {
		log(sprintf("Git '%s', dest '%s' not found.", p.repo, p.dest));
	    } else if (mithras.verbose) {
		    log(sprintf("Git repo '%s': status %d; out %s", 
				p.repo, status, out));
	    }
	    return false;
	}
	
	install: function(resource, inst, user, key) {
	    var p = resource.params;
	    var cmd = "";
	    switch (p.ensure) {
	    case "present":
		cmd = sprintf("git clone %s %s && cd %s && git reset --hard %s && git rev-parse %s", 
			      p.repo, p.dest, p.dest, p.version, p.version);
		break;
	    case "absent":
		cmd = sprintf("rm -rf %s", p.dest);
		break;
	    }

	    cmd = become(p.become, p.becomeUser, p.becomeMethod, cmd);
	    cmd = "GIT_SSH_COMMAND='ssh -o ForwardAgent=yes -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no' " + cmd;

	    var result = mithras.remote.shell(inst.PublicIpAddress, user, key, null, cmd, null, false);

	    var out = result[0];
	    var err = result[1];
	    var ok = result[2];
	    var status = result[3];
	    if (ok && status == 0) {
		if (mithras.verbose) {
		    log(sprintf("Git '%s' dest '%s': %s", p.repo, p.dest, p.ensure));
		}
		return true;
	    } else if (status == 255) {
		log(sprintf("Remote SSH error communicating with remote system '%s', repo '%s': %s %s",
			    inst.PublicIpAddress, 
                            p.repo, 
                            err.trim(), 
                            out.trim()));
		os.exit(3);
	    } else if (status == 1) {
		if (mithras.verbose) {
		    log(sprintf("Remote Git host '%s' error: %s %s", 
                                inst.PublicIpAddress,
                                err, 
                                out));
		}
		os.exit(3);
	    } else {
		if (mithras.verbose) {
		    log(sprintf("Remote Git host '%s': status %d out %s", 
				inst.PublicIpAddress,
				status, 
				out));
		}
		os.exit(3);
	    }
	    return false;
	}

	handle: function(catalog, resources, resource) {
	    if (!_.find(handler.moduleNames, function(m) { 
		return resource.module === m; 
	    })) {
		return [null, false];
	    }
		
	    var params = resource.params;
	    var ensure = params.ensure;
		    
	    // Sanity
	    if (!params || !params.repo || !params.version) {
		console.log("Invalid git resource params")
		exit(1);
	    }
	    
	    // Loop over hosts
	    if (!Array.isArray(params.hosts)) {
		return [null, true];
	    }
	    _.each(params.hosts, function(host) {
		var key = mithras.sshKeyPathForInstance(resource, host);
		var user = mithras.sshUserForInstance(resource, host);
		if (mithras.verbose) {
		    log(sprintf("Host: '%s' (%s)", host.PublicIpAddress, 
				host.InstanceId));
		}
		
		// update for by-host variance
		_.find(resources, function(r) {
		    return r.name === resource.name;
		})._currentHost = host;
		var updated = mithras.updateResource(resource, 
						     catalog, 
						     resources, 
						     resource.name);
		params = updated.params;
		
		var version = handler.check(updated, host, user, key);
		
		switch(ensure) {
		case "absent":
		    if (version) {
			handler.install(updated, host, user, key);
		    }
		    break;
		case "present":
		    if (!version) {
			handler.install(updated, host, user, key);
		    }
		    break;
		case "latest":
		    console.log("Git resource: ensure 'latest' not supported yet.")
		    os.exit(3);
		    break;
		}
	    });
	    
	    return [null, true];
	}
    }

    handler.init = function () {
        _.each(handler.moduleNames, function(name) {
            mithras.modules.handlers.register(name, handler.handle);
        });
	return handler;
    };
    
    return handler;
});
