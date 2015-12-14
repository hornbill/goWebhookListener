package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var callTests = []struct {
	status int
	key    string
	body   string
}{
	//-- Invalid key
	{500, "123", `Doesnt matter`},
	//-- Valid Response
	{200, "123456", `{"onEntityEvent":{"eventSource":"urn:webhook:entity:Contact:update","callingSessionId":"U2015111600017554","eventTime":"2015-11-16 14:50:47Z","actionBy":"admin","actionByName":"System Administrator","entity":"Contact","record":{"h_pk_id":"103","h_jobtitle":"test","h_datelastmodified":"2015-11-16T14:50:47"}}}`},
	//-- EmptyBody
	{500, "123456", ``},
	//-- Invalid JSON
	{500, "123456", `{`},
	//-- No Key
	{500, "nil", ``},
	//-- Empty Key
	{500, "", ``},
}

func TestHTTPAuth(t *testing.T) {

	for _, tt := range callTests {

		userJSON := strings.NewReader(tt.body)
		requestURL := "http://something.com/api/"
		if tt.key != "nil" {
			requestURL = requestURL + "?key=" + tt.key
		}

		request, err := http.NewRequest("POST", requestURL, userJSON)
		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()
		webhookCatcher(res, request)

		if res.Code != tt.status {
			t.Errorf("Expected %d got: %d\n", tt.status, res.Code)
		}

	}
}
