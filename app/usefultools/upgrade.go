package usefultools

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"useful-tools/common/config"
	"useful-tools/utils"
)

type UpgradeParam struct {
	PkgDownloadURL string `json:"pkg_download_url"`
	ProcessName    string `json:"process_name"`
	ZipPkgName     string `json:"zip_pkg_name"`
	Version        string `json:"version"`
	Goos           string `json:"goos"`
	Arch           string `json:"arch"`
}

func upgrade() error {
	resp, err := http.Get(getVersionUrl())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logrus.Debugf("upgrade content: %v", string(content))
	var upgradeParam UpgradeParam
	err = json.Unmarshal(content, &upgradeParam)
	if err != nil {
		return err
	}
	logrus.Infof("upgrade param: %+v", upgradeParam)
	if !isUpgrade(upgradeParam, config.Version) {
		return nil
	}
	filename := filepath.Join(os.TempDir(), fmt.Sprintf("useful-tools_%v.zip", upgradeParam.Version))
	downloadUrl := buildDownloadUrl(upgradeParam.Version, upgradeParam.PkgDownloadURL, upgradeParam.ZipPkgName)
	logrus.Infof("download url: %v", downloadUrl)
	err = DownloadPkg(downloadUrl, filename)
	if err != nil {
		return err
	}
	logrus.Infof("download pkg success: %v", filename)

	cmd := exec.Command(filepath.Join(config.GetProjectsPath(), config.ProcessUpgradeName), filename, upgradeParam.ProcessName)
	err = cmd.Run()
	if err != nil {
		logrus.Warnf("cmd run upgrade error: %v", err)
	}
	return nil
}

func getVersionUrl() string {
	if config.IsTest() {
		return config.VersionURL + "/version_test"
	} else {
		return config.VersionURL + "/version_release"
	}
}

func buildDownloadUrl(version, pkgDownloadURL, zipPkgName string) string {
	return fmt.Sprintf("%v/v%v/%v", pkgDownloadURL, version, fmt.Sprintf("%v_%v_%v", runtime.GOOS, runtime.GOARCH, zipPkgName))
}

func isUpgrade(upgradeParam UpgradeParam, oldVersion string) bool {
	if !utils.InArray(runtime.GOOS, strings.Split(upgradeParam.Goos, ",")) {
		return false
	}
	if !utils.InArray(runtime.GOARCH, strings.Split(upgradeParam.Arch, ",")) {
		return false
	}
	newVersion := cast.ToString(strings.ReplaceAll(upgradeParam.Version, ".", ""))
	oldVersion = cast.ToString(strings.ReplaceAll(oldVersion, ".", ""))
	if newVersion <= oldVersion {
		return false
	}
	return true
}

func DownloadPkg(url string, filename string) error {
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code error: %d", resp.StatusCode)
	}

	zipfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, _ = io.Copy(zipfile, resp.Body)
	return nil
}
