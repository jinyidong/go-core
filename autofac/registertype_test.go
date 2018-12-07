package autofac

import (
	"fmt"
	"github.com/jinyidong/go-core/util"
	"testing"
)

type testStruct struct {
	Name string
	Sex  int
}

func TestGetStruct(t *testing.T) {

	tempTestStruct := util.StructToJson(testStruct{
		Name: "1",
		Sex:  1,
	})

	RegisterType((*testStruct)(nil))

	structName := "testStruct"

	s, ok := NewStructByType(structName)
	if !ok {
		return
	}

	t1, ok1 := s.(testStruct)
	if !ok1 {
		return
	}

	err := util.ByteToStruct([]byte(tempTestStruct), &t1)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(s)
}
