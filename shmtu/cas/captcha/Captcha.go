package captcha

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func GetImageDataFromUrl(imageUrl string) ([]byte, error) {
	response, err := http.Get(imageUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"failed to fetch image, status code: %d",
			response.StatusCode,
		)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

func OcrByRemoteTcpServer(host string, port int, imageData []byte) (string, error) {
	// 连接到服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// 设置超时时间为 5 秒（如果需要）
	// conn.SetDeadline(time.Now().Add(5 * time.Second))

	// 发送图像数据
	_, err = conn.Write(imageData)
	if err != nil {
		return "", err
	}

	// 发送特殊标记，表示图像数据发送完毕
	endMarker := []byte("<END>")
	_, err = conn.Write(endMarker)
	if err != nil {
		return "", err
	}

	// 读取响应
	response := bytes.Buffer{}
	_, err = io.Copy(&response, conn)
	if err != nil {
		return "", err
	}

	return response.String(), nil
}

func OcrByRemoteTcpServerAutoRetry(host string, port int, imageData []byte, retryTimes int) (string, error) {
	var result string
	var err error

	for i := 0; i < retryTimes; i++ {
		result, err = OcrByRemoteTcpServer(host, port, imageData)
		if err != nil {
			fmt.Printf("第%d次尝试远程识别验证码失败\n", i+1)
			fmt.Printf("错误信息：%v\n", err)
			time.Sleep(1 * time.Second) // 等待1秒后再重试
			continue
		}

		if result != "" {
			break
		}
	}

	return result, err
}

func GetExprResultByExprString(expr string) string {
	index := strings.Index(expr, "=")
	if index != -1 {
		result := strings.TrimSpace(expr[index+1:])
		return result
	}
	return ""
}

func GetImageDataFromUrlUsingGet(cookie string) ([]byte, string, error) {
	imageUrl := "https://cas.shmtu.edu.cn/cas/captcha"

	client := resty.New()

	req := client.R()
	if cookie != "" {
		req.SetCookie(&http.Cookie{Name: "Cookie", Value: cookie})
	}

	resp, err := req.Get(imageUrl)
	if err != nil {
		return nil, "", err
	}

	responseCode := resp.StatusCode()

	if responseCode != http.StatusOK {
		return nil, "", fmt.Errorf("请求失败，状态码：%d", responseCode)
	}

	returnCookie := resp.Header().Get("Set-Cookie")
	if returnCookie == "" {
		returnCookie = cookie
	}

	bodyBytes := resp.Body()

	return bodyBytes, returnCookie, nil
}
