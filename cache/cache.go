package cache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	expiration *int64
	// CacheInstance -
	CacheInstance *cache.Cache
)

func init() {
	expiration = kingpin.Flag("expiration", "default expiration time of ? hours").Default("6").Int64()
	CacheInstance = cache.New(time.Duration(*expiration)*time.Hour, 10*time.Minute)
	// cd.Set("22", "2", time.Minute)
	fmt.Println("==cache==")
	log.Infof("expiration: %d Hours", *expiration)
}
