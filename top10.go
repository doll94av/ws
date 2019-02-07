/*
*Send a GET request to reddit to get the top 10 images from a given subreddit.
*
*Usage:
*go build top10.go
*top10.exe "https://old.reddit.com/r/wallpapers/"
*
*/

package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"time"
	"os"
	"strings"
	"net/http"
	"io"
//	"math/rand"
	"strconv"
)

type item struct {
	StoryURL string
	Source string
	comments string
	CrawledAt time.Time
	Comments string
	Title string
	
}



func DownloadFile(url string) error {
	
	fmt.Println(url)
	//random number as save file
	//randomName := rand.Intn(1000)
	randomName := count
	var x = "images/" + strconv.Itoa(randomName) + ".jpg"
	//create file
	out, err := os.Create(x)
	if err != nil {
		return err
	}
	defer out.Close()

	//get data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//write data
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	//imageTE.AppendText(walk.NewIconFromImage(x))
	return nil

}

//globals count to make sure we only download 10, and visitedUrl to keep track of
//where we've been to avoid duplicates
var count = 0
var visitedUrl map[string]bool




func main() {


	//initalize map
	visitedUrl = make(map[string]bool)

	stories := []item{}
	outputDir := "/images"
	c :=colly.NewCollector(
		colly.AllowedDomains("old.reddit.com"),
		colly.UserAgent("Chrome:com.learngo.top10download:v1 (by /u/myHoneyBaked)"),
		colly.Async(true),
	)




	//attached functions
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		r.Ctx.Put("url", r.URL.String())

	})
	c.OnResponse(func(r *colly.Response){
		if strings.Index(r.Headers.Get("Content-type"), "image") > -1 {
			r.Save(outputDir + r.FileName())
			return
		}
		fmt.Println("Visited",r.Request.URL)


	})

	c.OnHTML(".top-matter",func(e *colly.HTMLElement) {
		temp := item{}
		temp.StoryURL = e.ChildAttr("a[data-event-action=title]", "href")
		temp.Source = "https://old.reddit.com/r/wallpapers/"
		temp.Title = e.ChildText("a[data-event-action=comments]")
		temp.Comments = e.ChildAttr("a[data-event-action=comments]", "href")
		temp.CrawledAt = time.Now()
		stories = append(stories, temp)

	})

	c.OnHTML("span.next-button", func(h *colly.HTMLElement) {
		t := h.ChildAttr("a", "href")
		c.Visit(t)
	})

	//url to download image is the child of a <div>, always ends in .jpg 
	c.OnHTML("div", func(i *colly.HTMLElement){

		t := i.ChildAttr("a", "href")

		if strings.Contains(t, ".jpg") {
			if count < 10 {
				if _, exists := visitedUrl[t]; !exists {
					err := DownloadFile(t)
					visitedUrl[t] = true
					count++
					if err != nil {
						panic(err)
					}
				}

			}
		}
	})

reddits := os.Args[1:]
	for _, reddit := range reddits {
		c.Visit(reddit)
	}
	c.Wait()
}
