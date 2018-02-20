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
	FN_REQ
)

var cmdParser = regexp.MustCompile(CMD_REGEXP)

var funcMap = map[string]catFunc{
	"AC": AC, // ANTENNA TUNER CONTROL
	"AG": AG, // AUDIO GAIN
	"AI": AI,    // AUTO INFORMATION
	"IF": IF, // INFORMATION
	"FA": FA, // VFOA
	"FB": FB, // VFOB
	"FT": FT, // SET TX VFO
	"KS": KS, // KEY SPEED
	"MD": MD, // MODE
	"NB": NB, // NOISE BLANKER STATUS
	"NL": NL, // NOISE BLANKER LEVEL
	"PA": PA, // PREAMP
	"PC": PC, // Power
	"PS": PS, // Power switch
	"RA": RA, // RX Attenuator
	"RG": RG, // RF Gain
	"RM": RM, // Meter
	"RL": RL, // Noise Reduction
	"SM": SM, // SMeter
	"VX": VX, // VOX

}

type attributeMap map[string]string

type ft991a struct {

}

type catFunc func(param string, dir int) (*rigs.RigCommand, error)


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
	op := rigs.RIG_RSP
	return switchVal(re, param, key, op)
}

func switchValSet(re string, param string, key string) (*rigs.RigCommand, error) {
	op := rigs.RIG_SET
	return switchVal(re, param, key, op)
}

func switchVal(re string, param string, key string, op string) (*rigs.RigCommand, error) {

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
		Id: op,
		Params: cmdParams ,
	}, nil
}

func AC(param string, dir int) (*rigs.RigCommand, error) {
	// Valid values for response is 0 or 1 (2 is to start tuning)
	const RE="^00(?P<P1>0|1)"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return switchValRsp(RE, param, rigs.PRM_TUNER)

}

func AG(param string, dir int) (*rigs.RigCommand, error) {
	// Valid values are 0 to 255
	const RE="^0(?P<P1>[01][0-9][0-9]|2[0-4][0-9]|25[0-5])"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return singleValRsp(RE, param, rigs.PRM_AUDIOG, "0")
}

func AI(param string, dir int) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>0|1)"

	// rsp
	if dir == rigs.CAT_DIR_UP {
		return switchValRsp(RE, param, rigs.PRM_AUTOINFO)

	}

	cmdParams := make(map[string]string)

	// get
	if param == "" {
		cmdParams[rigs.PRM_AUTOINFO] = ""
		return &rigs.RigCommand{
			Id: rigs.RIG_GET,
			Params : cmdParams,
		}, nil
	}

	// set
	return switchValSet(RE, param, rigs.PRM_AUTOINFO)

}

func FA(param string, dir int) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>[0-9]{9})"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return singleValRsp(RE, param, rigs.PRM_VFOA, "0")
}

func FB(param string, dir int) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>[0-9]{9})"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return singleValRsp(RE, param, rigs.PRM_VFOB, "0")
}

func FT(param string, dir int) (*rigs.RigCommand, error) {
	const RE = "^(?P<P1>0|1)"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	cmdParams := make(map[string]string)

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


func IF(param string, dir int) (*rigs.RigCommand, error) {
	const RE="(?P<P1>[0-9]{3})(?P<P2>[0-9]{9})"+
		"(?P<P3>[+|-][0-9]{4})(?P<P4>[0-9a-eA-E]{1})"+
		"(?P<P5>[0-9]{1})(?P<P6>[0-9]{1})"+
		"(?P<P7>[0-9]{1})(?P<P8>[0-9]{1})"+
		"(?P<P9>[0-9]{2})(?P<P10>[0-9]{1})"

	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
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

func KS(param string, dir int) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>[0-9]{3})"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return singleValRsp(RE, param, rigs.PRM_KEYSPEED, "0")
}

func MD(param string, dir int) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>0)(?P<P2>[0-9a-eA-E]{1})"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
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
func NL(param string, dir int) (*rigs.RigCommand, error) {
	// Valid values are 0 to 10
	const RE="^00(?P<P1>0[0-9]|10)"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return singleValRsp(RE, param, rigs.PRM_NBLEVEL, "0")
}


//NOISE BLANKER STATUS
func NB(param string, dir int) (*rigs.RigCommand, error) {
	const RE="^0(?P<P1>0|1)"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return switchValRsp(RE, param, rigs.PRM_NBSTATUS)
}

//PREAMP
func PA(param string, dir int) (*rigs.RigCommand, error) {
	const RE="^0(?P<P1>[0-2])"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return singleValRsp(RE, param, rigs.PRM_PREAMP, "")
}

// TX POWER
func PC(param string, dir int) (*rigs.RigCommand, error) {
	// Valid values are 5 to 100
	const RE="^(?P<P1>00[5-9]|0[1-9][0-9]|100)"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return singleValRsp(RE, param, rigs.PRM_POWER, "0")
}

// POWER SWITCH
func PS(param string, dir int) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>0|1)"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return switchValRsp(RE, param, rigs.PRM_PWRSWITCH)
}

// RF ATTENUATOR
func RA(param string, dir int) (*rigs.RigCommand, error) {
	const RE="^0(?P<P1>0|1)"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return switchValRsp(RE, param, rigs.PRM_ATT)
}

// RF GAIN
func RG(param string, dir int) (*rigs.RigCommand, error) {
	// Valid values are 0 to 255
	const RE="^0(?P<P1>[01][0-9][0-9]|2[0-4][0-9]|25[0-5])"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return singleValRsp(RE, param, rigs.PRM_RFGAIN, "0")
}

// READ METERS
func RM(param string, dir int) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>[0-9]{1})(?P<P2>.*)"
	const RE2="^(?P<P1>[0-9]{3})"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
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
func RL(param string, dir int) (*rigs.RigCommand, error) {
	// Valid values are 0 to 15
	const RE="^0(?P<P1>0[0-9]|1[0-5])"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return singleValRsp(RE, param, rigs.PRM_NRLVEL, "0")
}

//SMETER
func SM(param string, dir int) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>[0-9]{4})"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return singleValRsp(RE, param, rigs.PRM_SMETER, "0")
}

//VOX
func VX(param string, dir int) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>0|1)"
	if dir == rigs.CAT_DIR_DOWN {
		return nil, fmt.Errorf("Not Implemented")
	}
	return switchValRsp(RE, param, rigs.PRM_VOX)

}

func (r *ft991a) Open()  error {
	return fmt.Errorf("Not implemented")

}

// Converts a CAT response/autoinfo into a RIG command

func  (r *ft991a) OnCat(command string, direction int) (*rigs.RigCommand, error) {

	re  := cmdParser.FindStringSubmatch(command)

	if len(re) < 2 {
		return nil, fmt.Errorf("Not a valid command")
	}
	cmdId := re[1]
	log.Printf("[DEBUG] FT991A: Command %s, Direction %d", cmdId, direction)

	if funcs, ok := funcMap[cmdId]; ok {
		return funcs(re[2], direction)
	}

	return nil, fmt.Errorf("Not a known command")
}

func  (r *ft991a) OnRig(command *rigs.RigCommand, direction int) (string, error) {
	return "", fmt.Errorf("Not Implemented")
}

func New() rigs.Rig {

	return &ft991a {}
}
