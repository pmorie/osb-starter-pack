package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"k8s.io/api/authentication/v1"
	authorizationv1 "k8s.io/api/authorization/v1"
	authv1 "k8s.io/client-go/kubernetes/typed/authorization/v1"
)

// Decision - should an action be allowed or not
type Decision string

const (
	// DecisionAllowed - action should be allowed
	DecisionAllowed Decision = "allowed"
	// DecisionDeny  - action should not be allowed
	DecisionDeny Decision = "deny"
	// DecisionNoOpinion - up to the caller to allow action or not
	DecisionNoOpinion = "no opinion"
)

// UserInfoAuthorizer - Authorizes k8s user info for a request.
type UserInfoAuthorizer interface {
	Authorize(v1.UserInfo, *http.Request) (Decision, error)
}

// SARUserInfoAuthorizer - Authorizes a k8s user info with a
// Subject Access Review
type SARUserInfoAuthorizer struct {
	SAR authv1.SubjectAccessReviewExpansion
}

// Authorize - Subject Access Review authorize the user.
func (s SARUserInfoAuthorizer) Authorize(u v1.UserInfo, req *http.Request) (Decision, error) {
	review := &authorizationv1.SubjectAccessReview{
		Spec: authorizationv1.SubjectAccessReviewSpec{
			User:   u.Username,
			UID:    u.UID,
			Groups: u.Groups,
			Extra:  convertToSARExtraValue(u.Extra),
		},
	}

	// For the OSB we will never have resource URLs so we should only check
	// the NonResourceAttributes of the request.

	review.Spec.NonResourceAttributes = &authorizationv1.NonResourceAttributes{
		Path: req.URL.Path,
		Verb: req.Method,
	}

	review, err := s.SAR.Create(review)
	if err != nil {
		glog.Errorf("Failed to create subject access review: %v", err)
		return DecisionDeny, err
	}
	switch {
	case review.Status.Denied && review.Status.Allowed:
		return DecisionDeny, fmt.Errorf("review has both denied and allowed the request. defaulting to closed")
	case review.Status.Denied:
		return DecisionDeny, nil
	case review.Status.Allowed:
		return DecisionAllowed, nil
	default:
		// If both allowed and denied are false, then the review "has no opinion"
		return DecisionNoOpinion, nil
	}
}

func convertToSARExtraValue(extra map[string]v1.ExtraValue) map[string]authorizationv1.ExtraValue {
	if extra == nil {
		return nil
	}
	ext := map[string]authorizationv1.ExtraValue{}
	for k, v := range extra {
		ext[k] = authorizationv1.ExtraValue(v)
	}
	return ext
}
