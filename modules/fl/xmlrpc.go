package fl

import (
	"log"
	"net/http"
	"github.com/gorilla/rpc"
	"github.com/divan/gorilla-xmlrpc/xml"
)

type FldigiService struct{}

func (h *FldigiService) FLGetName(r *http.Request, args *struct{}, reply *struct{Name string}) error {
	log.Print("[DEBUG] RPCServer:FLGetName")
	reply.Name = "rhub"
	return nil
}


/*
<?xml.version='1.0'?>
<methodCall>
	<methodName>main.get_trx_status</methodName>
    <params><param><value><base64></base64></value></param></params>
</methodCall>

<?xml.version="1.0"?>
<methodResponse>
	<params><param><value><string>rx<string></value></param></params>
</methodResponse>
*/
// RUMLOG is sending an empty string, but is should not
func (h *FldigiService) GetTrxStatusRL(r *http.Request, args *struct{X string}, reply *struct{Status string}) error {
	log.Print("[DEBUG] RPCServer:GetTrxStatusRL")
	reply.Status = "tx"
	return nil
}

/*
<?xml.version='1.0'?>
<methodCall>
	<methodName>main.get_trx_status</methodName>
    <params></params>
</methodCall>

<?xml.version="1.0"?>
<methodResponse>
	<params><param><value><string>rx<string></value></param></params>
</methodResponse>
*/

func (h *FldigiService) GetTrxStatus(r *http.Request, args *struct{}, reply *struct{Status string}) error {
	log.Print("[DEBUG] RPCServer:GetTrxStatus")
	reply.Status = "tx"
	return nil
}
func (h *FldigiService) GetFrequency(r *http.Request, args *struct{}, reply *struct{Freq float64}) error {
	log.Print("[DEBUG] RPCServer: GetFrequency")
	reply.Freq = 14072000.000000
	return nil
}
/*
<?xml.version='1.0'?>
<methodCall>
	<methodName>tx.get_data</methodName>
    <params></params>
</methodCall>

<?xml.version="1.0"?>
<methodResponse>
	<params><param><value><base64></base64></value></param></params>
</methodResponse>
*/
func (h *FldigiService) TxGetData(r *http.Request, args *struct{}, reply *struct{Status []byte}) error {
	log.Printf("[DEBUG] RPCServer: TxGetData")
	reply.Status = []byte{}
	return nil
}
/*
<?xml.version='1.0'?>
<methodCall>
    <methodName>rig.get_frequency</methodName>
    <params></params>
</methodCall>

<?xml.version="1.0"?>
<methodResponse>
    <params>
        <param>
            <value><double>14071000.000000</double></value>
        </param>
    </params>
</methodResponse>
 */
func (h *FldigiService) RigGetFreq(r *http.Request, args *struct{}, reply *struct{Freq float64}) error {
	log.Print("[DEBUG] RPCServer: RigGetFreq")
	reply.Freq = 14071000.000000
	return nil
}

func (h *FldigiService) RigSetFreq(r *http.Request, args *struct{Freq float64}, reply *struct{Freq float64}) error {
	log.Print("[DEBUG] RPCServer: RigSetFreq")
	reply.Freq = args.Freq
	return nil
}

func (h *FldigiService) RigGetMode(r *http.Request, args *struct{}, reply *struct{Mode string}) error {
	log.Print("[DEBUG] RPCServer: RigGetMode")
	reply.Mode = "LSB"
	return nil
}



func newRPCServer() *rpc.Server {
	r := rpc.NewServer()
	xmlrpcCodec := xml.NewCodec()
	r.RegisterCodec(xmlrpcCodec, "text/xml")
	r.RegisterService(new(FldigiService), "")


	xmlrpcCodec.RegisterAlias("main.get_trx_status", "FldigiService.GetTrxStatus")
	xmlrpcCodec.RegisterAlias("main.get_frequency", "FldigiService.GetFrequency")

	xmlrpcCodec.RegisterAlias("tx.get_data", "FldigiService.TxGetData")
	xmlrpcCodec.RegisterAlias("fldigi.name", "FldigiService.FLGetName")
	xmlrpcCodec.RegisterAlias("rig.get_frequency", "FldigiService.RigGetFreq")
	xmlrpcCodec.RegisterAlias("rig.set_frequency", "FldigiService.RigSetFreq")
	xmlrpcCodec.RegisterAlias("rig.get_mode", "FldigiService.RigGetMode")

	return r
}