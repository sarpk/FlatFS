package FlatFS

type UUIDToQuery struct {
	uuid string
	querykeyValue QueryKeyValue
}

type QueryKeyValue struct {
	keyValue map[string]string
}
