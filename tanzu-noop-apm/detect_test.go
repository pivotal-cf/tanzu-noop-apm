package tanzu_noop_apm_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	tanzu_noop_apm "github.com/pivotal-cf/tanzu-noop-apm/tanzu-noop-apm"
	"github.com/sclevine/spec"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect tanzu_noop_apm.Detect
	)

	it("passes detection", func() {
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "sample-agent"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "sample-agent"},
						{Name: "jvm-application"},
					},
				},
			},
		}))
	})
}
