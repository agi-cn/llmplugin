// Package google is Google Search Plugin
// Get a Google Serach API key according to the Instruction.
// https://stackoverflow.com/questions/37083058/programmatically-searching-google-in-python-using-custom-search

package google

import (
	"context"
	"fmt"
	"strings"

	"github.com/agi-cn/llmplugin/llm"
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

	summarizer llm.Summarizer
}

func NewGoogle(customSearchID, apiToken string, options ...Option) *Google {
	g := &Google{
		customSearchID: customSearchID,
		apiToken:       apiToken,
	}

	for _, o := range options {
		o(g)
	}

	return g
}

func (g Google) Do(ctx context.Context, query string) (answer string, err error) {

	results, err := g.doSearch(ctx, query)
	if err != nil {
		return "", err
	}

	return g.makeResult(ctx, query, results)
}

func (g Google) doSearch(ctx context.Context, query string) (*customsearch.Search, error) {
	client, err := customsearch.NewService(ctx, option.WithAPIKey(g.apiToken))
	if err != nil {
		return nil, errors.Wrap(err, "new google service failed")
	}

	results, err := client.Cse.List().Q(query).Cx(g.customSearchID).Do()
	if err != nil {
		return nil, errors.Wrap(err, "google search failed")
	}

	return results, nil
}

func (g Google) makeResult(ctx context.Context, query string, results *customsearch.Search) (string, error) {

	items := results.Items
	if len(items) == 0 {
		return "Google don't known", nil
	}

	if g.summarizer == nil {
		return g.makeRawResult(ctx, items)
	} else {
		return g.makeResultBySummary(ctx, query, items)
	}
}

func (g Google) makeRawResult(ctx context.Context, items []*customsearch.Result) (string, error) {
	// Only return top1 result without llm summary
	item := items[0]

	content := fmt.Sprintf("%s\n%s\n%s",
		item.Title,
		item.Snippet,
		item.Link,
	)
	return content, nil
}

func (g Google) makeResultBySummary(ctx context.Context, query string, items []*customsearch.Result) (string, error) {

	prompt := `User query: %s.

Here is the google search result, you task is the following,

1. If there is a suspicion of advertising, then simply ignore the corresponding search results.
2. Select the top 3 most relevant search results to the user's query.
3. Summarize all the search results into one paragraph using Chinese.
4. In the corresponding results, give the relevant citation link.
5. Your job is to just summarize the search results, don't explain why you did it. If you ignore certain results, don't summarize those ignored results.


Summarize the following text delimited by triple Triple dashes. Each line below is a google search result and the corresponding link.

---
%s
---
`

	if len(items) > 10 {
		items = items[:10]
	}
	lines := make([]string, 0, len(items))
	for i, item := range items {
		number := i + 1

		line := fmt.Sprintf(
			"%d. %s: %s (%s)",
			number,
			item.Title,
			item.Snippet,
			item.Link,
		)

		lines = append(lines, line)
	}

	allResult := strings.Join(lines, "\n")
	content := fmt.Sprintf(prompt, query, allResult)

	return g.summarizer.Summary(ctx, content)
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
