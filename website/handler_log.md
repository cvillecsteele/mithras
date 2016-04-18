 
 
 # log
 
 Log is resource handler for writing messages to the mithras log.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"log"`
 
 Usage:
 
 `var log = require("log").init();`
 
  ## Example Resource
 
 ```javascript
 var rLog = {
      name: "log"
      module: "log"
      params: {
          messsage:  "Hello, world!"
      }
 };
 ```
 
 ## Parameter Properties
 
 ### `message`

 * Required: true
 * Allowed Values: any string

 The value of this property will be written to the Mithras log.
 

