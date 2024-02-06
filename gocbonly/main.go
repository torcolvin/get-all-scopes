package main

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

const (
	cbsUserName = "Administrator"
	cbsPassword = "password"
	cbsAddr     = "localhost"
	bucketName  = "bucket1"
)

func main() {
	//gocb.SetLogger(gocb.VerboseStdioLogger())
	log.Println("Start test program")
	cluster := getCluster()
	defer func() {
		err := cluster.Close(nil)
		if err != nil {
			log.Fatalf("Error closing cluster %+v", err)
		}
	}()
	log.Println("Creating bucket")
	err := createBucket(cluster, bucketName)
	if err != nil {
		log.Fatalf("Error creating bucket %s", err)
	}
	log.Println("Created bucket")
	bucket := cluster.Bucket(bucketName)
	log.Println("Wait for bucket to be ready")
	err = bucket.WaitUntilReady(10*time.Second, &gocb.WaitUntilReadyOptions{
		RetryStrategy: gocb.NewBestEffortRetryStrategy(nil),
	})
	if err != nil {
		log.Fatalf("Error waiting for bucket to be ready %s", err)
	}
	log.Println("Done wait for bucket to be ready")
	fmt.Println("Start GetAllScopes")
	scopes, err := bucket.Collections().GetAllScopes(nil)
	if err != nil {
		log.Fatalf("Error getting scopes %s", err)
	}
	fmt.Println("Scopes:", scopes)

	fmt.Println("Start flush")
	err = cluster.Buckets().FlushBucket(bucketName, nil)
	if err != nil {
		log.Fatalf("Error flushing bucket %s", err)
	}

	time.Sleep(3 * time.Second)
	fmt.Println("Start GetAllScopes after flush")
	scopes, err = bucket.Collections().GetAllScopes(nil)
	if err != nil {
		log.Fatalf("Error getting scopes %s", err)
	}
	fmt.Println("Scopes after flush:", scopes)

	fmt.Println("SUCCESS")
}

func getCluster() *gocb.Cluster {
	DefaultGocbV2OperationTimeout := 10 * time.Second

	clusterOptions := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: cbsUserName,
			Password: cbsPassword,
		},
		SecurityConfig: gocb.SecurityConfig{
			TLSSkipVerify: false,
		},
		TimeoutsConfig: gocb.TimeoutsConfig{
			ConnectTimeout:    DefaultGocbV2OperationTimeout,
			KVTimeout:         DefaultGocbV2OperationTimeout,
			ManagementTimeout: DefaultGocbV2OperationTimeout,
			QueryTimeout:      90 * time.Second,
			ViewTimeout:       90 * time.Second,
		},
	}
	connStr := fmt.Sprintf("couchbase://%s?idle_http_connection_timeout=90000&kv_pool_size=2&max_idle_http_connections=64000&max_perhost_idle_http_connections=256", cbsAddr)
	cluster, err := gocb.Connect(connStr, clusterOptions)
	if err != nil {
		log.Fatalf("Error connecticting to cluster %+v", err)
	}
	err = cluster.WaitUntilReady(15*time.Second, nil)
	if err != nil {
		log.Fatalf("Can't connect to cluster %+v", err)
	}

	err = cluster.WaitUntilReady(90*time.Second,
		&gocb.WaitUntilReadyOptions{ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeQuery}},
	)
	if err != nil {
		log.Fatalf("Query service not online")
	}
	return cluster
}

func createBucket(cluster *gocb.Cluster, name string) error {
	quotaMB := 200
	settings := gocb.CreateBucketSettings{
		BucketSettings: gocb.BucketSettings{
			Name:         name,
			RAMQuotaMB:   uint64(quotaMB),
			BucketType:   gocb.CouchbaseBucketType,
			FlushEnabled: true,
			NumReplicas:  0,
		},
	}

	options := &gocb.CreateBucketOptions{
		Timeout: 10 * time.Second,
	}
	return cluster.Buckets().CreateBucket(settings, options)
}
