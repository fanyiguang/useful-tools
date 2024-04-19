package useful

import (
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

func upgrade() error {
	resp, err := http.Get(config.VersionURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	version, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	wlog.Info("new version: %v, current version: %v", string(version), config.Version)
	if !isUpgrade(string(version), config.Version) {
		return nil
	}

	filename := filepath.Join(os.TempDir(), fmt.Sprintf("useful-tools_%v.zip", string(version)))
	downloadUrl := buildDownloadUrl(string(version))
	wlog.Info("download url: %v", downloadUrl)
	err = DownloadPkg(downloadUrl, filename)
	if err != nil {
		return err
	}
	wlog.Info("download pkg success")

	cmd := exec.Command(filepath.Join(config.GetProjectsPath(), "upgrade.exe"), filename)
	if err != nil {
		return err
	}
	err = cmd.Run()
	if err != nil {
		wlog.Warm("cmd run upgrade error: %v", err)
	}
	return nil
}

func buildDownloadUrl(version string) string {
	return fmt.Sprintf("%v/v%v/useful-tools.zip", config.PkgDownloadURL, version)
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
