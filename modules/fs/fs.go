//
// # CORE FUNCTIONS: FS
//

package fs

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
// `filepath.create();`
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
