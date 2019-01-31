package iffy_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/iffy"
	"github.com/loopfz/gadgeto/tonic"
)

func helloHandler(c *gin.Context) error { return nil }
func newFoo(c *gin.Context) error       { return nil }
func delFoo(c *gin.Context) error       { return nil }

type Foo struct{}

func ExpectValidFoo(r *http.Response, body string, respObject interface{}) error { return nil }

func Test_Tester_Run(t *testing.T) {
	// Instantiate & configure anything that implements http.Handler
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/hello", tonic.Handler(helloHandler, 200))
	r.POST("/foo", tonic.Handler(newFoo, 201))
	r.DELETE("/foo/:fooid", tonic.Handler(delFoo, 204))

	tester := iffy.NewTester(t, r)

	// Variadic list of checker functions = func(r *http.Response, body string, responseObject interface{}) error
	//
	// Checkers can use closures to trap checker-specific configs -> ExpectStatus(200)
	// Some are provided in the iffy package, but you can use your own Checker functions
	tester.AddCall("helloworld", "GET", "/hello?who=world", "").Checkers(iffy.ExpectStatus(200), iffy.ExpectJSONFields("msg", "bla"))
	tester.AddCall("badhello", "GET", "/hello", "").Checkers(iffy.ExpectStatus(400))

	// Optionally, pass an instantiated response object ( &Foo{} )
	// The response body will be unmarshaled into it, then it will be presented to the Checker functions (parameter 'responseObject')
	// That way your custom checkers can directly use your business objects (ExpectValidFoo)
	tester.AddCall("createfoo", "POST", "/foo", `{"bar": "baz"}`).ResponseObject(&Foo{}).Checkers(iffy.ExpectStatus(201), ExpectValidFoo)

	// You can template query string and/or body using partial results from previous calls
	// e.g.: delete the object that was created in a previous step
	tester.AddCall("deletefoo", "DELETE", "/foo/{{.createfoo.id}}", "").Checkers(iffy.ExpectStatus(204))

	tester.Run()
}
