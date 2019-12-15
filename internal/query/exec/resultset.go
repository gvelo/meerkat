package exec

type ResultSet struct {
	rowAffected int
	rowScanned  int
	colsName    []string
	cols        []interface{}
}
