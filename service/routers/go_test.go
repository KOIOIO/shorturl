package routers

import (
	"fmt"
	"net/http"
	"os"
	"shorturl/service/routers"
	"time"
)

func main() {
	// 启动服务
	go func() {
		router := routers.InitRouter()
		_ = router.Run(":8080")
	}()

	// 等待服务启动
	time.Sleep(2 * time.Second)

	// 定义压测目标
	targeter := func() (*lib.Target, error) {
		// 这里可以根据需要修改请求的方法、URL、头部和正文
		// 压测 /generate 接口
		req, err := http.NewRequest("POST", "http://localhost:8080/generate", nil)
		if err != nil {
			return nil, err
		}
		return &lib.Target{
			Method: req.Method,
			URL:    req.URL.String(),
			Header: req.Header,
		}, nil
	}

	// 定义压测速率和持续时间
	rate := lib.Rate{Freq: 100, Per: time.Second} // 每秒 100 个请求
	duration := 30 * time.Second                  // 持续 30 秒

	// 执行压测
	attacker := lib.NewAttacker()
	var metrics lib.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Load testing") {
		metrics.Add(res)
	}
	metrics.Close()

	// 输出压测结果
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("Requests/sec: %.2f\n", metrics.Rate)
	fmt.Printf("Success ratio: %.2f%%\n", metrics.Success*100)

	// 压测 /:shortURL 接口
	targeterShortURL := func() (*lib.Target, error) {
		req, err := http.NewRequest("GET", "http://localhost:8080/testShortURL", nil)
		if err != nil {
			return nil, err
		}
		return &lib.Target{
			Method: req.Method,
			URL:    req.URL.String(),
			Header: req.Header,
		}, nil
	}

	var metricsShortURL lib.Metrics
	for res := range attacker.Attack(targeterShortURL, rate, duration, "Load testing for shortURL") {
		metricsShortURL.Add(res)
	}
	metricsShortURL.Close()

	// 输出 /:shortURL 接口压测结果
	fmt.Printf("99th percentile for shortURL: %s\n", metricsShortURL.Latencies.P99)
	fmt.Printf("Requests/sec for shortURL: %.2f\n", metricsShortURL.Rate)
	fmt.Printf("Success ratio for shortURL: %.2f%%\n", metricsShortURL.Success*100)

	// 将结果保存到文件
	f, err := os.Create("results.bin")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := lib.NewEncoder(f).Encode(metrics); err != nil {
		panic(err)
	}
}
