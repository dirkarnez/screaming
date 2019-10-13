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
				if written >= res.ContentLength || err != nil{
					return true // continue write
				} else  {
					return false // close and flush
				}
			})
		}
	})

	//type messageNumber struct {
	//	Number int `json:"number"`
	//}
	//
	//app.Get("/alternative", func(ctx iris.Context) {
	//	ctx.ContentType("application/json")
	//	ctx.Header("Transfer-Encoding", "chunked")
	//	i := 0
	//	ints := []int{1, 2, 3, 5, 7, 9, 11, 13, 15, 17, 23, 29}
	//	// Send the response in chunks and wait for half a second between each chunk.
	//	for {
	//		ctx.JSON(messageNumber{Number: ints[i]})
	//		ctx.WriteString("\n")
	//		time.Sleep(500 * time.Millisecond) // simulate delay.
	//		if i == len(ints)-1 {
	//			break
	//		}
	//		i++
	//		ctx.ResponseWriter().Flush()
	//	}
	//})

	app.Run(iris.Addr(":8080"))
}