// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/c9s/gomon/logger"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

func main() {
	options := []chromedp.ExecAllocatorOption{
		chromedp.ExecPath("/Users/ian/Desktop/Chromium.app/Contents/MacOS/Chromium"),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.Flag("headless", false),
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

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var htmlContent string
	runErr := chromedp.Run(ctx,
		// chromedp.Navigate(`https://www.google.com.tw/search?q=台中二手精品`),
		chromedp.Navigate(`https://www.google.com.tw`),
		chromedp.WaitVisible(`.SDkEP input[name="q"]`, chromedp.BySearch), // waiting for element exist

		chromedp.Sleep(1*time.Second),
		chromedp.SetValue(`.SDkEP input[name="q"]`, "台中二手精品", chromedp.BySearch),
		chromedp.SendKeys(`.SDkEP input[name="q"]`, kb.Enter), // 按下Enter
		chromedp.WaitVisible(`#search div[id="rso"]`),
		chromedp.OuterHTML(`#search div[id="rso"]`, &htmlContent, chromedp.BySearch),
	)
	if err != nil {
		logger.Info("Run err : %v\n", runErr)
		return
	}

	sleepErr := chromedp.Run(ctx,
		chromedp.Sleep(10*time.Second),
	)
	if sleepErr != nil {
		log.Fatal(sleepErr)
	}
	log.Printf("OK\n")
	log.Printf("1 \n")
	log.Printf("2 \n")
	log.Printf("3 \n")
	log.Printf("4 \n")
	log.Printf("5 \n")
	log.Printf("：" + htmlContent)
}

// //獲取網站上爬取的資料
// func GetHttpHtmlContent(url string, selector string, sel interface{}) (string, error) {
// 	options := []chromedp.ExecAllocatorOption{
// 		chromedp.Flag("headless", true), // debug使用
// 		chromedp.Flag("blink-settings", "imagesEnabled=false"),
// 		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"),
// 	}
// 	//初始化引數，先傳一個空的資料
// 	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

// 	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

// 	// create context
// 	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
// 	// 執行一個空task, 用提前建立Chrome例項
// 	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

// 	//建立一個上下文，超時時間為40s
// 	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 40*time.Second)
// 	defer cancel()

// 	var htmlContent string
// 	err := chromedp.Run(timeoutCtx,
// 		chromedp.Navigate(url),
// 		chromedp.WaitVisible(selector),
// 		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByJSPath),
// 	)
// 	if err != nil {
// 		logger.Info("Run err : %v\n", err)
// 		return "", err
// 	}
// 	//log.Println(htmlContent)

// 	return htmlContent, nil
// }
