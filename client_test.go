package salesforce_client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestServer struct {
	Mux    *http.ServeMux
	Server *httptest.Server
	Client *Client
}

func (t *TestServer) Setup() func() {
	t.Mux = http.NewServeMux()
	t.Server = httptest.NewServer(t.Mux)

	t.Mux.HandleFunc("/services/data/v20.0/sobjects/Task/00T0x000005hOcYEAU", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		f, err := os.Open("./test_fixtures/task.json")
		if err != nil {
			fmt.Println("Error: ", err)
		}

		c, _ := ioutil.ReadAll(f)
		_, _ = w.Write(c)
	})

	t.Mux.HandleFunc("/services/data/v20.0/sobjects/Task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`[{"message": "Bad request", "errorCode": "BAD_REQUEST"}]`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		f, err := os.Open("./test_fixtures/task.json")
		if err != nil {
			fmt.Println("Error: ", err)
		}

		c, _ := ioutil.ReadAll(f)
		_, _ = w.Write(c)
	})

	t.Client, _ = NewClient("stg", "v20.0")
	t.Client.baseUrl = t.Server.URL

	return func() {
		t.Server.Close()
	}
}

func TestClientMethods_Find(t *testing.T) {
	server := &TestServer{}
	serverCloser := server.Setup()
	defer serverCloser()

	obj := &struct {
		ID     string `json:"Id"`
		Status string
	}{}
	_ = server.Client.Find("Task", "00T0x000005hOcYEAU", obj)

	test := assert.New(t)
	test.Equal("00T0x000005hOcYEAU", obj.ID)
	test.Equal("Completed", obj.Status)
}

func TestClientMethods_Create(t *testing.T) {
	server := &TestServer{}
	serverCloser := server.Setup()
	defer serverCloser()

	obj := &struct {
		Id     string
		Status string
	}{}

	p := &Params{}
	p.AddField("Name", "Tests")
	_ = server.Client.Create("Task", p, obj)

	test := assert.New(t)
	test.Equal("00T0x000005hOcYEAU", obj.Id)
	test.Equal("Completed", obj.Status)
}
