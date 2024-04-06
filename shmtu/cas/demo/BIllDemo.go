package demo

import (
	"fmt"
	"os"
	"shmtu-cas-go/shmtu/cas/auth"
)

func DemoBill() {

	// 创建 EpayAuth 实例
	epayAuth := &auth.EpayAuth{}

	username := os.Getenv("SHMTU_USERNAME")
	password := os.Getenv("SHMTU_PASSWORD")

	//fmt.Println("Username:", username)
	//fmt.Println("Password:", password)

	loginStatus, err := epayAuth.Login(username, password)
	if err != nil {
		fmt.Println("Error logging in:", err)
		return
	}

	println("Login status:", loginStatus)

	_, statusCode, htmlCode, _, err :=
		epayAuth.GetBill("1", "1", "")

	// 处理响应
	if err != nil {
		fmt.Println("Error retrieving bill:", err)
		return
	}

	fmt.Printf("Status Code: %d\n", statusCode)
	fmt.Printf("HTML Code: %s\n", htmlCode)

	// 如果需要，可以处理 response.Body() 中的数据

}
