package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// TESTDIR points to the directory with all test yaml files
const TESTDIR = "tests"

var store Store
var verbose *bool
var quiet *bool
var stopOnError *bool
var requestID *int

func init() {
	verbose = flag.Bool("v", false, "verbose")
	quiet = flag.Bool("q", false, "quiet")
	requestID = flag.Int("r", 0, "only request ID")
	stopOnError = flag.Bool("s", false, "stop on first error")
	flag.Parse()
	store = *NewStore()
}

func main() {
	testFiles, err := ioutil.ReadDir(TESTDIR)
	if err != nil {
		log.Fatal(err)
	}
	var totalFailed, totalSuccess, totalSkipped int = 0, 0, 0
	startTime := time.Now()
	for _, testFile := range testFiles {
		if !strings.HasPrefix(testFile.Name(), "test") {
			continue
		}
		fileContent, err := ioutil.ReadFile(TESTDIR + "/" + testFile.Name())
		ERRORFATAL(err)
		apiTest := APITest{}
		err = yaml.Unmarshal(fileContent, &apiTest)
		ERRORFATAL(err)
		apiTest.FileName = testFile.Name()
		intermittedStartTime := time.Now()
		result := execTests(apiTest)
		if !*quiet && *requestID == 0 && !apiTest.TestRun.Skip {
			fmt.Printf("\n%s\n", testFile.Name())
			fmt.Println("--------------------------")
			green.Printf("%d TESTS OK\n", result.success)
			failedstr := fmt.Sprintf("%d TESTS FAILED", result.failed)
			if result.failed > 0 {
				red.Println(failedstr + " :(")
			} else {
				green.Println(failedstr + " :)")
			}
			if result.skipped > 0 {
				yellow.Printf("%d TESTS SKIPPED\n", result.skipped)
			}
			fmt.Printf("--------------------------\n")
			fmt.Printf("Test duration: %s\n", time.Since(intermittedStartTime).String())
		}
		totalSuccess += result.success
		totalFailed += result.failed
		totalSkipped += result.skipped
		if totalFailed > 0 && *stopOnError {
			break
		}
	}
	fmt.Println()
	fmt.Println("--------------------------")
	fmt.Println("---      ALL TESTS     ---")
	fmt.Println("--------------------------")
	green.Printf("%d TESTS OK\n", totalSuccess)
	failedstr := fmt.Sprintf("%d TESTS FAILED", totalFailed)
	if totalFailed > 0 {
		red.Println(failedstr + " :(")
	} else {
		green.Println(failedstr + " :)")
	}
	if totalSkipped > 0 {
		yellow.Printf("%d TESTS SKIPPED\n", totalSkipped)
	}
	fmt.Printf("--------------------------\n")
	fmt.Printf("Total duration: %s\n", time.Since(startTime).String())

}

func execTests(apiTest APITest) TestResult {
	url := buildURL(apiTest)
	success := 0
	failed := 0
	skipped := 0
	if !*quiet {
		fmt.Printf("\n\nFilename: %s\n", apiTest.FileName)
		fmt.Println("--------------------------")
		t := apiTest.TestRun
		fmt.Printf("\n%s\n(%s)\n", t.Name, t.Description)
		fmt.Printf("running %d Tests\n", len(t.Tests))
		h := apiTest.Host
		fmt.Printf("Host:\n  Name: %s\n  Port: %s\n  Prot: %s\n  Pref: %s\n", h.Name, h.Port, h.Protocol, h.Prefix)
		if apiTest.TestRun.Skip {
			fmt.Println("Test skipped")
		}
	}
	if apiTest.TestRun.Skip {
		return TestResult{skipped: len(apiTest.TestRun.Tests)}
	}
	for _, test := range apiTest.TestRun.Tests {
		if test.Request.Skip {
			yellow.Printf("Skipping '%s'\n\n", test.Request.Description)
			skipped++
		} else {
			if *requestID > 0 && *requestID != test.Request.ID {
				// skip request
				continue
			}
			if test.Request.Skip {
				continue
			}
			if execRequest(url, test) {
				success++
			} else {
				failed++
				if *stopOnError {
					break
				}
			}
		}
	}
	return TestResult{success: success, failed: failed, skipped: skipped}
}
