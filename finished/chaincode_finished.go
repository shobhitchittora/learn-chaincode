package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// PremiumPayment model
type PremiumPayment struct {
	PolicyNumber  int    `json:"policynumber"`
	DOB           int    `json:"dob"` //utc date
	Email         string `json:"email"`
	ContactNumber string `json:"contactnumber"`
	Name          string `json:"name"`
	DueDate       int    `json:"duedate"` //utc date
	Amount        int    `json:"amount"`
}

var accountPrefix = "acct:"

// Account model for user
type Account struct {
	ID       string   `json:"id"`
	DOB      int      `json:"dob"`
	Email    string   `json:"email"`
	Balance  int      `json:"balance"`
	Policies []string `json:"policies"`
}

var claimPrefix = "claim:"

// ClaimInsurance model
type ClaimInsurance struct {
	AccID        string `json:"accountID"`
	PolicyNumber int    `json:"policynumber"`
	TypeOfClaim  string `json:"type"`
	DocVerified  bool   `json:"docverified"`
	Amount       int    `json:"amount"`
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
	} else if function == "init_payment" {
		return t.initPayment(stub, args)
	} else if function == "add_balance" {
		return t.addBalance(stub, args)
	} else if function == "create_account" {
		return t.createAccount(stub, args)
	} else if function == "claim_insurance" {
		return t.claimInsurance(stub, args)
	} else if function == "buy_policy" {
		return t.buyPolicy(stub, args)
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
	} else if function == "get_balance" {
		return t.getBalance(stub, args)
	} else if function == "get_claim" {
		return t.getClaim(stub, args)
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

func (t *SimpleChaincode) getBalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting id of the account")
	}

	fmt.Println("- start getBalance")

	if len(args[0]) <= 0 {
		return nil, errors.New("UserName must be non-empty int")
	}

	accid := strings.ToLower(args[0])

	valAsbytes, err := stub.GetState(accountPrefix + accid)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + accid + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *SimpleChaincode) getClaim(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting policynumber for claim")
	}

	fmt.Println("- start getClaim")

	if len(args[0]) <= 0 {
		return nil, errors.New("PolicyNumber must be non-empty int")
	}

	policyNumberString := args[0]

	valAsbytes, err := stub.GetState(claimPrefix + policyNumberString)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get claim for " + policyNumberString + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

//Init Payment for Premium
func (t *SimpleChaincode) initPayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 7 {
		return nil, errors.New("Incorrect number of arguments. Expecting 7")
	}

	fmt.Println("- start init payment")
	if len(args[0]) <= 0 {
		return nil, errors.New("Policy Number must be non-empty int")
	}

	if len(args[1]) <= 0 {
		return nil, errors.New("DOB must be non-empty int")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("Email must be non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("ContactNumber must be non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("Name must be non-empty string")
	}
	if len(args[5]) <= 0 {
		return nil, errors.New("DueDate must be non-empty int")
	}
	if len(args[6]) <= 0 {
		return nil, errors.New("Amount must be non-empty int")
	}

	//Preparing fields for payment struct
	//************************************
	PolicyNumber, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("PolicyNumber must be a numeric string")
	}

	DOB, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("DOB must be a numeric string")
	}

	Email := strings.ToLower(args[2])
	ContactNumber := strings.ToLower(args[3])
	Name := strings.ToLower(args[4])

	DueDate, err := strconv.Atoi(args[5])
	if err != nil {
		return nil, errors.New("DueDate must be a numeric string")
	}

	Amount, err := strconv.Atoi(args[6])
	if err != nil {
		return nil, errors.New("Amount must be a numeric string")
	}

	//********************************
	//Check if balance >= amount
	//********************************

	accountAsBytes, err := stub.GetState("account")
	if err != nil {
		return nil, errors.New("Failed to get the Account info for the PolicyNumber")
	}

	acc := Account{}
	json.Unmarshal(accountAsBytes, &acc)

	if acc.Balance < Amount {
		fmt.Println("Not Enough Balance for PolicyNumber: " + strconv.Itoa(PolicyNumber))
		fmt.Println(acc)
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
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//Create Account for a user
func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	fmt.Println("- start create_account")

	if len(args[0]) <= 0 {
		return nil, errors.New("UserName must be non-empty int")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("DOB must be non-empty")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("Email must be non-empty")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("Balance must be non-empty")
	}

	username := strings.ToLower(args[0])

	dob, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("DOB must be a numeric string")
	}

	email := strings.ToLower(args[2])

	var policies []string

	var account = Account{ID: username, DOB: dob, Email: email, Balance: 0, Policies: policies}
	accountBytes, err := json.Marshal(&account)
	if err != nil {
		fmt.Println("error creating account" + account.ID)
		return nil, errors.New("Error creating account " + account.ID)
	}

	fmt.Println("Attempting to get state of any existing account for " + account.ID)
	existingBytes, err := stub.GetState(accountPrefix + account.ID)
	if err == nil {

		var user Account
		err = json.Unmarshal(existingBytes, &user)
		if err != nil {
			fmt.Println("Error unmarshalling account " + account.ID + "\n--->: " + err.Error())

			if strings.Contains(err.Error(), "unexpected end") {
				fmt.Println("No data means existing account found for " + account.ID + ", initializing account.")
				err = stub.PutState(accountPrefix+account.ID, accountBytes)

				if err == nil {
					fmt.Println("created account" + accountPrefix + account.ID)
					return nil, nil
				} else {
					fmt.Println("failed to create initialize account for " + account.ID)
					return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
				}
			} else {
				//return nil, errors.New("Error unmarshalling existing account " + account.ID)
			}
		} else {
			fmt.Println("Account already exists for " + account.ID + " " + user.ID)
			return nil, errors.New("Can't reinitialize existing user " + account.ID)
		}
	} else {

		fmt.Println("No existing account found for " + account.ID + ", initializing account.")
		err = stub.PutState(accountPrefix+account.ID, accountBytes)

		if err == nil {
			fmt.Println("created account" + accountPrefix + account.ID)
			return nil, nil
		} else {
			fmt.Println("failed to create initialize account for " + account.ID)
			return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
		}

	}

	return nil, nil
}

// Add balance for Account
func (t *SimpleChaincode) addBalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	fmt.Println("- start add_balance")

	if len(args[0]) <= 0 {
		return nil, errors.New("UserName must be non-empty int")
	}

	if len(args[1]) <= 0 {
		return nil, errors.New("Balance must be non-empty")
	}

	id := strings.ToLower(args[0])

	balance, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Balance must be a numeric string")
	}

	//******************************
	//Check if Account exists or not
	//******************************

	accountBytes, err := stub.GetState(accountPrefix + id)
	if err == nil {
		//Account found
		var account Account
		err = json.Unmarshal(accountBytes, &account)
		if err != nil {
			return nil, errors.New("Account reading problem")
		}

		//Increment the balance
		account.Balance = account.Balance + balance

		//Create byte stream of new Account details
		newAccountBytes, err := json.Marshal(&account)
		if err != nil {
			fmt.Println("error creating account" + account.ID)
			return nil, errors.New("Error creating account " + account.ID)
		}
		err = stub.PutState(accountPrefix+id, newAccountBytes)
		if err != nil {
			return nil, errors.New("Error adding balance")
		}

	} else {
		//Account not found
		return nil, errors.New("No account found for ID -->" + id)
	}

	return nil, nil
}

//Buy Policy
func (t *SimpleChaincode) buyPolicy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	fmt.Println("- start buy_policy")

	if len(args[0]) <= 0 {
		return nil, errors.New("AccID must be non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("PolicyNumber must be non-empty int")
	}

	id := strings.ToLower(args[0])

	policyNumberString := args[1]

	// 	policynumber, err:= strconv.Atoi(args[1])
	// 	if err!=nil{
	// 		return nil, errors.New("PolicyNumber must be a numeric string")
	// 	}

	accountBytes, err := stub.GetState(accountPrefix + id)
	if err == nil {
		//Account found
		var account Account
		err = json.Unmarshal(accountBytes, &account)
		if err != nil {
			return nil, errors.New("Account reading problem")
		}

		policies := account.Policies

		if !stringInSlice(policyNumberString, policies) {
			account.Policies = append(account.Policies, policyNumberString)

			//Create byte stream of new Account details
			newAccountBytes, err := json.Marshal(&account)
			if err != nil {
				fmt.Println("error creating account" + account.ID)
				return nil, errors.New("Error creating account " + account.ID)
			}

			err = stub.PutState(accountPrefix+id, newAccountBytes)
			if err != nil {
				return nil, errors.New("Error adding balance")
			}

		} else {
			return nil, errors.New("policy already bought")
		}

	} else {
		//Account not found
		return nil, errors.New("No account found for ID -->" + id)
	}

	return nil, nil
}

//Clain Insurance
func (t *SimpleChaincode) claimInsurance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var jsonResp string

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	fmt.Println("- start claim_insurance")

	if len(args[0]) <= 0 {
		return nil, errors.New("AccID must be non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("PolicyNumber must be non-empty int")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("Type must be non-empty")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("DocVerified must be non-empty")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("Amount must be non-empty")
	}

	id := strings.ToLower(args[0])

	policyNumberString := args[1]

	policynumber, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("PolicyNumber must be a numeric string")
	}

	insuranceType := strings.ToLower(args[2])

	docverified, err := strconv.ParseBool(args[3])
	if err != nil {
		return nil, errors.New("docverified must be a bool string")
	}

	amount, err := strconv.Atoi(args[4])
	if err != nil {
		return nil, errors.New("amount must be a numeric string")
	}

	//Create a claim object
	var claim = ClaimInsurance{AccID: id, PolicyNumber: policynumber, TypeOfClaim: insuranceType, DocVerified: docverified, Amount: amount}
	claimBytes, err := json.Marshal(&claim)
	if err != nil {
		fmt.Println("error claiming insurance for " + id + " and policy no - " + policyNumberString)
		return nil, errors.New("Error claiming insurance " + id)
	}

	//****************************************
	//If has policy and not already claimed
	//****************************************
	accountBytes, err := stub.GetState(accountPrefix + id)
	if err == nil {
		//Account found
		var account Account
		err = json.Unmarshal(accountBytes, &account)
		if err != nil {
			return nil, errors.New("Account reading problem")
		}

		policies := account.Policies

		//Has policy
		if stringInSlice(policyNumberString, policies) {
			//Not already claimed

			// Reading the Claim from BlockChain and checking for validity
			existingClaimBytes, err := stub.GetState(claimPrefix + policyNumberString)
			if err != nil {
				jsonResp = "{\"Error\":\"Already Claimed for policymnumber " + policyNumberString + "\"}"
				return nil, errors.New(jsonResp)
			} else {
				//=> Claim not already found
				var claimObject ClaimInsurance
				err = json.Unmarshal(existingClaimBytes, &claimObject)
				if err != nil {
					fmt.Println("Filing Claim for policyNumber " + policyNumberString)
					err = stub.PutState(claimPrefix+policyNumberString, claimBytes)
					if err != nil {
						jsonResp = "{\"Error\":\"Error adding claim for - " + policyNumberString + "\"}"
						return nil, errors.New(jsonResp)
					}
				} else {
					return nil, errors.New("Sucess unmarshalling existing Claim " + claimPrefix + policyNumberString)
				}

			}
		} else {
			jsonResp = "{\"Error\":\"Policy Not bought for  - " + policyNumberString + "\"}"
			return nil, errors.New(jsonResp)
		}

	} else {
		//Account not found
		jsonResp = "{\"Error\":\"No account found for ID --> " + id + "\"}"
		return nil, errors.New(jsonResp)
	}

	return nil, nil
}
