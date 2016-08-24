// integral
package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("integral")

type IntegralChaincode struct {
}

func (t *IntegralChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return t.init(stub, args)
}

func (t *IntegralChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	b, _ := stub.GetState(stub.UUID)
	if len(b) != 0 {
		return nil, errors.New("the uuid transaction alredy invoked")
	} else {
		stub.PutState(stub.UUID, []byte("ok"))
	}
	if function == "init" {
		return t.init(stub, args)
	} else if function == "distribute" {
		return t.distribute(stub, args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	} else if function == "inspect" {
		return t.inspect(stub, args)
	} else if function == "payment" {
		return t.payment(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}
func (t *IntegralChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "getown" {
		return t.getown(stub, args)
	} else if function == "getuser" {
		return t.getuser(stub, args)
	} else if function == "getusers" {
		return t.getusers(stub, args)
	} else if function == "getunionuser" {
		return t.getunionuser(stub, args)
	} else if function == "getunionusers" {
		return t.getunionusers(stub, args)
	} else if function == "gettradein" {
		return t.gettradein(stub, args)
	} else if function == "gettradeout" {
		return t.gettradeout(stub, args)
	} else if function == "gettrade" {
		return t.gettrade(stub, args)
	}

	return nil, errors.New("Invalid query unknown function invocation")
}

func main() {
	err := shim.Start(new(IntegralChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
