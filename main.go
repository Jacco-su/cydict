package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	VERSION = "0.1.0"
	KEY = "3975l6lr5pcbvidl6jl2"
)
var (
	Verbose bool
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

func main() {

	if len(os.Args) == 1 {
		ShowUsage()
		return 
	}
	
	var rootCmd = &cobra.Command{
		Use: "cydict",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 && args[0] == "help" {
				if err := cmd.Usage(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return
			}

			if Verbose {
				ShowVersion()
				return 
			}

			queryWords := strings.Join(args, " ")
			tranlate(queryWords, "auto2zh")
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")


	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// 请求连接并返回响应
func tranlate(source, direction string) {
	cyurl := "http://api.interpreter.caiyunai.com/v1/translator"
	token := "token " + KEY
	client := &http.Client{}

	type Payload struct {
		Source    string `json:"source"`
		TransType string `json:"trans_type"`
		RequestID string `json:"request_id"`
		Detect    bool   `json:"detect"`
	}
	payload := Payload{source, direction, "demo", true}
	pd, _ := json.Marshal(payload)
	pdString := strings.NewReader(string(pd))
	req, _ := http.NewRequest("POST", cyurl, pdString)
	//设置header
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-authorization", token)

	res, _ := client.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	type TransResult struct {
		Confidence float32 `json: "confidence"`
		Target     string  `json: "target"`
		Isdict     int     `json: "isdict"`
		Rc         int     `json: "rc"`
	}

	var transResult TransResult
	err := json.Unmarshal(body, &transResult)
	if err != nil {
		return
	}

	translation := transResult.Target
	translation = fmt.Sprintf(ErrorColor, translation)
	fmt.Printf(">>  %s: %s\n\n", source, translation)
}

// 打印版本号
func ShowVersion() {
	fmt.Printf("cydict %s\n", VERSION)
}

func ShowUsage() {
	usage := `cy <words>
Query words meanings via the command line.

Example:
words could be word or sentence.

cy hello
cy php is the best language in the world

Usage:
cy <words>...
cy -h | --help
cy -v | --version

Options:
-h --help         show this help message and exit.
-v  --version     displays the current version of cydict.`
	fmt.Printf(usage)
}


