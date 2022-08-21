package runners

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/altscore/test-by-example.git/internal/contexts"
	"bitbucket.org/altscore/test-by-example.git/internal/evaluators"
	"bitbucket.org/altscore/test-by-example.git/internal/model"
	"bitbucket.org/altscore/test-by-example.git/pkg/jsonx"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type TestRunner interface {
	Run() error
}

type testRunner struct {
	testFlow *model.TestFlow
	contexts.RunningContext
	client *resty.Client
	logger *zap.SugaredLogger
}

func NewTestRunner(testFlow *model.TestFlow, logger *zap.SugaredLogger) *testRunner {
	return &testRunner{
		testFlow,
		contexts.NewRunningContext(logger),
		resty.New(),
		logger,
	}
}

func (r *testRunner) Run() error {
	r.initContext()
	return r.runSteps(r.testFlow.Spec.Steps)
}

func (r *testRunner) initContext() {
	r.initEnvironmentVars()
	r.initValues()
}

func (r *testRunner) initEnvironmentVars() {
	for name, envVarName := range r.testFlow.Spec.Environment {
		r.Set(name, viper.GetString(envVarName))
	}
}

func (r *testRunner) initValues() {
	for name, value := range r.testFlow.Spec.Values {
		r.Set(name, value)
	}
}

func (r *testRunner) runSteps(steps []model.StepSpec) error {
	for _, step := range steps {
		var stepRef *model.StepSpec
		if step.IsReference() {
			var found bool
			stepRef, found = r.testFlow.GetGlobalStepSpec(step.NameOrUrl())
			if !found {
				return fmt.Errorf("step '%s' not found", step.NameOrUrl())
			}
		} else {
			stepRef = &step
		}

		if err := r.runStep(stepRef); err != nil {
			return err
		}
	}
	return nil
}

func (r *testRunner) runStep(step *model.StepSpec) error {
	r.logger.Infof("Running '%s'%s", step.NameOrUrl(), referenceType(step))
	request := r.client.R()

	if err := r.setHeaders(request, step.Headers); err != nil {
		return err
	}
	if err := r.setBody(request, step.Body); err != nil {
		return err
	}

	var resultBody map[string]any

	err := r.setResult(request, &resultBody)

	if err != nil {
		return err
	}

	response, err := r.execute(request, step)

	if err != nil {
		return err
	}

	if err := r.processResult(response, resultBody, step); err != nil {
		return err
	}

	return err
}

func referenceType(step *model.StepSpec) string {
	if step.IsReference() {
		return " (reference)"
	}
	return ""
}

func (r *testRunner) setHeaders(request *resty.Request, headers map[string]string) error {
	eval := evaluators.NewJsonXEvaluator(r.RunningContext)
	for key, value := range headers {

		name := eval.EvaluateStr(key)
		valueStr := eval.EvaluateStr(value)
		fmt.Printf("Using Header: %s: %s\n", name, valueStr)
		request.SetHeader(name, valueStr)
	}

	return nil
}

func (r *testRunner) setBody(request *resty.Request, body *model.Json) error {
	parser := jsonx.NewParser()

	jsonXBody := parser.Parse(body)

	finalBody := jsonXBody.Eval(r.RunningContext)

	jsonBody, err := json.Marshal(finalBody)

	if err != nil {
		panic(err)
	}
	fmt.Printf("\nUsing Body: %s\n\n", jsonBody)
	request.SetBody(finalBody)
	return nil

}

func (r *testRunner) execute(request *resty.Request, step *model.StepSpec) (*resty.Response, error) {
	eval := evaluators.NewJsonXEvaluator(r.RunningContext)
	url := eval.EvaluateStr(step.Url())

	return request.Execute(step.Method(), url)
}

func (r *testRunner) processResult(response *resty.Response, actualBody map[string]any, step *model.StepSpec) error {
	if err := r.checkResponseCode(response, step); err != nil {
		return err
	}

	differ := jsonx.NewDiffer(r.RunningContext)

	if jsonStr, err := json.Marshal(actualBody); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Actual Body Result: %s\n\n", jsonStr)
	}

	return differ.Compare(step.Response.Body, actualBody)
}

func (r *testRunner) setResult(request *resty.Request, m *map[string]any) error {
	request.SetResult(m)

	return nil
}

func (r *testRunner) checkResponseCode(response *resty.Response, step *model.StepSpec) error {
	expected := step.Response.StatusCode
	actual := response.StatusCode()
	if actual != expected {
		return fmt.Errorf("[%s] expected status %d, received %d. Msg: %s", step.NameOrUrl(), expected, actual, response.Body())
	}
	return nil
}
