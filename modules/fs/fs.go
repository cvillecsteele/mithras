package fs

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
