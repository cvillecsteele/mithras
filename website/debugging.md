# Debugging Mithras

## Crank up the log

First run Mithras in verbose mode:

    mithras -v run -f site.js

Also, to see debug level messages:

    mithras.log.setLevel("debug");

## Use `on_handle`

One technique is to set an `on_handle` property on a resource you wish
to debug.  The value of the property should be a function.  For
example:

    on_handle: function (catalog, 
                         resources, 
                         targetResource, 
                         f) {
        console.log(JSON.stringify(targetResource, null, 2));
        reutrn f(catalog, resources, targetResource);
    }

The `on_handle` function is called when the resource is being handled,
and it gives you an opportunity to inspect and/or modify a resource
before the its handler runs.

## Use wrappers

For example:

    mithras.handlers.wrap("service", function(catalog, resources, target, f) {
        console.log(JSON.stringify(targetResource, null, 2));
        return f(catalog, resources, target);
    });

This approach lets you insert your code before _every_ resource that
is executed by a given handler.
