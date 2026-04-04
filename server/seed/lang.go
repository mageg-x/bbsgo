package seed

import (
	"bbsgo/handlers"
	"flag"
	"log"
)

var langFlag = flag.String("lang", "zh", "initialization language: zh, en")
var initializedLang string

// GetLang 获取当前语言设置
func GetLang() string {
	return initializedLang
}

// GetLangPtr 返回语言标志指针
func GetLangPtr() *string {
	return langFlag
}

// Init 初始化
func Init() {
	flag.Parse()
	initializedLang = *langFlag
	log.Printf("[seed] using language: %s", initializedLang)

	switch initializedLang {
	case "en":
		seedDataEn()
	default:
		seedDataZh()
	}

	// 初始化勋章（传入当前语言）
	handlers.SeedBadgesWithLang(initializedLang)
}
