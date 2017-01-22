package cache

import (
	"log"
	"time"

	"github.com/FrozenOrb/lapitar/mc"
)

type memorySkinCache struct {
	sessionServer string
}

func Memory(sessionServer string) SkinCache {
	result := &memorySkinCache{
		sessionServer,
	}
	return result
}

type memorySkinMeta struct {
	mc.SkinMeta
	fallback  SkinMeta
	profile   mc.Profile
	timestamp time.Time
	skin      mc.Skin
}

func (cache *memorySkinCache) FetchByName(realName string) (skin SkinMeta) {
	name := mc.ToLower(realName)

	profile, err := mc.FetchProfile(name)
	if profile == nil || err != nil {
		log.Println("Failed to fetch UUID for", realName, err)
		skin = FallbackByName(realName)
		return
	}

	skin = Fetch(profile.UUID())
	return
}

func (cache *memorySkinCache) Fetch(uuid string) (skin SkinMeta) {
	uuid = mc.ParseUUID(uuid)

	sk, err := mc.FetchSkin(cache.sessionServer, uuid)
	if sk == nil || err != nil {
		log.Println("Failed to skin profile for", uuid, err)
		skin = FallbackByUUID(uuid)
		return
	}

	skin = &memorySkinMeta{
		SkinMeta: sk.Skin(),
		profile:  sk.Profile(),
	}
	return
}

func (meta *memorySkinMeta) Profile() mc.Profile {
	return meta.profile
}

func (meta *memorySkinMeta) Fetch() (m SkinMeta, sk mc.Skin) {
	m = meta

	if meta.fallback == nil {
		sk = meta.skin
	} else {
		m, sk = meta.fallback.Fetch()
	}

	if sk != nil {
		return
	}

	// We need to download the skin first
	if meta.fallback == nil {
		sk = meta.skin
	} else {
		m, sk = meta.fallback.Fetch()
	}

	if sk != nil {
		return
	}

	sk, err := meta.SkinMeta.Download()
	if err != nil {
		// Meh, we can't download the skin right now
		log.Println("Failed to fetch skin from", meta.URL(), err)
		m, sk = Fallback(meta.Profile()).Fetch()
		meta.fallback = m
		return
	}

	meta.skin = sk
	return
}

func (meta *memorySkinMeta) Download() (mc.Skin, error) {
	_, sk := meta.Fetch()
	return sk, nil
}

func (meta *memorySkinMeta) Timestamp() time.Time {
	return meta.timestamp
}
