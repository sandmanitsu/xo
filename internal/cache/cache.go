package cache

import "xo/internal/repository"

var Cache map[string]repository.User

func InitCache() {
	Cache = make(map[string]repository.User)
}
