// trade
// 交易记录，包括入账、出账记录
package model

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/examples/chaincode/go/integral/util"
)

const (
	Const_Trade_Type_In  string = "1"
	Const_Trade_Type_Out string = "2"
)

type Trade struct {
	TxUuid      string  `json:"TxUuid"`      //交易uuid，主键
	Amount      float64 `json:"Amount"`      //交易积分数量
	Type        string  `json:"Type"`        //交易类型 1:入账；2:出账
	FromToUuid  string  `json:"FromToUuid"`  //交易对方账户uuid
	BeforeTrade float64 `json:"BeforeTrade"` //交易之前
	AfterTrade  float64 `json:"AfterTrade"`  //交易之后
	Time        string  `json:"Time"`        //交易时间
	ExtAttr     string  `json:"ExtAttr"`     //扩展属性
}

// 获取交易表，返回(*shim.Table, error)
func (t *Trade) GetTable(stub *shim.ChaincodeStub, uuid string) (*shim.Table, error) {
	table_name := uuid + "_all"
	table, err := stub.GetTable(table_name)
	if err != nil {
		return nil, err
	}
	return table, nil
}

// 创建交易表，返回error
func (t *Trade) CreateTable(stub *shim.ChaincodeStub, uuid string) error {
	table_name := uuid + "_all"
	var colDefs []*shim.ColumnDefinition
	txuuid := shim.ColumnDefinition{Name: "TxUuid",
		Type: shim.ColumnDefinition_STRING, Key: true}
	amount := shim.ColumnDefinition{Name: "Amount",
		Type: shim.ColumnDefinition_STRING, Key: false}
	tType := shim.ColumnDefinition{Name: "Type",
		Type: shim.ColumnDefinition_STRING, Key: false}
	fromuuid := shim.ColumnDefinition{Name: "FromUuid",
		Type: shim.ColumnDefinition_STRING, Key: false}
	before := shim.ColumnDefinition{Name: "BeforeTrade",
		Type: shim.ColumnDefinition_STRING, Key: false}
	after := shim.ColumnDefinition{Name: "AfterTrade",
		Type: shim.ColumnDefinition_STRING, Key: false}
	time := shim.ColumnDefinition{Name: "Time",
		Type: shim.ColumnDefinition_STRING, Key: false}
	extAttr := shim.ColumnDefinition{Name: "ExtAttr",
		Type: shim.ColumnDefinition_STRING, Key: false}

	colDefs = append(colDefs, &txuuid, &amount, &tType, &fromuuid, &before, &after, &time, &extAttr)

	return stub.CreateTable(table_name, colDefs)
}

// 插入交易记录，返回error
func (t *Trade) InsertRow(stub *shim.ChaincodeStub, uuid string) error {
	table_name := uuid + "_all"
	var cols []*shim.Column
	txuuid := shim.Column{Value: &shim.Column_String_{String_: t.TxUuid}}
	amount := shim.Column{Value: &shim.Column_String_{String_: fmt.Sprint(t.Amount)}}
	tType := shim.Column{Value: &shim.Column_String_{String_: t.Type}}
	fromuuid := shim.Column{Value: &shim.Column_String_{String_: t.FromToUuid}}
	before := shim.Column{Value: &shim.Column_String_{String_: fmt.Sprint(t.BeforeTrade)}}
	after := shim.Column{Value: &shim.Column_String_{String_: fmt.Sprint(t.AfterTrade)}}
	time := shim.Column{Value: &shim.Column_String_{String_: t.Time}}
	extAttr := shim.Column{Value: &shim.Column_String_{String_: t.ExtAttr}}

	cols = append(cols, &txuuid, &amount, &tType, &fromuuid, &before, &after, &time, &extAttr)

	row := shim.Row{cols}
	fmt.Println(&row)
	ok, err := stub.InsertRow(table_name, row)
	if err != nil {
		return fmt.Errorf("Insert row failed. %s", err)
	}
	if !ok {
		return errors.New("Insert row failed. Row with given key already exists")
	}
	return nil
}

// 判断交易信息是否存在，返回(bool, error)
func (t *Trade) IsExist(stub *shim.ChaincodeStub, uuid string) (bool, error) {
	if t.TxUuid == "" {
		return false, errors.New("Primary key is nil!")
	}
	table_name := uuid + "_all"
	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: t.TxUuid}}
	columns = append(columns, col1)

	row, err := stub.GetRow(table_name, columns)
	if err != nil {
		return false, fmt.Errorf("Get table row failed. %s", err)
	}

	if len(row.Columns) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// 获取交易信息，返回(*Trade, error)
func (u *Trade) GetRow(stub *shim.ChaincodeStub, uuid string) (*Trade, error) {
	if u.TxUuid == "" {
		return nil, errors.New("Primary key is nil!")
	}
	table_name := uuid + "_all"
	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: u.TxUuid}}
	columns = append(columns, col1)

	row, err := stub.GetRow(table_name, columns)
	if err != nil {
		return nil, fmt.Errorf("Get table row failed. %s", err)
	}

	trade := new(Trade)
	if len(row.Columns) > 0 {
		trade.TxUuid = row.Columns[0].GetString_()
		trade.Amount, _ = strconv.ParseFloat(row.Columns[1].GetString_(), 64)
		trade.Type = row.Columns[2].GetString_()
		trade.FromToUuid = row.Columns[3].GetString_()
		trade.BeforeTrade, _ = strconv.ParseFloat(row.Columns[4].GetString_(), 64)
		trade.AfterTrade, _ = strconv.ParseFloat(row.Columns[5].GetString_(), 64)
		trade.Time = row.Columns[6].GetString_()
		trade.ExtAttr = row.Columns[7].GetString_()

		return trade, nil
	} else {
		return nil, errors.New("Row is nil")
	}
}

// 获取交易列表，返回([]*Trade, error)
func (t *Trade) GetRows(stub *shim.ChaincodeStub, uuid string, pagesize, pagenum int64) ([]*Trade, error) {
	table_name := uuid + "_all"
	rowChannel, err := stub.GetRows(table_name, nil)
	if err != nil {
		return nil, fmt.Errorf("Get Table rows failed. %s", err)
	}

	trades := make([]*Trade, 0)
	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {
				trade := new(Trade)
				trade.TxUuid = row.Columns[0].GetString_()
				trade.Amount, _ = strconv.ParseFloat(row.Columns[1].GetString_(), 64)
				trade.Type = row.Columns[2].GetString_()
				trade.FromToUuid = row.Columns[3].GetString_()
				trade.BeforeTrade, _ = strconv.ParseFloat(row.Columns[4].GetString_(), 64)
				trade.AfterTrade, _ = strconv.ParseFloat(row.Columns[5].GetString_(), 64)
				trade.Time = row.Columns[6].GetString_()
				trade.ExtAttr = row.Columns[7].GetString_()

				trades = append(trades, trade)
			}
		}
		if rowChannel == nil {
			list := TradeList(trades)
			sort.Sort(list)

			if pagenum <= 0 {
				return list, nil
			} else {
				begin, end := util.PageRow(pagesize, pagenum, int64(len(list)))
				return list[begin:end], nil
			}
			//return trades, nil
		}
	}
}

// 删除交易信息，返回error
func (t *Trade) DeleteRow(stub *shim.ChaincodeStub, uuid string) error {
	if t.TxUuid == "" {
		return errors.New("Primary key is nil!")
	}
	table_name := uuid + "_all"
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: t.TxUuid}}
	columns = append(columns, col1)

	err := stub.DeleteRow(table_name, columns)
	if err != nil {
		return fmt.Errorf("Delete row failed. %s", err)
	}
	return nil
}

type TradeList []*Trade

//排序规则：按照交易时间倒叙排列
func (list TradeList) Len() int {
	return len(list)
}

func (list TradeList) Less(i, j int) bool {
	if list[i].Time > list[j].Time {
		return true
	} else if list[i].Time < list[j].Time {
		return false
	} else {
		return list[i].TxUuid < list[j].TxUuid
	}
}

func (list TradeList) Swap(i, j int) {
	var temp *Trade = list[i]
	list[i] = list[j]
	list[j] = temp
}
