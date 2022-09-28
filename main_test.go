package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/tuneinsight/lattigo/v3/bfv"
	"github.com/tuneinsight/lattigo/v3/rlwe"
)

var paramDef = bfv.PN13QP218

func genKeyPair(t *testing.T) (*rlwe.SecretKey, *rlwe.PublicKey) {
	paramDef := bfv.PN13QP218
	param, err := bfv.NewParametersFromLiteral(paramDef)
	if err != nil {
		fmt.Println("Generate key pair failed", err.Error())
		t.FailNow()
	}

	kgen := bfv.NewKeyGenerator(param)
	return kgen.GenKeyPair()
}

func decrypt(t *testing.T, data []byte, sk *rlwe.SecretKey) int64 {
	param, err := bfv.NewParametersFromLiteral(paramDef)
	if err != nil {
		fmt.Println("decrypt: set of BFV parameters failed", err.Error())
		t.FailNow()
	}
	ciphertext := bfv.NewCiphertext(param, 1)
	err = ciphertext.UnmarshalBinary(data)
	if err != nil {
		fmt.Println("failed to unmarshal value")
		t.FailNow()
	}
	// Decrypt data with private key
	decryptor := bfv.NewDecryptor(param, sk)
	text := decryptor.DecryptNew(ciphertext)
	// decode data
	encoder := bfv.NewEncoder(param)
	return encoder.DecodeIntNew(text)[0]
}

func encrypt(t *testing.T, data int64, pubkey *rlwe.PublicKey) []byte {
	param, err := bfv.NewParametersFromLiteral(paramDef)
	if err != nil {
		fmt.Println("encrypt: set of BFV parameters failed", err.Error())
		t.FailNow()
	}

	text := bfv.NewPlaintext(param)

	// Generate a encoder and encode the data
	encoder := bfv.NewEncoder(param)
	encoder.Encode([]int64{data}, text)

	// Encrypt data with public key
	encryptor := bfv.NewEncryptor(param, pubkey)
	ciphertext := encryptor.EncryptNew(text)
	raw, err := ciphertext.MarshalBinary()
	if err != nil {
		fmt.Println("marshalBinary encodes a Ciphertext failed", err.Error())
		t.FailNow()
	}
	return raw
}

func checkInit(t *testing.T, stub *shimtest.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkeCreateReport(t *testing.T, stub *shimtest.MockStub, subject string, pubkey []byte) {
	args := [][]byte{[]byte("CreateReport"), []byte(subject), pubkey}
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Create report", subject, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkeSubmitData(t *testing.T, stub *shimtest.MockStub, subject, department string, revenue []byte) {
	args := [][]byte{[]byte("SubmitData"), []byte(subject), []byte(department), revenue}
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("submit data", subject, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkeQueryData(t *testing.T, stub *shimtest.MockStub, subject string, value int64, sk *rlwe.SecretKey) {
	args := [][]byte{[]byte("QueryData"), []byte(subject)}
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Query data", subject, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query data", subject, "failed to get value")
		t.FailNow()
	}

	count := decrypt(t, res.Payload, sk)
	if count != value {
		fmt.Println("Query value", count, "was not", value, "as expected")
		t.FailNow()
	}
}

func TestLedgerInit(t *testing.T) {
	cc := new(Ledger)
	stub := shimtest.NewMockStub("ledger", cc)

	checkInit(t, stub, nil)
}

func TestLedgerCreateReport(t *testing.T) {
	cc := new(Ledger)
	stub := shimtest.NewMockStub("ledger", cc)

	checkInit(t, stub, nil)

	_, pk := genKeyPair(t)
	pkBinary, err := pk.MarshalBinary()
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	checkeCreateReport(t, stub, "October", pkBinary)
}

func TestLedgerSubmitData(t *testing.T) {
	cc := new(Ledger)
	stub := shimtest.NewMockStub("ledger", cc)

	checkInit(t, stub, nil)

	_, pk := genKeyPair(t)
	pkBinary, err := pk.MarshalBinary()
	if err != nil {
		fmt.Println("marshal binary encodes a public key", err.Error())
		t.FailNow()
	}
	checkeCreateReport(t, stub, "October", pkBinary)

	revenue := encrypt(t, 2000, pk)
	checkeSubmitData(t, stub, "October", "VoneChain", revenue)
}

func TestLedgerQueryData(t *testing.T) {
	cc := new(Ledger)
	stub := shimtest.NewMockStub("ledger", cc)

	checkInit(t, stub, nil)

	sk, pk := genKeyPair(t)
	pkBinary, err := pk.MarshalBinary()
	if err != nil {
		fmt.Println("marshal binary encodes a public key", err.Error())
		t.FailNow()
	}
	checkeCreateReport(t, stub, "October", pkBinary)

	revenue1 := encrypt(t, 2000, pk)
	checkeSubmitData(t, stub, "October", "VoneChain", revenue1)

	revenue2 := encrypt(t, -1000, pk)
	checkeSubmitData(t, stub, "October", "GitHub", revenue2)

	checkeQueryData(t, stub, "October", 1000, sk)
}
