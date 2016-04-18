


# CORE FUNCTIONS: WEB




This package exports several entry points into the JS environment,
including:

> * [user.run](#run)
> * [user.stop](#stop)
> * [user.get](#get)
> * [user.handler](#handler)

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

```


