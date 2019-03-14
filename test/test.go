package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	iu "githup.com/jifengg/GoImageUtil"
)

func main() {
	// wd, _ := os.Getwd()
	// fmt.Printf("os.Getwd:%s\n", wd)
	var conf iu.Config
	if runtime.GOOS == "linux" {
		conf = iu.Config{
			ImageMagickConvertPath:  "convert",
			ImageMagickIdentifyPath: "identify",
			PngquantPath:            "pngquant",
			ShowDebug:               true,
			ShowError:               true}
	} else {
		conf = iu.Config{
			ImageMagickConvertPath:  "C:\\Program Files\\ImageMagick-7.0.8-Q16\\convert.exe",
			ImageMagickIdentifyPath: "C:\\Program Files\\ImageMagick-7.0.8-Q16\\identify.exe",
			PngquantPath:            "E:\\Program Files (x86)\\pngquant\\pngquant.exe",
			ShowDebug:               true,
			ShowError:               true}
	}
	//修改命令行工具的位置，可以使用相对路径，也可以使用绝对路径。
	//Init并非必须要调用的方法，但是可以使用此方法测试命令行工具是否可用。
	err := iu.Init(conf)
	if err != nil {
		fmt.Printf("GIU ERROR:%s", err)
		return
	}

	var testFile = "./test/test_image/me.jpg"
	var tempDir = "./test/test_output"
	_, err = os.Stat(tempDir)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(tempDir, 0755)
	}

	info, err := iu.Info(testFile)
	if err != nil {
		fmt.Printf("info error:%s\n", err)
	} else {
		fmt.Printf("info:%+v\n", info)
	}

	//转换一个jpg文件，生成一个300x200，质量参数为60的jpg文件
	opt := iu.Option{
		Quality: 60,
		Width:   300,
		Heigth:  200}
	succ, err := iu.Convert(testFile, path.Join(tempDir, "out"+getFileNameSuffix(opt)+".jpg"), opt)
	fmt.Println("return", succ, err)

	//转换一个jpg文件，生成一个宽度300，高度自动等比例缩放，质量参数为60的jpg文件
	opt.Heigth = 0
	succ, err = iu.Convert(testFile, path.Join(tempDir, "out"+getFileNameSuffix(opt)+".jpg"), opt)
	fmt.Println("return", succ, err)

	//转换一个jpg文件为一个bmp文件，不附加任何参数
	opt = iu.Option{}
	succ, err = iu.Convert(testFile, path.Join(tempDir, "out"+getFileNameSuffix(opt)+".bmp"), opt)
	fmt.Println("return", succ, err)

	opt = iu.Option{
		Quality:       75,
		PngQunlityMin: 30}
	//转换一个jpg文件成png文件，不修改分辨率，pngquant压缩参数为30-75
	succ, err = iu.Convert(testFile, path.Join(tempDir, "out"+getFileNameSuffix(opt)+".png"), opt)
	fmt.Println("return", succ, err)

	opt = iu.Option{
		Quality:       80,
		PngQunlityMin: 40}
	//转换一个jpg文件成png文件，不修改分辨率，pngquant压缩参数为40-80
	succ, err = iu.Convert(testFile, path.Join(tempDir, "out"+getFileNameSuffix(opt)+".png"), opt)
	fmt.Println("return", succ, err)

	opt = iu.Option{
		Quality:       80,
		PngQunlityMin: 40,
		Width:         130,
		Heigth:        240}
	//转换一个jpg文件成png文件，分辨率修改为130x240，pngquant压缩参数为40-80
	succ, err = iu.Convert(testFile, path.Join(tempDir, "out"+getFileNameSuffix(opt)+".png"), opt)
	fmt.Println("return", succ, err)

	opt = iu.Option{}
	pngTempFile := path.Join(tempDir, "temp.png")
	iu.Convert(testFile, pngTempFile, opt)
	opt = iu.Option{
		Quality:       90,
		PngQunlityMin: 50,
		Width:         233,
		Heigth:        322}
	//压缩一个png文件，分辨率修改为233x322，pngquant压缩参数为50-90
	succ, err = iu.Convert(pngTempFile, path.Join(tempDir, "out"+getFileNameSuffix(opt)+".png"), opt)
	fmt.Println("return", succ, err)

	fmt.Println("end")
}

func getFileNameSuffix(opt iu.Option) string {
	var name = ""
	if opt.Width > 0 {
		name += fmt.Sprintf(`_w%d`, opt.Width)
	}
	if opt.Heigth > 0 {
		name += fmt.Sprintf(`_h%d`, opt.Heigth)
	}
	if opt.PngQunlityMin > 0 {
		name += fmt.Sprintf(`_qmin%d`, opt.PngQunlityMin)
	}
	if opt.Quality > 0 {
		name += fmt.Sprintf(`_q%d`, opt.Quality)
	}
	return name
}
