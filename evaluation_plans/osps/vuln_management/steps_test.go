package vuln_management

import (
	"testing"

	"github.com/ossf/si-tooling/v2/si"
	"github.com/revanite-io/pvtr-github-repo/data"
	"github.com/revanite-io/sci/pkg/layer4"
	"github.com/stretchr/testify/assert"
)

type testingData struct {
	expectedResult layer4.Result
	expectedMessage string
	payloadData interface{}
	assertionMessage string
}


func TestSastToolDefined(t *testing.T) {
	
	testData := []testingData{
		{
			expectedResult: layer4.Passed,
			expectedMessage: "Static Application Security Testing documented in Security Insights",
			payloadData:    data.Payload{
				RestData: &data.RestData {
					Insights: si.SecurityInsights{
						Repository: si.Repository{
							Security: si.SecurityInfo{
								Tools: []si.Tool{
									{
										Type: "SAST",
									},
								},
							},
						},
					},
				},
			},
			assertionMessage: "Static Application Security Testing documented in Security Insights",
		},
		{
			expectedResult: layer4.Failed,
			expectedMessage: "No Static Application Security Testing documented in Security Insights",
			payloadData:    data.Payload{
				RestData: &data.RestData {
					Insights: si.SecurityInsights{
						Repository: si.Repository{
							Security: si.SecurityInfo{
								Tools: []si.Tool{
									{
										Type: "NotSast",
									},
								},
							},
						},
					},
				},
			},
			assertionMessage: "No Static Application Security Testing documented in Security Insights",
		},
	}
	
	for _, test := range testData {
		result, message := sastToolDefined(test.payloadData, nil)

		assert.Equal(t, test.expectedResult, result, test.assertionMessage)
		assert.Equal(t, test.expectedMessage, message, test.assertionMessage)
	}
	
}