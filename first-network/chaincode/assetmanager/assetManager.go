package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type assetManagerChaincode struct {
}

//Init assetManager & create first user whose id is "public"
func (t *assetManagerChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response{
	fmt.Println("assetManager Init....")
	err := stub.PutState("public", []byte(strconv.Itoa(0)))
	if err != nil{
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *assetManagerChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("assetManager Invoke")
	f, args := stub.GetFunctionAndParameters()
	if f == "createUserAsset"{
		return t.createUserAsset(stub, args)
	}else if f == "deleteUserAsset"{
		return t.deleteUserAsset(stub, args)
	}else if f == "queryUserAsset"{
		return t.queryUserAsset(stub, args)
	}else if f == "setUserAsset"{
		return t.setUserAsset(stub, args)
	}


	return shim.Error("Invalid invoke function name.")
}

//Add a user to our system.
//arg[0]: the user id(hash of {name + date}.
//arg[1]: asset that the user own.default 0
func (t *assetManagerChaincode) createUserAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	var id string
	var asset int
	var err error

	if len(args) == 1{
		id = args[0]
		asset = 0
	}else if len(args) == 2{
		id = args[0]
		asset, err = strconv.Atoi(args[1])
		if err != nil{
			return shim.Error("Invalid asset amount, expecting a integer value")
		}
	}else{
		return shim.Error("Incorrect number of arguments.")
	}

	err = stub.PutState(id, []byte(strconv.Itoa(asset)))
	if err != nil{
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//Delete a user account
//arg[0]: user id
func (t *assetManagerChaincode) deleteUserAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	var err error

	if len(args) != 1{
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	err = stub.DelState(args[0])
	if err != nil{
		return shim.Error("Failed to delete user: " + err.Error())
	}

	return shim.Success(nil)
}

//Query asset of the specific user
//arg[0]: user id
func (t *assetManagerChaincode) queryUserAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	var err error

	if len(args) != 1{
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	asset, err := stub.GetState(args[0])
	if err != nil{
		jsonResp := "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	if asset == nil{
		jsonResp := "{\"Error\":\"Nil to ammount for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"User id\":\"" + args[0] + "\",\"Amount\":\"" + string(asset) + "\"}"
	fmt.Printf("Query Response: %s\n", jsonResp)

	return shim.Success(asset)
}

//Set asseet of the specific user
//arg[0]: user id
//arg[1]: old asset value
//arg[1]: new asset value
func (t *assetManagerChaincode) setUserAsset(stub shim.ChaincodeStubInterface, args []string) pb.Reponse{
	var err error
	var oldValue int
	var newValue int

	if len(args) != 3{
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	oldValue = args[1]
	oldValue, err = strconv.Atoi(args[1])
	if err != nil{
		return shim.Error("Invalid asset amount, expecting a integer value")
	}
	newValue, err = strconv.Atoi(args[2])
	if err != nil{
		return shim.Error("Invalid asset amount, expecting a integer value")
	}
	if newValue < 0{
		jsonResp := "{\"Error\":\"Incorrect new asset amount for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	old, err := stub.GetState(args[0])
	if err != nil{
		jsonResp := "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	if old == nil{
		jsonResp := "{\"Error\":\"Nil to ammount for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	if oldValue != old{
		jsonResp := "{\"Error\":\"Incorrect old asset amount for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.PutState(args[0], []byte(strconv.Itoa(newValue)))
	if err != nil{
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(assetManagerChaincode))
	if err != nil {
		fmt.Printf("Error starting assetManagerChaincode chaincode: %s", err)
	}
}