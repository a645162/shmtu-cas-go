package str

import "regexp"

func OnlyDigit(originString string) string {
	re := regexp.MustCompile(`\d+`)
	numbers := re.FindAllString(originString, -1)
	digitString := ""
	for _, number := range numbers {
		digitString += number
	}
	return digitString
}
