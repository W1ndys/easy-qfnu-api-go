package notify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// FeishuNotifier 飞书通知器
type FeishuNotifier struct {
	webhookURL string
	secret     string
	client     *resty.Client
}

var (
	defaultNotifier *FeishuNotifier
	once            sync.Once
)

// InitFeishu 初始化飞书通知器
func InitFeishu() {
	once.Do(func() {
		webhookURL := os.Getenv("FEISHU_WEBHOOK_URL")
		secret := os.Getenv("FEISHU_WEBHOOK_SECRET")

		if webhookURL == "" {
			slog.Warn("飞书通知未配置: FEISHU_WEBHOOK_URL 为空")
			return
		}

		defaultNotifier = &FeishuNotifier{
			webhookURL: webhookURL,
			secret:     secret,
			client:     resty.New().SetTimeout(10 * time.Second),
		}
		slog.Info("飞书通知初始化完成")
	})
}

// genSign 生成签名
func (f *FeishuNotifier) genSign(timestamp int64) (string, error) {
	if f.secret == "" {
		return "", nil
	}

	stringToSign := fmt.Sprintf("%v\n%s", timestamp, f.secret)
	h := hmac.New(sha256.New, []byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

// CardMessage 消息卡片结构
type CardMessage struct {
	MsgType string      `json:"msg_type"`
	Card    interface{} `json:"card"`
}

// buildCard 构建消息卡片
func buildCard(title, content string, color string) map[string]interface{} {
	// 颜色映射
	colorMap := map[string]string{
		"green":  "green",
		"red":    "red",
		"orange": "orange",
		"blue":   "blue",
	}
	headerColor := colorMap[color]
	if headerColor == "" {
		headerColor = "blue"
	}

	return map[string]interface{}{
		"header": map[string]interface{}{
			"title": map[string]interface{}{
				"tag":     "plain_text",
				"content": title,
			},
			"template": headerColor,
		},
		"elements": []interface{}{
			map[string]interface{}{
				"tag":     "markdown",
				"content": content,
			},
			map[string]interface{}{
				"tag": "hr",
			},
			map[string]interface{}{
				"tag":     "note",
				"elements": []interface{}{
					map[string]interface{}{
						"tag":     "plain_text",
						"content": fmt.Sprintf("⏰ %s", time.Now().Format("2006-01-02 15:04:05")),
					},
				},
			},
		},
	}
}

// Send 发送消息卡片
func (f *FeishuNotifier) Send(title, content, color string) error {
	if f == nil {
		return nil
	}

	timestamp := time.Now().Unix()
	sign, err := f.genSign(timestamp)
	if err != nil {
		return fmt.Errorf("生成签名失败: %w", err)
	}

	card := buildCard(title, content, color)
	payload := map[string]interface{}{
		"msg_type": "interactive",
		"card":     card,
	}

	if sign != "" {
		payload["timestamp"] = fmt.Sprintf("%d", timestamp)
		payload["sign"] = sign
	}

	resp, err := f.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(f.webhookURL)

	if err != nil {
		return fmt.Errorf("发送飞书消息失败: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if code, ok := result["code"].(float64); ok && code != 0 {
		return fmt.Errorf("飞书返回错误: %v", result["msg"])
	}

	return nil
}

// SendAsync 异步发送消息（不阻塞主流程）
func (f *FeishuNotifier) SendAsync(title, content, color string) {
	go func() {
		if err := f.Send(title, content, color); err != nil {
			slog.Error("飞书通知发送失败", "error", err)
		}
	}()
}

// NotifyStartup 系统启动通知
func NotifyStartup(port string) {
	if defaultNotifier == nil {
		return
	}

	hostname, _ := os.Hostname()
	content := fmt.Sprintf(`**🚀 Easy QFNU API 服务已启动**

- **主机名**: %s
- **监听端口**: %s
- **启动时间**: %s`,
		hostname,
		port,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	defaultNotifier.SendAsync("系统启动通知", content, "green")
}

// NotifyNewRecommendation 新选课推荐提交通知
func NotifyNewRecommendation(courseName, teacher, recommender, reason string) {
	if defaultNotifier == nil {
		return
	}

	content := fmt.Sprintf(`**📚 收到新的选课推荐**

- **课程名称**: %s
- **授课教师**: %s
- **推荐人**: %s
- **推荐理由**: %s`,
		courseName,
		teacher,
		recommender,
		reason,
	)

	defaultNotifier.SendAsync("新选课推荐", content, "blue")
}

// NotifyError 系统错误通知
func NotifyError(errType, errMsg, stack string) {
	if defaultNotifier == nil {
		return
	}

	content := fmt.Sprintf(`**❌ 系统发生错误**

- **错误类型**: %s
- **错误信息**: %s
- **堆栈信息**:
\`\`\`
%s
\`\`\``,
		errType,
		errMsg,
		stack,
	)

	defaultNotifier.SendAsync("系统错误告警", content, "red")
}

// NotifyCustom 自定义通知
func NotifyCustom(title, content, color string) {
	if defaultNotifier == nil {
		return
	}
	defaultNotifier.SendAsync(title, content, color)
}

// GetNotifier 获取通知器实例
func GetNotifier() *FeishuNotifier {
	return defaultNotifier
}
