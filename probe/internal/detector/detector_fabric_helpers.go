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

import (
	"strings"

	"lucy/dependency"
	"lucy/types"
)

// parseFabricVersionRange parses Fabric version range strings into 2D constraint arrays.
//
// According to Fabric spec:
// - Space-separated ranges are AND relationships (all must be satisfied)
// - Multiple versions in array form would be OR relationships
//
// Returns [][]VersionConstraint where:
// - outer array represents OR alternatives
// - inner array represents AND constraints
func parseFabricVersionRange(s string) types.VersionConstraintExpression {
	s = strings.TrimSpace(s)
	if s == "*" {
		return nil
	}
	if s == "" {
		return nil
	}

	// Parse all space-separated constraints (AND relationship)
	var andConstraints []types.VersionConstraint
	parts := strings.Fields(s)
	for _, part := range parts {
		// Handle comma-separated alternatives (not typical in Fabric, but for safety)
		if strings.Contains(part, ",") {
			subParts := strings.Split(part, ",")
			for _, subPart := range subParts {
				andConstraints = append(andConstraints, parseSingleFabricVersion(subPart))
			}
		} else {
			andConstraints = append(andConstraints, parseSingleFabricVersion(part))
		}
	}

	// Wrap in outer array for OR relationship
	if len(andConstraints) > 0 {
		return [][]types.VersionConstraint{andConstraints}
	}
	return nil
}

func parseSingleFabricVersion(version string) types.VersionConstraint {
	version = strings.TrimSpace(version)
	op := types.OpEq
	if strings.HasPrefix(version, "<") {
		op = types.OpLt
		version = strings.TrimPrefix(version, "<")
	} else if strings.HasPrefix(version, "<=") {
		op = types.OpLte
		version = strings.TrimPrefix(version, "<=")
	} else if strings.HasPrefix(version, ">") {
		op = types.OpGt
		version = strings.TrimPrefix(version, ">")
	} else if strings.HasPrefix(version, ">=") {
		op = types.OpGte
		version = strings.TrimPrefix(version, ">=")
	}

	return types.VersionConstraint{
		Value:    dependency.Parse(types.RawVersion(version), types.Semver),
		Operator: op,
	}
}
