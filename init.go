package validators

const (
	ValidNotExist = "ValidNotExist"
	ValidError    = "ValidError"
	ValidIsNumber = "ValidIsNumber"
	ValidIsPhone  = "ValidIsPhone"
	ValidIsUrl    = "ValidIsUrl"
)

var Lang map[string]map[string]string

func init() {
	Lang = map[string]map[string]string{
		"zh": zh,
		"en": en,
	}
}
