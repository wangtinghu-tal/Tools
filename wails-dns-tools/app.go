package main

import (
	"context"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/valyala/fasthttp"
	"net"
	"strconv"
	"time"
)

var OPT = Options{}

const (
	count                = 10               // ping次数
	downloadTestDuration = 10 * time.Second // 下载测试时长
	testSpeedURL         = "https://www.sohu.com"
	testNetworkHost      = "www.sohu.com"
	testSpeedTimes       = 50 // 下载测试次数
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// 检测网络上下行速率
func (a *App) CheckSpeed() (downloadSpeed float64) {
	url := testSpeedURL
	duration := downloadTestDuration
	totalSize := 0
	totalDuration := 0.0

	for i := 0; i < testSpeedTimes; i++ {
		req := fasthttp.AcquireRequest()
		req.SetRequestURI(url)

		resp := fasthttp.AcquireResponse()

		startTime := time.Now()
		err := fasthttp.DoTimeout(req, resp, duration)

		if err != nil {
			fmt.Println("无法进行网速测试:", err)
			return -1.0
		}

		endTime := time.Now()

		totalSize += resp.Header.ContentLength()
		totalDuration += endTime.Sub(startTime).Seconds()
	}

	downloadSpeed = Decimal(float64(totalSize) / totalDuration / 1024 / 1024)

	return
}

// 检测网络延迟和丢包率
func (a *App) CheckLatency() (late []float64) {
	if OPT.TestHost == "" {
		OPT.TestHost = testNetworkHost
	}
	pingClient, err := ping.NewPinger(OPT.TestHost)
	if err != nil {
		fmt.Println("无法创建ping实例:", err)
		return
	}

	pingClient.Count = count
	pingClient.Interval = time.Millisecond * 100
	pingClient.Timeout = time.Second * 5
	pingClient.SetPrivileged(true)
	err = pingClient.Run()
	if err != nil {
		fmt.Printf("ping失败: %s\n", err)
		return []float64{-1.0, -1.0}
	}

	stats := pingClient.Statistics()
	fmt.Printf("包数：%d 发送：%d 接收：%d 丢失率：%.2f%% 平均延迟：%s \n", stats.PacketsSent, stats.PacketsRecv, stats.PacketsSent-stats.PacketsRecv, stats.PacketLoss, stats.AvgRtt.String())

	return []float64{Decimal(stats.PacketLoss), Decimal(float64(stats.AvgRtt.Milliseconds()))}
}

// 解析DNS
func (a *App) CheckDNS() (ips []string) {
	if OPT.TestDNSHost == "" {
		OPT.TestDNSHost = "www.baidu.com"
	}
	ip, err := net.ResolveIPAddr("ip", OPT.TestDNSHost)
	if err != nil {
		fmt.Println("无法解析主机:", err)
		return nil
	}
	ips = append(ips, ip.String())
	fmt.Printf("dns res is :%v\n", ips)
	return ips
}

// GetOptions 获取配置文件JSON
func (a *App) GetOptions() Options {
	return OPT.GetData()
}

// SetOptions 解析传来的配置
func (a *App) SetOptions(val Options) {
	OPT.SetData(&val)
}

// Decimal float64保留2位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
