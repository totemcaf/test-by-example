/*
Copyright © 2022 totemcaf@gmail.com

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/totemcaf/test-by-example.git/internal/model"
	"github.com/totemcaf/test-by-example.git/internal/parsers"
	"github.com/totemcaf/test-by-example.git/internal/runners"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const longDescription = `
"run" reads a test suite description files and executes the tests in it reporting the success or errors.

Test suite description files are YAML files that contain a list of tests to execute.

You can list the files to process, use globs, or point to folders with test suite description files.

All files will be read and combined.
`

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:     "run pathToFiles [pathToFiles ...]",
	Aliases: []string{"test"},
	Short:   "Executes a test suite",
	Long:    longDescription,
	Args:    cobra.MatchAll(cobra.MinimumNArgs(1), validateFilesOrFolders),
	Run:     executeRun,
}

func validateFilesOrFolders(_ *cobra.Command, args []string) error {
	for _, arg := range args {
		if info, err := os.Stat(arg); os.IsNotExist(err) || !(info.IsDir() || info.Mode().IsRegular()) {
			return fmt.Errorf("%s is not a valid file or folder", arg)
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)

	_ = runCmd.Flags().IntP("repetitions", "r", 1, "times to execute the test suite")
	_ = runCmd.Flags().StringP("suite", "s", "", "if multiple suites are found, only run the suite with the given name")
	_ = runCmd.Flags().BoolP("debug", "d", false, "enable debug logging")

	err := viper.BindPFlag("repetitions", runCmd.Flags().Lookup("repetitions"))
	if err != nil {
		panic(err)
	}
}

func executeRun(_ *cobra.Command, paths []string) {

	files := expandPaths(paths)

	repetitions := viper.GetInt("repetitions")
	suiteToExecute := viper.GetString("suite")
	debug := viper.GetBool("debug")

	fmt.Println("Echo: " + strings.Join(files, " "))

	l := makeLogger(debug)

	defer func() {
		_ = l.Sync()
	}()

	logger := l.Sugar()
	testFlowCollection, err := parsers.ReadTestFlowCollectionFrom(logger, files)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	suiteNames, err := verifySuitesToExecute(suiteToExecute, testFlowCollection)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	for repetition := 1; repetition <= repetitions; repetition++ {
		for _, suiteName := range suiteNames {
			testFlow, _ := testFlowCollection.GetTestFlow(suiteName)
			testRunner := runners.NewTestRunner(testFlow, logger)

			logger.Infof("Start running %s (%d/%d)", testFlow.Metadata.Name, repetition, repetitions)
			err = testRunner.Run()

			if err != nil {
				logger.Error(err)
				break
			}

			logger.Infof("Success running %s", testFlow.Metadata.Name)
		}
	}
}

func makeLogger(debug bool) *zap.Logger {

	var config zap.Config

	if debug {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
		config.DisableStacktrace = true
		config.Encoding = "console"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	return zap.Must(config.Build())
}

func verifySuitesToExecute(suiteToExecute string, testFlowCollection model.TestFlowCollection) ([]string, error) {
	var suiteNames []string

	if suiteToExecute != "" {
		suiteNames = strings.Split(suiteToExecute, ",")
		err := checkSuiteNames(testFlowCollection, suiteNames)
		if err != nil {
			return nil, err
		}
	} else {
		suiteNames = testFlowCollection.GetFlowNames()
	}

	return suiteNames, nil
}

func checkSuiteNames(flow model.TestFlowCollection, names []string) error {
	var missingSuiteNames []string

	for _, name := range names {
		if !flow.HasFlow(name) {
			missingSuiteNames = append(missingSuiteNames, name)
		}
	}

	if len(missingSuiteNames) > 0 {
		return fmt.Errorf("the following suites are missing: %s", strings.Join(missingSuiteNames, ", "))
	}

	return nil
}

func expandPaths(paths []string) []string {
	var files []string
	for _, path := range paths {
		files = getFiles(path)
	}

	return files
}

func getFiles(path string) []string {
	var files []string
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return files
	}
	for _, fileInfo := range fileInfos {
		fullPath := filepath.Join(path, fileInfo.Name())
		if fileInfo.IsDir() {
			files = append(files, getFiles(fullPath)...)
		} else if strings.HasSuffix(fileInfo.Name(), ".yaml") {
			files = append(files, fullPath)
		}
	}
	return files
}
