package demo

import (
	"fmt"
	"os"

	"shmtu-cas-go/shmtu/cas/auth"
)

func DemoBill() {

	// 创建 EpayAuth 实例
	epayAuth := &auth.EpayAuth{}

	println("Test Result:", epayAuth.TestLoginStatus())

	username := os.Getenv("SHMTU_USERNAME")
	password := os.Getenv("SHMTU_PASSWORD")

	fmt.Println("Username:", username)
	fmt.Println("Password:", password)

	loginStatus, err := epayAuth.Login(username, password)
	if err != nil {
		fmt.Println("Error logging in:", err)
		return
	}

	println("Login status:", loginStatus)

	_, statusCode, htmlCode, _, _ :=
		epayAuth.GetBill("1", "1", "")

	fmt.Printf("Status Code: %d\n", statusCode)
	fmt.Printf("HTML Code: %s\n", htmlCode)

	println("Test login status:", epayAuth.TestLoginStatus())

	println("finish")

	// HtmlCode 处理相关逻辑

}
