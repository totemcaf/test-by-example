package model

type TestFlowCollection struct {
	Flows       map[string]*TestFlow
	GlobalSteps map[string]*Step
}

func (c *TestFlowCollection) GetFlowNames() []string {
	var names []string
	for name := range c.Flows {
		names = append(names, name)
	}
	return names
}

func (c *TestFlowCollection) GetTestFlow(name string) (*TestFlow, bool) {
	flow, found := c.Flows[name]
	if found {
		flow.SetParent(c)
	}
	return flow, found
}

func (c *TestFlowCollection) HasFlow(name string) bool {
	_, found := c.Flows[name]
	return found
}

func (c *TestFlowCollection) GetGlobalStepSpec(name string) (*StepSpec, bool) {
	step, found := c.GlobalSteps[name]
	if found {
		return &step.Spec, true
	}
	return nil, false
}
