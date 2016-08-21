// In memory attribute mapper

package first

import (
	"log"
)

type MemAttrMapper struct {
	AttrMapper
}

func NewMemAttrMapper() *MemAttrMapper {
	memAttrMapper := &MemAttrMapper{}
	return memAttrMapper
}

func (attrMapper *MemAttrMapper) Foo() string {
	log.Println("Mocking middleware")
	return ""
}

func init() {
	AttrMapperManagerInjector = *NewAttrMapperManager()
	AttrMapperManagerInjector.Set("default", NewMemAttrMapper())
}
