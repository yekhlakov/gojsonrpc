package server

import (
	"math"
	"reflect"
	"testing"

	"github.com/yekhlakov/gojsonrpc/common"
)

type test_ExtractHandler struct{}

func (h test_ExtractHandler) Method_twoParams(params struct{}, _ string) (response struct{}) {
	return
}

func (h test_ExtractHandler) Method_twoReturns(params struct{}) (response struct{}, jsonRpcError common.Error) {
	return
}

func (h test_ExtractHandler) Method_wrongReturn1(params struct{}) (response struct{}, jsonRpcError common.Request, err error) {
	return
}

func (h test_ExtractHandler) Method_wrongReturn2(params struct{}) (response struct{}, jsonRpcError common.Error, err common.Error) {
	return
}

func (h test_ExtractHandler) Method_wrongFourReturns(params struct{}) (response struct{}, jsonRpcError common.Error, err error, err2 error) {
	return
}

func (h test_ExtractHandler) Handle_test1(params struct{}) (response struct{}, jsonRpcError common.Error, err error) {
	return
}
func (h test_ExtractHandler) Handle_test2(params struct{}) (response struct{}, jsonRpcError common.Error, err error) {
	return
}
func (h test_ExtractHandler) Handle_test3(params struct{}) (response struct{}, jsonRpcError common.Error, err error) {
	return
}

// Testing method extraction
func TestExtractMethods(t *testing.T) {

	handler := test_ExtractHandler{}

	m0 := ExtractMethods(handler, "Method_")

	if len(m0) != 0 {
		t.Errorf("Extracted more methods than needed")
	}

	m3 := ExtractMethods(handler, "Handle_")

	if len(m3) > 3 {
		t.Errorf("Extracted more methods than needed")
	} else if len(m3) < 3 {
		t.Errorf("Extracted less methods than needed")
	}
}

type test_EmptyHandler struct{}

func (c test_EmptyHandler) Handle_empty(params struct{}) (response struct{}, jsonRpcError common.Error, err error) {
	return
}

type test_ConstHandler struct{}

func (c test_ConstHandler) Handle_const(params struct{}) (response struct {
	Value string `json:"value"`
}, jsonRpcError common.Error, err error) {
	response.Value = "test"
	return
}

type test_PassHandler struct{}

func (c test_PassHandler) Handle_pass(params struct {
	Name string `json:"name"`
}) (response struct {
	Value string `json:"value"`
}, jsonRpcError common.Error, err error) {
	response.Value = params.Name
	return
}

type test_WrongHandler struct{}

func (c test_WrongHandler) Handle_wrong(params struct{}) (response struct {
	Value float64 `json:"value,int"`
}, jsonRpcError common.Error, err error) {
	response.Value = math.Inf(0)
	return
}

type test_ErrorHandler struct{}

func (c test_ErrorHandler) Handle_error(params struct{}) (response struct{}, jsonRpcError common.Error, err error) {
	jsonRpcError.Code = "666"
	jsonRpcError.Message = "error"

	return
}

func TestJsonRpcMethod_BindParams(t *testing.T) {
	m := ExtractMethods(test_PassHandler{}, "Handle_")

	// OK
	params := []byte(`{"name":"test"}`)
	v, e := m[0].BindParams(params)
	if e != nil {
		t.Errorf("Could not bind params")
		t.Errorf(e.Error())
	} else if v[1].Field(0).String() != "test" {
		t.Errorf("Params weren't bound properly")
	}

	// Invalid
	params = []byte(`{"name":[]}`)
	_, e = m[0].BindParams(params)
	if e == nil {
		t.Errorf("Params were bound when they shouldn't be")
	}
}

// Testing error extraction
func TestJsonRpcMethod_ExtractError(t *testing.T) {
	m := ExtractMethods(test_PassHandler{}, "Handle_")

	v := reflect.ValueOf(common.Error{
		Code:    "666",
		Message: "test",
		Data:    nil,
	})

	rawError, err := m[0].ExtractError([]reflect.Value{v, v})
	if err != nil {
		t.Errorf("Could not extract error")
		t.Errorf(err.Error())
	} else if string(rawError) != `{"code":666,"message":"test"}` {
		t.Errorf("Error was not extracted properly")
	}

	invalid := reflect.ValueOf(common.Request{})
	rawError, err = m[0].ExtractError([]reflect.Value{v, invalid})
	if rawError != nil {
		t.Errorf("Error was extracted when it shouldn't")
	}
	if err == nil {
		t.Errorf("Error extraction generated no error while it should")
	}

	empty := reflect.ValueOf(common.Error{})
	rawError, err = m[0].ExtractError([]reflect.Value{v, empty})
	if rawError != nil {
		t.Errorf("Error was extracted when it shouldn't")
	}
	if err != nil {
		t.Errorf("Error extraction errored while it shouldn't")
	}
}

// Testing result extraction
func TestJsonRpcMethod_ExtractResult(t *testing.T) {
	m := ExtractMethods(test_PassHandler{}, "Handle_")

	v := reflect.ValueOf(struct {
		Value string `json:"value"`
	}{"test"})

	rawResult, err := m[0].ExtractResult([]reflect.Value{v, v})
	if err != nil {
		t.Errorf("Could not extract result")
		t.Errorf(err.Error())
	} else if string(rawResult) != `{"value":"test"}` {
		t.Errorf("The result was not extracted properly")
	}

	invalid := reflect.ValueOf(common.Error{})
	rawResult, err = m[0].ExtractResult([]reflect.Value{invalid, v})
	if rawResult != nil {
		t.Errorf("Result was extracted while it shouldn't")
	}
	if err == nil {
		t.Errorf("Result extraction generated no error while it should")
	}

}
