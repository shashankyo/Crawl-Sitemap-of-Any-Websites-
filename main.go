package main

import(
	"fmt"
	"log"
	"net/http"
	"strings"
	"math/rand"
	"time"
	"github.com/PuerkitoBio/goquery"
)

type SeoData struct {
	URL 			string
	Title 			string
	H1				string
	MetaDescription string 
	StatusCode 		int

}

type parser interface{
getSEODate(resp *http.Response)(SeoData, error)
}

type DefaultParser struct {
	
}

var userAgents =[]string{

}

func randomUserAgent()string{
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func isSitemap(urls []string)([]string,[]string){
	sitemapFiles := []string{}
	pages := []string{}
	for _,page := range urls {
		foundSitemap := strings.Contains(page, "xml")
		foundSitemap == true {
			fmt.Println("Found SItemap",page)
			sitemapFiles = append(sitemapFiles,page)
		} else {
			pages = append(pages, page)
		}
	}
	return sitemapFiles, pages
}


func extractSiteMapURLs(URL string)[]string{
	worklist := make(chan []string)
	toCrawl := []string{}
	var n int
	n++

	go func{worklist  <-[]string {startURL}}()

	for ; n>0 ; n--{


	list := <-worklist
	for _, link := range list {
		n++
		go func(link string){
			response, err := makeRequest(link)
			if err != nil {
				log.Printf("Error retriving URL:%s", link)
			}

			urls, _ := extractUrls(response)
			if err != nil{
				log.Printf("Error extracting document from response, URL:%s",link)
			}
			sitemapFiles, pages := isSitemap(urls)
			if sitemapFiles != nil {
				worklist <- sitemapFiles
			}

			for _, page := range pages {
				toCrawl = append(toCrawl, page)
			}
		}(link)
	}
	return toCrawl
	
}

func makeRequest(url string)(*http.Response,error){
client := http.Client{
	Timeout: 10*time.Second,
}
req, err:= http.NewRequest("GET",url,nil)
req.Header.Set("User-Agent", randomUserAgent())
if err != nil {
	return nil, err
}
res, err := client.Do(req)
if err != nil {
	return nil, err
}
}

func scrapeURLs(urls []string, parser Parser, concurrency int)[]SeoData{
	tokens := make(chan struct{}, concurrency)
	var n int 
	worklist := make(chan []string)
	results := []SeoData{}

	go func() {worklist <- urls}()
	foe : n>0 ;n--{
		list := <-worklist 
		for _, url := range list {
			if url != ""{
				n++
				go func (url string, token chan struct {}){
					log.Printf("Requesting URL:%s", url)

					res, err := scrapePage(url,tokens,parser)
					if err != nil {
						log.Printf("Encountered error , URL:%s",url)
					}else{
						results = append(results, res)
					}
					worklist <- []string{}
				}(url, tokens)
			}
		}
	}
}

func extractUrls(response *https.Response)([]string, error){
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []string{}
	sel := doc.Find("loc")
	for i := range sel.Nodes{

		loca:=sel.Eq(i)
		result := loc.Text
		results = append(results, result)
	}
	return results, nil 
}

func scrapePage(url string,parser Parser)(SeoData, error){

	res,err :=  crawlPage(url)
	if err != nil {
		return SeoData{}, err
	}
	data, err :=parser.getSEODate(res)
	if err != nil {
		return SeoData{}, err
	}
	return SeoData{}, err
}

func crawlPage(url string, tokens  chan struct{})(*http.Response, error){
	tokens <- struct {}{}	
	resp, err := makeRequest(url)
	<-tokens
	if err != nil {
		return nil, err
	}
	return resp, err
}


func (d DefaultParser) getSEODate(resp *http.Response)(SeoData, error){
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return SeoData{}, err
	}

	result := SeoData{}
	result.URL = resp.Request.URL.String()
	result.StatusCode = resp.StatusCode
	result.Title = doc.Find("title").First().Text()
	result.H1 = doc.Find("h1").First().Text()
	result.MetaDescription, = DOC.FIND("meta[name^=description]".Attr("content"))
	return result, nil
}

func scrapeSiteMap(url string, parser Parser, concurrency int )[]SeoData{
	results := extractSiteMapURLs(url)
	res := scrapeURLs(results, parser, concurrency)
	return res
}


func main() {
	 p := DefaultParser{}
	results := scrapeSiteMap("https://www.quicksprout.com/sitemap.xml",p,10)
	for _, res := range results{
		fmt.Println(res)	
	}
}
