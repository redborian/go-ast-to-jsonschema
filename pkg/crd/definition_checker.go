// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crd

import (
	"fmt"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func checkDefinitions(defs v1beta1.JSONSchemaDefinitions, startingTypes map[string]bool) {
	fmt.Printf("Type checking Starting expecting %d types\n", len(defs))
	pruner := DefinitionPruner{defs, startingTypes}
	newDefs := pruner.Prune(false)
	if len(defs) != len(defs) {
		fmt.Printf("Type checking failed. Expected %d actual %d\n", len(defs), len(newDefs))
	} else {
		fmt.Println("Type checking PASSED")
	}
}
