 
 
 # packager
 
 Packager is resource handler for working with packager connections.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"packager"`
 
 Usage:
 
 `var packager = require("packager").init();`
 
  ## Example Resource
 
 ```javascript
 var rUpdatePkgs = {
     name: "updatePackages"
     module: "packager"
     dependsOn: [otherResource.name]
     params: {
         ensure: "latest"
         name: ""
         become: true
         becomeMethod: "sudo"
         becomeUser: "root"
         hosts: [ec2 instances objects]
     }
 };
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present", "absent", "latest"

 If `"present"`, the package will be installed if is not already.
 If `"absent"`, the package will be removed if it is present.  If
 `"latest"`, the package will be installed if not present, and
 updated if it is already on the remote system.
 
 ### `name`

 * Required: true
 * Allowed Values: any valid package name; eg "nginx"

 If set to the empty string, `""`, this resource handler will omit
 the package name from the command line invocation of `yum`.  Use
 `""` in conjunction with `ensure: "latest"` to upgrade all packages
 on the remote system.
 
 ### `become`

 * Required: false
 * Allowed Values: true or false

 If `true`, the copy will attempt to run with escalated privs, as
 specified in the properties `becomeMethod` and `becomeUser`.
 
 ### `becomeMethod`

 * Required: false
 * Allowed Values: "su" or "sudo"

 The method of privilege escalation.
 
 ### `becomeUser`

 * Required: false
 * Allowed Values: any string specifying a username suitable for use by `becomeMethod`

 ### `hosts`

 * Required: true
 * Allowed Values: an array of ec2 instance objects

 This property specifies the hosts on which this resource is to be applied.


