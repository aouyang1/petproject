package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	pf "github.com/aouyang1/go-petfinder/petfinder"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	rpcRequestCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_client_requests_total",
			Help: "Counter for number of requests made by a client to a downstream.",
		},
		[]string{"app", "code", "instance", "method", "handler"},
	)
	rpcDBCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_requests_total",
			Help: "Counter for number of db call made to the database.",
		},
		[]string{"app", "status", "database", "method"},
	)
	rpcReqRespTimeSummaryVec = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "request_response_times",
			Help:       "Summary of downstream request response times in milliseconds.",
			Objectives: map[float64]float64{0.5: 0.05, 0.90: 0.10, 0.95: 0.05, 0.99: 0.001},
		},
		[]string{"app", "instance", "method", "handler"},
	)
	rpcDBRespTimeSummaryVec = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "db_response_times",
			Help:       "Summary of downstream database response times in milliseconds.",
			Objectives: map[float64]float64{0.5: 0.05, 0.90: 0.10, 0.95: 0.05, 0.99: 0.001},
		},
		[]string{"app", "database", "method"},
	)
)

func init() {
	prometheus.MustRegister(
		rpcRequestCounterVec,
		rpcDBCounterVec,
		rpcReqRespTimeSummaryVec,
		rpcDBRespTimeSummaryVec,
	)
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/metrics", func(c *gin.Context) {
		h := promhttp.Handler()
		h.ServeHTTP(c.Writer, c.Request)
	})

	apiKey := os.Getenv("PETFINDER_API_KEY")
	if apiKey == "" {
		panic("Could not get petfinder api key from environment variable, PETFINDER_API_KEY")
	}

	pfClient := pf.NewClient(apiKey)

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for {
		err = db.Ping()
		if err != nil {
			fmt.Printf("Trying to connect to postgres db with postgres info, %s\n", psqlInfo)
			time.Sleep(time.Duration(3) * time.Second)
		} else {
			break
		}
	}

	go func() {
		startBreedCollection(pfClient, db)
	}()

	r.Run()
}
