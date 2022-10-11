package pluginengine

import (
	"errors"
	"fmt"
	"os"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/pluginmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/interact"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Remove(configFile string, continueDirectly bool, isForceDelete bool) error {
	cfg, err := configmanager.NewManager(configFile).LoadConfig()
	if err != nil {
		return err
	}

	if cfg == nil {
		return fmt.Errorf("failed to load the config file")
	}

	if err := file.SetPluginDir(cfg.PluginDir); err != nil {
		log.Errorf("Error: %s.", err)
	}

	err = pluginmanager.CheckLocalPlugins(cfg)
	if err != nil {
		log.Errorf(`Error checking required plugins. Maybe you forgot to run "dtm init" first?`)
		return err
	}

	smgr, err := statemanager.NewManager(*cfg.State)
	if err != nil {
		log.Debugf("Failed to get the manager: %s.", err)
		return err
	}

	changes, err := GetChangesForDelete(smgr, cfg, isForceDelete)

	if err != nil {
		return err
	}
	if len(changes) == 0 {
		log.Info("Nothing needs to be deleted. There is nothing to do.")
		return nil
	}
	for _, change := range changes {
		log.Info(change.Description)
	}

	if !continueDirectly {
		continued := interact.AskUserIfContinue(askUserIfContinue)
		if !continued {
			os.Exit(0)
		}
	}

	errsMap := execute(smgr, changes, true)
	if len(errsMap) != 0 {
		err := errors.New("some error(s) occurred during plugins delete process")
		for k, e := range errsMap {
			log.Infof("%s -> %s.", k, e)
		}
		return err
	}

	log.Success("All plugins deleted successfully.")
	return nil
}
