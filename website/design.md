# Mithras Design and Concepts

## System Metaphors

Mithras uses two basic [system
metaphors](http://c2.com/cgi/wiki?SystemMetaphor): the catalog, and
resources.  The catalog is a JS object with an array for every
resource type that it scans on AWS, such as instances and VPCs.

The format of the objects found in these arrays can be found
[here](api.html#objects).  Have a look through it, and you'll see that
the format of the objects corresponds to the structures in the [AWS Go
SDK](https://github.com/aws/aws-sdk-go), only translated directly into
JSON.

## Resources

In Mithras, the "stuff" that you need to either be present, absent, or
in some particular state is known as a resource.  A resource might be
a file, an instance itself, an installed package, an AWS subnet, or a
DNS entry.  All the "moving pieces" of a running system are resources.

These resources are expressed as Javascript objects in Mithras.  Your
script file.  For example:

    var rVpc = {
        name: "VPC"
        module: "vpc"
        params: {
            region: defaultRegion
            ensure: ensure

            vpc: {
                CidrBlock:       "172.33.0.0/16"
            }
            gateway: true
            tags: {
                Name: "my-vpc"
            }
        }
    };

Resources have a few pieces in common.  They all have a `name`
property.  The name is for you - it's a way to identify a resource
during Mithras execution, and to let one resource depend on another
(but more on that later).  With rare exceptions, they also have a
`module` property.  This property identies which _resource_ _handler_
will take care of translating your declaration into actions.

Generally, resources have a `params` property.  Params are
module-specific bits of configuration.  In the example above, the
`params.vpc` property holds an object used to identify and create VPCs
by the "vpc" resource handler.

The `params.ensure` property is common to almost all resource types.
Nearly all allow the values `"present"` and `"absent"`.  If
`"present"`, the handler will ensure that either the identified
resource already exists, or it is created. If `"absent"`, it will make
sure it's deleted.  In some cases, the `ensure` property can take
other values, and each resource handler module documents its accepted
values and their meanings.

Handlers that talk to AWS use the The `params.region` property to
identify how to communicate with the AWS API.  There are two
additional "special" properties for resources:

### Dependencies

When creating a resource, you may set the `dependsOn` property.  You
sent it to a Javascript array of strings.  Each string in the array
corresponds to the `name` property of another resource upon which this
one depends.  For example:

    var rSubnetA = {
        name: "subnetA"
        module: "subnet"
        dependsOn: [rVpc.name]
        params: {
            region: defaultRegion
            ensure: ensure
	    // more params properties here...
        }
    }
    
The resource in the `rSubnetA` variable depends on the `rVPC` resource
we defined above.  Therefore, the `dependsOn` property is set to an
array with the value of `rVPC.name`.

When you call [`mithras.apply(...)`](handler_mithras.html#apply), Mithras will build a dependency
graph of your resources, and then resolve things based on that
dependency ordering.

One problem I find particularly prominent when working with AWS is
that of resources that depend on other resources.  For example, you
can't delete a subnet until all of the instances in it are gone.

Mithras helps with this issue by letting you invert the dependency
graph easily, when you are tearing down resources.  The last argument
to the [`mithras.apply(...)`](handler_mithras.html#apply) call is
boolean.  If it's true, the graph is inverted.

So, you always express dependency relationships in the `dependsOn`
property in _forward_ order.

### Included Resources

One resource can include others.  By setting the `includes` property
of a resource to an array whose elements are other resources, you
instruct Mithras to operate on all of them.  This inclusion is
recursive, so A can include B which includes C, and so on.

Using included resources enables some abstraction that makes reasoning
about your system easier.  For example, the [nginx](mod_nginx.html)
module uses this feature of Mithras to package up many complex
sub-resources into one umbrella resource, which it returns to the
caller.

## The Catalog

When you call `mithras.run(...)`, Mithras goes off and talks to AWS,
asking it, "What's out there?"  It assembles all of the things it
finds on AWS into a Javascript object for you.  This is the
`catalog`.  Resource handler modules interrogate the catalog to
determine if the resource needs to be created or deleted, according to
the `ensure` property on the resource.

The catalog is keyed by entity type: VPCs, instances, etc.  Each value
is an array of the corresponding objects, which are documented
[here](api.html#objects).

If you want to restrict the regions that `mithras.run(...)` operates on,
see [`mithras.activeRegions()`](handler_mithras.html#active).

During a run of `mithras.apply(...)`, resource handler modules will
manipulate the catalog, adding (or removing) entries from it as their
corresponding entities are added or removed in AWS.

The catalog isn't anything fancy, and can be serialized with
`JSON.stringify`, saved to disk, etc.

## Applying Resources to the Catalog

When you call `mithras.apply(...)`, here's what happens:

First, Mithras pulls in all the resources that have been included in
other resources, via their `includes` property.

Then, Mithras builds and dependency graph of your resource set, in
forward order.  It also figures out what the inverted graph looks
like, in case you're tearing down resources.

### Preflight Phase

Next, Mithras _preflights_ every resource in _forward_ dependency
order.  During the preflight phase, every resource is given an
opportunity to find the corresponding entities in the catalog, and
grab a reference.  That reference is stored (by convention) on the
`_target` property of the resource.  So after preflighting, the `rVPC`
object above will have a (new) property set: `_target` will have the
corresponding VPC object stored in it.

#### Runtime Values: the `_target` Resource Property

During preflight, each resource is first recursively traversed, and
each property found is examined to see if it's a function.  If it is,
it's called with two arguments: the catalog, and a dictionary of all
of the resource set, indexed by resource name.

Here's an example snippet of a resource using a runtime-captured
value:

    var rSubnetA = {
        name: "subnetA"
        module: "subnet"
        dependsOn: [rVpc.name]
        params: {
            region: defaultRegion
            ensure: ensure

            subnet: {
                CidrBlock:        "172.33.1.0/24"
                VpcId:            mithras.watch("VPC._target.VpcId")
                AvailabilityZone: defaultZone
            }
            tags: {
                Name: "primary-subnet"
            }
            routes: [
                {
                    DestinationCidrBlock: "0.0.0.0/0"
                    GatewayId:            mithras.watch("VPC._target.VpcId", mithras.findGWByVpcId)
                }
            ]
        }
    };

Using a parameter property that is a function, your resource can look
into other resources, grab their `_target` property, and pull values
out of the AWS objects they've identified.  The
[`mithras.watch()`](handler_mithras.html#watch) higher-order function
(that is, a function that returns a function) gives you an easy way to
do this.

If the params function returns any non-undefined value, the property
is replaced with that value.  If it returns undefined, the property is
left unchanged, as a function.

Why?

Resources depend on values in other resources.  Generally, those
values are not known until runtime.  For example, to create a VPC
subnet, you need the id of the VPC.  That id isn't known when you're
building your resource script; it has to be determined at runtime.

After this recursive traversal, each preflight handler function
registered with `mithras.preflight.register` is called, giving it an
opportunity to handle the resource.  (See
[`mithras.preflight.register`](handler_mithras#modules.preflight.register) and 
[`mithras.preflight.run`](handler_mithras#modules.preflight.run).) 

### Execution Phase

Next, in either forward or inverted dependency order (depending on how
`apply` was called), Mithras will a) again recursively traverse the
resource object, calling property functions (as above), and then b)
give each resource handler module registered with
`mithras.handlers.register()` a chance to handle execution of the
resource.

In turn, every registered handler function is called with three
arguments: `catalog`, `resources` and `targetResource`.  The catalog
is the current value of the AWS resources gathered by `mithras.run()`,
which may have been modified by previous resource execution.  The
`resources` argument is the set of resources that was passed into
`mithras.apply()`.  The final argument is the resource object
currently being evaluated.

Handler functions return a Javascript array.  The second element is a
boolean. If `true`, it indicates that the handler function "owns" the
resource and has handled it.  If so, and the first element of the
return value is defined, the value of the first element is set as the
`_target` property of the resource object being evaluated.

Finally, if the resource being evaluated has a `delay` property set,
resource evaluation will stop as Mithas calls `os.sleep` with the
value of the `delay` property, effectively pausing evaluation for a
number of seconds.  Sometimes this is necessary to anticipate
propagation of information in AWS.
