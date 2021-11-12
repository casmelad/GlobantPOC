package users

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a profilesvc server.
func MakeHTTPHandler(s UserProxy, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodPost).Path(PostUser).Handler(httptransport.NewServer(
		e.PostUserEndpoint,
		decodePostProfileRequest,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodGet).Path(GetUser).Handler(httptransport.NewServer(
		e.GetUserEndpoint,
		decodeGetUserRequest,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodPut).Path(PutUser).Handler(httptransport.NewServer(
		e.PutUserEndpoint,
		decodePutProfileRequest,
		encodeResponse,
		options...,
	))
	r.Methods(http.MethodDelete).Path(DeleteUser).Handler(httptransport.NewServer(
		e.DeleteUserEndpoint,
		decodeDeleteProfileRequest,
		encodeResponse,
		options...,
	))

	return r
}

//decoders

func decodeGetUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getUserRequest{Email: email}, nil
}

func decodePostProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postUserRequest
	if e := json.NewDecoder(r.Body).Decode(&req.User); e != nil {
		return nil, e
	}
	return req, nil
}

func decodePutProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		return nil, ErrBadRouting
	}
	var usr User
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		return nil, err
	}
	usr.Email = email
	return putUserRequest{
		User: usr,
	}, nil
}

func decodeDeleteProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := strconv.Atoi(vars["id"])
	if ok != nil {
		return nil, ErrBadRouting
	}
	return deleteUserRequest{UserID: id}, nil
}

func decodePostProfileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postUserResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetProfileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getUserResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePutProfileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response putUserResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

//encoders

func encodePostProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/users/")
	req.URL.Path = "/users/"
	return encodeRequest(ctx, req, request)
}

func encodeGetProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/users/{id}")
	r := request.(getUserRequest)
	profileID := url.QueryEscape(r.Email)
	req.URL.Path = "/users/" + profileID
	return encodeRequest(ctx, req, request)
}

func encodePutProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PUT").Path("/users/{id}")
	r := request.(putUserRequest)
	profileID := url.QueryEscape(r.User.Email)
	req.URL.Path = "/users/" + profileID
	return encodeRequest(ctx, req, request)
}

func encodeDeleteProfileRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/users/{id}")
	r := request.(deleteUserRequest)
	profileID := url.QueryEscape(strconv.Itoa(r.UserID))
	req.URL.Path = "/users/" + profileID
	return encodeRequest(ctx, req, request)
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// profilesvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	// case ErrNotFound:
	// 	return http.StatusNotFound
	// case ErrAlreadyExists, ErrInconsistentIDs:
	// 	return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
