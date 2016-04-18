# QUICKSTART GUIDE

Use this document to get your first Mithras recipe working quickly and
easily, and get introduced to some of the concepts behind Mithras.

* [Part One](quickstart1.html): An EC2 instance
* [Part Two](quickstart2.html): Configuring our instance
* [Part Three](quickstart3.html): A complete application stack

## Part One: An EC2 Instance

This part of the guide will show you how to use Mithras to stand up an
instance in EC2 in your default VPC.  In addition, you learn about
some key concepts Mithras uses, including resources, handlers, the
catalog and idempotent operations.

Before you get going, make sure that you've [installed](usage.html)
Mithras first.  Also, double check that your AWS credentials are set
up correctly.

### The `run()` function

For this section of the guide, we'll be using the script found in the
`example/simple.js` file (here on
[github](https://github.com/cvillecsteele/mithras/blob/master/example/simple.js))
of the Mithras distribution, typically found in your
`$GOPATH/src/github.com/cvillecsteele/mithras` folder.  Pull that file
up in your favorite editor, and take a look around.  One thing you'll
see is the `run()` function:

```javascript
function run() {...}

```

Mithras scripts define a `run()` function, which Mithras will call
when invoked.  Broadly speaking, Mithras scripts consist of three
parts: initialization, during which the script does some setup, and
AWS is queried to create a `catalog`; definition of `resources`, in
which you lay out the various components of your application, and
relate them to one another using the resource property `dependsOn`;
and finally application of your resources to the catalog.  All of this
work is triggered by Mithras' invocation of your `run()` function.

### Config and setup

Next, you'll see a line of configuration setup:

```javascript
// Filter regions
mithras.activeRegions = function (catalog) { return ["us-east-1"]; };

```

This code restricts the scanning of AWS resources to a single region,
in this case `us-east-1`.  It is not necessary to provide this filter,
and if you do not override the `mithras.activeRegions` function, all
regions will be considered active.

Next, our script calls `mithras.run()`, which interrogrates the AWS
APIs to get a snapshot of all of the AWS resources currently active in
your account:

```javascript
// Talk to AWS
var catalog = mithras.run();

```

The return value from `mithras.run()` is captured in the `catalog`
variable.  Later, this set of Javascript objects will be used to
determine which resources need to be created or deleted in order to
satisfy your script.

Next up comes some configuration code:

```javascript
// Setup, variables, etc.
var ensure = "present";
var reverse = false;
if (mithras.ARGS[0] === "down") { 
  var ensure = "absent";
  var reverse = true;
}
var defaultRegion = "us-east-1";
var defaultZone = "us-east-1d";
var altZone = "us-east-1b";
var keyName = "mithras"
var ami = "ami-22111148";

// We tag (and find) our instance based on this tag
var instanceNameTag = "instance";

```

Note here that our script examines `mithras.ARGS`.  When you invoke
`mithras run` from teh command line, any additional arguments you
provide are passed through to your script via the `ARGS` property on
the `mithras` object.  You can use these, as this example script does,
to alter data or behavior in your script.  Here, if you invoke
`mithras run up`, the var `ensure` will be set to `"present"`, and if
you invoke mithras with `mithras run down`, the var `ensure` will have
the value `"absent"`.

Most often, `resources` indicate their desired state via a parameter
`ensure`, which generally take the values `present` or `absent`.
(Though in some cases, the associated resource handler may permit
additional values, such as `latest`, for `package` resources.)

You will also notice the script sets a few additional variables which
following resource definitions will use.

### Keypair resource

Next come resource definitions, the first of which is for an SSH keypair:

```javascript
var rKey = {
  name: "key"                 // required: resources always have a name
  module: "keypairs"          // required: resources always have a handler
  skip: (ensure === 'absent') // Optional.  Don't delete keys
  params: {
    region: defaultRegion
    ensure: ensure
    key: {
        KeyName: keyName
    }
    savePath: os.expandEnv("$HOME/.ssh/" + keyName + ".pem")
  }
};

```

Note that this code doesn't *do* anything in AWS-land.  It's not
executing any API calls; it's just creating a Javascript object and
filling in a few properties.  Nothing magic.

There are a few things to note.  The two most basic properties of a
resource are its name and it's associated handler, which are specified
by the `name` and the `module` properties, respectively.

Next, we set a `skip` property on the resource.  This is an optional
property, and in this case it's set to `true` if the `ensure` variable
has a value ot `"absent"`.  When this resource is applied to a catalog
later in this script, this `skip` property will indicate that it
should not be applied to the catalog when the caller invoked mithras
with `mithras run down`.  The effect is that the resource will create
a keypair if needed, but it will not delete the keypair when the
script is invoked for teardown.

Next up is the `params` property.  Every handler has a different set
of parameters it needs to accomplish its work, and these settings are
communicated from the resource definition to the handler via the
`params` property.  In this case, our `params` set a region (to
`defaultRegion`, defined above), the `ensure` property, and a `key`
property.  Most AWS-oriented resources require a `region`, and also
require resource-specific properites.  In this case we need a `key`
property with an object as its value, and that object has a `KeyName`
property which specifies the name of the SSH keypair in AWS.

Finally, our resource definition also inludes a `savePath` property,
which the resource handler uses to save the contents of the keypair
when (if) it is created, so it can be used to SSH into any instances
that are created with this key.

### Instance resource

Next, our script defines a resource for the EC2 instance itself.  This
definition is a bit longer, but there are a couple of important things
to note.

```
// This will launch an instance into your default (classic) VPC
var rInstance = {
        name: "instance"
        module: "instance"
        dependsOn: [rKey.name]
        params: {
            region: defaultRegion
            ensure: ensure
            on_find: function(catalog) {
                var matches = _.filter(catalog.instances, function (i) {
                    if (i.State.Name != "running") {
                        return false;
                    }
                    return (_.where(i.Tags, {"Key": "Name", 
                                             "Value": instanceNameTag}).length > 0);
                });
                return matches;
            }
            instance: {
                ImageId:                           ami
                MaxCount:                          1
                MinCount:                          1
                DisableApiTermination:             false
                EbsOptimized:                      false
                InstanceInitiatedShutdownBehavior: "terminate"
                KeyName:                           keyName
                Monitoring: {
                    Enabled: false
                }
            } // instance
            tags: {
                Name: instanceNameTag
            }
        } // params
};
```

Since you're new to Mithras, the first one to note is the `dependsOn`
property.  The value of this property should be a Javascript array of
strings.  Each element of the array tells Mithras that the current
resource has a dependency on the resource with the name as the array
element.  So in our case:

```
dependsOn = [rKey.name]
```

indicates that the resource we're defining (named `instance`) depends
on the resource `rKey`.

The effect of dependency declaration is to order the application of
resource handlers.  In AWS-land, generally you a graph of resource
dependencies, often with persistent or container-like resources as the
root(s) of the dependency tree.  For example, you may have a VPC in
which all of your AWS resources will live in.

This explicit dependency declaration allows you to control when
resources are created and deleted in relation to one another.  This is
requirement of dealing with AWS, which doesn't allow you to delete
resources with active dependencies, for example.

Moving on, the `on_find`. property deserves a special mention for this
resource.  This property is specific to resources with a `module`
value of `"instance"`, indicating that they will be handled by the
`"instance"` handler.  This handler cooperates with your resource
definitions by providing a programmatic interface to informing the
handler about which running EC2 instances are matches for this
resource definition.  The resource handler calls the `on_find` method,
invoking your code, which may perform complex (or simple) convolutions
to determine which instances in the supplied `catalog` object satisfy
your notion of instances which match this resource definition.

```javascript
on_find: function(catalog) {
    // Use underscore.js to filter catalog.instances	 
    var matches = _.filter(catalog.instances, function (i) {
        // Our iterator first looks at instance state
        if (i.State.Name != "running") {
            return false;
        }
        // Then tags are considered
        return (_.where(i.Tags, {"Key": "Name", 
                                 "Value": instanceNameTag}).length > 0);
    });
    // return an array of matching instance objects
    return matches;
}
```

Typically, as seen in this example, your `on_find` callback will look
through the `catalog.instances` array for instances matching criteria
such as a set of tags.  In this case, we consider "matches" to be
running instances with a tag `Name` with value matching the variable
`instanceNameTag`, which is `"mithras-instance"`.

Next up is the portion of the `params` object dedicated to informing
the handler how we want to create EC2 instances.  It corresponds to
the format and property names of the AWS Go SDK.  Each resource
handler's documentation will include links to the appropriate
documentation, which will outline all of the available properties and
values.

For our definition of the `instance` property, we set the `ImageId` to
the value of the `ami` variable, defined above, the `keyName` to the
value of the `keyName` variable, also defined above, and some other
properties germain to the AWS Go SDK, which as
`DisableApiTermination`.

Finally, our resource definition indicates that it wants EC2 instances
that match to be tagged with a `Name` of the value of the
`instanceNameTag` variable, which is `mithras-instance`.

### Applying resources to the catalog

Phew!  That's it for resource definition.  Now that our script has set
things up with resources, it tells Mithras to do the work of applying
those resources to the current catalog of existing AWS resources:

```
mithras.apply(catalog, [ rKey, rInstance ], reverse);
```

This code tells mithras to `apply` the resources in the second
argument, `[rKey, rInstance]`, to the `catalog` argument.  The final
argument is a boolean and if `true`, it tells Mithras that the order
of dependencies is reversed, which is appropriate for *deleting* AWS
resources, which is done in the opposite order as *creating* them.

### Running the script

Last but not least, here's how you tell Mithras to run this script.  Make sure you've set everything up according to the [usage](usage.html) instructions, first.  Then, in your terminal, run:

    mithras -v run -f example/simple.js

Since we specified a global Mithras CLI option of `-v`, we see some
pretty verbose output about what Mirthas does.  You should see it
create a keypair, and then stand up an EC2 instance.

Woot!

Now that you've seen the basics, it's time to move on to [Part Two](quickstart2.html).



