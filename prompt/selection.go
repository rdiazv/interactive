package prompt

import (
	"interactive/structs"
)

func Selection(options []*structs.Option) {
	r := structs.NewRenderer()
	r.Options = options
	r.Init()
}
