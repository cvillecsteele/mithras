


# CORE FUNCTIONS: EXEC




This package exports one entry point into the JS environment:

> * [exec.run](#run)

This API allows the caller to exec and run a program, collecting its output.

## EXEC.RUN
<a name="run"></a>
`exec.run(command);`

Exec and run a program, collecting its output.

Example:

```

var results = exec.run("pwd");
var out = results[0];
var err = results[1];
var ok = results[2];
var status = results[3];

```


