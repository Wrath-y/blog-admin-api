package friendlink

import (
	"blog-admin-api/pkg/goredis"
)

const (
	ListStrKey = "blog:friend-link:list"
)

func DelList() error {
	return goredis.Client.Del(ListStrKey).Err()
}
