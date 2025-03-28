package importh_test

import (
	"bytes"
	"encoding/json"
	"github.com/aviseu/jobs-backoffice/internal/app/application/http"
	"github.com/aviseu/jobs-backoffice/internal/app/domain"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/gateway"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/job"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/api/arbeitnow"
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
	server := testutils.NewArbeitnowServer()
	lbuf, log := testutils.NewLogger()
	chr := testutils.NewChannelRepository()
	chs := configuring.NewService(chr)
	ir := testutils.NewImportRepository()
	is := imports.NewService(ir)
	jr := testutils.NewJobRepository()
	js := job.NewService(jr, 10, 10)
	f := gateway.NewFactory(
		js,
		is,
		oghttp.DefaultClient,
		gateway.Config{
			Arbeitnow: arbeitnow.Config{URL: server.URL},
			Import: struct {
				ResultBufferSize int `split_words:"true" default:"10"`
				ResultWorkers    int `split_words:"true" default:"10"`
			}{
				ResultBufferSize: 10,
				ResultWorkers:    10,
			},
		},
		log,
	)
	ia := domain.NewImportAction(chs, is, f)
	h := http.ImportRootHandler(ia, log)

	chID := uuid.New()
	chr.Add(configuring.New(chID, "Channel Name", base.IntegrationArbeitnow, base.ChannelStatusActive).DTO())

	iID := uuid.New()
	ir.Add(imports.New(iID, chID))

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
	h.ServeHTTP(rr, req)

	// Assert response
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Empty(rr.Body.String())

	// Assert state change
	suite.Len(ir.Imports, 1)
	var i *imports.Import
	for _, v := range ir.Imports {
		i = v
	}
	suite.NotNil(i)
	suite.Equal(iID, i.ID())
	suite.Equal(chID, i.ChannelID())
	suite.Equal(imports.StatusCompleted, i.Status())

	// Assert log
	lines := testutils.LogLines(lbuf)
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
	server := testutils.NewArbeitnowServer()
	lbuf, log := testutils.NewLogger()
	chr := testutils.NewChannelRepository()
	chs := configuring.NewService(chr)
	ir := testutils.NewImportRepository()
	is := imports.NewService(ir)
	jr := testutils.NewJobRepository()
	js := job.NewService(jr, 10, 10)
	f := gateway.NewFactory(
		js,
		is,
		oghttp.DefaultClient,
		gateway.Config{
			Arbeitnow: arbeitnow.Config{URL: server.URL},
			Import: struct {
				ResultBufferSize int `split_words:"true" default:"10"`
				ResultWorkers    int `split_words:"true" default:"10"`
			}{
				ResultBufferSize: 10,
				ResultWorkers:    10,
			},
		},
		log,
	)
	ia := domain.NewImportAction(chs, is, f)
	h := http.ImportRootHandler(ia, log)

	chID := uuid.MustParse(testutils.ArbeitnowMethodNotFound)
	chr.Add(configuring.New(chID, "Channel Name", base.IntegrationArbeitnow, base.ChannelStatusActive).DTO())

	iID := uuid.New()
	ir.Add(imports.New(iID, chID))

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
	h.ServeHTTP(rr, req)

	// Assert response
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("skipped message\n", rr.Body.String())

	// Assert state change
	suite.Len(ir.Imports, 1)
	var i *imports.Import
	for _, v := range ir.Imports {
		i = v
	}
	suite.NotNil(i)
	suite.Equal(iID, i.ID())
	suite.Equal(chID, i.ChannelID())
	suite.Equal(imports.StatusFailed, i.Status())
	suite.True(i.Error().Valid)
	suite.Contains(i.Error().String, "failed to get jobs page 1 on channel")
	suite.Contains(i.Error().String, "<title>An Error Occurred: Method Not Allowed</title>")

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 3)
	suite.Contains(lines[0], `"level":"INFO"`)
	suite.Contains(lines[0], "received message!")
	suite.Contains(lines[1], `"level":"INFO"`)
	suite.Contains(lines[1], "processing import "+iID.String())
}
