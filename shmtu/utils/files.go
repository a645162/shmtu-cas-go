package utils

import (
	"bytes"
	"image"
	"image/png"
	"os"
)

func SaveImageDataToFile(data []byte, filePath string) error {
	return os.WriteFile(filePath, data, 0644)
}

func ReadImageFromFile(fileName string) ([]byte, error) {
	// 打开文件
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 解码图像
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// 将图像转换为PNG格式
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	// 返回图像的字节数据
	return buf.Bytes(), nil
}
