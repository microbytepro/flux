package dubbo

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/bytepowered/flux/flux-node"
	"github.com/bytepowered/flux/flux-node/common"
	"github.com/bytepowered/flux/flux-node/context"
	"github.com/bytepowered/flux/flux-node/ext"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultAssembleFunc(t *testing.T) {
	ext.SetArgumentLookupFunc(common.LookupMTValue)
	serializer := flux.NewJsonSerializer()
	ext.RegisterSerializer(ext.TypeNameSerializerDefault, serializer)
	ext.RegisterSerializer(ext.TypeNameSerializerJson, serializer)
	cases := []struct {
		arguments      []flux.Argument
		expectedTypes  []string
		expectedValues []hessian.Object
	}{
		{
			arguments: []flux.Argument{
				ext.NewStringArgument("username"),
				ext.NewIntegerArgument("year"),
				ext.NewStringArgument("stringmap"),
				func() flux.Argument {
					arg := ext.NewComplexArgument("net.bytepowreed.test.POJO", "pojo")
					arg.Fields = []flux.Argument{
						ext.NewStringArgument("username"),
						ext.NewIntegerArgument("year"),
						ext.NewHashMapArgument("hashmap"),
					}
					return arg
				}(),
			},
			expectedTypes: []string{"java.lang.String", "java.lang.Integer", "java.lang.String", "net.bytepowreed.test.POJO"},
			expectedValues: []hessian.Object{
				"yongjiachen",
				2020,
				"{\"int\":123,\"key\":\"value\"}",
				map[string]interface{}{
					"class":    "net.bytepowreed.test.POJO",
					"hashmap":  map[string]interface{}{"int": 123, "key": "value"},
					"username": "yongjiachen",
					"year":     2020,
				},
			},
		},
	}
	assert := assert2.New(t)
	ctx := context.NewMockWith("@rid", map[string]interface{}{
		"stringmap": map[string]interface{}{
			"key": "value",
			"int": 123,
		},
		"hashmap": map[string]interface{}{
			"key": "value",
			"int": 123,
		},
		"user": map[string]interface{}{
			"username": "yongjiachen",
			"year":     2020,
		},
		"username": "yongjiachen",
		"year":     2020,
	})
	for _, tcase := range cases {
		types, values, err := DefaultArgAssembleFunc(tcase.arguments, ctx)
		assert.Nil(err)
		assert.Equal(tcase.expectedTypes, types, "types matches")
		assert.Equal(tcase.expectedValues, values, "values matches")
	}
}
