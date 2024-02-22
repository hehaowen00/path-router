package pathrouter

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestRouterEmpty(t *testing.T) {
	url := "/"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(w, req)
}

func TestRouterDefault(t *testing.T) {
	success := false

	url := "/"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success = true
	}

	router := NewRouter()
	router.Get("/", h)
	router.ServeHTTP(w, req)

	if !success {
		t.FailNow()
	}
}

func TestRouterGet(t *testing.T) {
	success := false

	req1 := httptest.NewRequest("GET", "/b", nil)

	url := "/a"
	req2 := httptest.NewRequest("GET", url, nil)

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success = true
	}

	router := NewRouter()
	router.Get("/a/*", h)
	router.Get("/b", h)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req1)

	if !success {
		t.FailNow()
	}

	success = false

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req2)

	if success {
		t.FailNow()
	}
}

func TestRouterPost(t *testing.T) {
	success := false

	url := "/post"
	req := httptest.NewRequest(http.MethodPost, url, nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success = true
	}

	router := NewRouter()
	router.Post(url, h)
	router.ServeHTTP(w, req)

	if !success {
		t.FailNow()
	}
}

func TestRouterPut(t *testing.T) {
	success := false

	url := "/put"
	req := httptest.NewRequest(http.MethodPut, url, nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success = true
	}

	router := NewRouter()
	router.Put(url, h)
	router.ServeHTTP(w, req)

	if !success {
		t.FailNow()
	}
}

func TestRouterPatch(t *testing.T) {
	success := false

	url := "/patch"
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success = true
	}

	router := NewRouter()
	router.Patch(url, h)
	router.ServeHTTP(w, req)

	if !success {
		t.FailNow()
	}
}

func TestRouterDelete(t *testing.T) {
	success := false

	url := "/delete"
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success = true
	}

	router := NewRouter()
	router.Delete(url, h)
	router.ServeHTTP(w, req)

	if !success {
		t.FailNow()
	}
}

func TestRouterConnect(t *testing.T) {
	success := false

	url := "/connect"
	req := httptest.NewRequest(http.MethodConnect, url, nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success = true
	}

	router := NewRouter()
	router.Connect(url, h)
	router.ServeHTTP(w, req)

	if !success {
		t.FailNow()
	}
}

func TestRouterOptions(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
	}

	r := NewRouter()
	r.Get("/a", h)
	r.Post("/a", h)
	r.Put("/a", h)
	r.Patch("/a", h)
	r.Delete("/a", h)

	req := httptest.NewRequest(http.MethodOptions, "/a", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	methods := resp.Header.Get("Allow")
	methodsArray := strings.Split(strings.ReplaceAll(methods, " ", ""), ",")

	for _, method := range []string{"GET", "POST", "PUT", "PATCH", "DELETE"} {
		if !slices.Contains(methodsArray, method) {
			log.Printf("error missing %s in %v\n", method, methodsArray)
			t.FailNow()
		}
	}
}

func TestRouterParams(t *testing.T) {
	success := false

	url := "/param/:value"
	req := httptest.NewRequest(http.MethodGet, "/param/true", nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		value := ps.Get("value")
		success = value == "true"
	}

	router := NewRouter()
	router.Get("/param/*", h)
	router.Get(url, h)
	router.ServeHTTP(w, req)

	if !success {
		t.FailNow()
	}
}

func TestRouterParams2(t *testing.T) {
	success := false

	url := "/*"
	req := httptest.NewRequest(http.MethodGet, "/true", nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		value := ps.Get("*")
		success = value == "true"
	}

	router := NewRouter()
	router.Get(url, h)
	router.ServeHTTP(w, req)

	if !success {
		t.FailNow()
	}
}

func TestRouterParams3(t *testing.T) {
	success := false

	url := "/*"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		value := ps.Get("*")
		success = value == ""
	}

	router := NewRouter()
	router.Get(url, h)
	router.ServeHTTP(w, req)

	if !success {
		t.FailNow()
	}
}

func TestRouterHandle(t *testing.T) {
	success := false

	url := "/param/:value"
	req := httptest.NewRequest(http.MethodGet, "/param/true", nil)
	w := httptest.NewRecorder()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ps := r.Context().Value(ParamsKey).(*Params)
		value := ps.Get("value")
		success = value == "true"
	})

	router := NewRouter()
	router.Handle(http.MethodGet, url, h)
	router.ServeHTTP(w, req)

	if !success {
		t.FailNow()
	}
}

func TestRouterMiddleware(t *testing.T) {
	success := false

	url := "/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
	}

	middleware := func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, ps *Params) {
			success = true
		}
	}

	router := NewRouter()
	router.Use(middleware)
	router.Get(url, h)
	router.ServeHTTP(w, req)

	if !success {
		t.FailNow()
	}
}

func TestRouterGroup(t *testing.T) {
	success := false
	success2 := false
	routerMiddleware := false
	groupMiddleware := false

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success = true
	}

	h2 := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success2 = true
	}

	nilHandler := func(w http.ResponseWriter, r *http.Request, ps *Params) {
	}

	routerLevel := func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, ps *Params) {
			routerMiddleware = true
			next(w, r, ps)
		}
	}

	groupLevel := func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, ps *Params) {
			groupMiddleware = true
			next(w, r, ps)
		}
	}

	router := NewRouter()
	router.Use(routerLevel)

	router.Get("/hello", nilHandler)

	api := router.Scope("/api")
	api.Use(groupLevel)
	api.Get("/test", h)

	t1 := router.Scope("/t1")
	t1.Use(groupLevel)

	t2 := t1.Scope("/t2")
	t2.Get("/test", h2)

	url := "/api/test"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !routerMiddleware || !groupMiddleware || !success {
		log.Println(routerMiddleware, groupMiddleware, success)
		t.FailNow()
	}

	routerMiddleware = false
	groupMiddleware = false

	url = "/t1/t2/test"
	req = httptest.NewRequest(http.MethodGet, url, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !routerMiddleware || !groupMiddleware || !success2 {
		log.Println(routerMiddleware, groupMiddleware, success)
		t.FailNow()
	}
}

func TestRouterGzip(t *testing.T) {
	url := "/"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Accept-Encoding", "gzip,deflate")
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		w.Write([]byte("Hello, World!"))
	}

	router := NewRouter()
	router.Use(GzipMiddleware)
	router.Get("/", h)
	router.ServeHTTP(w, req)

	resp := w.Result()

	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.FailNow()
	}
	r, _ := gzip.NewReader(resp.Body)
	defer r.Close()

	bytes, err := io.ReadAll(r)
	if err != nil {
		t.FailNow()
	}

	if string(bytes) != "Hello, World!" {
		fmt.Println(string(bytes))
		t.FailNow()
	}
}

func TestPathOrder(t *testing.T) {
	success := false

	router := NewRouter()
	router.Get("/a/:b", func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success = true
	})
	router.Get("/a", func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success = true
	})

	req := httptest.NewRequest(http.MethodGet, "/a/b", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if !success {
		t.FailNow()
	}

	req = httptest.NewRequest(http.MethodGet, "/a", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if !success {
		t.FailNow()
	}
}
