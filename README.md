E2E Testing w/ Go and chromedp
====

May 8 2019

This repo has a contrived demo app for using chromedp and dockertest to drive both integration and end-to-end tests from go test.

Clone this repo into $GOPATH/src/github.com/fortytw2/gdgberlin (or run go get)

### run tests

`go test -v -failfast -run='TestDemoApp/fish/basic-no-chromedp'`

`go test -v -failfast -run='TestDemoApp/root'`

LICENSE
=====

MIT, see LICENSE for details