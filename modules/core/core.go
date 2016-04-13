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

//
// # CORE FUNCTIONS: CORE
//

package core

// This package exports several entry points into the JS environment,
// including:
//
// > * [sanitize](#sanitize)
//
// This API allows exposes "core" functions to JS.
//
// ## SANITIZE
// <a name="sanitize"></a>
// `sanitize(thing_from_go);`
//
// Make an object returned from a raw Go function and make it
// digestible by JS.
//
// TODO: kill this???
//
// Example:
//
// ```
//
//  sanitize(<something>);
//
// ```
//

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/robertkrimen/otto"
)

var InitFuncs []func(*otto.Otto)

func RegisterInit(f func(*otto.Otto)) {
	InitFuncs = append(InitFuncs, f)
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

func Tag(rt *otto.Otto, tags otto.Object, id string, region string, verbose bool) {
	svc := ec2.New(session.New(),
		aws.NewConfig().WithRegion(region).WithMaxRetries(5))

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
		rt.Set("sanitize", func(call otto.FunctionCall) otto.Value {
			val, err := call.Argument(0).Export()
			if err != nil {
				log.Fatalf("Sanitize export error: %s", err)
			}
			return Sanitize(rt, val)
		})
	})
}
