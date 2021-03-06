package create_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"tweeter/db"
	"tweeter/db/models/user"
	"tweeter/handlers/endpoints/users"
	usersCreate "tweeter/handlers/endpoints/users/create"
	"tweeter/handlers/responses"
	. "tweeter/handlers/testutil"
	. "tweeter/testutil"
)

func TestCreateEndpoint(t *testing.T) {
	db.InitForTests()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Users#create Endpoint Suite")
}

var _ = Describe("Users#create Endpoint", func() {
	var (
		server   *httptest.Server
		request  RequestArgs
		response *http.Response
	)

	var sendRequest = func(request RequestArgs) *http.Response {
		return MustSendRequest(server, request)
	}

	AfterEach(func() {
		server.Close()
		db.RollbackTransactionForTests()
	})

	BeforeEach(func() {
		// Transaction should be before all other actions
		db.BeginTransactionForTests()
	})

	JustBeforeEach(func() {
		r := mux.NewRouter()
		usersCreate.Endpoint.Attach(r)
		server = httptest.NewServer(r)
		response = sendRequest(request)
	})

	var successfulRequest = func() RequestArgs {
		return RequestArgs{
			Method:   http.MethodPost,
			Endpoint: "/api/users",
			JSONBody: map[string]interface{}{
				"email":    "darren.tsung@gmail.com",
				"password": "password",
			},
		}
	}

	Context("with a valid email and password", func() {
		BeforeEach(func() {
			request = successfulRequest()
		})

		It("has a success response with non-zero ID", func() {
			idData := struct{ ID user.ID }{}
			MustReadSuccessData(response, idData)
			Expect(idData.ID).NotTo(Equal(0))
		})

		It("errors for requests with the same email", func() {
			secondResponse := sendRequest(successfulRequest())
			errors := MustReadErrors(secondResponse)
			Expect(errors).To(Equal([]responses.Error{users.ErrEmailAlreadyExists("darren.tsung@gmail.com")}))
		})
	})

	Context("with too short of a password", func() {
		BeforeEach(func() {
			request = successfulRequest()
			request.JSONBody["password"] = "12345"
		})

		It("errors with users.ErrPasswordTooShort", func() {
			errors := MustReadErrors(response)
			Expect(errors).To(Equal([]responses.Error{users.ErrPasswordTooShort}))
		})
	})

	Context("with malformed json", func() {
		BeforeEach(func() {
			request = successfulRequest()
			request.JSONBody = nil
			request.RawBody = StrPtr("not valid json")
		})

		It("errors with ErrInvalidBody", func() {
			errors := MustReadErrors(response)
			Expect(errors).To(Equal([]responses.Error{users.ErrInvalidBody}))
		})
	})

	Context("with wrong method", func() {
		BeforeEach(func() {
			request = successfulRequest()
			request.Method = http.MethodGet
		})

		It("responds with StatusMethodNotAllowed", func() {
			Expect(response.StatusCode).To(Equal(http.StatusMethodNotAllowed))
		})
	})
})
