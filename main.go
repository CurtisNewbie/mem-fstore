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

	"github.com/curtisnewbie/mem-store/template"
	"github.com/curtisnewbie/miso/miso"
)

var (
	mu        sync.RWMutex
	mem_store map[string][]byte = make(map[string][]byte)
)

func init() {
	miso.SetDefProp(miso.PropAppName, "mem-store")
}

func main() {

	miso.RawGet("/", func(inb *miso.Inbound) {
		w, _ := inb.Unwrap()
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(template.IndexHtml))
	})

	miso.Get("/file/list", func(inb *miso.Inbound) (any, error) {
		mu.RLock()
		defer mu.RUnlock()
		var keys []string = []string{}
		for k := range mem_store {
			keys = append(keys, k)
		}
		return keys, nil
	})

	miso.RawPost("/file", func(inb *miso.Inbound) {
		rail := inb.Rail()
		file := inb.Query("name")
		rail.Infof("Reading data for file, %v", file)

		start := time.Now()
		w, r := inb.Unwrap()
		dat, e := io.ReadAll(r.Body)
		if e != nil {
			rail.Errorf("Read data, %v", e)
			return
		}
		took := time.Since(start)

		mu.Lock()
		defer mu.Unlock()
		mem_store[file] = dat

		url := fmt.Sprintf("http://%s:%s/file/%s",
			miso.GetLocalIPV4(),
			miso.GetPropStr(miso.PropServerPort),
			url.QueryEscape(file),
		)
		rail.Infof("File: %v, bytes: %v, url: '%v', took: %v", file, len(dat), url, took)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(url))
	})

	miso.RawGet("/file", func(c *miso.Inbound) {
		rail := c.Rail()
		file := c.Query("name")

		mu.RLock()
		defer mu.RUnlock()
		w, _ := c.Unwrap()

		if dat, ok := mem_store[file]; ok {
			w.Header().Set("Content-Disposition", `attachment; filename=`+url.QueryEscape(file))
			w.Header().Set("Content-Length", strconv.FormatInt(int64(len(dat)), 10))
			if _, e := io.Copy(w, bytes.NewReader(dat)); e != nil {
				rail.Errorf("Write data, %v", e)
				return
			}
		}
	})

	miso.PostServerBootstrapped(func(rail miso.Rail) error {
		rail.Infof("Upload file using cURL: 'curl 'http://%s:%s/file?name=YOUR_FILE_NAME' --data-binary @YOUR_FILE_NAME'",
			miso.GetLocalIPV4(), miso.GetPropStr(miso.PropServerPort))
		rail.Infof("Access index.html on 'http://%s:%s'",
			miso.GetLocalIPV4(), miso.GetPropStr(miso.PropServerPort))
		return nil
	})

	miso.BootstrapServer(os.Args)
}
