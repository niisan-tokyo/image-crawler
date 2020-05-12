package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/rakyll/statik/fs"

	//"fmt"
	"io/ioutil"
	"os"

	"math/rand"
	"mime"
	"strconv"

	_ "github.com/niisan-tokyo/image-crawler/statik"
)

type imageLinks struct {
	Links []string
}

func main() {
	r := gin.Default()

	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "/form.tmpl", gin.H{})
	})

	// スクレイピングしたものを表示
	r.GET("scrape", func(c *gin.Context) {
		url := c.Query("url")
		p := &imageLinks{Links: []string{}}

		coll := colly.NewCollector()
		coll.OnHTML("img", func(e *colly.HTMLElement) {
			src := e.Request.AbsoluteURL(getImgSrc(e))
			if src != "" {
				p.Links = append(p.Links, src)
			}
		})
		coll.Visit(url)
		c.HTML(http.StatusOK, "/scrape.tmpl", gin.H{
			"links": p.Links,
		})
	})

	// 画像を保存して戻る
	r.POST("save", func(c *gin.Context) {
		SaveImage(c.PostFormMap("urls"))
		c.HTML(http.StatusOK, "/complete.tmpl", gin.H{})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	// 仮想ファイルシステム上を走査する
	err = fs.Walk(statikFS, "/", func(path string, info os.FileInfo, err error) error {
		// ディレクトリはスキップ
		if info.IsDir() {
			return nil
		}
		r, err := statikFS.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()

		// データ読み出し
		h, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		// テンプレートにpath名で読み出したデータをパースして格納
		t, err = t.New(path).Parse(string(h))
		if err != nil {
			return err
		}

		return nil
	})

	return t, err
}

func getImgSrc(e *colly.HTMLElement) string {
	src := e.Attr("data-lazy-src")
	if src != "" {
		return src
	}

	src2 := e.Attr("data-wpfc-original-src")
	if src2 != "" {
		return src2
	}

	return e.Attr("src")
}

func SaveImage(imgs map[string]string) {
	os.Mkdir("dist", 0777)
	random := "dist/" + RandomString(10)
	i := 0
	for _, val := range imgs {
		i++
		response, err := http.Get(val)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		ext, _ := mime.ExtensionsByType(response.Header.Get("Content-Type"))
		name := random + strconv.Itoa(i) + ext[0]

		file, err := os.Create(name)
		if err != nil {
			panic(err)
		}
		io.Copy(file, response.Body)
	}
}

const randomLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randomLetters[rand.Intn(len(randomLetters))]
	}
	return string(b)
}
