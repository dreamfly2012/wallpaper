package main

import (
	"crypto/md5"

	"encoding/hex"
	"encoding/json"

	"errors"

	"fmt"

	"io/ioutil"

	"net/http"

	"os"

	"path/filepath"

	"regexp"

	"strings"

	"syscall"

	"time"

	"unsafe"
)

const (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36"

	BingHomeURL = "https://wallhaven.cc/api/v1/search?sorting=toplist"

	CurrentPathDir = "cache/"
)

const (
	Size1k string = "1920,1080"

	Size2k string = "2560,1440"

	Size4k string = "3840,2160"
)

// id": "gp8pdq",
// "url": "https://wallhaven.cc/w/gp8pdq",
// "short_url": "https://whvn.cc/gp8pdq",
// "views": 71682,
// "favorites": 1459,
// "source": "https://www.artstation.com/artwork/Omr2kJ",
// "purity": "sfw",
// "category": "anime",
// "dimension_x": 5760,
// "dimension_y": 2880,
// "resolution": "5760x2880",
// "ratio": "2",
// "file_size": 1878740,
// "file_type": "image/jpeg",
// "created_at": "2022-11-15 00:35:19",
// "colors": [
// "#424153",
// "#999999",
// "#000000",
// "#996633",
// "#333399"
// ],
// "path": "https://w.wallhaven.cc/full/gp/wallhaven-gp8pdq.jpg",
// "thumbs": {
// "large": "https://th.wallhaven.cc/lg/gp/gp8pdq.jpg",
// "original": "https://th.wallhaven.cc/orig/gp/gp8pdq.jpg",
// "small": "https://th.wallhaven.cc/small/gp/gp8pdq.jpg"
// }

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
	Infos []Info `json:"infos"`
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

// GetBingBackgroundImageURL 获取bing主页的背景图片链接
//result map[string]interface{}

func GetBingBackgroundImageURL() (result string, err error) {

	client := http.Client{}

	request, err := http.NewRequest("GET", BingHomeURL, nil)

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

	return data.Infos[0].Path, nil

}

// DownloadImage 下载图片,保存并返回保存的文件名的绝对路径

func DownloadImage(imageURL string, size *ImageSize) (string, error) {

	wRegexp := regexp.MustCompile("w=\\d+")

	hRegexp := regexp.MustCompile("h=\\d+")

	imageURL = wRegexp.ReplaceAllString(imageURL, "w="+size.w)

	imageURL = hRegexp.ReplaceAllString(imageURL, "h="+size.h)

	client := http.Client{}

	fmt.Println(imageURL)

	request, err := http.NewRequest("GET", imageURL, nil)

	if err != nil {

		return "", err

	}

	response, err := client.Do(request)

	if err != nil {

		return "", err

	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {

		return "", err

	}

	day := time.Now().Format("2006-01-02")

	fileName := EncodeMD5(imageURL)

	path := CurrentPathDir + fmt.Sprintf("[%sx%s][%s]%s", size.w, size.h, day, fileName) + ".jpg"

	err = ioutil.WriteFile(path, body, 0755)

	if err != nil {

		return "", err

	}

	absPath, err := filepath.Abs(path)

	if err != nil {

		return "", err

	}

	return absPath, nil

}

func main() {

	fmt.Println("获取必应背景图中...")

	imageURL, err := GetBingBackgroundImageURL()

	if err != nil {

		fmt.Println("获取背景图片链接失败: " + err.Error())

		return

	}

	fmt.Println("获取成功: " + imageURL)

	fmt.Println("下载图片...")

	imagePath, err := DownloadImage(imageURL, &ImageSize{

		w: strings.Split(Size4k, ",")[0],

		h: strings.Split(Size4k, ",")[1],
	})

	if err != nil {

		fmt.Println("下载图片失败: " + err.Error())

		return

	}

	fmt.Println("设置桌面...")

	err = SetWindowsWallpaper(imagePath)

	if err != nil {

		fmt.Println("设置桌面背景失败: " + err.Error())

		return

	}

}
