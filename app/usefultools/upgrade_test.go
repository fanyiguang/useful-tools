package usefultools

import "testing"

func TestDownloadPkg(t *testing.T) {
	url := "https://gitee.com/qingdeng_ancient_wine/useful-tools-upgrade/releases/download/v2.1.0/2_useful-tools.zip"
	err := DownloadPkg(url, "2_useful-tools.zip")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("download success")
}
