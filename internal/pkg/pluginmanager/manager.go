package pluginmanager

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/merico-dev/stream/internal/pkg/configloader"
)

func DownloadPlugins(conf *configloader.Config) error {
	// create plugins dir if not exist
	pluginDir := viper.GetString("plugin-dir")
	if pluginDir == "" {
		return fmt.Errorf("plugins directory should not be \"\"")
	}
	log.Printf("Using dir <%s> to store plugins.", pluginDir)

	// download all plugins that don't exist locally
	dc := NewPbDownloadClient()

	for _, tool := range conf.Tools {
		pluginFileName := configloader.GetPluginFileName(&tool)
		if _, err := os.Stat(filepath.Join(pluginDir, pluginFileName)); errors.Is(err, os.ErrNotExist) {
			// plugin does not exist
			err := dc.download(pluginDir, pluginFileName, tool.Version)
			if err != nil {
				return err
			}
			continue
		}
		// check md5
		if dup, _ := checkFileMD5(filepath.Join(pluginDir, pluginFileName), dc, pluginFileName, tool.Version); dup {
			err := dc.download(pluginDir, pluginFileName, tool.Version)
			if err != nil {
				return err
			}
			continue
		}
		log.Printf("Plugin: %s already exists, no need to download.", pluginFileName)

	}

	return nil
}

func checkFileMD5(file string, dc *PbDownloadClient, pluginFileName string, version string) (bool, error) {
	localmd5, err := LocalContentMD5(file)
	if err != nil {
		return false, err
	}
	remotemd5, err := dc.fetchContentMD5(pluginFileName, version)
	if err != nil {
		return false, err
	}

	if strings.Compare(localmd5, remotemd5) == 0 {
		return true, nil
	}
	return false, nil
}
