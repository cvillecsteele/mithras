 


 # CORE FUNCTIONS: FS


 

 This package exports entry points into the JS environment:

 > * [fs.chtimes](#chtimes)
 > * [fs.link](#link)
 > * [fs.symlink](#symlink)
 > * [fs.create](#create)
 > * [fs.close](#close)
 > * [fs.read](#read)
 > * [fs.bread](#bread)
 > * [fs.write](#write)
 > * [fs.copy](#copy)
 > * [fs.chdir](#chdir)
 > * [fs.getwd](#getwd)
 > * [fs.mkdirAll](#mkdirAll)
 > * [fs.remove](#remove)
 > * [fs.removeAll](#removeAll)
 > * [fs.rename](#rename)
 > * [fs.chown](#chown)
 > * [fs.lChown](#lchown)
 > * [fs.chmod](#chmod)
 > * [fs.dir](#dir)
 > * [fs.stat](#stat)

 This API allows the caller to work with files.

 ## FS.CHTIMES
 <a name="chtimes"></a>
 `filepath.chtimes(path);`

 Sets the file at `path` to have current atime and mtime.  TODO: rename this to `touch`.

 Example:

 ```

  var err = fs.chtimes("/tmp/foo");

 ```

 ## FS.LINK
 <a name="link"></a>
 `filepath.link(old, new);`

 Link creates newname as a hard link to the oldname file. If there
 is an error, it will be of type LinkError.

 Example:

 ```

  var err = fs.link("/tmp/old" "/tmp/new");

 ```

 ## FS.SYMLINK
 <a name="symlink"></a>
 `filepath.symlink(old, new);`

 Symlink creates newname as a symbolic link to oldname. If there is
 an error, it will be of type LinkError.

 Example:

 ```

  var err = fs.symlink("/tmp/old" "/tmp/new");

 ```

 ## FS.CREATE
 <a name="create"></a>
 `filepath.create(path);`

 Create creates the named file with mode 0666 (before umask),
 truncating it if it already exists. If successful, methods on the
 returned File can be used for I/O; the associated file descriptor
 has mode O_RDWR. If there is an error, it will be of type
 PathError.

 Example:

 ```

  var results = fs.create("/tmp/foo");
  var file = results[0];
  var error = results[1];

 ```

 ## FS.CLOSE
 <a name="close"></a>
 `filepath.close(file);`

 Close closes the File, rendering it unusable for I/O. It returns an
 error, if any.

 Example:

 ```

  var results = fs.create("/tmp/foo");
  var file = results[0];
  var error = results[1];
  if (error) {
   ...
  }
  var error = fs.close(file);

 ```

 ## FS.READ
 <a name="read"></a>
 `filepath.read(path);`

 Read the contents of the file at `path`.

 Example:

 ```

  var results = fs.read("/tmp/foo");
  var contents = results[0];
  var error = results[1];

 ```

 ## FS.BREAD
 <a name="bread"></a>
 `filepath.bread(path);`

 Read the contents of the file at `path` and return an array of
 `[content, error]`, where `content` is an array of numbers.

 Example:

 ```

  var results = fs.bread("/tmp/foo");
  var contents = results[0];
  var error = results[1];

 ```

 ## FS.WRITE
 <a name="write"></a>
 `filepath.write(path, contents, perms);`

 Write the contents of the file at `path`.

 Example:

 ```

  var error = fs.read("/tmp/foo", "contents", 0644);

 ```

 ## FS.COPY
 <a name="copy"></a>
 `filepath.copy(src, dest, perms);`

 Copy the file from `src` to `dest`.

 Example:

 ```

  var error = fs.copy("/tmp/foo", "/tmp/bar", 0644);

 ```

 ## FS.CHDIR
 <a name="chdir"></a>
 `filepath.chdir(dir);`

 Change working director to `dir`.

 Example:

 ```

  var error = fs.chdir("/tmp");

 ```

 ## FS.GETWD
 <a name="getwd"></a>
 `filepath.getwd();`

 Get the current working directory.

 Example:

 ```

  var results = fs.getwd();
  var where = results[0];
  var err = results[1];

 ```

 ## FS.MKDIRALL
 <a name="mkdirAll"></a>
 `filepath.mkdirALL(path, perm);`

 MkdirAll creates a directory named path, along with any necessary
 parents, and returns nil, or else returns an error. The permission
 bits perm are used for all directories that MkdirAll creates. If
 path is already a directory, MkdirAll does nothing and returns nil.

 Example:

 ```

  var error = fs.mkdirAll("/tmp/a/b/c", 0777);

 ```

 ## FS.REMOVE
 <a name="remove"></a>
 `filepath.remove(path);`

 Remove the file at `path`.

 Example:

 ```

  var error = fs.remove("/tmp/a/b/c");

 ```

 ## FS.REMOVEALL
 <a name="removeAll"></a>
 `filepath.removeAll(path);`

 RemoveAll removes path and any children it contains. It removes
 everything it can but returns the first error it encounters. If the
 path does not exist,

 Example:

 ```

  var error = fs.removeAll("/tmp/a/b/c");

 ```

 ## FS.RENAME
 <a name="rename"></a>
 `filepath.rename(oldname, newname);`

 Rename renames (moves) oldpath to newpath. If newpath already
 exists, Rename replaces it. OS-specific restrictions may apply when
 oldpath and newpath are in different directories. If there is an
 error, it will be of type *LinkError.

 Example:

 ```

  var error = fs.rename("/tmp/foo", "/tmp/bar");

 ```

 ## FS.CHOWN
 <a name="chown"></a>
 `filepath.chown(path, uid, gid);`

 Chown changes the numeric uid and gid of the named file. If the
 file is a symbolic link,it changes the uid and gid of the link's
 target. If there is an error, it will be of type *PathError.

 Example:

 ```

  var error = fs.chown("/tmp/foo", 12, 34);

 ```

 ## FS.LCHOWN
 <a name="lchown"></a>
 `filepath.lChown(path, uid, gid);`

 Lchown changes the numeric uid and gid of the named file. If the
 file is a symbolic link, it changes the uid and gid of the link
 itself. If there is an error, it will be of type *PathError.

 Example:

 ```

  var error = fs.lchown("/tmp/foo", 12, 34);

 ```

 ## FS.CHMOD
 <a name="chmod"></a>
 `filepath.chmod(path, mode);`

 Chmod changes the mode of the named file to mode. If the file is a
 symbolic link, it changes the mode of the link's target. If there
 is an error, it will be of type *PathError.

 Example:

 ```

  var error = fs.chmod("/tmp/foo", 0777);

 ```

 ## FS.DIR
 <a name="dir"></a>
 `filepath.dir(path);`

 Return the directory entries from `path`

 Example:

 ```

  var results = fs.chmod("/tmp");
  var entries = results[0];
  var error = results[1];

 ```

 ## FS.STAT
 <a name="stat"></a>
 `filepath.stat(path);`

 Stat returns the FileInfo structure describing file. If there is an
 error, it will be of type *PathError.

 Example:

 ```

  var results = fs.stat("/tmp/foo");
  var info = results[0];
  var error = results[1];

 ```


