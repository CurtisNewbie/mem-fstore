package main

import (
	"bytes"
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
	server.RawPut("/file", func(c *gin.Context, ec common.ExecContext) {
		headers := c.Request.Header
		file := headers["File"][0]
		mu.Lock()
		defer mu.Unlock()

		dat, e := io.ReadAll(c.Request.Body)
		ec.Log.Infof("len: %v", len(dat))
		if e != nil {
			ec.Log.Errorf("io.ReadAll, %v", e)
			return
		}
		mem_store[file] = dat
	})
	server.RawGet("/file", func(c *gin.Context, ec common.ExecContext) {
		file := c.Query("file")
		mu.RLock()
		defer mu.RUnlock()
		if dat, ok := mem_store[file]; ok {
			c.Writer.Header().Set("Content-Disposition", `attachment; filename=`+url.QueryEscape(file))
			c.Writer.Header().Set("Content-Length", strconv.FormatInt(int64(len(dat)), 10))
			if _, e := io.Copy(c.Writer, bytes.NewReader(dat)); e != nil {
				ec.Log.Errorf("c.Writer.Write, %v", e)
				return
			}
		}
	})
	server.BootstrapServer(os.Args)
}
