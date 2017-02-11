package selection

func Ask(question string, options []*Option) ([]interface{}, bool) {
	r := NewRenderer(question, options)
	return r.Init()
}
