## Usage
### Quick start
```go
func main() {
	app := sdk.NewApp("appkey", "appsecret")
	// specified server addr
	// app := sdk.NewAppWithAddr("http://127.0.0.1:8092", "appkey", "appsecret")

	result, _ := app.CreateAddress("ETH")
	printResult(result)
	
	result, _ = app.CreateAddressWithMode("ETH", "auto")
	printResult(result)

	result, _ = app.VerifyAddress("ETH", "0x9bf65CDF5729b9588F6bAEBb2Aa2926472D4a035")
	printResult(result)

	result, _ = app.GetAssets()
	printResult(result)

	result, _ = app.GetBalance("ETH")
	printResult(result)

	result, _ = app.GetOrder("rNXBQGJlw09apVyg4nDo")
	printResult(result)

	result, _ = app.Withdraw("1569225735", "ETH", "0xF0706B7Cab38EA42538f4D8C279B6F57ad1d4072", "0.05")
	printResult(result)

	result, _ = app.GetValidators("IRIS2")
	printResult(result)

	result, _ = app.Delegate("1569231558", "IRIS2", "1")
	printResult(result)

	result, _ = app.UnDelegate("1569231809", "IRIS2", "1")
	printResult(result)

	result, _ = app.AddUrgentStakingFunding("1569292076", "IRIS2", "1", 1569302076)
	printResult(result)

	company := sdk.NewCompany("appkey", "appsecret")
	// specified server addr
	// company := sdk.NewCompanyWithAddr("appkey", "appsecret", "http://127.0.0.1:8092")

	result, _ = company.GetFundingWallets()
	printResult(result)

	result, _ = company.FundingTransfer("QrLxg3XgKPMR1O8", "ZKz8XpwXGPBAQVE", "BTC", "0.0001")
	printResult(result)

	result, _ = company.GetFundingRecords(1, 10)
	printResult(result)
}
    
func printResult(result *sdk.Result) {
	fmt.Println("code:", result.Code)
	fmt.Println("message:", result.Message)
	fmt.Println("data:")

	b, err := json.MarshalIndent(result.Data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}
```

## CLI
Usage:
`ctl <key> <secret> <action> [<params>...] [-a <host>]`

e.g.
```bash
# create address
$ go run cmd/ctl/main.go "appkey" "appsecret" "CreateAddress" "ETH"
# verify address
$ go run cmd/ctl/main.go "appkey" "appsecret" "VerifyAddress" "ETH" "0x9bf65CDF5729b9588F6bAEBb2Aa2926472D4a035"
# get assets
$ go run cmd/ctl/main.go "appkey" "appsecret" "GetAssets"
# get balance
$ go run cmd/ctl/main.go "appkey" "appsecret" "GetBalance" "ETH"
# get order
$ go run cmd/ctl/main.go "appkey" "appsecret" "GetOrder" "rNXBQGJlw09apVyg4nDo"
# withdraw
$ go run cmd/ctl/main.go "appkey" "appsecret" "Withdraw" "$(date +%s)" "ETH" "0xF0706B7Cab38EA42538f4D8C279B6F57ad1d4072" "0.05"
# get validator
$ go run cmd/ctl/main.go "appkey" "appsecret" "GetValidators" "IRIS2"
# delegate
$ go run cmd/ctl/main.go "appkey" "appsecret" "Delegate" $(date +%s) "IRIS2" "1"
# undelegate
$ go run cmd/ctl/main.go "appkey" "appsecret" "UnDelegate" $(date +%s) "IRIS2" "1"
# get staking interest
$ go run cmd/ctl/main.go "appkey" "appsecret" "GetStakingInterest" "IRIS2" "2019-09-26"
# urgent staking funding
$ go run cmd/ctl/main.go "appkey" "appsecret" "AddUrgentStakingFunding" $(date +%s) "IRIS2" "1" "1569302076"
# get funding wallets
$ go run cmd/ctl/main.go "companykey" "companysecret" "GetFundingWallets"
# funding transfer
$ go run cmd/ctl/main.go "companykey" "companysecret" "FundingTransfer" "QrLxg3XgKPMR1O8" "ZKz8XpwXGPBAQVE" "BTC" "0.0001"
# get funding records
$ go run cmd/ctl/main.go "companykey" "companysecret" "GetFundingRecords" 1 10
```
