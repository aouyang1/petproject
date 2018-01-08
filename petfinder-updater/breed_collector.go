package main

import (
	"database/sql"
	"log"
	"time"

	pf "github.com/aouyang1/go-petfinder/petfinder"
	"github.com/prometheus/client_golang/prometheus"
)

func startBreedCollection(pfClient pf.Client, db *sql.DB) {
	var err error
	var breeds pf.Breeds
	var counter prometheus.Counter
	var timer *prometheus.Timer
	var summary prometheus.Observer
	var currentTime time.Time

	// Define animals to collect
	animals := []string{"barnyard", "bird", "cat", "dog", "horse", "reptile", "smallfurry"}

	requestLabel := prometheus.Labels{
		"app":      "petfinder-updator",
		"method":   "get",
		"instance": "api.petfinder.com",
		"handler":  "/breed.list",
	}

	requestSummaryLabel := prometheus.Labels{
		"app":      "petfinder-updator",
		"method":   "get",
		"instance": "api.petfinder.com",
		"handler":  "/breed.list",
	}

	dbLabel := prometheus.Labels{
		"app":      "petfinder-updator",
		"database": "petfinder",
	}

	dbSummaryLabel := prometheus.Labels{
		"app":      "petfinder-updator",
		"database": "petfinder",
	}

	for {
		for _, animal := range animals {
			summary, err = rpcReqRespTimeSummaryVec.GetMetricWith(requestSummaryLabel)
			if err != nil {
				log.Printf("Error retrieving summary for breed list, %+v\n", err)
				continue
			}

			timer = prometheus.NewTimer(summary)
			breeds, err = pfClient.ListBreeds(pf.Options{Animal: animal})
			timer.ObserveDuration()

			if err != nil {
				log.Printf("Error getting breed for animal, %s, %v\n", animal, err)
				requestLabel["code"] = "500"
			} else {
				requestLabel["code"] = "200"
			}

			counter, err = rpcRequestCounterVec.GetMetricWith(requestLabel)
			if err != nil {
				log.Printf("Error retrieving counter for breed list, %+v\n", err)
				continue
			}

			counter.Add(1)

			currentTime = time.Now().UTC()
			sqlStatement := `
INSERT INTO breed (name, animal, created_on, updated_on)
VALUES ($1, $2, $3, $4)
ON CONFLICT ON CONSTRAINT unq_name_animal DO UPDATE
SET updated_on = $4`

			stmt, err := db.Prepare(sqlStatement)
			if err != nil {
				log.Printf("Failed to prepare sql statement, %s\n", sqlStatement)
				return
			}

			dbSummaryLabel["method"] = "write"
			dbLabel["method"] = "write"
			for _, name := range breeds {
				summary, err = rpcDBRespTimeSummaryVec.GetMetricWith(dbSummaryLabel)
				if err != nil {
					log.Printf("Error retrieving summary for breed list, %+v\n", err)
					continue
				}

				timer = prometheus.NewTimer(summary)
				_, err = stmt.Exec(name, animal, currentTime, currentTime)
				timer.ObserveDuration()

				if err != nil {
					log.Println("Failed to insert breed, animal into db, %+v", err)
					dbLabel["status"] = "error"
					counter, err = rpcDBCounterVec.GetMetricWith(dbLabel)
					if err != nil {
						log.Println("Error retrieving db counter for breed list updater, %+v", err)
						continue
					}
					counter.Add(1)
					continue
				}

				dbLabel["status"] = "success"
				counter, err = rpcDBCounterVec.GetMetricWith(dbLabel)
				if err != nil {
					log.Println("Error retrieving db counter for breed list updater, %+v", err)
				}
				counter.Add(1)
			}

			time.Sleep(time.Duration(1) * time.Second)
		}
		time.Sleep(time.Duration(60) * time.Second)
	}
}
