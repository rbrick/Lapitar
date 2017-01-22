package server

import (
	"net/http"

	"github.com/FrozenOrb/lapitar/util"
	"github.com/zenazn/goji/web"
)

func serveSkin(c web.C, w http.ResponseWriter, r *http.Request) {
	watch := util.StartedWatch()

	player := c.URLParams["player"]
	meta := loadSkinMeta(player, watch)

	// Check if we can return 304 Not Modified
	if serveCached(w, r, meta) {
		return
	}

	meta, skin := meta.Fetch()
	prepareResponse(w, r, meta)

	sendResult(w, meta.Profile(), skin.Image(), watch)
	watch.Stop()
}
