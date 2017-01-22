package server

import (
	"log"
	"net/http"

	"github.com/FrozenOrb/lapitar/render"
	"github.com/FrozenOrb/lapitar/util"
	"github.com/zenazn/goji/web"
)

func serveRender(c web.C, w http.ResponseWriter, r *http.Request, size int, conf *renderConfig, portrait, full bool) {
	watch := util.StartedWatch()

	if size < render.MinimumSize {
		size = render.MinimumSize
	} else if size > conf.Size.Max {
		size = conf.Size.Max
	}

	player := c.URLParams["player"]
	meta := loadSkinMeta(player, watch)

	meta, skin := meta.Fetch()

	watch.Mark()
	sizeY := size
	if full {
		sizeY = int(float64(sizeY) * 1.625)
	}

	result, err := render.Render(skin, conf.Angle, conf.Tilt, conf.Zoom, size, sizeY, conf.SuperSampling, portrait, full, conf.Overlay, conf.Shadow, conf.Lighting, conf.Scale.Get())
	if err == nil {
		log.Println("Rendered head:", meta.Profile().Name(), watch)
	} else {
		printError(err, "Failed to render head:", meta.Profile().Name(), watch)
		return
	}

	sendResult(w, meta.Profile(), result, watch)
	watch.Stop()
}

func serveHeadNormal(c web.C, w http.ResponseWriter, r *http.Request) {
	serveRender(c, w, r, defaults.Head.Size.Def, defaults.Head, false, false)
}

func serveHeadWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	serveRender(c, w, r, parseSize(c, defaults.Head.Size.Def), defaults.Head, false, false)
}

func servePortraitNormal(c web.C, w http.ResponseWriter, r *http.Request) {
	// serveRender(c, w, r, defaults.Portrait.Size.Def, defaults.Portrait, true, false)
}

func servePortraitWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	// serveRender(c, w, r, parseSize(c, defaults.Portrait.Size.Def), defaults.Portrait, true, false)
}

func servePlayerNormal(c web.C, w http.ResponseWriter, r *http.Request) {
	serveRender(c, w, r, defaults.Body.Size.Def, defaults.Body, false, true)
}

func servePlayerWithSize(c web.C, w http.ResponseWriter, r *http.Request) {
	serveRender(c, w, r, parseSize(c, defaults.Body.Size.Def), defaults.Body, false, true)
}
