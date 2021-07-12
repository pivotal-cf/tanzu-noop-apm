/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tanzu_noop_apm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
)

type tanzuNoopApm struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewTanzuNoopApm(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) (tanzuNoopApm, libcnb.BOMEntry) {

	// Call libpak method to create a new 'contributor' which contributes our dependency to a 'Launch' layer
	contributor, entry := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Launch: true,
	})
	return tanzuNoopApm{LayerContributor: contributor}, entry
}

func (w tanzuNoopApm) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	w.LayerContributor.Logger = w.Logger

	return w.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {

		// Untar/Unzip the dependency to the layer path
		w.Logger.Bodyf("Expanding to %s", layer.Path)
		if err := crush.ExtractTarXz(artifact, layer.Path, 1); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to expand sample-agent\n%w", err)
		}
		// Create a bin directory so that the dependency is automatically added to $PATH at launch
		binDir := filepath.Join(layer.Path, "bin")

		if err := os.MkdirAll(binDir, 0755); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to mkdir\n%w", err)
		}
		// Create a symlink from bin/sample-agent.jar to where we un-zipped the agent
		if err := os.Symlink(filepath.Join(layer.Path, "sample-agent"), filepath.Join(binDir, "sample-agent")); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to symlink sample-agent\n%w", err)
		}
		// Finally add the agent to the JAVA_TOOL_OPTIONS env var via '-javaagent' flag - this points to the agent path
		layer.LaunchEnvironment.Appendf("JAVA_TOOL_OPTIONS", " ",
			"-javaagent:%s", filepath.Join(layer.Path, "agent", "sampleagent.jar"))

		return layer, nil
	})
}

func (w tanzuNoopApm) Name() string {
	return w.LayerContributor.LayerName()
}
