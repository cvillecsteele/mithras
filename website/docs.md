## General

See [Installation & Usage](usage.html) to get Mithras installed on
your system.
 
Then, jump into [Quickstart](quickstart.html) to get up and running.

For a more in-depth walkthrough, see [The Guide](guide1.html)

Read the API [Reference](api.html) to learn more.

Or, check out the [FAQ](faq.html).

## Running Mithras

You can run Mithras with a variety of sub-commands, including:

* [build](#build)
* [run](#run)
* [daemon](#daemon)
* [repl](#repl)
* [get](#get)
* [help](#help)

### Build <a name="build"></a>

Most Mithras [handlers](api.html) require a copy of Mithras itself
to be present on the target instance.  

(In your script, you'll use `mithras.bootstrap(...)`, which is documented
[here](handler_mithras.html#bootstrap).)

Using the `build` command, you create local copies of Mithras suitable
for remote instances.

Example usage:

    mithras build

Or if you want to specify architectures and OSs:

    mithras build -arch amd64 -arch 386 -os linux

### Run <a name="run"></a>

Run a Mithras script.

Example usage:

    mithras -v run -f example/simple.js up

Make sure that you've followed [Installation & Usage](usage.html) to
get Mithras property set up to run, first.

### Daemon <a name="daemon"></a>

Run a Mithras script in the background.

Example usage:

    mithras -v daemon -f example/asg_daemon.js

Make sure that you've followed [Installation & Usage](usage.html) to
get Mithras property set up to run, first.

### Repl <a name="repl"></a>

Sometimes you need to mess around with some Javascript.  The Mithras
`repl` command starts up a JS interpreter for you, and away you go:

    mithras repl

### Get <a name="get"></a>

Install a JS package from the NPM registry.

Example usage:

    mithras get zipcodes-regex

Note that just because Mithras can _install_ some random package from
the NPM registry, definitely does NOT mean that it will run in
Mithras.  Mithras isn't NodeJS.

### Help <a name="help"></a>

See Mithras command line flags, etc.

Example usage:

    mithras help

You'll get something like this:

    NAME:
       mithras - Manage resources in AWS

    USAGE:
       mithras [global options] command [command options] [arguments...]

    VERSION:
       0.1.0

    COMMANDS:
       run, r                Run a mithras script
       build, b              Build Mithras remote helper binaries.  Run this first.
       get, install, g       This command installs a package, and any packages that it depends on.
       repl                  Run a Mithras JS repl
       daemon, d, demon      Run a Mithras daemon
       help, h               Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --mithras, -m        Mithras home directory [$MITHRASHOME]
       --verbose, -v        Verbose output
       --help, -h           show help
       --version, -V        print the version
