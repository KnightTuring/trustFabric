package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {

}

type Vote struct {
  Voter string `json:"voter"`
	VotedFor string `json:"votedFor"`
  Timestamp string `json:"timestamp"`
	Location string `json:"location"`
}

//called during instantiation of smart contract by the network
//best practice to define every
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response{
	return shim.Success(nil)
}

//called when application requests to run the Smart Contract "voting-chaincode"
//also specifies which function to call within the SC
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response{

	function, args := APIstub.GetFunctionAndParameters()

	if function == "recordVote" {
		return s.recordVote(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryAllVotes"{
		return s.queryAllVotes(APIstub)
	} else if function == "queryVote" {
		return s.queryVote(APIstub, args)
	}

	return shim.Error("Function does not exist in Smart Contract")
}

//initLedger is called when ledger is instantiated by the network
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response{
	vote := []Vote{
		Vote{Voter: "Abdul", VotedFor: "CESA", Timestamp: "1510417417", Location: "PCCOE"},
		Vote{Voter: "Ankush", VotedFor: "MESA", Timestamp: "1510417417", Location: "PCCOE"},
		Vote{Voter: "Chandan", VotedFor: "ITSA", Timestamp: "1510417417", Location: "PCCOE"},
		Vote{Voter: "Ritvik", VotedFor: "CiESA", Timestamp: "1510417417", Location: "PCCOE"},
	}

	i := 0
	//iterate over all JSON array
	for i < len(vote){
		fmt.Println("i is ", i)

		//convert respective array object to bytes
		voteInBytes,_ := json.Marshal(vote[i])
		//put vote
		APIstub.PutState(strconv.Itoa(i+1), voteInBytes)
		fmt.Println("Added ", vote[i])

		i = i+1
	}

	return shim.Success(nil)
}

//record Vote
func (s *SmartContract) recordVote(APIstub shim.ChaincodeStubInterface, args[]string) sc.Response{

	if len(args) != 5 {
		return shim.Error("Incorrect no. of arguments passed,  expecting 5 including serial no")
	}
	var vote = Vote{ Voter: args[1], VotedFor: args[2], Timestamp: args[3], Location: args[4]}
	voteAsBytes,_ := json.Marshal(vote)

	 //TODO check if user has permission to cast vote
	 err := APIstub.PutState(args[0], voteAsBytes)
	 if err !=nil {
		 return shim.Error(fmt.Sprintf("Failed to record Vote %s", args[0]))
	 }

	 return shim.Success(nil)
}

//query single vote by giving ID
func (s *SmartContract) queryVote(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments passed")
	}

	voteAsBytes, _ := APIstub.GetState(args[0])
	if voteAsBytes == nil {
		return shim.Error("Could not locate specified vote")
	}

	return shim.Success(voteAsBytes)
}

//query all the votes in the ledger
func (s *SmartContract) queryAllVotes(APIstub shim.ChaincodeStubInterface) sc.Response{
	start := "0"
	end := "999"

	resultsIterator, err := APIstub.GetStateByRange(start, end)

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bufMemberAlreadyWritten := false

	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()

		if err!= nil {
			return shim.Error(err.Error())
		}

		if bufMemberAlreadyWritten == true{
			buffer.WriteString(",")
		}
		buffer.WriteString("{\" Key\" :")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bufMemberAlreadyWritten = true

	}
	buffer.WriteString("]")

	fmt.Printf(" All Votes are: \n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main(){
	err := shim.Start(new(SmartContract))

	if err == nil {
		fmt.Printf("Error creating Smart Contract")
	}
}
