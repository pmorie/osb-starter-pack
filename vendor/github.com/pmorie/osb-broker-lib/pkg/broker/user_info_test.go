package broker

import (
	"reflect"
	"testing"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

func TestParseIdentity(t *testing.T) {
	cases := []struct {
		name              string
		o                 osb.OriginatingIdentity
		shouldErr         bool
		expectingIdentity Identity
	}{
		{
			name: "kubernetes invalid value json",
			o: osb.OriginatingIdentity{
				Platform: osb.PlatformKubernetes,
				Value:    `{ser_id": "123"}`,
			},
			shouldErr: true,
		},
		{
			name: "kubernetes valid identity",
			o: osb.OriginatingIdentity{
				Platform: osb.PlatformKubernetes,
				Value:    `{"username": "foo", "groups":[], "extra": {}}`,
			},
			expectingIdentity: Identity{
				Platform: "kubernetes",
				Kubernetes: &KubernetesUserInfo{
					Username: "foo",
					Groups:   []string{},
					Extra:    map[string][]string{},
				},
			},
		},
		{
			name: "cloud foundry invalid value json",
			o: osb.OriginatingIdentity{
				Platform: "cloudfoundry",
				Value:    `{ser_id": "123"}`,
			},
			shouldErr: true,
		},
		{
			name: "cloud foundry no user_id",
			o: osb.OriginatingIdentity{
				Platform: "cloudfoundry",
				Value:    `{"username": "123"}`,
			},
			shouldErr: true,
		},
		{
			name: "cloud foundry invalid user_id",
			o: osb.OriginatingIdentity{
				Platform: "cloudfoundry",
				Value:    `{"username": 123}`,
			},
			shouldErr: true,
		},
		{
			name: "cloud foundry valid identity",
			o: osb.OriginatingIdentity{
				Platform: "cloudfoundry",
				Value:    `{"user_id": "123"}`,
			},
			expectingIdentity: Identity{
				Platform: "cloudfoundry",
				CloudFoundry: &CloudFoundryUserInfo{
					UserID: "123",
					Extras: map[string]interface{}{},
				},
			},
		},
		{
			name: "unknown identity",
			o: osb.OriginatingIdentity{
				Platform: "Unknown",
				Value:    `{"user_id": "123"}`,
			},
			expectingIdentity: Identity{
				Platform: "Unknown",
				Unknown: map[string]interface{}{
					"user_id": "123",
				},
			},
		},
	}
	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ParseIdentity(tc.o)
			if err != nil {
				if tc.shouldErr {
					return
				}
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if !reflect.DeepEqual(tc.expectingIdentity, actual) {
				t.Errorf("Unexpected response\n\nExpecting identity:%#+v\nGot: %#+v", tc.expectingIdentity, actual)
			}
		})
	}
}
