package poststart

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/node"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/kubelet"
)

const StepName = "poststart"

type Step struct {
	script *template.Template
}

func Init() {
	tpl, err := tm.GetTemplate(StepName)

	if err != nil {
		panic(fmt.Sprintf("template %s not found", StepName))
	}

	steps.RegisterStep(StepName, New(tpl))
}

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	ctx2, _ := context.WithTimeout(ctx, config.PostStartConfig.Timeout)
	config.PostStartConfig.IsMaster = config.IsMaster

	if config.IsMaster {
		config.PostStartConfig.Host = config.Node.PrivateIp
	} else {
		if masterNode := config.GetMaster(); masterNode != nil {
			config.PostStartConfig.Host = masterNode.PrivateIp
		}
	}

	err := steps.RunTemplate(ctx2, s.script, config.Runner, out, config.PostStartConfig)

	if err != nil {
		return errors.Wrap(err, "run post start script step")
	}

	// Mark current node as active to allow cluster check task select it for cluster wide task
	config.Node.State = node.StateActive
	config.NodeChan() <- config.Node
	config.SshConfig.BootstrapPrivateKey = ""

	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *Step) Description() string {
	return ""
}

func (s *Step) Depends() []string {
	return []string{kubelet.StepName}
}
