//
package main

/*
Packages must be imported:
    "core/common/page"
    "core/spider"
Pckages may be imported:
    "core/pipeline": scawler result persistent;
    "github.com/PuerkitoBio/goquery": html dom parser.
*/
import (
	"fmt"
	"strings"

	"github.com/Bulusideng/go_spider/core/common/page"
	"github.com/Bulusideng/go_spider/core/pipeline"
	"github.com/Bulusideng/go_spider/core/spider"
	"github.com/PuerkitoBio/goquery"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
	return &MyPageProcesser{}
}

// Parse html dom here and record the parse result that we want to Page.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (this *MyPageProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		println(p.Errormsg())
		return
	}

	query := p.GetHtmlParser()
	var urls []string
	//query.Find("div.d-inline-block.mb-1 h3 a").Each(func(i int, s *goquery.Selection) {
	query.Find("div#user-repositories-list ul li div h3 a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		urls = append(urls, "http://github.com/"+href)
		fmt.Printf("URL:%s\n", href)
	})
	// these urls will be saved and crawed by other coroutines.
	p.AddTargetRequests(urls, "html")

	name := query.Find(".entry-title .author").Text()
	name = strings.Trim(name, " \t\n")
	repository := query.Find(".entry-title .js-current-repository").Text()
	repository = strings.Trim(repository, " \t\n")
	//readme, _ := query.Find("#readme").Html()
	if name == "" {
		p.SetSkip(true)
	}
	// the entity we want to save by Pipeline
	p.AddField("author", name)
	p.AddField("project", repository)
	//p.AddField("readme", readme)
}

func (this *MyPageProcesser) Finish() {
	fmt.Printf("TODO:before end spider \r\n")
}

func main() {
	// Spider input:
	//  PageProcesser ;
	//  Task name used in Pipeline for record;
	spider.NewSpider(NewMyPageProcesser(), "TaskName").
		AddUrl("https://github.com/hu17889?tab=repositories", "html"). // Start url, html is the responce type ("html" or "json" or "jsonp" or "text")
		AddPipeline(pipeline.NewPipelineConsole()).                    // Print result on screen
		SetThreadnum(3).                                               // Crawl request by three Coroutines
		Run()
}
