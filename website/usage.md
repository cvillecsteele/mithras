# INSTALLATION & USAGE

## BUILDING FROM SOURCE

1. Set up a standard [Go 1.6 environment](http://golang.org/doc/code.html) and ensure that your `$GOPATH` environment variable properly set.
2. `go get github.com/cvillecsteele/mithras`
3. `cd $GOPATH/src/github.com/cvillecsteele/mithras`
4. `go install` to build the mithras binary
5. Make sure `$GOPATH/bin` is in your `$PATH`.

## RUNNING

0. Mithras depends on using `ssh-agent`.  Make sure it is set up and has the right keys added.
1. Set up a site directory. Let's say it's `~/projects/my_site`
2. `EXPORT MITHRASHOME=~/projects/my_site`.
3. Create your site file, `site.js`.
4. Set up your AWS credentials file and specify a profile to use with `export AWS_PROFILE=...`.  (See: [here](https://github.com/aws/aws-sdk-go/wiki/configuring-sdk))
5. Build remote runner and wrapper with `mithras build`
6. `mithras -v run`

To run the example from the mithras repo:

    `$ mithras -v run -f example/site.js`

To run a JS repl:

    `$ mithras repl`

