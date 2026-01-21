package i18n

type Locale string

const (
	EN Locale = "en"
	ZH Locale = "zh"
	JA Locale = "ja"
)

type Key string

func (k Key) String() string { return string(k) }
