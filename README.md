# API Tester Runner

A simple test runner, thats lets you descripe API test and helps find differences when responses have changed.

The testruns are written as yaml files and have a simple structure.

All test files are located in the sub directory 'tests'.

They have the following structure:

host:
  protocol: <http/https>
  name: <hostname>
  port: <port number>
testrun:
  skip: <true/false>
  name: <name of the testrun>
  description: <description of the testrun>
  tests: <array of tests>
    - request: 
        description:  <description of the test>
        method: GET
        path: <path of the api call and url params>
      response:
        code: <response code e.g. 200 or 404>
        message: <message text>
    - request:
        skip: <true/false>
        id: <int ID of the test>
        description: <description ot the test request>
        method: POST
        path: add_receipt_item
        contenttype: application/json
        body: <JSON request body>
      response:
        description: <description of what is expected as response>
        body: <JSON response string>

After testruns are written the api_test_runner can be started with following flags:
	-v verbose
        all descriptions, request and response bodies are printed out
	-q quiet
        only the results and failed tests are printed out
	-r test ID
        Runs only a test with the given test ID. Multiple tests can share the same ID.
	-s stop on first error
        stops testing after first failed test


# api_test_runner