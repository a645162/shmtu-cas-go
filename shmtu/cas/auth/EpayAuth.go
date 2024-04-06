package auth

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"shmtu-cas-go/shmtu/cas/auth/common"
	"shmtu-cas-go/shmtu/cas/captcha"
	"shmtu-cas-go/shmtu/utils"
)

// EpayAuth provides methods for Epay authentication and bill retrieval
type EpayAuth struct {
	savedCookie string
	htmlCode    string

	loginUrl    string
	loginCookie string
}

// GetBill retrieves bill information with pagination
func (ea *EpayAuth) GetBill(
	pageNo, tabNo, cookie string,
) (*resty.Response, int, string, string, error) {
	url := fmt.Sprintf(
		"https://ecard.shmtu.edu.cn/epay/consume/query?pageNo=%s&tabNo=%s",
		pageNo,
		tabNo,
	)

	client := resty.New().
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetCloseConnection(true)

	finalCookie := cookie
	if finalCookie == "" {
		finalCookie = ea.savedCookie
	}

	resp, _ := client.R().
		SetHeader("Cookie", cookie).
		Get(url)

	responseCode := resp.StatusCode()

	if responseCode == http.StatusOK {
		ea.htmlCode = string(resp.Body())
		return resp, responseCode, ea.htmlCode, finalCookie, nil
	} else if responseCode == http.StatusFound {
		location := resp.Header().Get("Location")
		setCookiesHeader := resp.Header().Get("Set-Cookie")
		newCookie := setCookiesHeader
		//println(setCookiesHeader)
		//setCookies := resp.Cookies()
		//var newCookie string
		//for _, cookie := range setCookies {
		//	if strings.Contains(cookie.Name, "JSESSIONID") {
		//		newCookie = cookie.String()
		//		break
		//	}
		//}
		ea.savedCookie = newCookie
		return resp,
			responseCode, location, newCookie,
			fmt.Errorf("redirect required, location: %s", location)
	}

	return resp,
		responseCode, "", "",
		fmt.Errorf("failed to get bill, status code: %d", responseCode)
}

// TestLoginStatus checks if the user is logged in
func (ea *EpayAuth) TestLoginStatus() bool {
	_, responseCode, loginUrl, savedCookie, _ := ea.GetBill("1", "1", ea.savedCookie)

	if responseCode == http.StatusOK {
		// OK
		return true
	} else if responseCode == http.StatusFound {
		// Update cookie after redirection in GetBill
		ea.loginUrl = loginUrl
		ea.savedCookie = savedCookie
		return false
	}

	return false
}

// Login performs login with username and password
func (ea *EpayAuth) Login(username, password string) (bool, error) {
	if ea.loginUrl == "" || ea.savedCookie == "" {
		loggedIn := ea.TestLoginStatus()
		if loggedIn {
			return true, nil
		}
	}

	// Call CasAuth functions (replace with actual implementation)
	executionStr, err := common.GetExecutionString(ea.loginUrl, ea.savedCookie)

	imageData, loginCookie, err := captcha.GetImageDataFromUrlUsingGet(ea.savedCookie)
	if err != nil {
		return false, fmt.Errorf("failed to get captcha image: %w", err)
	}
	if imageData == nil {
		return false, fmt.Errorf("failed to get captcha image data")
	}

	_ = utils.SaveImageDataToFile(imageData, "test.png")

	validateCode, err := captcha.OcrByRemoteTcpServer(
		"127.0.0.1", 21601, imageData,
	)
	if err != nil {
		return false, fmt.Errorf("failed to recognize captcha: %w", err)
	}
	println("Captcha:", validateCode)
	exprResult := captcha.GetExprResultByExprString(validateCode)

	// Call CasAuth functions (replace with actual implementation)
	resultCas, err :=
		common.CasLogin(
			ea.loginUrl,
			username, password,
			exprResult,
			executionStr,
			loginCookie,
		)

	if err != nil {
		return false, fmt.Errorf("cas login failed: %w", err)
	}

	if resultCas.ResponseCode != http.StatusFound {
		return false,
			fmt.Errorf(
				"cas login failed, status code: %d",
				resultCas.ResponseCode,
			)
	}

	ea.loginCookie = resultCas.Cookie

	resultCas, err =
		common.CasRedirect(resultCas.Location, ea.savedCookie)

	if err != nil {
		return false, fmt.Errorf("cas redirect failed: %w", err)
	}
	if resultCas.ResponseCode != http.StatusFound {
		return false,
			fmt.Errorf(
				"cas redirect failed, status code: %d",
				resultCas.ResponseCode,
			)
	}

	finalTestStatus := ea.TestLoginStatus()
	if !finalTestStatus {
		return false, fmt.Errorf("login failed")
	}

	return true, nil
}
