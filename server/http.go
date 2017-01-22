package server

import (
	"crypto/sha1"
	"flag"
	"log"
	"net/http"
	"os"

	"fmt"

	"github.com/FrozenOrb/lapitar"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

const (
	sha1Length = 40
)

var (
	defaults *renderDefaults
	useHash  bool
	hashSalt string
	//decoder  = schema.NewDecoder()
)

func start(conf *config, www string) {
	defaults = conf.Defaults
	useHash = conf.UseHash
	hashSalt = conf.HashSalt
	flag.Set("bind", conf.Address) // Uh, I guess that's a bit strange
	if conf.Proxy {
		goji.Insert(middleware.RealIP, middleware.Logger)
	}

	goji.Use(serveLapitar)

	register("/skin/:player", serveSkin)

	register("/face/:player", serveFaceNormal)
	register("/face/:size/:player", serveFaceWithSize)
	register("/helm/:player", serveHelmNormal)
	register("/helm/:size/:player", serveHelmWithSize)

	register("/head/:player", serveHeadNormal)
	register("/head/:size/:player", serveHeadWithSize)
	register("/portrait/:player", servePortraitNormal)
	register("/portrait/:size/:player", servePortraitWithSize)
	register("/player/:player", servePlayerNormal)
	register("/player/:size/:player", servePlayerWithSize)

	if exists(www) {
		goji.Get("/*", http.FileServer(http.Dir(www)))
	} else {
		log.Println("Failed to find website files at", www)
	}

	goji.Serve()
}

func exists(dir string) bool {
	stat, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func register(pattern string, handler web.HandlerFunc) {
	if useHash {
		wrappedHandler := func(c web.C, w http.ResponseWriter, r *http.Request) {
			if len(c.URLParams["hash"]) != sha1Length {
				http.NotFound(w, r)
				return
			}
			uri := r.RequestURI[sha1Length+1:] // uri without hash part
			hash := sha1.Sum([]byte(uri + hashSalt))
			if fmt.Sprintf("%x", hash) != c.URLParams["hash"] {
				http.NotFound(w, r)
				return
			}
			handler(c, w, r)
		}
		goji.Get("/:hash"+pattern+".png", wrappedHandler)
		goji.Get("/:hash"+pattern, wrappedHandler)
	} else {
		goji.Get(pattern+".png", handler)
		goji.Get(pattern, handler)
	}
}

func serveLapitar(c *web.C, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", lapitar.DisplayName)
		h.ServeHTTP(w, r)
	})
}
