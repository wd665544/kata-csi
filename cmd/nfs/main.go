package main

import (
	"flag"
	"os"

	"github.com/wd/kata-csi/pkg/kata/nfs"
	"k8s.io/klog/v2"
)

var (
	endpoint         = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	nodeID           = flag.String("nodeid", "", "node id")
	mountPermissions = flag.Uint64("mount-permissions", 0777, "mounted folder permissions")
	driverName       = flag.String("drivername", nfs.DefaultDriverName, "name of the driver")
	workingMountDir  = flag.String("working-mount-dir", "/tmp", "working directory for provisioner to mount nfs shares temporarily")
)

func init() {
	_ = flag.Set("logtostderr", "true")
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	if *nodeID == "" {
		klog.Warning("nodeid is empty")
	}

	handle()
	os.Exit(0)
}

func handle() {
	driverOptions := nfs.DriverOptions{
		NodeID:           *nodeID,
		DriverName:       *driverName,
		Endpoint:         *endpoint,
		MountPermissions: *mountPermissions,
		WorkingMountDir:  *workingMountDir,
	}
	d := nfs.NewDriver(&driverOptions)
	d.Run(false)
}

//GOPATH="/opt/go"
// GOPRIVATE=""
// GOPROXY="https://proxy.golang.org,direct"
// GOROOT="/usr/local/go"
