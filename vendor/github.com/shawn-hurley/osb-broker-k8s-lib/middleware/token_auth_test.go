package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubernetes/client-go/kubernetes/typed/authentication/v1"
	authenticationapi "k8s.io/api/authentication/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type fakeTokenReview struct {
	TokenReview *authenticationapi.TokenReview
}

func (ftr fakeTokenReview) Create(tr *authenticationapi.TokenReview) (*authenticationapi.TokenReview, error) {
	if tr.Spec.Token != ftr.TokenReview.Spec.Token {
		return nil, fmt.Errorf("token was not the same")
	}
	return ftr.TokenReview, nil
}

func TestTokenReviewMiddleware(t *testing.T) {
	cases := []struct {
		name         string
		tokenReview  v1.TokenReviewInterface
		header       string
		responseCode int
		errorMessage string
	}{
		{
			name:         "no auth string",
			tokenReview:  fake.NewSimpleClientset().Authentication().TokenReviews(),
			header:       "",
			responseCode: http.StatusUnauthorized,
			errorMessage: "unable to find authentication token",
		},
		{
			name:         "only bearer in auth string",
			tokenReview:  fake.NewSimpleClientset().Authentication().TokenReviews(),
			header:       "bearer",
			responseCode: http.StatusUnauthorized,
			errorMessage: "invalid authentication",
		},
		{
			name:         "no bearer in auth string",
			tokenReview:  fake.NewSimpleClientset().Authentication().TokenReviews(),
			header:       "faker newsimpletoken",
			responseCode: http.StatusUnauthorized,
			errorMessage: "invalid authentication",
		},
		{
			name:         "unauthenticated user",
			tokenReview:  fake.NewSimpleClientset().Authentication().TokenReviews(),
			header:       "bearer newsimpletoken",
			responseCode: http.StatusUnauthorized,
			errorMessage: "user was not authenticated",
		},
		{
			name:         "token review failure",
			tokenReview:  fakeTokenReview{&authenticationapi.TokenReview{Spec: authenticationapi.TokenReviewSpec{Token: "newsimpletoken"}, Status: authenticationapi.TokenReviewStatus{Authenticated: true}}},
			header:       "bearer anothertoken",
			responseCode: http.StatusUnauthorized,
			errorMessage: "unable to authenticate token",
		},
		{
			name:         "authenticated user",
			tokenReview:  fakeTokenReview{&authenticationapi.TokenReview{Spec: authenticationapi.TokenReviewSpec{Token: "newsimpletoken"}, Status: authenticationapi.TokenReviewStatus{Authenticated: true}}},
			header:       "bearer newsimpletoken",
			responseCode: http.StatusOK,
			errorMessage: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			trm := TokenReviewMiddleware{
				TokenReview: tc.tokenReview,
			}
			req := httptest.NewRequest("GET", "http://example.com/foo", nil)
			req.Header.Add("Authorization", tc.header)
			w := httptest.NewRecorder()
			trm.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				return
			})).ServeHTTP(w, req)
			resp := w.Result()
			if resp.StatusCode == http.StatusOK && tc.responseCode == resp.StatusCode {
				return
			}
			if resp.StatusCode != tc.responseCode {
				t.Fatalf("invalid response code expected %v, got: %v", tc.responseCode, w.Code)
			}
			if resp.Header.Get("Content-Type") != "application/json" {
				t.Fatalf("invalid content type expected %v, got: %v", "application/json", w.Header().Get("Content-Type"))
			}
			defer resp.Body.Close()
			e := osbError{}
			err := json.NewDecoder(resp.Body).Decode(&e)
			if err != nil {
				t.Fatalf("invalid json data in response body: %v", err)
			}
			if e.Description != tc.errorMessage {
				t.Fatalf("invalid error description expected %v, got: %v", tc.errorMessage, e.Description)
			}

		})
	}
}
