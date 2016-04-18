 
 
 # Copy
 
 Copy is resource handler for creating and copying files.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"copy"`
 
 Usage:
 
 `var copy = require("copy").init();`
 
  ## Example Resource
 
 ```javascript
 var rFile = {
   name: "someFile"
   module: "copy"
   dependsOn: [otherResource.name]
   params: {
     ensure: "present"          // "present" or "absent"
     become: true               // priv escalation
     becomeMethod: "sudo"       // escalation method
     becomeUser: "root"         // desired user
     dest: "/tmp/foo"           // destination file
     src: "/etc/hosts"          // source file
     mode: 0644                 // permissions for destination file
     hosts: [<array of ec2 instance objects>]
   }
 };
 ```
 
 ## Copy Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"`, the file specified in `dest` will be created,
 either with the contents of the file at `src`, or the file contents
 specified in `content`.  If `"absent"`, the file at `dest` will be
 removed.
 
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

 ### `dest`

 * Required: true
 * Allowed Values: a valid path on the target host

 This property specifies the path to the file to be copied into.

 ### `src`

 * Required: false
 * Allowed Values: a valid path on the target host

 One of `src` or `content` must be specified.  If `src` is set, it
 is the path to the file whose contents are to be copied into
 `dest`.

 ### `content`

 * Required: false
 * Allowed Values: string

 One of `src` or `content` must be specified.  If `content` is set, it
 is the path to the file whose contents are to be copied into
 `dest`.

 ### `mode`

 * Required: true
 * Allowed Values: octal number specifying a valid permission mask


