package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

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

type replayResultInstance struct {
	RespBody  string `json:"resp_body"`
	LatencyMs int64  `json:"latency_ms"`
}

type replayResult struct {
	Request int                             `json:"request_idx"`
	Results map[string]replayResultInstance `json:"results"`
}

func performReplay(r replayRequest, addr string) (replayResultInstance, error) {
	request, err := http.NewRequest("POST", addr, bytes.NewBuffer([]byte(r.ReqBody)))
	if err != nil {
		fmt.Println(err)
		return replayResultInstance{}, nil
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	start := time.Now()
	response, err := client.Do(request)
	if err != nil {
		return replayResultInstance{}, err
	}
	end := time.Now()
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return replayResultInstance{
		LatencyMs: end.Sub(start).Milliseconds(),
		RespBody:  string(body),
	}, nil
}

func performReplays(r replayRequest, addrs []string, reqIdx int) (string, error) {
	resMap := make(map[string]replayResultInstance, len(addrs))
	for _, addr := range addrs {
		res, err := performReplay(r, addr)
		if err != nil {
			return "", fmt.Errorf("Error for addr %s: %v", addr, err)
		}
		resMap[addr] = res
	}
	out := replayResult{
		Request: reqIdx,
		Results: resMap,
	}
	resBytes, err := json.Marshal(out)
	if err != nil {
		return "", err
	}
	return string(resBytes), nil
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
