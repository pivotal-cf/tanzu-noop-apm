package tanzu_noop_apm_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnit(t *testing.T) {
	suite := spec.New("tanzu-noop-apm", spec.Report(report.Terminal{}))
	suite("Build", testBuild)
	suite("Detect", testDetect)
	suite("tanzu-noop-apm", testJavaAgent)
	suite.Run(t)
}
