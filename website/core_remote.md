

# CORE FUNCTIONS: REMOTE




This package exports several entry points into the JS environment,
including:

> * [mithras.remote.scp](#scp)
> * [mithras.remote.shell](#shell)
> * [mithras.remote.wrapper](#wrapper)
> * [mithras.remote.mithras](#mithras)

This API allows resource handlers to execute tasks on remote hosts
in a variety of ways.

## MITHRAS.REMOTE.SCP
<a name="scp"></a>
`mithras.remote.scp(ip, user, keypath, src, dest);`

Copy a file to a remote host.

Example:

```

mithras.remote.scp("52.90.244.101",
"ec2-user",
"/home/user/.ssh/key.pem",
"/tmp/sourcefile",
"/etc/hosts");
```

## MITHRAS.REMOTE.SHELL
<a name="shell"></a>
`mithras.remote.shell(ip, user, keypath, input, cmd, env);`

Execute command(s) in a shell on a remote system.  The arg `env`
specifies an object mapping environment variables to values for the
*local* execution of the ssh command.  If the `input` arg is a
string, the stdin of the locally-executed ssh command will be set
to the contents of the argument and the locally executed ssh
command will not use the `-tt` command line option for setting a
tty.

Example:

```

mithras.remote.shell("52.90.244.101",
"ec2-user",
"/home/user/.ssh/key.pem",
"hello, world!\n",
"cat > /tmp/foo",
{"envVar": "value"});

```

## MITHRAS.REMOTE.WRAPPER
<a name="wrapper"></a>
`mithras.remote.wrapper(ip, user, keypath, args, env);`

Execute a single command in a shell on a remote system.  The arg `env`
specifies an object mapping environment variables to values for the
*remote* execution of the caller-supplied command.

Example:

```

mithras.remote.wrapper("52.90.244.101",
"ec2-user",
"/home/user/.ssh/key.pem",
["ls", "-l"],
{"envVar": "value"});

```
## MITHRAS.REMOTE.MITHRAS

<a name="mithras"></a>
`mithras.remote.mithras(instance, user, keypath, js, become, becomeUser, becomeMethod);`


Example:

```

mithras.remote.wrapper(<ec2 instance object>
"ec2-user",
"/home/user/.ssh/key.pem",
"(function run() { console.log('hi'); })"
true,
"root"
"sudo");

```


