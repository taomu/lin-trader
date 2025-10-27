package lintypes

const (
	// 订单方向
	SIDE_BUY  = "BUY"  // 买入
	SIDE_SELL = "SELL" // 卖出
	// 持仓方向
	POS_SIDE_LONG  = "LONG"  // 多单持仓方式
	POS_SIDE_SHORT = "SHORT" // 空单持仓方式
	POS_SIDE_BOTH  = "BOTH"  // 单向持仓方式
	// 订单类型
	ORDER_TYPE_LIMIT       = "LIMIT"       // 限价单
	ORDER_TYPE_MARKET      = "MARKET"      // 市价单
	ORDER_TYPE_STOP        = "STOP"        // 止损单
	ORDER_TYPE_TAKE_PROFIT = "TAKE_PROFIT" // 止盈单
	ORDER_TYPE_LIQUIDATION = "LIQUIDATION" // 强平单
	// 订单状态
	ORDER_STATUS_NEW              = "NEW"              // 已接受订单
	ORDER_STATUS_FILLED           = "FILLED"           // 已成交订单
	ORDER_STATUS_CANCELED         = "CANCELED"         // 已取消订单
	ORDER_STATUS_PARTIALLY_FILLED = "PARTIALLY_FILLED" // 部分成交订单
	ORDER_STATUS_EXPIRED          = "EXPIRED"          // 已过期订单
	ORDER_STATUS_EXPIRED_IN_MATCH = "EXPIRED_IN_MATCH" // 已过期且未成交订单
	// 订单事件
	ORDER_EVENT_NEW        = "NEW"        // 新订单事件
	ORDER_EVENT_CANCELED   = "CANCELED"   // 已撤订单事件
	ORDER_EVENT_CALCULATED = "CALCULATED" // 订单 ADL 或爆仓事件
	ORDER_EVENT_EXPIRED    = "EXPIRED"    // 订单失效事件
	ORDER_EVENT_TRADE      = "TRADE"      // 交易事件
	// 有效方式
	ORDER_TIME_IN_FORCE_GTC = "GTC" // 成交为止有效
	ORDER_TIME_IN_FORCE_FOK = "FOK" // 全部成交有效
	ORDER_TIME_IN_FORCE_IOC = "IOC" // 立即成交剩余有效
	ORDER_TIME_IN_FORCE_GTX = "GTX" // 仅做挂单 post-only
)
