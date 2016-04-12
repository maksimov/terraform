package command

import (
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/cli"
)

func TestStateMv(t *testing.T) {
	state := &terraform.State{
		Modules: []*terraform.ModuleState{
			&terraform.ModuleState{
				Path: []string{"root"},
				Resources: map[string]*terraform.ResourceState{
					"test_instance.foo": &terraform.ResourceState{
						Type: "test_instance",
						Primary: &terraform.InstanceState{
							ID: "bar",
							Attributes: map[string]string{
								"foo": "value",
								"bar": "value",
							},
						},
					},

					"test_instance.baz": &terraform.ResourceState{
						Type: "test_instance",
						Primary: &terraform.InstanceState{
							ID: "foo",
							Attributes: map[string]string{
								"foo": "value",
								"bar": "value",
							},
						},
					},
				},
			},
		},
	}

	statePath := testStateFile(t, state)

	p := testProvider()
	ui := new(cli.MockUi)
	c := &StateMvCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(p),
			Ui:          ui,
		},
	}

	args := []string{
		"-state", statePath,
		"test_instance.foo",
		"test_instance.bar",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Test it is correct
	testStateOutput(t, statePath, testStateMvOutput)

	// Test we have backups
	backups := testStateBackups(t, filepath.Dir(statePath))
	if len(backups) != 1 {
		t.Fatalf("bad: %#v", backups)
	}
	testStateOutput(t, backups[0], testStateMvOutputOriginal)
}

func TestStateMv_stateOutNew(t *testing.T) {
	state := &terraform.State{
		Modules: []*terraform.ModuleState{
			&terraform.ModuleState{
				Path: []string{"root"},
				Resources: map[string]*terraform.ResourceState{
					"test_instance.foo": &terraform.ResourceState{
						Type: "test_instance",
						Primary: &terraform.InstanceState{
							ID: "bar",
							Attributes: map[string]string{
								"foo": "value",
								"bar": "value",
							},
						},
					},
				},
			},
		},
	}

	statePath := testStateFile(t, state)
	stateOutPath := statePath + ".out"

	p := testProvider()
	ui := new(cli.MockUi)
	c := &StateMvCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(p),
			Ui:          ui,
		},
	}

	args := []string{
		"-state", statePath,
		"-state-out", stateOutPath,
		"test_instance.foo",
		"test_instance.bar",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Test it is correct
	testStateOutput(t, stateOutPath, testStateMvOutput_stateOut)
	testStateOutput(t, statePath, testStateMvOutput_stateOutSrc)

	// Test we have backups
	backups := testStateBackups(t, filepath.Dir(statePath))
	if len(backups) != 1 {
		t.Fatalf("bad: %#v", backups)
	}
	testStateOutput(t, backups[0], testStateMvOutput_stateOutOriginal)
}

func TestStateMv_noState(t *testing.T) {
	tmp, cwd := testCwd(t)
	defer testFixCwd(t, tmp, cwd)

	p := testProvider()
	ui := new(cli.MockUi)
	c := &StateMvCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(p),
			Ui:          ui,
		},
	}

	args := []string{"from", "to"}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

const testStateMvOutputOriginal = `
test_instance.baz:
  ID = foo
  bar = value
  foo = value
test_instance.foo:
  ID = bar
  bar = value
  foo = value
`

const testStateMvOutput = `
test_instance.bar:
  ID = bar
  bar = value
  foo = value
test_instance.baz:
  ID = foo
  bar = value
  foo = value
`

const testStateMvOutput_stateOut = `
test_instance.bar:
  ID = bar
  bar = value
  foo = value
`

const testStateMvOutput_stateOutSrc = `
<no state>
`

const testStateMvOutput_stateOutOriginal = `
test_instance.foo:
  ID = bar
  bar = value
  foo = value
`
