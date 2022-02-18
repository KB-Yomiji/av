package av_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kb-yomiji/av"
)

type Mytest struct {
	Matrix string
	Test   []byte
}

func TestAVtoJSON(t *testing.T) {
	av1, _ := attributevalue.Marshal(
		Mytest{
			Matrix: "testValue",
			Test:   []byte{1, 2, 4},
		},
	)

	js, err := av.ToJSON(av1)
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	if !strings.Contains(string(js), "testValue") {
		t.Fail()
	}
}

func TestAVtoJSON_SoloVar(t *testing.T) {
	av1, _ := attributevalue.Marshal(
		"valueTest",
	)

	js, err := av.ToJSON(av1)
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	if string(js) != "\"valueTest\"" {
		t.Fatalf("Expected \"valueTest\", got %s", string(js))
	}
}

func TestAVtoJSON_SoloVarBinary(t *testing.T) {
	bVal := []byte{1, 2, 4, 5}
	av1, _ := attributevalue.Marshal(
		bVal,
	)

	js, err := av.ToJSON(av1)
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	expected := make([]byte, base64.StdEncoding.EncodedLen(len(bVal)))

	base64.StdEncoding.Encode(expected, bVal)

	actual := js[1 : len(js)-1]

	if !bytes.Equal(expected, actual) {
		t.Fatalf("Expected %b, got %b", expected, actual)
	}
}

func TestJSONToAv_Unguided(t *testing.T) {
	av1, _ := attributevalue.Marshal(
		Mytest{
			Matrix: "testValue",
			Test:   []byte{1, 2, 4},
		},
	)

	js, err := av.ToJSON(av1)
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	if !strings.Contains(string(js), "testValue") {
		t.Fail()
	}

	av1, err = av.FromJSON(js, nil)
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	value := av1.(*types.AttributeValueMemberM).Value["Matrix"].(*types.AttributeValueMemberS).Value
	if value != "testValue" {
		t.Fatalf("expecting testValue, got: %s", value)
	}
}
func TestJSONToAv_Guided(t *testing.T) {
	av1, _ := attributevalue.Marshal(
		Mytest{
			Matrix: "testValue",
			Test:   []byte{1, 2, 4},
		},
	)

	js, err := av.ToJSON(av1)
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	if !strings.Contains(string(js), "testValue") {
		t.Fail()
	}

	av1, err = av.FromJSON(js, &Mytest{})
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	value := av1.(*types.AttributeValueMemberM).Value["Matrix"].(*types.AttributeValueMemberS).Value
	if value != "testValue" {
		t.Fatalf("expecting testValue, got: %s", value)
	}
}

func BenchmarkRegularJson(b *testing.B) {
	ourStruct := struct {
		TestBench string
	}{TestBench: "Some longer string than usual"}

	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(ourStruct)
	}
}

func BenchmarkRegularAVMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = attributevalue.Marshal(
			struct {
				TestBench string
			}{TestBench: "Some longer string than usual"},
		)
	}
}

func BenchmarkAVtoInterface(b *testing.B) {
	av1, err := attributevalue.Marshal(
		struct {
			TestBench string
		}{TestBench: "Some longer string than usual"},
	)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		_ = av.ToInterface(av1)
	}
}

func BenchmarkAVtoJSON_1_Value(b *testing.B) {
	av1, err := attributevalue.Marshal(
		struct {
			TestBench string
		}{TestBench: "Some longer string than usual"},
	)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		_, _ = av.ToJSON(av1)
	}
}

func BenchmarkJsonToAv_1_Value_Unguided(b *testing.B) {
	av1, err := attributevalue.Marshal(
		struct {
			TestBench string
		}{TestBench: "Some longer string than usual"},
	)
	if err != nil {
		panic(err)
	}

	js, err := av.ToJSON(av1)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		_, _ = av.FromJSON(js, nil)
	}
}

func BenchmarkJsonToAv_1_Value_Guided(b *testing.B) {
	av1, err := attributevalue.Marshal(
		struct {
			TestBench string
		}{TestBench: "Some longer string than usual"},
	)
	if err != nil {
		panic(err)
	}

	js, err := av.ToJSON(av1)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		_, _ = av.FromJSON(js, struct{ TestBench string }{TestBench: ""})
	}
}

func BenchmarkJSONCombo(b *testing.B) {
	b.Run("BenchmarkRegularJson", BenchmarkRegularJson)
	b.Run("BenchmarkRegularAVMarshal", BenchmarkRegularAVMarshal)
	b.Run("BenchmarkAVtoInterface", BenchmarkAVtoInterface)
	b.Run("BenchmarkAVtoJSON_1_Value", BenchmarkAVtoJSON_1_Value)
	b.Run("BenchmarkJsonToAv_1_Value_Unguided", BenchmarkJsonToAv_1_Value_Unguided)
	b.Run("BenchmarkJsonToAv_1_Value_Guided", BenchmarkJsonToAv_1_Value_Guided)
}
