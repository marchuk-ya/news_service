package setup

import (
	"context"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient   *mongo.Client
	pool          *dockertest.Pool
	mongoResource *dockertest.Resource
)

// SetupTestMongoDB initializes MongoDB for testing
func SetupTestMongoDB(t *testing.T) (*mongo.Client, func()) {
	t.Helper()

	if MongoClient == nil {
		var err error
		pool, err = dockertest.NewPool("")
		if err != nil {
			t.Fatalf("Could not connect to docker: %s", err)
		}

		pool.MaxWait = time.Minute * 2

		mongoResource, err = pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "mongo",
			Tag:        "latest",
			Env: []string{
				"MONGO_INITDB_ROOT_USERNAME=root",
				"MONGO_INITDB_ROOT_PASSWORD=example",
			},
		}, func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})
		if err != nil {
			t.Fatalf("Could not start MongoDB: %s", err)
		}

		if err := pool.Retry(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			clientOptions := options.Client().
				SetHosts([]string{mongoResource.GetHostPort("27017/tcp")}).
				SetAuth(options.Credential{
					Username: "root",
					Password: "example",
				})

			MongoClient, err = mongo.Connect(ctx, clientOptions)
			if err != nil {
				return err
			}

			return MongoClient.Ping(ctx, nil)
		}); err != nil {
			t.Fatalf("Could not connect to MongoDB: %s", err)
		}
	}

	cleanup := func() {
		ctx := context.Background()
		if err := MongoClient.Database("test_db").Drop(ctx); err != nil {
			t.Logf("Failed to drop test database: %v", err)
		}
	}

	return MongoClient, cleanup
}
