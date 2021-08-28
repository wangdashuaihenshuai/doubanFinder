package douban

import (
	"regexp"
	"strings"
)

type MovieNameInfo struct {
	Name        string
	EnglistName string
	Meta        []string
	Type        string
	Raw         string
}

var char = []string{}

var urlReg = regexp.MustCompile(`[a-zA-Z]+\.[a-zA-Z]+\.[com|cn|net]\]`)

var replaceRegs = []*regexp.Regexp{
	regexp.MustCompile(`\d+届-`),
	regexp.MustCompile(`\d+x\d+`),
	regexp.MustCompile(`\d+x\d+`),
}

var replaceWords = []string{
	"【十万度v信 shiwandus】",
	"【十万v信 shiwandus】",
	"【",
	"】",
	"-",
	"]",
	"[",
	"(",
	")",
}

var metas = []string{"aac", "10bit", "中字", "mnhd-frds", "3audio", "1080p", "x265", "x264", "2audio", "hd中英双字", "bd1080p", "x264", "chd_eng", "双语", "720p", "chi_eng", "bdrip", "双语", "字幕", "hr-hdtv", "双音轨", "ac3", "完整版", "加长版", "bluray", "x264", "国英音轨", "flac-cmct", "flac", "dvdrip", "unrated", "bluray", "ac3", "hr-hdtv", "4audios", "cmct", "dc", "repack", "人人影视"}

func includeMeta(word string) (string, bool) {
	for _, m := range metas {
		if strings.Contains(word, m) {
			return m, true
		}
	}

	return "", false
}

func splitChars(r rune) bool {
	return r == '.' || r == '(' || r == ')'
}

func splitName(name string) []string {
	ret := []string{}
	for _, v := range strings.FieldsFunc(name, splitChars) {
		if v != "" {
			ret = append(ret, v)
		}
	}

	if len(ret) <= 0 {
		return []string{name}
	}

	return ret
}

func replaceName(word string) string {
	for _, w := range replaceWords {
		word = strings.ReplaceAll(word, w, " ")
	}

	for _, r := range replaceRegs {
		word = r.ReplaceAllString(word, " ")
	}

	return word
}

func ParseMovieInfo(name string) *MovieNameInfo {
	name = strings.ToLower(name)
	name = urlReg.ReplaceAllString(name, "")
	words := splitName(name)
	meta := []string{}
	if len(words) <= 1 {
		return &MovieNameInfo{
			Name: replaceName(name),
			Meta: meta,
			Raw:  name,
		}
	}

	t := words[len(words)-1]
	if len(words) == 2 {
		return &MovieNameInfo{
			Name: replaceName(words[0]),
			Meta: meta,
			Raw:  name,
			Type: t,
		}
	}

	formatWords := []string{}
	for _, w := range words[0 : len(words)-1] {
		m, ok := includeMeta(w)
		if ok {
			meta = append(meta, m)
		} else {
			formatWords = append(formatWords, w)
		}
	}

	return &MovieNameInfo{
		Name: replaceName(strings.Join(formatWords, " ")),
		Meta: meta,
		Raw:  name,
		Type: t,
	}
}
