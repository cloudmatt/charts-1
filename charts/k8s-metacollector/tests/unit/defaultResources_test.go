package unit

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/strings/slices"
)

const chartPath = "../../"

// Using the default values we want to test that all the expected resources are rendered.
func TestRenderedResourcesWithDefaultValues(t *testing.T) {
	t.Parallel()

	helmChartPath, err := filepath.Abs(chartPath)
	require.NoError(t, err)

	releaseName := "rendered-resources"

	// Template files that we expect to be rendered.
	templateFiles := []string{
		"clusterrole.yaml",
		"clusterrolebinding.yaml",
		"deployment.yaml",
		"service.yaml",
		"serviceaccount.yaml",
	}

	require.NoError(t, err)

	options := &helm.Options{}

	// Template the chart using the default values.yaml file.
	output, err := helm.RenderTemplateE(t, options, helmChartPath, releaseName, nil)
	require.NoError(t, err)

	// Extract all rendered files from the output.
	pattern := `# Source: k8s-metacollector/templates/([^\n]+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(output, -1)

	var renderedTemplates []string
	for _, match := range matches {
		// Filter out test templates.
		if !strings.Contains(match[1], "test-") {
			renderedTemplates = append(renderedTemplates, match[1])
		}
	}

	// Assert that the rendered resources are equal tho the expected ones.
	require.Equal(t, len(renderedTemplates), len(templateFiles), "should be equal")

	for _, rendered := range renderedTemplates {
		require.True(t, slices.Contains(templateFiles, rendered), "template files should contain all the rendered files")
	}
}
