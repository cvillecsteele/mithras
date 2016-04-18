 
 
 # network
 
 Network is resource handler for working with network connections.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"network"`
 
 Usage:
 
 `var network = require("network").init();`
 
  ## Example Resource
 
 ```javascript
 var available = {
   name: "sshAvailable"
   module: "network"
   dependsOn: [otherResource.name]
   params: {
     timeout: 120
     port: 22
     hosts: [<array of ec2 instance objects>]
   }
 };
 ```
 
 ## Parameter Properties
 
 ### `port`

 * Required: false
 * Allowed Values: integer

 The port on the remote host to attempt to connect to.  Defaults to
 22.
 
 ### `timeout`

 * Required: false
 * Allowed Values: integer

 A number of seconds to attempt to connect to the remote host.
 Defaults to 120.
 

