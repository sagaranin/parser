package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const rootUrl = "https://online.globus.ru/catalog"

var store Store

func init() {
	store = Store{}
	store.Open()
}

func main() {
	categoryList := loadCategoryList()
	for _, v := range categoryList {
		loadCategory(v)
	}
}

func loadCategoryList() []string {
	log.Println("loading category list...")
	var response []string

	resp, err := http.Get(rootUrl)
	if err != nil {
		log.Println(err)
	}

	re, _ := regexp.Compile(`\<a href=\"/catalog(/[a-z\-]*/)\"\>`)
	body, err := io.ReadAll(resp.Body)
	linksMatch := re.FindAllStringSubmatch(string(body), -1)

	for _, v := range linksMatch {
		response = append(response, v[1])
	}

	log.Printf("found %d categories...\n", len(response))
	return response
}

func loadCategory(categoryUrl string) {
	var pageNum int = 1

	for {
		pagedUrl := fmt.Sprintf("%s%s?PAGEN_1=%d", rootUrl, categoryUrl, pageNum)
		log.Printf("get data from url: %s\n", pagedUrl)

		resp, err := http.Get(pagedUrl)
		if err != nil {
			log.Println(err)
		}

		re, _ := regexp.Compile(`GlobusProducts.addList\((.*)\);`)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		found := re.FindAllStringSubmatch(string(body), -1)
		var itemList string

		if len(found) > 0 {
			if len(found[0]) > 1 {
				itemList = strings.ReplaceAll(found[0][1], "'", "\"")
			}
		} else {
			log.Printf("Page %d - void response\n", pageNum)
			break
		}

		allItems := make(map[string]GlobusItem)
		json.Unmarshal([]byte(itemList), &allItems)

		saveItems(allItems)
		pageNum++
		time.Sleep(time.Second)
	}
}

func saveItems(items map[string]GlobusItem) {
	log.Printf("saving %d items...\n", len(items))
	if len(items) > 0 {
		store.SaveGlobus(items)
	}
}
