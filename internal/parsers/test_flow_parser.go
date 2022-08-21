package parsers

import (
	"os"

	"bitbucket.org/altscore/test-by-example.git/internal/model"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type DocumentType interface {
	*model.TestFlow | *model.Step
	Validate() error
}

func ReadSpec[Spec DocumentType](fileName string) (Spec, error) {
	bytes, err := os.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	var spec Spec

	err = yaml.UnmarshalStrict(bytes, &spec)

	if err == nil {
		err = spec.Validate()
	}

	return spec, err
}

func ReadTestFlowCollectionFrom(logger *zap.SugaredLogger, files []string) (model.TestFlowCollection, error) {
	collection := model.TestFlowCollection{
		Flows:       make(map[string]*model.TestFlow, 0),
		GlobalSteps: make(map[string]*model.Step, 0),
	}

	for _, file := range files {
		testFlow, err := ReadSpec[*model.TestFlow](file)

		if err == nil {
			name := testFlow.Metadata.Name
			logger.Infof("Read '%s' from %s", testFlow.FullName(), file)
			collection.Flows[name] = testFlow
			continue
		}

		step, err := ReadSpec[*model.Step](file)
		if err == nil {
			name := step.Metadata.Name
			logger.Infof("Read '%s' from %s", step.FullName(), file)
			collection.GlobalSteps[name] = step
			continue
		}

		logger.Errorf("Failed to read %s, skipping it. %s", file, err.Error())
	}

	return collection, nil
}
