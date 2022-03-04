package configloader

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"
)

func validateConfig(config *Config) []error {
	errors := make([]error, 0)

	for _, t := range config.Tools {
		errors = append(errors, validateTool(&t)...)
	}

	errors = append(errors, validateDependency(config.Tools)...)

	return errors
}

func validateTool(t *Tool) []error {
	errors := make([]error, 0)

	// Name
	if t.Name == "" {
		errors = append(errors, fmt.Errorf("name is empty"))
	}

	errs := validation.IsDNS1123Subdomain(t.Name)
	for _, e := range errs {
		errors = append(errors, fmt.Errorf("name %s is invalid: %s", t.Name, e))
	}

	// Plugin
	if t.Plugin.Kind == "" {
		errors = append(errors, fmt.Errorf("plugin.kind is empty"))
	}

	return errors
}

func validateDependency(tools []Tool) []error {
	errors := make([]error, 0)

	// config "set" (map)
	toolMap := make(map[string]bool)
	// creating the set
	for _, tool := range tools {
		key := fmt.Sprintf("%s.%s", tool.Name, tool.Plugin.Kind)
		toolMap[key] = true
	}

	for _, tool := range tools {
		// no dependency, pass
		if tool.DependsOn == "" {
			continue
		}

		dependencies := strings.Split(tool.DependsOn, ",")
		// for each dependency
		for _, dependency := range dependencies {
			dependency = strings.TrimSpace(dependency)
			// skip empty string
			if dependency == "" {
				continue
			}

			// generate an error if the dependency isn't in the config set,
			if _, ok := toolMap[dependency]; !ok {
				errors = append(errors, fmt.Errorf("tool %s's dependency %s doesn't exist in the config", tool.Name, dependency))
			}
		}
	}

	return errors
}
