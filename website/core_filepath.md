


# CORE FUNCTIONS: FILEPATH




This package exports entry points into the JS environment:

> * [filepath.dir](#dir)
> * [filepath.base](#base)
> * [filepath.ext](#ext)
> * [filepath.glob](#glob)
> * [filepath.match](#match)
> * [filepath.split](#split)
> * [filepath.splitList](#splitlist)
> * [filepath.rel](#rel)
> * [filepath.clean](#clean)
> * [filepath.abs](#abs)
> * [filepath.join](#join)
> * [filepath.walk](#walk)

This API allows the caller to work with filesystem paths.

## FILEPATH.DIR
<a name="dir"></a>
`filepath.dir(path);`

Dir returns all but the last element of path, typically the path's
directory. After dropping the final element, the path is Cleaned
and trailing slashes are removed. If the path is empty, Dir returns
".". If the path consists entirely of separators, Dir returns a
single separator. The returned path does not end in a separator
unless it is the root directory.

Example:

```

var dir = filepath.dir("/tmp/foo");

```

## FILEPATH.BASE
<a name="base"></a>
`filepath.base(path);`

Base returns the last element of path. Trailing path separators are
removed before extracting the last element. If the path is empty,
Base returns ".". If the path consists entirely of separators, Base
returns a single separator.

Example:

```

var base = filepath.base("/tmp/foo");

```

## FILEPATH.EXT
<a name="ext"></a>
`filepath.ext(path);`

Ext returns the file name extension used by path. The extension is
the suffix beginning at the final dot in the final element of path;
it is empty if there is no dot.

Example:

```

var ext = filepath.ext("/tmp/foo.jpg");

```

## FILEPATH.GLOB
<a name="glob"></a>
`filepath.glob(pattern);`

Glob returns the names of all files matching pattern or nil if
there is no matching file. The syntax of patterns is the same as in
Match. The pattern may describe hierarchical names such as
/usr/*/bin/ed (assuming the Separator is '/').

Glob ignores file system errors such as I/O errors reading
directories. The only possible returned error is ErrBadPattern,
when pattern is malformed.

Example:

```

var results = filepath.glob("/tmp/*.jpg");
var matches = results[0];
var error = results[1];

```

## FILEPATH.MATCH
<a name="match"></a>
`filepath.match(pattern, name);`

Match reports whether name matches the shell file name pattern.
See [here](https://golang.org/pkg/path/filepath/#Dir) for syntax.

Example:

```

var results = filepath.match("/tmp/*.jpg", "/var/other");
var matched = results[0];
var error = results[1];

```

## FILEPATH.split
<a name="split"></a>
`filepath.split(path);`

Split splits path immediately following the final Separator,
separating it into a directory and file name component. If there is
no Separator in path, Split returns an empty dir and file set to
path. The returned values have the property that path = dir+file.

Example:

```

var results = filepath.split("/tmp/foo.jpg");
var dir = results[0];
var file = results[1];

```

## FILEPATH.splitList
<a name="splitlist"></a>
`filepath.splitList(path);`

SplitList splits a list of paths joined by the OS-specific
ListSeparator, usually found in PATH or GOPATH environment
variables. Unlike strings.Split, SplitList returns an empty slice
when passed an empty string. SplitList does not replace slash
characters in the returned paths.

Example:

```

var parts = filepath.splitList("/tmp/bar/foo.jpg");

```

## FILEPATH.rel
<a name="rel"></a>
`filepath.rel(basePath, targetPath);`

Rel returns a relative path that is lexically equivalent to
targpath when joined to basepath with an intervening
separator. That is, Join(basepath, Rel(basepath, targpath)) is
equivalent to targpath itself. On success, the returned path will
always be relative to basepath, even if basepath and targpath share
no elements. An error is returned if targpath can't be made
relative to basepath or if knowing the current working directory
would be necessary to compute it.

Example:

```

var results = filepath.rel("/a", "/a/b/c");
var path = results[0];  // => "b/c"
var error = results[1];

```

## FILEPATH.clean
<a name="clean"></a>
`filepath.clean(path);`

Clean returns the shortest path name equivalent to path by purely lexical processing.
See [here](https://golang.org/pkg/path/filepath/#Dir) for rules.

Example:

```

var squeaky = filepath.clean("/a//b"); // => "/a/b"

```

## FILEPATH.abs
<a name="abs"></a>
`filepath.abs(path);`

Abs returns an absolute representation of path. If the path is not
absolute it will be joined with the current working directory to
turn it into an absolute path. The absolute path name for a given
file is not guaranteed to be unique.

Example:

```

var results = filepath.abs("/a/../b");
var path = results[0]; // => "/b"
var err = results[1];

```

## FILEPATH.join
<a name="join"></a>
`filepath.join();`

Join joins any number of path elements into a single path, adding a
Separator if necessary. The result is Cleaned.

Example:

```

var joined = filepath.join("/a", "/b");

```

## FILEPATH.walk
<a name="walk"></a>
`filepath.walk(path, walker);`

Walk walks the file tree rooted at root, calling `walker` for each
file or directory in the tree, including root. All errors that
arise visiting files and directories are filtered by `walker`. The
files are walked in lexical order, which makes the output
deterministic but means that for very large directories Walk can be
inefficient. Walk does not follow symbolic links.

Example:

```

filepath.walk("/var", function(path,info,err) { console.log(path); });

```


