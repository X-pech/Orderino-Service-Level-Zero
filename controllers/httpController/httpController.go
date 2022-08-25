package httpController

import (
	"log"
	"net/http"
)

const jsonPath = "/json/"
const webfacePath = "/face/"

type HTTPController struct {
	apiCallback func(string) string
}

func New(apiCallback func(string) string) HTTPController {
	result := new(HTTPController)
	result.apiCallback = apiCallback
	http.HandleFunc(jsonPath, result.apiHandler)
	http.HandleFunc(webfacePath, result.webfaceHandler)
	return *result
}

func (hc *HTTPController) Start(port string) error {
	return http.ListenAndServe(":"+port, nil)
}

func (hc *HTTPController) webfaceHandler(rw http.ResponseWriter, rq *http.Request) {
	if rq.Method != "GET" {
		msg := "Wrong MEthod"
		log.Println(msg)
		rw.Write([]byte("Nothing here"))
		return
	}

	http.ServeFile(rw, rq, "./webface/index.html")
}

func (hc *HTTPController) apiHandler(rw http.ResponseWriter, rq *http.Request) {
	if rq.Method != "GET" {
		msg := "[HTTPController] Wrong request method"
		log.Println(msg)
		rw.Write([]byte(msg))
		return
	}

	data := hc.apiCallback(rq.URL.Query().Get("order_uid"))
	if data == "" {
		data = "{}"
	}
	rw.Write([]byte(data))
	rw.Header().Add("Content-Type", "application/json")
}
