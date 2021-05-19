// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/davecgh/go-spew/spew"
)

const (
	// PageRange -
	PageRange = 5

	// 搜尋關鍵字
	keyWord = "二手精品"

	// 記錄資料檔名
	FileName = "crawlerRecord.json"
)

// CrawlerRecordData - json記錄資料
type CrawlerRecordData struct {
	AT    string `json:"at"`
	Rank  int    `json:"rank"`
	Page  int    `json:"page"`
	Index int    `json:"index"`
	Title string `json:"title"`
}

var CrawlerRecord = make([]*CrawlerRecordData, 0) // 記錄資料陣列

// Result -
type Result struct {
	AT     string `json:"at"`
	Rank   int    `json:"rank"`
	Page   int    `json:"page"`
	Index  int    `json:"index"`
	Target bool   `json:"target"`
	Title  string `json:"title"`
}

var wg sync.WaitGroup

func main() {
	options := []chromedp.ExecAllocatorOption{
		chromedp.ExecPath("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"), // 不加載圖片
		chromedp.Flag("headless", true),                        // 無頭模式
		chromedp.WindowSize(1920, 1080),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(
		context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			options...,
		)...,
	)
	defer cancel()

	wg.Add(PageRange)

	res := make([]*Result, 0)

	for x := 1; x <= PageRange; x++ {

		go func(page int) {
			// open chrome
			ctx, cancel := chromedp.NewContext(allocCtx)
			defer cancel()

			url := fmt.Sprintf("https://www.google.com.tw/search?q=%s&start=%s", keyWord, strconv.Itoa(((page - 1) * 10)))
			fmt.Printf("visit: %s\n", url)

			var htmlContent string
			err := chromedp.Run(
				ctx,
				chromedp.Navigate(url),
				chromedp.WaitReady(`#search div[id="rso"]`),
				chromedp.OuterHTML(`#search div[id="rso"]`, &htmlContent, chromedp.BySearch),
			)
			if err != nil {
				panic(err)
			}

			res = append(res, ParsingData(htmlContent, page)...)

			wg.Done()
		}(x)
	}
	wg.Wait()

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Page < res[j].Page
	})

	for _, r := range res {
		if r.Target {
			fmt.Printf("目前排名第%d名，第%d頁 第%d個\n", r.Rank, r.Page, r.Index)
			fmt.Printf("AT：%s\n", r.AT)
			fmt.Printf("Title：%s\n", r.Title)
			fmt.Printf("%s", time.Now().String())

			CrawlerRecord = append(CrawlerRecord, &CrawlerRecordData{
				AT:    r.AT,
				Rank:  r.Rank,
				Page:  r.Page,
				Index: r.Index,
				Title: r.Title,
			})
			writeFile(FileName, CrawlerRecord) // 寫入json檔
		}
	}
}

// ParsingData - 解析資料
func ParsingData(res string, page int) []*Result {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(res)))
	if err != nil {
		log.Fatal(err)
	}

	result := make([]*Result, 0)
	doc.Find(".yuRUbf").Each(func(i int, s *goquery.Selection) {
		itemTitle := s.Find(".LC20lb").Text()
		result = append(result, &Result{
			AT:    time.Now().String(),
			Rank:  ((page - 1) * 10) + i + 1,
			Page:  page,
			Index: i + 1,
			Title: s.Find(".LC20lb").Text(),
			Target: func() bool {
				return strings.Contains(itemTitle, "Relithe")
			}(),
		})
	})

	return result
}

// 讀取檔案
func openFile(FileName string) {
	_, err := os.Open("crawlerRecord.json") // 開啟檔案
	if err != nil {
		// 檢查檔案不存在則建立
		if os.IsNotExist(err) {
			newFile(FileName)
			return
		}
		log.Fatal(err)
	}
}

// 檢查檔案是否存在
func checkFileIsExist(FileName string) bool {
	_, err := os.Open(FileName) // 開啟檔案
	return !os.IsNotExist(err)
}

// 建立檔案
func newFile(FileName string) {
	f, err := os.Create(FileName)
	defer f.Close()
	if err != nil {
		fmt.Println("newFile faild：", err.Error())
	}
}

// 寫入資料
func writeFile(fileName string, list []*CrawlerRecordData) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Println("\nwriting data \n", list)
	spew.Dump(list)
	data := jsonMarshal(list)                    // 資料轉檔
	err = ioutil.WriteFile(fileName, data, 0644) // 寫入json檔
	if err != nil {
		log.Fatal(err)
	}
}

// 資料轉檔
func jsonMarshal(list []*CrawlerRecordData) []byte {
	data, err := json.Marshal(list)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
