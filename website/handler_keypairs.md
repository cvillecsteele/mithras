 
 
 # keypairs
 
 Keypairs is resource handler for manipulating AWS SSH keypairs.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"keypairs"`
 
 Usage:
 
 `var keypairs = require("keypairs").init();`
 
  ## Example Resource
 
 ```javascript
 var rKey = {
      name: "key"
      module: "keypairs"
      skip: (ensure === 'absent') // Don't delete keys
      params: {
          region: "us-east-1"
          ensure: "present"
          key: {
              KeyName: "my-fancy-key"
          }
          savePath: os.expandEnv("$HOME/.ssh/" + keyName + ".pem")
      }
 };
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "absent", "present"

 If `"present"`, the handler will ensure the keypair exists, and it
 not, it will be created.  If `"absent"`, the keypair is removed.
 
 ### `key`

 * Required: true
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2.html#type-CreateKeyPairInput)

 Specifies parameters for keypair creation.

 ### `savepath`

 * Required: true
 * Allowed Values: A valid path for saving the pemfile when a keypair is created

 When the handler creates a new keypair, the contents of the key are saved to this path.


