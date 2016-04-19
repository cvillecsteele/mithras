
 # RESOURDCE HANDLER: SHELL
 
 
 
 # shell
 
 Shell is a resource handler for dealing with shells running on instances.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"shell"`
 
 Usage:
 
 `var shell = require("shell").init();`
 
  ## Example Resource
 
 ```javascript
 var foo = {
     name: "dostuff"
     module: "shell"
     params: {
       command: "whoami"
       hosts: [ec2 instance objects...]
     }
 };
 ```
 
 ## Parameter Properties
 
 ### `command`

 * Required: true
 * Allowed Values: a valid shell command for the remote instance

 The command in this property is run on remote instances.
 
 ### `hosts`

 * Required: false
 * Allowed Values: an array of ec2 instance objects

 The command is executed on these instances.


