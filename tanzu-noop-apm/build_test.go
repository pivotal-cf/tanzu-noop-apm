package tanzu_noop_apm_test

import (
	"testing"

	tanzu_noop_apm "github.com/pivotal-cf/tanzu-noop-apm/tanzu-noop-apm"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {

	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it("contributes sample Java agent", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "sample-agent"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "sample-agent",
					"version": "1.0.0",
					"stacks":  []interface{}{"io.buildpacks.stacks.bionic"},
				},
			},
		}
		ctx.StackID = "io.buildpacks.stacks.bionic"

		result, err := tanzu_noop_apm.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(1))
		Expect(result.Layers[0].Name()).To(Equal("sample-agent"))
		Expect(result.BOM.Entries).To(HaveLen(1))
		Expect(result.BOM.Entries[0].Name).To(Equal("sample-agent"))
	})

}
