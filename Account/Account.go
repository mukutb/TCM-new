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
	Securities string `json:"securities"`
}

type Securities struct{
	SecurityId string `json:"securityId"`
	AccountNumber string `json:"accountNumber"`
	SecuritiesName string `json:"securityName"`
	SecuritiesQuantity string `json:"securityQuantity"`
	SecurityType string `json:"securityType"`
	Category string `json:"category"`
	Totalvalue string `json:"totalvalue"`
	ValuePercentage string `json:"valuePercentage"`
	MTM string `json:"mtm"`
	EffectivePercentage string `json:"effectivePercentage"`
	EffectiveValueinUSD string `json:"effectiveValueinUSD"`
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
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
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
	fmt.Println("len(AccountIndex) : ")
	fmt.Println(len(AccountIndex))
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
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
	if len(args) != 7 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 7\", \"code\" : \"503\"}"
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
		res.Securities				=args[6]

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
		`"securities": "` + res.Securities + `" `+
		`}`
	err = stub.PutState(res.AccountName, []byte(order))									//store Account with id as key
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
	if len(args) != 7 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 7\", \"code\" : \"503\"}"
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
	securities				:=args[6]
		
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
	if res.AccountID == accountId{
		fmt.Println("This Account already exists: " + accountId)
		fmt.Println(res);
		errMsg := "{ \"message\" : \"This Account already exists\", \"code\" : \"503\"}"
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
		`"securities": "` + securities + `" `+
		`}`
		fmt.Println("order: " + order)
		fmt.Print("order in bytes array: ")
		fmt.Println([]byte(order))
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
	fmt.Println("! Account index after appending AccountId: ", AccountIndex)
	jsonAsBytes, _ := json.Marshal(AccountIndex)
	fmt.Print("jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(AccountIndexStr, jsonAsBytes)						//store name of Account
	if err != nil {
		return nil, err
	}

	tosend := "{ \"AccountId\" : \""+accountId+"\", \"message\" : \"Account created succcessfully\", \"code\" : \"200\"}"
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
	if len(args) != 10 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 10\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("start add_security")
	/*if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	*/

	_securityId				:= args[0]
	_accountNumber 				:= args[1]
		
	AccountAsBytes, err := stub.GetState(_securityId)
		if err != nil {
			return nil, errors.New("Failed to get Security " + _securityId)
		}
	res := Securities{}
	json.Unmarshal(AccountAsBytes, &res)
	if res.SecurityId == _securityId{
		errMsg := "{ \"message\" : \"This Security already exists\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
	return nil, nil				//all stop a Account by this name exists
	}
	
	//build the Account json string manually
	order := 	`{`+

		`"securityId": "` + args[0] + `" ,`+
		`"accountNumber": "` + args[1] + `" ,`+
		`"securityName": "` + args[2] + `" ,`+
		`"securityQuantity": "` + args[3] + `" ,`+
		`"securityType": "` + args[4] + `" ,`+
		`"category": "` + args[5] + `" ,`+
		`"totalvalue": "` + args[6] + `" ,`+
		`"valuePercentage": "` + args[7] + `" ,`+
		`"mtm": "` + args[8] + `" ,`+
		`"effectivePercentage": "` + args[9] + `" `+
		`}`
		fmt.Println("order: " + order)
		fmt.Print("order in bytes array: ")
		fmt.Println([]byte(order))
	err = stub.PutState(_securityId, []byte(order))									//store Account with AccountId as key
	if err != nil {
		return nil, err
	}

	//Adding Security to the account
	res2 := Accounts{}
	json.Unmarshal(AccountAsBytes, &res2)
	if res2.AccountNumber == _accountNumber{
		fmt.Println("Account found with AccountNumber : " + _accountNumber)
		fmt.Println(res2);

	}else{
		errMsg := "{ \"message\" : \""+ _accountNumber+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	res2.Securities = res2.Securities+ "," + _securityId
	
	order2 := 	`{`+
		`"accountId": "` + res2.AccountID + `" ,`+
		`"accountName": "` + res2.AccountName + `" ,`+
		`"accountNumber": "` + res2.AccountNumber + `" ,`+
		`"accountType": "` + res2.AccountType + `" ,`+
		`"totalValue": "` + res2.TotalValue + `" ,`+
		`"currency": "` + res2.Currency + `" ,`+
		`"securities": "` + res2.Securities + `" `+
		`}`
	err = stub.PutState(res2.AccountNumber, []byte(order2))									//store Account with id as key
	if err != nil {
		return nil, err
	}

	tosend := "{ \"SecurityId\" : \""+_securityId+"\", \"message\" : \"Security updated succcesfully\", \"code\" : \"200\"}"
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

	_accountNumber				:=args[0]
	
		
	AccountAsBytes, err := stub.GetState(_accountNumber)
		if err != nil {
			return nil, errors.New("Failed to get Account " + _accountNumber)
		}
	res := Accounts{}
	json.Unmarshal(AccountAsBytes, &res)

	
	//build the Account json string manually
	order := 	`{`+
		`"accountId": "` + res.AccountID + `" ,`+
		`"accountName": "` + res.AccountName + `" ,`+
		`"accountNumber": "` + res.AccountNumber + `" ,`+
		`"accountType": "` + res.AccountType + `" ,`+
		`"totalValue": "` + res.TotalValue + `" ,`+
		`"currency": "` + res.Currency + `" ,`+
		`"securities": `+ "" +`" `+
		`}`
		fmt.Println("order: " + order)
		fmt.Print("order in bytes array: ")
		fmt.Println([]byte(order))
	err = stub.PutState(_accountNumber, []byte(order))									//store Account with _accountNumber as key
	if err != nil {
		return nil, err
	}
	tosend := "{ \"AccountNumber\" : \""+_accountNumber+"\", \"message\" : \"Securities deleted succcesfully!\", \"code\" : \"200\"}"
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
	_tempJson :=Accounts{}

	var res = Accounts{}
	AccountAsBytes, err := stub.GetState(_AccountNumber)
	if err != nil {
		return nil, errors.New("Failed to get Account index")
	}
	json.Unmarshal(AccountAsBytes, &res)
	_SecuritySplit := strings.Split(res.Securities, ",")

	jsonResp = "{"
	for i,val := range _SecuritySplit{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Account")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		json.Unmarshal(valueAsBytes, &_tempJson)
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		if _tempJson.AccountType == _AccountNumber {
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			if i < len(_SecuritySplit)-1 {
				jsonResp = jsonResp + ","
			}
		}
	}
	jsonResp = jsonResp + "}"
	fmt.Println("end getSecurities_byAccount")
	return []byte(jsonResp), nil
}
