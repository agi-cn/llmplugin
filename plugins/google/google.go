// Package google is Google Search Plugin
// Get a Google Serach API key according to the Instruction.
// https://stackoverflow.com/questions/37083058/programmatically-searching-google-in-python-using-custom-search

package google

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

const (
	pluginName = "Google"

	pluginInputExample = "Who is Google boss?"

	pluginDesc = `Search something by query input.`
)

type Google struct {
	customSearchID string
	apiToken       string
}

func NewGoogle(customSearchID, apiToken string) *Google {

	return &Google{
		customSearchID: customSearchID,
		apiToken:       apiToken,
	}
}

func (g Google) Do(ctx context.Context, query string) (answer string, err error) {

	client, err := customsearch.NewService(ctx, option.WithAPIKey(g.apiToken))
	if err != nil {
		return "", errors.Wrap(err, "new google service failed")
	}

	results, err := client.Cse.List().Q(query).Cx(g.customSearchID).Do()
	if err != nil {
		return "", errors.Wrap(err, "google search failed")
	}

	items := results.Items
	// only top-3
	if len(items) > 3 {
		items = items[:3]
	}
	lines := make([]string, 0, len(items))
	for i, item := range items {
		number := i + 1

		line := fmt.Sprintf(
			"<%d> %s: %s (%s)",
			number,
			item.Title,
			item.Snippet,
			item.Link,
		)

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}

func (g Google) GetName() string {
	return pluginName
}

func (g Google) GetInputExample() string {
	return pluginInputExample
}

func (g Google) GetDesc() string {
	return pluginDesc
}
