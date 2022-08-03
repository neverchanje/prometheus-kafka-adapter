// Copyright 2018 Telefónica
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("creating kafka producer")

	var compression kafka.Compression
	switch kafkaCompression {
	case "none":
		compression = compress.None
	case "snappy":
		compression = kafka.Snappy
	default:
		log.Fatalf("Unsupported compression type: %s", kafkaCompression)
	}

	writer := &kafka.Writer{
		Addr:        kafka.TCP(strings.Split(kafkaBrokerList, ",")...),
		Topic:       kafkaTopic,
		BatchSize:   kafkaBatchNumMessages,
		Compression: compression,
	}

	r := gin.New()

	r.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true), gin.Recovery())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "UP"}) })
	if basicauth {
		authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
			basicauthUsername: basicauthPassword,
		}))
		authorized.POST("/receive", receiveHandler(writer, serializer))
	} else {
		r.POST("/receive", receiveHandler(writer, serializer))
	}

	logrus.Fatal(r.Run())
}
