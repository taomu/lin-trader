package lintypes

import "fmt"

// PLAT 定义交易平台枚举类型
type PLAT int

const (
	PLAT_BINANCE PLAT = iota // binance
	PLAT_OKX                 // okx
	PLAT_BYBIT               // bybit
	PLAT_BITGET              // bitget
	PLAT_UNKNOWN             // unknown
)

// String 返回PLAT的字符串表示
func (p PLAT) String() string {
	switch p {
	case PLAT_BINANCE:
		return "BINANCE"
	case PLAT_OKX:
		return "OKX"
	case PLAT_BYBIT:
		return "BYBIT"
	case PLAT_BITGET:
		return "BITGET"
	default:
		return "UNKNOWN"
	}
}

// ParsePLAT 从字符串解析PLAT
func ParsePLAT(s string) (PLAT, error) {
	switch s {
	case "BINANCE":
		return PLAT_BINANCE, nil
	case "OKX":
		return PLAT_OKX, nil
	case "BYBIT":
		return PLAT_BYBIT, nil
	case "BITGET":
		return PLAT_BITGET, nil
	default:
		return PLAT_UNKNOWN, fmt.Errorf("unknown platform: %s", s)
	}
}
