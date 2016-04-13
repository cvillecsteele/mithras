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
    
    var sprintf = require("sprintf.js").sprintf;

    var handler = {
	moduleName: "s3"
	preflight: function(catalog, resource) {
	    if (resource.module != handler.moduleName) {
		return [null, false];
	    }
	    return [null, true];
	}
	handleBucket: function(catalog, resource) {
	    if (!resource.params.bucket) {
		return;
	    }
	    var buckets = aws.s3.buckets.describe(resource.params.region, "*");
	    var bucket = _.findWhere(buckets, 
				     {"Name": resource.params.bucket.Bucket});
	    if (bucket) {
		if (resource.params.ensure === 'absent') {
		    if (mithras.verbose) {
			log(sprintf("Deleting bucket '%s'", 
				    resource.params.bucket.Bucket));
		    }
		    aws.s3.buckets.delete(resource.params.bucket.Bucket,
					  resource.params.region);
		}
	    } else {
		if (resource.params.ensure === 'present') {
		    if (mithras.verbose) {
			log(sprintf("Creating bucket '%s'", 
				    resource.params.bucket.Bucket));
		    }
		    var res = aws.s3.buckets.create(resource.params.bucket,
						    resource.params.region,
						    mithras.verbose);
		}
	    }
	}
	runObject: function (params) {
	    var sprintf = require("sprintf.js").sprintf;
	
	    var objects = aws.s3.objects.describe(params.region, 
						  params.object.Bucket,
						  params.object.Key);

	    var obj = _.findWhere(objects, 
				  {"Key": params.object.Key});
	    if (obj) {
		if (params.ensure === 'absent') {
		    if (mithras.verbose) {
			log(sprintf("Deleting object '%s'", 
				    params.object.Key));
		    }
		    aws.s3.objects.delete(params.object.Bucket,
					  params.object.Key,
					  params.region);
		}
	    } else {
		if (params.ensure === 'present') {
		    if (mithras.verbose) {
			log(sprintf("Creating object '%s'", 
				    params.object.Key));
		    }
		    var res = aws.s3.objects.create(params.object,
						    params.region,
						    mithras.verbose);
		}
	    }
	}
	handleObject: function(catalog, resource) {
	    if (!resource.params.object) {
		return;
	    }
	    var params = resource.params;
	    if (params.hosts) {
		var js = sprintf("var run = function() {\n (%s)(%s); };\n", 
				 handler.runObject.toString(),
				 JSON.stringify(_.omit(params, 'hosts')));
		for (var i in params.hosts) {
		    var instance = params.hosts[i];
		    var result = mithras.remote.mithras(instance, 
							mithras.sshUserForInstance(resource, instance), 
							mithras.sshKeyPathForInstance(resource, instance), 
							js,
							params.become,
							params.becomeUser,
							params.becomeMethod);
		    if (result[3] == 0) {
			log(sprintf("S3 object '%s' %s", 
				    params.object.Key, 
				    params.ensure));
		    }
		}
	    } else {
		handler.runObject(params);
	    }
	}
	handle: function(catalog, resources, resource) {
	    if (resource.module != handler.moduleName) {
		return [null, false];
	    }
	    handler.handleBucket(catalog, resource);
	    handler.handleObject(catalog, resource);
	    return [null, true];
	}
    };
    
    handler.init = function () {
	mithras.modules.handlers.register("s3", handler.handle);
    };
    
    return handler;
});
