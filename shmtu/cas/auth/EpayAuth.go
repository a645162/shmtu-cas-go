package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"

	"shmtu-cas-go/shmtu/cas/auth/common"
	"shmtu-cas-go/shmtu/cas/captcha"
	"shmtu-cas-go/shmtu/utils"
)

// EpayAuth provides methods for Epay authentication and bill retrieval
type EpayAuth struct {
	SavedCookie string
	htmlCode    string

	loginUrl    string
	loginCookie string
}

// GetBill retrieves bill information with pagination
func (ea *EpayAuth) GetBill(
	pageNo string, tabNo string,
	cookie string,
) (*resty.Response, int, string, string, error) {
	if pageNo == "" {
		pageNo = "1"
	}
	if tabNo == "" {
		tabNo = "1"
	}

	// https://ecard.shmtu.edu.cn/epay/consume/query?pageNo=1&tabNo=1
	url := fmt.Sprintf(
		"https://ecard.shmtu.edu.cn/epay/consume/query?pageNo=%s&tabNo=%s",
		pageNo,
		tabNo,
	)

	client := resty.New().
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetCloseConnection(true)

	finalCookie := strings.TrimSpace(cookie)
	if finalCookie == "" {
		finalCookie = ea.SavedCookie
	}
	println("Cookie:", finalCookie)

	resp, _ := client.R().
		SetHeader("Cookie", finalCookie).
		Get(url)

	responseCode := resp.StatusCode()

	if responseCode == http.StatusOK {
		ea.htmlCode = strings.TrimSpace(string(resp.Body()))
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
		ea.SavedCookie = newCookie
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
	_, responseCode, loginUrl, savedCookie, _ :=
		ea.GetBill("1", "1", ea.SavedCookie)

	if responseCode == http.StatusOK {
		// OK
		return true
	} else if responseCode == http.StatusFound {
		// Update cookie after redirection in GetBill
		ea.loginUrl = loginUrl

		cookieParts := strings.Split(savedCookie, ";")
		for _, part := range cookieParts {
			if strings.Contains(strings.TrimSpace(part), "JSESSIONID") {
				savedCookie = strings.TrimSpace(part)
			}
		}

		ea.SavedCookie = savedCookie
		return false
	}

	return false
}

// Login performs login with username and password
func (ea *EpayAuth) Login(
	username, password string,
	ocrServerHost string,
) (bool, error) {

	if ocrServerHost == "" {
		ocrServerHost = "localhost"
	}

	if ea.loginUrl == "" || ea.SavedCookie == "" {
		loggedIn := ea.TestLoginStatus()
		if loggedIn {
			return true, nil
		}
	}

	// Call CasAuth functions (replace with actual implementation)
	executionStr, err :=
		common.GetExecutionString(ea.loginUrl, ea.SavedCookie)

	imageData, loginCookie, err :=
		captcha.GetImageDataFromUrlUsingGet(ea.SavedCookie)
	if err != nil {
		return false, fmt.Errorf("failed to get captcha image: %w", err)
	}
	if imageData == nil {
		return false, fmt.Errorf("failed to get captcha image data")
	}
	ea.loginCookie = loginCookie

	_ = utils.SaveImageDataToFile(imageData, "test.png")

	validateCode, err := captcha.OcrByRemoteTcpServer(
		ocrServerHost, 21601, imageData,
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
		common.CasRedirect(resultCas.Location, ea.SavedCookie)

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
		return false, fmt.Errorf("login final test failed")
	}

	return true, nil
}
