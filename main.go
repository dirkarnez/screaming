package main

import (
	"github.com/kataras/iris"
	"io"
	"log"
	"net/http"
)

func main() {
	app := iris.New()
	url := "http://mmrc.amss.cas.cn/tlb/201702/W020170224608149940643.pdf"

	app.Get("/", func(ctx iris.Context) {
		ctx.ContentType("text/html")
		ctx.Header("Transfer-Encoding", "chunked")

		res, err := http.Get(url)

		if err != nil {
			log.Fatal("http get error: ", err)
		} else {
			defer res.Body.Close()

			ctx.StreamWriter(func(w io.Writer) bool {
				written, err := io.Copy(w, res.Body)
				if written >= res.ContentLength || err != nil {
					return false // close and flush
				} else  {
					return true // continue write
				}
			})
		}
	})

	app.Run(iris.Addr(":8080"))
}
