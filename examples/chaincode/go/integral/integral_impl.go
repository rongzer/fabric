// integral_impl
package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	chaincode_ecdsa "github.com/hyperledger/fabric/core/chaincode/shim/crypto/ecdsa"
	"github.com/hyperledger/fabric/examples/chaincode/go/integral/model"
	"github.com/hyperledger/fabric/examples/chaincode/go/integral/util"
)

const (
	per_amount float64 = 1000000.00
	pagesize   int64   = 10
)

func (t *IntegralChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//创建普通用户表
	var user model.User
	tableUser, err := user.GetTable(stub)
	if err != nil {
		logger.Error("User GetTable", err)
	}
	if tableUser != nil {
		logger.Debug("User Table struct:", tableUser)
	} else {
		err = user.CreateTable(stub)
		if err != nil {
			return nil, err
		}
		logger.Info("Create User Table successfully!\n")
	}

	//创建联盟用户表
	var unionuser model.UnionUser
	tableUnionUser, err := unionuser.GetTable(stub)
	if err != nil {
		logger.Error("UnionUser GetTable", err)
	}
	if tableUnionUser != nil {
		logger.Debug("UnionUser Table struct:", tableUnionUser)
	} else {
		err = unionuser.CreateTable(stub)
		if err != nil {
			return nil, err
		}
		logger.Info("Create UnionUser Table successfully!\n")
	}

	//初始化发行参数
	//发行次数
	bdt, err := stub.GetState("distribute_times")
	if err != nil {
		logger.Error("获取发行次数:", err)
		return nil, err
	}
	idt, _ := strconv.Atoi(string(bdt))
	logger.Info("发行次数:", idt)
	if idt == 0 {
		err = stub.PutState("distribute_times", []byte("1"))
		if err != nil {
			return nil, err
		}
	}

	//发行总量
	b_distribute_total, err := stub.GetState("distribute_total")
	if err != nil {
		logger.Error("获取发行总量:", err)
		return nil, err
	}
	logger.Info("发行总量:", string(b_distribute_total))

	return nil, nil
}

func (t *IntegralChaincode) distribute(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	cert, err := stub.GetCallerCertificate()
	if err != nil {
		return nil, err
	}
	logger.Info("hex.EncodeToString(cert)", hex.EncodeToString(cert))
	scert := hex.EncodeToString(cert)
	affs, err := util.GetAffiliation(cert)
	if err != nil {
		logger.Error("GetAffiliation:", err)
		return nil, err
	}
	uuid := affs[0]
	aff := affs[1]
	if uuid == "" {
		logger.Error("uuid is nil")
		return nil, errors.New("uuid is nil")
	}
	logger.Error("affs", affs)
	if aff == "business" {
		var u model.UnionUser
		u.Uuid = uuid
		isExist, err := u.IsExist(stub)
		if err != nil {
			return nil, err
		}
		if isExist {
			return nil, errors.New("UnionUser already exists")
		}
		u.Cert = scert
		if len(args) > 0 {
			u.Affilication = uuid
		} else {
			u.Affilication = ""
		}
		bdt, err := stub.GetState("distribute_times")
		if err != nil {
			logger.Error("distribute_times:", err)
		}
		idt, _ := strconv.Atoi(string(bdt))
		u.Times = int32(idt)
		amount := float64(idt) * per_amount
		u.Total = amount
		u.Balance = amount
		u.UserStatus = model.Const_User_Status_On
		u.IntegralStatus = model.Const_Integral_Status_On
		err = u.InsertRow(stub)
		if err != nil {
			return nil, err
		}
		logger.Info("Insert UnionUser successfully!\n")

		b_distribute_total, err := stub.GetState("distribute_total")
		if err != nil {
			logger.Error("获取发行总量:", err)
			return nil, err
		}
		logger.Info("发行总量:", string(b_distribute_total))
		f_distribute_total, _ := strconv.ParseFloat(string(b_distribute_total), 64)
		f_distribute_total = f_distribute_total + amount

		err = stub.PutState("distribute_total", []byte(fmt.Sprint(f_distribute_total)))

		if err != nil {
			return nil, err
		}

		txuuid := util.GenerateUUID()
		err = u.InitIntegral(stub, txuuid, amount)
		if err != nil {
			return nil, err
		}

	} else {
		var u model.User
		u.Uuid = uuid
		isExist, err := u.IsExist(stub)
		if err != nil {
			return nil, err
		}
		if isExist {
			return nil, errors.New("User already exists")
		}
		u.Cert = scert
		u.UserStatus = model.Const_User_Status_On
		u.IntegralStatus = model.Const_Integral_Status_On
		err = u.InsertRow(stub)
		if err != nil {
			return nil, err
		}
		logger.Info("Insert User successfully!\n")

		err = u.InitIntegral(stub)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (t *IntegralChaincode) transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	uuid, aff, err := t.getaffs(stub)
	if err != nil {
		return nil, err
	}
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting >= 2")
	}
	cert := args[0]
	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, err
	}
	touuid, toaff, err := t.getaffsFromString(cert)
	if err != nil {
		return nil, err
	}
	txuuid := util.GenerateUUID()

	err = model.Transfer(stub, txuuid, uuid, aff, touuid, toaff, amount)
	if err != nil {
		return nil, err
	}
	//	if aff == model.Const_User_Aff_Business {
	//		u := new(model.UnionUser)
	//		u.Uuid = uuid
	//		err = u.TransferOut(stub, touuid, txuuid, amount)
	//		if err != nil {
	//			return nil, err
	//		}

	//	} else {
	//		u := new(model.User)
	//		u.Uuid = uuid
	//		err = u.TransferOut(stub, touuid, txuuid, amount)
	//		if err != nil {
	//			return nil, err
	//		}
	//	}

	//	if toaff == model.Const_User_Aff_Business {
	//		//转账给联盟用户
	//		to := new(model.UnionUser)
	//		to.Uuid = touuid
	//		err = to.TransferIn(stub, uuid, txuuid, amount)
	//		if err != nil {
	//			return nil, err
	//		}
	//	} else {
	//		//转账给普通用户
	//		to := new(model.User)
	//		to.Uuid = touuid
	//		err = to.TransferIn(stub, uuid, txuuid, amount)
	//		if err != nil {
	//			return nil, err
	//		}
	//	}
	return nil, nil
}

// 商户所剩积分总额低于一共发行积分的30%，积分发行一次
func (t *IntegralChaincode) inspect(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var u model.UnionUser
	users, err := u.GetRows(stub, 0, 0)
	if err != nil {
		return nil, err
	}
	var balanceTotal float64 = 0
	for _, v := range users {
		balanceTotal = balanceTotal + v.Balance
	}
	bdt, err := stub.GetState("distribute_times")
	if err != nil {
		logger.Error("获取发行次数:", err)
	}
	idt, _ := strconv.Atoi(string(bdt))

	b_distribute_total, err := stub.GetState("distribute_total")
	if err != nil {
		logger.Error("获取发行总量:", err)
		return nil, err
	}
	total, _ := strconv.ParseFloat(string(b_distribute_total), 64)

	logger.Info("发行次数:", string(bdt))
	logger.Info("每次每人发行额度:", per_amount)
	logger.Info("联盟用户数量:", len(users))
	logger.Info("发行总量:", total)
	logger.Info("联盟用户剩余总额:", balanceTotal)

	if balanceTotal < total*0.3 {
		logger.Info("<30%:再次发行")
		err = stub.PutState("distribute_times", []byte(strconv.Itoa(idt+1)))
		if err != nil {
			return nil, err
		}
		txuuid := util.GenerateUUID()
		for _, v := range users {
			err = v.Redistribute(stub, txuuid, per_amount)
			logger.Error(err)
		}

		amount := float64(len(users)) * per_amount
		total = total + amount
		err = stub.PutState("distribute_total", []byte(fmt.Sprint(total)))
		if err != nil {
			return nil, err
		}

	}
	return nil, nil
}
func (t *IntegralChaincode) payment(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) < 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting >= 4")
	}
	uuid, aff, err := t.getaffs(stub)
	if err != nil {
		return nil, err
	}
	logger.Info("uuid:", uuid)
	logger.Info("aff:", aff)

	var cert, sign, msg []byte
	if args[0] == "" {
		return nil, errors.New("cert is nil")
	} else {
		cert, _ = hex.DecodeString(args[0])
	}

	if args[1] == "" {
		return nil, errors.New("sign is nil")
	} else {
		sign, _ = hex.DecodeString(args[1])
	}
	msg, _ = hex.DecodeString(args[2])

	sv := chaincode_ecdsa.NewX509ECDSASignatureVerifier()
	// Verify the signature
	ok, err := sv.Verify(cert, sign, msg)
	if err != nil {
		logger.Error(err)
	}
	if ok {
		logger.Info("verify sign success")
	} else {
		logger.Error("verify sign fail")
	}

	fromuuid, fromaff, err := t.getaffsFromString(args[0])
	if err != nil {
		return nil, err
	}
	logger.Info("fromuuid:", fromuuid)
	logger.Info("fromaff:", fromaff)
	amount, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return nil, err
	}
	logger.Info("amount:", amount)

	txuuid := util.GenerateUUID()
	err = model.Transfer(stub, txuuid, fromuuid, fromaff, uuid, aff, amount)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
func (t *IntegralChaincode) getaffsFromString(scert string) (string, string, error) {
	cert, err := hex.DecodeString(scert)
	if err != nil {
		return "", "", err
	}
	logger.Info("hex.EncodeToString(cert)", hex.EncodeToString(cert))
	affs, err := util.GetAffiliation(cert)
	if err != nil {
		logger.Error("GetAffiliation:", err)
		return "", "", err
	}
	logger.Info("affs:", affs)
	uuid := affs[0]
	aff := affs[1]
	if uuid == "" {
		return "", "", errors.New("uuid is nil")
	}
	return uuid, aff, nil
}
func (t *IntegralChaincode) getaffs(stub *shim.ChaincodeStub) (string, string, error) {
	cert, err := stub.GetCallerCertificate()
	if err != nil {
		return "", "", err
	}
	logger.Info("hex.EncodeToString(cert)", hex.EncodeToString(cert))
	affs, err := util.GetAffiliation(cert)
	if err != nil {
		logger.Error("GetAffiliation:", err)
		return "", "", err
	}
	logger.Info("affs:", affs)
	uuid := affs[0]
	aff := affs[1]
	if uuid == "" {
		return "", "", errors.New("uuid is nil")
	}
	return uuid, aff, nil
}
func (t *IntegralChaincode) getown(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	uuid, aff, err := t.getaffs(stub)
	if err != nil {
		return nil, err
	}
	if aff == model.Const_User_Aff_Business {
		var u model.UnionUser
		u.Uuid = uuid
		res, err := u.GetRow(stub)
		if err != nil {
			return nil, err
		}
		jsonRow, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}

		return jsonRow, nil
	} else {
		var u model.User
		u.Uuid = uuid
		res, err := u.GetRow(stub)
		if err != nil {
			return nil, err
		}
		jsonRow, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}

		return jsonRow, nil
	}
	return nil, nil
}

func (t *IntegralChaincode) getuser(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting >0")
	}
	var u model.User
	u.Uuid = args[0]
	res, err := u.GetRow(stub)
	if err != nil {
		return nil, err
	}
	jsonRow, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return jsonRow, nil
}

func (t *IntegralChaincode) getusers(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var pagenum int64
	if len(args) == 0 {
		pagenum = 0
	} else {
		pagenum, _ = strconv.ParseInt(args[0], 10, 64)
	}
	var u model.User
	users, err := u.GetRows(stub, pagesize, pagenum)
	if err != nil {
		return nil, err
	}
	jsonRows, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}

	return jsonRows, nil
}

func (t *IntegralChaincode) getunionuser(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting >0")
	}
	var u model.UnionUser
	u.Uuid = args[0]
	res, err := u.GetRow(stub)
	if err != nil {
		return nil, err
	}
	jsonRow, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return jsonRow, nil
}

func (t *IntegralChaincode) getunionusers(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var pagenum int64
	if len(args) == 0 {
		pagenum = 0
	} else {
		pagenum, _ = strconv.ParseInt(args[0], 10, 64)
	}
	var u model.UnionUser
	users, err := u.GetRows(stub, pagesize, pagenum)
	if err != nil {
		return nil, err
	}
	jsonRows, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}

	return jsonRows, nil
}

func (t *IntegralChaincode) gettrade(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	uuid, _, err := t.getaffs(stub)
	if err != nil {
		return nil, err
	}
	var pagenum int64
	if len(args) == 0 {
		pagenum = 0
	} else {
		pagenum, _ = strconv.ParseInt(args[0], 10, 64)
	}
	logger.Info("页数：", pagenum)
	var trade model.Trade
	ts, err := trade.GetRows(stub, uuid, pagesize, pagenum)
	if err != nil {
		return nil, err
	}
	jsonRows, err := json.Marshal(ts)
	if err != nil {
		return nil, err
	}

	return jsonRows, nil
}
func (t *IntegralChaincode) gettradein(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	uuid, _, err := t.getaffs(stub)
	if err != nil {
		return nil, err
	}
	var pagenum int64
	if len(args) == 0 {
		pagenum = 0
	} else {
		pagenum, _ = strconv.ParseInt(args[0], 10, 64)
	}
	logger.Info("页数：", pagenum)
	var trade model.TradeIn
	ts, err := trade.GetRows(stub, uuid, pagesize, pagenum)
	if err != nil {
		return nil, err
	}
	jsonRows, err := json.Marshal(ts)
	if err != nil {
		return nil, err
	}

	return jsonRows, nil
}
func (t *IntegralChaincode) gettradeout(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	uuid, _, err := t.getaffs(stub)
	if err != nil {
		return nil, err
	}
	var pagenum int64
	if len(args) == 0 {
		pagenum = 0
	} else {
		pagenum, _ = strconv.ParseInt(args[0], 10, 64)
	}
	logger.Info("页数：", pagenum)
	var trade model.TradeOut
	ts, err := trade.GetRows(stub, uuid, pagesize, pagenum)
	if err != nil {
		return nil, err
	}
	jsonRows, err := json.Marshal(ts)
	if err != nil {
		return nil, err
	}

	return jsonRows, nil
}
