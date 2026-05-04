package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
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

	adapter := echoadapter.New(e)
	lambda.Start(adapter.ProxyWithContext)
}
