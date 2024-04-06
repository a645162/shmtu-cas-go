package common

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

//type CasAuthStatus int

const (
	SUCCESS int = iota
	VALIDATE_CODE_ERROR
	PASSWORD_ERROR
	FAILURE
)

func GetExecutionString(url string, cookie string) (string, error) {
	client := resty.New().
		SetRedirectPolicy(resty.NoRedirectPolicy())

	resp, err := client.R().
		SetHeader("Cookie", cookie).
		Get(url)

	if err != nil {
		return "", fmt.Errorf("获取execution时发生错误: %v", err)
	}

	if resp.StatusCode() == http.StatusOK {
		document, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
		if err != nil {
			return "", fmt.Errorf("解析HTML文档时出错: %v", err)
		}

		executionElement := document.Find("input[name=execution]")
		if executionElement.Length() > 0 {
			value, _ := executionElement.Attr("value")
			return strings.TrimSpace(value), nil
		}
		return "", fmt.Errorf("未找到execution元素")
	}

	return "", fmt.Errorf("获取execution失败，状态码: %d", resp.StatusCode())
}

type CasAuthResponse struct {
	ResponseCode int
	Location     string
	Cookie       string
	ErrorMessage string
}

func CasLogin(
	url string,
	username string,
	password string,
	validateCode string,

	execution string,
	cookie string,
) (*CasAuthResponse, error) {
	client := resty.New().
		SetRedirectPolicy(resty.NoRedirectPolicy())

	formData := make(map[string]string)
	formData["username"] = strings.TrimSpace(username)
	formData["password"] = strings.TrimSpace(password)
	formData["validateCode"] = strings.TrimSpace(validateCode)
	formData["execution"] = strings.TrimSpace(execution)
	formData["_eventId"] = "submit"
	formData["geolocation"] = ""

	resp, _ := client.R().
		SetHeader("Host", "cas.shmtu.edu.cn").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Connection", "keep-alive").
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetHeader("Accept", "*/*").
		SetHeader("Cookie", strings.TrimSpace(cookie)).
		SetFormData(formData).
		Post(url)

	responseCode := resp.StatusCode()

	var result *CasAuthResponse
	if responseCode == http.StatusFound { // Status code 302
		location := resp.Header().Get("Location")
		newCookie := resp.Header().Get("Set-Cookie")

		result = &CasAuthResponse{
			ResponseCode: responseCode,
			Location:     location,
			Cookie:       newCookie,
		}
	} else {
		document, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
		if err != nil {
			return nil, fmt.Errorf("解析登录错误页面时出错: %v", err)
		}

		errorPanel := document.Find("#loginErrorsPanel")
		errorText := errorPanel.Text()
		fmt.Println(errorText)

		if strings.Contains(errorText, "account is not recognized") {
			fmt.Println("用户名或密码错误")
			result = &CasAuthResponse{
				ResponseCode: PASSWORD_ERROR,
				ErrorMessage: errorText,
			}
		} else if strings.Contains(errorText, "reCAPTCHA") {
			fmt.Println("验证码错误")
			result = &CasAuthResponse{
				ResponseCode: VALIDATE_CODE_ERROR,
				ErrorMessage: errorText,
			}
		} else {
			result = &CasAuthResponse{
				ResponseCode: responseCode,
				ErrorMessage: errorText,
			}
		}
	}

	return result, nil
}

func CasRedirect(url string, cookie string) (*CasAuthResponse, error) {
	client := resty.New().
		SetRedirectPolicy(resty.NoRedirectPolicy())

	resp, _ := client.R().
		SetHeader("Cookie", cookie).
		Get(url)

	responseCode := resp.StatusCode()

	if responseCode == 302 {
		location := resp.Header().Get("Location")
		newCookie := resp.Header().Get("Set-Cookie")

		return &CasAuthResponse{
			ResponseCode: responseCode,
			Location:     location,
			Cookie:       newCookie,
		}, nil
	} else {
		return &CasAuthResponse{
			ResponseCode: responseCode,
		}, fmt.Errorf("重定向失败，状态码: %d", responseCode)
	}
}
