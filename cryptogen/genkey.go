package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"os"

	"github.com/spf13/cobra"
	"github.com/tuneinsight/lattigo/v3/bfv"
	"github.com/tuneinsight/lattigo/v3/rlwe"
)

var (
	outpath string
)

var genkey = &cobra.Command{
	Short: "genkey",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func initGenkey() {
	genkey.Flags().StringVarP(&outpath, "outpath", "o", ".out", "save the path to the key pair")
}

func GenKey() error {
	var err error
	sk, pk, err := NewKeyPair()
	if err != nil {
		return err
	}
	_, err := os.Stat(outpath)
	if 
	pem.Encode()
}

type SecretKey struct {
	Param string          `json:"param"`
	SK    *rlwe.SecretKey `json:"secretKey"`
}

func (sk *SecretKey) Marshal() ([]byte, error) {
	var src, dst []byte
	var err error
	src, err = json.Marshal(sk)
	if err != nil {
		return dst, err
	}
	base64.RawStdEncoding.Encode(dst, src)
	return dst, err
}

func (sk *SecretKey) Unmarshal(src []byte) error {
	var dst []byte
	var err error
	_, err = base64.RawStdEncoding.Decode(dst, src)
	return err
}

type PublicKey struct {
	Param string          `json:"param"`
	PK    *rlwe.PublicKey `json:"publicKey"`
}

func (sk *PublicKey) Marshal() ([]byte, error) {
	var src, dst []byte
	var err error
	src, err = json.Marshal(sk)
	if err != nil {
		return dst, err
	}
	base64.RawStdEncoding.Encode(dst, src)
	return dst, err
}

func (sk *PublicKey) Unmarshal(src []byte) error {
	var dst []byte
	var err error
	_, err = base64.RawStdEncoding.Decode(dst, src)
	return err
}

func NewKeyPair() (*SecretKey, *PublicKey, error) {
	p, err := bfv.NewParametersFromLiteral(paramDef[param])
	if err != nil {
		return nil, nil, err
	}

	kgen := bfv.NewKeyGenerator(p)
	sk, pk := kgen.GenKeyPair()
	return &SecretKey{Param: param, SK: sk}, &PublicKey{Param: param, PK: pk}, nil
}
