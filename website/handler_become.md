 
 
 # BECOME
 
 This module exports:
 
 > * `become`
 
 Usage:
 
 `var become = require("become").become;`
 

 
 
 ## `become`
 
 Returns a command suitable for running with escalated privs.
 
 `become(become, becomeUser, becomeMethod, command);`
 
 Example:
 ```javascript
 var cmd = become(true, "root", "sudo", "ls /var/log");
 ```
 

