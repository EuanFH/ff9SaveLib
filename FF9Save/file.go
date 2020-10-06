package FF9Save

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type File struct {
	Data struct {
		State_10000       State_10000       `json:"10000_State"`
		Event_20000       Event_20000       `json:"20000_Event"`
		MiniGame_30000    MiniGame_30000    `json:"30000_MiniGame"`
		Common_40000      Common_40000      `json:"40000_Common"`
		Setting_50000     Setting_50000     `json:"50000_Setting"`
		Sound_60000       Sound_60000       `json:"60000_Sound"`
		World_70000       World_70000       `json:"70000_World"`
		Achievement_80000 Achievement_80000 `json:"80000_Achievement"`
		State_91000       [128]int32FromStr `json:"91000_State"`
		Event_92000       [128]int32FromStr `json:"92000_Event"`
		MiniGame_93000    [128]int32FromStr `json:"93000_MiniGame"`
		Common_94000      Common_94000      `json:"94000_Common"`
		Setting_95000     Setting_95000     `json:"95000_Setting"`
		Sound_96000       [128]int32FromStr `json:"96000_Sound"`
		World_97000       [128]int32FromStr `json:"97000_World"`
		Achievement_98000 Achievement_98000 `json:"98000_Achievement"`
		Other_99000       [384]int32FromStr `json:"99000_Other"`
	}
}


func (f *File)MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	//handle none
	if err := binary.Write(buf, binary.LittleEndian, []byte{'S','A','V','E'}); err != nil{
		return nil, err
	}
	if err := ff9SaveToBinary(*f, buf, 0); err != nil{
		return nil, err
	}
	return buf.Bytes(), nil
}

func (f *File)UnmarshalBinary(data []byte) error{
	buf := bytes.NewBuffer(data[4:])
	//handle none
	if err := ff9SaveBinaryToStruct(f, buf, 0); err != nil{
		return err
	}
	return nil
}

func ff9SaveBinaryToStruct(save interface{}, buf *bytes.Buffer, depth int) error{
	//field := reflect.TypeOf(save)
	value := reflect.Indirect(reflect.ValueOf(save))
	switch reflect.Indirect(reflect.ValueOf(save)).Kind() {
	case reflect.Struct:
		var fieldSortedIndex []int
		if(depth > 1){
			var fieldNames []string
			for i := 0; i < value.NumField(); i++ {
				jsonName := value.Type().Field(i).Tag.Get("json")
				jsonName = strings.Split(jsonName, ",")[0]
				fieldNames = append(fieldNames, jsonName)
			}
			sortedFieldNames := make([]string, len(fieldNames))
			copy(sortedFieldNames, fieldNames)
			//Reason for the weird sorting here is because c sharp and go sort very
			//differently even though they use the same algorithm
			//Go takes into account case c sharp dosn't
			//not going to attempt to explain the numbers shit
			sort.Slice(sortedFieldNames, func(i int, j int) bool{
				//check if start with number
				_, containsNumberI := strconv.Atoi(string(sortedFieldNames[i][0]))
				_, containsNumberJ := strconv.Atoi(string(sortedFieldNames[i][0]))
				if containsNumberI == nil && containsNumberJ == nil {
					numberI, _ := strconv.Atoi(strings.Split(sortedFieldNames[i], "_")[0])
					numberj, _ := strconv.Atoi(strings.Split(sortedFieldNames[j], "_")[0])
					return numberI < numberj
				}
				//lowercase to avoid uppercase letters messing with order
				return strings.ToLower(sortedFieldNames[i]) < strings.ToLower(sortedFieldNames[j])
			})
			for _, sortedFieldName := range sortedFieldNames{
				for j, fieldName := range fieldNames {
					if(sortedFieldName == fieldName) {
						fieldSortedIndex = append(fieldSortedIndex, j)
					}
				}
			}
		}
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			if(fieldSortedIndex != nil){
				field = value.Field(fieldSortedIndex[i])
			}
			if err := ff9SaveBinaryToStruct(field.Addr().Interface(), buf, depth + 1); err != nil{
				return err
			}
		}
	case reflect.Array:
		for i := 0; i < value.Len(); i++ {
			arrayValue := value.Index(i)
			if err := ff9SaveBinaryToStruct(arrayValue.Addr().Interface(), buf, depth + 1); err != nil{
				return err
			}
		}
	default:
		variableBytes := make([]byte, reflect.Indirect(reflect.ValueOf(save)).Type().Size())
		if _, err := buf.Read(variableBytes); err != nil{
			return err
		}
		if err := binary.Read(bytes.NewBuffer(variableBytes), binary.LittleEndian, save); err != nil {
			return err
		}
	}
	return nil
}

//depth is needed so we dont sort the first or second layer of structs
//data and the order of the main structs should not be sorted
func ff9SaveToBinary(save interface{}, buf *bytes.Buffer, depth int) error {
	field := reflect.TypeOf(save)
	value := reflect.ValueOf(save)
	switch reflect.TypeOf(save).Kind() {
	case reflect.Struct:
		var fieldSortedIndex []int
		if(depth > 1){
			var fieldNames []string
			for i := 0; i < field.NumField(); i++ {
				jsonName := value.Type().Field(i).Tag.Get("json")
				jsonName = strings.Split(jsonName, ",")[0]
				fieldNames = append(fieldNames, jsonName)
			}
			sortedFieldNames := make([]string, len(fieldNames))
			copy(sortedFieldNames, fieldNames)
			//Reason for the weird sorting here is because c sharp and go sort very
			//differently even though they use the same algorithm
			//Go takes into account case c sharp dosn't
			//not going to attempt to explain the numbers shit
			sort.Slice(sortedFieldNames, func(i int, j int) bool{
				//check if start with number
				_, containsNumberI := strconv.Atoi(string(sortedFieldNames[i][0]))
				_, containsNumberJ := strconv.Atoi(string(sortedFieldNames[i][0]))
				if containsNumberI == nil && containsNumberJ == nil {
					numberI, _ := strconv.Atoi(strings.Split(sortedFieldNames[i], "_")[0])
					numberj, _ := strconv.Atoi(strings.Split(sortedFieldNames[j], "_")[0])
					return numberI < numberj
				}
				//lowercase to avoid uppercase letters messing with order
				return strings.ToLower(sortedFieldNames[i]) < strings.ToLower(sortedFieldNames[j])
			})
			for _, sortedFieldName := range sortedFieldNames{
				for j, fieldName := range fieldNames {
					if(sortedFieldName == fieldName) {
						fieldSortedIndex = append(fieldSortedIndex, j)
					}
				}
			}
		}
		for i := 0; i < field.NumField(); i++ {
			field := value.Field(i)
			if(fieldSortedIndex != nil){
				field = value.Field(fieldSortedIndex[i])
			}
			if err := ff9SaveToBinary(field.Interface(), buf, depth + 1); err != nil{
				return err
			}
		}
	case reflect.Array:
		for i := 0; i < value.Len(); i++ {
			arrayValue := value.Index(i)
			if err := ff9SaveToBinary(arrayValue.Interface(), buf, depth + 1); err != nil{
				return err
			}
		}
	default:
		if err := binary.Write(buf, binary.LittleEndian, reflect.ValueOf(save).Interface()); err != nil {
			return err
		}
	}
	return nil
}

func(f *File) UnmarshalJSON(data []byte) error{
	//aliasing type to remove unmarshal function to stop infinite loop
	type Alias File
	var alias Alias
	//fixing boolean values to convert correctly
	data = bytes.ReplaceAll(data, []byte("True"), []byte("true"))
	data = bytes.ReplaceAll(data, []byte("False"), []byte("false"))
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*f = File(alias)
	return nil
}

type State_10000 struct {
	Mode         byte        `json:"mode,string"`
	PrevMode     byte        `json:"prevMode,string"`
	FldMapNo     int16       `json:"fldMapNo,string"`
	FldLocNo     int16       `json:"fldLocNo,string"`
	BtlMapNo     int16       `json:"btlMapNo,string"`
	BtlSubMapNo  int8        `json:"btlSubMapNo,string"` //sbyte
	WldMapNo     int16       `json:"wldMapNo,string"`
	WldLocNo     int16       `json:"wldLocNo,string"`
	TimeCounter  float32     `json:"timeCounter,string"` //math.float32bits might be needed don't think its encoding with the right standard
	TimerControl boolFromStr `json:"timerControl"`
	TimerDisplay boolFromStr `json:"timerDisplay"`
}

type Event_20000 struct {
	GStepCount   int32    `json:"gStepCount,string"`
	GEventGlobal String4K `json:"gEventGlobal"` //size of large string //this could also be rune
}

type MiniGame_30000 struct {
	SWin int16                     `json:"sWin,string"`
	SLose int16                    `json:"sLose,string"`
	SDraw int16                    `json:"sDraw,string"`
	MiniGameCard [100]MiniGameCard `json:"MiniGameCard,string"`
}

type MiniGameCard struct {
	ID byte `json:"id,string"`
	Side byte `json:"side,string"`
	Atk byte `json:"atk,string"`
	Type int32 `json:"type,string"`
	Pdef byte `json:"pdef,string"`
	Mdef byte `json:"mdef,string"`
	Cpoint byte `json:"cpoint,string"`
	Arrow byte `json:"arrow,string"`
}

type Common_40000 struct {
	Player [9]Player          `json:"players,string"`
	Slot [4]byteFromStr       `json:"slot,string"`
	Escape_no uint16          `json:"escape_no,string"`
	Summon_flag uint16        `json:"summon_flag,string"`
	Gil uint32                `json:"gil,string"`
	Frog_no int16             `json:"frog_no,string"`
	Steal_no int16            `json:"steal_no,string"`
	Dragon_no int16           `json:"dragon_no,string"`
	Items [256]Item           `json:"items,string"`
	Rareitems [64]byteFromStr `json:"rareItems,string"`
}

type Player struct {
	Name     String128 `json:"name,string"` //assuming 128 could be wrong
	Category byte      `json:"category,string"`
	Level    byte      `json:"level,string"`
	Exp      uint32    `json:"exp,string"`
	Cur      struct {
		Hp uint16 `json:"hp,string"`
		Mp int16 `json:"mp,string"`
		At int16 `json:"at,string"`
		At_coef int8 `json:"at_coef,string"` //sbyte
		Capa byte `json:"capa,string"`
	} `json:"cur"`
	Max struct { //max and cur are the same type
		Hp uint16 `json:"hp,string"`
		Mp int16 `json:"mp,string"`
		At int16 `json:"at,string"`
		At_coef int8 `json:"at_coef,string"` //sbyte
		Capa byte `json:"capa,string"`
	} `json:"max"`
	Trance byte `json:"trance,string"`
	Web_bone byte `json:"web_bone,string"` //should this always be 0?
	Elem struct {
		Dex byte `json:"dex,string"`
		Str byte `json:"str,string"`
		Mgc byte `json:"mgc,string"`
		Wpr byte `json:"wpr,string"`
	} `json:"elem"`
	Defence struct {
		P_def byte `json:"p_def,string"`
		P_ev byte `json:"p_ev,string"`
		M_def byte `json:"m_def,string"`
		M_ev byte `json:"m_ev,string"`
	} `json:"defence"`
	Basis struct {
		Max_hp int16 `json:"max_hp,string"`
		Max_mp int16 `json:"max_mp,string"`
		Dex byte `json:"dex,string"`
		Str byte `json:"str,string"`
		Mgc byte `json:"mgc,string"`
		Wpr byte `json:"wpr,string"`
	} `json:"basis"`
	Info struct {
		Slot_no byte `json:"slot_no,string"`
		Serial_no byte `json:"serial_no,string"`
		Row byte `json:"row,string"`
		Win_pose byte `json:"win_pose,string"`
		Party byte `json:"party,string"`
		Menu_type byte `json:"menu_type,string"`
	} `json:"info"`
	Status byte          `json:"status,string"`
	Equip [5]byteFromStr `json:"equip,string"`
	Bonus struct {
		Dex uint16 `json:"dex,string"`
		Str uint16 `json:"str,string"`
		Mgc uint16 `json:"mgc,string"`
		Wpr uint16 `json:"wpr,string"`
	} `json:"bonus"`
	Pa [48]byteFromStr  `json:"pa,string"`
	Sa [2]uint32FromStr `json:"sa,string"`
}

type Item struct {
	ID byte `json:"id,string"`
	Count byte `json:"count,string"`
}

type Setting_50000 struct {
	Cfg struct{
		Sound uint64 `json:"sound,string"`
		Sound_effect uint64 `json:"sound_effect,string"`
		Control uint64 `json:"control,string"`
		Cursor uint64 `json:"cursor,string"`
		Atb uint64                             `json:"atb,string"`
		Camera uint64                          `json:"camera,string"`
		Move uint64                            `json:"move,string"`
		Vibe uint64                            `json:"vibe,string"`
		Btl_speed uint64                       `json:"btl_speed,string"`
		Fld_msg uint64                         `json:"fld_msg,string"`
		Here_icon uint64                       `json:"here_icon,string"`
		Win_type uint64                        `json:"win_type,string"`
		Target_win uint64                      `json:"target_win,string"`
		Control_data uint64                    `json:"control_data,string"`
		Control_data_keyboard [10]int32FromStr `json:"control_data_keyboard,string"` //most likely int32 but not sure
		Control_data_joystick [10]String128    `json:"control_data_joystick,string"` //assuming 128
		Skip_btl_camera uint64                 `json:"skip_btl_camera,string"`
	} `json:"cfg"`
	Time float32 `json:"time,string"`
}

type Sound_60000 struct {
	Auto_save_bgm_id int32 `json:"auto_save_bgm_id,string"`
}

type World_70000 struct {
	DataCameraStateRotationMax       float32     `json:"data.cameraState.rotationMax,string"`
	DataCameraStateUpperCounter      int16       `json:"data.cameraState.upperCounter,string"`
	DataCameraStateUpperCounterSpeed int32       `json:"data.cameraState.upperCounterSpeed,string"`
	DataCameraStateUpperCounterForce boolFromStr `json:"data.cameraState.upperCounterForce"`
	DataCameraStateRotation          float32     `json:"data.cameraState.rotation,string"`
	DataCameraStateRotationRev       float32     `json:"data.cameraState.rotationRev,string"`
	DataHintmap                      uint32      `json:"data.hintmap,string"`
}

type Achievement_80000 struct {
	AteCheckArray [100]int32FromStr   `json:"AteCheckArray,string"`
	EvtReservedArray [17]int32FromStr `json:"EvtReservedArray,string"`
	BlkMag_no int32                   `json:"blkMag_no,string"`
	WhtMag_no int32                   `json:"whtMag_no,string"`
	BluMag_no int32                   `json:"bluMag_no,string"`
	Summon_no int32                   `json:"summon_no,string"`
	Enemy_no int32                    `json:"enemy_no,string"`
	BackAtk_no int32                  `json:"backAtk_no,string"`
	Defence_no int32                  `json:"defence_no,string"`
	Trance_no int32                   `json:"trance_no,string"`
	Abilities [221]int32FromStr       `json:"abilities,string"`
	PassiveAbilities [63]int32FromStr `json:"passiveAbilities,string"`
	SynthesisCount int32              `json:"synthesisCount,string"`
	AuctionTime int32                 `json:"AuctionTime,string"`
	StiltzkinBuy int32                `json:"StiltzkinBuy,string"`
	QuadmistWinList [300]int32FromStr `json:"QuadmistWinList,string"`
}

type Setting_95000 struct {
	Time_00001 float64                     `json:"00001_time,string"`
	ReservedBuffer_99999 [126]int32FromStr `json:"99999_ReservedBuffer,string"`
}

type Achievement_98000 struct {
	Abnormal_status_00001            uint32            `json:"00001_abnormal_status,string"`
	Summon_shiva_00002               boolFromStr       `json:"00002_summon_shiva"`
	Summon_ifrit_00003               boolFromStr       `json:"00003_summon_ifrit"`
	Summon_ramuh_00004               boolFromStr       `json:"00004_summon_ramuh"`
	Summon_carbuncle_reflector_00005 boolFromStr       `json:"00005_summon_carbuncle_reflector"`
	Summon_carbuncle_haste_00006     boolFromStr       `json:"00006_summon_carbuncle_haste"`
	Summon_carbuncle_protect_00007   boolFromStr       `json:"00007_summon_carbuncle_protect"`
	Summon_carbuncle_shell_00008     boolFromStr       `json:"00008_summon_carbuncle_shell"`
	Summon_fenrir_earth_00009        boolFromStr       `json:"00009_summon_fenrir_earth"`
	Summon_fenrir_wind_000010        boolFromStr       `json:"000010_summon_fenrir_wind"`
	Summon_atomos_000011             boolFromStr       `json:"000011_summon_atomos"`
	Summon_phoenix_000012            boolFromStr       `json:"000012_summon_phoenix"`
	Summon_leviathan_000013          boolFromStr       `json:"000013_summon_leviathan"`
	Summon_odin_000014               boolFromStr       `json:"000014_summon_odin"`
	Summon_madeen_000015             boolFromStr       `json:"000015_summon_madeen"`
	Summon_bahamut_000016            boolFromStr       `json:"000016_summon_bahamut"`
	Summon_arc_000017                boolFromStr       `json:"000017_summon_arc"`
	ReservedBuffer_99999             [123]int32FromStr `json:"99999_ReservedBuffer,string"`
}

type Common_94000 struct {
	Player_bonus_00001 [9]uint32FromStr    `json:"00001_player_bonus,string"`
	ReservedBuffer_99999 [119]int32FromStr `json:"99999_ReservedBuffer,string"`
}
