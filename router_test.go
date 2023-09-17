package pathrouter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterGet(t *testing.T) {
	success := false

	url := "/get"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success = true
	}

	router := NewRouter()
	router.Get(url, h)
	router.ServeHTTP(w, req)

	if !success {
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

func TestRouterParams(t *testing.T) {
	success := false

	url := "/param/:value"
	req := httptest.NewRequest(http.MethodGet, "/param/true", nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		value := ps.Get(r.URL.Path, "value")
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

func TestRouterHandle(t *testing.T) {
	success := false

	url := "/param/:value"
	req := httptest.NewRequest(http.MethodGet, "/param/true", nil)
	w := httptest.NewRecorder()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ps := r.Context().Value(ParamsKey).(*Params)
		value := ps.Get(r.URL.Path, "value")
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

	url := "/api/test"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		success = true
	}

	nilHandler := func(w http.ResponseWriter, r *http.Request, ps *Params) {
	}

	router := NewRouter()
	router.Get("/hello", nilHandler)
	router.Group("/api", func(g *Group) {
		g.Get("/test", h)
	})

	router.ServeHTTP(w, req)

	if !success {
		t.FailNow()
	}
}
