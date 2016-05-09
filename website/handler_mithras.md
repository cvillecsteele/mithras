 
 
 # mithras
 
 Mithras is a core object for building scripts to manage your AWS
 stacks.  The `mithras` object is loaded prior to loading your
 scripts.
 
 The `mithras` object has the following properties:

 > * [ARGS](#args)
 > * [HOMe](#home)
 > * [CWD](#cwd)
 > * [GOPATH](#gopath)
 > * [MODULES](#modules)
 > * [VERSION](#version)
 > * [activeRegions](#active)
 > * [apply](#apply)
 > * [become](#become)
 > * [bootstrap](#bootstrap)
 > * [buildDeps](#buildDeps)
 > * [depGraph](#depGraph)
 > * [doIncludes](#doIncludes)
 > * [findGWByVpcId](#findGWByVpcId)
 > * [modules.handlers.register](#modules.handlers.register)
 > * [modules.handlers.run](#modules.handlers.run)
 > * [modules.preflight.register](#modules.preflight.register)
 > * [modules.preflight.run](#modules.preflight.run)
 > * [objectPath](#objectPath)
 > * [remote](#remote)
 > * [resourceMap](#resourceMap)
 > * [run](#run) 
 > * [sshKeyPathForInstance](#sshKeyPathForInstance)
 > * [sshUserForInstance](#sshUserForInstance)
 > * [traverse](#traverse)
 > * [updateResource](#updateResource)
 > * [verbose](#verbose)
 > * [watch](#watch)
 
 ## Properties
 
 ### `remote`

 See the documentation for the [remote](core_remote.html) core.
 
 ### `traverse` <a name="traverse"></a>

 See the documentation for the [traverse js](https://github.com/substack/js-traverse) module.
 
 ### `objectPath` <a name="objectPath"></a>

 See the documentation for the [object-path](https://github.com/mariocasciaro/object-path) module.
 
 ### `depGraph` <a name="depGraph"></a>

 See the documentation for the [dep-graph.js](https://github.com/TrevorBurnham/dep-graph) module.
 
 ### `resourceMap(resources) {...}` <a name="resourceMap"></a>

 Helper function.  Returns a map of resources by their names.
 
 ### `become(command, become, becomeUser, becomeMethod) {...}` <a name="become"></a>

 Helper function.  Returns a string with the command wrapped in a privilege escalation.
 
 ### `modules.handlers.register(name, cb) {...}` <a name="modules.handlers.register"></a>

 Register a resource handler function.
 
 A resource handler function takes three arguments:
 
 > * `catalog`: the current value of the AWS resources found by `mithras.run()`, possibly modified by other resource execution.
 > * `resources`: the set of resources passed into `mithras.apply()`
 > * `targetResource`: the resource object being evaluated for execution.
 
 Handler functions return a Javascript array.  The second element is
 a boolean. If `true`, it indicates that the handler function "owns"
 the resource and has handled it.  If so, and the first element of
 the return value is defined, the value of the first element is set
 as the `_target` property of the resource object being evaluated.
 
 See [Design and Concepts](design.html) for a more detailed
 explanation of how this all works.
 
 ### `modules.handlers.run(catalog, resources, targetResource, dict) {...}` <a name="modules.handlers.run"></a>

 Run a resource through handler functions.

 In turn, every registered handler function is called with three
 arguments: `catalog`, `resources` and `targetResource`.  The
 catalog is the current value of the AWS resources gathered by
 `mithras.run()`, which may have been modified by previous resource
 execution.  The `resources` argument is the set of resources that
 was passed into `mithras.apply()`.  The final argument is the
 resource object currently being evaluated.

 Handler functions return a Javascript array.  The second element is
 a boolean. If `true`, it indicates that the handler function "owns"
 the resource and has handled it.  If so, and the first element of
 the return value is defined, the value of the first element is set
 as the `_target` property of the resource object being evaluated.
 
 See [Design and Concepts](design.html) for a more detailed
 explanation of how this all works.
 
 ### `modules.preflight.register(name, cb) {...}` <a name="modules.preflight.register"></a>

 Register a resource handler preflight function.
 
 A resource handler preflight function takes three arguments:
 
 > * `catalog`: the current value of the AWS resources found by `mithras.run()`, possibly modified by other resource execution.
 > * `resources`: the set of resources passed into `mithras.apply()`
 > * `targetResource`: the resource object being evaluated for preflight.
 
 Preflight functions return a Javascript array.  The second element
 is a boolean. If `true`, it indicates that the function "owns" the
 resource and has handled it.  If so, and the first element of the
 return value is defined, the value of the first element is set as
 the `_target` property of the resource object being evaluated.
 
 See [Design and Concepts](design.html) for a more detailed
 explanation of how this all works.
 
 ### `modules.preflight.run(catalog, resources, order) {...}` <a name="modules.preflight.run"></a>

 Run preflight functions on resources.
 
 For each resource, first the resource is recursively traversed by
 [`updateResource`](#updateResource), and given an opportunity to
 access runtime values in other resources.
 
 Next, every registered preflight function is called with three
 arguments: `catalog`, `resources` and `targetResource`.  The
 catalog is the current value of the AWS resources gathered by
 `mithras.run()`, which may have been modified by previous resource
 execution.  The `resources` argument is the set of resources that
 was passed into `mithras.apply()`.  The final argument is the
 resource object currently being evaluated.

 Handler functions return a Javascript array.  The second element is
 a boolean. If `true`, it indicates that the handler function "owns"
 the resource and has succesfully preflighted it.  If so, and the
 first element of the return value is defined, the value of the
 first element is set as the `_target` property of the resource
 object being evaluated.
 
 See [Design and Concepts](design.html) for a more detailed
 explanation of how this all works.
 
 ### `MODULES` <a name="modules"></a>

 A map of loaded core module names to their version strings.
 
 ### `GOPATH` <a name="gopath"></a>

 The current GOPATH when mithras is invoked (if set).
 
 ### `HOME` <a name="home"></a>

 The value of $MITHRASHOME or -m if set on the command line.
 
 ### `CWD` <a name="cwd"></a>

 The current working directory of Mithras.
 
 ### `ARGS` <a name="args"></a>

 When `mithras run` is invoked at the command line, any additional
 non-flag parameters supplied on the command line are passed through
 to the user script in this array.
 
 ### `verbose` <a name="verbose"></a>

 Set to `true` if the `-v` global flag is used to invoke Mithras on
 the command line.  Eg., `mithras -v run ...`
 

 
 <a name="bootstrap"></a>
 
 ### `bootstrap(template) {...}`

 Returns a single resource object, which encapsulates all of
 the resources need to bootstrap an instance for use with
 Mithras.

 Example template supplied by caller:
 ```
 { 
   dependsOn: ["webserver"]
   params: {
         become: true
         becomeMethod: "sudo"
         becomUser: "root"
         hosts: mithras.watch("webserver._target")
  }
 }
 ```

 
 <a name="apply"></a>
 
 ### `apply(catalog, resources, reverse) {...}`

 The "core" function of Mithras.  Given a `catalog`, an
 array of resource objects in `resources`, and a boolean
 (`reverse`), apply the resources to the catalog.

 First a dependency graph is built.  In forward order, all
 resources are preflighted.

 Next, in the desired order (reversed, if `reverse` is
 `true`), the resources are run through their handlers in
 dependency order.

 The catalog, after update by handlers, is returned.

 See [Design and Concepts](design.html) for a more detailed
 explanation of how this all works.
 

 
 
 ### `buildDeps(resources, add_node, add_dep) {...}`

 Helper function.  Loop through resources, calling the
 `add_node` and `add_dep` functions supplied by the caller
 to create adependency graph.


 
 
 ### `updateResource(resource, catalog, resources, name) {...}`

 Given a `resource`, the `catalog` and a list of all
 `resources`, update the resource allowing it to reach into
 the catalog and/or other resources to set its properties.

 Returns a COPY of the resource with updated fields.


 
 <a name="doIncludes"></a>
 
 ### `doIncludes(resources) {...}`

 Recursively descend through supplied `resources`, adding
 their dependencies via their `includes` property.


 
 <a name="sshKeyPathForInstance"></a>
 
 ### `sshKeyPathForInstance(resource, instance) {}`

 Given a `resource` and an ec2 `instance` object, return the
 appropriate path to the SSH key for the instance.

 If the resource has a property named
 `sshKeyPathForInstance`, it is invoked and its return value
 used.

 The default return value is:

 `"~/.ssh/" + instance.KeyName + ".pem"`


 
 <a name="sshUserForInstance"></a>
 
 ### `sshUserForInstance(resource, instance) {}`

 Given a `resource` and an ec2 `instance` object, return the
 appropriate SSH username for the instance.

 If the resource has a property named
 `sshUserForInstance`, it is invoked and its return value
 used.

 The default return value is:

 `"ec2-user"`


 
 <a name="watch"></a>
 
 ### `watch(path, cb) {...}`

 A resource may set any property value to a function,
 instead of a string, array, etc.  When the resource is
 preflight'ed, that function will be called with two
 arguments, the current `catalog`, and an array of
 `resources`.  The parameter function may return `undefined`,
 and if so, it will remain a function.  If it returns any
 other value, the value of the property to which it is
 attached is retplaced with the parameter function's return
 value.

 This allows parameter properties to be evaluated at
 runtime, not just when the resource is defined.  The use
 case is appropriate to AWS, when a given resource needs the
 value from some other resource in order to be handled.  For
 example, instances may be placed into subnets, which are
 defined on some other resource.  This gives the target
 resource the ability to reach into resources it depends on
 and extract values for use in parameters.

 All of this functionality is wrapped up in a neat little
 package here, in `mithras.wrapper`.  Using this function,
 one resource can examine another for its runtime
 properties.  Here's an example:

 ```
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
 ```

 This example resource defines a subnet.  Subnets are
 created in VPCs.  Here, the subnet property `VpcId` uses
 `mithras.watch` at runtime to obtain the id of the VPC on
 which it depends, since the id is not known when the
 resource is defined.

 Supply two parameters to `watch`: a path (expressed in
 [object-path](https://github.com/mariocasciaro/object-path)
 form), and an optional callback.  If the object path is
 defined, the callback will be invoked `cb(catalog,
 resources, value_at_path)` The return value of the callback
 will be returned by the watch function during preflight,
 and follows the same rules outlined above.

 See [Design and Concepts](design.html) for a more detailed
 explanation of how the `_target` property gets set on
 resource objects.
 

 
 <a name="findGWByVpcId"></a>
 
 ### `findGWByVpcId(cat, resources, vpcId) {...}`

 Given a `vpcId`, look through the `catalog` and find a
 matching internet gateway.  If one is found, return its
 `InternetGatewayId` property.


 
 <a name="run"></a>
 
 ### `run() {...}`

 Called by user scripts to interrogate AWS and return a
 `catalog` of resources.


 
 <a name="active"></a>
 ### `activeRegions(catalog) {...}`

 Returns an array of AWS regions.  User scripts may replace
 this function with their own to limit the scope of queries
 that `mithras.run()` will execute in looking for resources
 on AWS.

 If not replaced, `mithras.activeRegions` will return all
 regions.


