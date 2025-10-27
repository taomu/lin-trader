package types

type ApiInfo struct {
	Key        string
	Secret     string
	Passphrase string
}

const (
	SIDE_BUY             = "buy"
	SIDE_SELL            = "sell"
	POS_SIDE_LONG        = "long"
	POS_SIDE_SHORT       = "short"
	ORDER_TYPE_LIMIT     = "limit"     // 限价单
	ORDER_TYPE_MARKET    = "market"    // 市价单
	ORDER_TYPE_IOC       = "ioc"       // 立即成交剩余部分，否则取消
	ORDER_TYPE_FOK       = "fok"       // 全部成交否则取消
	ORDER_TYPE_POST_ONLY = "post_only" // 仅做maker单
)
