package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/omarahr/huffman-encoder-decoder/encoder"
	"github.com/spf13/cobra"
)

type State struct {
	InputFilePath  string
	OutputFilePath string
	DecodeMode     bool
}

var (
	state = State{
		InputFilePath:  "",
		OutputFilePath: "",
	}
)

func timeElapsed(f func()) {
	currenTime := time.Now()
	f()
	elapsedTime := time.Since(currenTime)
	fmt.Printf("Elapsed time: %s\n", elapsedTime)
}

func handleOutputFileName(outputFilePath string) {
	if outputFilePath == "" {
		if state.DecodeMode {
			state.OutputFilePath = state.InputFilePath + ".decomp.txt"
		} else {
			state.OutputFilePath = state.InputFilePath + ".comp"
		}
	} else {
		state.OutputFilePath = outputFilePath
	}
}

var rootCmd = &cobra.Command{
	Use:   "coco",
	Short: "coco is a encoder / decoder tool",
	Long:  `coco is a encoder / decoder tool`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Please input a file path")
			return
		}

		// get input file
		state.InputFilePath = args[0]

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			log.Fatal(err)
			return
		}

		handleOutputFileName(output)

		if state.InputFilePath == state.OutputFilePath {
			log.Fatal("Input file path and output file path can't be same")
			return
		}

		if state.DecodeMode {
			timeElapsed(func() {
				encoder.Decompress(state.InputFilePath, state.OutputFilePath)
			})
		} else {
			timeElapsed(func() {
				encoder.Compress(state.InputFilePath, state.OutputFilePath)
			})
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("output", "o", "", "output file path")

	rootCmd.PersistentFlags().BoolVarP(
		&state.DecodeMode,
		"decode mode",
		"d",
		false,
		"decode mode",
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
