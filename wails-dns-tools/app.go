package main

import (
	"context"
	"errors"
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
func (a *App) CheckSpeed() (downloadSpeed string, err error) {
	url := "https://www.sohu.com"
	duration := downloadTestDuration
	totalSize := 0
	totalDuration := 0.0

	for i := 0; i < 10; i++ {
		req := fasthttp.AcquireRequest()
		req.SetRequestURI(url)

		resp := fasthttp.AcquireResponse()

		startTime := time.Now()
		err = fasthttp.DoTimeout(req, resp, duration)

		if err != nil {
			fmt.Println("无法进行网速测试:", err)
			return "", errors.New("无法进行网速测试")
		}

		endTime := time.Now()

		totalSize += resp.Header.ContentLength()
		totalDuration += endTime.Sub(startTime).Seconds()
	}

	downloadSpeed = fmt.Sprintf("%.5f Mb/s", float64(totalSize)/totalDuration/1024/1024)

	return
}

// 检测网络延迟和丢包率
func (a *App) CheckLatency() (loss string, delay string, err error) {
	if OPT.TestHost == "" {
		return "", "", fmt.Errorf("未设置测试域名")
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
		errName := fmt.Sprintf("ping失败: %s", err)
		return "", "", errors.New(errName)
	}

	stats := pingClient.Statistics()
	//fmt.Printf("包数：%d 发送：%d 接收：%d 丢失率：%.2f%% 平均延迟：%s \n", stats.PacketsSent, stats.PacketsRecv, stats.PacketsSent-stats.PacketsRecv, stats.PacketLoss, stats.AvgRtt.String())

	return fmt.Sprintf("%.2f%%", stats.PacketLoss), fmt.Sprintf("%s", stats.AvgRtt.String()), nil
}

// 解析DNS
func (a *App) CheckDNS() (ips []string, err error) {
	host := "www.baidu.com"

	ips, err = net.LookupHost(host)
	if err != nil {
		fmt.Println("无法解析主机:", err)
		return nil, errors.New("无法解析主机")
	}
	return
}

// GetOptions 获取配置文件JSON
func (a *App) GetOptions() Options {
	return OPT.GetData()
}

// SetOptions 解析传来的配置
func (a *App) SetOptions(val Options) {
	OPT.SetData(&val)
}

// string to float64
func s2f(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
