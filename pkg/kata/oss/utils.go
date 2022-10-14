package oss

import (
	"crypto/sha256"
	"fmt"
	utils2 "github.com/wd/kata-csi/pkg/utils"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	// MetadataURL is metadata url
	MetadataURL = "http://100.100.100.200/latest/meta-data/"
	// InstanceID is instance ID
	InstanceID = "instance-id"
	// RAMRoleResource is ram-role url subpath
	RAMRoleResource = "ram/security-credentials/"
)

// GetMetaData get host regionid, zoneid
func GetMetaData(resource string) string {
	resp, err := http.Get(MetadataURL + resource)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

func GetGlobalMountPath(volumeId string) string {

	result := sha256.Sum256([]byte(fmt.Sprintf("%s", volumeId)))
	volSha := fmt.Sprintf("%x", result)

	globalFileVer1 := filepath.Join(utils2.KubeletRootDir, "/plugins/kubernetes.io/csi/pv/", volumeId, "/globalmount")
	globalFileVer2 := filepath.Join(utils2.KubeletRootDir, "/plugins/kubernetes.io/csi/", driverName, volSha, "/globalmount")

	if utils2.IsFileExisting(globalFileVer1) {
		return globalFileVer1
	} else {
		return globalFileVer2
	}
}

// GetRAMRoleOption get command line's ram_role option
func GetRAMRoleOption() string {
	ramRole := GetMetaData(RAMRoleResource)
	ramRoleOpt := MetadataURL + RAMRoleResource + ramRole
	mntCmdRamRole := fmt.Sprintf("-oram_role=%s", ramRoleOpt)
	return mntCmdRamRole
}

// IsOssfsMounted return if oss mountPath is mounted
func IsOssfsMounted(mountPath string) bool {
	checkMountCountCmd := fmt.Sprintf("%s mount | grep %s | grep -E 'fuse.ossfs|fuse.jindo-fuse' | grep -v grep | wc -l", NsenterCmd, mountPath)
	out, err := utils2.Run(checkMountCountCmd)
	if err != nil {
		return false
	}
	if strings.TrimSpace(out) == "0" {
		return false
	}
	return true
}

// IsLastSharedVol return code status to help check if this oss volume uses UseSharedPath and is the last one
func IsLastSharedVol(pvName string) (string, error) {
	keyStr := fmt.Sprintf("volumes/kubernetes.io~csi/%s/mount", pvName)
	checkMountCountCmd := fmt.Sprintf("%s mount | grep %s | grep -E 'fuse.ossfs|fuse.jindo-fuse' | grep -v grep | wc -l", NsenterCmd, keyStr)
	out, err := utils2.Run(checkMountCountCmd)
	if err != nil {
		return "0", err
	}
	return strings.TrimSpace(out), nil
}
