package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

var whiteUnderlined color.Color = *color.New(color.FgWhite).Add(color.Underline)
var blue color.Color = *color.New(color.FgBlue)
var green color.Color = *color.New(color.FgHiGreen)
var yellow color.Color = *color.New(color.FgYellow)
var redBG color.Color = *color.New(color.BgRed).Add(color.Bold)
var cyan color.Color = *color.New(color.FgCyan)
var red color.Color = *color.New(color.FgRed)

func execRequest(url string, test SingleTest) bool {
	req := test.Request
	res := test.Response
	uri := url + fmt.Sprintf("/%s", req.Path)
	if len(req.Params) > 0 {
		uri = fmt.Sprintf("%s?", uri)
		for _, param := range req.Params {
			uri = fmt.Sprintf("%s&%s", uri, param)
		}
	}
	if *verbose {
		whiteUnderlined.Printf("Request ID: %3d - '%s'\n", req.ID, req.Description)
		blue.Printf("requesting ")
		yellow.Printf("%s: '%s'\n", req.Method, uri)
	}
	m := "GET"
	if len(req.Method) > 2 {
		m = strings.ToUpper(req.Method)
	}
	var err error
	var resp *http.Response
	switch m {
	case "GET":
		resp, err = http.Get(uri)
	case "POST":
		var request string
		request, err = format(req.Body)
		if *verbose {
			blue.Println("request body:")
			yellow.Printf("%s\n", request)
		}
		ERRORFATAL(err)
		bodyReader := bytes.NewBuffer([]byte(request))
		resp, err = http.Post(uri, req.ContentType, bodyReader)
	}
	ERRORFATAL(err)
	if resp == nil {
		log.Fatal("*** No response from server")
	}
	body, err := ioutil.ReadAll(resp.Body)
	ERRORFATAL(err)
	bodystr := string(body)
	if *verbose {
		whiteUnderlined.Printf("%s\n", res.Description)
		blue.Printf("response:\n")
		cyan.Printf("%s\n", bodystr)
	}
	success := false
	if len(res.Body) > 0 {
		expect, err := format(res.Body)
		expect = escapeSquareBracket(expect)
		ERRORFATAL(err)
		if testExpectation(bodystr, expect) {
			success = true
		} else {
			printDiff(expect, bodystr, &test, findDiff(bodystr, expect))
		}
		saveResult(req.Store, bodystr)
	} else {
		var result ResponseT
		err = json.Unmarshal(body, &result)
		ERRORFATAL(err)
		if reflect.DeepEqual(result, test.Response) {
			success = true
		} else {
			printDiff(test.Response, result, &test, 0)
		}
		resultJSON, err := json.Marshal(result)
		if err != nil {
			red.Println("trouble marshalling result object")
		} else {
			saveResult(req.Store, string(resultJSON))
		}
	}
	return success
}

func printDiff(exp interface{}, got interface{}, test *SingleTest, pos int) {
	red.Printf("*** FAILED: ID: %3d\n", test.Request.ID)
	red.Println("*")
	red.Print("* ")
	blue.Printf("expected:\n")
	red.Print("* ")
	whiteUnderlined.Println((*test).Response.Description)
	red.Print("* ")
	expect := fmt.Sprintf("%s", exp)
	yellow.Printf("'%v", expect[:pos])
	redBG.Printf(expect[pos : pos+1])
	yellow.Printf("%v'\n", expect[pos+1:])
	red.Println("*")
	red.Print("* ")
	blue.Println("actual:")
	red.Print("* ")
	cyan.Printf("'%v'", got)
	red.Print("\n***\n\n")
}
