// customer_impl
package main

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/examples/chaincode/go/rbccustomer/model"
)

// 验证ID、Cert是否匹配
func (t *CustomerChaincode) verifyIDCert(id, cert string) error {
	//	uuid, _, err := util.GetAffiliationFromString(cert)
	//	if err != nil {
	//		return err
	//	}
	//	if id != uuid {
	//		return errors.New("ID does not match with Cert")
	//	}
	return nil
}

// 将 customerNo,customerSignCert 索引至 customerId
func (t *CustomerChaincode) createCustomerIndex(stub *shim.ChaincodeStub, cus *model.Customer) error {
	if cus.CustomerId == "" {
		return errors.New("customerId is nil")
	}
	var err error
	var key string
	var b []byte
	if cus.CustomerNo != "" {
		key = "customerNo" + cus.CustomerNo
		b, _ = stub.GetState(key)
		if len(b) != 0 {
			if string(b) != cus.CustomerId {
				return errors.New("customerNo already exists")
			}
		}
		err = stub.PutState(key, []byte(cus.CustomerId))
		if err != nil {
			return err
		}
	}
	if cus.CustomerSignCert != "" {
		key = "customerSignCert" + cus.CustomerSignCert
		b, _ = stub.GetState(key)
		if len(b) != 0 {
			if string(b) != cus.CustomerId {
				return errors.New("customerSignCert already exists")
			}
		}
		err = stub.PutState(key, []byte(cus.CustomerId))
		if err != nil {
			return err
		}
	}
	return nil
}

// 根据 customerNo,customerSignCert 获取 customerId
func (t *CustomerChaincode) getCustomerIdByIndex(stub *shim.ChaincodeStub, key, value string) (string, error) {
	var b []byte
	var err error
	switch key {
	case "customerId":
		return value, nil
	case "customerNo":
		b, err = stub.GetState("customerNo" + value)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case "customerSignCert":
		b, err = stub.GetState("customerSignCert" + value)
		if err != nil {
			return "", err
		}
		return string(b), nil
	default:
		return "", errors.New("index does not exists")
	}
}

// 删除 customerNo,customerSignCert 索引
func (t *CustomerChaincode) delCustomerIndex(stub *shim.ChaincodeStub, key, value string) error {
	var err error
	switch key {
	case "customerNo":
		err = stub.DelState("customerNo" + value)
		if err != nil {
			return err
		}
		return nil
	case "customerSignCert":
		err = stub.DelState("customerSignCert" + value)
		if err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}

// args:[]
func (t *CustomerChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var cus model.Customer
	table_cus, err := cus.GetTable(stub)
	if err != nil {
		logger.Error("Customer GetTable", err)
	}
	if table_cus != nil {
		logger.Debug("Customer Table struct:", table_cus)
	} else {
		err = cus.CreateTable(stub)
		if err != nil {
			return nil, err
		}
		logger.Info("Customer Chaincode Table successfully!\n")
	}

	return nil, nil
}

// 会员注册
// args[0]:会员信息json字符串
func (t *CustomerChaincode) register(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	//获取参数，转移成map[string]string
	var m map[string]string
	err := json.Unmarshal([]byte(args[0]), &m)
	if err != nil {
		return nil, err
	}

	//初始化Customer对象
	cus := new(model.Customer)
	cus.CustomerStatus = model.CONST_CUSTOMER_STATUS_ON
	cus.Dict = make(map[string]string)
	for k, v := range m {
		switch k {
		case "customerId":
			cus.CustomerId = v
		case "customerNo":
			cus.CustomerNo = v
		case "customerName":
			cus.CustomerName = v
		case "customerSignCert":
			cus.CustomerSignCert = v
		case "customerRole":
			cus.CustomerRole = v
		case "customerAuth":
			cus.CustomerAuth = v
		default:
			cus.Dict[k] = v
		}
	}

	//验证ID与Cert是否匹配
	err = t.verifyIDCert(cus.CustomerId, cus.CustomerSignCert)
	if err != nil {
		return nil, err
	}

	err = t.createCustomerIndex(stub, cus)
	if err != nil {
		return nil, err
	}

	isExist, err := cus.IsExist(stub)
	if err != nil {
		logger.Error("customer register fail:", err)
		return nil, err
	}
	if isExist {
		logger.Error("customer register fail:customer already exists")
		return nil, errors.New("customer already exists")
	}
	err = cus.Insert(stub)
	if err != nil {
		logger.Error("customer register fail:", err)
		return nil, err
	}
	return nil, nil
}

// 会员信息修改
// args[0]:会员信息json字符串
func (t *CustomerChaincode) modify(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	//获取参数，转移成map[string]string
	var m map[string]string
	err := json.Unmarshal([]byte(args[0]), &m)
	if err != nil {
		return nil, err
	}
	//获取要修改的会员customerId
	var customerId string
	for k, v := range m {
		switch k {
		case "customerId":
			customerId = v
			break
		}
	}
	//查询会员信息
	cus := new(model.Customer)
	cus.CustomerId = customerId
	cus, err = cus.GetRow(stub)
	if err != nil {
		return nil, err
	}

	for k, v := range m {
		switch k {
		case "customerId":
			continue
		case "customerNo":
			if cus.CustomerNo != v {
				err = t.delCustomerIndex(stub, "customerNo", cus.CustomerNo)
				if err != nil {
					return nil, err
				}
				cus.CustomerNo = v
			}
		case "customerName":
			cus.CustomerName = v
		case "customerSignCert":
			continue
		case "customerRole":
			continue
		case "customerAuth":
			cus.CustomerAuth = v
		default:
			cus.Dict[k] = v
		}
	}
	err = t.createCustomerIndex(stub, cus)
	if err != nil {
		return nil, err
	}
	err = cus.Update(stub)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// 会员状态变更
// args:[ID,Status]
func (t *CustomerChaincode) modifyStatus(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	//获取参数，转移成map[string]string
	var m map[string]string
	err := json.Unmarshal([]byte(args[0]), &m)
	if err != nil {
		return nil, err
	}
	//获取customerId, customerStatus的值
	var customerId, customerStatus string
	for k, v := range m {
		switch k {
		case "customerId":
			customerId = v
		case "customerStatus":
			customerStatus = v
		}
	}
	//判断customerStatus是否合法
	if !model.IsInCUSTOMER_STATUS(customerStatus) {
		return nil, errors.New("customer status is illegal!")
	}

	cus := new(model.Customer)
	cus.CustomerId = customerId
	cus, err = cus.GetRow(stub)
	if err != nil {
		return nil, err
	}

	//判断customerStatus是否与原值相同
	if cus.CustomerStatus == customerStatus {
		return nil, errors.New("customer status is the same!")
	}

	cus.CustomerStatus = customerStatus
	err = cus.Update(stub)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// 重置验证证书
// args:[ID,Cert]
func (t *CustomerChaincode) resetCert(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	//获取参数，转移成map[string]string
	var m map[string]string
	err := json.Unmarshal([]byte(args[0]), &m)
	if err != nil {
		return nil, err
	}
	//获取 customerId, customerSignCert 的值
	var customerId, customerSignCert string
	for k, v := range m {
		switch k {
		case "customerId":
			customerId = v
		case "customerSignCert":
			customerSignCert = v
		}
	}

	//验证 customerId,customerSignCert 是否匹配
	err = t.verifyIDCert(customerId, customerSignCert)
	if err != nil {
		return nil, err
	}

	cus := new(model.Customer)
	cus.CustomerId = customerId
	cus, err = cus.GetRow(stub)
	if err != nil {
		return nil, err
	}

	//判断传入cert是否与原值相同
	if cus.CustomerSignCert == customerSignCert {
		return nil, errors.New("customer SignCert is the same!")
	}

	err = t.delCustomerIndex(stub, "customerSignCert", cus.CustomerSignCert)
	if err != nil {
		return nil, err
	}
	cus.CustomerSignCert = customerSignCert

	err = t.createCustomerIndex(stub, cus)
	if err != nil {
		return nil, err
	}
	err = cus.Update(stub)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// 重置加密公私钥
// args[0]:{"customerId":"customerId"}
func (t *CustomerChaincode) resetKey(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	//获取参数，转移成 map[string]string
	var m map[string]string
	err := json.Unmarshal([]byte(args[0]), &m)
	if err != nil {
		return nil, err
	}
	//获取 customerId
	var customerId string
	for k, v := range m {
		switch k {
		case "customerId":
			customerId = v
			break
		}
	}

	cus := new(model.Customer)
	cus.CustomerId = customerId
	cus, err = cus.GetRow(stub)
	if err != nil {
		return nil, err
	}

	logger.Info("resetKey", cus.CustomerId)
	return nil, nil
}

// 根据用户ID获取用户详细
// args[0]:{"customerId":"customerId"}|{"customerNo":"customerNo"}|{"customerSignCert":"customerSignCert"}
func (t *CustomerChaincode) queryOne(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	//获取参数，转移成 map[string]string
	var m map[string]string
	err := json.Unmarshal([]byte(args[0]), &m)
	if err != nil {
		return nil, err
	}
	//获取 customerId
	var customerId string
	for k, v := range m {
		switch k {
		case "customerId":
			customerId = v
		case "customerNo":
			customerId, err = t.getCustomerIdByIndex(stub, "customerNo", v)
		case "customerSignCert":
			customerId, err = t.getCustomerIdByIndex(stub, "customerSignCert", v)
		}
	}
	if err != nil {
		return nil, err
	}

	cus := new(model.Customer)
	cus.CustomerId = customerId
	cus, err = cus.GetRow(stub)
	if err != nil {
		return nil, err
	}

	res, err := json.Marshal(cus)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// 获取所有用户详细
// args:[{"pagenum":"pagenum"}]
// pagenum:页数，可以为空
func (t *CustomerChaincode) queryAll(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	//获取参数，转移成 map[string]string
	var m map[string]string
	err := json.Unmarshal([]byte(args[0]), &m)
	if err != nil {
		return nil, err
	}
	//获取 pagenum
	var pagenum string
	for k, v := range m {
		switch k {
		case "pagenum":
			pagenum = v
		}
	}

	var num int64
	num, _ = strconv.ParseInt(pagenum, 10, 64)

	var cus model.Customer
	cs, err := cus.GetRows(stub, pagesize, num)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(cs)
	if err != nil {
		return nil, err
	}

	return res, nil
}
