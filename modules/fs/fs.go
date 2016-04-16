// MITHRAS: Javascript configuration management tool for AWS.
// Copyright (C) 2016, Colin Steele
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//   the Free Software Foundation, either version 3 of the License, or
//                  (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//              GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// @public
//
//
// # CORE FUNCTIONS: FS
//

package fs

// @public
//
// This package exports entry points into the JS environment:
//
// > * [fs.chtimes](#chtimes)
// > * [fs.link](#link)
// > * [fs.symlink](#symlink)
// > * [fs.create](#create)
// > * [fs.close](#close)
// > * [fs.read](#read)
// > * [fs.write](#write)
// > * [fs.copy](#copy)
// > * [fs.chdir](#chdir)
// > * [fs.getwd](#getwd)
// > * [fs.mkdirAll](#mkdirAll)
// > * [fs.remove](#remove)
// > * [fs.removeAll](#removeAll)
// > * [fs.rename](#rename)
// > * [fs.chown](#chown)
// > * [fs.lChown](#lchown)
// > * [fs.chmod](#chmod)
// > * [fs.dir](#dir)
// > * [fs.stat](#stat)
//
// This API allows the caller to work with files.
//
// ## FS.CHTIMES
// <a name="chtimes"></a>
// `filepath.chtimes(path);`
//
// Sets the file at `path` to have current atime and mtime.  TODO: rename this to `touch`.
//
// Example:
//
// ```
//
//  var err = fs.chtimes("/tmp/foo");
//
// ```
//
// ## FS.LINK
// <a name="link"></a>
// `filepath.link(old, new);`
//
// Link creates newname as a hard link to the oldname file. If there
// is an error, it will be of type LinkError.
//
// Example:
//
// ```
//
//  var err = fs.link("/tmp/old" "/tmp/new");
//
// ```
//
// ## FS.SYMLINK
// <a name="symlink"></a>
// `filepath.symlink(old, new);`
//
// Symlink creates newname as a symbolic link to oldname. If there is
// an error, it will be of type LinkError.
//
// Example:
//
// ```
//
//  var err = fs.symlink("/tmp/old" "/tmp/new");
//
// ```
//
// ## FS.CREATE
// <a name="create"></a>
// `filepath.create(path);`
//
// Create creates the named file with mode 0666 (before umask),
// truncating it if it already exists. If successful, methods on the
// returned File can be used for I/O; the associated file descriptor
// has mode O_RDWR. If there is an error, it will be of type
// PathError.
//
// Example:
//
// ```
//
//  var results = fs.create("/tmp/foo");
//  var file = results[0];
//  var error = results[1];
//
// ```
//
// ## FS.CLOSE
// <a name="close"></a>
// `filepath.close(file);`
//
// Close closes the File, rendering it unusable for I/O. It returns an
// error, if any.
//
// Example:
//
// ```
//
//  var results = fs.create("/tmp/foo");
//  var file = results[0];
//  var error = results[1];
//  if (error) {
//   ...
//  }
//  var error = fs.close(file);
//
// ```
//
// ## FS.READ
// <a name="read"></a>
// `filepath.read(path);`
//
// Read the contents of the file at `path`.
//
// Example:
//
// ```
//
//  var results = fs.read("/tmp/foo");
//  var contents = results[0];
//  var error = results[1];
//
// ```
//
// ## FS.WRITE
// <a name="write"></a>
// `filepath.write(path, contents, perms);`
//
// Write the contents of the file at `path`.
//
// Example:
//
// ```
//
//  var error = fs.read("/tmp/foo", "contents", 0644);
//
// ```
//
// ## FS.COPY
// <a name="copy"></a>
// `filepath.copy(src, dest, perms);`
//
// Copy the file from `src` to `dest`.
//
// Example:
//
// ```
//
//  var error = fs.copy("/tmp/foo", "/tmp/bar", 0644);
//
// ```
//
// ## FS.CHDIR
// <a name="chdir"></a>
// `filepath.chdir(dir);`
//
// Change working director to `dir`.
//
// Example:
//
// ```
//
//  var error = fs.chdir("/tmp");
//
// ```
//
// ## FS.GETWD
// <a name="getwd"></a>
// `filepath.getwd();`
//
// Get the current working directory.
//
// Example:
//
// ```
//
//  var results = fs.getwd();
//  var where = results[0];
//  var err = results[1];
//
// ```
//
// ## FS.MKDIRALL
// <a name="mkdirAll"></a>
// `filepath.mkdirALL(path, perm);`
//
// MkdirAll creates a directory named path, along with any necessary
// parents, and returns nil, or else returns an error. The permission
// bits perm are used for all directories that MkdirAll creates. If
// path is already a directory, MkdirAll does nothing and returns nil.
//
// Example:
//
// ```
//
//  var error = fs.mkdirAll("/tmp/a/b/c", 0777);
//
// ```
//
// ## FS.REMOVE
// <a name="remove"></a>
// `filepath.remove(path);`
//
// Remove the file at `path`.
//
// Example:
//
// ```
//
//  var error = fs.remove("/tmp/a/b/c");
//
// ```
//
// ## FS.REMOVEALL
// <a name="removeAll"></a>
// `filepath.removeAll(path);`
//
// RemoveAll removes path and any children it contains. It removes
// everything it can but returns the first error it encounters. If the
// path does not exist,
//
// Example:
//
// ```
//
//  var error = fs.removeAll("/tmp/a/b/c");
//
// ```
//
// ## FS.RENAME
// <a name="rename"></a>
// `filepath.rename(oldname, newname);`
//
// Rename renames (moves) oldpath to newpath. If newpath already
// exists, Rename replaces it. OS-specific restrictions may apply when
// oldpath and newpath are in different directories. If there is an
// error, it will be of type *LinkError.
//
// Example:
//
// ```
//
//  var error = fs.rename("/tmp/foo", "/tmp/bar");
//
// ```
//
// ## FS.CHOWN
// <a name="chown"></a>
// `filepath.chown(path, uid, gid);`
//
// Chown changes the numeric uid and gid of the named file. If the
// file is a symbolic link,it changes the uid and gid of the link's
// target. If there is an error, it will be of type *PathError.
//
// Example:
//
// ```
//
//  var error = fs.chown("/tmp/foo", 12, 34);
//
// ```
//
// ## FS.LCHOWN
// <a name="lchown"></a>
// `filepath.lChown(path, uid, gid);`
//
// Lchown changes the numeric uid and gid of the named file. If the
// file is a symbolic link, it changes the uid and gid of the link
// itself. If there is an error, it will be of type *PathError.
//
// Example:
//
// ```
//
//  var error = fs.lchown("/tmp/foo", 12, 34);
//
// ```
//
// ## FS.CHMOD
// <a name="chmod"></a>
// `filepath.chmod(path, mode);`
//
// Chmod changes the mode of the named file to mode. If the file is a
// symbolic link, it changes the mode of the link's target. If there
// is an error, it will be of type *PathError.
//
// Example:
//
// ```
//
//  var error = fs.chmod("/tmp/foo", 0777);
//
// ```
//
// ## FS.DIR
// <a name="dir"></a>
// `filepath.dir(path);`
//
// Return the directory entries from `path`
//
// Example:
//
// ```
//
//  var results = fs.chmod("/tmp");
//  var entries = results[0];
//  var error = results[1];
//
// ```
//
// ## FS.STAT
// <a name="stat"></a>
// `filepath.stat(path);`
//
// Stat returns the FileInfo structure describing file. If there is an
// error, it will be of type *PathError.
//
// Example:
//
// ```
//
//  var results = fs.stat("/tmp/foo");
//  var info = results[0];
//  var error = results[1];
//
// ```
//
import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "fs"

func dir(src string) (names []string, err error) {
	f, err := os.Open(src)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()
	return f.Readdirnames(0)
}

func read(fileName string) (string, error) {
	content, err := ioutil.ReadFile(fileName)
	return string(content), err
}

func write(fileName, data string, perm uint64) error {
	return ioutil.WriteFile(fileName, []byte(data), os.FileMode(perm))
}

func copy(src, dst string, perm uint64) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	if err = d.Close(); err != nil {
		return err
	}

	return os.Chmod(dst, os.FileMode(perm))
}

func chdir(dir string) error {
	return os.Chdir(dir)
}

func getwd() (string, error) {
	return os.Getwd()
}

func mkdir(dir string, perm uint64) error {
	return os.Mkdir(dir, os.FileMode(perm))
}

func mkdirAll(path string, perm uint64) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

func remove(name string) error {
	return os.Remove(name)
}

func removeAll(name string) error {
	return os.RemoveAll(name)
}

func rename(old, new string) error {
	return os.Rename(old, new)
}

func chown(name string, uid int, gid int) error {
	return os.Chown(name, uid, gid)
}

func lChown(name string, uid int, gid int) error {
	return os.Lchown(name, uid, gid)
}

func chmod(name string, perm uint64) error {
	return os.Chmod(name, os.FileMode(perm))
}

func stat(rt *otto.Otto, name string) otto.Value {
	s, err := os.Stat(name)
	type r struct {
		Name       string
		Size       int64
		Mode       os.FileMode
		ModeString string
		ModTime    time.Time
		IsDir      bool
		IsRegular  bool
		Perm       os.FileMode
		Error      error
	}
	if err == nil {
		result := r{
			Error:      err,
			Name:       s.Name(),
			Size:       s.Size(),
			Mode:       s.Mode(),
			ModeString: s.Mode().String(),
			ModTime:    s.ModTime(),
			IsDir:      s.IsDir(),
			IsRegular:  s.Mode().IsRegular(),
			Perm:       s.Mode().Perm(),
		}
		return mcore.Sanitize(rt, result)
	} else {
		return mcore.Sanitize(rt, err)
	}

}

func create(path string) (*os.File, error) {
	return os.Create(path)
}

func close(f *os.File) error {
	return f.Close()
}

func symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

func link(oldname, newname string) error {
	return os.Link(oldname, newname)
}

func chtimes(name string) error {
	// get current timestamp
	currenttime := time.Now().Local()
	return os.Chtimes(name, currenttime, currenttime)
}

func init() {
	mcore.RegisterInit(func(rt *otto.Otto) {
		fsObj, _ := rt.Object(`fs = {}`)
		fsObj.Set("chtimes", func(call otto.FunctionCall) otto.Value {
			result := chtimes(call.Argument(0).String())
			return mcore.Sanitize(rt, result)
		})
		fsObj.Set("link", func(call otto.FunctionCall) otto.Value {
			result := link(call.Argument(0).String(), call.Argument(1).String())
			return mcore.Sanitize(rt, result)
		})
		fsObj.Set("symlink", func(call otto.FunctionCall) otto.Value {
			result := symlink(call.Argument(0).String(), call.Argument(1).String())
			return mcore.Sanitize(rt, result)
		})
		fsObj.Set("create", func(path string) (otto.Value, otto.Value) {
			f := mcore.Sanitizer(rt)
			x, err := create(path)
			val, err := rt.ToValue(x)
			return val, f(err)
		})
		fsObj.Set("close", func(f *os.File) otto.Value {
			x := mcore.Sanitizer(rt)
			return x(close(f))
		})
		fsObj.Set("read", func(fileName string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(read(fileName))
		})
		fsObj.Set("write", func(fileName, data string, perm uint64) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(write(fileName, data, perm))
		})
		fsObj.Set("copy", copy)
		fsObj.Set("chdir", chdir)
		fsObj.Set("getwd", getwd)
		fsObj.Set("mkdir", mkdir)
		fsObj.Set("mkdirAll", func(call otto.FunctionCall) otto.Value {
			mode, err := call.Argument(1).ToInteger()
			if err != nil {
				log.Fatalf("Invalid argument to 'mkdirAll': %s", err)
			}
			result := mkdirAll(call.Argument(0).String(), uint64(mode))
			return mcore.Sanitize(rt, result)
		})
		fsObj.Set("remove", remove)
		fsObj.Set("removeAll", func(path string) otto.Value {
			f := mcore.Sanitizer(rt)
			return f(removeAll(path))
		})
		fsObj.Set("rename", rename)
		fsObj.Set("chown", chown)
		fsObj.Set("lChown", lChown)
		fsObj.Set("chmod", chmod)
		fsObj.Set("dir", dir)
		fsObj.Set("stat", func(name string) otto.Value {
			return stat(rt, name)
		})
	})
}
