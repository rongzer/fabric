// controller
package main

import (
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("controller")

type ControllerChaincode struct {
}

func (t *ControllerChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return t.init(stub, args)
}

func (t *ControllerChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	b, _ := stub.GetState(stub.UUID)
	if len(b) != 0 {
		return nil, errors.New("the uuid transaction alredy invoked")
	} else {
		stub.PutState(stub.UUID, []byte("ok"))
	}
	b2, _ := stub.GetCallerCertificate()
	fmt.Println("cert", hex.EncodeToString(b2))
	cert, err := x509.ParseCertificate(b2)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	fmt.Println(cert.Subject.CommonName)
	if function == "register" {
		return t.register(stub, args)
	} else if function == "setDefault" {
		return t.setDefault(stub, args)
	} else if function == "invoke" {
		return t.invoke(stub, args)
	}
	return nil, nil
}

func (t *ControllerChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	b2, _ := stub.GetCallerCertificate()
	fmt.Println("cert", hex.EncodeToString(b2))
	cert, err := x509.ParseCertificate(b2)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	fmt.Println(cert.Subject.CommonName)
	if function == "getDefault" {
		return t.getDefault(stub, args)
	} else if function == "getDefaults" {
		return t.getDefaults(stub, args)
	} else if function == "getChaincode" {
		return t.getChaincode(stub, args)
	} else if function == "getChaincodes" {
		return t.getChaincodes(stub, args)
	} else if function == "getChaincodeVersions" {
		return t.getChaincodeVersions(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	}
	return nil, nil
}

func main() {
	err := shim.Start(new(ControllerChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
