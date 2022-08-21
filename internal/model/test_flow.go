package model

import "fmt"

type TestFlowSpec struct {
	Environment map[string]string `yaml:"fromEnvironment,omitempty"`
	Values      map[string]any    `yaml:"values,omitempty"`
	Steps       []StepSpec        `yaml:"steps,omitempty"`
}

type TestFlow struct {
	ApiVersion string       `yaml:"apiVersion" description:"Group and version of this API"`
	Kind       string       `yaml:"kind" description:"Kind of API"`
	Metadata   Metadata     `yaml:"metadata"`
	Spec       TestFlowSpec `yaml:"spec"`
	parent     *TestFlowCollection
}

func (t *TestFlow) Validate() error {
	if t.ApiVersion != ApiVersion || t.Kind != TestFlowKind {
		return fmt.Errorf("invalid API version or kind")
	}

	if err := t.Metadata.Validate(); err != nil {
		return err
	}

	return nil
}

func (t *TestFlow) FullName() string {
	return fmt.Sprintf("%s@%s/%s", t.Kind, t.ApiVersion, t.Metadata.Name)
}

func (t *TestFlow) GetGlobalStepSpec(name string) (*StepSpec, bool) {
	return t.parent.GetGlobalStepSpec(name)
}

func (t *TestFlow) SetParent(parent *TestFlowCollection) {
	t.parent = parent
}
