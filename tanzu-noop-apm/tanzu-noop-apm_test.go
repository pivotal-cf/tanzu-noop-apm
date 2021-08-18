package tanzu_noop_apm_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	tanzu_noop_apm "github.com/pivotal-cf/tanzu-noop-apm/tanzu-noop-apm"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"
)

func testJavaAgent(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Buildpack.Path, err = ioutil.TempDir("", "sample-agent-buildpack")
		Expect(err).NotTo(HaveOccurred())

		ctx.Layers.Path, err = ioutil.TempDir("", "sample-agent-layers")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Buildpack.Path)).To(Succeed())
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	it("contributes Java agent", func() {

		dep := libpak.BuildpackDependency{
			ID:     "sample-agent",
			URI:    "https://localhost/stub-sample-agent.zip",
			SHA256: "fae6690187c2ea9473bce63d8618040567c7273625d4e0486de2326fae7ba946",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		j, bomEntry := tanzu_noop_apm.NewTanzuNoopApm(dep, dc)
		Expect(bomEntry).To(HaveLen(1))
		Expect(bomEntry.Name).To(Equal("sample-agent"))
		Expect(bomEntry.Metadata["layer"]).To(Equal("sample-agent"))
		Expect(bomEntry.Launch).To(BeTrue())
		Expect(bomEntry.Build).To(BeFalse())

		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = j.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.Launch).To(BeTrue())
		Expect(filepath.Join(layer.Path, "sampleagent.jar")).To(BeARegularFile())
		Expect(layer.LaunchEnvironment["JAVA_TOOL_OPTIONS.delim"]).To(Equal(" "))
		Expect(layer.LaunchEnvironment["JAVA_TOOL_OPTIONS.append"]).To(Equal(fmt.Sprintf("-javaagent:%s",
			filepath.Join(layer.Path, "sampleagent.jar"))))
	})
}
