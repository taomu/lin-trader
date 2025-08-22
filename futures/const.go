package futures

import "fmt"

// PLAT 定义交易平台枚举类型
type PLAT int

const (
	PLAT_BINANCE PLAT = iota // binance
	PLAT_OKX                 // okx
	PLAT_BYBIT               // bybit
	PLAT_BITGET              // bitget
)

// String 返回PLAT的字符串表示
func (p PLAT) String() string {
	switch p {
	case PLAT_BINANCE:
		return "binance"
	case PLAT_OKX:
		return "okx"
	case PLAT_BYBIT:
		return "bybit"
	case PLAT_BITGET:
		return "bitget"
	default:
		return "unknown"
	}
}

// ParsePLAT 从字符串解析PLAT
func ParsePLAT(s string) (PLAT, error) {
	switch s {
	case "binance":
		return PLAT_BINANCE, nil
	case "okx":
		return PLAT_OKX, nil
	case "bybit":
		return PLAT_BYBIT, nil
	case "bitget":
		return PLAT_BITGET, nil
	default:
		return PLAT_BINANCE, fmt.Errorf("unknown platform: %s", s)
	}
}

type ApiInfo struct {
	Key        string
	Secret     string
	Passphrase string
}
