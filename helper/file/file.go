package file

import (
	"bufio"
	"io"
	"os"
)

func CopyFile(dstFileName string, srcFileName string) (written int64, err error) {
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	//通过srcFile，获取到READER
	reader := bufio.NewReader(srcFile)

	//打开dstFileName
	dstFile, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}

	//通过dstFile，获取到WRITER
	writer := bufio.NewWriter(dstFile)
	//writer.Flush()

	defer dstFile.Close()

	return io.Copy(writer, reader)
}
