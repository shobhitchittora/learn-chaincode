package main

import (
	"errors"
	"fmt"
	"strings"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type PremiumPayment struct {
	PolicyNumber int `json:"policynumber"`
	DOB int `json:"dob"`				//utc date
	Email string `json:"email"`
	ContactNumber string `json:"contactnumber"`
	Name string `json:"name"`
	DueDate int `json:"duedate"` 	//utc date
	Amount int `json:"amount"`
}

type Account struct {
	PolicyNumber int `json:"policynumber"`
	Balance int `json:"balance"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}else if function == "init_payment"{
		return t.init_payment(stub,args)
	}else if function == "generate_balance"{
		return t.generate_balance(stub,args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

//Init Payment for Premium
func (t *SimpleChaincode) init_payment(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	var err error
	
	if len(args)!=7{
		return nil, errors.New("Incorrect number of arguments. Expecting 7")
	}
	
	fmt.Println("- start init payment")
	if len(args[0])<=0{
		return nil, errors.New("Policy Number must be non-empty int")
	}
	
	if len(args[1])<=0{
		return nil, errors.New("DOB must be non-empty int")
	}
	if len(args[2])<=0{
		return nil, errors.New("Email must be non-empty string")
	}
	if len(args[3])<=0{
		return nil, errors.New("ContactNumber must be non-empty string")
	}
	if len(args[4])<=0{
		return nil, errors.New("Name must be non-empty string")
	}
	if len(args[5])<=0{
		return nil, errors.New("DueDate must be non-empty int")
	}
	if len(args[6])<=0{
		return nil, errors.New("Amount must be non-empty int")
	}
	
	//Preparing fields for payment struct
	//************************************
	PolicyNumber, err:= strconv.Atoi(args[0]);
	if err!=nil{
		return nil, errors.New("PolicyNumber must be a numeric string")
	}
	
	DOB, err:= strconv.Atoi(args[1]);
	if err!=nil{
		return nil, errors.New("DOB must be a numeric string")
	}
	
	Email := strings.ToLower(args[2]);
	ContactNumber := strings.ToLower(args[3]);
	Name := strings.ToLower(args[4]);
	
	DueDate ,err:= strconv.Atoi(args[5]);
	if err!=nil{
		return nil, errors.New("DueDate must be a numeric string")
	}
	
	Amount ,err:= strconv.Atoi(args[6]);
	if err!=nil{
		return nil, errors.New("Amount must be a numeric string")
	}
	
	//********************************
	//Check if balance >= amount 
	//********************************
	
	accountAsBytes, err := stub.GetState(strconv.Itoa(PolicyNumber))
	if err != nil {
		return nil, errors.New("Failed to get the Account info for the PolicyNumber")
	}
	
	acc := Account{}
	json.Unmarshal(accountAsBytes, &acc)
	
	if(acc.Balance < Amount){
		fmt.Println("Not Enough Balance for PolicyNumber: " + strconv.Itoa(PolicyNumber))
		fmt.Println(acc);
		return nil, errors.New("Transaction Cancelled")	
	}
	
	res := `{"policynumber" :` + strconv.Itoa(PolicyNumber) +
					`, "dob": ` + strconv.Itoa(DOB) +
					`, "email": "` + Email +
					`", "contactnumber": "` + ContactNumber +
					`", "name": "` + Name +
					`", "duedate": ` + strconv.Itoa(DueDate) +
					`, "amount": ` + strconv.Itoa(Amount) + `}`
	
	err = stub.PutState(strconv.Itoa(PolicyNumber), []byte(res))
	if err!=nil{
		return nil, err
	}
	
	return nil,nil
}

// Generate balance for Payment
func (t *SimpleChaincode) generate_balance(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	var err error
	
	if len(args)!=2{
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	fmt.Println("- start generate_balance")
	
	if len(args[0])<=0{
		return nil, errors.New("PolicyNumber must be non-empty int")
	}
	
	if len(args[1])<=0{
		return nil, errors.New("Balance must be non-empty")
	}
	
	PolicyNumber, err := strconv.Atoi(args[0])
	if err!=nil{
		return nil, errors.New("PolicyNumber must be a numeric string")
	}
	
	Balance, err:= strconv.Atoi(args[1])
	if err!=nil{
		return nil, errors.New("Balance must be a numeric string")
	}
	
	res :=  `{"policynumber" :` + strconv.Itoa(PolicyNumber) +
					`, "balance": ` + strconv.Itoa(Balance) + `}`
	
	err = stub.PutState(strconv.Itoa(PolicyNumber), []byte(res))
	if err!=nil{
		return nil, err
	}
	
	return nil,nil
}