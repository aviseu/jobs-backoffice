package testutils

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PubSubSuite struct {
	suite.Suite

	container testcontainers.Container

	client    *pubsub.Client
	badClient *pubsub.Client

	ImportTopic        *pubsub.Topic
	BadImportTopic     *pubsub.Topic
	ImportSubscription *pubsub.Subscription

	JobTopic        *pubsub.Topic
	BadJobTopic     *pubsub.Topic
	JobSubscription *pubsub.Subscription
}

func (suite *PubSubSuite) SetupSuite() {
	suite.container, suite.client, suite.badClient, suite.ImportTopic, suite.BadImportTopic, suite.ImportSubscription, suite.JobTopic, suite.BadJobTopic, suite.JobSubscription = suite.createDependencies(context.Background())
}

func (suite *PubSubSuite) createDependencies(ctx context.Context) (testcontainers.Container, *pubsub.Client, *pubsub.Client, *pubsub.Topic, *pubsub.Topic, *pubsub.Subscription, *pubsub.Topic, *pubsub.Topic, *pubsub.Subscription) {
	req := testcontainers.ContainerRequest{
		Image:        "google/cloud-sdk:emulators",
		ExposedPorts: []string{"8085/tcp"},
		Cmd:          []string{"gcloud", "beta", "emulators", "pubsub", "start", "--host-port=0.0.0.0:8085"},
		WaitingFor:   wait.ForLog("Server started").WithStartupTimeout(10 * time.Second),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	suite.NoError(err)

	host, err := c.Host(ctx)
	suite.NoError(err)

	port, err := c.MappedPort(ctx, "8085")
	suite.NoError(err)

	client, err := pubsub.NewClient(
		ctx,
		"test-project",
		option.WithEndpoint(fmt.Sprintf("%s:%s", host, port.Port())),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		option.WithoutAuthentication(),
		option.WithTelemetryDisabled(),
		internaloption.SkipDialSettingsValidation(),
	)
	suite.NoError(err)

	badClient, err := pubsub.NewClient(
		ctx,
		"test-project",
		option.WithEndpoint(fmt.Sprintf("%s:%s", host, port.Port())),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		option.WithoutAuthentication(),
		option.WithTelemetryDisabled(),
		internaloption.SkipDialSettingsValidation(),
	)
	suite.NoError(err)

	importTopic, err := client.CreateTopic(ctx, "import-test-topic")
	suite.NoError(err)

	badImportTopic, err := badClient.CreateTopic(ctx, "import-bad-test-topic")
	suite.NoError(err)

	importSubscription, err := client.CreateSubscription(ctx, "import-test-sub", pubsub.SubscriptionConfig{
		Topic: importTopic,
	})
	suite.NoError(err)

	jobTopic, err := client.CreateTopic(ctx, "job-test-topic")
	suite.NoError(err)

	badJobTopic, err := badClient.CreateTopic(ctx, "job-bad-test-topic")
	suite.NoError(err)

	jobSubscription, err := client.CreateSubscription(ctx, "job-test-sub", pubsub.SubscriptionConfig{
		Topic: jobTopic,
	})
	suite.NoError(err)

	suite.NoError(badClient.Close())

	return c, client, badClient, importTopic, badImportTopic, importSubscription, jobTopic, badJobTopic, jobSubscription
}

func (suite *PubSubSuite) TearDownSuite() {
	go func() {
		suite.NoError(suite.client.Close())
		suite.NoError(suite.container.Terminate(context.Background()))
	}()
}
