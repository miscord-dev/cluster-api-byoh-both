package installer

type NodeInfo struct {
	OS         string
	OSImage    string
	Arch       string
	K8sVersion string
}

type InstallerConfig struct {
	NodeInfo          NodeInfo
	InstallTemplate   string
	UninstallTemplate string
	Repository        string
	TagTemplate       string
}
