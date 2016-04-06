DESIGN
###############################################################################

* AWS first and only
* javascript DSL
* simple (abstract) aws interface - not a port of the aws sdk to js

BUILDING FROM SOURCE
###############################################################################

1. Setup standard Go 1.6 environment http://golang.org/doc/code.html and ensure that $GOPATH environment variable properly set.
2. `go get github.com/cvillecsteele/mithras`.
3. `cd $GOPATH/src/github.com/cvillecsteele/mithras`
4. `go install` to get binary
5. Make sure $GOPATH/bin is in your $PATH.

RUNNING
###############################################################################

0. Mithras depends on using ssh-agent.  Make sure it is set up and has the right keys added.
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

* DOCUMENTATION

* goroutines
* write select
* repl should have mithras env loaded
* copy - localFile so we don't read it into memory??
* test base config in nginx
* ssh pipelines
* versioned js modules
* better nodejs package support

ISSUES
###############################################################################

* damn it's a big hunk of code... 28Mb???
* phew is otto sloowwww
