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

package main

import (
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

// Recursively flattens "allOf" tags. If there is cyclic
// dependency, execution is aborted.
func recursiveFlatten(schema *v1beta1.JSONSchemaProps, definition *v1beta1.JSONSchemaProps, defName string, visited *map[string]bool) *v1beta1.JSONSchemaProps {
	if len(definition.AllOf) == 0 {
		return definition
	}
	isAlreadyVisited := (*visited)[defName]
	if isAlreadyVisited {
		panic("Cycle detected in definitions")
	}
	(*visited)[defName] = true

	aggregatedDef := new(v1beta1.JSONSchemaProps)
	for _, allOfDef := range definition.AllOf {
		var newDef *v1beta1.JSONSchemaProps
		if allOfDef.Ref != nil && len(*allOfDef.Ref) > 0 {
			// If the definition has $ref url, fetch the referred resource
			// after flattening it.
			nameOfRef := getNameFromURL(*allOfDef.Ref)
			def := schema.Definitions[nameOfRef]
			newDef = recursiveFlatten(schema, &def, nameOfRef, visited)
		} else {
			newDef = &allOfDef
		}
		mergeDefinitions(aggregatedDef, newDef)
	}

	delete(*visited, defName)
	return aggregatedDef
}

// Merges the properties from the 'rhsDef' to the 'lhsDef'.
// Also transfers the description as well.
func mergeDefinitions(lhsDef *v1beta1.JSONSchemaProps, rhsDef *v1beta1.JSONSchemaProps) {
	if lhsDef == nil || rhsDef == nil {
		return
	}
	// At this point, both defs will not have any 'AnyOf' defs.
	// 1. Add all the properties from rhsDef to lhsDef
	if lhsDef.Properties == nil {
		lhsDef.Properties = make(map[string]v1beta1.JSONSchemaProps)
	}
	for propKey := range rhsDef.Properties {
		lhsDef.Properties[propKey] = rhsDef.Properties[propKey]
	}
	// 2. Transfer the description
	lhsDef.Description = rhsDef.Description
}

// Flattens the schema by inlining 'allOf' tags.
func flattenAllOf(schema *v1beta1.JSONSchemaProps) {
	for nameOfDef := range schema.Definitions {
		visited := make(map[string]bool)
		def := schema.Definitions[nameOfDef]
		schema.Definitions[nameOfDef] = *recursiveFlatten(schema, &def, nameOfDef, &visited)
	}
}
