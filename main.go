/**
使用多个命令行工具，进行图片相关的处理
*/

package goimageutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
)

// VERSION 版本号
var VERSION = "1.0.0"

// imageMagickConvertPath ImageMagick 包的 convert 工具的路径
var imageMagickConvertPath = "convert"

// imageMagickIdentifyPath ImageMagick 包的 identify 工具的路径
var imageMagickIdentifyPath = "identify"

// pngquantPath pngquant命令行工具的路径
var pngquantPath = "pngquant"

// 是否打印出调试信息
var showDebug = false

// 是否打印出错误信息
var showError = false

// JPEG 图片格式JPEG
var JPEG = "JPEG"

// PNG 图片格式PNG
var PNG = "PNG"

// GIF 图片格式GIF
var GIF = "GIF"

// BMP 图片格式BMP
var BMP = "BMP"

// ImageInfo 图片信息
type ImageInfo struct {
	IsKnowImage bool   `json:"isimg"` // 是否是已知的图片格式
	Width       uint   `json:"w"`     // 宽度像素
	Height      uint   `json:"h"`
	FileSize    int64  `json:"size"`
	FilePath    string `json:"path"`
	Format      string `json:"m"`
}

// Option 处理选项
type Option struct {
	Width         uint // 要缩放成的宽度，为0表示默认
	Heigth        uint // 要缩放成的高度，为0表示默认
	Quality       uint // 不等于0时，表示要压缩，设置图片的质量参数，1-100。如果要压缩的是png格式，则这个值表示pngquant的质量参数的最大值
	PngQunlityMin uint // pngquant压缩png图片时，质量参数的最小值
}

// Config 工具的配置信息
type Config struct {
	ImageMagickConvertPath  string //ImageMagick 包的 convert 工具的路径
	ImageMagickIdentifyPath string //ImageMagick 包的 identify 工具的路径
	PngquantPath            string //pngquant命令行工具的路径
	ShowDebug               bool   //是否打印出调试信息
	ShowError               bool   //是否打印出错误信息
}

// Init 进行初始化，检查各个命令是否可用。
func Init(conf Config) error {
	if conf.ImageMagickConvertPath != "" {
		imageMagickConvertPath = conf.ImageMagickConvertPath
	}
	if conf.ImageMagickIdentifyPath != "" {
		imageMagickIdentifyPath = conf.ImageMagickIdentifyPath
	}
	if conf.PngquantPath != "" {
		pngquantPath = conf.PngquantPath
	}
	exitCode, _, err := run(imageMagickIdentifyPath, "--version")
	if exitCode != 0 || err != nil {
		return err
	}
	exitCode, _, err = run(imageMagickConvertPath, "--version")
	if exitCode != 0 || err != nil {
		return err
	}
	exitCode, _, err = run(pngquantPath, "--version")
	if exitCode != 0 || err != nil {
		return err
	}
	showDebug = conf.ShowDebug
	showError = conf.ShowError
	if showDebug {
		fmt.Printf("ImageMagickConvertPath set to :%s\n", imageMagickConvertPath)
		fmt.Printf("ImageMagickIdentifyPath set to :%s\n", imageMagickIdentifyPath)
		fmt.Printf("PngquantPath set to :%s\n", pngquantPath)
		fmt.Printf("ShowError set to :%t\n", showError)
	}
	return nil
}

// run 执行一个外部程序，返回：退出码，输出，异常
func run(name string, arg ...string) (int, []byte, error) {
	if showDebug {
		fmt.Println(name, arg)
	}
	cmd := exec.Command(name, arg...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 1, nil, err
	}
	exitCode := (cmd.ProcessState.Sys().(syscall.WaitStatus)).ExitStatus()

	return exitCode, output, nil
}

// Info 获取图片的相关信息
func Info(file string) (ImageInfo, error) {
	info := ImageInfo{}
	fi, err := os.Stat(file)
	if err == nil {
		format := `{"w":%[width],"h":%[height],"m":"%[magick]"}` //,"size":%[B]}	B=文件字节数，在7.0版本之前不支持
		exitCode, output, runerr := run(imageMagickIdentifyPath, "-format", format, file)
		if exitCode == 0 {
			if showDebug {
				fmt.Printf("get info output:%s\n", output)
			}
			err = json.Unmarshal(output, &info)
			info.FileSize = fi.Size()
			info.IsKnowImage = true
			info.FilePath = file
		} else {
			info.IsKnowImage = false
			err = runerr
		}
	}
	if err != nil && showError {
		fmt.Printf("get info error:%s\n", err)
	}
	return info, err
}

// Convert 图片转换
// fileIn 要处理的文件的绝对路径
func Convert(fileIn string, fileOut string, opt Option) (bool, error) {
	var args []string
	//处理过程中使用的临时文件
	tempFile := fileIn
	outFormat := strings.ToUpper(strings.Replace(path.Ext(fileOut), ".", "", -1))
	info, err := Info(fileIn)
	if err != nil {
		return false, err
	}
	if info.IsKnowImage == false {
		return false, errors.New("UnknowImageFormat")
	}
	//如果要修改分辨率
	if opt.Width > 0 || opt.Heigth > 0 {
		var w, h uint = 0, 0
		if opt.Width == 0 {
			h = opt.Heigth
			w = info.Width * h / info.Height
		} else if opt.Heigth == 0 {
			w = opt.Width
			h = info.Height * w / info.Width
		} else {
			w = opt.Width
			h = opt.Heigth
		}
		//加上感叹号!表示强制转换到这个分辨率，否则会按照等比例缩放
		args = append(args, "-resize", fmt.Sprintf("%dx%d!", w, h))
	}
	//如果需要进行压缩
	if opt.Quality > 0 && outFormat != PNG {
		args = append(args, "-quality", strconv.Itoa(int(opt.Quality)))
	}
	//如果没有特殊参数，判断输入输出文件格式是否相同，如果相同，则不做处理
	if len(args) > 0 || path.Ext(tempFile) != path.Ext(fileOut) {
		args = append(args, tempFile, fileOut)
	}

	//如果需要使用ImageMagick处理
	if len(args) > 0 {
		exitCode, output, err := run(imageMagickConvertPath, args...)
		if exitCode != 0 {
			return false, err
		}
		if showDebug {
			fmt.Printf("ImageMagick exit with code (%d)\noutput: %s\nerror:%s\n", exitCode, output, err)
		}
		//将临时文件路径指向输出文件，便于后续处理
		tempFile = fileOut
	}

	//如果是png，且需要压缩
	if opt.Quality > 0 && outFormat == PNG {
		if opt.PngQunlityMin > opt.Quality {
			opt.PngQunlityMin = opt.Quality
		}
		args = []string{
			"--force", "--quiet", "--ordered", "--speed=1", fmt.Sprintf("--quality=%d-%d", opt.PngQunlityMin, opt.Quality), tempFile, "--output", fileOut}
		exitCode, output, err := run(pngquantPath, args...)
		if showDebug {
			fmt.Printf("Pngquant exit with code (%d)\noutput: %s\nerror:%s\n", exitCode, output, err)
		}
		if exitCode != 0 {
			return false, err
		}
	}

	return true, nil
}
