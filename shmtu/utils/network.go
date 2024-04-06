package utils

import (
	"regexp"
	"strconv"
)

// ValidateIPAddress 验证IP地址是否有效
func ValidateIPAddress(ip string) bool {
	ipAddressPattern := regexp.MustCompile(`^([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.` +
		`([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.` +
		`([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.` +
		`([01]?\\d\\d?|2[0-4]\\d|25[0-5])$`)
	return ipAddressPattern.MatchString(ip)
}

// ValidatePortString 验证端口号是否有效（字符串输入）
func ValidatePortString(port string) bool {
	integerPort, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return ValidatePort(integerPort)
}

// ValidatePort 验证端口号是否有效（整数输入）
func ValidatePort(port int) bool {
	return port >= 0 && port <= 65535
}
