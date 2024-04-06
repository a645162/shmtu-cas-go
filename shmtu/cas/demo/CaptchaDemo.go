package demo

import (
	"fmt"

	"shmtu-cas-go/shmtu/cas/captcha"
	"shmtu-cas-go/shmtu/utils"
)

func DemoCaptcha() {
	imageUrl := "https://cas.shmtu.edu.cn/cas/captcha"
	imageData, err := captcha.GetImageDataFromUrl(imageUrl)
	if err != nil {
		fmt.Println("Error fetching image data:", err)
		return
	}

	err = utils.SaveImageDataToFile(imageData, "test.png")
	if err != nil {
		fmt.Println("Error saving image data:", err)
		return
	}

	result, _ :=
		captcha.OcrByRemoteTcpServerAutoRetry(
			"127.0.0.1", 21601,
			imageData,
			1,
		)

	fmt.Println("Ocr Result:", result)

	answer := captcha.GetExprResultByExprString(result)
	fmt.Println(answer)

	fmt.Println("Image data saved.")
}
