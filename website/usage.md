# INSTALLATION & USAGE

## Configuring AWS Credentials

Before using Mithras, ensure that you've configured AWS
credentials. The best way to configure credentials on a development
machine is to use the `~/.aws/credentials` file, which might look
like:

```
[default]
aws_access_key_id = AKID1234567890
aws_secret_access_key = MY-SECRET-KEY
```

You can learn more about the credentials file from this
[blog post](http://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs).

Alternatively, you can set the following environment variables:

```
AWS_ACCESS_KEY_ID=AKID1234567890
AWS_SECRET_ACCESS_KEY=MY-SECRET-KEY
```

## Building From Source

1. Set up a standard [Go 1.6 environment](http://golang.org/doc/code.html) and ensure that your `$GOPATH` environment variable properly set.
2. `go get github.com/cvillecsteele/mithras`
3. `cd $GOPATH/src/github.com/cvillecsteele/mithras`
4. `go install` to build the mithras binary
5. Make sure `$GOPATH/bin` is in your `$PATH`.

## Running

0. Mithras depends on using `ssh-agent`.  Make sure it is set up and has the right keys added.
1. Set up a site directory. Let's say it's `~/projects/my_site`
2. `EXPORT MITHRASHOME=~/projects/my_site`.
3. Create your site file, `site.js`.
4. Set up your AWS credentials file and specify a profile to use with `export AWS_PROFILE=...`.  
5. Build remote runner and wrapper with `mithras build`
6. `mithras -v run`

To run the example from the mithras repo:

    $ mithras -v run -f example/site.js

To run a JS repl:

    $ mithras repl

