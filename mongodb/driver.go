package mongodb

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

// MaxPool, RetryTime ...
const (
	MaxPool   = 20
	RetryTime = 5
)

// DBNAME ...
var DBNAME = os.Getenv(`MDB_NAME`) // actual name of Mongo Database Name, local or 'Rewards' for dev/test/prod

// CheckAndInitServiceConnection ...
func CheckAndInitServiceConnection() {
	urlChan := make(chan string)
	go MongoConnect(urlChan)
	service.URL = os.Getenv(`MDB_URL`)
	urlChan <- service.URL
}

// MongoConnect ... go routine to connect to mongo
func MongoConnect(urlChan chan string) {
	serviceURL := <-urlChan
	const mongoStartRetries = 3
	mongoRe := regexp.MustCompile(`^mongodb://(.*)@`)
	url := mongoRe.ReplaceAllLiteralString(serviceURL, "mongodb://username:password@")
	for retries := 0; retries <= mongoStartRetries; retries++ {
		fmt.Printf("Trying to connect to mongo:[%s] for %d time\n", url, retries+1)
		err := service.New() // Attempt to connect to Mongo
		if err == nil {
			log.Printf("GOOD: Connected to mongo [%s] in [%d] tries\n", url, retries+1)
			log.Printf("GOOD: Mongodb Initialized; Session Created...\n")
			return
		} else if retries == mongoStartRetries {
			fmt.Printf("Could not connect to mongo:[%s] after %d retries\n", url, retries+1)
			fmt.Printf(err.Error())
			fmt.Printf("Going to retry with a longer timeout [%dm]\n", RetryTime)
			var round = 0
			for {
				err := service.New()
				if err == nil {
					log.Printf("GOOD: Mongodb Initialized; Session Created...\n")
					return
				}
				time.Sleep(RetryTime * time.Minute)
				fmt.Printf("Retry %v: Failed to create service.New()\n", round)
				round++
			}
		} else {
			fmt.Printf(err.Error())
		}
		time.Sleep(5 * time.Second)
		fmt.Printf("Retry: start mongo service\n")
	}
}
