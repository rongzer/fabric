// controller_impl
package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	//	"encoding/hex"
	//	"encoding/json"
	//	"errors"
	//	"fmt"
	//	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	// chaincode_ecdsa "github.com/hyperledger/fabric/core/chaincode/shim/crypto/ecdsa"
	"github.com/hyperledger/fabric/examples/chaincode/go/controller/model"
	_ "github.com/hyperledger/fabric/examples/chaincode/go/util"
)

const (
	pagesize   int64  = 10
	timeformat string = "2006-01-02 15:04:05"
)

// args:[]
func (t *ControllerChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var chaincode model.Chaincode
	table_chaincode, err := chaincode.GetTable(stub)
	if err != nil {
		logger.Error("Chaincode GetTable", err)
	}
	if table_chaincode != nil {
		logger.Debug("Chaincode Table struct:", table_chaincode)
	} else {
		err = chaincode.CreateTable(stub)
		if err != nil {
			return nil, err
		}
		logger.Info("Create Chaincode Table successfully!\n")
	}

	var defaultchaincode model.DefaultChaincode
	table_default, err := defaultchaincode.GetTable(stub)
	if err != nil {
		logger.Error("Chaincode GetTable", err)
	}
	if table_default != nil {
		logger.Debug("Chaincode Table struct:", defaultchaincode)
	} else {
		err = defaultchaincode.CreateTable(stub)
		if err != nil {
			return nil, err
		}
		logger.Info("Create Chaincode Table successfully!\n")
	}
	return nil, nil
}

// args:[alias,version,code,text,name,time,extend]
func (t *ControllerChaincode) register(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 7 {
		return nil, errors.New("Incorrect number of arguments. Expecting 7")
	}
	stime := args[5]
	_, err := time.ParseInLocation(timeformat, stime, time.Local)
	if err != nil {
		return nil, err
	}
	cert, _ := stub.GetCallerCertificate()

	chaincode := new(model.Chaincode)
	chaincode.Alias = args[0]
	chaincode.Version = args[1]
	chaincode.Cert = hex.EncodeToString(cert)
	chaincode.Code = args[2]
	chaincode.Text = args[3]
	chaincode.Name = args[4]
	chaincode.Time = stime
	chaincode.Extend = args[6]

	isExist, err := chaincode.IsExist(stub)
	if err != nil {
		logger.Error("Chaincode register fail:", err)
		return nil, err
	}
	if isExist {
		logger.Error("Chaincode register fail:Row already exists")
		return nil, errors.New("Row already exists")
	}
	err = chaincode.Insert(stub)
	if err != nil {
		logger.Error("Chaincode register fail:", err)
		return nil, err
	}

	return nil, nil
}

// args:[alias,version,time]
func (t *ControllerChaincode) setDefault(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}
	stime := args[2]
	_, err := time.ParseInLocation(timeformat, stime, time.Local)
	if err != nil {
		return nil, err
	}
	cert, _ := stub.GetCallerCertificate()

	chaincode := new(model.Chaincode)
	chaincode.Alias = args[0]
	chaincode.Version = args[1]
	chaincode, err = chaincode.GetRow(stub)
	if err != nil {
		return nil, err
	}
	if chaincode == nil {
		return nil, errors.New("chaincode getrow is nil")
	}
	defaultchaincode := new(model.DefaultChaincode)
	defaultchaincode.Alias = chaincode.Alias
	defaultchaincode.Version = chaincode.Version
	defaultchaincode.Cert = hex.EncodeToString(cert)
	defaultchaincode.Code = chaincode.Code
	defaultchaincode.Text = chaincode.Text
	defaultchaincode.Name = chaincode.Name
	defaultchaincode.Time = stime
	defaultchaincode.Extend = chaincode.Extend

	err = defaultchaincode.Insert(stub)
	if err != nil {
		logger.Error("Chaincode register fail:", err)
		return nil, err
	}

	return nil, nil
}

// args:[alias,function...]
func (t *ControllerChaincode) invoke(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) <= 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting > 2")
	}
	defaultchaincode := new(model.DefaultChaincode)
	defaultchaincode.Alias = args[0]
	defaultchaincode, err := defaultchaincode.GetRow(stub)
	if err != nil {
		return nil, err
	}
	name := defaultchaincode.Name
	function := args[1]
	newargs := args[2:]

	if name == "" {
		return nil, errors.New("chaincodename is nil")
	}
	response, err := stub.InvokeChaincode(name, function, newargs)
	logger.Infof("chaincode:%s,function:%s,args:%v", name, function, newargs)
	if err != nil {
		logger.Errorf("Failed to invoke chaincode. Got error: %s", err.Error())
		return nil, fmt.Errorf("Failed to invoke chaincode. Got error: %s", err.Error())
	}

	logger.Infof("Invoke chaincode successful. Got response %s", string(response))

	return nil, nil
}

// args:[alias]
func (t *ControllerChaincode) getDefault(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting >0")
	}
	d := new(model.DefaultChaincode)
	d.Alias = args[0]
	d, err := d.GetRow(stub)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// args:[pagenum]
// pagenum:页数，可以不传
func (t *ControllerChaincode) getDefaults(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var pagenum int64
	if len(args) == 0 {
		pagenum = 0
	} else {
		pagenum, _ = strconv.ParseInt(args[0], 10, 64)
	}
	var d model.DefaultChaincode
	ds, err := d.GetRows(stub, pagesize, pagenum)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// args:[alias,version]
func (t *ControllerChaincode) getChaincode(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting >0")
	}
	c := new(model.Chaincode)
	c.Alias = args[0]
	c.Version = args[1]
	c, err := c.GetRow(stub)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// args:[alias,pagenum]
// alias:别名  pagenum:页数，可以不传
func (t *ControllerChaincode) getChaincodeVersions(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var pagenum int64
	if len(args) <= 1 {
		pagenum = 0
	} else {
		pagenum, _ = strconv.ParseInt(args[1], 10, 64)
	}
	c := new(model.Chaincode)
	c.Alias = args[0]
	cs, err := c.GetRowsByAlias(stub, pagesize, pagenum)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(cs)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// args:[pagenum]
// pagenum:页数，可以不传
func (t *ControllerChaincode) getChaincodes(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var pagenum int64
	if len(args) == 0 {
		pagenum = 0
	} else {
		pagenum, _ = strconv.ParseInt(args[0], 10, 64)
	}
	c := new(model.Chaincode)
	cs, err := c.GetRows(stub, pagesize, pagenum)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(cs)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// args:[alias,function...]
func (t *ControllerChaincode) query(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) <= 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting > 2")
	}
	defaultchaincode := new(model.DefaultChaincode)
	defaultchaincode.Alias = args[0]
	defaultchaincode, err := defaultchaincode.GetRow(stub)
	if err != nil {
		return nil, err
	}
	name := defaultchaincode.Name
	function := args[1]
	newargs := args[2:]

	if name == "" {
		return nil, errors.New("chaincodename is nil")
	}
	res, err := stub.QueryChaincode(name, function, newargs)
	logger.Infof("chaincode:%s,function:%s,args:%v", name, function, newargs)
	if err != nil {
		logger.Errorf("Failed to query chaincode. Got error: %s", err.Error())
		return nil, fmt.Errorf("Failed to query chaincode. Got error: %s", err.Error())
	}

	logger.Infof("Query chaincode successful. Got response %s", string(res))

	return res, nil
}
