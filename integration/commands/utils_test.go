package commands

import (
	"testing"
)

func TestCheckAPIInstance(t *testing.T) {
	cases := []struct {
		name     string
		envLog   string
		instance int
		spec     []string
		expected bool
	}{
		{
			name: "check envs for community",
			envLog: `
                TERM=xterm
                RECORD_RETENTION=0
                PRIVATE_KEY=/run/secrets/api_private_key
                MAXMIND_LICENSE=
                SHELLHUB_ENTERPRISE=false
                SHELLHUB_BILLING=false
                STORE_CACHE=false
                PUBLIC_KEY=/run/secrets/api_public_key
                GEOIP=false
                SHELLHUB_CLOUD=false
                GOLANG_VERSION=1.17.6
                GOPATH=/go
                GOPROXY=http://localhost:3333
                MONGO_DB_NAME=test
            `,
			instance: Community,
			spec:     []string{},
			expected: true,
		},
		{
			name: "check the enterprise instance",
			envLog: `
                RECORD_RETENTION=0
                PRIVATE_KEY=/run/secrets/api_private_key
                PUBLIC_KEY=/run/secrets/api_public_key
                STORE_CACHE=false
                GEOIP=false
                MAXMIND_LICENSE=
                GOLANG_VERSION=1.17.6
                GOPATH=/go
                GOPROXY=
                MONGO_DB_NAME=test
                SHELLHUB_ENTERPRISE=true
                SHELLHUB_CLOUD=false
                SHELLHUB_STORE_CACHE=false
                HOME=/root
            `,
			instance: Enterprise,
			spec:     []string{},
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if r := CheckAPIInstance(tc.envLog, tc.instance, tc.spec); r != tc.expected {
				t.Errorf("Error, expected: %v, got: %v\n", tc.expected, r)
			}
		})
	}
}
