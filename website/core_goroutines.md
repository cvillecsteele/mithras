 


 # CORE FUNCTIONS: GOROUTINES (GO)


 

 This package exports entry points into the JS environment:

 > * [go.run](#run)

 This API allows the caller to work with goroutines

 ## GOROUTINES.RUN
 <a name="run"></a>
 `go.run(f);`

 Run the function `f` in a goroutine.

 Example:

 ```

  go.run(function() { console.log("hello from a goroutine"); });

 ```


