package upgrade

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"github.com/astaxie/beego/utils"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"useful-tools/common/config"
)

var (
	bakDir string
	runDir string
	upDir  string
)

func Upgrade(file string, processName string) error {
	if !utils.FileExists(file) {
		return fmt.Errorf("%v file not exist", file)
	}
	upDir = filepath.Join(config.GetProjectsPath(), "useful-tools_new")
	if !utils.FileExists(upDir) {
		_ = os.MkdirAll(upDir, 0766)
	}
	defer func() {
		err := os.RemoveAll(upDir)
		if err != nil {
			logrus.Warnf("remove dir error: %v", err)
		}
		err = os.Remove(file)
		if err != nil {
			logrus.Warnf("remove file error: %v", err)
		}
	}()
	err := Unzip(file, upDir)
	if err != nil {
		return err
	}
	logrus.Infof("unzip dir: %v", upDir)
	proc, err := process.NewProcess(int32(os.Getppid()))
	if err != nil {
		return err
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	err = proc.KillWithContext(timeout)
	if err != nil {
		return err
	}
	logrus.Infof("kill process: %v", proc.Pid)
	time.Sleep(time.Second)

	runDir, err = readRunDir()
	if err != nil {
		return err
	}
	if !utils.FileExists(runDir) {
		return fmt.Errorf("%v runDir not exist", file)
	}

	bakDir = fmt.Sprintf("%v_bak", runDir)
	err = CopyDir(runDir, bakDir)
	if err != nil {
		return err
	}
	defer os.RemoveAll(bakDir)
	logrus.Infof("copy run dir: %v", runDir)

	err = CopyDir(upDir, runDir)
	if err != nil {
		resurrection()
		return err
	}
	logrus.Infof("copy new dir: %v", upDir)

	err = runProc(filepath.Join(runDir, fmt.Sprintf("%v.exe", processName)))
	if err != nil {
		return err
	}

	logrus.Infof("run new process: %v", filepath.Join(runDir, fmt.Sprintf("%v.exe", processName)))
	return nil
}

func runProc(filename string) error {
	command := exec.Command(filename)
	err := command.Start()
	if err != nil {
		return err
	}
	return nil
}

func resurrection() {
	err := CopyDir(bakDir, runDir)
	if err != nil {
		logrus.Warnf("resurrection copy dir error: %v", err)
	}
	_ = os.RemoveAll(bakDir)
}

func readRunDir() (string, error) {
	open, err := os.Open(filepath.Join(config.GetConfigPath(), "path"))
	if err != nil {
		return "", err
	}
	defer open.Close()
	all, err := io.ReadAll(open)
	if err != nil {
		return "", err
	}
	return string(all), nil
}

func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CopyDir(from, to string) error {
	f, e := os.Stat(from)
	if e != nil {
		return e
	}
	if f.IsDir() {
		//from是文件夹，那么定义to也是文件夹
		if list, e := ioutil.ReadDir(from); e == nil {
			for _, item := range list {
				if e = CopyDir(filepath.Join(from, item.Name()), filepath.Join(to, item.Name())); e != nil {
					return e
				}
			}
		}
	} else {
		//from是文件，那么创建to的文件夹
		p := filepath.Dir(to)
		if _, e = os.Stat(p); e != nil {
			if e = os.MkdirAll(p, 0777); e != nil {
				return e
			}
		}
		//读取源文件
		file, e := os.Open(from)
		if e != nil {
			return e
		}
		defer file.Close()
		bufReader := bufio.NewReader(file)
		// 创建一个文件用于保存
		out, e := os.Create(to)
		if e != nil {
			return e
		}
		defer out.Close()
		// 然后将文件流和文件流对接起来
		_, e = io.Copy(out, bufReader)
	}
	return e
}
