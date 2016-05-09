 
 
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
     src: scp://localhost/file.txt
     ensure: "file"
     mode: 0644
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

 ### `owner`

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

 If ensure=`"file"` or `"directory"`, the value of thie property may
 take one of three forms.  If of the form
 `"scp://localhost/foo/bar"`, then the *local* file specified by the
 `src` is SCP'd to the remote host, to the value of `dest`.  If of
 the form `"http://www.someplace.com/foo/bar"`, then from the remote
 instance, an HTTP GET request is performed to the value of `src`,
 and the contents of the response are written to `dest`.

 If ensure=`"link"`, specifies the path of the file to link to.
 Will accept absolute, relative and nonexisting paths. Relative
 paths are not expanded.

 ### `content`

 * Required: false
 * Allowed Values: a string of file contents to be written

 If ensure=`"file"`, the value of this property (presumably a
 string) will be written to `dest`.

 ### `force`

 * Required: false
 * Allowed Values: boolean

 If `true`, any `file` will be overwritten, even if it already exists.


