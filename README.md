# anticaptcha [![PkgGoDev](https://pkg.go.dev/badge/github.com/aidenesco/anticaptcha)](https://pkg.go.dev/github.com/aidenesco/anticaptcha) [![Go Report Card](https://goreportcard.com/badge/github.com/aidenesco/anticaptcha)](https://goreportcard.com/report/github.com/aidenesco/anticaptcha)
This package is an anti-captcha.net client library. See official documentation [here](https://anticaptcha.atlassian.net/wiki/spaces/API/pages/196635/Documentation+in+English)

## Installation
```sh
go get -u github.com/aidenesco/anticaptcha
```

## Usage

```go
import "github.com/aidenesco/anticaptcha"

func main() {
    client := anticaptcha.NewClient("your-key-here")
    
    balance, _ := client.GetBalance()
    
    fmt.Println(balance) // 4.77

}
```
