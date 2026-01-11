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

package mcdr

import (
	"lucy/dependency"
	"lucy/types"
	"strings"
)

func parseRequiredVersion(s string) (reqs []types.DependencyExpression) {
	split := strings.Split(s, ",")
	for _, expr := range split {
		expr = strings.TrimSpace(expr)

		if expr == "*" {
			// No specific requirement
			continue
		}

		version := strings.TrimLeft(expr, "<>=!^~")
		req := types.DependencyExpression{
			Value: dependency.Parse(
				types.RawVersion(version),
				types.Semver,
			),
		}

		// Currently, I did not see the x.x.x or *.*.* pattern in MCDR's plugin
		// requirements, so I will not implement it for now.
		if strings.HasPrefix(expr, "=") {
			req.Operator = types.OpEq
		} else if strings.HasPrefix(expr, "~") ||
			strings.HasPrefix(expr, "~=") {
			req.Operator = types.OpWeakEq
		} else if strings.HasPrefix(expr, "<=") {
			req.Operator = types.OpLe
		} else if strings.HasPrefix(expr, "<") {
			req.Operator = types.OpLt
		} else if strings.HasPrefix(expr, ">=") {
			req.Operator = types.OpGeq
		} else if strings.HasPrefix(expr, ">") {
			req.Operator = types.OpGt
		} else if strings.HasPrefix(expr, "!=") {
			req.Operator = types.OpNeq
		} else if strings.HasPrefix(expr, "^") {
			req.Operator = types.OpWeakGt
		} else {
			req.Operator = types.OpEq
		}

		reqs = append(reqs, req)
	}
	return
}
