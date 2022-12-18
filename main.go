package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/kataras/iris"
)

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

func main() {
	app := iris.New()

	// tokens := strings.Split(url, "/")
	// fileName := tokens[len(tokens)-1]

	app.Get("/download.pdf", func(ctx iris.Context) {
		ctx.ContentType("application/octet-stream")
		ctx.Header("Transfer-Encoding", "chunked")

		res, err := http.Get("http://mmrc.amss.cas.cn/tlb/201702/W020170224608149940643.pdf")

		if err != nil {
			log.Fatal("http get error: ", err)
		} else {
			defer res.Body.Close()

			ctx.StreamWriter(func(w io.Writer) bool {
				written, err := io.Copy(w, res.Body)
				if written >= res.ContentLength || err != nil {
					return false // close and flush
				} else {
					return true // continue write
				}
			})
		}
	})

	app.Get("/tasklist", func(ctx iris.Context) {
		ctx.ContentType("text/plain")
		ctx.Header("Transfer-Encoding", "chunked")

		cmd := exec.Command("C:\\Windows\\System32\\tasklist.exe")

		var stdout, stderr []byte
		var errStdout, errStderr error
		stdoutPipe, _ := cmd.StdoutPipe()
		stderrIn, _ := cmd.StderrPipe()
		err := cmd.Start()
		if err != nil {
			log.Fatalf("cmd.Start() failed with '%s'\n", err)
		}

		// cmd.Wait() should be called only after we finish reading
		// from stdoutIn and stderrIn.
		// wg ensures that we finish
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			stdout, errStdout = copyAndCapture(os.Stdout, stdoutPipe)
			wg.Done()
		}()

		//stderr, errStderr =
		copyAndCapture(os.Stderr, stderrIn)

		wg.Wait()

		err = cmd.Wait()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		if errStdout != nil || errStderr != nil {
			log.Fatal("failed to capture stdout or stderr\n")
		}

		outStr, errStr := string(stdout), string(stderr)
		fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

		myReader := strings.NewReader(outStr)
		ctx.StreamWriter(func(w io.Writer) bool {
			written, err := io.Copy(w, myReader)
			if int(written) >= len(outStr) || err != nil {
				return false // close and flush
			} else {
				return true // continue write
			}
		})
	})

	app.Get("/ffmpeg", func(ctx iris.Context) {
		ctx.ContentType("video/mp4")
		ctx.Header("Transfer-Encoding", "chunked")

		cmd := exec.Command("P:\\Downloads\\ffmpeg-2021-10-28-git-e84c83ef98-full_build\\bin\\ffmpeg.exe",
			"-i",
			`https://ewcdn04.nowe.com/session/09-ea4a143a0908c150afdfa113f8005/nowstreamer/getnowmediahls/NNEW00372892/sd/index.m3u8?token=2772bdf2ffc2cd8444ecade15a7e996b_1671374384`,
			"-movflags", "frag_keyframe", "-f", "mp4", "pipe:1")

		var stdout, stderr []byte
		var errStdout, errStderr error
		stdoutPipe, _ := cmd.StdoutPipe()
		stderrIn, _ := cmd.StderrPipe()
		err := cmd.Start()
		if err != nil {
			log.Fatalf("cmd.Start() failed with '%s'\n", err)
		}

		// cmd.Wait() should be called only after we finish reading
		// from stdoutIn and stderrIn.
		// wg ensures that we finish
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			//stdout, errStdout = copyAndCapture(os.Stdout, stdoutPipe)
			defer wg.Done()

			ctx.StreamWriter(func(w io.Writer) bool {
				var out []byte
				buf := make([]byte, 1024, 1024)
				for {
					n, err := stdoutPipe.Read(buf[:])
					if n > 0 {
						d := buf[:n]
						out = append(out, d...)
						_, err := w.Write(d)
						if err != nil {
							return false
						} else {
							return true
						}
					} else {
						if err != nil {
							// Read returns io.EOF at the end of file, which is not an error for us
							if err == io.EOF {
								return false // close and flush
							}
							return false
							//return out, err
						} else {
							return true
						}
					}

				}

				// _, err := io.Copy(w, )
				// if err != nil {
				// 	if err == io.EOF {
				// 		fmt.Println("EOF!!!!!")
				// 		return false // close and flush
				// 	}

				// 	fmt.Println("err != nil!!!!!")
				// 	return true
				// } else {
				// 	return true // continue write
				// }
			})
		}()

		//stderr, errStderr =
		copyAndCapture(os.Stderr, stderrIn)

		wg.Wait()

		err = cmd.Wait()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		if errStdout != nil || errStderr != nil {
			log.Fatal("failed to capture stdout or stderr\n")
		}

		// outStr, errStr := string(stdout), string(stderr)
		// fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
		errStr := string(stderr)
		fmt.Printf("\nlen(out):\n%d\nerr:\n%s\n", len(stdout), errStr)
		//reader := bytes.NewReader(stdout)

	})

	app.Run(iris.Addr(":8080"))
}
