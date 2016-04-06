package core

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/robertkrimen/otto"
)

var InitFuncs []func(*otto.Otto)
var HandlerFuncs []func(*otto.Otto, *otto.Value, *otto.Value) (*otto.Value, bool)
var PreFlightFuncs = map[string]func(*otto.Otto, *otto.Value, *otto.Value) (*otto.Value, bool){}

func RegisterInit(f func(*otto.Otto)) {
	InitFuncs = append(InitFuncs, f)
}

func RegisterHandler(f func(*otto.Otto, *otto.Value, *otto.Value) (*otto.Value, bool)) {
	HandlerFuncs = append(HandlerFuncs, f)
}

func RegisterPreFlight(name string, f func(*otto.Otto, *otto.Value, *otto.Value) (*otto.Value, bool)) {
	PreFlightFuncs[name] = f
}

func StringsToArray(rt *otto.Otto, ary []*string) *otto.Object {
	o, err := rt.Object(`([])`)
	if err != nil {
		panic(err)
	}
	for _, e := range ary {
		v, err := rt.ToValue(*e)
		if err != nil {
			panic(err)
		}
		o.Call("push", v)
	}
	return o
}

func IsModule(rt *otto.Otto, resource *otto.Value, moduleName string) bool {
	r := *resource.Object()
	js := `(function () { return this["module"] || ""; })`
	name, err := rt.Call(js, r)
	if err != nil {
		log.Fatalf("Resource missing param 'name'.")
	}
	return (moduleName == name.String())
}

func Module(rt *otto.Otto, resource *otto.Value) string {
	r := *resource.Object()
	js := `(function () { return this["module"]; })`
	name, err := rt.Call(js, r)
	if err != nil {
		log.Fatalf("Resource missing param 'name'.")
	}
	return name.String()
}

func IsVerbose(rt *otto.Otto) bool {
	js := `(function () { return mithras["verbose"]; })`
	v, err := rt.Call(js, nil)
	if err != nil {
		log.Fatalf("Mithras object is missing: %s", err)
	}
	if v.IsBoolean() {
		verbose, err := v.ToBoolean()
		if err != nil {
			log.Fatalf("Mithras object error: %s", err)
		}
		return verbose
	}
	log.Fatalf("Mithras object error: 'verbose' is not a boolean")
	return false
}

func ResourceAttr(rt *otto.Otto, resource otto.Object, key string) string {
	js := fmt.Sprintf(`(function () { return this["%s"]; })`, key)
	v, err := rt.Call(js, resource)
	if err != nil {
		log.Fatalf("Error reading resource attribute '%s': %s", key, err)
	}
	if v.IsString() {
		return v.String()
	}
	log.Fatalf("Resource object error: '%s' is not a string", key)
	return ""
}

func ResourceParams(rt *otto.Otto, resource otto.Object) otto.Value {
	js := fmt.Sprintf(`(function () { return this["params"] || {}; })`)
	v, err := rt.Call(js, resource)
	if err != nil {
		log.Fatalf("Error reading resource attribute 'params': %s", err)
	}
	return v
}

func GetParams(rt *otto.Otto, params otto.Object, keys []string) map[string]string {
	paramMap := map[string]string{}
	for _, key := range keys {
		js := fmt.Sprintf(`(function () { return this["%s"] || {}; })`, key)
		val, err := rt.Call(js, params)
		if err == nil && !val.IsUndefined() {
			paramMap[key] = val.String()
		}
	}
	return paramMap
}

func DoResourceCallback(rt *otto.Otto, c otto.Object, r otto.Object, cbName string, args ...interface{}) bool {
	n, err := r.Get("name")
	if err != nil {
		log.Fatalf("Error on resource callback '%s': ", cbName, err)
	}
	rName := n.String()

	m, err := r.Get(cbName)
	if m.IsUndefined() {
		return true
	}

	result, err := r.Call(cbName, args...)
	if err != nil {
		log.Fatalf("Error in resource '%s', callback '%s': ", rName, cbName, err)
	}
	if !result.IsBoolean() {
		log.Fatalf("Error in resource '%s', callback '%s': No boolean returned",
			rName,
			cbName)
	}
	ok, _ := result.ToBoolean()
	return ok
}

func Tag(rt *otto.Otto, tags otto.Object, id string, region string, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

	if verbose {
		log.Printf("  ### Tagging '%s'\n", id)
	}

	params := &ec2.CreateTagsInput{
		Resources: []*string{aws.String(id)},
		Tags:      []*ec2.Tag{},
	}

	js := `(function (cb) {  for (var key in this) {
	                            cb(key, this[key]);
	                         }
	                      })`
	_, err := rt.Call(js,
		tags,
		func(call otto.FunctionCall) otto.Value {
			t := ec2.Tag{
				Key:   aws.String(call.Argument(0).String()),
				Value: aws.String(call.Argument(1).String()),
			}
			params.Tags = append(params.Tags, &t)
			return otto.FalseValue()
		})
	if err != nil {
		log.Fatalf("Error obtaining tags from object.")
	}

	_, err = svc.CreateTags(params)
	if err != nil {
		log.Fatalf("Error tagging '%s': %s", id, err)
	}
}

func InstanceToObject(rt *otto.Otto, d *ec2.Instance) *otto.Object {
	j, err := json.Marshal(*d)
	if err != nil {
		log.Fatalf("Instance marshal error: %s", err)
	}
	js := `(function (json) { return JSON.parse(json); })`
	v, err := rt.Call(js, nil, string(j))
	if err != nil {
		log.Fatalf("Json parse error: %s", err)
	}

	return v.Object()
}

func SSHKeypath(rt *otto.Otto, r *otto.Object, instance *ec2.Instance) string {
	js := `(function (instance) {
                 if (this.params && typeof(this.params.sshKeyPathForInstance) === 'function') {
                   return this.params.sshKeyPathForInstance(instance);
                 } else {
                   return mithras.sshKeyPathForInstance({}, instance);
                 }
               })`
	val, err := rt.Call(js, r, InstanceToObject(rt, instance))
	if err != nil || val.IsUndefined() {
		log.Fatalf("Error obtaining keypath for instance '%s': %s.", instance.InstanceId, err)
	}
	return val.String()
}

func SSHUser(rt *otto.Otto, r *otto.Object, instance *ec2.Instance) string {
	js := `(function (instance) {
                 if (this.params && typeof(this.params.sshUserForInstance) === 'function') {
                   return this.params.sshUserForInstance(instance);
                 } else {
                   return mithras.sshUserForInstance({}, instance);
                 }
               })`
	val, err := rt.Call(js, r, InstanceToObject(rt, instance))
	if err != nil || val.IsUndefined() {
		log.Fatalf("Error obtaining user for instance '%s': %s.", instance.InstanceId, err)
	}
	return val.String()
}

func Sanitize(rt *otto.Otto, o interface{}) otto.Value {
	j, err := json.Marshal(o)
	if err != nil {
		log.Fatalf("Sanitize marhsall error:", err)
	}
	js := `(function (json) { return JSON.parse(json); })`
	val, err := rt.Call(js, nil, string(j))
	if err != nil {
		log.Fatalf("Sanitize can't create object: %s", err)
	}
	return val
}

func Sanitizer(rt *otto.Otto) func(objs ...interface{}) otto.Value {
	return func(objs ...interface{}) otto.Value {
		marshalled := []string{}
		for _, o := range objs {
			j, err := json.Marshal(o)
			if err != nil {
				log.Fatalf("Sanitizer marhsall error:", err)
			}
			marshalled = append(marshalled, string(j))
		}
		js := `(function (things) {
           var parsed = _.map(things, function(thing){ return JSON.parse(thing); });
           if (parsed.length == 1) {
             return parsed[0];
           } else {
             return parsed;
           }
          })`
		val, err := rt.Call(js, nil, marshalled)
		if err != nil {
			log.Fatalf("Sanitizer error sanitizing objects: %s", err)
		}

		return val
	}
}

func init() {
	RegisterInit(func(rt *otto.Otto) {
		var o1 *otto.Object
		var mObj *otto.Object
		if a, err := rt.Get("mithras"); err != nil || a.IsUndefined() {
			mObj, _ = rt.Object(`mithras = {}`)
		} else {
			mObj = a.Object()
		}

		if b, err := mObj.Get("modules"); err != nil || b.IsUndefined() {
			o1, _ = rt.Object(`mithras.modules = {}`)
		} else {
			o1 = b.Object()
		}

		f := func(call otto.FunctionCall) otto.Value {
			name := call.Argument(0).String()
			cb := call.Argument(1)
			f := func(rt *otto.Otto, c *otto.Value, r *otto.Value) (*otto.Value, bool) {
				js := `(function (f, c, r) {  return f(c, r); })`
				result, err := rt.Call(js, nil, cb, *c, *r)
				if err != nil {
					log.Fatalf("Error in preflight '%s': %s", name, err)
				}

				js = `(function () {  return this[0]; })`
				target, err := rt.Call(js, result)
				if err != nil {
					log.Fatalf("Error in preflight '%s': %s", name, err)
				}

				js = `(function () {  return this[1]; })`
				v, err := rt.Call(js, result)
				if err != nil {
					log.Fatalf("Error in preflight '%s': %s", name, err)
				}
				handled, err := v.ToBoolean()
				if err != nil {
					log.Fatalf("Error in handler '%s': %s", name, err)
				}

				if target.Class() != "" {
					return &target, handled
				} else {
					return nil, handled
				}
			}
			RegisterPreFlight(name, f)
			return otto.Value{}
		}
		o1.Set("preflight", f)

		f = func(call otto.FunctionCall) otto.Value {
			name := call.Argument(0).String()
			cb := call.Argument(1)
			f := func(rt *otto.Otto, c *otto.Value, r *otto.Value) (*otto.Value, bool) {
				js := `(function (f, c, r) {  return f(c, r); })`
				result, err := rt.Call(js, nil, cb, *c, *r)
				if err != nil {
					log.Fatalf("Error in handler '%s': %s", name, err)
				}

				js = `(function () {  return this[0]; })`
				target, err := rt.Call(js, result)
				if err != nil {
					log.Fatalf("Error in handler '%s': %s", name, err)
				}

				js = `(function () {  return this[1]; })`
				v, err := rt.Call(js, result)
				if err != nil {
					log.Fatalf("Error in handler '%s': %s", name, err)
				}
				handled, err := v.ToBoolean()
				if err != nil {
					log.Fatalf("Error in handler '%s': %s", name, err)
				}

				if !target.IsUndefined() {
					return &target, handled
				} else {
					return nil, handled
				}
			}
			RegisterHandler(f)
			return otto.Value{}
		}
		o1.Set("handle", f)

		rt.Set("sanitize", func(call otto.FunctionCall) otto.Value {
			val, err := call.Argument(0).Export()
			if err != nil {
				log.Fatalf("Sanitize export error: %s", err)
			}
			return Sanitize(rt, val)
		})
	})
}
