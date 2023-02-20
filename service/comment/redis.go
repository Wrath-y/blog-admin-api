package comment

import (
	"blog-admin-api/pkg/goredis"
	"fmt"
)

const (
	ListStrKey  = "blog:comment:list:%d:%d"
	CountStrKey = "blog:comment:count:%d"
)

func DelList() error {
	var cursor uint64
	for {
		keys, nextCursor, err := goredis.Client.Scan(cursor, "blog:comment:list:*", 10).Result()
		if err != nil {
			return err
		}
		for _, key := range keys {
			err := goredis.Client.Del(key).Err()
			if err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

func DelByArticleId(articleId int) error {
	return goredis.Client.Del(fmt.Sprintf(CountStrKey, articleId)).Err()
}
