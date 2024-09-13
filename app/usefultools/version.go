package usefultools

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"useful-tools/common/config"
)

func printVersion() {
	logrus.Infof("version info: %v", versionInfo())
}

func versionInfo() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Version: %s\n", config.Version))
	builder.WriteString(fmt.Sprintf("Runtime: %s - %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH))
	builder.WriteString(fmt.Sprintf("Env: %s\n", config.Env()))
	return builder.String()
}
