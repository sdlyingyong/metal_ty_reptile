package controllers

import (
	"beego_reptile_ty/models"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"sync"
)

type CrawlController struct {
	beego.Controller
}

//@router /reptile/movie [get]
func (c *CrawlController) CrawlMovie() {
	var startUrl = "https://movie.douban.com/subject/33420285/?tag=%E7%83%AD%E9%97%A8&from=gaia"
	for {
		resp := httplib.Get(startUrl)
		resp.Header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:50.0) Gecko/20100101 Firefox/50.0")
		resp.Header("Cookie", `bid=gFP9qSgGTfA; __utma=30149280.1124851270.1482153600.1483055851.1483064193.8; __utmz=30149280.1482971588.4.2.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; ll="118221"; _pk_ref.100001.4cf6=%5B%22%22%2C%22%22%2C1483064193%2C%22https%3A%2F%2Fwww.douban.com%2F%22%5D; _pk_id.100001.4cf6=5afcf5e5496eab22.1482413017.7.1483066280.1483057909.; __utma=223695111.1636117731.1482413017.1483055857.1483064193.7; __utmz=223695111.1483055857.6.5.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; _vwo_uuid_v2=BDC2DBEDF8958EC838F9D9394CC5D9A0|2cc6ef7952be8c2d5408cb7c8cce2684; ap=1; viewed="1006073"; gr_user_id=e5c932fc-2af6-4861-8a4f-5d696f34570b; __utmc=30149280; __utmc=223695111; _pk_ses.100001.4cf6=*; __utmb=30149280.0.10.1483064193; __utmb=223695111.0.10.1483064193`)
		resD, err := resp.String()
		if err != nil {
			panic(err)
		}
		logs.Info(resD)
	}
	c.Ctx.WriteString("end of crawl!")
}

//@router /reptile/sync/go [get]
func (c *CrawlController) CrawlGoBlobSync() {
	//setting down article count
	down_count := 3
	//get all article list
	//save in db
	for i := 1; i < 10; i++ {
		_, err := getPageArticles(i)
		if err != nil {
			logs.Error(err)
		}
		logs.Info("all page list read end.")
	}
	//get all need request url from db
	var article_url_list []models.RequestLog
	o := orm.NewOrm()
	_, err := o.QueryTable("request_log").Filter("is_visited", 0).
		OrderBy("-id").All(&article_url_list)
	if err != nil {
		logs.Error(err)
	}
	logs.Info("count article urls is " + strconv.Itoa(len(article_url_list)))
	//three execute same time
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		for _, article_url := range article_url_list[:(down_count - 0)] {
			//get one page detail
			_, err := getArticleDetail(article_url.Url)
			if err != nil {
				logs.Error(err)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for _, article_url := range article_url_list[:(down_count - 1)] {
			//get one page detail
			_, err := getArticleDetail(article_url.Url)
			if err != nil {
				logs.Error(err)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for _, article_url := range article_url_list[:(down_count - 2)] {
			//get one page detail
			_, err := getArticleDetail(article_url.Url)
			if err != nil {
				logs.Error(err)
			}
		}
	}()
	wg.Wait()
	logs.Info("all article read end.")
}

//@router /reptile/go [get]
func (c *CrawlController) CrawlGoBlob() {
	//setting down article count
	down_count := 3
	//get all article list
	//save in db
	for i := 1; i < 10; i++ {
		_, err := getPageArticles(i)
		if err != nil {
			logs.Error(err)
		}
		logs.Info("all page list read end.")
	}
	//get all need request url from db
	var article_url_list []models.RequestLog
	o := orm.NewOrm()
	_, err := o.QueryTable("request_log").Filter("is_visited", 0).
		OrderBy("-id").All(&article_url_list)
	if err != nil {
		logs.Error(err)
	}
	logs.Info("count article urls is " + strconv.Itoa(len(article_url_list)))
	//save all article detail
	for _, article_url := range article_url_list[:down_count] {
		//get one page detail
		_, err := getArticleDetail(article_url.Url)
		if err != nil {
			logs.Error(err)
		}
	}
	logs.Info("all article read end.")
}

//@router	/migrate/blob	[get]
func (c *CrawlController) MigrateBlob() {
	res, _ := syncArticle(3)
	logs.Info("sync blob success.Sync count :" + strconv.Itoa(res))
}

func (c *CrawlController) MigrateBlobSync() {
	//import sync package
	var wg sync.WaitGroup
	//need three routine
	wg.Add(2)
	go func() {
		//when done set count --1
		defer wg.Done()
		_, ok := syncArticle(1)
		if !ok {
			logs.Error("do syncArticle error")
		}
	}()
	go func() {
		//when done set count --1
		defer wg.Done()
		_, ok := syncArticle(1)
		if !ok {
			logs.Error("do syncArticle error")
		}
	}()
	//wait wg count is 0
	wg.Wait()
	logs.Info("sync blob success.Sync count :")
}

func syncArticle(sync_count int) (int, bool) {
	//setting sync article count
	//sync_count := 1
	//get content data
	var article_list []models.CrawlWeb
	o := orm.NewOrm()
	_, err := o.QueryTable("crawl_web").
		Filter("is_sync", 0).
		OrderBy("-id").All(&article_list)
	if err != nil {
		logs.Error(err)
		return 0, false
	}
	if len(article_list) < sync_count {
		logs.Error("the sync count large than data in mysql.send " + strconv.Itoa(sync_count))
		logs.Error("get from data : " + strconv.Itoa(len(article_list)))
		return 0, false
	}
	//create article
	for _, row := range article_list[:sync_count] {
		createArticle(row.Title, row.Content+"\n origin : "+row.Url)
		//update is_sync
		o.QueryTable("crawl_web").Filter("url", row.Url).Update(orm.Params{"is_sync": 1})
	}
	return len(article_list), true
}

func createArticle(title string, content string) (int, bool) {
	//setting
	api_url := "http://localhost:8080/admin/api/article"
	//send to create article endpoint
	args := &fasthttp.Args{}
	args.Add("title", title)
	args.Add("content", content)
	status, _, err := fasthttp.Post(nil, api_url, args)
	if err != nil {
		logs.Error(err)
		return 0, false
	}
	if status != fasthttp.StatusOK {
		logs.Error("request fail,code : " + strconv.Itoa(status))
		return 0, false
	}
	logs.Info("success ")
	return 1, true
}

func sendJsonRequest(api_url string) (int, bool) {
	//json type
	req := &fasthttp.Request{}
	req.SetRequestURI(api_url)

	requestBody := []byte(`{"content":"test sync content","title":"test sync"}`)
	req.SetBody(requestBody)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")

	resp := &fasthttp.Response{}

	client := &fasthttp.Client{}

	err := client.Do(req, resp)
	if err != nil {
		logs.Error(err)
		return 0, false
	}
	if resp.StatusCode() != 200 {
		logs.Error("request err.status is :" + strconv.Itoa(resp.StatusCode()))
		logs.Error(resp.Body())
		return 0, false
	}
	return 1, true
}

/**
get article detail
content,title,origin_url
*/
func getArticleDetail(url string) (bool, error) {
	logs.Info("request url is :" + url)
	//var url = "https://learnku.com/go/t/62060"
	//article list
	res, err := http.Get(url)
	if err != nil {
		logs.Error(err)
		return false, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logs.Error("status code not 200 while request: " + url)
		return false, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logs.Error(err)
		return false, err
	}
	//parse article content
	content := doc.Find("div.ui.readme.markdown-body.content-body.fluidbox-content").Text()
	title := doc.Find("div.extra-padding").Find("div.pull-left").Find("span").Text()

	//save in crawl_web
	o := orm.NewOrm()
	crawl_web := models.CrawlWeb{Url: url, Content: content, Title: title}
	created, _, err := o.ReadOrCreate(&crawl_web, "Url")
	if err != nil {
		logs.Error(err)
		return false, err
	}
	if created {
		logs.Info("save in req log")
	} else {
		logs.Info("was exist")
	}
	//change request log status
	o.QueryTable("request_log").Filter("url", url).Update(orm.Params{"is_visited": 1})
	return true, nil
}

func getPageArticles(pageId int) (bool, error) {
	//get web site start page
	var startUrl = "https://learnku.com/go?page=" + strconv.Itoa(pageId)
	//get url
	logs.Info("start request :" + startUrl)
	res, err := http.Get(startUrl)
	if err != nil {
		logs.Error(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logs.Error("status code not 200 while request: " + startUrl)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logs.Error(err)
	}
	//article list save  in db
	o := orm.NewOrm()
	doc.Find("a.topic-title-wrap").Each(func(i int, selection *goquery.Selection) {
		article_url, ok := selection.Attr("href")
		if ok {
			fmt.Printf("Review %d %s \n", i, article_url)
			//if url not exist
			//save in request_log with no request
			req_log := models.RequestLog{Url: article_url}
			created, _, err := o.ReadOrCreate(&req_log, "Url")
			if err != nil {
				logs.Error(err)
			}
			if created {
				logs.Info("save in req log")
			} else {
				logs.Info("was exist")
			}
		}
	})
	logs.Info("request end")
	return true, nil
}

func hasExist(url string) (bool, error) {
	o := orm.NewOrm()
	req_log := models.RequestLog{Url: url}

	err := o.Read(&req_log)
	if err == orm.ErrNoRows {
		return false, nil
	} else {
		return true, nil
	}
}

//@router	/reptile/beego	[get]
func (c *CrawlController) CrawlBeeDoc() {
	var url = "http://beego.me/"
	result, _ := getData(url)
	logs.Info(result)
	logs.Info("success")
}

//@router	/reptile/fish	[get]
func (c *CrawlController) CrawlFish() {
	var url = "https://api.tophub.fun/v2/GetAllInfoGzip?id=60&page=0&type=pc"
	//wait work flow
	wg := sync.WaitGroup{}
	wg.Add(2)
	//open another request handle
	go func() {
		defer wg.Done()
		getData(url)
		//logs.Info(result)
		logs.Info("success")
	}()
	//open another request handle
	go func() {
		defer wg.Done()
		getData(url)
		//logs.Info(result)
		logs.Info("success")
	}()
	//wait for all go function execute end
	wg.Wait()
	//save all data in mysql

}

func getData(url string) (string, error) {
	//set request options
	logs.Info("start request")
	req := httplib.Get(url)
	req.Debug(true)
	//send request
	resp_str, err := req.String()
	if err != nil {
		logs.Error(err)
		return "", err
	}
	//notice
	logs.Info("response success")
	return resp_str, nil
}
