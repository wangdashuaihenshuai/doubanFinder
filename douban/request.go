package douban

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var cookie = ""

// 使用百度来搜索豆瓣,豆瓣自己搜索难爬,这里设置百度cookie不然会被墙
func SetBaiduCookie(c string) {
	cookie = c
}

func request(url string, qs map[string]string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	for k, v := range qs {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	if err != nil {
		return "", err
	}
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7,zh-TW;q=0.6")
	req.Header.Set("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("reqeust ret http code %d not 200", resp.StatusCode)
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyText), nil
}
