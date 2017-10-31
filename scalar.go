package bigiot

import (
	"github.com/shurcooL/graphql"
)

type (
	// Long represents a 64 bit signed integer. Used for expiration times on
	// activations
	Long int64

	// String represents our simple graphql string type
	String string

	// BigDecimal is used for handling currency values. This is not a precise value
	// so should be revisited with a more precise numeric type
	BigDecimal float64

	// Boolean represents a true or false value
	Boolean graphql.Boolean
)
