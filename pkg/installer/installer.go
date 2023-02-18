package installer

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/miscord-dev/cluster-api-byoh-both/pkg/machinetemplate"
)

type Installer interface {
	Generate(InstallerConfig) (install, uninstall string, err error)
}

func New() Installer {
	return &installer{}
}

type installer struct {
}

func (i *installer) generateImageTag(ic InstallerConfig) (string, error) {
	if ic.Repository == "" {
		return "", nil
	}

	tmpl, err := machinetemplate.NewTemplate(ic.TagTemplate)

	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", ic.TagTemplate, err)
	}

	tag, err := tmpl.GenerateTag(machinetemplate.NodeInfo{
		OS:         ic.NodeInfo.OS,
		OSImage:    ic.NodeInfo.OSImage,
		Arch:       ic.NodeInfo.Arch,
		K8sVersion: ic.NodeInfo.K8sVersion,
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate tag: %w", err)
	}

	if tag == "" {
		return ic.Repository, nil
	}

	return fmt.Sprintf("%s:%s", ic.Repository, tag), nil
}

func (i *installer) ensureImage(tag string) error {
	ref, err := name.ParseReference(tag)

	if err != nil {
		return fmt.Errorf("failed to parse tag %s: %w", tag, err)
	}

	_, err = remote.Image(ref)

	if err != nil {
		return fmt.Errorf("failed to find image %s: %w", tag, err)
	}

	return nil
}

func (i *installer) generateScript(ic InstallerConfig, tag, template string) (string, error) {
	tmpl, err := machinetemplate.NewTemplate(template)

	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", ic.TagTemplate, err)
	}

	script, err := tmpl.GenerateScript(machinetemplate.ScriptInfo{
		ImageTag:           tag,
		BundleDownloadPath: "{{.BundleDownloadPath}}",
		NodeInfo: machinetemplate.NodeInfo{
			OS:         ic.NodeInfo.OS,
			OSImage:    ic.NodeInfo.OSImage,
			Arch:       ic.NodeInfo.Arch,
			K8sVersion: ic.NodeInfo.K8sVersion,
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate script: %w", err)
	}

	return script, nil
}

func (i *installer) Generate(ic InstallerConfig) (install, uninstall string, err error) {
	tag, err := i.generateImageTag(ic)

	if err != nil {
		return "", "", fmt.Errorf("failed to generate image tag: %w", err)
	}

	if tag != "" {
		if err := i.ensureImage(tag); err != nil {
			return "", "", fmt.Errorf("failed to ensure the image: %w", err)
		}
	}

	install, err = i.generateScript(ic, tag, ic.InstallTemplate)

	if err != nil {
		return "", "", fmt.Errorf("failed to generate install script: %w", err)
	}

	uninstall, err = i.generateScript(ic, tag, ic.UninstallTemplate)

	if err != nil {
		return "", "", fmt.Errorf("failed to generate uninstall script: %w", err)
	}

	return install, uninstall, nil
}
