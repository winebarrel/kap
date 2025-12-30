package kap_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/kap"
)

func TestSecretHas(t *testing.T) {
	tests := []struct {
		secrets  []string
		input    string
		expected bool
	}{
		{
			secrets:  []string{"foo", "bar", "zoo"},
			input:    "foo",
			expected: true,
		},
		{
			secrets:  []string{"foo", "bar", "zoo"},
			input:    "bar",
			expected: true,
		},
		{
			secrets:  []string{"foo", "bar", "zoo"},
			input:    "baz",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v has '%s' -> %v", tt.secrets, tt.input, tt.expected), func(t *testing.T) {
			assert := assert.New(t)
			secret := kap.Secret(tt.secrets)
			assert.Equal(tt.expected, secret.Has(tt.input))
		})
	}
}

func TestSecretHasWithHash(t *testing.T) {
	tests := []struct {
		secrets  []string
		input    string
		expected bool
	}{
		{
			secrets: []string{
				"$2y$05$mY3jIrHMxws14rbG2EJRSeB7SzZbcEQdk2fEG9nTN5ILGoYR05U/.", // foo
				"$2y$05$Gdin/n.0ipS.XEODHu6Vxufv1sIpHC37MyU6Vf0zWfSYZGfUAiXRS", // bar
				"$2y$05$JHzeO43Y8sXvMSvhTQlO5ObnInLSnhoiDHK0IZo4oaXZyit.vrM3O", // zoo
			},
			input:    "foo",
			expected: true,
		},
		{
			secrets: []string{
				"$2y$05$MLxPAOAaoivyZFZGsC9XoOEHVuSFk0Nrs3w44jnR05iFYVqvc1Tga", // foo
				"$2y$05$JkvrjeCFHRuhHUyidnxcju9b5R6zaiPJorfadiHSm28VnviyVxoCq", // bar
				"$2y$05$z4xTs3F597RkFIAv23frdeJSmEEkM.PuLOLWYigoiGldBCKC4OFYC", // zoo
			},
			input:    "bar",
			expected: true,
		},
		{
			secrets: []string{
				"$2y$05$L1yKQ2OkDINLryTbdgQ7ue6r139HSrQIqIfPPf0w3.Zd5XoEW9Mee", // foo
				"$2y$05$wr8cofZ4ciH0uMeV44li/OZGbTz3z92CnLKFhWiDlomQu75MASNEW", // bar
				"$2y$05$fNXLheZwYvRSjlUDWvpK8e.VztZQBLL5a7hQAGs2AqaIf1UZfegEa", // zoo
			},
			input:    "baz",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v has '%s' -> %v", tt.secrets, tt.input, tt.expected), func(t *testing.T) {
			assert := assert.New(t)
			secret := kap.Secret(tt.secrets)
			assert.Equal(tt.expected, secret.Has(tt.input))
		})
	}
}
