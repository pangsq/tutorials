package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/pangsq/tutorials/servicemesh/tracing/go-service/tracing"
	"net/http"
	"os"
	"strconv"
)

var nextTarget string
var finalTarget string

func main() {
	var ok bool
	if nextTarget, ok = os.LookupEnv("NEXT_TARGET"); !ok {
		nextTarget = "http://java-service"
	}
	if finalTarget, ok = os.LookupEnv("FINAL_TARGET"); !ok {
		finalTarget = "http://httpbin"
	}
	engine := gin.New()
	engine.GET("/*path", forward)
	server := &http.Server{
		Addr:    fmt.Sprintf(":8080"),
		Handler: engine,
	}
	if err := server.ListenAndServe(); err != nil {
		logrus.Error(err)
	}
}

func forward(ginCtx *gin.Context) {
	// extract headers to get b3 propagation (including trace-id/span-id/sampled)
	// then inject these headers into the request to the next service
	ctx := tracing.ExtractGinCtx(ginCtx)

	// just for debug
	//logrus.Info(ginCtx.Request.Header)
	//logrus.Infof("b3: %s/%s", ginCtx.GetHeader("X-B3-TraceId"), ginCtx.GetHeader("X-B3-SpanId"))

	ttl := ginCtx.Query("ttl")
	ttlNum, err := strconv.Atoi(ttl)
	if err != nil {
		logrus.Warningf("failed to get ttl: %s", err.Error())
		ttlNum = 0
	}

	path := ginCtx.Param("path")
	var requestUrl string
	//logrus.Infof("ttl=%d", ttlNum)
	if ttlNum == 0 {
		requestUrl = finalTarget + path
	} else {
		requestUrl = nextTarget + path
	}
	request, err := http.NewRequestWithContext(ctx, "GET", requestUrl, nil)
	if err != nil {
		logrus.Error(err)
		ginCtx.JSON(http.StatusBadRequest, err)
		return
	}
	if ttlNum != 0 {
		q := request.URL.Query()
		q.Add("ttl", strconv.Itoa(ttlNum - 1))
		request.URL.RawQuery = q.Encode()
	}
	tracing.Inject(ctx, request)
	//logrus.Info(request.Header)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logrus.Error(err)
		ginCtx.JSON(http.StatusBadGateway, err)
		return
	}
	ginCtx.DataFromReader(resp.StatusCode, resp.ContentLength, gin.MIMEJSON, resp.Body, nil)
}
