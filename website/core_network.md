 


 # CORE FUNCTIONS: NETWORK


 

 This package exports several entry points into the JS environment,
 including:

 > * [network.check](#check)

 This API allows resource handlers to open TCP connections.

 ## NETWORK.check
 <a name="check"></a>
 `network.check(host, port, timeout);`

 Returns true if a TCP connection can be established.

 Example:

 ```

  var ok = network.check("10.0.22.33", 22, 100);

 ```


