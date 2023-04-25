package worker

import (
	"context"
	"io"
	"log"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"os"
	"fmt"
)

type Job struct {
	Cancel context.CancelFunc
	// ETCD key, length : 160
	ID string
	JobVal
}

type JobVal struct {
	URL     string `json:"url"`
	Pattern string `json:"pattern"`
	// SM OracleWiter related, not sure yet
}

// TODO: add timeout
func (j Job) Scrap() (string, error) {
	log.Println("start scraping")
	res, err := http.Get(j.URL)
	if err != nil {
		return "", err
	}
	data, err := j.resolve(res)
	if err != nil {
		return "", err
	}
	return data, nil
}

// Not implemented
// TODO: add resolver
func (j Job) resolve(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// TODO: ADD RESOLVER
	return string(body), nil
}

//Create directories and write data
func CreateDateDir(path string,name string, content string) {
	_, err := os.Stat(path)
  
	if err != nil {
		fmt.Println("文件不存在！")

		if os.IsNotExist(err) {
			err := os.Mkdir(path, os.ModePerm)
			
			if err != nil {
			fmt.Printf("创建失败![%v]\n", err)
			}else {
				fmt.Println("创建成功!")
			}			
		}		
	}

	file, _ := os.OpenFile(path + "/" + name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	file.Write([]byte(content + "\n \n"))		
}

//web crawler
func crawler() {
	client := http.Client{}
    req, err := http.NewRequest("GET", "https://movie.douban.com/chart", nil)
	
    if err != nil {
        fmt.Println(err)
    }

	req.Header.Set("Connection", "keep-alive")
    req.Header.Set("Content-Type", "text/html; charset=utf-8")
    req.Header.Set("Keep-Alive", "timeout=30")
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	
	resp, err := client.Do(req)

    if err != nil {
        fmt.Println(err)
    }

    //Resolve the web page
    dom, err := goquery.NewDocumentFromReader(resp.Body)

    if err != nil {
        fmt.Println(err)
    }
	
	dom.Find("#content > div > div.article > div  tr.item").Each(func(i int, selection *goquery.Selection) {
		src, _ := selection.Find("td").First().Find("a.nbg > img").Attr("src")
		text := selection.Find("span").Text()
		CreateDateDir("/home/songguokele/桌面/GO/Data", "test.txt", text + "\n" + src )//File path
		fmt.Println(text + src)
	})
}
