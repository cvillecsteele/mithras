 
 
 # mithras
 
 Mithras is a core object for building scripts to manage your AWS
 stacks.  The `mithras` object is loaded prior to loading your
 scripts.
 
 The `mithras` object has the following properties:

 > * [remote](#remote)
 > * [traverse](#traverse)
 > * [objectPath](#objectPath)
 > * [depGraph](#depGraph)
 > * [resourceMap](#resourceMap)
 > * [become](#become)
 > * [modules.handlers.register](#modules.handlers.register)
 > * [modules.handlers.run](#modules.handlers.run)
 > * [modules.preflight.register](#modules.preflight.register)
 > * [modules.preflight.run](#modules.preflight.run)
 > * [bootstrap](#bootstrap)
 > * [apply](#apply)
 > * [buildDeps](#buildDeps)
 > * [updateResource](#updateResource)
 > * [doIncludes](#doIncludes)
 > * [sshKeyPathForInstance](#sshKeyPathForInstance)
 > * [sshUserForInstance](#sshUserForInstance)
 > * [watch](#watch)
 > * [findGWByVpcId](#findGWByVpcId)
 > * [run](#run) 
 > * [activeRegions](#active)
 > * [MODULES](#modules)
 > * [VERSION](#version)
 > * [GOPATH](#gopath)
 > * [ARGS](#args)
 > * [verbose](#verbose)
 
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
 
 ### `modules.handlers.run(catalog, resources, targetResource, dict) {...}` <a name="modules.handlers.run"></a>

 Run a resource through handler functions.
 
 ### `modules.preflight.register(name, cb) {...}` <a name="modules.preflight.register"></a>

 Register a resource handler preflight function.
 
 ### `modules.preflight.run(catalog, resources, order) {...}` <a name="modules.preflight.run"></a>

 Run preflight functions on resources.
 
 ### `MODULES` <a name="modules"></a>

 A map of loaded core module names to their version strings.
 
 ### `GOPATH` <a name="gopath"></a>

 The current GOPATH when mithras is invoked (if set).
 
 ### `ARGS` <a name="args"></a>

 When `mithras run` is invoked at the command line, any additional
 non-flag parameters supplied on the command line are passed through
 to the user script in this array.
 
 ### `verbose` <a name="verbose"></a>

 Set to `true` if the `-v` global flag is used to invoke Mithras on
 the command line.  Eg., `mithras -v run ...`
 

