# WALKTHROUGH, PART 2

Use this document to get up and working quickly and easily with
Mithras.

* [Part One](quickstart1.html): An EC2 instance
* [Part Two](quickstart2.html): Configuring our instance
* [Part Three](quickstart3.html): A complete application stack

## Part Two: Configuring our EC2 Instance

This part of the guide will show you how to use Mithras to configure
various resources on an EC2 instance, such as package upgrades and
installation, service initialization, and more. We'll also dig a
little deeper into key concepts used by Mithras.

Before you get going, make sure that you've [installed](usage.html)
Mithras first.  Also, double check that your AWS credentials are set
up correctly and that you've run the first part of this script. (See
[Part One](quickstart1.html).)

### A Security Group

### Coordinating Resources

### Keypair resource

### Instance resource

### Applying resources to the catalog

Phew!  That's it for resource definition.  Now that our script has set
things up with resources, it tells Mithras to do the work of applying
those resources to the current catalog of existing AWS resources:

```
mithras.apply(catalog, [  ], reverse);
```

This code tells mithras to `apply` the resources in the second
argument, `[rKey, rInstance]`, to the `catalog` argument.  The final
argument is a boolean and if `true`, it tells Mithras that the order
of dependencies is reversed, which is appropriate for *deleting* AWS
resources, which is done in the opposite order as *creating* them.

### Running the script

Last but not least, here's how you tell Mithras to run this script.
Make sure you've set everything up according to the
[usage](usage.html) instructions, first.  Then, in your terminal, run:

    mithras -v run -f example/intermediate.js

Since we specified a global Mithras CLI option of `-v`, we see some
pretty verbose output about what Mirthas does.






