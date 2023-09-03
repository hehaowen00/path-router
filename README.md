# path-router-go

path router is a http router for go using a trie data structure for routing

```go
r := pathrouter.NewRouter()

r.Get("/", func(w http.ResponseWriter, req *http.Request, ps *pathrouter.Params) {
    fmt.Fprintf(w, "Hello, World!")
})

log.Fatal(http.ListenAndServe(":8000", r))
```

- middleware

middleware can be added using the `use` method

```go
func logger(next HandlerFunc) HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
        log.Println(r.Method, r.URL.String())
        next(w, r)
    }
}

r := pathrouter.NewRouter()
r.Use(logger)
```

- groups

groups in path router allow url path prefixing and middleware scoped to urls

```go
r := pathrouter.NewRouter()

r.Group("/api", func(g *Group) {
    g.Use(logger)

    // GET /api/users
    g.Get("/users", getUsers)
})
```

- url params

url parameters can be specified by prefixing the url element with a `:`

```
/hello/:user

/hello/hehaowen00
```

the matched values will be stored in the parameter struct in the request context

```go
func (w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
    value := ps.Get(":user") // hehaowen00
}
```

url parameters will only match one url segment unlike wildcards

- wildcards

wildcards are specified using `*` and must be at the end of the url

urls with wildcards will terminate matching and consume the rest of the url

```
/static/js/*

/static/js/app/index.min.js
```

Wildcards are stored in the parameter struct in the request context

```go
func (w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
    value := ps.Get("*") // app/index.min.js
}
```

- routing conflicts

when a parameter and wildcard node are in the same position within the URL,
only the parameter node will be matched and the wildcard route ignored

```
/special/:id
/special/*
```

adding another route with the same special path segment will overwrite the previous

```
/special/:a
/special/:b
```

- compatibility

path router implements the `Handle` method to add `http.Handler` routes

```go
r := pathrouter.NewRouter()

r.Handle("GET", "/handle", http.HandleFunc(func (w http.ResponseWriter, r *http.Request) {
    ps := r.Context().Value(pathrouter.ParamsKey).(*pathrouter.Params)
})
```
