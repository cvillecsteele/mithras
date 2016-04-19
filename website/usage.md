# INSTALLATION & USAGE

## Building From Source

1. Set up a standard [Go 1.6 environment](http://golang.org/doc/code.html) and ensure that your `$GOPATH` environment variable properly set.
2. `go get github.com/cvillecsteele/mithras`
3. `cd $GOPATH/src/github.com/cvillecsteele/mithras`
4. `go install` to build the mithras binary
5. Make sure `$GOPATH/bin` is in your `$PATH`.

## Configuring AWS Credentials

Before using Mithras, ensure that you've configured AWS
credentials. The best way to configure credentials on a development
machine is to use the `~/.aws/credentials` file, which might look
like:

    [default]
    aws_access_key_id = AKID1234567890
    aws_secret_access_key = MY-SECRET-KEY

You can learn more about the credentials file from this
[blog post](http://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs).

Alternatively, you can set the following environment variables:

    export AWS_ACCESS_KEY_ID=AKID1234567890
    export AWS_SECRET_ACCESS_KEY=MY-SECRET-KEY


You can select a credentials file profile to use by setting an environment variable:

    export AWS_PROFILE=nifty

Or, on the command line with Mithras:

    AWS_PROFILE=nifty mithras -v run

Read more about setting up credentials [here](https://github.com/aws/aws-sdk-go/wiki/configuring-sdk).

## Running

1. Mithras depends on using `ssh-agent`.  Make sure it is set up and has the right keys added.  Github has a [good explanation](https://developer.github.com/guides/using-ssh-agent-forwarding/).
2. `export MITHRASHOME=$GOPATH/src/github.com/cvillecsteele/mithras`
3. Set up your AWS credentials file and specify a profile to use with `export AWS_PROFILE=...`.  
4. Build remote runner and wrapper with `mithras build`
5. Set up a site directory. Let's say it's `~/projects/my_site`
6. Create your site file, `site.js`.  A good example to start with is `cp $MITHRASHOME/example/simple.js ~/projects/my_site/site.js`

To run your script:

    cd ~/project/my_site
    mithras -v run -f 

To run the example from the mithras repo:

    mithras -v run -f $MITHRASHOME/example/simple.js

To run a JS repl:

    mithras repl

