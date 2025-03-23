package build_release

import (
	"encoding/base64"
	"regexp"
	"slices"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/rhysd/actionlint"
)


var encodedWorkflowFile string = "bmFtZTogT1NQUyBCYXNlbGluZSBTY2FuCgpvbjogW3dvcmtmbG93X2Rpc3Bh\ndGNoXQoKam9iczoKICBzY2FuOgogICAgcnVucy1vbjogdWJ1bnR1LWxhdGVz\ndAoKICAgIHN0ZXBzOgogICAgICAtIG5hbWU6IENoZWNrb3V0IHJlcG9zaXRv\ncnkKICAgICAgICB1c2VzOiBhY3Rpb25zL2NoZWNrb3V0QHY0CgogICAgICAt\nIG5hbWU6IFB1bGwgdGhlIHB2dHItZ2l0aHViLXJlcG8gaW1hZ2UKICAgICAg\nICBydW46IGRvY2tlciBwdWxsIGVkZGlla25pZ2h0L3B2dHItZ2l0aHViLXJl\ncG86bGF0ZXN0CgogICAgICAtIG5hbWU6IEFkZCBHaXRIdWIgU2VjcmV0IHRv\nIGNvbmZpZyBmaWxlIHNvIGl0IGlzIHByb3RlY3RlZCBpbiBvdXRwdXRzCiAg\nICAgICAgcnVuOiB8CiAgICAgICAgICBzZWQgLWkgJ3Mve3sgVE9LRU4gfX0v\nJHt7IHNlY3JldHMuVE9LRU4gfX0vZycgJHt7IGdpdGh1Yi53b3Jrc3BhY2Ug\nfX0vLmdpdGh1Yi9wdnRyLWNvbmZpZy55bWwKCiAgICAgIC0gbmFtZTogU2Nh\nbiBhbGwgcmVwb3Mgc3BlY2lmaWVkIGluIC5naXRodWIvcHZ0ci1jb25maWcu\neW1sCiAgICAgICAgcnVuOiB8CiAgICAgICAgICBkb2NrZXIgcnVuIC0tcm0g\nXAogICAgICAgICAgICAtdiAke3sgZ2l0aHViLndvcmtzcGFjZSB9fS8uZ2l0\naHViL3B2dHItY29uZmlnLnltbDovLnByaXZhdGVlci9jb25maWcueW1sIFwK\nICAgICAgICAgICAgLXYgJHt7IGdpdGh1Yi53b3Jrc3BhY2UgfX0vZG9ja2Vy\nX291dHB1dDovZXZhbHVhdGlvbl9yZXN1bHRzIFwKICAgICAgICAgICAgZWRk\naWVrbmlnaHQvcHZ0ci1naXRodWItcmVwbzpsYXRlc3QK\n"

var goodWorkflowFile = 
`name: OSPS Baseline Scan

on: [workflow_dispatch]

jobs:
  scan:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Pull the pvtr-github-repo image
        run: docker pull eddieknight/pvtr-github-repo:latest

      - name: Add GitHub Secret to config file so it is protected in outputs
        run: |
          sed -i 's/{{ TOKEN }}/${{ secrets.TOKEN }}/g' ${{ github.workspace }}/.github/pvtr-config.yml

      - name: Scan all repos specified in .github/pvtr-config.yml
        run: |
          docker run --rm \
            -v ${{ github.workspace }}/.github/pvtr-config.yml:/.privateer/config.yml \
            -v ${{ github.workspace }}/docker_output:/evaluation_results \
            eddieknight/pvtr-github-repo:latest`


var badWorkflowFile =
`name: OSPS Baseline Scan

on: [workflow_dispatch]

jobs:
  scan:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Pull the pvtr-github-repo image
        run: docker pull eddieknight/pvtr-github-repo:latest

      - name: Add GitHub Secret to config file so it is protected in outputs
        run: |
          sed -i 's/{{ TOKEN }}/${{ secrets.TOKEN }}/g' ${{ github.workspace }}/.github/pvtr-config.yml

      - name: Scan all repos specified in .github/pvtr-config.yml
        run: |
          docker run --rm \
            -v ${{ github.event.issue.title }}/.github/pvtr-config.yml:/.privateer/config.yml \
            -v ${{ github.workspace }}/docker_output:/evaluation_results \
            eddieknight/pvtr-github-repo:latest`


type testingData struct {
	expectedResult bool
	workflowFile string
	assertionMessage string
}



var testScript = 
` echo ${{github.event.issue.title }}
  if ${{ github.event.commits.arbitrary.data.message}} -ne 0
  then
	echo "Checkout report image" ${{ githubnodotevent.commits.arbitrary.data.message}}
	run: docker pull the pvt-r-github-repo image
  fi`



func TestCicdSanitizedInputParameters (t * testing.T) {

	testData := []testingData {
		{
			expectedResult: false,
			workflowFile: badWorkflowFile,
			assertionMessage: "Untrusted input not detected",
		},
		{
			expectedResult: true,
			workflowFile: goodWorkflowFile,
			assertionMessage: "Untrusted input detected where it should not have been",
		},
	}

	for _, data := range testData {

		workflow, _ := actionlint.Parse([]byte(data.workflowFile))

		result, _ := checkWorkflowFileForUntrustedInputs(workflow)

		assert.Equal(t, result, data.expectedResult, data.assertionMessage)
	}
}


func TestVariableExtraction(t *testing.T) {

	varNames := pullVariablesFromScript(testScript)

	assert.Equal(t, slices.Contains(varNames, "github.event.issue.title"), true, "Variable extraction failed")
	assert.Equal(t, slices.Contains(varNames, "github.event.commits.arbitrary.data.message"), true, "Variable extraction failed")

}



func TestRegex ( t * testing.T ) {

	expression, err := regexp.Compile(regex)
	if err != nil {
		t.Errorf("Error compiling regex: %v", err)
		return
	}

	assert.Equal(t, expression.Match([]byte("github.event.issue.title")), true, "regex match failed" )
	assert.Equal(t, expression.Match([]byte("github.event.commits.arbitrary.data.message")), true, "regex match failed" )
}


//TODO remove this test
func TestParse( t *testing.T) {

	decoded, _ := base64.StdEncoding.DecodeString(encodedWorkflowFile)
	// fmt.Printf("Decoded String: %s", decoded)

	workflow, err := actionlint.Parse(decoded)

	if err != nil {
		t.Errorf("Error parsing workflow: %v", err)
	}

	assert.Equal(t, workflow.Name.Value, "OSPS Baseline Scan", "workflow parsing failed")

}
