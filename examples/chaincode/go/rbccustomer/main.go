// main
package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("user")

const (
	pagesize   int64  = 10
	timeformat string = "2006-01-02 15:04:05"
)

type CustomerChaincode struct {
}

func (t *CustomerChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return t.init(stub, args)
}

func (t *CustomerChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	//确保transaction只执行一次
	b, _ := stub.GetState(stub.UUID)
	if len(b) != 0 {
		return nil, errors.New("the uuid transaction alredy invoked")
	} else {
		stub.PutState(stub.UUID, []byte("ok"))
	}

	if function == "register" { //注册会员
		return t.register(stub, args)
	} else if function == "modify" { //修改会员信息
		return t.modify(stub, args)
	} else if function == "modifyStatus" { //会员状态变更
		return t.modifyStatus(stub, args)
	} else if function == "resetCert" { //重置验证证书
		return t.resetCert(stub, args)
	} else if function == "resetKey" { //重置加密公私钥
		return t.resetKey(stub, args)
	}
	return nil, nil
}

func (t *CustomerChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "queryOne" { //根据用户ID获取用户详细
		return t.queryOne(stub, args)
	} else if function == "queryAll" { //获取所有用户详细
		return t.queryAll(stub, args)
	}
	return nil, nil
}

func main() {
	err := shim.Start(new(CustomerChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
