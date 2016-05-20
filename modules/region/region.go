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
// # CORE FUNCTIONS: REGIONS
//

package region

// @public
//
// This package exports several entry points into the JS environment,
// including:
//
// > * [aws.regions.scan](#scan)
//
// This API allows resource handlers to query AWS regions
//
// ## AWS.REGIONS.SCAN
// <a name="scan"></a>
// `aws.regions.scan(region);`
//
// Returns a list of regionss.
//
// Example:
//
// ```
//
//  var regions =  aws.regions.scan("us-east-1");
//
// ```
//

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/robertkrimen/otto"

	mcore "github.com/cvillecsteele/mithras/modules/core"
)

var Version = "1.0.0"
var ModuleName = "region"

func scan(rt *otto.Otto) otto.Value {
	svc := ec2.New(session.New(), aws.NewConfig().WithRegion("us-east-1").WithMaxRetries(5))

	resp, err := svc.DescribeRegions(nil)
	if err != nil {
		panic(err)
	}

	regions := []ec2.Region{}
	for _, r := range resp.Regions {
		regions = append(regions, *r)
	}
	return mcore.Sanitize(rt, regions)
}

func init() {
	mcore.RegisterInit(func(context *mcore.Context) {
		rt := context.Runtime

		rt.Object(`aws = aws || {}`)
		o2, _ := rt.Object(`aws.regions = {}`)
		o2.Set("scan", func() otto.Value {
			return scan(rt)
		})
	})
}
