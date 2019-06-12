package lang

const (
	VALID_NOT_EXIST = "VALID_NOT_EXIST"
	VALID_ERROR     = "VALID_ERROR"
)

var Lang map[string]map[string]string

func init() {
	Lang = map[string]map[string]string{
		"zh": zh,
		"en": en,
	}
}
