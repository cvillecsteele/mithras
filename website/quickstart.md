# QUICKSTART

We're gonna get you up and running, quick.

## Install

Follow the [usage guide](usage.html) to get Mithras installed.

Install the system from source, and then get your AWS credentials set
up.  Make sure that your AWS credentials give you authorization to
create security groups and instances in AWS.

Then:

    export MITHRASHOME=$GOPATH/src/github.com/cvillecsteele/mithras
    cd ~
    mkdir mysite
    cp $MITHRASHOME/example/simple.js mysite/site.js
    cd mysite
    AWS_PROFILE=<your-profile-name> mithras -v run -f site.js up

Watch as mithras creates a security group, an SSH key, and then an instance.

To shut it all down:

    AWS_PROFILE=<your-profile-name> mithras -v run -f site.js down


