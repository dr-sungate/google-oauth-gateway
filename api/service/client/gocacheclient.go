package client

import (
	"fmt"
	log "github.com/dr-sungate/google-oauth-gateway/api/service/logger"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	NUM_INSTANCE       = 5
	PURGE_EXPIRED_TIME = 10
)

type GoCacheClient struct {
	CacheInstance []*cache.Cache
}

// インスタンス作成用のメソッド
func NewGoCacheClient(ttl int) *GoCacheClient {
	ret := &GoCacheClient{
		CacheInstance: make([]*cache.Cache, NUM_INSTANCE, NUM_INSTANCE),
	}
	for i := 0; i < NUM_INSTANCE; i++ {
		ret.CacheInstance[i] = cache.New(time.Duration(ttl)*time.Minute, PURGE_EXPIRED_TIME*time.Minute)
	}
	return ret
}

func (cc *GoCacheClient) Get(key string) (interface{}, bool) {
	log.Info(key)
	return cc.getInstance(key).Get(key)
}

func (cc *GoCacheClient) Set(key string, value interface{}, ttl int) {
	log.Info(key)
	cc.getInstance(key).Set(key, value, time.Duration(ttl)*time.Minute)
}

func (cc *GoCacheClient) getInstance(key string) *cache.Cache {
	// djb2アルゴリズム
	i, hash := 0, uint32(5381)
	for _, c := range key {
		hash = ((hash << 5) + hash) + uint32(c)
	}
	i = int(hash) % NUM_INSTANCE
	log.Info(fmt.Sprintf("get_instance: %d", i))
	return cc.CacheInstance[i]
}
