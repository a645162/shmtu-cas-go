package billparser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func getCoreBillListHtmlElement(htmlCode string) (string, error) {

	document, err :=
		goquery.NewDocumentFromReader(strings.NewReader(htmlCode))

	if err != nil {
		return "", fmt.Errorf("解析HTML文档时出错: %v", err)
	}

	errorPanel := document.Find("//*[@id=\"all\"]")
	errorText := strings.TrimSpace(errorPanel.Text())
	fmt.Println(errorText)

	return "table", nil
}

func GetTotalPagesCount(htmlCode string) (int, error) {
	// Get the total pages count from the html code
	return 0, nil
}
