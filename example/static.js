function run() {

    s3 = require("s3");

    // Filter regions
    mithras.activeRegions = function (catalog) { return ["us-east-1"]; };

    catalog = mithras.run();
  

    ///////////////////////////////////////////////////////////////////////////
    // Variables
    ///////////////////////////////////////////////////////////////////////////

    var ensure = "present";
    var reverse = false;
    if (mithras.ARGS[0] === "down") { 
        var ensure = "absent";
        var reverse = true;
    }
    var defaultRegion = "us-east-1";
    var bucketName = "mithras.io"

    ///////////////////////////////////////////////////////////////////////////
    // Resource Definitions
    ///////////////////////////////////////////////////////////////////////////

    var bucket = {
        name: "s3bucket"
        module: "s3"
        params: {
            ensure: ensure
            region: defaultRegion
            bucket: {
                Bucket: bucketName
                ACL:    "public-read"
		LocationConstraint: defaultRegion
            }
            website: {
                Bucket: bucketName
                WebsiteConfiguration: {
                    ErrorDocument: {
                        Key: "error.html"
                    }
                    IndexDocument: {
                        Suffix: "index.html"
                    }
		}
	    } // website
	} // params
    };

    var objects = [];
    filepath.walk("website/www", function(path, info, err) {
        if (!info.IsDir) {
            var type = s3.contentTypeMap[filepath.ext(path).substring(1)];
            objects.push({
                name: path
                module: "s3"
                dependsOn: [bucket.name]
                params: {
                    ensure: "latest"
                    region: defaultRegion
		    stat: fs.stat(path)
                    object: {
                        Bucket:             bucketName
                        Key:                filepath.rel("website/www", path)[0]
                        ACL:                "public-read"
                        Body:               fs.read(path)
                        ContentType:        type
                    }
                }
            });
        }
    });

    var dns = {
    	name: "dns"
    	module: "route53"
	dependsOn: [bucket.name]
    	params: {
	    region: defaultRegion
    	    ensure: ensure
	    domain: "mithras.io."
	    resource: {
		Name:         "mithras.io."
		Type:         "A"
		AliasTarget: {
		    DNSName:              "s3-website-us-east-1.amazonaws.com"
		    EvaluateTargetHealth: false
		    HostedZoneId:         "Z3AQBSTGFYJSTF"
		}
	    }
	} // params
    };


    objects.push(bucket);
    objects.push(dns);
    mithras.apply(catalog, objects, reverse);

    return true;
}
