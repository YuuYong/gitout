package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// 调试模式
var debug bool

func main() {

	// 接收命令行参数
	var dir string
	var outDir string
	var version string
	var currentId string
	var lastId string
	flag.BoolVar(&debug, "debug", false, "Print debug info")
	flag.StringVar(&dir, "dir", "D:/www/XXX", "Your program dir")
	flag.StringVar(&outDir, "outDir", "./", "Output dir,default current dir")
	flag.StringVar(&version, "version", "XXX_1.0.0", "Your version like XXX_1.0.0")
	flag.StringVar(&currentId, "currentId", "HEAD", "Current commit id, default HEAD")
	flag.StringVar(&lastId, "lastId", "HEAD~1", "Last version commit id, default HEAD~1")

	// 解析命令行标志
	flag.Parse()

	// git项目所在目录，如："D:/www/XXX"
	programDir := dir
	if _, err := os.Stat(programDir); os.IsNotExist(err) {
		fmt.Printf("Program Dir: %s not exists.\n", programDir)
		return
	}

	// 目录所在磁盘
	disk := string(programDir[0]) + ":"

	// 输出目录，如："C:/Users/Admin/Desktop/XXX_V1.0.5"
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		fmt.Printf("Out Dir: %s not exists.\n", outDir)
		return
	}
	destDir := strings.TrimRight(outDir, "/") + "/" + version

	// zip文件路径，如："C:/Users/Admin/Desktop/XXX_V1.0.5.zip"
	destFilePath := strings.TrimRight(outDir, "/") + "/" + version + ".zip"

	// 当前版本commit-id，如："HEAD"
	currentVer := currentId

	// 对比版本commit-id，如："HEAD~1"
	lastVer := lastId

	// Added (A), Copied (C), Deleted (D), Modified (M), Renamed (R),
	// have their type (i.e. regular file, symlink, submodule, …) changed (T),
	// are Unmerged (U), are Unknown (X), or have had their pairing Broken (B).
	// 小写字母含义为展示所有非指定类型的变化，这里打增量包只需要排查删除的文件
	programExecArgs := []string{"/C", disk + " && cd " + programDir + " && git diff " + lastVer + " " + currentVer + " --name-only --diff-filter=d"}

	if debug {
		fmt.Println("[debug] cmd params:\n", programExecArgs)
	}

	// 获取需要打包的文件列表
	cmd := exec.Command("cmd.exe", programExecArgs...)
	output, _ := cmd.Output()

	if debug {
		fmt.Println("[debug] cmd exec result:\n", string(output))
	}

	// 获取完整路径
	fileList := getFullFileList(string(output), programDir)

	if debug {
		fmt.Println("[debug] cmd result resolve:\n", fileList)
	}

	if len(fileList) == 0 {
		fmt.Println("Resolve change files error.")
		return
	}
	fmt.Println("Resolve change files full path successful.")

	// 导出到指定目录
	exportFiles(programDir, destDir, fileList)

	fmt.Println("Export successful. Saved in: ", destDir)

	// 压缩为zip
	createZip(destDir, destFilePath)
	fmt.Println("Compress successful. Saved in: ", destFilePath)
}

// 将 git diff 输出的文件路径转为可用切片
func getFullFileList(output string, programDir string) []string {
	// 按行转切片
	fileList := strings.Split(output, "\n")

	// 清理空字符串，同时转义unicode字符
	for i := 0; i < len(fileList); i++ {
		if fileList[i] == "" {
			// 删除空字符串
			fileList = append(fileList[:i], fileList[i+1:]...)
		} else if strings.Contains(fileList[i], "\"") {
			// 包含`"`的是经过 Unicode 转义的中文字符串
			fileList[i], _ = strconv.Unquote(fileList[i])
			fileList[i] = strings.TrimRight(programDir, "/") + "/" + fileList[i]
		} else {
			// 补充完整路径
			fileList[i] = strings.TrimRight(programDir, "/") + "/" + fileList[i]
		}
	}
	return fileList
}

// 导出文件到指定目录
func exportFiles(programDir, destDir string, fileList []string) bool {
	// 遍历文件路径列表，将每个文件移动到指定的目标文件夹中
	for _, filePath := range fileList {
		// 获取相对于programDir的路径
		relPath, err := filepath.Rel(programDir, filePath)
		if err != nil {
			panic(err)
		}

		// 拼接目标文件夹中的路径
		destPath := filepath.Join(destDir, relPath)

		// 创建目标文件夹中的目录（如果不存在）
		err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		if err != nil {
			panic(err)
		}

		// 打开原文件
		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// 创建目标文件
		destFile, err := os.Create(destPath)
		if err != nil {
			panic(err)
		}
		defer destFile.Close()

		// 将原文件的内容拷贝到目标文件中
		_, err = io.Copy(destFile, file)
		if err != nil {
			panic(err)
		}

		// 获取文件名
		if debug {
			fileName := filepath.Base(filePath)
			fmt.Printf("[debug] 文件[%s] => [%s]\n", fileName, destPath)
		}

	}

	return true
}

func createZip(originDir string, destFilePath string) {
	// 创建zip文件
	zipFile, err := os.Create(destFilePath)
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	// 创建zip写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	//根目录
	rootDir := ""
	parts := strings.Split(originDir, "/")
	if len(parts) > 1 {
		rootDir = parts[len(parts)-1] + "/"
	}

	// 遍历目录中的文件和子目录，将每个文件和子目录加入到zip文件中
	filepath.Walk(originDir, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果是目录则跳过
		if fileInfo.IsDir() {
			return nil
		}

		// 创建zip文件中的文件
		zipFile, err := zipWriter.Create(rootDir + relPath(filePath, originDir))
		if err != nil {
			return err
		}

		// 打开原文件
		fileToZip, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		// 将原文件内容拷贝到zip文件中
		_, err = io.Copy(zipFile, fileToZip)
		if err != nil {
			return err
		}

		return nil
	})
}

// 计算相对于目录的路径
func relPath(filePath string, dirPath string) string {
	relPath, err := filepath.Rel(dirPath, filePath)
	if err != nil {
		panic(err)
	}
	return relPath
}
