package xegrpcgw

import (
	"github.com/gotomicro/ego/server/egin"
)

// Component ...
type Component struct {
	*egin.Component // gin框架
}

// newComponent ...
func newComponent(eg *egin.Component) *Component {
	return &Component{
		Component: eg,
	}
}
