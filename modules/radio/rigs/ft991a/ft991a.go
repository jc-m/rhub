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

var cmdParser = regexp.MustCompile(CMD_REGEXP)

var rspMap = map[string]catFunc{
	"AI": AIRsp, // AUTO INFORMATION
	"IF": IFRsp, // INFORMATION
	"FA": FARsp, // VFOA
	"FB": FBRsp, // VFOB
	"FT": FTRsp, // SET TX VFO
	"KS": KSRsp, // KEY SPEED
	"MD": MDRsp, // MODE
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
	var cmdParams map[string]string

	m, err := parse(re, param)
	if err != nil {
		return nil, err
	}
	if len(trim) > 0 {
		cmdParams[key] = strings.TrimLeft(m["P1"],trim)
	} else {
		cmdParams[key] = m["P1"]
	}

	return &rigs.RigCommand{
		Id: rigs.RIG_RSP,
		Params : cmdParams,
	}, nil

}


func AIRsp(param string) (*rigs.RigCommand, error) {
	const RE="^(?P<P1>[0|1])"
	var cmdParams map[string]string

	m, err := parse(RE, param)
	if err != nil {
		return nil, err
	}
	cmdParams[rigs.PRM_AUTOINFO] = getSwitch(m["P1"])

	return &rigs.RigCommand{
		Id: rigs.RIG_RSP,
		Params: cmdParams ,
	}, nil
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
	const RE = "^(?P<P1>[0-3]{1})"

	var cmdParams map[string]string

	m, err := parse(RE, param)
	if err != nil {
		return nil, err
	}

	switch m["P1"] {
	case "0":
		cmdParams[rigs.PRM_TXVFO] = rigs.VAL_VFOA
	case "1":
		cmdParams[rigs.PRM_TXVFO] = rigs.VAL_VFOB
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
	//		P4 : RX Clarifier on/off
	//		P5 : TX Clarifier on/off
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
	id := re[1]

	if f, ok := rspMap[id]; ok {
		return f(re[2])
	}

	return nil, fmt.Errorf("Not a known command")
}

func New() rigs.Rig {

	return &ft991a {}
}
