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

// 計算排名
func computedRank(ctx context.Context, ch chan bool) {
	var err error
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
			ctx.Done()
			return
		}
		ch <- false
	}
}

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

	//建立一個上下文
	// ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	// ch := make(chan int32)

	// go func(ctx context.Context, ch chan int32) {
	// 	ch <- 1
	// }(ctx, ch)

	// go func(ctx context.Context, ch chan int32) {
	// 	ctx.Done()
	// }(ctx, ch)

	// var count int32 = 0
	// for {
	// 	select {
	// 	case value := <-ch:
	// 		count += value
	// 	case <-ctx.Done():
	// 		fmt.Println("finish")
	// 		break
	// 	}
	// }

	// spew.Dump(count)

	ch := make(chan string)
	Rank := 0 // 排名
	hasCurrentItem := false
	keyWord := "二手精品" // 搜尋關鍵字

	for i := 1; i <= 10; i++ {
		time.Sleep(time.Second)
		go func(i int) {
			fmt.Printf("now in Page: %d\n", i)
			// navigate
			err = chromedp.Run(ctx,
				chromedp.Navigate("https://www.google.com.tw/search?q="+keyWord+"&start="+strconv.Itoa(((i-1)*10))),
				// chromedp.Sleep(1*time.Second),
				chromedp.WaitReady(`#search div[id="rso"]`),
			)
			// if i == 1 {
			// 	err = chromedp.Run(ctx,
			// 		chromedp.Navigate("https://www.google.com.tw/search?q="+keyWord+"&start="+((i-1)*10)),
			// 		// chromedp.Sleep(1*time.Second),
			// 		chromedp.WaitReady(`#search div[id="rso"]`),
			// 	)
			// } else {
			// 	pageSelection := `a[aria-label="Page ` + strconv.Itoa(i) + `"]`
			// 	err = chromedp.Run(ctx,
			// 		chromedp.Click(pageSelection),
			// 		// chromedp.Sleep(1*time.Second),
			// 		chromedp.WaitReady(`#search div[id="rso"]`),
			// 	)
			// }
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
				ctx.Done()
				return
			}
			ch <- "Not Found"
		}(i)
	}

	for {
		select {
		case v := <-ch:
			fmt.Println(v)
		case <-ctx.Done():
			fmt.Println("finish")
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
