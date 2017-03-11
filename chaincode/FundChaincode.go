// FundChaincode
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//定义合约
type FundChaincode struct {
}

//定义基金结构体
type FundInfo struct { //基金信息
	Id                 string
	AppSheetSerialNo   string
	FundCode2          string
	TransactionDate    string
	TransactionTime    string
	DistributorCode    string
	BusinessCode       string
	ApplicationVol     string
	ApplicationAmount  string
	TaAccountID2       string
	CurrencyType       string
	CodeOfTargetFund2  string
	SpecifyRateFee     string
	RateFee            string
	TransactionCfmDate string
	ReturnCode         string
	TaSerialNO         string
	ConfirmedVol       string
	ConfirmedAmount    string
	Nav                string
	PayAmount          string
}

//合约初始化
func (t *FundChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	t.createTable(stub)
	return nil, nil
}

//invoke基金数据
func (t *FundChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("------------start invoke chaincode------------")
	if len(args) != 1 {
		errors.New("Incorrect number of arguments. Expecting 1")
	}
	//获取传递的json数据
	arg := args[0]
	fmt.Println(arg)
	var data FundInfo
	var err = json.Unmarshal([]byte(arg), &data)
	if err != nil {
		return nil, err
	}

	if function == "update" {
		return t.update(stub, data)
	} else if function == "insert" {
		ok, err := stub.InsertRow("fund", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: data.Id}},
				&shim.Column{Value: &shim.Column_String_{String_: data.AppSheetSerialNo}},
				&shim.Column{Value: &shim.Column_String_{String_: data.FundCode2}},
				&shim.Column{Value: &shim.Column_String_{String_: data.TransactionDate}},
				&shim.Column{Value: &shim.Column_String_{String_: data.TransactionTime}},
				&shim.Column{Value: &shim.Column_String_{String_: data.DistributorCode}},
				&shim.Column{Value: &shim.Column_String_{String_: data.BusinessCode}},
				&shim.Column{Value: &shim.Column_String_{String_: data.ApplicationVol}},
				&shim.Column{Value: &shim.Column_String_{String_: data.ApplicationAmount}},
				&shim.Column{Value: &shim.Column_String_{String_: data.TaAccountID2}},
				&shim.Column{Value: &shim.Column_String_{String_: data.CurrencyType}},
				&shim.Column{Value: &shim.Column_String_{String_: data.CodeOfTargetFund2}},
				&shim.Column{Value: &shim.Column_String_{String_: data.SpecifyRateFee}},
				&shim.Column{Value: &shim.Column_String_{String_: data.RateFee}},

				&shim.Column{Value: &shim.Column_String_{String_: data.TransactionCfmDate}}, //交易确认日期
				&shim.Column{Value: &shim.Column_String_{String_: data.ReturnCode}},         //交易处理返回代码
				&shim.Column{Value: &shim.Column_String_{String_: data.TaSerialNO}},         //TA确认交易流水号
				&shim.Column{Value: &shim.Column_String_{String_: data.ConfirmedVol}},       //基金账户交易确认份数
				&shim.Column{Value: &shim.Column_String_{String_: data.ConfirmedAmount}},    //每笔交易确认金额
				&shim.Column{Value: &shim.Column_String_{String_: data.Nav}},                //基金单位净值
				&shim.Column{Value: &shim.Column_String_{String_: data.PayAmount}}},         //交收金额
		})
		fmt.Println("------------end invoke contract chaincode------------")
		if !ok {
			return nil, err
		}
	} else {
		return nil, errors.New("Incorrect method of request.")
	}

	return nil, nil
}

//查询基金表数据
func (t *FundChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("---------start query chaincode-----------")

	var result string

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	ids := args[0]
	arr := strings.Split(ids, ",")
	fmt.Println("pram data:" + ids)
	result = "["

	for i := 0; i < len(arr); i++ {
		var columns []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: arr[i]}}
		columns = append(columns, col1)
		row, err := stub.GetRow("fund", columns)
		if err != nil {
			return nil, err
		}
		jsonResp := `{"id":"` + row.Columns[0].GetString_() + `","appSheetSerialNo":"` + row.Columns[1].GetString_() + `","fundCode2":"` + row.Columns[2].GetString_() +
			`","transactionDate":"` + row.Columns[3].GetString_() + `","transactionTime":"` + row.Columns[4].GetString_() +
			`","distributorCode":"` + row.Columns[5].GetString_() + `","businessCode":"` + row.Columns[6].GetString_() +
			`","applicationVol":"` + row.Columns[7].GetString_() + `","applicationAmount":"` + row.Columns[8].GetString_() +
			`","taAccountID2":"` + row.Columns[9].GetString_() + `","currencyType":"` + row.Columns[10].GetString_() +
			`","codeOfTargetFund2":"` + row.Columns[11].GetString_() + `","specifyRateFee":"` + row.Columns[12].GetString_() +
			`","rateFee":"` + row.Columns[13].GetString_() + `","transactionCfmDate":"` + row.Columns[14].GetString_() +
			`","returnCode":"` + row.Columns[15].GetString_() + `","taSerialNO":"` + row.Columns[16].GetString_() +
			`","confirmedVol":"` + row.Columns[17].GetString_() + `","confirmedAmount":"` + row.Columns[18].GetString_() +
			`","nav":"` + row.Columns[19].GetString_() + `","payAmount":"` + row.Columns[20].GetString_() + `"}`

		result = result + jsonResp + ","
	}
	if len(result) > 1 {
		result = t.Substr(result, 0, len(result)-1) + "]"
	}
	return []byte([]byte(`{"status":"OK","data":` + result + `}`)), nil
}

//创建表
func (t *FundChaincode) createTable(stub shim.ChaincodeStubInterface) error {
	fmt.Println("------------start create table-------------------")
	err := stub.CreateTable("fund", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "appSheetSerialNo", Type: shim.ColumnDefinition_STRING, Key: false},   //申请单编号
		&shim.ColumnDefinition{Name: "fundCode2", Type: shim.ColumnDefinition_STRING, Key: false},          //基金代码
		&shim.ColumnDefinition{Name: "transactionDate", Type: shim.ColumnDefinition_STRING, Key: false},    //交易发生日期
		&shim.ColumnDefinition{Name: "transactionTime", Type: shim.ColumnDefinition_STRING, Key: false},    //交易发生时间
		&shim.ColumnDefinition{Name: "distributorCode", Type: shim.ColumnDefinition_STRING, Key: false},    //销售人代码
		&shim.ColumnDefinition{Name: "businessCode", Type: shim.ColumnDefinition_STRING, Key: false},       //业务代码
		&shim.ColumnDefinition{Name: "applicationVol", Type: shim.ColumnDefinition_STRING, Key: false},     //申请基金份数
		&shim.ColumnDefinition{Name: "applicationAmount", Type: shim.ColumnDefinition_STRING, Key: false},  //申请金额
		&shim.ColumnDefinition{Name: "taAccountID2", Type: shim.ColumnDefinition_STRING, Key: false},       //投资人基金帐号
		&shim.ColumnDefinition{Name: "currencyType", Type: shim.ColumnDefinition_STRING, Key: false},       //结算币种
		&shim.ColumnDefinition{Name: "codeOfTargetFund2", Type: shim.ColumnDefinition_STRING, Key: false},  //转换时的目标基金代码
		&shim.ColumnDefinition{Name: "specifyRateFee", Type: shim.ColumnDefinition_STRING, Key: false},     //代理费率
		&shim.ColumnDefinition{Name: "rateFee", Type: shim.ColumnDefinition_STRING, Key: false},            //总费率
		&shim.ColumnDefinition{Name: "transactionCfmDate", Type: shim.ColumnDefinition_STRING, Key: false}, //交易确认日期
		&shim.ColumnDefinition{Name: "returnCode", Type: shim.ColumnDefinition_STRING, Key: false},         //交易处理返回代码
		&shim.ColumnDefinition{Name: "taSerialNO", Type: shim.ColumnDefinition_STRING, Key: false},         //TA确认交易流水号
		&shim.ColumnDefinition{Name: "confirmedVol", Type: shim.ColumnDefinition_STRING, Key: false},       //基金账户交易确认份数
		&shim.ColumnDefinition{Name: "confirmedAmount", Type: shim.ColumnDefinition_STRING, Key: false},    //每笔交易确认金额
		&shim.ColumnDefinition{Name: "nav", Type: shim.ColumnDefinition_STRING, Key: false},                //基金单位净值
		&shim.ColumnDefinition{Name: "payAmount", Type: shim.ColumnDefinition_STRING, Key: false},          //交收金额
	})

	if err != nil {
		return errors.New("create table:stock error")
	}
	fmt.Println("------------end create table-------------------")
	return nil
}

//基金表数据更新
func (t *FundChaincode) update(stub shim.ChaincodeStubInterface, data FundInfo) ([]byte, error) {
	//根据ID查询 基金交易数据
	var columns []shim.Column
	col := shim.Column{Value: &shim.Column_String_{String_: data.Id}}
	columns1 := append(columns, col)
	row, err := stub.GetRow("fund", columns1) //row是否为空
	if err != nil {
		fmt.Println("---------------5------------------")
		return nil, errors.New("找不到数据")
	}
	appSheetSerialNo := row.Columns[1].GetString_()
	fundCode2 := row.Columns[2].GetString_()
	transactionDate := row.Columns[3].GetString_()
	transactionTime := row.Columns[4].GetString_()
	distributorCode := row.Columns[5].GetString_()
	businessCode := row.Columns[6].GetString_()
	applicationVol := row.Columns[7].GetString_()
	applicationAmount := row.Columns[8].GetString_()
	taAccountID2 := row.Columns[9].GetString_()
	currencyType := row.Columns[10].GetString_()
	codeOfTargetFund2 := row.Columns[11].GetString_()
	specifyRateFee := row.Columns[12].GetString_()
	rateFee := row.Columns[13].GetString_()

	ok, err := stub.ReplaceRow("fund", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: data.Id}},           //id
			&shim.Column{Value: &shim.Column_String_{String_: appSheetSerialNo}},  //申请单编号
			&shim.Column{Value: &shim.Column_String_{String_: fundCode2}},         //基金代码
			&shim.Column{Value: &shim.Column_String_{String_: transactionDate}},   //交易发生日期
			&shim.Column{Value: &shim.Column_String_{String_: transactionTime}},   //交易发生时间
			&shim.Column{Value: &shim.Column_String_{String_: distributorCode}},   //销售人代码
			&shim.Column{Value: &shim.Column_String_{String_: businessCode}},      //业务代码
			&shim.Column{Value: &shim.Column_String_{String_: applicationVol}},    //申请基金份数
			&shim.Column{Value: &shim.Column_String_{String_: applicationAmount}}, //申请金额
			&shim.Column{Value: &shim.Column_String_{String_: taAccountID2}},      //投资人基金帐号
			&shim.Column{Value: &shim.Column_String_{String_: currencyType}},      //结算币种
			&shim.Column{Value: &shim.Column_String_{String_: codeOfTargetFund2}}, //转换时费率
			&shim.Column{Value: &shim.Column_String_{String_: specifyRateFee}},    //代理费率
			&shim.Column{Value: &shim.Column_String_{String_: rateFee}},           //总费率

			&shim.Column{Value: &shim.Column_String_{String_: data.TransactionCfmDate}}, //交易确认日期
			&shim.Column{Value: &shim.Column_String_{String_: data.ReturnCode}},         //交易处理返回代码
			&shim.Column{Value: &shim.Column_String_{String_: data.TaSerialNO}},         //TA确认交易流水号
			&shim.Column{Value: &shim.Column_String_{String_: data.ConfirmedVol}},       //基金账户交易确认份数
			&shim.Column{Value: &shim.Column_String_{String_: data.ConfirmedAmount}},    //每笔交易确认金额
			&shim.Column{Value: &shim.Column_String_{String_: data.Nav}},                //基金单位净值
			&shim.Column{Value: &shim.Column_String_{String_: data.PayAmount}}},         //交收金额

	})
	fmt.Println("---------------8------------------")
	if !ok && err == nil {
		fmt.Println("---------------9------------------")
		return nil, errors.New("operation failed")
	}
	return nil, nil
}

//字符串截取
func (t *FundChaincode) Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
func main() {
	err := shim.Start(new(FundChaincode))
	if err != nil {
		fmt.Printf("Error starting Save State chaincode: %s", err)
	}
}
