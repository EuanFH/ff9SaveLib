package FF9Save

type SlotPreview struct {
	IsPreviewCorrupted bool
	HasData bool
	Gil int //find out type
	PlayDuration int //find out type
	WinType int //find out type `json:"win_type"`
	Location string //find out type
	CharacterInfoList [4]CharacterInfo
	Timestamp float64 //find out type
	ReservedData [64]int //find out type
}

type CharacterInfo struct {
	SerialID int //find out type
	Level int //find out type
	Name string //find out type
}

