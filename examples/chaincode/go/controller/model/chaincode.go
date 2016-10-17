// chaincode
package model

import (
	"errors"
	"fmt"
	//	"strconv"
	//	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/examples/chaincode/go/util"
	"github.com/op/go-logging"
)

const const_table_chaincode = "chaincode"

var logger_chaincode = logging.MustGetLogger("chaincode")

type Chaincode struct {
	Alias   string `json:"Alias"`   //chaincode别名(别名+版本组合主键)
	Version string `json:"Version"` //chaincode版本(别名+版本组合主键)
	Cert    string `json:"Cert"`    //发布人Cert
	Code    string `json:"Code"`    //chaincode代码内容
	Text    string `json:"Text"`    //chaincode文本描述
	Name    string `json:"Name"`    //chaincode名
	Time    string `json:"Time"`    //注册时间
	Extend  string `json:"Extend"`  //扩展属性
}

// 查找链码表，返回(*shim.Table,error)
func (c *Chaincode) GetTable(stub *shim.ChaincodeStub) (*shim.Table, error) {
	table, err := stub.GetTable(const_table_chaincode)
	if err != nil {
		return nil, err
	}
	return table, nil
}

// 创建链码表，返回error
func (c *Chaincode) CreateTable(stub *shim.ChaincodeStub) error {
	var colDefs []*shim.ColumnDefinition
	alias := shim.ColumnDefinition{Name: "Alias",
		Type: shim.ColumnDefinition_STRING, Key: true}
	version := shim.ColumnDefinition{Name: "Version",
		Type: shim.ColumnDefinition_STRING, Key: true}
	cert := shim.ColumnDefinition{Name: "Cert",
		Type: shim.ColumnDefinition_STRING, Key: false}
	code := shim.ColumnDefinition{Name: "Code",
		Type: shim.ColumnDefinition_STRING, Key: false}
	text := shim.ColumnDefinition{Name: "Text",
		Type: shim.ColumnDefinition_STRING, Key: false}
	name := shim.ColumnDefinition{Name: "Name",
		Type: shim.ColumnDefinition_STRING, Key: false}
	stime := shim.ColumnDefinition{Name: "Time",
		Type: shim.ColumnDefinition_STRING, Key: false}
	extend := shim.ColumnDefinition{Name: "Extend",
		Type: shim.ColumnDefinition_STRING, Key: false}

	colDefs = append(colDefs, &alias, &version, &cert, &code, &text, &name, &stime, &extend)

	return stub.CreateTable(const_table_chaincode, colDefs)
}

func (c *Chaincode) keys() ([]shim.Column, error) {
	if c.Alias == "" {
		return nil, errors.New("Alias key is nil!")
	}
	if c.Version == "" {
		return nil, errors.New("Version key is nil!")
	}

	var cols []shim.Column
	alias := shim.Column{Value: &shim.Column_String_{String_: c.Alias}}
	version := shim.Column{Value: &shim.Column_String_{String_: c.Version}}
	cols = append(cols, alias, version)
	return cols, nil
}

// 获取链码表的字段列表
func (c *Chaincode) columns() ([]*shim.Column, error) {
	if c.Alias == "" {
		return nil, errors.New("Alias key is nil!")
	}
	if c.Version == "" {
		return nil, errors.New("Version key is nil!")
	}

	var cols []*shim.Column
	alias := shim.Column{Value: &shim.Column_String_{String_: c.Alias}}
	version := shim.Column{Value: &shim.Column_String_{String_: c.Version}}
	cert := shim.Column{Value: &shim.Column_String_{String_: c.Cert}}
	code := shim.Column{Value: &shim.Column_String_{String_: c.Code}}
	text := shim.Column{Value: &shim.Column_String_{String_: c.Text}}
	name := shim.Column{Value: &shim.Column_String_{String_: c.Name}}
	stime := shim.Column{Value: &shim.Column_String_{String_: c.Time}}
	extend := shim.Column{Value: &shim.Column_String_{String_: c.Extend}}

	cols = append(cols, &alias, &version, &cert, &code, &text, &name, &stime, &extend)
	return cols, nil
}

// 插入链码版本，返回error
func (c *Chaincode) Insert(stub *shim.ChaincodeStub) error {
	cols, err := c.columns()
	if err != nil {
		return err
	}
	row := shim.Row{cols}
	fmt.Println(&row)
	ok, err := stub.InsertRow(const_table_chaincode, row)
	if err != nil {
		return fmt.Errorf("Insert Chaincode failed. %s", err)
	}
	if !ok {
		return errors.New("Insert Chaincode failed. Chaincode with given key already exists")
	}
	return nil
}

// 更新链码版本，返回error
func (c *Chaincode) Update(stub *shim.ChaincodeStub) error {
	cols, err := c.columns()
	if err != nil {
		return err
	}
	row := shim.Row{cols}
	fmt.Println(&row)
	ok, err := stub.ReplaceRow(const_table_chaincode, row)
	if err != nil {
		return fmt.Errorf("Update Chaincode failed. %s", err)
	}
	if !ok {
		return errors.New("Update Chaincode failed. Chaincode with given key not exist")
	}
	return nil
}

// 判断链码版本是否存在，返回(bool,error)
func (c *Chaincode) IsExist(stub *shim.ChaincodeStub) (bool, error) {
	cols, err := c.keys()
	if err != nil {
		return false, err
	}
	row, err := stub.GetRow(const_table_chaincode, cols)
	if err != nil {
		return false, fmt.Errorf("Get Chaincode row failed. %s", err)
	}

	if len(row.Columns) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// 根据(别名,版本)获取链码版本详细
func (c *Chaincode) GetRow(stub *shim.ChaincodeStub) (*Chaincode, error) {
	cols, err := c.keys()
	if err != nil {
		return nil, err
	}
	row, err := stub.GetRow(const_table_chaincode, cols)
	if err != nil {
		return nil, fmt.Errorf("Get Chaincode row failed. %s", err)
	}

	chaincode := new(Chaincode)
	if len(row.Columns) > 0 {
		chaincode.Alias = row.Columns[0].GetString_()
		chaincode.Version = row.Columns[1].GetString_()
		chaincode.Cert = row.Columns[2].GetString_()
		chaincode.Code = row.Columns[3].GetString_()
		chaincode.Text = row.Columns[4].GetString_()
		chaincode.Name = row.Columns[5].GetString_()
		chaincode.Time = row.Columns[6].GetString_()
		chaincode.Extend = row.Columns[7].GetString_()

		return chaincode, nil
	} else {
		return nil, errors.New("Chaincode is nil")
	}
}

// 根据链码别名获取版本列表
func (c *Chaincode) GetRowsByAlias(stub *shim.ChaincodeStub, pagesize, pagenum int64) ([]*Chaincode, error) {
	if c.Alias == "" {
		return nil, errors.New("Alias key is nil!")
	}
	var cols []shim.Column
	alias := shim.Column{Value: &shim.Column_String_{String_: c.Alias}}
	cols = append(cols, alias)
	rowChannel, err := stub.GetRows(const_table_chaincode, cols)
	if err != nil {
		return nil, fmt.Errorf("Get Chaincode rows failed. %s", err)
	}

	chaincodes := make([]*Chaincode, 0)
	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {
				chaincode := new(Chaincode)
				chaincode.Alias = row.Columns[0].GetString_()
				chaincode.Version = row.Columns[1].GetString_()
				chaincode.Cert = row.Columns[2].GetString_()
				chaincode.Code = row.Columns[3].GetString_()
				chaincode.Text = row.Columns[4].GetString_()
				chaincode.Name = row.Columns[5].GetString_()
				chaincode.Time = row.Columns[6].GetString_()
				chaincode.Extend = row.Columns[7].GetString_()
				chaincodes = append(chaincodes, chaincode)
			}
		}
		if rowChannel == nil {
			if pagenum <= 0 {
				return chaincodes, nil
			} else {
				begin, end := util.PageRow(pagesize, pagenum, int64(len(chaincodes)))
				return chaincodes[begin:end], nil
			}
			return chaincodes, nil
		}
	}
}

// 获取所有链码版本列表
func (c *Chaincode) GetRows(stub *shim.ChaincodeStub, pagesize, pagenum int64) ([]*Chaincode, error) {
	rowChannel, err := stub.GetRows(const_table_chaincode, nil)
	if err != nil {
		return nil, fmt.Errorf("Get Chaincode rows failed. %s", err)
	}

	chaincodes := make([]*Chaincode, 0)
	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {
				chaincode := new(Chaincode)
				chaincode.Alias = row.Columns[0].GetString_()
				chaincode.Version = row.Columns[1].GetString_()
				chaincode.Cert = row.Columns[2].GetString_()
				chaincode.Code = row.Columns[3].GetString_()
				chaincode.Text = row.Columns[4].GetString_()
				chaincode.Name = row.Columns[5].GetString_()
				chaincode.Time = row.Columns[6].GetString_()
				chaincode.Extend = row.Columns[7].GetString_()
				chaincodes = append(chaincodes, chaincode)
			}
		}
		if rowChannel == nil {
			if pagenum <= 0 {
				return chaincodes, nil
			} else {
				begin, end := util.PageRow(pagesize, pagenum, int64(len(chaincodes)))
				return chaincodes[begin:end], nil
			}
			return chaincodes, nil
		}
	}
}

// 删除链码版本
func (c *Chaincode) DeleteRow(stub *shim.ChaincodeStub) error {
	cols, err := c.keys()
	if err != nil {
		return err
	}
	err = stub.DeleteRow(const_table_chaincode, cols)
	if err != nil {
		return fmt.Errorf("Delete Chaincode row failed. %s", err)
	}
	return nil
}
