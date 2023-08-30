package articleseo

import (
	"blog-admin-api/pkg/goredis"
	"fmt"
)

const (
	ListStrKey = "blog:article-seo:%d"
)

func DelList(articleID int) error {
	return goredis.Client.Del(fmt.Sprintf(ListStrKey, articleID)).Err()
}
