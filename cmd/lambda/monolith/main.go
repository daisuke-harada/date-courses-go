package main

import (
	"log"

	iface "github.com/daisuke-harada/date-courses-go/internal/interface"
	"github.com/daisuke-harada/date-courses-go/pkg/logger"
)

func main() {
	logger.Init("date-courses-go", false)
	defer logger.Close()

	e, err := iface.NewEchoApp()
	if err != nil {
		log.Fatal(err)
	}

	// TODO(Step 3): github.com/awslabs/aws-lambda-go-api-proxy/echo を使って
	// Lambda ハンドラーとして起動する。現在はスケルトン。
	_ = e
}
