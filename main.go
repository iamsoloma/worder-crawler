package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TinajXD/worder-crawler/config"
	"github.com/gocolly/colly"
	"github.com/saintfish/chardet"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

func main() {
	fmt.Println("Starting...")
	cfg := config.GetConf()
	url := "https://ru.wikipedia.org/wiki/Челябинск"
	//url := "http://az.lib.ru/t/tolstoj_lew_nikolaewich/text_0040.shtml"
	data := fetch(url, cfg.UserAgent)
	tags := tagger(data)
	links := linker(url, cfg.UserAgent)
	for _, tag := range tags{
		if tag != "" {
			fmt.Println(fmt.Sprint(len(tag)) + " : " + tag)
		}
	}
	for name, link := range links{
		if link != "" {
			fmt.Printf("Link found: %s -> %s\n", name, link)
		}
	}
}

func fetch (url, userAgent string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err!=nil{
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err!=nil{
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println(err)
	}
	return string(body)
}

func tagger(strHtml string) (tags []string) {
	detector := chardet.NewHtmlDetector()
	result, err := detector.DetectBest([]byte(strHtml))
	if err != nil {
    	return
	}
	tokenizer := html.NewTokenizer(strings.NewReader(strHtml))
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return
			}
			fmt.Printf("Error: %v", tokenizer.Err())
			return
		}

		switch token.Data {
		case "script":
			tokenizer.Next()
		case "/script":
			tokenizer.Next()
		case "style":
			tokenizer.Next()
		case "link":
			tokenizer.Next()
		}

		if tokenType == html.TextToken {
			content := strings.TrimSpace(token.Data)
			clearTag := convertToUTF8(content, result.Charset)
			tags = append(tags, clearTag)
		}
		
	}
}

func linker(url, userAgent string) (links map[string]string){
	c := colly.NewCollector(
		colly.DetectCharset(),
		colly.UserAgent(userAgent),
	)
	links = make(map[string]string)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link != "" && e.Request.AbsoluteURL(link) != ""{
			//fmt.Printf("Link found: %q -> %s\n", e.Text, e.Request.AbsoluteURL(link))
			links[e.Text] = e.Request.AbsoluteURL(url)
		}
	})
	c.Visit(url)
	return links
}

func convertToUTF8(str string, origEncoding string) string {
    strBytes := []byte(str)
    byteReader := bytes.NewReader(strBytes)
    reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
    strBytes, _ = io.ReadAll(reader)
    return string(strBytes)
}