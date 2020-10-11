package FF9Save

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
		//-1: "System", not valid but game uses it internally might re add later
		0:  "English(US)",
		1:  "English(UK)",
		2:  "Japanese",
		3:  "German",
		4:  "French",
		5:  "Italian",
		6:  "Spanish",
	}
}

//assume system language if not found

//possibly enum this with all possible languages
type PrefsLangauge struct {
	Value string //find out type //int32 on marshal
}
