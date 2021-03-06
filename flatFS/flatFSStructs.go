package FlatFS

type UUIDToQuery struct {
	uuid string
	querykeyValue QueryKeyValue
}

type QueryKeyValue struct {
	keyValue map[string]string
}

type QueryType struct {
	addSpec bool
	querySpec bool
	replaceSpec bool
	deleteSpec bool
	fileSpec bool
	emptyType bool
}
