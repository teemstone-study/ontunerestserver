# go backend test [![go](https://miro.medium.com/max/700/1*Ifpd_HtDiK9u6h68SZgNuA.png)](https://go.dev/)

> start go backend

I will fill in the future.

## Description
모자이크 화면의 데이터 정보를 저장 하고 업데이트 한다.

## Environment
* go lang 
* windows11
* vscode 

## Prerequisite



* [[router and dispatcher] gorilla mux](https://github.com/gorilla/mux)

---

  Package `gorilla/mux` implements a request router and dispatcher for matching incoming requests to
  their respective handler.

   The name mux stands for "HTTP request multiplexer". Like the standard `http.ServeMux`, `mux.Router` matches incoming requests against a list of registered routes   and calls a handler for the route that matches the URL or other conditions. The main features are:

* It implements the `http.Handler` interface so it is compatible with the standard `http.ServeMux`.
* Requests can be matched based on URL host, path, path prefix, schemes, header and query values, HTTP methods or using custom matchers.
* URL hosts, paths and query values can have variables with an optional regular expression.
* Registered URLs can be built, or "reversed", which helps maintaining references to resources.
* Routes can be used as subrouters: nested routes are only tested if the parent route matches. This is useful to define groups of routes that share common conditions like a host, a path prefix or other repeated attributes. As a bonus, this optimizes request matching.

---

* [[cors] rs/cors](https://github.com/rs/cors) 
  - CORS is a `net/http` handler implementing [Cross Origin Resource Sharing W3 specification](http://www.w3.org/TR/cors/) in Golang.

---
* [[middleware] negroni](https://github.com/urfave/negroni)

   Negroni is an idiomatic approach to web middleware in Go. It is tiny,
   non-intrusive, and encourages use of `net/http` Handlers.
   
   ## `negroni.Classic()`

	`negroni.Classic()` provides some default middleware that is useful for most
	applications:

	* [`negroni.Recovery`](#recovery) - Panic Recovery Middleware.
	* [`negroni.Logger`](#logger) - Request/Response Logger Middleware.
	* [`negroni.Static`](#static) - Static File serving under the "public"
	  directory.

	This makes it really easy to get started with some useful features from Negroni.
   
   
---   
* [[render] unrolled](https://github.com/unrolled/render)
  - Render is a package that provides functionality for easily rendering JSON, XML, text, binary data, and HTML templates.
---     
* [[.env] godotenv](https://github.com/joho/godotenv)
From the original Library:

> Storing configuration in the environment is one of the tenets of a twelve-factor app. Anything that is likely to change between deployment environments–such as resource handles for databases or credentials for external services–should be extracted from the code into environment variables.
>
> But it is not always practical to set environment variables on development machines or continuous integration servers where multiple projects are run. Dotenv load variables from a .env file into ENV when the environment is bootstrapped.

   It can be used as a library (for loading in env for your own daemons etc.) or as a bin command.

   There is test coverage and CI for both linuxish and Windows environments, but I make no guarantees about the bin version working on Windows.
  
---

* [[db] pq](https://github.com/lib/pq) 

  - A pure Go postgres driver for Go's database/sql package
   
---   

* [[tdd] goconvey](https://github.com/smartystreets/goconvey)

GoConvey supports the current versions of Go (see the official Go
[release policy](https://golang.org/doc/devel/release#policy)). Currently
this means Go 1.16 and Go 1.17 are supported.

**Features:**

- Directly integrates with `go test`
- Fully-automatic web UI (works with native Go tests, too)
- Huge suite of regression tests
- Shows test coverage
- Readable, colorized console output (understandable by any manager, IT or not)
- Test code generator
- Desktop notifications (optional)
- Immediately open problem lines in [Sublime Text](http://www.sublimetext.com) ([some assembly required](https://github.com/asuth/subl-handler))

[Quick start](https://github.com/smartystreets/goconvey/wiki#get-going-in-25-seconds)
-----------

Make a test, for example:

```go
package package_name

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given some integer with a starting value", t, func() {
		x := 1

		Convey("When the integer is incremented", func() {
			x++

			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 2)
			})
		})
	})
}
```


#### [In the browser](https://github.com/smartystreets/goconvey/wiki/Web-UI)

Start up the GoConvey web server at your project's path:

	$ $GOPATH/bin/goconvey

Then watch the test results display in your browser at:

	http://localhost:8080


If the browser doesn't open automatically, please click [http://localhost:8080](http://localhost:8080) to open manually.

There you have it.
![](http://d79i1fxsrar4t.cloudfront.net/goconvey.co/gc-1-dark.png)
As long as GoConvey is running, test results will automatically update in your browser window.

![](http://d79i1fxsrar4t.cloudfront.net/goconvey.co/gc-5-dark.png)
The design is responsive, so you can squish the browser real tight if you need to put it beside your code.


The [web UI](https://github.com/smartystreets/goconvey/wiki/Web-UI) supports traditional Go tests, so use it even if you're not using GoConvey tests.



#### [In the terminal](https://github.com/smartystreets/goconvey/wiki/Execution)

Just do what you do best:

    $ go test

Or if you want the output to include the story:

    $ go test -v


---

* [[tdd] assert](https://github.com/stretchr/testify/tree/master/assert)

[`assert`](http://godoc.org/github.com/stretchr/testify/assert "API documentation") package
-------------------------------------------------------------------------------------------

The `assert` package provides some helpful methods that allow you to write better test code in Go.

  * Prints friendly, easy to read failure descriptions
  * Allows for very readable code
  * Optionally annotate each assertion with a message

See it in action:

```go
package yours

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {

  // assert equality
  assert.Equal(t, 123, 123, "they should be equal")

  // assert inequality
  assert.NotEqual(t, 123, 456, "they should not be equal")

  // assert for nil (good for errors)
  assert.Nil(t, object)

  // assert for not nil (good when you expect something)
  if assert.NotNil(t, object) {

    // now we know that object isn't nil, we are safe to make
    // further assertions without causing any errors
    assert.Equal(t, "Something", object.Value)

  }

}
```

  * Every assert func takes the `testing.T` object as the first argument.  This is how it writes the errors out through the normal `go test` capabilities.
  * Every assert func returns a bool indicating whether the assertion was successful or not, this is useful for if you want to go on making further assertions under certain conditions.

if you assert many times, use the below:

```go
package yours

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
  assert := assert.New(t)

  // assert equality
  assert.Equal(123, 123, "they should be equal")

  // assert inequality
  assert.NotEqual(123, 456, "they should not be equal")

  // assert for nil (good for errors)
  assert.Nil(object)

  // assert for not nil (good when you expect something)
  if assert.NotNil(object) {

    // now we know that object isn't nil, we are safe to make
    // further assertions without causing any errors
    assert.Equal("Something", object.Value)
  }
}
```

## Files
* pages.go 
```
type DBHandler interface {
	GetPages() []*Page
	AddPage(page *Page) bool	
	UpdatePage(page *Page) bool		
	GetPage(index int) *Page
	DeletePage() bool
	Close()
}

func NewDBHandler(dbConn string) DBHandler {
	//return newMemHandler()
	//return newSqliteHandler(filepath)
	return newPgHandler(dbConn)
}
```
* 원하는 저장소 memory, db, file 등에 대하여 DBHandler interface 를 구현하면  다른 소스파일에 영향 없이 Adapter 처럼 변경 사용 가능 합니다.
  (memHandler.go, pgHandler.go) 
## Usage
* tdd
```
\bego\app>goconvey      

\bego\app>go test
```
