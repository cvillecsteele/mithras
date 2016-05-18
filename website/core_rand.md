 


 # CORE FUNCTIONS: RAND


 

 This package exports several entry points into the JS environment,
 including:

 > * [rand](#rand)

 This API generates random numbers

 ## LOG
 <a name="rand"></a>
 `rand(n);`

 Intn returns, as an int, a non-negative pseudo-random number in
 [0,n) from the default Source. It panics if n <= 0.

 Example:

 ```

  var x = rand(42);

 ```


