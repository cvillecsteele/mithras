 
 
 # Cache
 
 This module exports:
 
 > * `Cache` constructor function
 
 Usage:
 
 `var Cache = (new (require("cache").Cache)).init();`
 

 
 
 ## `Cache`
 
 Constructor function
 
 `new Cache(path)`
 
 Example:
 ```javascript
 var cache = new Cache(".cache");
 ```
 

 
 
 ## `init`
 
 Initialize a cache
 
 `cache.init()`
 
 Example:
 ```javascript
 var cache = new Cache(".cache");
 cache.init();
 ```
 

 
 
 ## `get`
 
 Get an object from cache.
 
 `cache.get(key)`
 
 Example:
 ```javascript
 var cache = new Cache(".cache");
 cache.init();
 var foo = cache.get("foo");
 ```
 

 
 
 ## `put`
 
 Put an object into the cache.  The `value` arg is run
 through `JSON.stringify` before storage in the cache.
 
 `cache.put(key, value, expiry)`
 
 Example:
 ```javascript
 var cache = new Cache(".cache");
 cache.init();
 cache.put("foo", "somevalue", (60 * 5));
 ```
 

