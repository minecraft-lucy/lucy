/*
Copyright 2024 4rcadia

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package detector

// detectorRegistry manages registered detectors
type detectorRegistry struct {
	executableDetectors  []ExecutableDetector
	modDetectors         []ModDetector
	environmentDetectors []EnvironmentDetector
}

// Global registry instance
var registry = &detectorRegistry{
	executableDetectors:  make([]ExecutableDetector, 0),
	modDetectors:         make([]ModDetector, 0),
	environmentDetectors: make([]EnvironmentDetector, 0),
}

// RegisterExecutableDetector adds a new executable detector to the registry
func RegisterExecutableDetector(detector ExecutableDetector) {
	registry.executableDetectors = append(
		registry.executableDetectors,
		detector,
	)
	// Priority mechanism removed: retain insertion order
}

// RegisterModDetector adds a new mod detector to the registry
func RegisterModDetector(detector ModDetector) {
	registry.modDetectors = append(registry.modDetectors, detector)
	// Priority mechanism removed: retain insertion order
}

// RegisterEnvironmentDetector adds a new environment detector to the registry
func RegisterEnvironmentDetector(detector EnvironmentDetector) {
	registry.environmentDetectors = append(
		registry.environmentDetectors,
		detector,
	)
	// Priority mechanism removed: retain insertion order
}

// GetExecutableDetectors returns all registered executable detectors
func GetExecutableDetectors() []ExecutableDetector {
	return registry.executableDetectors
}

// GetModDetectors returns all registered mod detectors
func GetModDetectors() []ModDetector {
	return registry.modDetectors
}

// GetEnvironmentDetectors returns all registered environment detectors
func GetEnvironmentDetectors() []EnvironmentDetector {
	return registry.environmentDetectors
}
