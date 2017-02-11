package prompt

import (
	"interactive/structs"
)

func Selection(question string, options []*structs.Option) {
	r := structs.NewRenderer(question, options)
	r.Init()
}
