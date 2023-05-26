package google

import "github.com/agi-cn/llmplugin/llm"

type Option func(g *Google)

// WithSummarizer 总结内容
func WithSummarizer(summarizer llm.Summarizer) Option {

	return func(g *Google) {
		g.summarizer = summarizer
	}
}
