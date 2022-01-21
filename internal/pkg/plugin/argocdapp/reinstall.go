package argocdapp

import (
	"fmt"
	"os"

	"github.com/merico-dev/stream/internal/pkg/log"

	"github.com/mitchellh/mapstructure"
)

// Reinstall an ArgoCD app
func Reinstall(options *map[string]interface{}) (bool, error) {
	var param Param
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		return false, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s", e)
		}
		return false, fmt.Errorf("params are illegal")
	}

	file := defaultYamlPath

	// delete resource
	err = kubectlAction(ActionDelete, file)
	if err != nil {
		return false, err
	}

	// remove app.yaml file
	if err = os.Remove(file); err != nil {
		return false, err
	}

	// recreate app.yaml file
	err = writeContentToTmpFile(file, appTemplate, &param)
	if err != nil {
		return false, err
	}

	err = kubectlAction(ActionApply, file)
	if err != nil {
		return false, err
	}

	return true, nil
}
