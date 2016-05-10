 


 # CORE FUNCTIONS: WEB


 

 This package exports several entry points into the JS environment,
 including:

 > * [web.run](#run)
 > * [web.stop](#stop)
 > * [web.get](#get)
 > * [web.handler](#handler)
 > * [web.url.parse](#uparse)

 This API allows JS to fetch from the web and to create a web server.

 ## WEB.RUN
 <a name="run"></a>
 `web.run(addr);`

 Run a web server.

 Example:

 ```

 var server = web.run(":http");

 ```

 ## WEB.STOP
 <a name="stop"></a>
 `web.stop(server);`

 Stop a web server.

 Example:

 ```

 var server = web.run(":http");
 web.stop(server);

 ```

 ## WEB.GET
 <a name="get"></a>
 `web.get(url);`

 Fetch an URL and return its contents.

 Example:

 ```

 var html = web.get("http://www.cnn.com");

 // To write the contents to a file:

 web.get("http://www.cnn.com", "/tmp/cnn", 0644);

 ```

 ## WEB.URL.PARSE
 <a name="uparse"></a>
 `web.url.parse(url);`

 Parse and url and return its component parts.

 Example:

 ```

 var url = web.url.parse("http://www.cnn.com");

 ```


