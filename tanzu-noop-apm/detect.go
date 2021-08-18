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
	"github.com/buildpacks/libcnb"
)

const (
	PlanEntrySample = "sample-agent"
)

type Detect struct {
}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	return libcnb.DetectResult{
		Pass: true,
		Plans: []libcnb.BuildPlan{
			{
				// Indicates that our Buildpack 'provides' a dependency called 'sample-agent'
				Provides: []libcnb.BuildPlanProvide{
					{Name: PlanEntrySample},
				},
				// All provided requirements must be 'required' by some buildpack -
				// in our case we will specify that it is required by our own buildpack
				// itself since it provides a standalone agent
				Requires: []libcnb.BuildPlanRequire{
					{Name: PlanEntrySample},
					{Name: "jvm-application"},
				},
			},
		},
	}, nil
}
