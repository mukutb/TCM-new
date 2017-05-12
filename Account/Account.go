/*/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
"errors"
"fmt"
"strconv"
"encoding/json"
"strings"
"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ManageAccounts example simple Chaincode implementation
type ManageAccounts struct {
}

var AccountIndexStr = "_AccountIndex"				//name for the key/value that will store a list of all known RQv's
var SecurityIndexStr = "_SecurityIndex"

type Accounts struct{
	AccountID string `json:"accountId"`
	AccountName string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	AccountType string `json:"accountType"`
	TotalValue string `json:"totalValue"`
	Currency string `json:"currency"`
	Pledger string `json:"pledger"`
	Securities string `json:"securities"`
}

type Securities struct{
	SecurityId string `json:"securityId"`
	AccountNumber string `json:"accountNumber"`
	SecurityName string `json:"securityName"`
	SecurityQuantity string `json:"securityQuantity"`
	SecurityType string `json:"securityType"`
	CollateralForm string `json:"collateralForm"`
	Totalvalue string `json:"totalvalue"`
	ValuePercentage string `json:"valuePercentage"`
	MTM string `json:"mtm"`
	EffectivePercentage string `json:"effectivePercentage"`
	EffectiveValueinUSD string `json:"effectiveValueinUSD"`
	Currency string `json:"currency"`
}
// ============================================================================================================================
// Main - start the chaincode for Account management
// ============================================================================================================================
func main() {			
	err := shim.Start(new(ManageAccounts))
	if err != nil {
		fmt.Printf("Error starting Account management chaincode: %s", err)
	}
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManageAccounts) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var msg string
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting ' ' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// Initialize the chaincode
	msg = args[0]
	// Write the state to the ledger
	err = stub.PutState("abc", []byte(msg))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(AccountIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	tosend := "{ \"message\" : \"ManageAccounts chaincode is deployed successfully.\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// Run - Our entry Accountint for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *ManageAccounts) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}
// ============================================================================================================================
// Invoke - Our entry Accountint for Invocations
// ============================================================================================================================
func (t *ManageAccounts) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "create_account" {											//create a new Account
		return t.create_Account(stub, args)
	}else if function == "update_account" {									
		return t.update_Account(stub, args)
	}else if function == "add_security" {									
		return t.add_security(stub, args)
	}else if function == "remove_securitiesFromAccount" {									
		return t.remove_securitiesFromAccount(stub, args)
	}else if function == "update_security" {									
		return t.update_security(stub, args)
	}else if function == "delete_security" {									
		return t.delete_security(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	errMsg := "{ \"message\" : \"Received unknown function invocation\", \"code\" : \"503\"}"
	err := stub.SetEvent("errEvent", []byte(errMsg))
	if err != nil {
		return nil, err
	} 
	return nil, nil			//error
}
// ============================================================================================================================
// Query - Our entry Accountint for Queries
// ============================================================================================================================
func (t *ManageAccounts) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getAccount_byName" {													//Read a Account by name
		return t.getAccount_byName(stub, args)
	} else if function == "getAccount_byType" {													//Read a Account by Type
		return t.getAccount_byType(stub, args)
	} else if function == "getAccount_byNumber" {													//Read a Account by Number
		return t.getAccount_byNumber(stub, args)
	} else if function == "get_AllAccount" {													//Read all Accounts
		return t.get_AllAccount(stub, args)
	}else if function == "getSecurities_byAccount" {									//update a Account
		return t.getSecurities_byAccount(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//errors
	errMsg := "{ \"message\" : \"Received unknown function query\", \"code\" : \"503\"}"
	err := stub.SetEvent("errEvent", []byte(errMsg))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
//  getAccount_byName- get details of all Account from chaincode state
// ============================================================================================================================
func (t *ManageAccounts) getAccount_byName(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var AccountIndex []string
	fmt.Println("start getAccount_byName")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 1\" \" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	_AccountName := args[0]
	var _tempJson Accounts

	AccountAsBytes, err := stub.GetState(AccountIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Account index")
	}
	json.Unmarshal(AccountAsBytes, &AccountIndex)								//un stringify it aka JSON.parse()
	jsonResp = "{"
	for i,val := range AccountIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Account")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		json.Unmarshal(valueAsBytes, &_tempJson)
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		if _tempJson.AccountName == _AccountName {
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			if i < len(AccountIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}
	}
	fmt.Println("len(AccountIndex) : ")
	fmt.Println(len(AccountIndex))
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	if jsonResp == "{}" {
        fmt.Println("Account not found for  " + _AccountName)
        jsonResp =  "{\"AccountName\" : \"" + _AccountName + "\", \"message\" : \"Account not found.\", \"code\" : \"503\"}"
        errMsg:= "{ \"AccountName\" : \"" + _AccountName + "\", \"message\" : \"Account not found.\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
        	return nil, err
        }
    }
    if strings.Contains(jsonResp,"},}"){
    	jsonResp = strings.Replace(jsonResp, "},}", "}}", -1)
    }
	fmt.Println("end getAccount_byName")
	return []byte(jsonResp), nil
											//send it onward
}
// ============================================================================================================================
//  getAccount_byType- get details of all Account from chaincode state
// ============================================================================================================================
func (t *ManageAccounts) getAccount_byType(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var AccountIndex []string
	fmt.Println("start getAccount_byType")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 1\" \" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	_AccountType := args[0]
	_tempJson :=Accounts{}

	AccountAsBytes, err := stub.GetState(AccountIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Account index")
	}
	fmt.Print("AccountAsBytes : ")
	fmt.Println(AccountAsBytes)
	json.Unmarshal(AccountAsBytes, &AccountIndex)								//un stringify it aka JSON.parse()
	fmt.Print("AccountIndex : ")
	fmt.Println(AccountIndex)
	fmt.Println("len(AccountIndex) : ")
	fmt.Println(len(AccountIndex))
	jsonResp = "{"
	for i,val := range AccountIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Account")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		json.Unmarshal(valueAsBytes, &_tempJson)
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		if _tempJson.AccountType == _AccountType {
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			if i < len(AccountIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}
	}
	jsonResp = jsonResp + "}"
	if jsonResp == "{}" {
        fmt.Println(_AccountType + " account not found")
        jsonResp = "{ \"AccountType\" : \"" + _AccountType + "\", \"message\" : \"Account not found.\", \"code\" : \"503\"}"
        errMsg:= "{ \"AccountType\" : \"" + _AccountType + "\", \"message\" : \"Account Not Found.\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
    	    return nil, err
        }
    }
    if strings.Contains(jsonResp,"},}"){
    	jsonResp = strings.Replace(jsonResp, "},}", "}}", -1)
    }
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Println("end getAccount_byType")
	return []byte(jsonResp), nil
											//send it onward
}
// ============================================================================================================================
//  getAccount_byNumber- get details of all Account from chaincode state
// ============================================================================================================================
func (t *ManageAccounts) getAccount_byNumber(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	fmt.Println("start getAccount_byNumber")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 1\" \" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	_AccountNumber := args[0]
	_tempJson :=Accounts{}

	jsonResp = "{"
	fmt.Println(_AccountNumber + " - looking at " + _AccountNumber + " for Account")
	valueAsBytes, err := stub.GetState(_AccountNumber)
	if err != nil {
		errResp = "{\"Error\":\"Failed to get state for " + _AccountNumber + "\"}"
		return nil, errors.New(errResp)
	}
	json.Unmarshal(valueAsBytes, &_tempJson)
	fmt.Print("valueAsBytes : ")
	fmt.Println(valueAsBytes)
	if _tempJson.AccountNumber == _AccountNumber {
		jsonResp = jsonResp + "\""+ _AccountNumber + "\":" + string(valueAsBytes[:])
	}else{
        fmt.Println(_AccountNumber + " not found")
        jsonResp = jsonResp + "\"AccountNumber\" : \"" + _AccountNumber + "\",\"message\" : \"" + "Account not found.\", \"code\" : \"503\""
        errMsg:= "{ \"AccountNumber\" : \"" + _AccountNumber + "\",\"message\" : \"" + "Account Not Found.\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
    }
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getAccount_byNumber")
	return []byte(jsonResp), nil
											//send it onward
}
// ============================================================================================================================
//  get_AllAccount- get details of all Account from chaincode state
// ============================================================================================================================
func (t *ManageAccounts) get_AllAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var AccountIndex []string
	fmt.Println("start get_AllAccount")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 1\" \" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	AccountAsBytes, err := stub.GetState(AccountIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Account index")
	}
	fmt.Print("AccountAsBytes : ")
	fmt.Println(AccountAsBytes)
	json.Unmarshal(AccountAsBytes, &AccountIndex)								//un stringify it aka JSON.parse()
	fmt.Print("AccountIndex : ")
	fmt.Println(AccountIndex)
	jsonResp = "{"
	for i,val := range AccountIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Account")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i < len(AccountIndex)-1 {
			jsonResp = jsonResp + ","
		}
	}
	fmt.Println("len(AccountIndex) : ")
	fmt.Println(len(AccountIndex))
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end get_AllAccount")
	return []byte(jsonResp), nil
											//send it onward
}
// ============================================================================================================================
// update_Account - update Account into chaincode state
// ============================================================================================================================
func (t *ManageAccounts) update_Account(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating Account")
	if len(args) != 8 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 8\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set accountNumber
	accountNumber := args[2]
	AccountAsBytes, err := stub.GetState(accountNumber)									//get the Account for the specified AccountId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to get state for " + accountNumber + "\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Print("AccountAsBytes in update Account")
	fmt.Println(AccountAsBytes);
	res := Accounts{}
	json.Unmarshal(AccountAsBytes, &res)
	if res.AccountNumber == accountNumber{
		fmt.Println("Account found with AccountNumber : " + accountNumber)
		fmt.Println(res);
		res.AccountID				=args[0]
		res.AccountName				=args[1]
		res.AccountNumber			=args[2]
		res.AccountType				=args[3]
		res.TotalValue				=args[4]
		res.Currency				=args[5]
		res.Pledger				    =args[6]
		res.Securities				=args[7]

	}else{
		errMsg := "{ \"message\" : \""+ accountNumber+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	//build the Account json string manually
	order := 	`{`+
		`"accountId": "` + res.AccountID + `" ,`+
		`"accountName": "` + res.AccountName + `" ,`+
		`"accountNumber": "` + res.AccountNumber + `" ,`+
		`"accountType": "` + res.AccountType + `" ,`+
		`"totalValue": "` + res.TotalValue + `" ,`+
		`"currency": "` + res.Currency + `" ,`+
		`"pledger": "` + res.Pledger + `" ,`+
		`"securities": "` + res.Securities + `" `+
		`}`
	err = stub.PutState(res.AccountNumber, []byte(order))									//store Account with id as key
	if err != nil {
		return nil, err
	}

	tosend := "{ \"AccountNumber\" : \"" + accountNumber + "\", \"message\" : \"Account updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("Account updated succcessfully")
	return nil, nil
}
// ============================================================================================================================
// create Account - create a new Account, store into chaincode state
// ============================================================================================================================
func (t *ManageAccounts) create_Account(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 8 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 8\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("start create_Account")

	accountId				:=args[0]
	accountName				:=args[1] 
	accountNumber			:=args[2]
	accountType				:=args[3]
	totalValue				:=args[4]
	currency				:=args[5]
	pledger 				:=args[6]
	securities				:=args[7]
	
	AccountAsBytes, err := stub.GetState(accountNumber)
	if err != nil {
		return nil, errors.New("Failed to get Account " + accountNumber)
	}
	fmt.Print("AccountAsBytes: ")
	fmt.Println(AccountAsBytes)
	res := Accounts{}
	json.Unmarshal(AccountAsBytes, &res)
	fmt.Print("res: ")
	fmt.Println(res)
	if res.AccountNumber == accountNumber{
		fmt.Println("This Account already exists: " + accountNumber)
		fmt.Println(res);
		errMsg := "{ \"AccountNumber\" : \""+accountNumber+"\", \"message\" : \"This account already exists\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil				//all stop a Account by this name exists
	}
	
	//build the Account json string manually
	order := 	`{`+
		`"accountId": "` + accountId + `" ,`+
		`"accountName": "` + accountName + `" ,`+
		`"accountNumber": "` + accountNumber + `" ,`+
		`"accountType": "` + accountType + `" ,`+
		`"totalValue": "` + totalValue + `" ,`+
		`"currency": "` + currency + `" ,`+
		`"pledger": "` + pledger + `" ,`+
		`"securities": "` + securities + `" `+
		`}`
	fmt.Println("order: " + order)
	err = stub.PutState(accountNumber, []byte(order))									//store Account with AccountId as key
	if err != nil {
		return nil, err
	}
	//get the Account index
	AccountIndexAsBytes, err := stub.GetState(AccountIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Account index")
	}
	var AccountIndex []string
	fmt.Print("AccountIndexAsBytes: ")
	fmt.Println(AccountIndexAsBytes)
	
	json.Unmarshal(AccountIndexAsBytes, &AccountIndex)							//un stringify it aka JSON.parse()
	fmt.Print("AccountIndex after unmarshal..before append: ")
	fmt.Println(AccountIndex)
	//append
	AccountIndex = append(AccountIndex, accountNumber)									//add Account AccountId to index list
	fmt.Println("! Account index after appending accountNumber: ", AccountIndex)
	jsonAsBytes, _ := json.Marshal(AccountIndex)
	fmt.Print("jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(AccountIndexStr, jsonAsBytes)						//store name of Account
	if err != nil {
		return nil, err
	}

	tosend := "{ \"AccountNumber\" : \""+accountNumber+"\", \"message\" : \"Account created succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("end create_Account")
	return nil, nil
}
// ============================================================================================================================
// Add Securities - Update Securities for an account, store into chaincode state
// ============================================================================================================================
func (t *ManageAccounts) add_security(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) !=  12{
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 12\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("start add_security")
	
	_securityId				:= args[0]
	_accountNumber 			:= args[1]
	_securityName			:= args[2]
	_securityQuantity		:= args[3]
	_securityType			:= args[4]
	_collateralForm			:= args[5]
	_totalValue			    := args[6]
	_valuePercentage		:= args[7]
	_mtm					:= args[8]
	_effectivePercentage	:= args[9]
	_effectiveValueinUSD	:= args[10];
	_currency			    := args[11]
	

	SecurityAsBytes, err := stub.GetState(_accountNumber+"-"+_securityId)
		if err != nil {
			return nil, errors.New("Failed to get Security " + _accountNumber+"-"+_securityId)
		}
	res := Securities{}
	json.Unmarshal(SecurityAsBytes, &res)

	// NOTE:: This is not required as Securities can be added, hence remove check for already existing
	/*if res.SecurityId == _securityId{
		errMsg := "{ \"SecurityId\" : \""+_securityId+"\",\"message\" : \"This Security already exists\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil				//all stop a Account by this name exists
	}*/
	
	//build the Account json string manually
	order := 	`{`+
		`"securityId": "` + _securityId + `" ,`+
		`"accountNumber": "` + _accountNumber + `" ,`+
		`"securityName": "` + _securityName + `" ,`+
		`"securityQuantity": "` + _securityQuantity + `" ,`+
		`"securityType": "` + _securityType + `" ,`+
		`"collateralForm": "` + _collateralForm + `" ,`+
		`"totalvalue": "` + _totalValue + `" ,`+
		`"valuePercentage": "` + _valuePercentage + `" ,`+
		`"mtm": "` + _mtm + `" ,`+
		`"effectivePercentage": "` + _effectivePercentage + `" ,`+
		`"effectiveValueinUSD": "` + _effectiveValueinUSD + `" ,`+
		`"currency": "` + _currency + `"`+
		`}`
	fmt.Println("order: " + order)
	err = stub.PutState(_accountNumber+"-"+_securityId, []byte(order))									//store Account with AccountId as key
	if err != nil {
		return nil, err
	}
	AccountAsBytes, err := stub.GetState(_accountNumber)
	if err != nil {
		return nil, errors.New("Failed to get account " + _accountNumber)
	}
	//Adding Security to the account
	res2 := Accounts{}
	json.Unmarshal(AccountAsBytes, &res2)
	fmt.Println(res2);
	if res2.AccountNumber == _accountNumber{
		fmt.Println("Account found with AccountNumber : " + _accountNumber)
		_SecuritySplit := strings.Split(res2.Securities, ",")
		fmt.Print("_SecuritySplit: " )
		fmt.Println(_SecuritySplit)
		for i := range _SecuritySplit{
			fmt.Println("_SecuritySplit[i]: " + _SecuritySplit[i])
			if _SecuritySplit[i] == _accountNumber+"-"+_securityId {
				errMsg := "{ \"SecurityId\" : \""+_accountNumber+"-"+_securityId+"\",\"message\" : \" SecurityId already exists in the account.\", \"code\" : \"503\"}"
				err = stub.SetEvent("errEvent", []byte(errMsg))
				if err != nil {
					return nil, err
				} 
				return nil, nil
			}
		}
	}else{
		errMsg := "{ \"message\" : \""+ _accountNumber+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// Convert account's totalValue(String) to float
	tempTotalValue1, errBool := strconv.ParseFloat(res2.TotalValue, 32)
	if errBool != nil {
		fmt.Println(errBool)
	}
	// Convert security's totalvalue(String) to float
	tempTotalvalue2, errBool := strconv.ParseFloat(_totalValue, 32)
	if errBool != nil {
		fmt.Println(errBool)
	}
	if res2.Securities == " " || res2.Securities == "" {
		res2.Securities = _accountNumber+"-"+_securityId;
		_tempTotal := tempTotalValue1 + tempTotalvalue2
		res2.TotalValue = strconv.FormatFloat(_tempTotal, 'f', -1, 64)
	}else {
		res2.Securities = res2.Securities+ "," + _accountNumber+"-"+_securityId;
		_tempTotal := tempTotalValue1 + tempTotalvalue2
		res2.TotalValue = strconv.FormatFloat(_tempTotal, 'f', -1, 64)
	}
	order2 := 	`{`+
		`"accountId": "` + res2.AccountID + `" ,`+
		`"accountName": "` + res2.AccountName + `" ,`+
		`"accountNumber": "` + res2.AccountNumber + `" ,`+
		`"accountType": "` + res2.AccountType + `" ,`+
		`"totalValue": "` + res2.TotalValue + `" ,`+
		`"currency": "` + res2.Currency + `" ,`+
		`"pledger": "` + res2.Pledger + `" ,`+
		`"securities": "` + res2.Securities + `" `+
		`}`
	fmt.Println("order2: " + order2)
	err = stub.PutState(res2.AccountNumber, []byte(order2))									//store Account with id as key
	if err != nil {
		return nil, err
	}
	tosend := "{ \"SecurityId\" : \""+_accountNumber+"-"+_securityId+"\", \"message\" : \"Security updated succcesfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("end add_security")
	return nil, nil
}
// ============================================================================================================================
// Remove Securities - Update Securities for an account, store into chaincode state
// ============================================================================================================================
func (t *ManageAccounts) remove_securitiesFromAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 1\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("start remove_securitiesFromAccount")

	_accountNumber	:=args[0]
		
	AccountAsBytes, err := stub.GetState(_accountNumber)
	if err != nil {
		return nil, errors.New("Failed to get Account " + _accountNumber)
	}
	totalValueOfTheDeletedSecurities := float64(0.0)
	res := Accounts{}
	res_Security := Securities{}
	json.Unmarshal(AccountAsBytes, &res)
	_SecuritySplit := strings.Split(res.Securities, ",")
	fmt.Print("_SecuritySplit: " )
	fmt.Println(_SecuritySplit)
	for i:=0;i<len(_SecuritySplit);{
		fmt.Println("_SecuritySplit[i]: " + _SecuritySplit[i])

		//Get the total value of the securities
		SecuritiesAsBytes, err := stub.GetState(_SecuritySplit[i])
		if err != nil {
			return nil, errors.New("Failed to get Security " + _SecuritySplit[i])
		}
		json.Unmarshal(SecuritiesAsBytes, &res_Security)
		valToBeAdded, _ := strconv.ParseFloat(res_Security.Totalvalue, 64)
		totalValueOfTheDeletedSecurities = totalValueOfTheDeletedSecurities + valToBeAdded

		//Got the info. now delete
		err = stub.DelState(_SecuritySplit[i])													//remove the key from chaincode state
		if err != nil {
			errMsg := "{ \"security\" : \"" + _SecuritySplit[i] + "\", \"message\" : \"Failed to delete state\", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		}
		_SecuritySplit = append(_SecuritySplit[:i], _SecuritySplit[i+1:]...)			//remove it
		//fmt.Println(_SecuritySplit[:i])
		//fmt.Println(_SecuritySplit[i+1])
		fmt.Println(_SecuritySplit)
		for x:= range _SecuritySplit{											//debug prints...
			fmt.Println(string(x) + " - " + _SecuritySplit[x])
		}
	}

	res.Securities = strings.Join(_SecuritySplit,",");
	fmt.Println(res.Securities)
	fmt.Println(_SecuritySplit)
	fmt.Println("totalValueOfTheDeletedSecurities::")
	fmt.Println(totalValueOfTheDeletedSecurities)
	
	//build the Account json string manually
	account_json := 	`{`+
		`"accountId": "` + res.AccountID + `" ,`+
		`"accountName": "` + res.AccountName + `" ,`+
		`"accountNumber": "` + res.AccountNumber + `" ,`+
		`"accountType": "` + res.AccountType + `" ,`+
		`"totalValue": "` + strconv.FormatFloat(totalValueOfTheDeletedSecurities, 'f', 2, 64) + `" ,`+
		`"currency": "` + res.Currency + `" ,`+
		`"pledger": "` + res.Pledger + `" ,`+
		`"securities": "`+ res.Securities +`" `+
		`}`
	fmt.Println("account_json: " + account_json)
	err = stub.PutState(_accountNumber, []byte(account_json))									//store Account with _accountNumber as key
	if err != nil {
		return nil, err
	}

	tosend := "{ \"AccountNumber\" : \""+_accountNumber+"\", \"message\" : \"All securities deleted succcesfully!\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	}
	fmt.Println("end remove_securitiesFromAccount")
	return nil, nil
}
// ============================================================================================================================
//  getSecurities_byAccount- get details of all Account from chaincode state
// ============================================================================================================================
func (t *ManageAccounts) getSecurities_byAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	fmt.Println("start getSecurities_byAccount")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"AccountNumber\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	_AccountNumber := args[0]
	_tempJson :=Securities{}

	var res = Accounts{}
	AccountAsBytes, err := stub.GetState(_AccountNumber)
	if err != nil {
		return nil, errors.New("Failed to get Account index")
	}
	json.Unmarshal(AccountAsBytes, &res)
	fmt.Print("account details: ");
	fmt.Println(res)
	_SecuritySplit := strings.Split(res.Securities, ",")
	fmt.Print("_SecuritySplit: " )
	fmt.Println(_SecuritySplit)
	jsonResp = "["
	for i := range _SecuritySplit{
		fmt.Println("_SecuritySplit[i]: " + _SecuritySplit[i])
		valueAsBytes, err := stub.GetState(_SecuritySplit[i])
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + _SecuritySplit[i] + "\"}"
			return nil, errors.New(errResp)
		}
		json.Unmarshal(valueAsBytes, &_tempJson)
		fmt.Print("_tempJson : ")
		fmt.Println(_tempJson)
		jsonResp = jsonResp + string(valueAsBytes[:])
		if i < len(_SecuritySplit)-1 {
			jsonResp = jsonResp + ","
		}
	}
	jsonResp = jsonResp + "]"
	fmt.Print("jsonResp: ")
	fmt.Println(jsonResp)
	if jsonResp == "[\"\":]" || jsonResp == "[\" \":]" {
		jsonResp = "{ \"AccountNumber\" : \"" + _AccountNumber + "\", \"message\" : \"No securities found.\", \"code\" : \"503\"}"
	}
	fmt.Println("end getSecurities_byAccount")
	return []byte(jsonResp), nil
}
// ============================================================================================================================
// update_security - update Security into chaincode state
// ============================================================================================================================
func (t *ManageAccounts) update_security(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating Security")
	if len(args) != 12 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 12\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set accountNumber
	securityId := args[0]
	accountNumber := args[1]
	securityAsBytes, err := stub.GetState(accountNumber + "-" + securityId)									//get the Security for the specified accountNumber-securityId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to get state for " + accountNumber + "-" + securityId + "\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Print("securityAsBytes in update Security")
	fmt.Println(securityAsBytes);
	res := Securities{}
	json.Unmarshal(securityAsBytes, &res)
	if res.SecurityId == securityId{
		fmt.Println("Security found with SecurityId : " + securityId)
		fmt.Println(res);
		//build the Account json string manually
		order := 	`{`+
			`"securityId": "` + res.SecurityId + `" ,`+
			`"accountNumber": "` + res.AccountNumber + `" ,`+
			`"securityName": "` + args[2] + `" ,`+
			`"securityQuantity": "` + args[3] + `" ,`+
			`"securityType": "` + args[4] + `" ,`+
			`"collateralForm": "` + args[5] + `" ,`+
			`"totalvalue": "` + args[6]+ `" ,`+
			`"valuePercentage": "` + args[7]+ `" ,`+
			`"mtm": "` + args[8]+ `" ,`+
			`"effectivePercentage": "` + args[9]+ `" ,`+
			`"effectiveValueinUSD": "` + args[10]+ `" ,`+
			`"currency": "` + args[11] + `"`+
			`}`
		fmt.Println(order);
		err = stub.PutState(accountNumber + "-" + securityId, []byte(order))									//store security with id as key
		if err != nil {
			return nil, err
		}
		fmt.Println("Security updated succcessfully")
	}else{
		errMsg := "{ \"message\" : \""+ securityId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	tosend := "{ \"Security\" : \"" + accountNumber + "-" + securityId + "\", \"message\" : \"Security updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	return nil, nil
}
// ============================================================================================================================
// Delete - remove a Security from state and then remove from account
// ============================================================================================================================
func (t *ManageAccounts) delete_security(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"accountNumber,securityId\" arguments.\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set security
	_securityId := args[0];
	_accountNumber := args[1];
	security := _accountNumber + "-" + _securityId;
	fmt.Println(security);
	err := stub.DelState(security)													//remove the key from chaincode state
	if err != nil {
		errMsg := "{ \"security\" : \"" + security + "\", \"message\" : \"Failed to delete state\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	//get the account Number details
	accountAsBytes, err := stub.GetState(_accountNumber)
	if err != nil {
		return nil, errors.New("Failed to get Account number")
	}
	valIndex := Accounts{}
	json.Unmarshal(accountAsBytes, &valIndex)	
	_SecuritySplit := strings.Split(valIndex.Securities, ",")
	fmt.Print("_SecuritySplit: " )
	fmt.Println(_SecuritySplit)
	for i := range _SecuritySplit{
		fmt.Println("_SecuritySplit[i]: " + _SecuritySplit[i])
		fmt.Println(_SecuritySplit[i] == (_accountNumber+"-"+_securityId))
		if _SecuritySplit[i] == (_accountNumber+"-"+_securityId) {
			fmt.Println("Security Found.");
			_SecuritySplit = append(_SecuritySplit[:i], _SecuritySplit[i+1:]...)			//remove it
			fmt.Println(_SecuritySplit[:i])
			fmt.Println(_SecuritySplit)
			for x:= range _SecuritySplit{											//debug prints...
				fmt.Println(string(x) + " - " + _SecuritySplit[x])
			}
			break
		}
	}
	fmt.Println(_SecuritySplit);
	valIndex.Securities = strings.Join(_SecuritySplit,",");
	fmt.Println(_SecuritySplit);
	//build the Account json string manually
	order := 	`{`+
		`"accountId": "` + valIndex.AccountID + `" ,`+
		`"accountName": "` + valIndex.AccountName + `" ,`+
		`"accountNumber": "` + valIndex.AccountNumber + `" ,`+
		`"accountType": "` + valIndex.AccountType + `" ,`+
		`"totalValue": "` + valIndex.TotalValue + `" ,`+
		`"currency": "` + valIndex.Currency + `" ,`+
		`"pledger": "` + valIndex.Pledger + `" ,`+
		`"securities": "`+ valIndex.Securities +`" `+
		`}`
		
	fmt.Println("order: " + order)
	err = stub.PutState(_accountNumber, []byte(order))									//store Account with _accountNumber as key
	if err != nil {
		return nil, err
	}
	tosend := "{ \"security\" : \""+security+"\", \"message\" : \"Security deleted succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
