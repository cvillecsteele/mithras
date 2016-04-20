WHAT IS MITHRAS?
###############################################################################

See the Mithras main site at [mithras.io](http://mithras.io) for more
information.

Use Mithras to deploy and configure full-stack applications in AWS,
with agentless setup and the flexibility and power of a Javascript
configuration language.

Mithras focuses on AWS, making design choices that make simple AWS
tasks easy, including features like explicit dependency declaration
and idempotent operations.

Mithras presents a simplified interface, reducing the surface area of
the AWS SDK to a digestible and easily managed size.

Using brittle data languages such as YAML make customizing your
configuration unecessarily convoluted.  Mithras chooses a well-known
and powerful language for descriving your AWS configurations:
Javascript.  You won't have to twist an unsuitable language to do
what's needed for your configuration.

Here are some uses cases for Mithras:

* Provisioning
* Configuration Management
* App Deployment
* Continuous Delivery
* Orchestration

Sounds pretty nifty, eh?  But it's not perfect.  If you need
enterprise configuration management features, Mithras isn't a good
choice.  (Not yet, anyway.)  If you have the need to manage hosts and
resources outside AWS, Mithras probably isn't a good choice.  If you
love YAML, don't bother with Mithras.

Finally, Mithras is *new*.  It's currently alpha quality software,
with bugs and design choices still being shaken out.  Proceed with
caution.

DESIGN
###############################################################################

* AWS first
* Javascript DSL
* Simple (abstract) aws interface - not a port of the aws sdk to js
* Agentless 
* Idempotent  
* Declarative resources
* Explicit dependencies
* Immutable infrastructure

BUILDING FROM SOURCE
###############################################################################

1. Set up a standard Go 1.6 environment http://golang.org/doc/code.html and ensure that $GOPATH environment variable properly set.
2. `go get github.com/cvillecsteele/mithras`.
3. `cd $GOPATH/src/github.com/cvillecsteele/mithras`
4. `go install` to build the mithras binary
5. Make sure $GOPATH/bin is in your $PATH.

RUNNING
###############################################################################

0. Mithras depends on using `ssh-agent`.  Make sure it is set up and has the right keys added.
1. Set up a site directory. Let's say it's '~/projects/my_site'
2. `EXPORT MITHRASHOME=~/projects/my_site`.
3. Create your site file, `site.js`.
4. Set up your AWS credentials file and specify a profile to use with `export AWS_PROFILE=...`.  (See: https://github.com/aws/aws-sdk-go/wiki/configuring-sdk)
5. Build remote runner and wrapper with `mithras build`
6. `mithras -v run`

To run the example from the mithras repo:

    `$ mithras -v run -f example/site.js`

To run a JS repl:

    `$ mithras repl`

TODO
###############################################################################

* website
  * quickstart 2
  * quickstart 3
* better null arg handling in exec.run  
* make sure no handlers are interacting with AWS outside the catalog (s3)
* fix arg order in S3 exposed JS functions  
* rework / robustify web core module

* more s3 configuration support for buckets
* test goroutines interaction with otto runtime
* select for writes
* repl should have mithras env loaded
* copy - localFile so we don't read it into memory??
* test base config in nginx
* ssh pipelines
* versioned js modules
* better nodejs package support

BUILDING DOCS
###############################################################################

Run:

    mithras -v run dev/website.js

To serve docs locally, first install [harp](http://harpjs.com), then:

    cd website && harp server

ISSUES
###############################################################################

* damn it's a big hunk of code... 28Mb???... NO... Now 30!!
* phew is otto sloowwww
