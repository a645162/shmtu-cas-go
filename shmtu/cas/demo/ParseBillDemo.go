package demo

import (
	"shmtu-cas-go/shmtu/cas/parser/billparser"
	"shmtu-cas-go/shmtu/cas/parser/billparser/export"
	"shmtu-cas-go/shmtu/utils"
)

func ParseBillDemo() {

	htmlCode, err := utils.ReadTextFromFile("./result.html")
	if err != nil {
		panic(err)
	}

	tabElement, _ := billparser.GetTargetTypeHtmlElement(
		htmlCode, billparser.BillALL,
	)

	pageCount, err := billparser.GetTotalPagesCount(tabElement)
	if err != nil {
		panic(err)
	}

	println("pageCount:", pageCount)

	billList, err := billparser.GetBillList(tabElement)
	if err != nil {
		panic(err)
	}
	err = export.ToCsvFile(
		"./result.csv",
		billList,
	)
	if err != nil {
		panic(err)
	}

}
