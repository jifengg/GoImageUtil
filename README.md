# GoImageUtil
A image command line tool util for Golang。

对常用图像处理命令行工具的golang封装。

目前使用[ImageMagick](https://github.com/ImageMagick/ImageMagick)来进行图片的裁剪、格式转换、压缩（非PNG）。
使用[pngquant](https://github.com/kornelski/pngquant)来进行PNG图片压缩。

要使用这个包，请保证系统中已经安装了ImageMagick和pngquant。安装方式请自行查找。


------

# install   安装

## go package

```shell
go get github.com/jifengg/GoImageUtil
```

# usage 使用

可参考[`GoImageUtilDemo`](https://github.com/jifengg/GoImageUtilDemo)；

```go
import (
	iu "github.com/jifengg/GoImageUtil"
)
//修改命令行工具的位置，可以使用相对路径，也可以使用绝对路径。
//Init并非必须要调用的方法，但是可以使用此方法测试命令行工具是否可用。
err := iu.Init(conf)

//获取图片信息
info, err := iu.Info(testFile)
if err != nil {
    fmt.Printf("info error:%s\n", err)
} else {
    fmt.Printf("info:%+v\n", info)
}
/*
{
    "IsKnowImage": true,
    "Width": 573,
    "Height": 573,
    "FileSize": 81166,
    "FilePath": "./test/test_image/me.jpg",
    "Format": "JPEG"
}
*/


//转换一个jpg文件，生成一个300x200，质量参数为60的jpg文件
opt := iu.Option{
    Quality: 60,
    Width:   300,
    Heigth:  200}
succ, err := iu.Convert(testFile, path.Join(tempDir, "out.jpg"), opt)



//转换一个jpg文件成png文件，不修改分辨率，pngquant压缩参数为40-80
opt = iu.Option{
    Quality:       80,
    PngQunlityMin: 40}
//转换一个jpg文件成png文件，不修改分辨率，pngquant压缩参数为40-80
succ, err = iu.Convert(testFile, path.Join(tempDir, "out.png"), opt)

```

# test 测试

```shell
go run test\test.go
```

运行之后，如果没有报错，可以到`test/test_output/`中查看输出的文件。