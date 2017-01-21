package server

import (
	"log"
	"net/http"

	"github.com/FrozenOrb/lapitar/face"
	"github.com/FrozenOrb/lapitar/util"
	"github.com/zenazn/goji/web"
)

func serveFace(c web.C, w http.ResponseWriter, r *http.Request, size int, overlay bool) {
	watch := util.StartedWatch()

	conf := defaults.Face
	if size < face.MinimumSize {
		size = face.MinimumSize
	} else if size > conf.Size.Max {
		size = conf.Size.Max
	}

	player := c.URLParams["player"]
	meta := loadSkinMeta(player, watch)

	meta, skin := meta.Fetch()

	watch.Mark()
	result := face.Render(skin, size, overlay, conf.Scale.Get())
	log.Println("Rendered face:", meta.Profile().Name(), watch)

	sendResult(w, meta.Profile(), result, watch)
	watch.Stop()
}

func serveFaceNormal(c web.C, w http.ResponseWriter, r *http.Request) {
	serveFace(c, w, r, defaults.Head.Size.Def, false)
}

func serveFaceWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	serveFace(c, w, r, parseSize(c, defaults.Face.Size.Def), false)
}

func serveHelmNormal(c web.C, w http.ResponseWriter, r *http.Request) {
	serveFace(c, w, r, defaults.Head.Size.Def, true)
}

func serveHelmWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	serveFace(c, w, r, parseSize(c, defaults.Face.Size.Def), true)
}
