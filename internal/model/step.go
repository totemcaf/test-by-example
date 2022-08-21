package model

import "fmt"

type Headers map[string]string

type Step struct {
	ApiVersion string   `yaml:"apiVersion" description:"Group and version of this API"`
	Kind       string   `yaml:"kind" description:"Kind of API"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       StepSpec `yaml:"spec"`
}

func (s Step) FullName() string {
	return fmt.Sprintf("%s@%s/%s", s.Kind, s.ApiVersion, s.Metadata.Name)
}

func (s Step) Validate() error {
	if s.ApiVersion != ApiVersion || s.Kind != TestStepKind {
		return fmt.Errorf("invalid API version or kind")
	}

	if err := s.Metadata.Validate(); err != nil {
		return err
	}

	return nil
}

type StepSpec struct {
	Get      *string
	Post     *string
	Name     *string
	Headers  Headers
	Body     *Json
	Response *Response
}

func (s StepSpec) Method() string {
	if s.Get != nil {
		return "GET"
	}
	if s.Post != nil {
		return "POST"
	}

	return "????" // TODO
}

func (s StepSpec) Url() string {
	if s.Get != nil {
		return *s.Get
	}
	if s.Post != nil {
		return *s.Post
	}

	return "????" // TODO
}

func (s StepSpec) NameOrUrl() string {
	if s.Name != nil {
		return *s.Name
	}
	return s.Url()
}

// IsReference returns true if this Step is not defined here, but references a
// global defined step
func (s StepSpec) IsReference() bool {
	return s.Get == nil && s.Post == nil && s.Body == nil && s.Response == nil
}
