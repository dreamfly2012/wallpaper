package main

import (
	"crypto/md5"
	"log"
	"strings"

	"encoding/hex"
	"encoding/json"

	"errors"

	"fmt"

	"io/ioutil"

	"net/http"

	"os"

	"path/filepath"

	"regexp"

	"syscall"

	"time"

	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

const (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36"

	ApiURL = "https://wallhaven.cc/api/v1/search?ratios=16x10&page=12&q=landscape"

	CurrentPathDir = "cache/"
)

const (
	Size1k string = "1920,1080"

	Size2k string = "2560,1440"

	Size4k string = "3840,2160"
)

type Info struct {
	ID         string   `json:"id"`
	URL        string   `json:"url"`
	ShortURL   string   `json:"short_url"`
	Views      int      `json:"views"`
	Favorites  int      `json:"favorites"`
	Source     string   `json:"source"`
	Purity     string   `json:"purity"`
	Category   string   `json:"category"`
	DimensionX int      `json:"dimension_x"`
	DimensionY int      `json:"dimension_y"`
	Resolution string   `json:"resolution"`
	Ratio      string   `json:"ratio"`
	FileSize   int      `json:"file_size"`
	FileType   string   `json:"file_type"`
	CreatedAt  string   `json:"created_at"`
	Colors     []string `json:"colors"`
	Path       string   `json:"path"`
	Thumbs     struct {
		Large    string `json:"large"`
		Original string `json:"original"`
		Small    string `json:"small"`
	} `json:"thumbs"`
}

type Data struct {
	Data []Info `json:"data"`
}

// ImageSize 图片大小
type ImageSize struct {
	w string

	h string
}

func init() {

	_ = os.Mkdir(CurrentPathDir, 0755)

}

// EncodeMD5 MD5编码

func EncodeMD5(value string) string {

	m := md5.New()

	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))

}

// SetWindowsWallpaper 设置windows壁纸

func SetWindowsWallpaper(imagePath string) error {

	dll := syscall.NewLazyDLL("user32.dll")

	proc := dll.NewProc("SystemParametersInfoW")

	_t, _ := syscall.UTF16PtrFromString(imagePath)

	ret, _, _ := proc.Call(20, 1, uintptr(unsafe.Pointer(_t)), 0x1|0x2)

	if ret != 1 {

		return errors.New("系统调用失败")

	}

	return nil

}

func GetImageList() (info []Info) {

	client := http.Client{}

	request, err := http.NewRequest("GET", ApiURL, nil)

	if err != nil {

		panic(err)

	}

	request.Header.Set("user-agent", UserAgent)

	response, err := client.Do(request)

	if err != nil {

		panic(err)

	}

	jsonByte, _ := ioutil.ReadAll(response.Body)

	var data Data
	err = json.Unmarshal(jsonByte, &data)
	if err != nil {
		panic(err)
	}

	return data.Data

}

// DownloadImage 下载图片,保存并返回保存的文件名的绝对路径

func DownloadImage(imageURL string, size *ImageSize, result chan<- string) {

	wRegexp := regexp.MustCompile("w=\\d+")

	hRegexp := regexp.MustCompile("h=\\d+")

	imageURL = wRegexp.ReplaceAllString(imageURL, "w="+size.w)

	imageURL = hRegexp.ReplaceAllString(imageURL, "h="+size.h)

	client := http.Client{}

	fmt.Println(imageURL)

	request, err := http.NewRequest("GET", imageURL, nil)

	if err != nil {

		log.Println(err)

	}

	response, err := client.Do(request)

	if err != nil {

		log.Println(err)

	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {

		log.Println(err)

	}

	day := time.Now().Format("2006-01-02")

	fileName := EncodeMD5(imageURL)

	path := CurrentPathDir + fmt.Sprintf("[%sx%s][%s]%s", size.w, size.h, day, fileName) + ".jpg"

	err = ioutil.WriteFile(path, body, 0755)

	if err != nil {

		log.Println(err)

	}

	absPath, err := filepath.Abs(path)

	if err != nil {

		log.Println(err)

	}

	result <- absPath

}

func draw(infoList []Info) {
	imageList := []fyne.CanvasObject{}

	myApp := app.New()
	myWindow := myApp.NewWindow("桌面背景修改器")

	ch := make(chan string, 12)

	for i := 0; i < 12; i++ {
		//todo download and assign image
		go DownloadImage(infoList[i].Path, &ImageSize{

			w: strings.Split(Size4k, ",")[0],

			h: strings.Split(Size4k, ",")[1],
		}, ch)
	}

	for i := 0; i < 12; i++ {
		path := <-ch
		image := canvas.NewImageFromFile(path)
		image.FillMode = canvas.ImageFillOriginal

		clickImage := NewClickImage()
		clickImage.Image = image
		clickImage.OnTapped = func() {
			setDeskBackgroud(path)
			println("clicked", path)
		}

		imageList = append(imageList, clickImage)
	}

	content1 := container.New(layout.NewGridLayout(4), imageList...)

	myWindow.SetContent(content1)

	myWindow.Resize(fyne.NewSize(800, 400))

	myWindow.SetIcon(theme.FyneLogo())

	myWindow.ShowAndRun()
}

func setDeskBackgroud(imagePath string) {

	err := SetWindowsWallpaper(imagePath)

	if err != nil {

		fmt.Println("设置桌面背景失败: " + err.Error())

		return

	}
}

func main() {
	imageList := GetImageList()

	draw(imageList)

}
