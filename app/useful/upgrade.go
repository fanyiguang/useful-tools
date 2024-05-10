package useful

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"useful-tools/common/config"
	"useful-tools/pkg/wlog"
)

type UpgradeParam struct {
	PkgDownloadURL string `json:"pkg_download_url"`
	ProcessName    string `json:"process_name"`
	ZipPkgName     string `json:"zip_pkg_name"`
	Version        string `json:"version"`
}

func upgrade() error {
	resp, err := http.Get(config.VersionURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var upgradeParam UpgradeParam
	err = json.Unmarshal(content, &upgradeParam)
	if err != nil {
		return err
	}
	wlog.Info("new version: %v, current version: %v", upgradeParam.Version, config.Version)
	if !isUpgrade(upgradeParam.Version, config.Version) {
		return nil
	}

	filename := filepath.Join(os.TempDir(), fmt.Sprintf("useful-tools_%v.zip", upgradeParam.Version))
	downloadUrl := buildDownloadUrl(upgradeParam.Version, upgradeParam.PkgDownloadURL, upgradeParam.ZipPkgName)
	wlog.Info("download url: %v", downloadUrl)
	err = DownloadPkg(downloadUrl, filename)
	if err != nil {
		return err
	}
	wlog.Info("download pkg success")

	cmd := exec.Command(filepath.Join(config.GetProjectsPath(), "upgrade.exe"), filename, upgradeParam.ProcessName)
	err = cmd.Run()
	if err != nil {
		wlog.Warm("cmd run upgrade error: %v", err)
	}
	return nil
}

func buildDownloadUrl(version, pkgDownloadURL, zipPkgName string) string {
	return fmt.Sprintf("%v/v%v/%v", pkgDownloadURL, version, zipPkgName)
}

func isUpgrade(newVersion string, oldVersion string) bool {
	newVersion = cast.ToString(strings.ReplaceAll(newVersion, ".", ""))
	oldVersion = cast.ToString(strings.ReplaceAll(oldVersion, ".", ""))
	if newVersion > oldVersion {
		return true
	}
	return false
}

func DownloadPkg(url string, filename string) error {
	client := &http.Client{Timeout: 2 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	zipfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, _ = io.Copy(zipfile, resp.Body)
	return nil
}
