package api_test

import (
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/application/http"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/base"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/imports"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	oghttp "net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestImportHandler(t *testing.T) {
	suite.Run(t, new(ImportHandlerSuite))
}

type ImportHandlerSuite struct {
	suite.Suite
}

func (suite *ImportHandlerSuite) Test_List_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	ir := testutils.NewImportRepository()
	is := imports.NewService(ir)
	chr := testutils.NewChannelRepository()
	chs := configuring.NewService(chr)
	h := http.APIRootHandler(chs, is, nil, http.Config{}, log)

	ch := configuring.New(uuid.New(), "Channel Name", base.IntegrationArbeitnow, base.ChannelStatusActive)
	chr.Add(ch.DTO())

	id1 := uuid.New()
	sAt1 := time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC)
	eAt1 := time.Date(2020, 1, 1, 0, 0, 4, 0, time.UTC)
	i1 := imports.New(
		id1,
		ch.ID(),
		imports.WithStatus(base.ImportStatusPublishing),
		imports.WithMetadata(1, 2, 3, 4, 5),
		imports.WithStartAt(sAt1),
		imports.WithEndAt(eAt1),
		imports.WithError("happened this error"),
	)
	ir.Add(i1.ToDTO())

	id2 := uuid.New()
	sAt2 := time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)
	i2 := imports.New(id2, ch.ID(), imports.WithStartAt(sAt2))
	ir.Add(i2.ToDTO())

	id3 := uuid.New()
	sAt3 := time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)
	i3 := imports.New(id3, ch.ID(), imports.WithStartAt(sAt3))
	ir.Add(i3.ToDTO())

	req, err := oghttp.NewRequest("GET", "/api/imports", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"imports":[{"id":"`+id1.String()+`","channel_id":"`+ch.ID().String()+`","channel_name":"Channel Name","integration":"arbeitnow","status":"publishing","started_at":"2020-01-01T00:00:03Z","ended_at":"2020-01-01T00:00:04Z","error":"happened this error","new_jobs":1,"updated_jobs":2,"no_change_jobs":3,"missing_jobs":4,"failed_jobs":5,"total_jobs":15},{"id":"`+id2.String()+`","channel_id":"`+ch.ID().String()+`","channel_name":"Channel Name","integration":"arbeitnow","status":"pending","started_at":"2020-01-01T00:00:02Z","ended_at":null,"error":null,"new_jobs":0,"updated_jobs":0,"no_change_jobs":0,"missing_jobs":0,"failed_jobs":0,"total_jobs":0},{"id":"`+id3.String()+`","channel_id":"`+ch.ID().String()+`","channel_name":"Channel Name","integration":"arbeitnow","status":"pending","started_at":"2020-01-01T00:00:01Z","ended_at":null,"error":null,"new_jobs":0,"updated_jobs":0,"no_change_jobs":0,"missing_jobs":0,"failed_jobs":0,"total_jobs":0}]}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ImportHandlerSuite) Test_List_RepositoryFail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	ir := testutils.NewImportRepository()
	ir.FailWith(errors.New("boom!"))
	is := imports.NewService(ir)
	chr := testutils.NewChannelRepository()
	chs := configuring.NewService(chr)
	h := http.APIRootHandler(chs, is, nil, http.Config{}, log)

	req, err := oghttp.NewRequest("GET", "/api/imports", nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"error":{"message":"Internal Server Error"}}`+"\n", rr.Body.String())

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], "failed to get imports: boom!")
}

func (suite *ImportHandlerSuite) Test_Find_Success() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	ir := testutils.NewImportRepository()
	is := imports.NewService(ir)
	chr := testutils.NewChannelRepository()
	chs := configuring.NewService(chr)
	h := http.APIRootHandler(chs, is, nil, http.Config{}, log)

	ch := configuring.New(uuid.New(), "Channel Name", base.IntegrationArbeitnow, base.ChannelStatusActive)
	chr.Add(ch.DTO())

	id := uuid.New()
	sAt := time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC)
	eAt := time.Date(2020, 1, 1, 0, 0, 2, 0, time.UTC)
	i := imports.New(
		id,
		ch.ID(),
		imports.WithStatus(base.ImportStatusCompleted),
		imports.WithMetadata(1, 2, 3, 4, 5),
		imports.WithStartAt(sAt),
		imports.WithEndAt(eAt),
		imports.WithError("happened this error"),
	)
	ir.Add(i.ToDTO())

	req, err := oghttp.NewRequest("GET", "/api/imports/"+id.String(), nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusOK, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"id":"`+id.String()+`","channel_id":"`+ch.ID().String()+`","channel_name":"Channel Name","integration":"arbeitnow","status":"completed","started_at":"2020-01-01T00:00:01Z","ended_at":"2020-01-01T00:00:02Z","error":"happened this error","new_jobs":1,"updated_jobs":2,"no_change_jobs":3,"missing_jobs":4,"failed_jobs":5,"total_jobs":15}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ImportHandlerSuite) Test_Find_NotFound() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	ir := testutils.NewImportRepository()
	is := imports.NewService(ir)
	chr := testutils.NewChannelRepository()
	chs := configuring.NewService(chr)
	h := http.APIRootHandler(chs, is, nil, http.Config{}, log)

	id := uuid.New()
	req, err := oghttp.NewRequest("GET", "/api/imports/"+id.String(), nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusNotFound, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"error":{"message":"import not found"}}`+"\n", rr.Body.String())

	// Assert log
	suite.Empty(lbuf.String())
}

func (suite *ImportHandlerSuite) Test_Find_RepositoryFail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	ir := testutils.NewImportRepository()
	ir.FailWith(errors.New("boom!"))
	is := imports.NewService(ir)
	chr := testutils.NewChannelRepository()
	chs := configuring.NewService(chr)
	h := http.APIRootHandler(chs, is, nil, http.Config{}, log)

	id := uuid.New()
	req, err := oghttp.NewRequest("GET", "/api/imports/"+id.String(), nil)
	suite.NoError(err)
	rr := httptest.NewRecorder()

	// Execute
	h.ServeHTTP(rr, req)

	// Assert
	suite.Equal(oghttp.StatusInternalServerError, rr.Code)
	suite.Equal("application/json", rr.Header().Get("Content-Type"))
	suite.Equal(`{"error":{"message":"Internal Server Error"}}`+"\n", rr.Body.String())

	// Assert log
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"ERROR"`)
	suite.Contains(lines[0], `failed to get import: failed to find import `+id.String()+`: boom!`)
}
