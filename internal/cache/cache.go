package cache

import "xo/internal/repository"

var Cache map[string]repository.User
var CacheRooms []int

func InitCache() {
	Cache = make(map[string]repository.User)
}
