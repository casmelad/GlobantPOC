package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("bad request")
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
	r.Methods(http.MethodGet).Path(UsersBaseUri).Handler(httptransport.NewServer(
		e.GetAllUsersEndpoint,
		decodeGetAllUsersRequest,
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

func decodeGetAllUsersRequest(_ context.Context, r *http.Request) (request interface{}, err error) {

	return getAllUsersRequest{}, nil
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

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client.
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
	case ErrNotFound:
		return http.StatusNotFound
	case ErrUserAlreadyExists:
		return http.StatusConflict
	case ErrInvalidInput:
		return http.StatusUnprocessableEntity
	case ErrBadRouting:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
