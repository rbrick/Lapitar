package cache

import (
	"time"

	"github.com/FrozenOrb/lapitar/mc"
)

var skinCache SkinCache

type SkinCache interface {
	Fetch(uuid string) SkinMeta
	FetchByName(name string) SkinMeta
}

type SkinMeta interface {
	mc.SkinMeta
	Profile() mc.Profile
	Timestamp() time.Time
	Fetch() (SkinMeta, mc.Skin)
}

func Init(cache SkinCache) {
	skinCache = cache
}

func Fetch(uuid string) SkinMeta {
	return skinCache.Fetch(uuid)
}

func FetchByName(name string) SkinMeta {
	return skinCache.FetchByName(name)
}
