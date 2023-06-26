package main

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/gin-gonic/gin"
)

var (
	mu        sync.RWMutex
	mem_store map[string][]byte = make(map[string][]byte)
)

func main() {
	server.RawPost("/file/:file", func(c *gin.Context, ec common.ExecContext) {
		file := c.Param("file")
		mu.Lock()
		defer mu.Unlock()

		dat, e := io.ReadAll(c.Request.Body)
		if e != nil {
			ec.Log.Errorf("Read data, %v", e)
			return
		}
		mem_store[file] = dat

		url := fmt.Sprintf("http://%s:%s/file/%s",
			common.GetLocalIPV4(),
			common.GetPropStr(common.PROP_SERVER_PORT),
			url.QueryEscape(file),
		)
		ec.Log.Infof("File: %v, bytes: %v, url: '%v'", file, len(dat), url)
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
	server.BootstrapServer(os.Args)
}
