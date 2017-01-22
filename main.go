package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/elango/go-amazon-product-api"
	"gopkg.in/kataras/iris.v4"
)

func main() {
	api := iris.New()
	api.Post("/search", productSearch)
	api.Listen(":8080")
}

func getAPIHandler() amazonproduct.AmazonProductAPI {
	var api amazonproduct.AmazonProductAPI

	api.AccessKey = "XXXXXXXXXXXXXXXXX"
	api.SecretKey = "XXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	api.Host = "webservices.amazon.com"
	api.AssociateTag = "XXXXXXXXXXXX"
	api.Client = &http.Client{} // optional

	return api
}

/*
	Post Handler to query Amazon product search API
*/
func productSearch(ctx *iris.Context) {
	keyword := ctx.FormValue("keyword")
	responseGroup := string(ctx.FormValue("responseGroup"))
	searchIndex := string(ctx.FormValue("searchIndex"))
	pageIndexStr := string(ctx.FormValue("pageIndex"))

	pageIndex := 1

	if pageIndexConv, err := strconv.Atoi(pageIndexStr); err == nil {
		pageIndex = pageIndexConv
	}

	if string(keyword) == "" {
		ctx.Error("Invalid Keyword", 400)
		return
	}

	//Images,ItemAttributes,Small,EditorialReview

	if string(responseGroup) == "" {
		responseGroup = "Images,ItemAttributes,Small,EditorialReview"
	}

	// Search Index - Defaults to 'All' ( Case sensitive )
	if searchIndex == "" {
		searchIndex = "All"
	}

	responseGroup = strings.Replace(responseGroup, " ", "", -1)

	api := getAPIHandler()
	result, err := api.ItemSearchByKeywordWithResponseGroupWithSearchIndex(string(keyword), string(responseGroup), searchIndex, pageIndex)
	if err != nil {
		fmt.Println(err)
	}

	xml := strings.NewReader(result)

	json, err := xj.Convert(xml)
	if err != nil {
		panic("That's embarrassing...")
	}

	//fmt.Println(result)
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Write(json.String())
}
