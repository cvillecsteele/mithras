# WALKTHROUGH, PART 3: A Complete Application Stack

Use this document to get up and working quickly and easily with
Mithras.

* [Part One](guide1.html): An EC2 instance
* [Part Two](guide2.html): VPC & Configuring our instance
* [Part Three](guide3.html): A complete application stack
* [Part Four](guide4.html): A dynamically-built script

## Part Three: A Complete Application Stack

This part of the guide demonstrates a somewhat "full" environment, including:

* VPC
* Subnets
* Security Group
* IAM Instance Profile
* RDS Cluster
* Elasticache Cluster
* ELB
* ELB Instance Membership
* Route53 DNS Entry for ELB
* Instance
* Instance Setup
  * Package Update
  * Package Installation
  * Git Clone
  * Nginx Installation and Configuration

Just for fun, though they don't do anything "real", it also includes:

* File Manipulation on the Instance
* S3 Bucket Creation
* S3 Object Creation

Since you've gotten this far through the Guide, you already have some
basic familiarity with Mithras.  This script does introduce a few new
things, however.

To get rolling:

    cp -r $MITHRASHOME/example ~/mysite

Then fire up your favorite editor and load `~/mysite/example/site.js`
to follow along.

### Dealing with AWS Slowness

    var rElb = {
        name: "elb"
        module: "elb"
        dependsOn: [rVpc.name, rSubnetA.name, rSubnetB.name, rwsSG.name]
        on_delete: function(elb) { 
            // Sometimes aws takes a bit to delete an elb, and we can't
            // proceed with deleting until it's GONE.
            this.delay = 30; 
            return true;
        }
        params: {
            region: defaultRegion
            ensure: ensure

            elb: {...}
            ...
       }
    }

Here, you'll notice this use of a new resource callback, `on_delete`.
The `"elb"` handler implements this callback, giving you a way to
modify the runtime value of the resource when the ELB has been
deleted.  Why?  Sometimes AWS is slow.

### Nginx

    var template = {dependsOn: [rBootstrap.name]
                    params: {
                        ensure: ensure 
                        hosts: mithras.watch(rWebServer.name+"._target")
                        become: true
                        becomeMethod: "sudo"
                        becomeUser: "root"
                    }
                   };
    var nginx = require("nginx")(template, 
                                 // base conf content - use default
                                 null, 
                                 // included configs - none in this case
                                 null,
                                 // config files for our sites
                                 {
                                     site1: fs.read("example/site1.conf")[0]
                                 });
    nginx.dependsOn = [rBootstrap.name]

Here, we are using a Mithras Module, a javascript package which builds
on core Mithras functions.  In this case it's the
[nginx](mod_nginx.html) module.  This module returns a set of
resources to the caller, which neatly package up all of the work to
install and configure Nginx on an instance.

We call it with a template, which it copies into all of its included
resources.  We also give it a configuration file to install on the
remote instance.

### Working with Git

    var rRepo = {
        name: "apiRepo"
        module: "git"
        dependsOn: [rGitPkg.name]
        params: {
            ensure: ensure
            repo: "git@github.com:cvillecsteele/mithras.git"
            version: apiSHA
            dest: "mithras"
            hosts: mithras.watch(rWebServer.name+"._target")
        }
    };

This resource will clone a Git repository from GitHub onto the remote
instance.  This functionality is one of the reasons that Mithras uses
`ssh-agent`.  It uses SSH Key Forwarding to pass your keys from the
machine running Mithras through the remote instance and on to GitHub,
so you don't have to mess with moving your sensitive keys all round
the 'net.  See
[here](https://developer.github.com/guides/using-ssh-agent-forwarding/)
for more information about SSH key forwarding.

### Run It

Try it:

    mithras -v run -f example/site.js up

Be prepared to wait a bit.  Spinning up RDS and Elasticache Clusters
can be slowwww.  To tear it all down:

    mithras -v run -f example/site.js down


