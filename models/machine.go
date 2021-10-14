package models

type MachineInput struct {
	// required
	MachineName      string   `json:"machineName"`
	SubnetIDs        []string `json:"subnetIDs"`
	SecurityGroupIDs []string `json:"securityGroupIDs"`

	// optional
	InstanceType            string                   `json:"instanceType,omitempty"`
	ByolOnDedicatedInstance bool                     `json:"byolOnDedicatedInstance,omitempty"`
	DedicatedHostIdentifier string                   `json:"dedicatedHostIdentifier,omitempty"`
	Disks                   []map[string]interface{} `json:"disks,omitempty"`
	ForceUEFI               bool                     `json:"forceUEFI,omitempty"`
	IamRole                 string                   `json:"iamRole,omitempty"`
	NetworkInterface        string                   `json:"networkInterface,omitempty"`
	PlacementGroup          string                   `json:"placementGroup,omitempty"`
	PrivateIPAction         string                   `json:"privateIPAction,omitempty"`
	PrivateIPs              []string                 `json:"privateIPs,omitempty"`
	PublicIPAction          string                   `json:"publicIPAction,omitempty"`
	RunAfterLaunch          bool                     `json:"runAfterLaunch,omitempty"`
	StaticIP                string                   `json:"staticIp,omitempty"`
	StaticIPAction          string                   `json:"staticIpAction,omitempty"`
	SubnetsHostProject      string                   `json:"subnetsHostProject,omitempty"`
	Tags                    []map[string]string      `json:"tags,omitempty"`
	Tenancy                 string                   `json:"tenancy,omitempty"`
	ScsiAdapterType         string                   `json:"scsiAdapterType,omitempty"`
	Cpus                    int                      `json:"cpus,omitempty"`
	MbRam                   int                      `json:"mbRam,omitempty"`
	CoresPerCpu             int                      `json:"coresPerCpu,omitempty"`
	LaunchOnInstanceId      string                   `json:"launchOnInstanceId,omitempty"`
	SecurityGroupAction     string                   `json:"securityGroupAction,omitempty"`
	ComputeLocationId       string                   `json:"computeLocationId,omitempty"`
	LogicalLocationId       string                   `json:"logicalLocationId,omitempty"`
	NetworkAdapterType      string                   `json:"networkAdapterType,omitempty"`
	UseSharedRam            bool                     `json:"useSharedRam,omitempty"`
}
