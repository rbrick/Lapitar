package server

import (
	"image"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/FrozenOrb/lapitar/mc"
	"github.com/FrozenOrb/lapitar/server/cache"
	"github.com/FrozenOrb/lapitar/util"
	"github.com/zenazn/goji/web"
)

func printError(err error, message ...interface{}) {
	if err == nil {
		return
	}

	log.Println(message...)
	log.Println(err)
}

func parseSize(c web.C, def int) (result int) {
	size := c.URLParams["size"]
	result, err := strconv.Atoi(size)
	if err != nil {
		printError(err, "Failed to parse size", size)
		return def
	}
	return
}

func loadSkinMeta(name string, watch *util.StopWatch) (skin cache.SkinMeta) {
	watch.Mark()

	if mc.IsName(name) {
		skin = cache.FetchByName(name)
	} else {
		name = mc.ParseUUID(name)
		if mc.IsUUID(name) {
			skin = cache.Fetch(name)
		} else {
			skin = cache.FallbackByName(name)
		}
	}

	log.Println("Loaded skin:", skin.Profile().Name(), watch)
	return
}

const (
	keepCache    = 24 * time.Hour
	cacheControl = "max-age=86400" // 24*60*60, one day in seconds
)

func serveCached(w http.ResponseWriter, r *http.Request, meta cache.SkinMeta) bool {
	if tag := r.Header.Get("If-None-Match"); tag == meta.ID() {
		prepareResponse(w, r, meta)
		w.WriteHeader(http.StatusNotModified)
		return true
	}

	return false
}

func prepareResponse(w http.ResponseWriter, r *http.Request, meta cache.SkinMeta) {
	w.Header().Add("Cache-Control", cacheControl)
	w.Header().Add("Expires", time.Now().Add(keepCache).UTC().Format(http.TimeFormat))
	w.Header().Add("ETag", meta.ID())
	w.Header().Add("Last-Modified", meta.Timestamp().UTC().Format(http.TimeFormat))
}

func sendResult(w http.ResponseWriter, profile mc.Profile, result image.Image, watch *util.StopWatch) (err error) {
	watch.Mark()
	w.Header().Add("Content-Type", "image/png")
	err = png.Encode(w, result)
	if err == nil {
		log.Println("Sent response:", profile.Name(), watch)
	} else {
		printError(err, "Failed to send response:", profile.Name(), watch)
	}

	return
}
