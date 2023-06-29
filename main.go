package main

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

var (
	mu        sync.RWMutex
	mem_store map[string][]byte = make(map[string][]byte)
	indexHtml []byte            = nil
)

func main() {

	var err error
	indexHtml, err = os.ReadFile("./index.html")
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	if indexHtml != nil {
		server.RawGet("/", func(c *gin.Context, ec common.ExecContext) {
			c.Header("Content-Type", "text/html")
			_, _ = c.Writer.Write(indexHtml)
		})
	}

	server.Get("/file/list", func(c *gin.Context, ec common.ExecContext) (any, error) {
		mu.RLock()
		defer mu.RUnlock()
		var keys []string = []string{}
		for k := range mem_store {
			keys = append(keys, k)
		}
		return keys, nil
	})

	server.RawPost("/file/:file", func(c *gin.Context, ec common.ExecContext) {
		file := c.Param("file")
		ec.Log.Infof("Reading data for file, %v", file)

		start := time.Now()
		dat, e := io.ReadAll(c.Request.Body)
		if e != nil {
			ec.Log.Errorf("Read data, %v", e)
			return
		}
		took := time.Since(start)

		mu.Lock()
		defer mu.Unlock()
		mem_store[file] = dat

		url := fmt.Sprintf("http://%s:%s/file/%s",
			common.GetLocalIPV4(),
			common.GetPropStr(common.PROP_SERVER_PORT),
			url.QueryEscape(file),
		)
		ec.Log.Infof("File: %v, bytes: %v, url: '%v', took: %v", file, len(dat), url, took)
		c.Data(200, "text/plain", []byte(url))
	})

	server.RawGet("/file/:file", func(c *gin.Context, ec common.ExecContext) {
		file := c.Param("file")
		mu.RLock()
		defer mu.RUnlock()

		if dat, ok := mem_store[file]; ok {
			c.Writer.Header().Set("Content-Disposition", `attachment; filename=`+url.QueryEscape(file))
			c.Writer.Header().Set("Content-Length", strconv.FormatInt(int64(len(dat)), 10))
			if _, e := io.Copy(c.Writer, bytes.NewReader(dat)); e != nil {
				ec.Log.Errorf("Write data, %v", e)
				return
			}
		}
	})

	server.PostServerBootstrapped(func(c common.ExecContext) error {
		c.Log.Infof("Upload file using cURL: 'curl 'http://%s:%s/file/YOUR_FILE_NAME' --data-binary @YOUR_FILE_NAME'",
			common.GetLocalIPV4(), common.GetPropStr(common.PROP_SERVER_PORT))
		return nil
	})

	server.BootstrapServer(os.Args)
}
