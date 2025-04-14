package importh_test

import (
	"bytes"
	"encoding/json"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/aviseu/jobs-protobuf/build/gen/commands/jobs"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	oghttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HandlerSuite))
}

type HandlerSuite struct {
	suite.Suite
}

type pubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func (suite *HandlerSuite) Test_Import_Success() {
	// Prepare
	chID := uuid.New()
	iID := uuid.New()

	dsl := testutils.NewDSL(
		testutils.WithArbeitnowEnabled(),
		testutils.WithChannel(
			testutils.WithChannelID(chID),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
		),
		testutils.WithImport(
			testutils.WithImportID(iID),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusPending),
		),
	)

	data, err := proto.Marshal(&jobs.ExecuteImportChannel{
		ImportId: iID.String(),
	})
	suite.NoError(err)
	msg := &pubSubMessage{
		Message: struct {
			Data []byte `json:"data,omitempty"`
			ID   string `json:"id"`
		}{
			Data: data,
			ID:   "1",
		},
	}
	msgJson, err := json.Marshal(msg)
	suite.NoError(err)

	req, err := oghttp.NewRequest("POST", "/import", bytes.NewBuffer(msgJson))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.ImportServer.ServeHTTP(rr, req)

	// Assert response
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Empty(rr.Body.String())

	// Assert state change
	suite.Len(dsl.Imports(), 1)
	i := dsl.FirstImport()
	suite.NotNil(i)
	suite.Equal(iID, i.ID)
	suite.Equal(chID, i.ChannelID)
	suite.Equal(aggregator.ImportStatusCompleted, i.Status)

	// Assert log
	lines := dsl.LogLines()
	suite.Len(lines, 3)
	suite.Contains(lines[0], `"level":"INFO"`)
	suite.Contains(lines[0], "received message!")
	suite.Contains(lines[1], `"level":"INFO"`)
	suite.Contains(lines[1], "processing import "+iID.String())
	suite.Contains(lines[2], `"level":"INFO"`)
	suite.Contains(lines[2], "completed import "+iID.String())
}

func (suite *HandlerSuite) Test_Import_ServerFail() {
	// Prepare
	iID := uuid.New()
	chID := uuid.MustParse(testutils.ArbeitnowMethodNotFound)
	dsl := testutils.NewDSL(
		testutils.WithArbeitnowEnabled(),
		testutils.WithChannel(
			testutils.WithChannelID(chID),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
		),
		testutils.WithImport(
			testutils.WithImportID(iID),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusPending),
		),
	)

	data, err := proto.Marshal(&jobs.ExecuteImportChannel{
		ImportId: iID.String(),
	})
	suite.NoError(err)
	msg := &pubSubMessage{
		Message: struct {
			Data []byte `json:"data,omitempty"`
			ID   string `json:"id"`
		}{
			Data: data,
			ID:   "1",
		},
	}
	msgJson, err := json.Marshal(msg)
	suite.NoError(err)

	req, err := oghttp.NewRequest("POST", "/import", bytes.NewBuffer(msgJson))
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	dsl.ImportServer.ServeHTTP(rr, req)

	// Assert response
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("skipped message\n", rr.Body.String())

	// Assert state change
	suite.Len(dsl.Imports(), 1)
	i := dsl.FirstImport()
	suite.NotNil(i)
	suite.Equal(iID, i.ID)
	suite.Equal(chID, i.ChannelID)
	suite.Equal(aggregator.ImportStatusFailed, i.Status)
	suite.True(i.Error.Valid)
	suite.Contains(i.Error.String, "failed to get jobs page 1 on channel")
	suite.Contains(i.Error.String, "<title>An Error Occurred: Method Not Allowed</title>")

	// Assert log
	lines := dsl.LogLines()
	suite.Len(lines, 3)
	suite.Contains(lines[0], `"level":"INFO"`)
	suite.Contains(lines[0], "received message!")
	suite.Contains(lines[1], `"level":"INFO"`)
	suite.Contains(lines[1], "processing import "+iID.String())
}
