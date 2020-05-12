package main

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/rakyll/statik/fs"

	//"fmt"
	"io/ioutil"
	"os"

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
