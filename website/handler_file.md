 
 
 # File
 
 File is resource handler for manipulating files.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"file"`
 
 Usage:
 
 `var file = require("file").init();`
 
  ## Example Resource
 
 ```javascript
 var rFile = {
   name: "someFile"
   module: "file"
   dependsOn: [otherResource.name]
   params: {
     dest: "/etc/foo/bar"
     ensure: "directory"
     mode: 0777
     hosts: [<array of ec2 instance objects>]
   }
 };
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "absent", "file", "directory", "link", "hard", "touch"

 If `"directory"`, all immediate subdirectories will be created if
 they do not exist. If `"file"`, the file will NOT be created if it
 does not exist, see the copy or template module if you want that
 behavior. If `"link"`, the symbolic link will be created or
 changed. Use `"hard"` for hardlinks. If `"absent"`, directories
 will be recursively deleted, and files or symlinks will be
 unlinked. If `"touch"`, an empty file will be created if the path
 does not exist, while an existing file or directory will receive
 updated file access and modification times (similar to the way
 `touch` works from the command line).
 
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

 ### `mode`

 * Required: true
 * Allowed Values: octal number specifying a valid permission mask

 This property specifies the path to the file/link/directory to be manipulated

 ### `chown`

 * Required: false
 * Allowed Values: username of the user to which the file will be `chown`'ed

 This property specifies the path to the file/link/directory to be manipulated

 ### `dest`

 * Required: true
 * Allowed Values: a valid path on the target host

 This property specifies the path to the file/link/directory to be manipulated

 ### `src`

 * Required: false
 * Allowed Values: a valid path on the target host

 The path of the file to link to (applies only to
 ensure=`"link"`). Will accept absolute, relative and nonexisting
 paths. Relative paths are not expanded.


