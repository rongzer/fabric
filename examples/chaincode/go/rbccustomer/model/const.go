// const.go
package model

const (
	CONST_TABLE_CUSTOMER = "Customer" //用户表
	//用户状态
	CONST_CUSTOMER_STATUS_ON   = "3" //正常状态
	CONST_CUSTOMER_STATUS_LOCK = "4" //锁定状态
	CONST_CUSTOMER_STATUS_OFF  = "5" //注销状态
	//用户角色
	CONST_CUSTOMER_ROLE_SUPPER  = "1" //超级用户
	CONST_CUSTOMER_ROLE_AUDITOR = "2" //审计用户
	CONST_CUSTOMER_ROLE_B       = "3" //B端用户
	CONST_CUSTOMER_ROLE_C       = "4" //C端用户

	//加密信息
	CONST_PRIVATE_DATA = "********" //C端用户
)

var (
	CUSTOMER_STATUSES = []string{CONST_CUSTOMER_STATUS_ON, CONST_CUSTOMER_STATUS_LOCK, CONST_CUSTOMER_STATUS_OFF}
	CUSTOMER_ROLES    = []string{CONST_CUSTOMER_ROLE_SUPPER, CONST_CUSTOMER_ROLE_AUDITOR, CONST_CUSTOMER_ROLE_B, CONST_CUSTOMER_ROLE_C}
)

func IsInCUSTOMER_STATUS(s string) bool {
	for _, v := range CUSTOMER_STATUSES {
		if v == s {
			return true
		}
	}
	return false
}
func IsInCUSTOMER_ROLE(s string) bool {
	for _, v := range CUSTOMER_ROLES {
		if v == s {
			return true
		}
	}
	return false
}
