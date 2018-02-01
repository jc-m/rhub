package ft991a


/*
Operates in two modes :
   - CAT : GET and SET are initiated as CAT commands, converted to RIG, RSP are RIG responses, converted in CAT
   - RIG : GET and SET are initiated as RIG commands, converted to CAT, RSP are CAT responses, converted in RIG
*/

import (
	"fmt"
	"github.com/jc-m/rhub/modules/radio/rigs"
	"regexp"
	"log"
	"strings"
)

const CMD_REGEXP  =  "^(?P<cmd>[a-zA-Z]{2})(?P<param>.*);"
const (
	FN_RSP = iota
	FN_GET
	FN_SET
)

var cmdParser = regexp.MustCompile(CMD_REGEXP)

var rspMap = map[string][]catFunc{
	"AC": {ACRsp}, // ANTENNA TUNER CONTROL
	"AG": {AGRsp}, // AUDIO GAIN
	"AI": {AIRsp}, // AUTO INFORMATION
	"IF": {IFRsp}, // INFORMATION
	"FA": {FARsp}, // VFOA
	"FB": {FBRsp}, // VFOB
	"FT": {FTRsp}, // SET TX VFO
	"KS": {KSRsp}, // KEY SPEED
	"MD": {MDRsp}, // MODE
	"NB": {NBRsp}, // NOISE BLANKER STATUS
	"NL": {NLRsp}, // NOISE BLANKER LEVEL
	"PA": {PARsp}, // PREAMP
	"PC": {PCRsp}, // Power
	"PS": {PSRsp}, // Power switch
	"RA": {RARsp}, // RX Attenuator
	"RG": {RGRsp}, // RF Gain
	"RM": {RMRsp}, // Meter
	"RL": {RLRsp}, // Noise Reduction
	"SM": {SMRsp}, // SMeter
	"VX": {VXRsp}, // VOX

}

type attributeMap map[string]string

type ft991a struct {

}

type catFunc func(param string) (*rigs.RigCommand, error)


func parse(r  string, param string ) (map[string]string, error) {

	parser := regexp.MustCompile(r)

	re  := parser.FindStringSubmatch(param)
	if re == nil {
		return nil, fmt.Errorf("Invalid parameter")
	}
	result := make(map[string]string)
	for i, name := range parser.SubexpNames() {
		if i != 0 { result[name] = re[i] }
	}
	return result, nil
}

func getSwitch(v string) string  {
	if v == "1" {
		return rigs.VAL_ON
	}
	if v == "0" {
		return rigs.VAL_OFF
	}
	return ""
}

func getModulation(v string) string {

	switch v{
	case "1":
		return rigs.MODUL_LSB
	case "2", "9", "C", "c":
		return rigs.MODUL_USB
	case "3":
		return rigs.MODUL_CW
	case "4", "A", "a", "B", "b":
		return rigs.MODUL_FM
	case "5", "D", "d":
		return rigs.MODUL_AM
	case "6":
		return rigs.MODUL_LSB
	case "7":
		return rigs.MODUL_CW_R
	case "8":
		return rigs.MODUL_LSB
	case "E","e":
		return rigs.MODUL_C4FM
	}
	return ""
}

func getMode(v string) string {

	switch v{
	case "6", "9":
		return rigs.MODE_RTTY
	case "8", "A", "a", "C", "c":
		return rigs.MODE_DATA
	}
	return rigs.MODE_VOICE
}

func getFilter(v string) string {

	switch v{
	case "B","b","D","d":
		return rigs.FILT_N
	}
	return rigs.FILT_W
}

func singleValRsp(re string, param string, key string, trim string) (*rigs.RigCommand, error) {

	m, err := parse(re, param)
	if err != nil {
		return nil, err
	}

	v, ok := m["P1"]
	if ! ok {
		return nil, fmt.Errorf("Missing P1")
	}

	cmdParams := make(map[string]string)

	if len(trim) > 0 {
		cmdParams[key] = strings.TrimLeft(v,trim)
	} else {
		cmdParams[key] = v
	}

	return &rigs.RigCommand{
		Id: rigs.RIG_RSP,
		Params : cmdParams,
	}, nil

}

func switchValRsp(re string, param string, key string) (*rigs.RigCommand, error) {

	m, err := parse(re, param)
	if err != nil {
		return nil, err
	}

	v, ok := m["P1"]
	if ! ok {
		return nil, fmt.Errorf("Missing P1")
	}

	var cmdParams map[string]string

	cmdParams[key] = getSwitch(v)

	return &rigs.RigCommand{
		Id: rigs.RIG_RSP,
		Params: cmdParams ,
	}, nil
}

func ACRsp(param string) (*rigs.RigCommand, error) {
	// Valid values for response is 0 or 1 (2 is to start tuning)
	const RE="^00(?P<P1>0|1)"

	return switchValRsp(RE, param, rigs.PRM_TUNER)

}

func AGRsp(param string) (*rigs.RigCommand, error) {
	// Valid values are 0 to 255
	const RE="^0(?P<P1>[01][0-9][0-9]|2[0-4][0-9]|25[0-5])"

	return singleValRsp(RE, param, rigs.PRM_AUDIOG, "0")
}

func AIRsp(param string) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>0|1)"

	return switchValRsp(RE, param, rigs.PRM_AUTOINFO)
}

func FARsp(param string) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>[0-9]{9})"

	return singleValRsp(RE, param, rigs.PRM_VFOA, "0")
}

func FBRsp(param string) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>[0-9]{9})"

	return singleValRsp(RE, param, rigs.PRM_VFOB, "0")
}

func FTRsp(param string) (*rigs.RigCommand, error) {
	const RE = "^(?P<P1>0|1)"

	var cmdParams map[string]string

	m, err := parse(RE, param)
	if err != nil {
		return nil, err
	}

	switch m["P1"] {
	case "0":
		cmdParams[rigs.PRM_TXVFO] = rigs.VAL_VFOA
		// IF VFOA is TX, SPLIT is OFF
		cmdParams[rigs.PRM_SPLIT] = rigs.VAL_OFF
	case "1":
		cmdParams[rigs.PRM_TXVFO] = rigs.VAL_VFOB
		// IF VFOB is TX, SPLIT is ON
		cmdParams[rigs.PRM_SPLIT] = rigs.VAL_ON
	}

	return &rigs.RigCommand{
		Id:     rigs.RIG_RSP,
		Params: cmdParams,
	}, nil
}


func IFRsp(param string) (*rigs.RigCommand, error) {
	const RE="(?P<P1>[0-9]{3})(?P<P2>[0-9]{9})"+
		"(?P<P3>[+|-][0-9]{4})(?P<P4>[0-9a-eA-E]{1})"+
		"(?P<P5>[0-9]{1})(?P<P6>[0-9]{1})"+
		"(?P<P7>[0-9]{1})(?P<P8>[0-9]{1})"+
		"(?P<P9>[0-9]{2})(?P<P10>[0-9]{1})"

	m, err := parse(RE, param)
	if err != nil {
		return nil, err
	}
	cmdParams := make(map[string]string)
	// just support
	// 		P2 : VFO A Freq
	//		P3 : Clarifier Direction + offset
	//		P6 : Mode
	cmdParams[rigs.PRM_VFOA] = strings.TrimLeft(m["P2"],"0")
	cmdParams[rigs.PRM_CLAR] = m["P3"]
	cmdParams[rigs.PRM_MODUL] = getModulation(m["P4"])
	cmdParams[rigs.PRM_MODE] = getMode(m["P4"])
	cmdParams[rigs.PRM_FILT] = getFilter(m["P4"])

	return &rigs.RigCommand{
		Id: rigs.RIG_RSP,
		Params : cmdParams,
	}, nil
}

func KSRsp(param string) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>[0-9]{3})"

	return singleValRsp(RE, param, rigs.PRM_KEYSPEED, "0")
}

func MDRsp(param string) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>0)(?P<P2>[0-9a-eA-E]{1})"

	m, err := parse(RE, param)
	if err != nil {
		return nil, err
	}
	cmdParams := make(map[string]string)

	cmdParams[rigs.PRM_MODUL] = getModulation(m["P2"])
	cmdParams[rigs.PRM_MODE] = getMode(m["P2"])
	cmdParams[rigs.PRM_FILT] = getFilter(m["P2"])

	return &rigs.RigCommand{
		Id: rigs.RIG_RSP,
		Params: cmdParams ,
	}, nil

}

//NOISE BLANKER LEVEL
func NLRsp(param string) (*rigs.RigCommand, error) {
	// Valid values are 0 to 10
	const RE="^00(?P<P1>0[0-9]|10)"

	return singleValRsp(RE, param, rigs.PRM_NBLEVEL, "0")
}


//NOISE BLANKER STATUS
func NBRsp(param string) (*rigs.RigCommand, error) {
	const RE="^0(?P<P1>0|1)"

	return switchValRsp(RE, param, rigs.PRM_NBSTATUS)
}

//PREAMP
func PARsp(param string) (*rigs.RigCommand, error) {
	const RE="^0(?P<P1>[0-2])"

	return singleValRsp(RE, param, rigs.PRM_PREAMP, "")
}

// TX POWER
func PCRsp(param string) (*rigs.RigCommand, error) {
	// Valid values are 5 to 100
	const RE="^(?P<P1>00[5-9]|0[1-9][0-9]|100)"

	return singleValRsp(RE, param, rigs.PRM_POWER, "0")
}

// POWER SWITCH
func PSRsp(param string) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>0|1)"

	return switchValRsp(RE, param, rigs.PRM_PWRSWITCH)
}

// RF ATTENUATOR
func RARsp(param string) (*rigs.RigCommand, error) {
	const RE="^0(?P<P1>0|1)"

	return switchValRsp(RE, param, rigs.PRM_ATT)
}

// RF GAIN
func RGRsp(param string) (*rigs.RigCommand, error) {
	// Valid values are 0 to 255
	const RE="^0(?P<P1>[01][0-9][0-9]|2[0-4][0-9]|25[0-5])"

	return singleValRsp(RE, param, rigs.PRM_RFGAIN, "0")
}

// READ METERS
func RMRsp(param string) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>[0-9]{1})(?P<P2>.*)"
	const RE2="^(?P<P1>[0-9]{3})"

	m, err := parse(RE, param)
	if err != nil {
		return nil, err
	}
	log.Printf("%+v",m)

	var key = ""
	switch m["P1"] {
	case "1": // SMeter
		key = rigs.PRM_SMETER
	case "3": // Comp
		key = rigs.PRM_COMP
	case "4": // ALC
		key = rigs.PRM_ALC
	case "5": // POWER
		key = rigs.PRM_POWER
	case "6": // SWR
		key = rigs.PRM_SWR
	case "7": // ID
		key = rigs.PRM_ID
	case "8": // VDD
		key = rigs.PRM_VDD
	}
	return singleValRsp(RE2, m["P2"], key, "0")
}

// NOISE REDUCTION LEVEL
func RLRsp(param string) (*rigs.RigCommand, error) {
	// Valid values are 0 to 15
	const RE="^0(?P<P1>0[0-9]|1[0-5])"

	return singleValRsp(RE, param, rigs.PRM_NRLVEL, "0")
}

//SMETER
func SMRsp(param string) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>[0-9]{4})"

	return singleValRsp(RE, param, rigs.PRM_SMETER, "0")
}

//VOX
func VXRsp(param string) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>0|1)"

	return switchValRsp(RE, param, rigs.PRM_VOX)

}

func (r *ft991a) Open()  error {
	return fmt.Errorf("Not implemented")

}

// Converts a CAT response/autoinfo into a RIG command

func  (r *ft991a) OnCatUpStream(command string) (*rigs.RigCommand, error) {

	re  := cmdParser.FindStringSubmatch(command)
	log.Printf("%d", len(re))

	if len(re) < 2 {
		return nil, fmt.Errorf("Not a valid command")
	}
	cmdId := re[1]

	if funcs, ok := rspMap[cmdId]; ok {
		if  funcs[FN_RSP] != nil {
			return funcs[FN_RSP](re[2])
		}
	}

	return nil, fmt.Errorf("Not a known command")
}

func New() rigs.Rig {

	return &ft991a {}
}
