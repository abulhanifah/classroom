package constants

var Operator = map[string]string{
	"$eq":      "=",
	"$ne":      "<>",
	"$gt":      ">",
	"$gte":     ">=",
	"$lte":     "<=",
	"$lt":      "<",
	"$like":    "like",
	"$ilike":   "ilike",
	"$nlike":   "not like",
	"$nilike":  "not ilike",
	"$in":      "in",
	"$nin":     "not in",
	"$regexp":  "regexp",
	"$nregexp": "not regexp",
}
