 


 # CORE FUNCTIONS: WORKERS (GOROUTINES)


 

 This package exports entry points into the JS environment:

 > * [workers.run](#run)
 > * [workers.stop](#stop)
 > * [workers.send](#send)
 > * [workers.receive](#receive)

 This API allows the caller to work with workers

 ## WORKERS.RUN
 <a name="run"></a>
 `workers.run(name, source);`

 Create a new worker identified by `name` and with a handler
 function with the JS source of `source` (a string, not a function).

 Example:

 ```
  function foo(input) { return "hi!"; }
  workers.run("test", foo.toString());

 ```

 ## WORKERS.SEND
 <a name="send"></a>
 `workers.send(name, value);`

 Send a (string) value to the worker named `name`.

 Example:

 ```
  workers.send("test", JSON.stringify({"a": 42}));

 ```

 ## WORKERS.RECEIVE
 <a name="receive"></a>
 `workers.receive(name);

 Read the output of a worker.

 Example:

 ```
  var out = JSON.parse(workers.receive("test"));

 ```

 ## WORKERS.STOP
 <a name="stop"></a>
 `workers.stop(name);

 Shut down the worker and remove it.

 Example:

 ```
  workers.stop("name");

 ```


