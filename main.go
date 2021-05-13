// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/c9s/gomon/logger"
	"github.com/chromedp/chromedp"
)

func main() {
	options := []chromedp.ExecAllocatorOption{
		chromedp.ExecPath("/Users/ian/Desktop/Chromium.app/Contents/MacOS/Chromium"),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"), // 不加載圖片
		chromedp.Flag("headless", true),                        // 無頭模式
		chromedp.WindowSize(1920, 1080),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"),
	}

	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	//建立一個上下文，超時時間為40s
	ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	Rank := 0 // 排名
	hasCurrentItem := false
	keyWord := "二手精品" // 搜尋關鍵字

	for i := 1; i <= 20; i++ {
		fmt.Printf("now in Page: %d\n", i)
		// navigate
		if i == 1 {
			err = chromedp.Run(ctx,
				chromedp.Navigate("https://www.google.com.tw/search?q="+keyWord),
				chromedp.Sleep(1*time.Second),
				chromedp.WaitVisible(`#search div[id="rso"]`),
			)
		} else {
			pageSelection := `a[aria-label="Page ` + strconv.Itoa(i) + `"]`
			err = chromedp.Run(ctx,
				chromedp.Click(pageSelection),
				chromedp.Sleep(1*time.Second),
				chromedp.WaitReady(`#search div[id="rso"]`),
			)
		}
		if err != nil {
			logger.Info("Run err :", err)
			return
		}

		// outer source data
		htmlContent, runErr := outerPageSourceData(ctx)
		if runErr != nil {
			logger.Info("Outer err : %v", runErr)
			return
		}

		sleepErr := chromedp.Run(ctx,
			chromedp.Sleep(2*time.Second),
		)
		if sleepErr != nil {
			log.Fatal(sleepErr)
		}

		Rank, hasCurrentItem = ParsingData(htmlContent, Rank, i)
		if hasCurrentItem {
			return
		}
	}
}

// 輸出page source data
func outerPageSourceData(ctx context.Context) (string, error) {
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.WaitVisible(`#search div[id="rso"]`),
		chromedp.OuterHTML(`#search div[id="rso"]`, &htmlContent, chromedp.BySearch),
	)

	return htmlContent, err
}

// 解析資料
func ParsingData(res string, rank int, page int) (int, bool) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(res)))
	if err != nil {
		log.Fatal(err)
	}

	var itemTitle string
	hasCurrentItem := false
	doc.Find(".yuRUbf").Each(func(i int, s *goquery.Selection) {
		itemTitle = s.Find(".LC20lb").Text()
		isCurrentItem := strings.Contains(itemTitle, "Relithe") // 模糊搜尋是否含有關鍵字
		rank++
		if isCurrentItem {
			hasCurrentItem = true
			fmt.Print("\n")
			fmt.Printf("Rank：%d,\nPage：%d,\nIndex：%d,\nTitle：%s\n", rank, page, (i + 1), itemTitle)
		}
		time.Sleep(1 * time.Second)
	})

	return rank, hasCurrentItem
}
