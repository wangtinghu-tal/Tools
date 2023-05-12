package main

import (
	"context"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var OPT = Options{}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SelectFile 选择需要处理的文件
func (a *App) SelectFile(filetype string) []string {
	if filetype == "" {
		filetype = "*.png;*.jpeg;*.jpg;*.bmp;*.gif;*.tif;*.tiff;*.webp"
	}
	selection, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "图片数据",
				Pattern:     filetype,
			},
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	return selection
}

// SelectSavePath 选择保存文件路径
func (a *App) SelectSavePath() string {
	selection, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "选择保存位置",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "文本数据",
				Pattern:     "*.pdf",
			},
		},
	})
	if err != nil {
		return fmt.Sprintf("err %s!", err)
	}
	return selection
}

// 保存图片到本地
func saveImage(path string, src image.Image) error {
	f, err := os.OpenFile(path, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	ext := filepath.Ext(path)
	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {
		err = jpeg.Encode(f, src, &jpeg.Options{Quality: 80})
	} else if strings.EqualFold(ext, ".png") {
		err = png.Encode(f, src)
	} else if strings.EqualFold(ext, ".gif") {
		err = gif.Encode(f, src, &gif.Options{NumColors: 256})
	}
	return err
}

// 读取图片文件数据
func getImageFromFile(filePath string) (img image.Image, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

// 图片拼接pdf
func (a *App) HandleImageToPDF(filePaths []string) string {
	if len(filePaths) == 0 {
		return "0"
	}
	timeStart := time.Now()
	// 创建pdf文档
	pdf := gofpdf.New("L", "pt", "A4", "")

	for _, path := range filePaths {
		fmt.Println(path)
		imgFile, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return "0"
		}
		defer imgFile.Close()

		// 解码图片文件
		img, err := jpeg.Decode(imgFile)
		if err != nil {
			fmt.Println(err)
			return "0"
		}

		// 将图片添加到pdf中
		pdf.AddPage()
		temp, _ := pdf.GetPageSize()
		imgWidth := temp - 20.0
		imgHeight := imgWidth * float64(img.Bounds().Dy()) / float64(img.Bounds().Dx())
		pdf.Image(path, 10.0, 10.0, imgWidth, imgHeight, false, "", 0, "")
	}

	// 保存pdf文件
	err := pdf.OutputFileAndClose(OPT.PDFExportPath)
	if err != nil {
		fmt.Println(err)
	}

	return time.Now().Sub(timeStart).String()
}

// GetOptions 获取配置文件JSON
func (a *App) GetOptions() Options {
	return OPT.GetData()
}

// SetOptions 解析传来的配置
func (a *App) SetOptions(val Options) {
	OPT.SetData(&val)
}

// string to float64
func s2f(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
