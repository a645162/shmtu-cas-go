package demo

import (
	"fmt"
	"os"

	"shmtu-cas-go/shmtu/cas/auth"
	"shmtu-cas-go/shmtu/utils"
)

func BillDemo() {

	// 创建 EpayAuth 实例
	epayAuth := &auth.EpayAuth{}

	println("Test Result:", epayAuth.TestLoginStatus())

	username := os.Getenv("SHMTU_USER_ID")
	password := os.Getenv("SHMTU_PASSWORD")
	ocrServerHost := os.Getenv("OCR_SERVER_HOST")

	//fmt.Println("Username:", username)
	//fmt.Println("Password:", password)

	loginStatus, err := epayAuth.Login(username, password, ocrServerHost)
	if err != nil {
		fmt.Println("Error logging in:", err)
		return
	}

	println("Login status:", loginStatus)

	_, statusCode, htmlCode, _, _ :=
		epayAuth.GetBill("1", "1", "")

	fmt.Printf("Status Code: %d\n", statusCode)
	//fmt.Printf("HTML Code: %s\n", htmlCode)

	err = utils.SaveTextToFile("result.html", htmlCode)
	if err != nil {
		fmt.Println("写入文件时发生错误:", err)
		return
	}

	println("Test login status:", epayAuth.TestLoginStatus())

	println("finish")

	// HtmlCode 处理相关逻辑

}
