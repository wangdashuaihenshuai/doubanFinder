package douban

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var numRegs = []*regexp.Regexp{
	regexp.MustCompile(`s\d+e(\d+)`),
	regexp.MustCompile(`ep(\d+)`),
	regexp.MustCompile(`e(\d+)`),
	regexp.MustCompile(`\((\d+)\)`),
	regexp.MustCompile(`(\d+)集`),
}

var numReg2 = regexp.MustCompile(`(\d+)`)

var dmp = diffmatchpatch.New()

func ParseMovieNumStr(name string, nextName string) (string, string) {
	name = strings.ToLower(name)
	if nextName != "" {
		nextName = strings.ToLower(nextName)
	}

	for _, r := range numRegs {
		allStr, str := getRegexString(name, r)
		if str != "" {
			return allStr, str
		}
	}

	allStr, allNumStr := getRegexString(name, numReg2)
	if len(allNumStr) == len(name) {
		return allStr, allNumStr
	}

	if nextName != "" {
		diffNum := parseDiffNum(name, nextName)
		if diffNum >= 0 {
			str := strconv.Itoa(diffNum)
			return str, str
		}
	}

	return getRegexString(name, numReg2)
}

func ParseMovieNum(name string, nextName string) int {
	_, numStr := ParseMovieNumStr(name, nextName)
	return parseInt(numStr)
}

var numRune = map[rune]bool{
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
}

func isNumChar(char rune) bool {
	_, exist := numRune[char]
	return exist
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func getLeftNumStr(str string) string {
	ret := []rune{}
	for _, r := range str {
		if isNumChar(r) {
			ret = append(ret, r)
		} else {
			return string(ret)
		}
	}
	return string(ret)
}

func parseDiffNum(name string, nextName string) int {
	leftStr := getDiffNumLeft(name, nextName)
	rightStr := name[len(leftStr):]
	rightNumStr := getLeftNumStr(rightStr)
	leftNumStr := reverse(getLeftNumStr(reverse(leftStr)))
	return getParseInt(string(leftNumStr)+string(rightNumStr), numReg2)
}

func getDiffNumLeft(name string, nextName string) string {
	diffs := dmp.DiffMain(nextName, name, false)
	left := ""
	for _, d := range diffs {
		if d.Type == diffmatchpatch.DiffInsert {
			diffNum := getParseInt(d.Text, numReg2)
			if diffNum >= 0 {
				return left
			}
		}

		if d.Type == diffmatchpatch.DiffEqual {
			left = left + d.Text
		}
	}

	return left
}

func getRegexString(str string, r *regexp.Regexp) (string, string) {
	info := r.FindStringSubmatch(str)
	if len(info) <= 1 {
		return "", ""
	}

	return info[0], info[1]
}

func getParseInt(str string, r *regexp.Regexp) int {
	_, parseStr := getRegexString(str, r)

	if parseStr == "" {
		return -1
	}

	return parseInt(parseStr)
}

func parseInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return -1
	}

	return i
}
