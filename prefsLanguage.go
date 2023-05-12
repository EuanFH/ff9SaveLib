package ff9Save

import "github.com/jeandeaual/go-locale"

func LanguageStringToLanguageInt() map[string]int32 {
	return map[string]int32{
		//"System":      -1, not valid but game uses it internally might re add later
		"English(US)": 0,
		"English(UK)": 1,
		"Japanese":    2,
		"German":      3,
		"French":      4,
		"Italian":     5,
		"Spanish":     6,
	}
}
func LanguageIntToLanguageString() map[int32]string {
	return map[int32]string{
		//System language is only valid for binary saves converted to system language when saved as json
		-1: getSystemLanguage(),
		0:  "English(US)",
		1:  "English(UK)",
		2:  "Japanese",
		3:  "German",
		4:  "French",
		5:  "Italian",
		6:  "Spanish",
	}
}

func getSystemLanguage() string {
	defaultSystemLanguage := "English(US)"
	userLanguage, err := locale.GetLanguage()
	if err != nil {
		return defaultSystemLanguage
	}
	switch userLanguage {
	case "en":
		return defaultSystemLanguage
	case "ja":
		return "Japanese"
	case "de":
		return "German"
	case "fr":
		return "French"
	case "it":
		return "Italian"
	case "es":
		return "Spanish"
	}
	return defaultSystemLanguage
}

//assume system language if not found

// possibly enum this with all possible languages
type PrefsLanguage struct {
	Value string //find out type //int32 on marshal
}
