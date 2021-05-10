package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// ERRORFATAL logs the error and quits the program
func ERRORFATAL(err error) {
	if err != nil {
		color.Red("***")
		log.Fatal(err)
	}
}
func buildURL(t APITest) string {
	h := t.Host
	url := fmt.Sprintf("%s://%s", h.Protocol, h.Name)
	if len(h.Port) > 0 {
		url += fmt.Sprintf(":%s", h.Port)
	}
	return url
}

func saveResult(key string, data string) {
	if len(key) > 0 {
		defer store.lock.Unlock()
		store.lock.Lock()
		store.data[key] = string(data)
	}
}

func escapeSquareBracket(s string) string {
	var result string = s
	result = strings.ReplaceAll(result, "[", "\\[")
	result = strings.ReplaceAll(result, "µ", "[")
	return result
}

func format(s string) (string, error) {
	mLeft := strings.Split(s, "{{")
	if len(mLeft) == 1 {
		return s, nil
	}
	var result string = s
	mLeft = mLeft[1:]
	for _, m := range mLeft {
		mRight := strings.Split(m, "}}")
		name := strings.Split(mRight[0], ".")
		doFormat(name, &result)
	}
	return result, nil
}

func doFormat(name []string, result *string) {
	var child string
	var grandchild string
	parent := strings.TrimSpace(name[0])
	if len(name) > 1 {
		child = strings.TrimSpace(name[1])
		if len(name) > 2 {
			grandchild = strings.TrimSpace(name[2])
		}
	}
	parentData := store.data[parent]
	var dataMap map[string]interface{}
	err := json.Unmarshal([]byte(parentData), &dataMap)
	if err != nil {
		fmt.Println(err.Error())
	}
	switch (dataMap[child]).(type) {
	case string:
		value := fmt.Sprintf("%v", dataMap[child])
		*result = strings.Replace(*result, fmt.Sprintf("{{%s.%s}}", parent, child), value, 1)
	case interface{}:
		var ddmap map[string]interface{}
		ddmap = dataMap[child].(map[string]interface{})
		value := fmt.Sprintf("%v", ddmap[grandchild])
		*result = strings.Replace(*result, fmt.Sprintf("{{%s.%s.%s}}", parent, child, grandchild), value, 1)
	}
}

func testExpectation(result string, regext string) bool {
	r := regexp.MustCompile(regext)
	return r.MatchString(result)
}

func findDiff(result string, regext string) int {
	regexlen := len(regext) - 1
	var backSlash byte = 92
	for i := range regext {
		newRegex := regext[:regexlen-i]
		if newRegex[len(newRegex)-1] == backSlash {
			continue
		}
		r, err := regexp.Compile(newRegex)
		if err == nil {
			if r.MatchString(result) {
				return len(newRegex)
			}
		}
	}
	return 0
}
