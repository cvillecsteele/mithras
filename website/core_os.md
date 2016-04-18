


# CORE FUNCTIONS: OS




This package exports several entry points into the JS environment,
including:

> * [os.exit](#exit)
> * [os.hostname](#hostname)
> * [os.getenv](#getenv)
> * [os.expandEnv](#expand)

This API allows resource handlers to execute various OS-related functions.

## OS.EXIT
<a name="exit"></a>
`os.exit(status);`

Terminate the program, returning the specificed status code.

Example:

```

os.exit(1);

```

## OS.EXIT
<a name="hostname"></a>
`os.hostname();`

Get the hostname.

Example:

```

var hostname = os.hostname();

```

## OS.GETENV
<a name="getenv"></a>
`os.getenv(key);`

Getenv retrieves the value of the environment variable named by the
key. It returns the value, which will be empty if the variable is
not present.

Example:

```

var home = os.getenv("HOME");

```

## OS.EXPANDENV
<a name="expand"></a>
`os.expandEnv(target);`

ExpandEnv replaces ${var} or $var in the string according to the
values of the current environment variables. References to
undefined variables are replaced by the empty string.

Example:

```

var where = os.getenv("$HOME/.ssh");

```


