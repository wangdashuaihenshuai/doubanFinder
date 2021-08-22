# README

一个豆瓣电影搜索接口,自己管理电影文件使用.

## 特点

- 使用百度搜索豆瓣电影(豆瓣搜索接口难弄)
- 会格式化处理电影名字

example 
```go
	douban.SetBaiduCookie("you baidu reqeust cookie")
	const name = "2001.A.Space.Odyssey.1968.2001太空漫游.双语字幕.HR-HDTV.AC3.1024X576.x264.mkv"
	formatInfo := douban.FormatMovieName(name)
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
```

## 注意⚠️

- 只能作为非盈利以及学习使用
- 禁止任何商业使用
