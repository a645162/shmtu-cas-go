package billparser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"shmtu-cas-go/shmtu/utils/program_time"
	"shmtu-cas-go/shmtu/utils/str"
	"strconv"
	"strings"
	"time"
)

const (
	BillALL     = "#all"
	BillWaitFor = "#waitfor"
	BillSuccess = "#succ"
	BillFailure = "#fail"
)

type BillItemInfo struct {
	dateStr             string
	timeStr             string
	timeStrFormat       string
	dateTimeStrFormated string

	datetime  time.Time
	timeStamp int64

	itemType   string
	number     string
	targetUser string
	money      string
	method     string
	status     string
}

func ConvertBillInfoToHashmap(billInfo *BillItemInfo) map[string]string {
	hashmap := make(map[string]string)

	hashmap["dateStr"] = billInfo.dateStr
	hashmap["timeStr"] = billInfo.timeStr
	hashmap["timeStrFormat"] = billInfo.timeStrFormat
	hashmap["dateTimeStrFormated"] = billInfo.dateTimeStrFormated
	hashmap["datetime"] = billInfo.datetime.String()
	hashmap["timeStamp"] = strconv.FormatInt(billInfo.timeStamp, 10)

	hashmap["itemType"] = billInfo.itemType
	hashmap["number"] = billInfo.number
	hashmap["targetUser"] = billInfo.targetUser
	hashmap["money"] = billInfo.money
	hashmap["method"] = billInfo.method
	hashmap["status"] = billInfo.status

	return hashmap
}

func GetTargetTypeHtmlElement(
	htmlCode string, Type string,
) (*goquery.Selection, error) {

	document, err :=
		goquery.NewDocumentFromReader(strings.NewReader(htmlCode))

	if err != nil {
		return nil, fmt.Errorf("解析HTML文档时出错: %v", err)
	}

	tabContent :=
		document.Find("#content > div:nth-child(2) > div > div.panel-body > div > div")
	targetDiv := tabContent.Find(Type)

	return targetDiv, nil
}

func GetTotalPagesCount(htmlElement *goquery.Selection) (int, error) {

	if htmlElement == nil {
		return 0, fmt.Errorf("htmlElement is null")
	}

	pageCountTextElement :=
		htmlElement.Find("div > table > tbody > tr > td:nth-child(1)")
	pageText := strings.TrimSpace(pageCountTextElement.Text())

	if pageText == "" {
		return 0, fmt.Errorf("pageText is empty")
	}

	//println("pageCountTextElement:", pageText)

	fullTemplate := "当前[0-9]+/[0-9]+页"
	fullRegex := regexp.MustCompile(fullTemplate)
	matches := fullRegex.FindAllString(pageText, -1)
	if len(matches) != 1 {
		return 0, fmt.Errorf("匹配页码失败")
	}

	pageCountTemplate := "/[0-9]+"
	pageCountRegex := regexp.MustCompile(pageCountTemplate)
	pageCountText := pageCountRegex.FindString(matches[0])
	pageCountText = strings.Replace(pageCountText, "/", "", 1)

	pageCount, err := strconv.Atoi(pageCountText)

	if err != nil {
		return 0, fmt.Errorf("解析页码失败: %v", err)
	}

	return pageCount, nil
}

func GetBillList(htmlElement *goquery.Selection) ([]BillItemInfo, error) {
	tbodyElement := htmlElement.Find("span > table > tbody")
	trElement := tbodyElement.Find("tr")
	println("trElement:", trElement.Length())

	billList := make([]BillItemInfo, trElement.Length())

	for i := 0; i < trElement.Length(); i++ {
		tr := trElement.Eq(i)
		children := tr.Children()
		if children.Length() != 7 {
			return nil, fmt.Errorf("tr > children.Length() != 7")
		}

		billItemInfo := BillItemInfo{}

		timeChildren := children.Eq(0).Children()
		if timeChildren.Length() != 2 {
			return nil, fmt.Errorf("timeChildren.Length() != 2")
		}
		billItemInfo.dateStr = strings.TrimSpace(timeChildren.Eq(0).Text())
		billItemInfo.timeStr = strings.TrimSpace(timeChildren.Eq(1).Text())
		billItemInfo.timeStrFormat =
			program_time.AddChatTo6DigitTime(billItemInfo.timeStr)
		billItemInfo.dateTimeStrFormated =
			billItemInfo.dateStr + " " + billItemInfo.timeStrFormat

		datetime, err :=
			program_time.ParseTimeFromString(billItemInfo.dateTimeStrFormated)
		if err != nil {
			return nil, fmt.Errorf("解析时间失败: %v", err)
		}
		billItemInfo.datetime = datetime
		billItemInfo.timeStamp = datetime.Unix()

		dealChildren := children.Eq(1).Children()
		if dealChildren.Length() != 2 {
			return nil, fmt.Errorf("dealChildren.Length() != 2")
		}
		billItemInfo.itemType = strings.TrimSpace(dealChildren.Eq(0).Text())
		billItemInfo.number =
			str.OnlyDigit(strings.TrimSpace(dealChildren.Eq(1).Text()))

		billItemInfo.targetUser = strings.TrimSpace(children.Eq(2).Text())
		billItemInfo.money = strings.TrimSpace(children.Eq(3).Text())
		billItemInfo.method = strings.TrimSpace(children.Eq(4).Text())
		billItemInfo.status = strings.TrimSpace(children.Eq(5).Text())

		billList[i] = billItemInfo
	}

	return billList, nil
}
