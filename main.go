package main

import (
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

	"github.com/cvillecsteele/mithras/modules/channels"
	"github.com/cvillecsteele/mithras/modules/exec"
	"github.com/cvillecsteele/mithras/modules/filepath"
	"github.com/cvillecsteele/mithras/modules/fs"
	"github.com/cvillecsteele/mithras/modules/goroutines"
	"github.com/cvillecsteele/mithras/modules/iam"
	"github.com/cvillecsteele/mithras/modules/network"
	"github.com/cvillecsteele/mithras/modules/os"
	"github.com/cvillecsteele/mithras/modules/routetables"
	"github.com/cvillecsteele/mithras/modules/s3"
	"github.com/cvillecsteele/mithras/modules/tag"
	"github.com/cvillecsteele/mithras/modules/time"
	"github.com/cvillecsteele/mithras/modules/user"
	"github.com/cvillecsteele/mithras/modules/web"
)

var Version = "1.0.0"

type ModuleVersion struct{ version, module string }

func main() {
	vers := []ModuleVersion{
		ModuleVersion{version: channels.Version, module: channels.ModuleName},
		ModuleVersion{version: goroutines.Version, module: goroutines.ModuleName},
		ModuleVersion{version: iam.Version, module: iam.ModuleName},
		ModuleVersion{version: tag.Version, module: tag.ModuleName},
		ModuleVersion{version: routetables.Version, module: routetables.ModuleName},
		ModuleVersion{version: filepath.Version, module: filepath.ModuleName},
		ModuleVersion{version: s3.Version, module: s3.ModuleName},
		ModuleVersion{version: user.Version, module: user.ModuleName},
		ModuleVersion{version: os.Version, module: os.ModuleName},
		ModuleVersion{version: time.Version, module: time.ModuleName},
		ModuleVersion{version: web.Version, module: web.ModuleName},
		ModuleVersion{version: exec.Version, module: exec.ModuleName},
		ModuleVersion{version: fs.Version, module: fs.ModuleName},
		ModuleVersion{version: network.Version, module: network.ModuleName},
		ModuleVersion{version: route53.Version, module: route53.ModuleName},
		ModuleVersion{version: elasticache.Version, module: elasticache.ModuleName},
		ModuleVersion{version: rds.Version, module: rds.ModuleName},
		ModuleVersion{version: elb.Version, module: elb.ModuleName},
		ModuleVersion{version: secgroup.Version, module: secgroup.ModuleName},
		ModuleVersion{version: require.Version, module: require.ModuleName},
		ModuleVersion{version: log.Version, module: log.ModuleName},
		ModuleVersion{version: vpc.Version, module: vpc.ModuleName},
		ModuleVersion{version: instance.Version, module: instance.ModuleName},
		ModuleVersion{version: region.Version, module: region.ModuleName},
		ModuleVersion{version: peek.Version, module: peek.ModuleName},
		ModuleVersion{version: subnet.Version, module: subnet.ModuleName},
	}
	Run(vers)
}
