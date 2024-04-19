package useful

import (
	"os"
	"path/filepath"
	"useful-tools/common/config"
	"useful-tools/helper/path"
)

func initConfig() error {
	s, err := path.Path()
	if err != nil {
		return err
	}
	create, err := os.Create(filepath.Join(config.GetConfigPath(), "path"))
	if err != nil {
		return err
	}
	_, err = create.WriteString(filepath.Dir(s))
	if err != nil {
		return err
	}
	return nil
}
