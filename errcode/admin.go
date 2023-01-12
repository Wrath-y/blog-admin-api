package errcode

var (
	AdminNetworkBusy   = &ErrCode{40000, "网络繁忙，请稍后重试", ""}
	AdminInvalidToken  = &ErrCode{40001, "无效的token，请重试", ""}
	AdminInvalidParam  = &ErrCode{40002, "无效的参数", ""}
	AdminInvalidSign   = &ErrCode{40003, "无效的签名", ""}
	AdminBodyTooLarge  = &ErrCode{40004, "请求消息体过大", ""}
	AdminUserGetFailed = &ErrCode{40005, "获取用户失败", ""}
	AdminUserPasswdErr = &ErrCode{40006, "密码错误", ""}
)

var (
	ArticleCreateFailed = &ErrCode{40100, "创建文章失败", ""}
	ArticleDelFailed    = &ErrCode{40101, "删除文章失败", ""}
	ArticleUpdateFailed = &ErrCode{40102, "更新文章失败", ""}
	ArticleGetFailed    = &ErrCode{40103, "获取文章失败", ""}

	CommentDelFailed = &ErrCode{40104, "删除评论失败", ""}
	CommentGetFailed = &ErrCode{40105, "获取评论失败", ""}

	FriendLinkCreateFailed = &ErrCode{40106, "创建友情链接失败", ""}
	FriendLinkDelFailed    = &ErrCode{40107, "删除友情链接失败", ""}
	FriendLinkUpdateFailed = &ErrCode{40108, "更新友情链接失败", ""}
	FriendLinkGetFailed    = &ErrCode{40109, "获取友情链接失败", ""}
)
