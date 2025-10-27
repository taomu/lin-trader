package lintypes

const (
	SIDE_BUY             = "buy"       // 买入
	SIDE_SELL            = "sell"      // 卖出
	POS_SIDE_LONG        = "long"      // 多单持仓方式
	POS_SIDE_SHORT       = "short"     // 空单持仓方式
	POS_SIDE_BOTH        = "both"      // 单向持仓方式
	ORDER_TYPE_LIMIT     = "limit"     // 限价单
	ORDER_TYPE_MARKET    = "market"    // 市价单
	ORDER_TYPE_IOC       = "ioc"       // 立即成交剩余部分，否则取消
	ORDER_TYPE_FOK       = "fok"       // 全部成交否则取消
	ORDER_TYPE_POST_ONLY = "post_only" // 仅做maker单
)
