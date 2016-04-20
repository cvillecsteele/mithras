# WALKTHROUGH, PART 4: Building Resources Dynamically

Use this document to get up and working quickly and easily with
Mithras.

* [Part One](guide1.html): An EC2 instance
* [Part Two](guide2.html): VPC & Configuring our instance
* [Part Three](guide3.html): A complete application stack
* [Part Four](guide4.html): A dynamically-built script

## Part Four: A Dynamically-Built Script

This part of the guide falls firmly into the "I'm eating my own
dogfood" category.  The main website for Mithras is
[http://mithras.io](http://mithras.io), which is a static site, built
from the source code in the Mithras repository, and uploaded to S3.

This script demonstrates how using Javascript and the Mithras core
functions, you can dynamically build a set of resource definitions.
In this case, we're looking at the local filesystem, but you could
just as easily use information in the catalog of AWS resources
returned by `mithras.run()`, or any other external source of
information. The sky is the limit.

To get rolling:

    cp -r $MITHRASHOME/example ~/mysite

Then fire up your favorite editor and load `~/mysite/example/static.js`
to follow along.

### A Hole in the Bucket

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

Check out the [documentation](handler_s3.html) for the `"s3"` handler.
We want an S3 bucket configured to serve our website.  Nothing fancy
here.

### Dynamic-ness

This is the kind of thing that Mithras excels at.  You can't do this
very easily with brittle data languages like YAML and JSON.

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

This code walks through the local filesystem, looking for files.  For
each file it finds, it creates a resource with some dynamically set
properties, and adds it to the array in the `objects` var.  Note the
use of the [filepath](core_filepath.html) and [fs](core_fs.html) functions.

### DNS

Last but not least, we need a DNS entry that points to our S3 bucket:

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

Voila!

### That's It

You have graduated.

For extra credit, have a look at the Nginx module
[here](https://github.com/cvillecsteele/mithras/blob/master/js/nginx/nginx.js).
Happy Mithras-ing!
