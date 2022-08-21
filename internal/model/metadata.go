package model

import "fmt"

type Metadata struct {
	Name        string
	Annotations map[string]string
}

func (m Metadata) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("name is mandatory")
	}
	return nil
}
