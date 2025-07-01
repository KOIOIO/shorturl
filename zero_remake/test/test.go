package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"sync"
	"time"
)

func goroutine1(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	url := "http://localhost:8080/generate"
	method := "POST"
	client := &http.Client{}

	for i := 0; i < 1000; i++ { // 每个协程发送100次请求，可自行调整
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		longUrl := fmt.Sprintf("https://example.com/page/%d?rand=%d", i, rand.Intn(1000000))
		_ = writer.WriteField("url", longUrl)
		_ = writer.WriteField("expiration", "1h")
		err := writer.Close()
		if err != nil {
			fmt.Printf("[Goroutine %d] %v\n", id, err)
			continue
		}

		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			fmt.Printf("[Goroutine %d] %v\n", id, err)
			continue
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("[Goroutine %d] %v\n", id, err)
			continue
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Printf("[Goroutine %d] %v\n", id, err)
			continue
		}
		fmt.Printf("[Goroutine %d] 第%d次请求，长链接：%s\n返回：%s\n", id, i+1, longUrl, string(body))
		time.Sleep(100 * time.Millisecond) // 间隔可调整
	}
}

func main() {
	var wg sync.WaitGroup
	goroutineNum := 5 // 并发协程数，可自行调整

	for i := 0; i < goroutineNum; i++ {
		wg.Add(1)
		go goroutine1(i+1, &wg)
	}

	wg.Wait()
}
