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

package main

import (
	"github.com/cvillecsteele/mithras/modules/cli"
	"github.com/cvillecsteele/mithras/modules/script"

	"github.com/cvillecsteele/mithras/modules/log"
	"github.com/cvillecsteele/mithras/modules/peek"

	"github.com/cvillecsteele/mithras/modules/elasticache"
	"github.com/cvillecsteele/mithras/modules/elb"
	"github.com/cvillecsteele/mithras/modules/instance"
	"github.com/cvillecsteele/mithras/modules/rds"
	"github.com/cvillecsteele/mithras/modules/region"
	"github.com/cvillecsteele/mithras/modules/require"
	"github.com/cvillecsteele/mithras/modules/route53"
	"github.com/cvillecsteele/mithras/modules/secgroup"
	"github.com/cvillecsteele/mithras/modules/subnet"
	"github.com/cvillecsteele/mithras/modules/vpc"

	"github.com/cvillecsteele/mithras/modules/autoscaling"
	"github.com/cvillecsteele/mithras/modules/exec"
	"github.com/cvillecsteele/mithras/modules/filepath"
	"github.com/cvillecsteele/mithras/modules/fs"
	"github.com/cvillecsteele/mithras/modules/iam"
	"github.com/cvillecsteele/mithras/modules/keypairs"
	"github.com/cvillecsteele/mithras/modules/network"
	"github.com/cvillecsteele/mithras/modules/os"
	"github.com/cvillecsteele/mithras/modules/routetables"
	"github.com/cvillecsteele/mithras/modules/s3"
	"github.com/cvillecsteele/mithras/modules/sns"
	"github.com/cvillecsteele/mithras/modules/sqs"
	"github.com/cvillecsteele/mithras/modules/tag"
	"github.com/cvillecsteele/mithras/modules/time"
	"github.com/cvillecsteele/mithras/modules/user"
	"github.com/cvillecsteele/mithras/modules/web"
	"github.com/cvillecsteele/mithras/modules/workers"
)

var Version = "0.1.0"

func main() {
	vers := []script.ModuleVersion{
		script.ModuleVersion{Version: autoscaling.Version, Module: autoscaling.ModuleName},
		script.ModuleVersion{Version: sqs.Version, Module: sqs.ModuleName},
		script.ModuleVersion{Version: sns.Version, Module: sns.ModuleName},
		script.ModuleVersion{Version: keypairs.Version, Module: keypairs.ModuleName},
		script.ModuleVersion{Version: workers.Version, Module: workers.ModuleName},
		script.ModuleVersion{Version: iam.Version, Module: iam.ModuleName},
		script.ModuleVersion{Version: tag.Version, Module: tag.ModuleName},
		script.ModuleVersion{Version: routetables.Version, Module: routetables.ModuleName},
		script.ModuleVersion{Version: filepath.Version, Module: filepath.ModuleName},
		script.ModuleVersion{Version: s3.Version, Module: s3.ModuleName},
		script.ModuleVersion{Version: user.Version, Module: user.ModuleName},
		script.ModuleVersion{Version: os.Version, Module: os.ModuleName},
		script.ModuleVersion{Version: time.Version, Module: time.ModuleName},
		script.ModuleVersion{Version: web.Version, Module: web.ModuleName},
		script.ModuleVersion{Version: exec.Version, Module: exec.ModuleName},
		script.ModuleVersion{Version: fs.Version, Module: fs.ModuleName},
		script.ModuleVersion{Version: network.Version, Module: network.ModuleName},
		script.ModuleVersion{Version: route53.Version, Module: route53.ModuleName},
		script.ModuleVersion{Version: elasticache.Version, Module: elasticache.ModuleName},
		script.ModuleVersion{Version: rds.Version, Module: rds.ModuleName},
		script.ModuleVersion{Version: elb.Version, Module: elb.ModuleName},
		script.ModuleVersion{Version: secgroup.Version, Module: secgroup.ModuleName},
		script.ModuleVersion{Version: require.Version, Module: require.ModuleName},
		script.ModuleVersion{Version: log.Version, Module: log.ModuleName},
		script.ModuleVersion{Version: vpc.Version, Module: vpc.ModuleName},
		script.ModuleVersion{Version: instance.Version, Module: instance.ModuleName},
		script.ModuleVersion{Version: region.Version, Module: region.ModuleName},
		script.ModuleVersion{Version: peek.Version, Module: peek.ModuleName},
		script.ModuleVersion{Version: subnet.Version, Module: subnet.ModuleName},
	}
	cli.Run(vers, Version)
}
