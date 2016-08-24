// unionuser
// 联盟用户
package model

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/examples/chaincode/go/integral/util"
	"github.com/op/go-logging"
)

var table_unionuser = "unionuser"
var logger_unionuser = logging.MustGetLogger("unionuser")

type UnionUser struct {
	Uuid           string  `json:"Uuid"`          //Cert的hash值（主键）
	Cert           string  `json:"Cert"`          //用户Cert
	Affilication   string  `json:"Affilication"`  //集团 例如：PHHS KFC
	Times          int32   `json:"Times"`         //发行次数
	Total          float64 `json:"Total"`         //发行总额
	Balance        float64 `json:"Balance"`       //余额
	UserStatus     string  `json:"UserStatus"`    //用户状态 1:正常 2:禁用
	IntegralStatus string  `json:"IntegralState"` //资金状态 1:正常 2:冻结
	ExtAttr        string  `json:"ExtAttr"`       //扩展属性
}

// 查找联盟用户表，返回(*shim.Table,error)
func (u *UnionUser) GetTable(stub *shim.ChaincodeStub) (*shim.Table, error) {
	table, err := stub.GetTable(table_unionuser)
	if err != nil {
		return nil, err
	}
	return table, nil
}

// 创建联盟用户表，返回error
func (u *UnionUser) CreateTable(stub *shim.ChaincodeStub) error {
	var colDefs []*shim.ColumnDefinition
	uuid := shim.ColumnDefinition{Name: "Uuid",
		Type: shim.ColumnDefinition_STRING, Key: true}
	cert := shim.ColumnDefinition{Name: "Cert",
		Type: shim.ColumnDefinition_STRING, Key: false}
	aff := shim.ColumnDefinition{Name: "Affilication",
		Type: shim.ColumnDefinition_STRING, Key: false}
	times := shim.ColumnDefinition{Name: "Times",
		Type: shim.ColumnDefinition_INT32, Key: false}
	total := shim.ColumnDefinition{Name: "Total",
		Type: shim.ColumnDefinition_STRING, Key: false}
	balance := shim.ColumnDefinition{Name: "Balance",
		Type: shim.ColumnDefinition_STRING, Key: false}
	ustatus := shim.ColumnDefinition{Name: "UserStatus",
		Type: shim.ColumnDefinition_STRING, Key: false}
	istatus := shim.ColumnDefinition{Name: "IntegralStatus",
		Type: shim.ColumnDefinition_STRING, Key: false}
	extAttr := shim.ColumnDefinition{Name: "ExtAttr",
		Type: shim.ColumnDefinition_STRING, Key: false}

	colDefs = append(colDefs, &uuid, &cert, &aff, &times, &total, &balance, &ustatus, &istatus, &extAttr)

	return stub.CreateTable(table_unionuser, colDefs)
}

// 插入联盟用户，返回error
func (u *UnionUser) InsertRow(stub *shim.ChaincodeStub) error {
	var cols []*shim.Column
	uuid := shim.Column{Value: &shim.Column_String_{String_: u.Uuid}}
	cert := shim.Column{Value: &shim.Column_String_{String_: u.Cert}}
	aff := shim.Column{Value: &shim.Column_String_{String_: u.Affilication}}
	times := shim.Column{Value: &shim.Column_Int32{Int32: u.Times}}
	total := shim.Column{Value: &shim.Column_String_{String_: fmt.Sprint(u.Total)}}
	balance := shim.Column{Value: &shim.Column_String_{String_: fmt.Sprint(u.Balance)}}
	ustatus := shim.Column{Value: &shim.Column_String_{String_: u.UserStatus}}
	istatus := shim.Column{Value: &shim.Column_String_{String_: u.IntegralStatus}}
	extAttr := shim.Column{Value: &shim.Column_String_{String_: u.ExtAttr}}

	cols = append(cols, &uuid, &cert, &aff, &times, &total, &balance, &ustatus, &istatus, &extAttr)

	row := shim.Row{cols}
	fmt.Println(&row)
	ok, err := stub.InsertRow(table_unionuser, row)
	if err != nil {
		return fmt.Errorf("Insert row failed. %s", err)
	}
	if !ok {
		return errors.New("Insert row failed. Row with given key already exists")
	}
	return nil
}

// 更新联盟用户，返回error
func (u *UnionUser) UpdateRow(stub *shim.ChaincodeStub) error {
	var cols []*shim.Column
	uuid := shim.Column{Value: &shim.Column_String_{String_: u.Uuid}}
	cert := shim.Column{Value: &shim.Column_String_{String_: u.Cert}}
	aff := shim.Column{Value: &shim.Column_String_{String_: u.Affilication}}
	times := shim.Column{Value: &shim.Column_Int32{Int32: u.Times}}
	total := shim.Column{Value: &shim.Column_String_{String_: fmt.Sprint(u.Total)}}
	balance := shim.Column{Value: &shim.Column_String_{String_: fmt.Sprint(u.Balance)}}
	ustatus := shim.Column{Value: &shim.Column_String_{String_: u.UserStatus}}
	istatus := shim.Column{Value: &shim.Column_String_{String_: u.IntegralStatus}}
	extAttr := shim.Column{Value: &shim.Column_String_{String_: u.ExtAttr}}

	cols = append(cols, &uuid, &cert, &aff, &times, &total, &balance, &ustatus, &istatus, &extAttr)

	row := shim.Row{cols}
	fmt.Println(&row)
	ok, err := stub.ReplaceRow(table_unionuser, row)
	if err != nil {
		return fmt.Errorf("Update row failed. %s", err)
	}
	if !ok {
		return errors.New("Update row failed. Row with given key not exist")
	}
	return nil
}

// 判断联盟用户是否存在，返回(bool, error)
func (u *UnionUser) IsExist(stub *shim.ChaincodeStub) (bool, error) {
	if u.Uuid == "" {
		return false, errors.New("Primary key is nil!")
	}

	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: u.Uuid}}
	columns = append(columns, col1)

	row, err := stub.GetRow(table_unionuser, columns)
	if err != nil {
		return false, fmt.Errorf("Get table row failed. %s", err)
	}

	if len(row.Columns) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// 根据uuid获取联盟用户，返回(*UnionUser, error)
func (u *UnionUser) GetRow(stub *shim.ChaincodeStub) (*UnionUser, error) {
	if u.Uuid == "" {
		return nil, errors.New("Primary key is nil!")
	}

	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: u.Uuid}}
	columns = append(columns, col1)

	row, err := stub.GetRow(table_unionuser, columns)
	if err != nil {
		return nil, fmt.Errorf("Get table row failed. %s", err)
	}

	user := new(UnionUser)
	if len(row.Columns) > 0 {
		user.Uuid = row.Columns[0].GetString_()
		user.Cert = row.Columns[1].GetString_()
		user.Affilication = row.Columns[2].GetString_()
		user.Times = row.Columns[3].GetInt32()
		user.Total, _ = strconv.ParseFloat(row.Columns[4].GetString_(), 64)
		user.Balance, _ = strconv.ParseFloat(row.Columns[5].GetString_(), 64)
		user.UserStatus = row.Columns[6].GetString_()
		user.IntegralStatus = row.Columns[7].GetString_()
		user.ExtAttr = row.Columns[8].GetString_()

		return user, nil
	} else {
		return nil, errors.New("Row is nil")
	}
}

// 获取所有联盟用户列表，返回([]*UnionUser, error)
func (u *UnionUser) GetRows(stub *shim.ChaincodeStub, pagesize, pagenum int64) ([]*UnionUser, error) {
	rowChannel, err := stub.GetRows(table_unionuser, nil)
	if err != nil {
		return nil, fmt.Errorf("Get Table rows failed. %s", err)
	}

	users := make([]*UnionUser, 0)
	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {
				user := new(UnionUser)
				user.Uuid = row.Columns[0].GetString_()
				user.Cert = row.Columns[1].GetString_()
				user.Affilication = row.Columns[2].GetString_()
				user.Times = row.Columns[3].GetInt32()
				user.Total, _ = strconv.ParseFloat(row.Columns[4].GetString_(), 64)
				user.Balance, _ = strconv.ParseFloat(row.Columns[5].GetString_(), 64)
				user.UserStatus = row.Columns[6].GetString_()
				user.IntegralStatus = row.Columns[7].GetString_()
				user.ExtAttr = row.Columns[8].GetString_()

				users = append(users, user)
			}
		}
		if rowChannel == nil {
			if pagenum <= 0 {
				return users, nil
			} else {
				begin, end := util.PageRow(pagesize, pagenum, int64(len(users)))
				return users[begin:end], nil
			}
			//return users, nil
		}
	}
}

// 删除联盟用户
func (u *UnionUser) DeleteRow(stub *shim.ChaincodeStub) error {
	if u.Uuid == "" {
		return errors.New("Primary key is nil!")
	}

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: u.Uuid}}
	columns = append(columns, col1)

	err := stub.DeleteRow(table_unionuser, columns)
	if err != nil {
		return fmt.Errorf("Delete row failed. %s", err)
	}
	return nil
}

func (u *UnionUser) InitIntegral(stub *shim.ChaincodeStub, txuuid string, amount float64) error {
	//创建入账表
	var trade_in TradeIn
	table, err := trade_in.GetTable(stub, u.Uuid)
	if err != nil {
		logger_unionuser.Error(err)
	}
	if table != nil {
		logger_unionuser.Info(table)
	} else {
		err = trade_in.CreateTable(stub, u.Uuid)
		if err != nil {
			logger_unionuser.Error(err)
		}
	}
	//入账记录
	trade_in.TxUuid = txuuid
	trade_in.Amount = amount
	trade_in.FromUuid = ""
	trade_in.BeforeTrade = 0
	trade_in.AfterTrade = amount
	trade_in.Time = time.Now().Format(layout)
	err = trade_in.InsertRow(stub, u.Uuid)
	if err != nil {
		logger_unionuser.Error(err)
	}

	//创建出账表
	var trade_out TradeOut
	table, err = trade_out.GetTable(stub, u.Uuid)
	if err != nil {
		logger_user.Error(err)
	}
	if table != nil {
		logger_user.Info(table)
	} else {
		err = trade_out.CreateTable(stub, u.Uuid)
		if err != nil {
			logger_user.Error(err)
		}
	}
	//创建所有交易表
	var trade_all Trade
	table, err = trade_all.GetTable(stub, u.Uuid)
	if err != nil {
		logger_unionuser.Error(err)
	}
	if table != nil {
		logger_unionuser.Info(table)
	} else {
		err = trade_all.CreateTable(stub, u.Uuid)
		if err != nil {
			logger_unionuser.Error(err)
		}
	}
	//所有记录
	trade_all.TxUuid = txuuid
	trade_all.Amount = amount
	trade_all.Type = Const_Trade_Type_In
	trade_all.FromToUuid = ""
	trade_all.BeforeTrade = 0
	trade_all.AfterTrade = amount
	trade_all.Time = time.Now().Format(layout)
	err = trade_all.InsertRow(stub, u.Uuid)
	if err != nil {
		logger_unionuser.Error(err)
	}
	return nil
}
func (u *UnionUser) Redistribute(stub *shim.ChaincodeStub, txuuid string, amount float64) error {
	var err error
	u, err = u.GetRow(stub)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.New("UnionUser is nil")
	}

	curBalance := u.Balance
	newBalance := curBalance + amount
	//转入
	u.Times = u.Times + 1
	u.Balance = newBalance
	u.Total = u.Total + amount
	err = u.UpdateRow(stub)
	if err != nil {
		logger_unionuser.Error(err)
	}

	//入账记录
	var trade_in TradeIn
	trade_in.TxUuid = txuuid
	trade_in.Amount = amount
	trade_in.FromUuid = ""
	trade_in.BeforeTrade = curBalance
	trade_in.AfterTrade = newBalance
	trade_in.Time = time.Now().Format(layout)
	err = trade_in.InsertRow(stub, u.Uuid)
	if err != nil {
		logger_unionuser.Error(err)
	}
	//所有记录
	var trade_all Trade
	trade_all.TxUuid = txuuid
	trade_all.Amount = amount
	trade_all.Type = Const_Trade_Type_In
	trade_all.FromToUuid = ""
	trade_all.BeforeTrade = curBalance
	trade_all.AfterTrade = newBalance
	trade_all.Time = time.Now().Format(layout)
	err = trade_all.InsertRow(stub, u.Uuid)
	if err != nil {
		logger_unionuser.Error(err)
	}
	return nil
}
func (u *UnionUser) transferIn(stub *shim.ChaincodeStub, fromuuid, txuuid string, amount float64) error {
	var err error
	u, err = u.GetRow(stub)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.New("UnionUser is nil")
	}
	if u.UserStatus != Const_User_Status_On || u.IntegralStatus != Const_Integral_Status_On {
		return errors.New("UnionUser UserStatus or IntegralStatus is not normal")
	}

	curBalance := u.Balance
	newBalance := curBalance + amount
	//转入
	u.Balance = newBalance
	err = u.UpdateRow(stub)
	if err != nil {
		logger_unionuser.Error(err)
	}

	//入账记录
	var trade_in TradeIn
	trade_in.TxUuid = txuuid
	trade_in.Amount = amount
	trade_in.FromUuid = fromuuid
	trade_in.BeforeTrade = curBalance
	trade_in.AfterTrade = newBalance
	trade_in.Time = time.Now().Format(layout)
	err = trade_in.InsertRow(stub, u.Uuid)
	if err != nil {
		logger_unionuser.Error(err)
	}
	//所有记录
	var trade_all Trade
	trade_all.TxUuid = txuuid
	trade_all.Amount = amount
	trade_all.Type = Const_Trade_Type_In
	trade_all.FromToUuid = fromuuid
	trade_all.BeforeTrade = curBalance
	trade_all.AfterTrade = newBalance
	trade_all.Time = time.Now().Format(layout)
	err = trade_all.InsertRow(stub, u.Uuid)
	if err != nil {
		logger_unionuser.Error(err)
	}
	return nil
}

func (u *UnionUser) transferOut(stub *shim.ChaincodeStub, touuid, txuuid string, amount float64) error {
	var err error
	u, err = u.GetRow(stub)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.New("UnionUser is nil")
	}
	if u.UserStatus != Const_User_Status_On || u.IntegralStatus != Const_Integral_Status_On {
		return errors.New("UnionUser UserStatus or IntegralStatus is not normal")
	}
	if u.Balance < amount {
		return errors.New("UnionUser Balance is not enough")
	}

	curBalance := u.Balance
	newBalance := curBalance - amount
	//转出
	u.Balance = newBalance
	err = u.UpdateRow(stub)
	if err != nil {
		logger_unionuser.Error(err)
	}

	//出账记录
	var trade_out TradeOut
	trade_out.TxUuid = txuuid
	trade_out.Amount = amount
	trade_out.ToUuid = touuid
	trade_out.BeforeTrade = curBalance
	trade_out.AfterTrade = newBalance
	trade_out.Time = time.Now().Format(layout)
	err = trade_out.InsertRow(stub, u.Uuid)
	if err != nil {
		logger_unionuser.Error(err)
	}
	//所有记录
	var trade_all Trade
	trade_all.TxUuid = txuuid
	trade_all.Amount = amount
	trade_all.Type = Const_Trade_Type_Out
	trade_all.FromToUuid = touuid
	trade_all.BeforeTrade = curBalance
	trade_all.AfterTrade = newBalance
	trade_all.Time = time.Now().Format(layout)
	err = trade_all.InsertRow(stub, u.Uuid)
	if err != nil {
		logger_unionuser.Error(err)
	}
	return nil
}
func (u *UnionUser) verifyfrom(stub *shim.ChaincodeStub, amount float64) (*UnionUser, error) {
	u, err := u.GetRow(stub)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("UnionUser is nil")
	}
	if u.UserStatus != Const_User_Status_On || u.IntegralStatus != Const_Integral_Status_On {
		return nil, errors.New("UnionUser UserStatus or IntegralStatus is not normal")
	}
	if u.Balance < amount {
		return nil, errors.New("UnionUser Balance is not enough")
	}
	return u, nil
}
func (u *UnionUser) verifyto(stub *shim.ChaincodeStub) (*UnionUser, error) {
	u, err := u.GetRow(stub)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("UnionUser is nil")
	}
	if u.UserStatus != Const_User_Status_On || u.IntegralStatus != Const_Integral_Status_On {
		return nil, errors.New("UnionUser UserStatus or IntegralStatus is not normal")
	}
	return u, nil
}
func Transfer(stub *shim.ChaincodeStub, txuuid, fromuuid, fromaff, touuid, toaff string, amount float64) error {
	var err error
	if fromaff == Const_User_Aff_Business {
		from := new(UnionUser)
		from.Uuid = fromuuid
		from, err = from.verifyfrom(stub, amount)
		if err != nil {
			return err
		}
		if toaff == Const_User_Aff_Business {
			to := new(UnionUser)
			to.Uuid = touuid
			to, err = to.verifyto(stub)
			if err != nil {
				return err
			}
			err = from.transferOut(stub, touuid, txuuid, amount)
			if err != nil {
				return err
			}
			err = to.transferIn(stub, fromuuid, txuuid, amount)
			if err != nil {
				return err
			}
		} else {
			to := new(User)
			to.Uuid = touuid
			to, err = to.verifyto(stub)
			if err != nil {
				return err
			}
			err = from.transferOut(stub, touuid, txuuid, amount)
			if err != nil {
				return err
			}
			err = to.transferIn(stub, fromuuid, txuuid, amount)
			if err != nil {
				return err
			}
		}

	} else {
		from := new(User)
		from.Uuid = fromuuid
		from, err = from.verifyfrom(stub, amount)
		if err != nil {
			return err
		}
		if toaff == Const_User_Aff_Business {
			to := new(UnionUser)
			to.Uuid = touuid
			to, err = to.verifyto(stub)
			if err != nil {
				return err
			}
			err = from.transferOut(stub, touuid, txuuid, amount)
			if err != nil {
				return err
			}
			err = to.transferIn(stub, fromuuid, txuuid, amount)
			if err != nil {
				return err
			}

		} else {
			to := new(User)
			to.Uuid = touuid
			to, err = to.verifyto(stub)
			if err != nil {
				return err
			}
			err = from.transferOut(stub, touuid, txuuid, amount)
			if err != nil {
				return err
			}
			err = to.transferIn(stub, fromuuid, txuuid, amount)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
