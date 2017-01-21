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

	meta, skin := meta.Fetch()

	sendResult(w, meta.Profile(), skin.Image(), watch)
	watch.Stop()
}
