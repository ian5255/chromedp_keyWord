// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	// PageRange -
	PageRange = 5

	// 搜尋關鍵字
	keyWord = "二手精品"
)

// Result -
type Result struct {
	page   int
	index  int
	target bool
	title  string
}

var wg sync.WaitGroup

func main() {
	options := []chromedp.ExecAllocatorOption{
		chromedp.ExecPath("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"), // 不加載圖片
		chromedp.Flag("headless", false),                       // 無頭模式
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

	Rank := 0 // 排
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

			res = append(res, ParsingData(htmlContent, Rank, page)...)

			wg.Done()
		}(x)
	}
	wg.Wait()

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].page < res[j].page
	})

	for _, r := range res {
		if r.target {
			fmt.Printf("我在第%d頁 第%d個\n", r.page, r.index)
		}
	}
}

// ParsingData - 解析資料
func ParsingData(res string, rank int, page int) []*Result {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(res)))
	if err != nil {
		log.Fatal(err)
	}

	result := make([]*Result, 0)
	doc.Find(".yuRUbf").Each(func(i int, s *goquery.Selection) {
		itemTitle := s.Find(".LC20lb").Text()
		result = append(result, &Result{
			page:  page,
			index: i + 1,
			title: s.Find(".LC20lb").Text(),
			target: func() bool {
				return strings.Contains(itemTitle, "Relithe")
			}(),
		})
	})

	return result
}
