package zapHelper

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
	
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const (
	// AlertWebHook https://developer.work.weixin.qq.com/document/path/91770
	AlertWebHook = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
)

func NewAlertEncoder(encoderConf zapcore.EncoderConfig, webhook string) zapcore.Encoder {
	return &AlertEncoder{
		AlertWebhook: webhook,
		Encoder:      zapcore.NewJSONEncoder(encoderConf),
	}
}

type AlertEncoder struct {
	AlertWebhook string
	zapcore.Encoder
}

// EncodeEntry alertEncoder encodeEntry
func (a *AlertEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (buf *buffer.Buffer, err error) {
	buf, err = a.Encoder.EncodeEntry(entry, fields)
	a.alarmToWeCom(buf.Bytes())
	return buf, err
}

// alarmToWeCom push error logs to WeCom（不能超过20条/分钟）
func (a *AlertEncoder) alarmToWeCom(buf []byte) {
	// get interface ip
	ipAddr := ""
	addrList, _ := net.InterfaceAddrs()
	for _, address := range addrList {
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				ipAddr = ip.IP.String()
				break
			}
		}
	}
	webHook := a.AlertWebhook
	var logContent LogContent
	err := json.Unmarshal(buf, &logContent)
	if err != nil {
		fmt.Println("Alert json Unmarshal content err：", err)
	}
	content := getContentStr(logContent, ipAddr)
	data := fmt.Sprintf(`{"msgtype":"markdown","markdown":{"content":%q}}`, content)
	resp, err := http.Post(webHook, "application/json", strings.NewReader(data))
	if err != nil {
		fmt.Println("Alert http post err:", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Alert ioutil readall body err：", err)
	}
	var result WeComWebHookReply
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Alert json Unmarshal body err：", err)
	}
	if result.Errcode == 0 {
		return
	}
}

// getContentStr splicing content
func getContentStr(logContent LogContent, ip string) (content string) {
	logContent.Color = "#99cc33"
	switch logContent.Level {
	case log.LevelDebug.String():
		logContent.Color = "#339900"
	case log.LevelInfo.String():
		logContent.Color = "#99cc33"
	case log.LevelWarn.String():
		logContent.Color = "#ffcc00"
	case log.LevelError.String():
		logContent.Color = "#ff9966"
	case log.LevelFatal.String():
		logContent.Color = "#cc3300"
	}
	// 定制日志通知格式
	content = `
        >level: <font color=` + logContent.Color + `>` + logContent.Level + `</font>
        >time: <font color="comment">` + logContent.Time.Format("2006-01-02 15:04:05") + `</font>
        >caller: <font color="comment">` + logContent.Caller + `</font>
        >massage: <font color="comment">` + logContent.Msg + `</font>
        >stack: <font color="comment">` + logContent.Stack + `</font>
        >operation: <font color="comment">` + logContent.Operation + `</font>
        >args: <font color="comment">` + logContent.Args + `</font>
        >ip: <font color="comment">` + ip + `</font>`
	return
}

type WeComWebHookReply struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

type LogContent struct {
	Level     string    `json:"level"`
	Time      time.Time `json:"time"`
	Caller    string    `json:"caller"`
	Msg       string    `json:"msg"`
	Stack     string    `json:"stack"`
	Color     string    `json:"color"`
	Operation string    `json:"operation"`
	Args      string    `json:"args"`
}
