// customer
package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/examples/chaincode/go/util"
)

type Customer struct {
	CustomerId       string            `json:"customerId"`       //会员CustomerId(主键)
	CustomerNo       string            `json:"customerNo"`       //会员编号
	CustomerName     string            `json:"customerName"`     //会员名称
	CustomerSignCert string            `json:"customerSignCert"` //会员公钥，用于会员验证
	CustomerRole     string            `json:"customerRole"`     //会员角色
	CustomerAuth     string            `json:"customerAuth"`     //会员权限
	CustomerStatus   string            `json:"customerStatus"`   //会员状态:3正常；4锁定；5注销
	DataEucCert      string            `json:"dataEucCert"`      //加密公钥，用于加解密
	DataEucPrivate   string            `json:"dataEucPrivate"`   //加密私钥，用于加解密
	Dict             map[string]string `json:"dict"`             //会员扩展信息Dict
}

// 查找会员表，返回(*shim.Table,error)
func (c *Customer) GetTable(stub *shim.ChaincodeStub) (*shim.Table, error) {
	table, err := stub.GetTable(CONST_TABLE_CUSTOMER)
	if err != nil {
		return nil, err
	}
	return table, nil
}

// 创建会员表，返回error
func (c *Customer) CreateTable(stub *shim.ChaincodeStub) error {
	var colDefs []*shim.ColumnDefinition
	CustomerId := shim.ColumnDefinition{Name: "CustomerId",
		Type: shim.ColumnDefinition_STRING, Key: true}
	CustomerNo := shim.ColumnDefinition{Name: "CustomerNo",
		Type: shim.ColumnDefinition_STRING, Key: false}
	CustomerName := shim.ColumnDefinition{Name: "CustomerName",
		Type: shim.ColumnDefinition_STRING, Key: false}
	CustomerSignCert := shim.ColumnDefinition{Name: "CustomerSignCert",
		Type: shim.ColumnDefinition_STRING, Key: false}
	CustomerRole := shim.ColumnDefinition{Name: "CustomerRole",
		Type: shim.ColumnDefinition_STRING, Key: false}
	CustomerAuth := shim.ColumnDefinition{Name: "CustomerAuth",
		Type: shim.ColumnDefinition_STRING, Key: false}
	CustomerStatus := shim.ColumnDefinition{Name: "CustomerStatus",
		Type: shim.ColumnDefinition_STRING, Key: false}
	DataEucCert := shim.ColumnDefinition{Name: "DataEucCert",
		Type: shim.ColumnDefinition_STRING, Key: false}
	DataEucPrivate := shim.ColumnDefinition{Name: "DataEucPrivate",
		Type: shim.ColumnDefinition_STRING, Key: false}
	Dict := shim.ColumnDefinition{Name: "Dict",
		Type: shim.ColumnDefinition_STRING, Key: false}

	colDefs = append(colDefs, &CustomerId, &CustomerNo, &CustomerName, &CustomerSignCert, &CustomerRole,
		&CustomerAuth, &CustomerStatus, &DataEucCert, &DataEucPrivate, &Dict)

	return stub.CreateTable(CONST_TABLE_CUSTOMER, colDefs)
}

func (c *Customer) keys() ([]shim.Column, error) {
	if c.CustomerId == "" {
		return nil, errors.New("Row key is nil!")
	}

	var cols []shim.Column
	CustomerId := shim.Column{Value: &shim.Column_String_{String_: c.CustomerId}}
	cols = append(cols, CustomerId)
	return cols, nil
}

// 获取会员表的字段列表
func (c *Customer) columns() ([]*shim.Column, error) {
	if c.CustomerId == "" {
		return nil, errors.New("Row key is nil!")
	}

	var cols []*shim.Column
	CustomerId := shim.Column{Value: &shim.Column_String_{String_: c.CustomerId}}
	CustomerNo := shim.Column{Value: &shim.Column_String_{String_: c.CustomerNo}}
	CustomerName := shim.Column{Value: &shim.Column_String_{String_: c.CustomerName}}
	CustomerSignCert := shim.Column{Value: &shim.Column_String_{String_: c.CustomerSignCert}}
	CustomerRole := shim.Column{Value: &shim.Column_String_{String_: c.CustomerRole}}
	CustomerAuth := shim.Column{Value: &shim.Column_String_{String_: c.CustomerAuth}}
	CustomerStatus := shim.Column{Value: &shim.Column_String_{String_: c.CustomerStatus}}
	DataEucCert := shim.Column{Value: &shim.Column_String_{String_: c.DataEucCert}}
	DataEucPrivate := shim.Column{Value: &shim.Column_String_{String_: c.DataEucPrivate}}
	var dict string
	if c.Dict != nil {
		b, err := json.Marshal(c.Dict)
		if err != nil {
			return nil, err
		}
		dict = string(b)
	} else {
		dict = "{}"
	}
	Dict := shim.Column{Value: &shim.Column_String_{String_: dict}}

	cols = append(cols, &CustomerId, &CustomerNo, &CustomerName, &CustomerSignCert, &CustomerRole,
		&CustomerAuth, &CustomerStatus, &DataEucCert, &DataEucPrivate, &Dict)
	return cols, nil
}

// 插入会员，返回error
func (c *Customer) Insert(stub *shim.ChaincodeStub) error {
	cols, err := c.columns()
	if err != nil {
		return err
	}
	row := shim.Row{cols}
	fmt.Println(&row)
	ok, err := stub.InsertRow(CONST_TABLE_CUSTOMER, row)
	if err != nil {
		return fmt.Errorf("Insert row failed. %s", err)
	}
	if !ok {
		return errors.New("Insert row failed. Row with given key already exists")
	}
	return nil
}

// 更新会员，返回error
func (c *Customer) Update(stub *shim.ChaincodeStub) error {
	cols, err := c.columns()
	if err != nil {
		return err
	}
	row := shim.Row{cols}
	fmt.Println(&row)
	ok, err := stub.ReplaceRow(CONST_TABLE_CUSTOMER, row)
	if err != nil {
		return fmt.Errorf("Update row failed. %s", err)
	}
	if !ok {
		return errors.New("Update row failed. Row with given key not exist")
	}
	return nil
}

// 判断会员是否存在，返回(bool,error)
func (c *Customer) IsExist(stub *shim.ChaincodeStub) (bool, error) {
	cols, err := c.keys()
	if err != nil {
		return false, err
	}
	row, err := stub.GetRow(CONST_TABLE_CUSTOMER, cols)
	if err != nil {
		return false, fmt.Errorf("Get table row failed. %s", err)
	}

	if len(row.Columns) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// 根据会员CustomerId获取会员详细
func (c *Customer) GetRow(stub *shim.ChaincodeStub) (*Customer, error) {
	cols, err := c.keys()
	if err != nil {
		return nil, err
	}
	row, err := stub.GetRow(CONST_TABLE_CUSTOMER, cols)
	if err != nil {
		return nil, fmt.Errorf("Get table row failed. %s", err)
	}

	res := new(Customer)
	if len(row.Columns) > 0 {
		res.CustomerId = row.Columns[0].GetString_()
		res.CustomerNo = row.Columns[1].GetString_()
		res.CustomerName = row.Columns[2].GetString_()
		res.CustomerSignCert = row.Columns[3].GetString_()
		res.CustomerRole = row.Columns[4].GetString_()
		res.CustomerAuth = row.Columns[5].GetString_()
		res.CustomerStatus = row.Columns[6].GetString_()
		res.DataEucCert = row.Columns[7].GetString_()
		res.DataEucPrivate = row.Columns[8].GetString_()

		res.Dict = make(map[string]string)
		err := json.Unmarshal([]byte(row.Columns[9].GetString_()), &res.Dict)
		if err != nil {
			return nil, err
		}

		return res, nil
	} else {
		return nil, errors.New("Row is nil")
	}
}

// 获取所有会员列表
func (c *Customer) GetRows(stub *shim.ChaincodeStub, pagesize, pagenum int64) ([]*Customer, error) {
	rowChannel, err := stub.GetRows(CONST_TABLE_CUSTOMER, nil)
	if err != nil {
		return nil, fmt.Errorf("Get Table rows failed. %s", err)
	}

	cs := make([]*Customer, 0)
	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {
				cus := new(Customer)
				cus.CustomerId = row.Columns[0].GetString_()
				cus.CustomerNo = row.Columns[1].GetString_()
				cus.CustomerName = row.Columns[2].GetString_()
				cus.CustomerSignCert = row.Columns[3].GetString_()
				cus.CustomerRole = row.Columns[4].GetString_()
				cus.CustomerAuth = row.Columns[5].GetString_()
				cus.CustomerStatus = row.Columns[6].GetString_()
				cus.DataEucCert = row.Columns[7].GetString_()
				cus.DataEucPrivate = row.Columns[8].GetString_()

				cus.Dict = make(map[string]string)
				err := json.Unmarshal([]byte(row.Columns[9].GetString_()), &cus.Dict)
				if err != nil {
					return nil, err
				}

				cs = append(cs, cus)
			}
		}
		if rowChannel == nil {
			if pagenum <= 0 {
				return cs, nil
			} else {
				begin, end := util.PageRow(pagesize, pagenum, int64(len(cs)))
				return cs[begin:end], nil
			}
			return cs, nil
		}
	}
}
