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
// # S3 For a Static Website
// 
// This is my "eat your own dogfood" script.  I use this to upload all
// of the [mithras.io](http://mithras.io) website, using
// dynamically-created resources.
// 
// Usage:
// 
//     mithras -v run -f example/static.js
// 
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
            var ext = filepath.ext(path).substring(1);
            var type = s3.contentTypeMap[ext];
            var result;
            if (ext != "html" && ext != "js" && ext != "css") {
                result = fs.bread(path);
            } else {
                result = fs.read(path);
            }
            if (result[1]) {
                console.log(sprintf("Error reading path '%s': %s", path, result[1]));
                os.exit(1);
            }
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
                        Body:               result[0]
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
