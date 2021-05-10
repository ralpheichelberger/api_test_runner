package main

// ResponseT is the response type
type ResponseT struct {
	Description string
	Code        int
	Message     string
	Body        string
}

// SingleTest is the single test type
type SingleTest struct {
	Request struct {
		ID          int
		Skip        bool
		Description string
		Method      string
		Path        string
		Params      []string
		ContentType string
		Body        string
		Store       string
		Omit        bool
	}
	Response ResponseT
}

// APITest is the structur of the test run files
type APITest struct {
	FileName string
	Host     struct {
		Protocol string
		Name     string
		Port     string
		Prefix   string
	}
	TestRun struct {
		Skip        bool
		Name        string
		Description string
		Tests       []SingleTest
	}
}

// TestResult holts test results
type TestResult struct {
	success int
	failed  int
	skipped int
}
