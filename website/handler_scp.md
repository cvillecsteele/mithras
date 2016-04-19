 
 
 # scp
 
 SCP is a resource handler for moving files between systems.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"scp"`
 
 Usage:
 
 `var scp = require("scp").init();`
 
  ## Example Resource
 
 ```javascript
 var testFile = {
     name: "test"
     module: "scp"
     dependsOn: [other.name]
     params: {
         region: defaultRegion
         ensure: ensure
         src: "/etc/hosts"
         dest: "/tmp/foo"
         hosts: [array of ec2 objects]
     }
 }
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"` and the file at `dest` does not exist, the file at
 `src`, locally, is copied to the remote system.  If `"absent"`, and
 the file at `dest` on the remore system exists, it is removed.
 
 ### `src`

 * Required: false
 * Allowed Values: a path to a file on the local system.

 This file is copied to remote hosts.

 ### `dest`

 * Required: false
 * Allowed Values: a path to a file on the remote system.

 The file from `src` on the local host is copied to this path on
 remote hosts.

 ### `hosts`

 * Required: false
 * Allowed Values: an array of ec2 instance objects

 This property specifies the hosts on which this resource is to be applied.


