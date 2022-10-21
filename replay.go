package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	Replay.Flags().String("addrs", "", "The addresses to hit (with path), comma separated")
	Replay.Flags().String("input", "", "A path to the captured JSON requests to replay.")
	Replay.Flags().String("output", "", "The path to the output results file")
}

type replayRequest struct {
    ReqBody string `json:"req_body"`
}

// Replay replays requests from a given model to another model.
var Replay = &cobra.Command{
	Use:   "replay",
	Short: "Replay a request to another instance of a model (local or remote)",
	Run: func(cmd *cobra.Command, args []string) {
		inputFile, err := cmd.Flags().GetString("input")
		if err != nil || inputFile == "" {
			fmt.Println("`mlm replay` requires flag --input")
			os.Exit(1)
		}
		outputFile, err := cmd.Flags().GetString("output")
		if err != nil || outputFile == "" {
			fmt.Println("`mlm replay` requires flag --output")
			os.Exit(1)
		}
		addrsStr, err := cmd.Flags().GetString("addrs")
		if err != nil || addrsStr == "" {
			fmt.Println("`mlm replay` requires flag --addrs")
			os.Exit(1)
		}
		addrs := strings.Split(addrsStr, ",")
		fmt.Println(addrs)

		in, err := os.Open(inputFile)
		if err != nil {
			fmt.Printf("Could not open input file %s: %v\n", inputFile, err)
			os.Exit(1)
		}
		outFile, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("Could not open output file %s: %v\n", outputFile, err)
			os.Exit(1)
		}
		defer outFile.Close()


		inputReader := bufio.NewReader(in)
		decoder := json.NewDecoder(inputReader)

		reqIdx := 0
		for decoder.More() {
		    var request replayRequest
		    if err := decoder.Decode(&request); err != nil {
		        fmt.Printf("Error parsing replay request: %v\n", err)
		        os.Exit(1)
		    }

		    resStr, err := performReplays(request, addrs, reqIdx)
		    if err != nil {
		    	fmt.Printf("Error performing replays: %v\n", err)
		    	os.Exit(1)
		    }

		    outFile.WriteString(resStr)
		    outFile.WriteString("\n")
		    reqIdx++
		}

		fmt.Printf("Wrote %d replays to %s", reqIdx, outputFile)
	},	
}