// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	options := []chromedp.ExecAllocatorOption{
		chromedp.ExecPath("/Users/ian/Desktop/Chromium.app/Contents/MacOS/Chromium"),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("headless", true),
		chromedp.WindowSize(1280, 1024),
		chromedp.UserAgent(`Mozilla/5.0 (iPhone; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.25 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1`),
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
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var example string
	err1 := chromedp.Run(ctx,
		chromedp.Navigate(`https://golang.org/pkg/time/`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`body > footer`),
		// find and click "Expand All" link
		chromedp.Click(`#pkg-examples > div`, chromedp.NodeVisible),
		// retrieve the value of the textarea
		chromedp.Value(`#example_After .play .input textarea`, &example),
	)
	if err1 != nil {
		log.Fatal(err1)
	}
	log.Printf("Go's time.After example:\n%s", example)
}
