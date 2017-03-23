/*/
*
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
import ("errors"
        "fmt"
        "strconv"
        "strings"
        "encoding/json"
        "github.com/hyperledger/fabric/core/chaincode/shim")

type ManageDeals struct {}

var DealIndexStr = "_Dealindex" //name for the key/value that will store a list of all known RQv's

type Transactions struct {
    TransactionId string `json:"transactionId"`
    TransactionDate string `json:"transactionDate"`
    DealID string `json:"dealId"`
    Pledger string `json:"pledger"`
    Pledgee string `json:"pledgee"`
    RQV string `json:"rqv"`
    Currency string `json:"currency"`
    CurrencyConversionRate string `json:"currencyConversionRate"`
    MarginCAllDate string `json:"marginCAllDate"`
    AllocationStatus string `json:"allocationStatus"`
    TransactionStatus string `json:"transactionStatus"`
}

type Deals struct { // Attributes of a Deal
    DealID string `json:"dealId"`
    Pledger string `json:"pledger"`
    Pledgee string `json:"pledgee"`
    MaxValue string `json:"maxValue"` //Maximum Value of all the securities of each Collateral Form 
    TotalValueLongBoxAccount string `json:"totalValueLongBoxAccount"`
    TotalValueSegregatedAccount string `json:"totalValueSegregatedAccount"`
    IssueDate string `json:"issueDate"`
    LastSuccessfulAllocationDate string `json:"lastSuccessfulAllocationDate"`
    Transactions string `json:"transactions"`
}
/*

type Pledger struct{
    PleadgerID string `json:"pledgerId"`
    PledgerName string `json:"PledgerName"`
    LongAccountNumber string `json:"longAccountNumber"`
    SegregatedAccountNumbers []string `json:"segregatedAccountNumbers"`
}*/
// ============================================================================================================================
// Main - start the chaincode for Deal management
// ============================================================================================================================
func main() {
    err:= shim.Start(new(ManageDeals))
    if err != nil {
        fmt.Printf("Error starting Deal management chaincode: %s", err)
    }
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func(t * ManageDeals) Init(stub shim.ChaincodeStubInterface, function string, args[] string)([] byte, error) {
    var msg string
    var err error
    if len(args) != 1 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting ' ' as an argument\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    // Initialize the chaincode
    msg = args[0]
    // Write the state to the ledger
    err = stub.PutState("abc", [] byte(msg)) //making a test var "abc", I find it handy to read/write to it right away to test the network
    if err != nil {
        return nil, err
    }
    var empty[] string
    jsonAsBytes, _:= json.Marshal(empty) //marshal an emtpy array of strings to clear the index
    err = stub.PutState(DealIndexStr, jsonAsBytes)
    if err != nil {
        return nil, err
    }
    tosend:= "{ \"message\" : \"ManageDeals chaincode is deployed successfully.\", \"code\" : \"200\"}"
    err = stub.SetEvent("evtsender", [] byte(tosend))
    if err != nil {
        return nil, err
    }
    return nil, nil
}
// ============================================================================================================================
// Run - Our entry Dealint for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func(t * ManageDeals) Run(stub shim.ChaincodeStubInterface, function string, args[] string)([] byte, error) {
    fmt.Println("run is running " + function)
    return t.Invoke(stub, function, args)
}
// ============================================================================================================================
// Invoke - Our entry Dealint for Invocations
// ============================================================================================================================
func(t * ManageDeals) Invoke(stub shim.ChaincodeStubInterface, function string, args[] string)([] byte, error) {
    fmt.Println("invoke is running " + function)
    // Handle different functions
    if
    function == "init" { //initialize the chaincode state, used as reset
        return t.Init(stub, "init", args)
    } else if
    function == "create_deal" { //create a new deal
        return t.create_deal(stub, args)
    } else if
    function == "update_deal" { //update a deal
        return t.update_deal(stub, args)
    } else if
    function == "addTransaction_inDeal" { //update a deal
        return t.addTransaction_inDeal(stub, args)
    } else if
    function == "create_transaction" { //create a new deal
        return t.create_transaction(stub, args)
    } else if
    function == "update_transaction" { //update a deal
        return t.update_transaction(stub, args)
    } else if
    function == "update_transaction_AllocationStatus" { //update a deal
        return t.update_transaction_AllocationStatus(stub, args)
    }
    fmt.Println("invoke did not find func: " + function)
    errMsg:= "{ \"message\" : \"Received unknown function invocation\", \"code\" : \"503\"}"
    err:= stub.SetEvent("errEvent", [] byte(errMsg))
    if err != nil {
        return nil, err
    }
    return nil, nil //error
}
// ============================================================================================================================
// Query - Our entry Dealint for Queries
// ============================================================================================================================
func(t * ManageDeals) Query(stub shim.ChaincodeStubInterface, function string, args[] string)([] byte, error) {
    fmt.Println("query is running " + function)
    // Handle different functions
    if
    function == "getDeal_byID" { //Read a Deal by dealId
        return t.getDeal_byID(stub, args)
    } else if
    function == "getDeal_byPledger" { //Read a Deal by Pledgee's name
        return t.getDeal_byPledger(stub, args)
    } else if
    function == "getDeal_byPledgee" { //Read a Deal by Pledgee's name
        return t.getDeal_byPledgee(stub, args)
    } else if
    function == "get_AllDeal" { //Read all Deals
        return t.get_AllDeal(stub, args)
    } else if
    function == "getTransaction_byDealID" { //Read all Deals
        return t.getTransaction_byDealID(stub, args)
    }
    fmt.Println("query did not find func: " + function) //errors
    errMsg:= "{ \"message\" : \"Received unknown function query\", \"code\" : \"503\"}"
    err:= stub.SetEvent("errEvent", [] byte(errMsg))
    if err != nil {
        return nil, err
    }
    return nil, nil
}
// ============================================================================================================================
// getDeal_byID - get Deal details for a specific ID from chaincode state
// ============================================================================================================================
func(t * ManageDeals) getDeal_byID(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var DealId string
    var err error
    fmt.Println("start getDeal_byID")
    if len(args) != 1 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting 'DealId' as an argument\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    // set dealId
    DealId = args[0]
    valAsbytes, err:= stub.GetState(DealId) //get the DealId from chaincode state
    if err != nil {
        errMsg:= "{ \"message\" : \"" + DealId + " not Found.\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    //fmt.Print("valAsbytes : ")
    //fmt.Println(valAsbytes)
    fmt.Println("end getDeal_byID")
    return valAsbytes, nil //send it onward
}
// ============================================================================================================================
//  getDeal_byPledger - get Deal details by Pledgee's name from chaincode state
// ============================================================================================================================
func(t * ManageDeals) getDeal_byPledger(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var jsonResp, pledgerName, errResp string
    var dealIndex[] string
    var valIndex Deals
    fmt.Println("start getDeal_byPledger")
    var err error
    if len(args) != 1 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting 'pledgerName' as an argument\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    // set Pledgee's name
    pledgerName = args[0]
    //fmt.Println("pledgerName" + pledgerName)
    dealAsBytes, err:= stub.GetState(DealIndexStr)
    if err != nil {
        return nil, errors.New("Failed to get Deal index string")
    }
    //fmt.Print("dealAsBytes : ")
    //fmt.Println(dealAsBytes)
    json.Unmarshal(dealAsBytes, &dealIndex) //un stringify it aka JSON.parse()
    fmt.Print("dealIndex : ")
    fmt.Println(dealIndex)
    //fmt.Println("len(dealIndex) : ")
    //fmt.Println(len(dealIndex))
    jsonResp = "{"
    for i, val:= range dealIndex {
        fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getDeal_byPledger")
        valueAsBytes, err:= stub.GetState(val)
        if err != nil {
            errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
            return nil, errors.New(errResp)
        }
        //fmt.Print("valueAsBytes : ")
        //fmt.Println(valueAsBytes)
        json.Unmarshal(valueAsBytes, &valIndex)
        fmt.Print("valIndex: ")
        fmt.Print(valIndex)
        if valIndex.Pledger == pledgerName {
            fmt.Println("Pledger found: " + val)
            jsonResp = jsonResp + "\"" + val + "\":" + string(valueAsBytes[: ])
            //fmt.Println("jsonResp inside if")
            //fmt.Println(jsonResp)
            if i < len(dealIndex) - 1 {
                jsonResp = jsonResp + ","
            }
        } 
    }
    jsonResp = jsonResp + "}"
    fmt.Println("jsonResp : " + jsonResp)
    if jsonResp == "{}" {
        fmt.Println("Pledger not found.")
        jsonResp =  "{ \"message\" : \"" + pledgerName + " Not Found.\", \"code\" : \"503\"}"
        errMsg:= "{ \"message\" : \"" + pledgerName + " Not Found.\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
        return nil, err
        }
    }
    //fmt.Print("jsonResp in bytes : ")
    //fmt.Println([]byte(jsonResp))
    fmt.Println("end getDeal_byPledger")
    return [] byte(jsonResp), nil //send it onward
}
// ============================================================================================================================
//  getDeal_byPledgee - get Deal details for a specific Pledgee from chaincode state
// ============================================================================================================================
func(t * ManageDeals) getDeal_byPledgee(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var jsonResp, pledgeeName, errResp string
    var dealIndex[] string
    var valIndex Deals
    fmt.Println("start getDeal_byPledgee")
    var err error
    if len(args) != 1 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting 'pledgeeName' as an argument\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    // set Pledgee name
    pledgeeName = args[0]
    //fmt.Println("pledgerName" + pledgeeName)
    dealAsBytes, err:= stub.GetState(DealIndexStr)
    if err != nil {
        return nil, errors.New("Failed to get Deal index")
    }
    //fmt.Print("dealAsBytes : ")
    //fmt.Println(dealAsBytes)
    json.Unmarshal(dealAsBytes, &dealIndex) //un stringify it aka JSON.parse()
    fmt.Print("dealIndex : ")
    fmt.Println(dealIndex)
    //fmt.Println("len(dealIndex) : ")
    //fmt.Println(len(dealIndex))
    jsonResp = "{"
    for i, val:= range dealIndex {
        fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getting pledgeeName")
        valueAsBytes, err:= stub.GetState(val)
        if err != nil {
            errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
            return nil, errors.New(errResp)
        }
        //fmt.Print("valueAsBytes : ")
        //fmt.Println(valueAsBytes)
        json.Unmarshal(valueAsBytes, &valIndex)
        fmt.Print("valIndex: ")
        fmt.Print(valIndex)
        if valIndex.Pledgee == pledgeeName {
            fmt.Println("Pledgee found")
            jsonResp = jsonResp + "\"" + val + "\":" + string(valueAsBytes[: ])
            //fmt.Println("jsonResp inside if")
            //fmt.Println(jsonResp)
            if i < len(dealIndex) - 1 {
                jsonResp = jsonResp + ","
            }
        }
    }
    jsonResp = jsonResp + "}"
    fmt.Println("jsonResp : " + jsonResp)
    if jsonResp == "{}" {
        fmt.Println("Pledgee not found.")
        jsonResp = jsonResp + "\"message\" : \"" + pledgeeName + " Not Found.\", \"code\" : \"503\"}"
        errMsg:= "{ \"message\" : \"" + pledgeeName + " Not Found.\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
    }
    //fmt.Print("jsonResp in bytes : ")
    //fmt.Println([]byte(jsonResp))
    fmt.Println("end getDeal_byPledgee")
    return [] byte(jsonResp), nil //send it onward
}
// ============================================================================================================================
//  get_AllDeal- get details of all Deal from chaincode state
// ============================================================================================================================
func(t * ManageDeals) get_AllDeal(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var jsonResp, errResp string
    var dealIndex[] string
    fmt.Println("start get_AllDeal")
    var err error
    if len(args) != 1 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting \" \" as an argument\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    dealAsBytes, err:= stub.GetState(DealIndexStr)
    if err != nil {
        return nil, errors.New("Failed to get Deal index")
    }
    //fmt.Print("dealAsBytes : ")
    //fmt.Println(dealAsBytes)
    json.Unmarshal(dealAsBytes, &dealIndex) //un stringify it aka JSON.parse()
    //fmt.Print("dealIndex : ")
    //fmt.Println(dealIndex)
    jsonResp = "{"
    for i, val:= range dealIndex {
        fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Deal")
        valueAsBytes, err:= stub.GetState(val)
        if err != nil {
            errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
            return nil, errors.New(errResp)
        }
        //fmt.Print("valueAsBytes : ")
        //fmt.Println(valueAsBytes)
        jsonResp = jsonResp + "\"" + val + "\":" + string(valueAsBytes[: ])
        if i < len(dealIndex) - 1 {
            jsonResp = jsonResp + ","
        }
    }
    //fmt.Println("len(dealIndex) : ")
    //fmt.Println(len(dealIndex))
    jsonResp = jsonResp + "}"
    //fmt.Println("jsonResp : " + jsonResp)
    //fmt.Print("jsonResp in bytes : ")
    //fmt.Println([]byte(jsonResp))
    fmt.Println("end get_AllDeal")
    return [] byte(jsonResp), nil
    //send it onward
}
// ============================================================================================================================
// Write - update Deal into chaincode state
// ============================================================================================================================
func(t * ManageDeals) update_deal(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var err error
    fmt.Println("Updating Deal")
    if len(args) != 9 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting 9\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    // set dealId
    dealId:= args[0]
    dealAsBytes, err:= stub.GetState(dealId) //get the Deal for the specified dealId from chaincode state
    if err != nil {
        errMsg:= "{ \"message\" : \"Failed to get state for " + dealId + "\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    //fmt.Print("dealAsBytes in update Deal")
    //fmt.Println(dealAsBytes);
    res:= Deals {}
    json.Unmarshal(dealAsBytes, &res)
    if res.DealID == dealId {
        fmt.Println("Deal found with dealId : " + dealId)
        //fmt.Println(res);
    } else {
        errMsg:= "{ \"message\" : \"" + dealId + " Not Found.\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    //build the Deal json string manually
    order:= `{` + `"dealId": "` + dealId + `" , ` + `"pledger": "` + args[1] + `" , ` + `"pledgee": "` + args[2] + `" , ` + `"maxValue": "` + args[3] + `" , ` + `"totalValueLongBoxAccount": "` + args[4] + `" , ` + `"totalValueSegregatedAccount": "` + args[5] + `" , ` + `"issueDate": "` + args[6] + `" , ` + `"lastSuccessfulAllocationDate": "` + args[7] + `" ` + `"transactions": "` + args[8] + `" , ` + `}`
    err = stub.PutState(dealId, [] byte(order)) //store Deal with id as key
    if err != nil {
        return nil, err
    }
    tosend:= "{ \"dealId\" : \"" + dealId + "\", \"message\" : \"Deal updated succcessfully\", \"code\" : \"200\"}"
    err = stub.SetEvent("evtsender", [] byte(tosend))
    if err != nil {
        return nil, err
    }
    fmt.Println("Deal updated succcessfully")
    return nil, nil
}
// ============================================================================================================================
// create Deal - create a new Deal, store into chaincode state
// ============================================================================================================================
func(t * ManageDeals) create_deal(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var err error
    if len(args) != 9 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting 9\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    fmt.Println("start create_deal")
    /*if len(args[0]) <= 0 {
        return nil, errors.New("1st argument must be a non-empty string")
    }
    */
    dealId:= args[0]
    Pledger:= args[1]
    Pledgee:= args[2]
    MaxValue:= args[3]
    TotalValueLongBoxAccount:= args[4]
    TotalValueSegregatedAccount:= args[5]
    IssueDate:= args[6]
    LastSuccessfulAllocationDate:= args[7]
    Transactions:= args[8]
    dealAsBytes, err:= stub.GetState(dealId)
    if err != nil {
        return nil, errors.New("Failed to get Deal dealId")
    }
    res:= Deals {}
    json.Unmarshal(dealAsBytes, &res)
    if res.DealID == dealId {
        //fmt.Println("This Deal arleady exists: " + dealId)
        //fmt.Println(res);
        errMsg:= "{ \"message\" : \"This Deal already exists\", \"code\" : \"503\"}"
        err:= stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil //all stop a Deal by this name exists
    }
    //build the Deal json string manually
    order:= `{` + `"dealId": "` + dealId + `" , ` + `"pledger": "` + Pledger + `" , ` + `"pledgee": "` + Pledgee + `" , ` + `"maxValue": "` + MaxValue + `" , ` + `"totalValueLongBoxAccount": "` + TotalValueLongBoxAccount + `" , ` + `"totalValueSegregatedAccount": "` + TotalValueSegregatedAccount + `" , ` + `"issueDate": "` + IssueDate + `" , ` + `"transactions": "` + Transactions + `" , ` + `"lastSuccessfulAllocationDate": "` + LastSuccessfulAllocationDate + `"  ` + `}`
    //fmt.Println("order: " + order)
    fmt.Print("order in bytes array: ")
    fmt.Println([] byte(order))
    err = stub.PutState(dealId, [] byte(order)) //store Deal with dealId as key
    if err != nil {
        return nil, err
    }
    //get the Deal index
    dealIndexAsBytes, err:= stub.GetState(DealIndexStr)
    if err != nil {
        return nil, errors.New("Failed to get Deal index")
    }
    var dealIndex[] string
    //fmt.Print("dealIndexAsBytes: ")
    //fmt.Println(dealIndexAsBytes)
    json.Unmarshal(dealIndexAsBytes, &dealIndex) //un stringify it aka JSON.parse()
    //fmt.Print("dealIndex after unmarshal..before append: ")
    //fmt.Println(dealIndex)
    //append
    dealIndex = append(dealIndex, dealId) //add Deal dealId to index list
    //fmt.Println("! Deal index after appending dealId: ", dealIndex)
    jsonAsBytes, _:= json.Marshal(dealIndex)
    //fmt.Print("jsonAsBytes: ")
    //fmt.Println(jsonAsBytes)
    err = stub.PutState(DealIndexStr, jsonAsBytes) //store name of Deal
    if err != nil {
        return nil, err
    }
    tosend:= "{ \"dealId\" : \"" + dealId + "\", \"message\" : \"Deal created succcessfully\", \"code\" : \"200\"}"
    err = stub.SetEvent("evtsender", [] byte(tosend))
    if err != nil {
        return nil, err
    }
    fmt.Println("end create_deal")
    return nil, nil
}
// ============================================================================================================================
// getTransaction_byDealID - get Transaction details for a specific Deal ID from chaincode state
// ============================================================================================================================
func(t * ManageDeals) getTransaction_byDealID(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var dealId string
    var err error
    var _tempJson Transactions
    fmt.Println("start getTransaction_byDealID")
    if len(args) != 1 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting 'dealId' as an argument\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    // set dealId
    dealId = args[0];
    dealAsBytes, err:= stub.GetState(dealId) //get the dealId from chaincode state
    if err != nil {
        errMsg:= "{ \"message\" : \"" + dealId + " not Found.\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    var dealIndex Deals
    json.Unmarshal(dealAsBytes, &dealIndex) //un stringify it aka JSON.parse()
    fmt.Print("dealIndex : ")
    fmt.Println(dealIndex)
    //fmt.Println("len(dealIndex) : ")
    fmt.Println(len(dealIndex.Transactions))
    _transactionSplit:= strings.Split(dealIndex.Transactions, ",")
    fmt.Print("_transactionSplit: ")
    fmt.Println(_transactionSplit)
    jsonResp := "{"
    for i:= range _transactionSplit {
        fmt.Println("_transactionSplit[i]: " + _transactionSplit[i])
        valueAsBytes, err:= stub.GetState(_transactionSplit[i])
        if err != nil {
            errResp := "{\"Error\":\"Failed to get state for " + _transactionSplit[i] + "\"}"
            return nil, errors.New(errResp)
        }
        json.Unmarshal(valueAsBytes, &_tempJson)
        fmt.Print("valueAsBytes : ")
        fmt.Println(valueAsBytes)
        jsonResp = jsonResp + "\"" + _transactionSplit[i] + "\":" + string(valueAsBytes[: ])
        if i < len(_transactionSplit) - 1 {
            jsonResp = jsonResp + ","
        }
    }
    jsonResp = jsonResp + "}"
    fmt.Print("jsonResp: ")
    fmt.Println(jsonResp)
    fmt.Println("end getTransaction_byDealID")
    return []byte(jsonResp), nil //send it onward
}
// ============================================================================================================================
// Write - update Deal into chaincode state
// ============================================================================================================================
func(t * ManageDeals) addTransaction_inDeal(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) { //to be done add margincalldate in arguments
    // dealid,transactionid
    var err error
    fmt.Println("addTransaction_inDeal")
    if len(args) != 5 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting 5\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    // set dealId
    dealId:= args[0]
    // set transactionId
    _transactionId:= args[1]
    _pledger:= args[2]
    _pledgee:= args[3]
    _marginCallDate := args[4]
    // converting string marginCallDate to int
    _newMarginCallDate,err := strconv.Atoi(_marginCallDate)
    if err != nil {
        fmt.Sprintf("Error while converting string '_marginCallDate' to int : %s", err.Error())
        return nil, errors.New("Error while converting string '_marginCallDate' to int ")
    }

    dealAsBytes, err:= stub.GetState(dealId) //get the Deal for the specified dealId from chaincode state
    if err != nil {
        errMsg:= "{ \"dealId\" : \"" + dealId + "\",\"message\" : \"Failed to get DealID\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    fmt.Print("dealAsBytes in update Deal")
    fmt.Println(dealAsBytes);
    res:= Deals {}
    _tempJson := Transactions{}
    json.Unmarshal(dealAsBytes, &res)
    fmt.Println(res);
    if res.DealID == dealId {
        fmt.Println("Deal found with dealId : " + dealId)
        _transactionSplit:= strings.Split(res.Transactions, ",")
        fmt.Print("_transactionSplit: ")
        fmt.Println(_transactionSplit)
        for i:= range _transactionSplit {
            fmt.Println("_transactionSplit[i]: " + _transactionSplit[i])
            valueAsBytes, err:= stub.GetState(_transactionSplit[i])
            if err != nil {
                errResp := "{\"Error\":\"Failed to get state for " + _transactionSplit[i] + "\"}"
                return nil, errors.New(errResp)
            }
            json.Unmarshal(valueAsBytes, &_tempJson)
            fmt.Print("valueAsBytes : ")
            fmt.Println(valueAsBytes)
            _oldMarginCallDate,err := strconv.Atoi(_tempJson.MarginCAllDate)
            if err != nil {
                fmt.Sprintf("Error while converting string '_tempJson.marginCAllDate' to int : %s", err.Error())
                return nil, errors.New("Error while converting string '_tempJson.marginCAllDate' to int ")
            }
            
            if _tempJson.Pledger == _pledger &&_tempJson.Pledgee == _pledgee{
                if _tempJson.AllocationStatus != "Allocation Succcessful" {
                    _tempJson.TransactionId = _transactionId;
                    _transactionSplit[i] = _transactionId;
                }
                if _oldMarginCallDate < _newMarginCallDate {
                    _tempJson.TransactionId = _transactionId;
                    _transactionSplit[i] = _transactionId;
                }
            }
            res.Transactions = res.Transactions + "," + _transactionSplit[i]
        }
    } else {
    errMsg:= "{ \"message\" : \"" + dealId + " Not Found.\", \"code\" : \"503\"}"
    err = stub.SetEvent("errEvent", [] byte(errMsg))
    if err != nil {
        return nil, err
    }
    return nil,nil
}
//build the Deal json string manually
order:= `{` + `"dealId": "` + res.DealID + `" , ` + `"pledger": "` + res.Pledger + `" , ` + `"pledgee": "` + res.Pledgee + `" , ` + `"maxValue": "` + res.MaxValue + `" , ` + `"totalValueLongBoxAccount": "` + res.TotalValueLongBoxAccount + `" , ` + `"totalValueSegregatedAccount": "` + res.TotalValueSegregatedAccount + `" , ` + `"issueDate": "` + res.IssueDate + `" , ` + `"transactions": "` + res.Transactions + `" , ` + `"lastSuccessfulAllocationDate": "` + res.LastSuccessfulAllocationDate + `" ` + `}`
err = stub.PutState(dealId, [] byte(order)) //store Deal with id as key
if err != nil {
    return nil, err
}
tosend:= "{ \"dealId\" : \"" + dealId + "\", \"message\" : \"Transaction added succcessfully\", \"code\" : \"200\"}"
err = stub.SetEvent("evtsender", [] byte(tosend))
if err != nil {
    return nil, err
}
fmt.Println("addTransaction_inDeal")
return nil, nil
}
// ============================================================================================================================
// update_transaction - update Transaction into chaincode state
// ============================================================================================================================
func(t * ManageDeals) update_transaction(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var err error
    fmt.Println(" update_transaction")
    if len(args) != 11 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting 11\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    // set dealId
    _transactionId:= args[0]
    dealAsBytes, err:= stub.GetState(_transactionId) //get the Deal for the specified dealId from chaincode state
    if err != nil {
        errMsg:= "{ \"message\" : \"Failed to get state for " + _transactionId + "\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    //fmt.Print("dealAsBytes in update Deal")
    //fmt.Println(dealAsBytes);
    res:= Transactions {}
    json.Unmarshal(dealAsBytes, &res)
    if res.TransactionId == _transactionId {
        fmt.Println("Deal found with _transactionId : " + _transactionId)
        //fmt.Println(res);
    } else {
        errMsg:= "{ \"message\" : \"" + _transactionId + " Not Found.\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    //build the Deal json string manually
    order:= `{` + `"transactionId": "` + args[0] + `" , ` + `"transactionDate": "` + args[1] + `" , ` + `"dealId": "` + args[2] + `" , ` + `"pledger": "` + args[3] + `" , ` + `"pledgee": "` + args[4] + `" , ` + `"rqv": "` + args[5] + `" , ` +
             `"currency": "` + args[6] + `" , ` + `"currencyConversionRate": "` + args[7] + `" , ` +  `"marginCAllDate": "` + args[8] + `" , ` + `"allocationStatus": "` + args[9] + `" , ` + `"transactionStatus": "` + args[10] + `" ` + `}`
    err = stub.PutState(args[2], [] byte(order)) //store Deal with id as key
    if err != nil {
        return nil, err
    }
    tosend:= "{ \"transactionId\" : \"" + _transactionId + "\", \"message\" : \"Transaction updated succcessfully\", \"code\" : \"200\"}"
    err = stub.SetEvent("evtsender", [] byte(tosend))
    if err != nil {
        return nil, err
    }
    fmt.Println("Transaction updated succcessfully")
    return nil, nil
}
// ============================================================================================================================
// update_transaction_AllocationStatus - update Transaction into chaincode state
// ============================================================================================================================
func(t * ManageDeals) update_transaction_AllocationStatus(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var err error
    fmt.Println(" update_transaction_AllocationStatus")
    if len(args) != 2 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting 2\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    // set dealId
    _transactionId:= args[0]
    dealAsBytes, err:= stub.GetState(_transactionId) //get the Deal for the specified dealId from chaincode state
    if err != nil {
        errMsg:= "{ \"message\" : \"Failed to get state for " + _transactionId + "\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    //fmt.Print("dealAsBytes in update Deal")
    //fmt.Println(dealAsBytes);
    res:= Transactions {}
    json.Unmarshal(dealAsBytes, &res)
    if res.TransactionId == _transactionId {
        fmt.Println("Deal found with _transactionId : " + _transactionId)
        //fmt.Println(res);
    } else {
        errMsg:= "{ \"message\" : \"" + _transactionId + " Not Found.\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    //build the Deal json string manually
    order:= `{` + `"transactionId": "` + _transactionId + `" , ` + `"transactionDate": "` + res.TransactionDate + `" , ` + `"dealId": "` + res.DealID + `" , ` + `"pledger": "` + res.Pledger + `" , ` + `"pledgee": "` + res.Pledgee + `" , ` + `"rqv": "` + res.RQV + `" , ` + `"currency": "` + res.Currency + `" , ` + `"marginCAllDate": "` + res.MarginCAllDate + `" , ` + `"allocationStatus": "` + args[1] + `" , ` + `"transactionStatus": "` + res.TransactionStatus + `" ` + `}`
    err = stub.PutState(_transactionId, [] byte(order)) //store Deal with id as key
    if err != nil {
        return nil, err
    }
    tosend:= "{ \"transactionId\" : \"" + _transactionId + "\", \"message\" : \"Transaction updated succcessfully\", \"code\" : \"200\"}"
    err = stub.SetEvent("evtsender", [] byte(tosend))
    if err != nil {
        return nil, err
    }
    fmt.Println("update_transaction_AllocationStatus")
    return nil, nil
}
// ============================================================================================================================
//  create_transaction - create a new Deal, store into chaincode state
// ============================================================================================================================
func(t * ManageDeals) create_transaction(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {
    var err error
    var _allocationStatus string
    if len(args) != 9 {
        errMsg:= "{ \"message\" : \"Incorrect number of arguments. Expecting 9\", \"code\" : \"503\"}"
        err = stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil
    }
    fmt.Println("start create_transaction")
    _transactionId:= args[0]
    _transactionStatus:= args[8];
    res:= Transactions {}
    dealAsBytes, err:= stub.GetState(_transactionId)
    json.Unmarshal(dealAsBytes, &res)
    if res.TransactionId == _transactionId {
        errMsg:= "{ \"transactionId\" : \"" + _transactionId + "\", \"message\" : \"This Transaction already exists\", \"code\" : \"503\"}"
        err:= stub.SetEvent("errEvent", [] byte(errMsg))
        if err != nil {
            return nil, err
        }
        return nil,nil //all stop a Deal by this name exists
    }
    if _transactionStatus == "Matched" {
        _allocationStatus = "Ready for allocation"
    } else if _transactionStatus == "Unmatched" {
        _allocationStatus = "Deal Unmatched. Can't be allocated"
    }
    //build the Deal json string manually
    order:= `{` + 
        `"transactionId": "` + args[0] + `" , ` +
        `"transactionDate": "` + args[1] + `" , ` + 
        `"dealId": "` + args[2] + `" , ` + 
        `"pledger": "` + args[3] + `" , ` + 
        `"pledgee": "` + args[4] + `" , ` + 
        `"rqv": "` + args[5] + `" , ` + 
        `"currency": "` + args[6] + `" , ` + 
        `"currencyConversionRate": "` + " " + `" , ` + 
        `"marginCAllDate": "` + args[7] + `" , ` + 
        `"allocationStatus": "` + _allocationStatus + `" , ` + 
        `"transactionStatus": "` + args[8] + `" ` + 
        `}`
    fmt.Println("order: " + order)
    //fmt.Print("order in bytes array: ")
    //fmt.Println([]byte(order))
    err = stub.PutState(_transactionId, [] byte(order)) //store Deal with dealId as key
    if err != nil {
        return nil, err
    }
    //Send DealID & transaction
    var temp[] string
    temp = append(temp, args[2], args[0], args[3], args[4],args[7])
    t.addTransaction_inDeal(stub, temp)
    tosend:= "{ \"dealId\" : \"" + args[2] + "\", \"message\" : \"Deal created succcessfully\", \"code\" : \"200\"}"
    err = stub.SetEvent("evtsender", [] byte(tosend))
    if err != nil {
        return nil, err
    }
    fmt.Println("end create_transaction")
    return nil, nil
}