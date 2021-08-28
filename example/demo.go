package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/wangdashuaihenshuai/doubanFinder/douban"
)

func main() {
	douban.SetBaiduCookie("you baidu reqeust cookie")
	const name = "2001.A.Space.Odyssey.1968.2001太空漫游.双语字幕.HR-HDTV.AC3.1024X576.x264.mkv"
	formatInfo := douban.ParseMovieInfo(name)
	fmt.Printf("开始搜索 %s", formatInfo.Name)
	id, ok, err := douban.SearchMovieId(formatInfo.Name)
	if err != nil {
		log.Fatal(err)
	}

	if !ok {
		log.Fatal(errors.New("搜索不到"))
	}

	details, err := douban.GetMovieDetailById(id)
	if err != nil {
		log.Fatal(err)
	}

	bs, _ := json.Marshal((details))
	fmt.Println(string(bs))
}
