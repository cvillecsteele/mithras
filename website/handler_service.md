 
 
 # service
 
 Service is a resource handler for dealing with services running on instances.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"service"`
 
 Usage:
 
 `var service = require("service").init();`
 
  ## Example Resource
 
 ```javascript
 var svc = {
     name: "myservice"
     module: "service"
     params: {
       ensure: "present" // or "absent"
       name: "nginx"
     }
 };
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present" or "absent"

 If `"present"` and the service `name` is not running, it is started.
 If `"absent"`, and the service is running, it is stopped.
 
 ### `name`

 * Required: true
 * Allowed Values: a string specifying a configured service on the instance.

 This is the service that is started/stopped.


