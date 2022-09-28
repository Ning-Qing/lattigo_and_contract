package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/tuneinsight/lattigo/v3/bfv"
	"github.com/tuneinsight/lattigo/v3/rlwe"
)

// report records a relevant report
type report struct {
	Subject string            `json:"subject"`
	Data    map[string][]byte `json:"data"`
	Count   []byte            `json:"count"`
	Pubkey  []byte            `json:"-"` // public key
}

//        report
//	------------------
// | subject: October |
//	------------------
// |   A    |   1000  |
//	------------------
// |   B    |   2000  |
//	------------------
// |   count: 3000    |
//  ------------------

func newReport(subject string, pubkey []byte) *report {
	return &report{
		Subject: subject,
		Data:    make(map[string][]byte),
		Pubkey:  pubkey,
	}
}

type Ledger struct {
	param bfv.Parameters // BFV parameters
}

func (l *Ledger) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// set of BFV parameters
	paramDef := bfv.PN13QP218
	paramDef.T = 0x3ee0001
	param, err := bfv.NewParametersFromLiteral(paramDef)
	if err != nil {
		return shim.Error(fmt.Sprintf("init chaincode failed: %s", err.Error()))
	}
	l.param = param
	return shim.Success(nil)
}

func (l *Ledger) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	var err error
	args := stub.GetArgs()
	function := string(args[0])
	args = args[1:]

	switch function {
	case "CreateReport":
		if err = checkeArgs(args, 2); err != nil {
			return shim.Error(fmt.Sprintf("create report failed: %s", err.Error()))
		}
		// args[0] is subject string
		// args[1] is public key []byte
		err = l.CreateReport(stub, string(args[0]), args[1])
		if err != nil {
			return shim.Error(fmt.Sprintf("create report failed: %s", err.Error()))
		}
		return shim.Success(nil)
	case "SubmitData":
		if err = checkeArgs(args, 3); err != nil {
			return shim.Error(fmt.Sprintf("submit data failed: %s", err.Error()))
		}
		err = l.SubmitData(stub, string(args[0]), string(args[1]), args[2])
		if err != nil {
			return shim.Error(fmt.Sprintf("submit data failed: %s", err.Error()))
		}
		return shim.Success(nil)
	case "QueryData":
		if err = checkeArgs(args, 1); err != nil {
			return shim.Error(fmt.Sprintf("query data failed: %s", err.Error()))
		}
		payload, err := l.QueryData(stub, string(args[0]))
		if err != nil {
			return shim.Error(fmt.Sprintf("query data failed: %s", err.Error()))
		}
		return shim.Success(payload)
	}
	return shim.Error(fmt.Sprintf("unknown method: %s", function))
}

// CreateReport create a report that records the subject and public key
func (l *Ledger) CreateReport(stub shim.ChaincodeStubInterface, subject string, pubkey []byte) error {
	var err error
	report := newReport(subject, pubkey)

	value, err := json.Marshal(report)
	if err != nil {
		return err
	}
	err = stub.PutState(subject, value)
	if err != nil {
		return err
	}
	return nil
}

func (l *Ledger) SubmitData(stub shim.ChaincodeStubInterface, subject, department string, revenue []byte) error {
	var err error
	value, err := stub.GetState(subject)
	if err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("no related subject retrieved")
	}
	report := &report{}
	err = json.Unmarshal(value, &report)
	if err != nil {
		return err
	}
	report.Data[department] = revenue
	value, err = json.Marshal(report)
	if err != nil {
		return err
	}
	err = stub.PutState(subject, value)
	if err != nil {
		return err
	}
	return nil
}

func (l *Ledger) QueryData(stub shim.ChaincodeStubInterface, subject string) ([]byte, error) {
	var err error
	// get report
	value, err := stub.GetState(subject)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("no related subject retrieved")
	}
	report := &report{}
	err = json.Unmarshal(value, &report)
	if err != nil {
		return nil, err
	}
	// creates a new Evaluator, that can be used to do homomorphic operations
	// on ciphertexts and/or plaintexts. It stores a memory buffer and ciphertexts
	// that will be used for intermediate values.
	evaluator := bfv.NewEvaluator(l.param, rlwe.EvaluationKey{})
	// create a Ciphertext to store the results
	results := bfv.NewCiphertext(l.param, 1)
	for _, v := range report.Data {
		op := bfv.NewCiphertext(l.param, 1)
		err = op.UnmarshalBinary(v)
		if err != nil {
			return nil, fmt.Errorf("calculation failed: %s", err.Error())
		}
		// equivalent to: results = results + op
		evaluator.Add(results, op, results)
	}
	return results.MarshalBinary()
}

func checkeArgs(args [][]byte, expect int) error {
	length := len(args)
	if length != expect {
		return fmt.Errorf("arguments not as expected, %d was required but %d", expect, length)
	}
	return nil
}

func main() {
	err := shim.Start(new(Ledger))
	if err != nil {
		fmt.Printf("Error starting ledger chaincode: %s", err)
	}
}
