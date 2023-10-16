package internalhttp_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	mockedapp "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/app/mocked"
	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	internalhttproutes "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/server/http/routes"
)

type client struct {
	httpClient *http.Client
	baseURL    string
}

func newClient(baseURL string) *client {
	return &client{
		httpClient: http.DefaultClient,
		baseURL:    baseURL,
	}
}

func (c *client) send(method string, path string, payload any) (*http.Response, error) {
	fullURL, _ := url.JoinPath(c.baseURL, path)

	var body io.Reader
	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(jsonPayload)
	}

	req, err := http.NewRequest(method, fullURL, body) //nolint:noctx
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func getContent(body io.Reader) []byte {
	content, _ := io.ReadAll(body)
	return content
}

func toJSON[T any](content []byte) (T, error) {
	var data T

	if err := json.Unmarshal(content, &data); err != nil {
		return data, err
	}

	return data, nil
}

type eventTestData struct {
	req          map[string]any
	incorrectReq map[string]any
	expectedResp map[string]any
	before       *models.Event
	after        *models.Event
}

func newEventTestData() *eventTestData {
	dt := time.Now().Round(0)
	duration, _ := time.ParseDuration("1m")
	userID := uuid.New()

	return &eventTestData{
		req: map[string]any{
			"title":       "test",
			"description": "description",
			"datetime":    dt.Format(time.RFC3339),
			"duration":    duration.String(),
			"user_id":     userID.String(),
		},
		incorrectReq: map[string]any{
			"title":       "t",
			"description": "desc",
			"datetime":    "incorrect dt",
			"user_id":     "incorrect uuid",
		},
		expectedResp: map[string]any{
			"id":          float64(1),
			"title":       "test",
			"description": "description",
			"datetime":    dt.Format(time.RFC3339),
			"duration":    duration.String(),
			"user_id":     userID.String(),
		},
		before: &models.Event{
			Title: "test",
			Description: sql.NullString{
				String: "description",
				Valid:  true,
			},
			Datetime: dt,
			Duration: duration,
			UserID:   userID,
		},
		after: &models.Event{
			ID:    1,
			Title: "test",
			Description: sql.NullString{
				String: "description",
				Valid:  true,
			},
			Datetime: dt,
			Duration: duration,
			UserID:   userID,
		},
	}
}

type EventsHandlesSuite struct {
	suite.Suite
	mockedApp *mockedapp.MockApp
	server    *httptest.Server
	client    *client

	eventData *eventTestData
}

func (s *EventsHandlesSuite) SetupSuite() {
	s.mockedApp = mockedapp.NewMockApp(s.T())
	router := internalhttproutes.SetupRoutes(s.mockedApp)
	s.server = httptest.NewServer(router)
	s.client = newClient(s.server.URL)

	s.eventData = newEventTestData()
}

func (s *EventsHandlesSuite) TearDownSuite() {
	s.server.Close()
}

func (s *EventsHandlesSuite) TestCreateEventHandler() {
	t := s.T()

	t.Run("201", func(t *testing.T) {
		s.mockedApp.EXPECT().CreateEvent(mock.Anything, s.eventData.before).Return(s.eventData.after, nil).Once()

		resp, err := s.client.send("POST", "/events", s.eventData.req)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusCreated, resp.StatusCode, string(content))

		respBody, err := toJSON[map[string]any](content)
		require.NoError(t, err)
		require.Equal(t, s.eventData.expectedResp, respBody)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("400", func(t *testing.T) {
		resp, err := s.client.send("POST", "/events", s.eventData.incorrectReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusBadRequest, resp.StatusCode, string(content))

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("500", func(t *testing.T) {
		s.mockedApp.EXPECT().CreateEvent(mock.Anything, s.eventData.before).Return(nil, fmt.Errorf("test")).Once()

		resp, err := s.client.send("POST", "/events", s.eventData.req)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusInternalServerError, resp.StatusCode, string(content))

		s.mockedApp.AssertExpectations(t)
	})
}

func (s *EventsHandlesSuite) TestGetEventsHandler() {
	t := s.T()

	t.Run("200", func(t *testing.T) {
		s.mockedApp.EXPECT().GetEvents(mock.Anything).Return([]models.Event{*s.eventData.after}, nil).Once()

		resp, err := s.client.send("GET", "/events", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusOK, resp.StatusCode, string(content))

		respBody, err := toJSON[[]map[string]any](content)
		require.NoError(t, err)
		require.Equal(t, []map[string]any{s.eventData.expectedResp}, respBody)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("500", func(t *testing.T) {
		s.mockedApp.EXPECT().GetEvents(mock.Anything).Return(nil, fmt.Errorf("test")).Once()

		resp, err := s.client.send("GET", "/events", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusInternalServerError, resp.StatusCode, string(content))

		s.mockedApp.AssertExpectations(t)
	})
}

func (s *EventsHandlesSuite) TestGetEventHandler() {
	t := s.T()

	t.Run("200", func(t *testing.T) {
		s.mockedApp.EXPECT().GetEvent(mock.Anything, int64(1)).Return(s.eventData.after, nil).Once()

		resp, err := s.client.send("GET", "/events/1", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusOK, resp.StatusCode, string(content))

		respBody, err := toJSON[map[string]any](content)
		require.NoError(t, err)
		require.Equal(t, s.eventData.expectedResp, respBody)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("404", func(t *testing.T) {
		s.mockedApp.EXPECT().GetEvent(mock.Anything, int64(1)).Return(nil, app.ErrDocumentNotFound).Once()

		resp, err := s.client.send("GET", "/events/1", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusNotFound, resp.StatusCode, string(content))

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("500", func(t *testing.T) {
		s.mockedApp.EXPECT().GetEvent(mock.Anything, int64(1)).Return(nil, fmt.Errorf("test")).Once()

		resp, err := s.client.send("GET", "/events/1", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusInternalServerError, resp.StatusCode, string(content))

		s.mockedApp.AssertExpectations(t)
	})
}

func (s *EventsHandlesSuite) TestUpdateEventHandler() {
	t := s.T()

	t.Run("200", func(t *testing.T) {
		s.mockedApp.EXPECT().UpdateEvent(mock.Anything, int64(1), s.eventData.before).Return(s.eventData.after, nil).Once()

		resp, err := s.client.send("PUT", "/events/1", s.eventData.req)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusOK, resp.StatusCode, string(content))

		respBody, err := toJSON[map[string]any](content)
		require.NoError(t, err)
		require.Equal(t, s.eventData.expectedResp, respBody)

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("400", func(t *testing.T) {
		resp, err := s.client.send("PUT", "/events/1", s.eventData.incorrectReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusBadRequest, resp.StatusCode, string(content))

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("404", func(t *testing.T) {
		s.mockedApp.EXPECT().
			UpdateEvent(mock.Anything, int64(1), s.eventData.before).
			Return(nil, app.ErrDocumentNotFound).
			Once()

		resp, err := s.client.send("PUT", "/events/1", s.eventData.req)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusNotFound, resp.StatusCode, string(content))

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("500", func(t *testing.T) {
		s.mockedApp.EXPECT().UpdateEvent(mock.Anything, int64(1), s.eventData.before).Return(nil, fmt.Errorf("test")).Once()

		resp, err := s.client.send("PUT", "/events/1", s.eventData.req)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusInternalServerError, resp.StatusCode, string(content))

		s.mockedApp.AssertExpectations(t)
	})
}

func (s *EventsHandlesSuite) TestDeleteEventHandler() {
	t := s.T()

	t.Run("200", func(t *testing.T) {
		s.mockedApp.EXPECT().DeleteEvent(mock.Anything, int64(1)).Return(nil).Once()

		resp, err := s.client.send("DELETE", "/events/1", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusNoContent, resp.StatusCode, string(content))

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("404", func(t *testing.T) {
		s.mockedApp.EXPECT().DeleteEvent(mock.Anything, int64(1)).Return(app.ErrDocumentNotFound).Once()

		resp, err := s.client.send("DELETE", "/events/1", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusNotFound, resp.StatusCode, string(content))

		s.mockedApp.AssertExpectations(t)
	})

	t.Run("500", func(t *testing.T) {
		s.mockedApp.EXPECT().DeleteEvent(mock.Anything, int64(1)).Return(fmt.Errorf("test")).Once()

		resp, err := s.client.send("DELETE", "/events/1", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		content := getContent(resp.Body)
		require.Equalf(t, http.StatusInternalServerError, resp.StatusCode, string(content))

		s.mockedApp.AssertExpectations(t)
	})
}
