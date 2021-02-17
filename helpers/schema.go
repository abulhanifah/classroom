package helpers

import (
	"math"
	"strconv"
	"strings"

	"github.com/abulhanifah/classroom/constants"
	"github.com/jinzhu/gorm"
)

func GetById(ctx Context, object, key, id string, params map[string][]string, schema map[string]interface{}) map[string]interface{} {
	params[key] = []string{id}
	params["is_skip_count"] = []string{"true"}
	params["is_include_has_many"] = []string{"true"}
	data := GetPaginated(ctx, params, schema)
	res := map[string]interface{}{}
	if data["results"] != nil {
		temp := data["results"].([]map[string]interface{})
		if len(temp) > 0 {
			res = temp[0]
		}
	}
	if res[key] != nil {
		return res
	} else {
		return NotFoundMessage(object, key, id)
	}
}

func GetPaginated(ctx Context, params map[string][]string, schema map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	count := int64(0)
	pageContext := map[string]int64{}

	t := schema["table"].(map[string]string)
	table := t["name"]
	if t["as"] != "" {
		table = t["name"] + " as " + t["as"]
	}
	db := GetDB(ctx).Table(table)
	db = db.Unscoped()
	db = setJoin(db, params, schema)
	db = setWhere(db, params, schema)
	if params["is_skip_count"] == nil {
		db.Count(&count)
		res["count"] = count
	}

	db = setSelect(db, params, schema)
	db = setOrder(db, params, schema)
	if params["is_skip_count"] == nil {
		db, pageContext = SetPage(db, params, count)
		res["page_context"] = pageContext
	}
	rows, _ := db.Rows()
	if rows != nil {
		results := DotToInterface(GetResults(rows), schema)
		if params["is_include_has_many"] != nil {
			for i, d := range results {
				results[i] = GetHasMayData(ctx, d, schema)
			}
		}
		res["results"] = results
	}
	return res
}

// fields=field_a,field_b,field_c
func setSelect(db *gorm.DB, params map[string][]string, schema map[string]interface{}) *gorm.DB {
	table := schema["table"].(map[string]string)
	fields := schema["fields"].(map[string]map[string]string)
	selectedField := []string{}
	if params["fields"] != nil {
		for _, field := range strings.Split(params["fields"][0], ",") {
			if fields[field] != nil {
				selectedField = append(selectedField, fields[field]["name"]+" as "+fields[field]["as"])
			}
		}
	} else if table["as"] != "" {
		for _, f := range fields {
			if f["is_hide"] != "true" {
				selectedField = append(selectedField, f["name"]+" as "+f["as"])
			}
		}
	}
	if len(selectedField) > 0 {
		db = db.Select(selectedField)
	}
	return db
}

func setJoin(db *gorm.DB, params map[string][]string, schema map[string]interface{}) *gorm.DB {
	if schema["relations"] != nil {
		joinType := "left join"
		for _, r := range schema["relations"].([]map[string]string) {
			if r["type"] == "BelongsTo" {
				joinType = "inner join"
			}
			db = db.Joins(joinType + " " + r["name"] + " as " + r["as"] + " on " + r["on"])
		}
	}
	return db
}

// field_a=value_a&field_b.$gte=10&field_c.$ilike
func setWhere(db *gorm.DB, params map[string][]string, schema map[string]interface{}) *gorm.DB {
	if schema["where"] != nil {
		where := schema["where"].([]map[string]interface{})
		for _, w := range where {
			if w["raw"] != nil {
				db = db.Where(w["raw"].(string))
			}
		}
	}
	fields := schema["fields"].(map[string]map[string]string)
	for param, value := range params {
		if fields[param] != nil {
			db = db.Where(fields[param]["name"] + " = '" + value[0] + "'")
		} else {
			temp := strings.Split(param, ".")
			if len(temp) > 1 {
				field := strings.Join(temp[:len(temp)-1], ".")
				operator := temp[len(temp)-1]
				if fields[field] != nil {
					if constants.Operator[operator] != "" && operator != "$in" && operator != "$nin" {
						db = db.Where(fields[field]["name"] + " " + constants.Operator[operator] + " '" + value[0] + "'")
					} else {
						db = db.Where(fields[field]["name"]+" "+constants.Operator[operator]+" (?)", strings.Split(value[0], ","))
					}
				}
			}
		}
	}
	return db
}

// sorts=field_asc,-field_desc,field_asc:i
func setOrder(db *gorm.DB, params map[string][]string, schema map[string]interface{}) *gorm.DB {
	if params["sorts"] != nil {
		fields := schema["fields"].(map[string]map[string]string)
		for _, sort := range strings.Split(params["sorts"][0], ",") {
			direction := "asc"
			descending := strings.Split(sort, "-")
			if len(descending) > 1 {
				direction = "desc"
				sort = descending[1]
			}
			caseInsensitive := strings.Split(sort, ":")
			sort = caseInsensitive[0]
			if fields[sort] != nil {
				field := fields[sort]["name"]
				if fields[sort]["as"] != "" {
					field = fields[sort]["as"]
				}
				if len(caseInsensitive) > 1 && caseInsensitive[1] == "i" {
					field = "lower(" + field + ")"
				}
				db = db.Order(field + " " + direction)
			}
		}
	}
	return db
}

// page=1&per_page=10
func SetPage(db *gorm.DB, params map[string][]string, count int64) (*gorm.DB, map[string]int64) {
	page := int64(1)
	if params["page"] != nil {
		page, _ = strconv.ParseInt(params["page"][0], 10, 64)
	}
	perPage := int64(10)
	if params["per_page"] != nil {
		perPage, _ = strconv.ParseInt(params["per_page"][0], 10, 64)
	}
	totalPages := int64(math.Ceil(float64(count) / float64(perPage)))
	offset := int64((page - 1) * perPage)
	pageContext := map[string]int64{
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	}
	return db.Limit(perPage).Offset(offset), pageContext
}

func GetHasMayData(ctx Context, data map[string]interface{}, schema map[string]interface{}) map[string]interface{} {
	if schema["has_many_relations"] != nil {
		for f, r := range schema["has_many_relations"].(map[string]map[string]interface{}) {
			filter := map[string][]string{}
			filter["is_skip_count"] = []string{"true"}
			filter["is_include_has_many"] = []string{"true"}
			if r["primary_key"] != nil && data[r["primary_key"].(string)] != nil {
				filter[r["foreign_key"].(string)] = []string{Iconvert{Val: data[r["primary_key"].(string)]}.String()}
				temp := GetPaginated(ctx, filter, r["schema"].(map[string]interface{}))
				if temp["results"] != nil {
					data[f] = temp["results"]
				}
			}
		}
	}
	return data
}
