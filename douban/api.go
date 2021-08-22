package douban

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var subjectIdReg = regexp.MustCompile(`https\:\/\/movie\.douban\.com\/subject\/(\d+)`)
var idReg = regexp.MustCompile(`(\d+)`)

const shrotTimeLayout = "2006-01-02"

func parseId(str string) string {
	info := idReg.FindStringSubmatch(str)
	if info == nil || len(info) <= 1 {
		return ""
	}
	return info[0]
}

func parseTime(str string) *time.Time {
	t, err := time.Parse(shrotTimeLayout, str)
	if err != nil {
		return nil
	}
	return &t
}

func parseNum(str string) int {
	info := idReg.FindStringSubmatch(str)
	if info == nil || len(info) <= 1 {
		return 0
	}
	num, err := strconv.Atoi(info[0])
	if err != nil {
		return 0
	}

	return num
}

func SearchMovieId(name string) (string, bool, error) {
	baseUrl := "https://www.baidu.com/s"
	qs := map[string]string{
		"wd": name + " site:movie.douban.com",
	}
	html, err := request(baseUrl, qs)
	if err != nil {
		return "", false, err
	}

	subjectUrl := subjectIdReg.FindString(html)
	if subjectUrl != "" {
		words := strings.Split(subjectUrl, "/")
		return words[len(words)-1], true, nil
	}

	return "", false, nil
}

type Person struct {
	ID   string
	Name string
}

type MovieDetail struct {
	ID            string
	Start         float32
	Title         string
	Comment       string
	Description   string
	Cover         string
	Rate          string
	Subtype       string // Movie TV
	Directors     []*Person
	Authors       []*Person
	Actors        []*Person
	Duration      int
	Region        string
	Types         []string
	ReleaseYear   int
	PublishedTime *time.Time
}

func GetMovieDetailById(id string) (*MovieDetail, error) {
	sd, err := getShortDetail(id)
	if err != nil {
		return nil, err
	}

	pd, err := getPageDetail(id)
	if err != nil {
		return nil, err
	}

	return &MovieDetail{
		ID:            id,
		Start:         sd.Start,
		Title:         pd.Name,
		Comment:       sd.ShortComment.Content,
		Description:   pd.Description,
		Cover:         pd.Image,
		Rate:          sd.Rate,
		Subtype:       sd.Subtype,
		Directors:     toPersons(pd.Director),
		Authors:       toPersons(pd.Author),
		Actors:        toPersons(pd.Actor),
		Duration:      parseNum(sd.Duration),
		Region:        sd.Region,
		Types:         sd.Types,
		ReleaseYear:   parseNum(sd.ReleaseYear),
		PublishedTime: parseTime(pd.DatePublished),
	}, nil
}

type sortDetail struct {
	Start        float32 `json:"star"`
	Title        string  `json:"title"`
	URL          string  `json:"url"`
	Rate         string  `json:"rate"`
	ShortComment struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	} `json:"short_comment"`
	Subtype     string   `json:"subtype"`
	Duration    string   `json:"duration"`
	Region      string   `json:"region"`
	Types       []string `json:"types"`
	ReleaseYear string   `json:"release_year"`
}

type sortRes struct {
	Res     int        `json:"r"`
	Subject sortDetail `json:"subject"`
}

func getShortDetail(id string) (*sortDetail, error) {
	url := "https://movie.douban.com/j/subject_abstract"
	jsonStr, err := request(url, map[string]string{"subject_id": id})
	if err != nil {
		return nil, fmt.Errorf("请求页面 %s 错误%w", url, err)
	}
	res := sortRes{}
	err = json.Unmarshal([]byte(jsonStr), &res)
	if err != nil {
		return nil, fmt.Errorf("解析详情页面json: %s 错误%w", jsonStr, err)
	}

	if res.Res != 0 {
		return nil, fmt.Errorf("请求豆瓣详情接口错误 %s", jsonStr)
	}

	return &res.Subject, nil
}

type pagePerson struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func (pp *pagePerson) toPerson() *Person {
	return &Person{
		ID:   parseId(pp.URL),
		Name: pp.Name,
	}
}

func toPersons(ps []*pagePerson) []*Person {
	ret := []*Person{}
	for _, p := range ps {
		ret = append(ret, p.toPerson())
	}

	return ret
}

type pageDetail struct {
	Name          string        `json:"name"`
	URL           string        `json:"url"`
	Image         string        `json:"image"`
	Description   string        `json:"description"`
	Director      []*pagePerson `json:"director"`
	Author        []*pagePerson `json:"author"`
	Actor         []*pagePerson `json:"actor"`
	DatePublished string        `json:"datePublished"`
}

var pageDetailJsonReg = regexp.MustCompile(`<script\ type\=\"application\/ld\+json\">([\s\S]*?)<\/script>`)

func getPageDetail(id string) (*pageDetail, error) {
	url := fmt.Sprintf("https://movie.douban.com/subject/%s/", id)
	htmlStr, err := request(url, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("请求页面 %s 错误", url)
	}

	jsonStrs := pageDetailJsonReg.FindStringSubmatch(htmlStr)
	if jsonStrs == nil || len(jsonStrs) <= 1 {
		return nil, errors.New("匹配详情页面json错误")
	}

	res := pageDetail{}
	err = json.Unmarshal([]byte(jsonStrs[1]), &res)
	if err != nil {
		return nil, fmt.Errorf("解析详情页面json: %s 错误%w", jsonStrs[1], err)
	}

	return &res, nil
}
