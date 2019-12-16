package flow_cleaner

import (
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "flow_cleaner"
const FlowLimitDefault = 3000

type Operator interface {
	Clean(limit uint64) error
}