/*
Ghttp stands for good http...
Is a VERY VERY simple http framework inspired in Echo.
Echo framework return error on its handlers, and thats something that I believe to be really clever. So thats what I did here
*/
package ghttp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Ghttp struct {
	*http.ServeMux
	Context
	middleware Middleware
}

type Context struct {
	http.ResponseWriter
	*http.Request
}

type httpError struct {
	Err    error
	Status int
}

type HandlerFunc func(Context) error

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Allows to chain middlewares
func chain(middleware ...Middleware) Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		for _, m := range middleware {
			handler = m(handler)
		}
		return handler
	}
}

// Add CORS headers to the request
func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func New() *Ghttp {
	return &Ghttp{
		http.NewServeMux(),
		Context{},
		chain(),
	}
}

// Middleware function for basic authentication
func BasicAuth(next http.HandlerFunc, username, password string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		credentials, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		pair := strings.SplitN(string(credentials), ":", 2)
		if len(pair) != 2 || pair[0] != username || pair[1] != password {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w, r}
}

// Add CORS to the chain of middlewares
func (g *Ghttp) CORS() *Ghttp {
	g.middleware = chain(cors)
	return g
}

// Allows to control the status been returned in the error, else it will default to internal server error
func newHttpError(err error, status int) httpError {
	return httpError{
		Err:    err,
		Status: status,
	}
}

func (e httpError) Error() string {
	return e.Err.Error()
}

func (g *Ghttp) Start(port string) {
	log.Printf("LISTENING ON PORT %s", port)
	log.Fatal(http.ListenAndServe(port, g))
}

// Simplifies the return for a function
func (c *Context) JSON(obj interface{}) error {
	return json.NewEncoder(c.ResponseWriter).Encode(obj)
}

// Returns a error with a custom status code
func (c *Context) FAIL(err error, status int) httpError {
	return newHttpError(err, status)
}

// This is the handler that allows all others handlers to simply return an error
func (g *Ghttp) defaultHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		g.Context = Context{w, r}
		err := handler(g.Context)
		if err != nil {
			if ghttperr, ok := err.(httpError); ok {
				w.WriteHeader(ghttperr.Status)
				g.JSON(map[string]string{"message": ghttperr.Error()})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				g.JSON(map[string]string{"message": err.Error()})
			}
		}
	}
}

func (g *Ghttp) GET(path string, handler HandlerFunc) {
	g.HandleFunc(fmt.Sprintf("GET %s", path), g.middleware(g.defaultHandler(handler)))
}

func (g *Ghttp) POST(path string, handler HandlerFunc) {
	g.HandleFunc(fmt.Sprintf("POST %s", path), g.middleware(g.defaultHandler(handler)))
}

func (g *Ghttp) PUT(path string, handler HandlerFunc) {
	g.HandleFunc(fmt.Sprintf("PUT %s", path), g.middleware(g.defaultHandler(handler)))
}

func (g *Ghttp) DELETE(path string, handler HandlerFunc) {
	g.HandleFunc(fmt.Sprintf("DELETE %s", path), g.middleware(g.defaultHandler(handler)))
}
