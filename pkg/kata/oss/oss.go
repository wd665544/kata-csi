package oss

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
	log "github.com/sirupsen/logrus"
	"github.com/wd/kata-csi/pkg/options"
	"github.com/wd/kata-csi/pkg/utils"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	k8smount "k8s.io/utils/mount"
	"sync"
)

const (
	driverName = "ossplugin.csi.alibabacloud.com"
)

var (
	version = "1.0.0"
)

// OSS the OSS object
type OSS struct {
	driver           *csicommon.CSIDriver
	endpoint         string
	idServer         *csicommon.DefaultIdentityServer
	nodeServer       *nodeServer
	controllerServer csi.ControllerServer
	cap              []*csi.VolumeCapability_AccessMode
	cscap            []*csi.ControllerServiceCapability
}

// NewDriver init oss type of csi driver
func NewDriver(nodeID, endpoint string) *OSS {
	log.Infof("Driver: %v version: %v", driverName, version)

	d := &OSS{}
	d.endpoint = endpoint

	if nodeID == "" {
		nodeID = utils.RetryGetMetaData(InstanceID)
		log.Infof("Use node id : %s", nodeID)
	}
	csiDriver := csicommon.NewCSIDriver(driverName, version, nodeID)
	csiDriver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER})
	csiDriver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_UNKNOWN,
	})

	d.driver = csiDriver
	d.controllerServer = NewControllerServer(d.driver)
	return d
}

// newNodeServer init oss type of csi nodeServer
func newNodeServer(d *OSS) *nodeServer {
	cfg, err := clientcmd.BuildConfigFromFlags(options.MasterURL, options.Kubeconfig)
	if err != nil {
		log.Fatalf("Build kubeconfig is failed, err: %s", err.Error())
	}
	crdClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Create crd client is failed, err: %v", err)
	}
	return &nodeServer{
		k8smounter:           k8smount.New(""),
		DefaultNodeServer:    csicommon.NewDefaultNodeServer(d.driver),
		writeCredentialMutex: sync.Mutex{},
		dynamicClient:        crdClient,
	}
}

// Run start a newNodeServer
func (d *OSS) Run() {
	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(d.endpoint,
		csicommon.NewDefaultIdentityServer(d.driver),
		d.controllerServer,
		newNodeServer(d))
	s.Wait()
}
