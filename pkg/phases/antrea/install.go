package antrea

import (
	"github.com/flanksource/karina/pkg/platform"
)

const (
	Namespace  = "kube-system"
	Controller = "antrea-controller"
	DaemonSet  = "anrtea-agent"
)

func Install(p *platform.Platform) error {
	if !p.Calico.IsDisabled() && !p.Antrea.IsDisabled() {
		p.Fatalf("both calico and antrea are enabled. Please disable one of them to continue")
	}

	if p.Antrea.IsDisabled() {
		return p.DeleteSpecs(Namespace, "antrea.yaml")
	}
	p.Antrea.IsCertReady = false
	if err := p.ApplySpecs(Namespace, "antrea.yaml"); err != nil {
		return err
	}

	if err := platform.WaitForDeployment(Namespace, Controller, 2*time.Minute); err != nil {
		return err
	}
	return platform.WaitForDaemonSet(Namespace, DaemonSet, 2*time.Minute)
}