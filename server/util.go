package server

import (
	"image"
	"image/png"
	"log"
	"net/http"
	"strconv"

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
