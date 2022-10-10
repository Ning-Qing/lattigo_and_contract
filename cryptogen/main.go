package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tuneinsight/lattigo/v3/bfv"
)

var (
	param string
)

var paramDef = map[string]bfv.ParametersLiteral{
	"PN12QP109": bfv.PN12QP109,
	"PN13QP218": bfv.PN13QP218,
	"PN14QP438": bfv.PN14QP438,
}

var cmd = &cobra.Command{
	Use:     "cryptogen",
	Short:   "cryptogen",
	Version: "1.0.0",
}

func init() {
	cmd.PersistentFlags().StringVarP(&param, "param", "p", "PN13QP218", "BFV param")
	initGenkey()
	cmd.AddCommand(genkey)
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
