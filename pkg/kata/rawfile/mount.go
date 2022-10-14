package nfs

type MountInfo struct {
	// The type of the volume (ie. block)
	VolumeType string`json:"volume-type"`

	// The device backing the volume.
	Device string `json:"device"`

	// The filesystem type to be mounted on the volume.
	FsType string `json:"fstype"`

	// Additional metadata to pass to the agentregarding this volume.
	Metadata map[string]string`json:"metadata,omitempty"`

	// Additional mount options.
	Options []string`json:"options,omitempty"`
}

func NewMountInfo(volumeType, device, fsType string, metadata map[string]string, options []string) MountInfo{
	return MountInfo{
		VolumeType: volumeType,
		FsType: fsType,
		Metadata: metadata,
		Options: options,
	}
}