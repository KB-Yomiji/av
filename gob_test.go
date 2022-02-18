package av_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kb-yomiji/av"
)

type TestStruct struct {
	TestBench string
}

func TestAVtoGobBytes(t *testing.T) {
	newAv, err := attributevalue.Marshal(
		TestStruct{TestBench: "Some longer string than usual"},
	)
	if err != nil {
		t.Fatalf("marshal attribute failure: %s", err)
	}

	gb, err := av.AVtoGobBytes(newAv)
	if err != nil {
		t.Fatalf("gob bytes error: %s", err)
	}

	if len(gb) == 0 {
		t.Fatalf("len(gb) == 0")
	}
}

func TestGobBytesToAv(t *testing.T) {
	av1, err := attributevalue.Marshal(
		TestStruct{TestBench: "Some longer string than usual"},
	)
	if err != nil {
		t.Fatalf("marshal attribute failure: %s", err)
	}

	gb, err := av.AVtoGobBytes(av1)
	if err != nil {
		t.Fatalf("gob bytes error: %s", err)
	}

	fmt.Println("AV", av1)

	if len(gb) == 0 {
		t.Fatalf("len(gb) == 0")
	}

	newAv, err := attributevalue.Marshal(TestStruct{})
	if err != nil {
		t.Fatalf("marshal attribute failure #2: %s", err)
	}

	err = av.GobBytesToAv(gb, newAv)
	if err != nil {
		t.Fatalf("gob bytes conversion error: %s", err)
	}

	// check the internal structure of the attribute
	value := newAv.(*types.AttributeValueMemberM).Value["TestBench"].(*types.AttributeValueMemberS).Value
	if "Some longer string than usual" != value {
		t.Fatalf("expected Some longer string than usual, got; %s", value)
	}
}

func TestAVtoGobStream(t *testing.T) {
	av1, err := attributevalue.Marshal(
		TestStruct{TestBench: "Some longer string than usual"},
	)
	if err != nil {
		t.Fatalf("marshal attribute failure: %s", err)
	}

	gb, err := av.AVtoGobStream(av1)
	if err != nil {
		t.Fatalf("gob stream error: %s", err)
	}

	if b, err := ioutil.ReadAll(gb); len(b) == 0 {
		t.Fatalf("len(b) == 0")
	} else if err != nil {
		t.Fatalf("error reading gb: %s", err)
	}
}

func BenchmarkAVtoGobBytes(b *testing.B) {
	av1, err := attributevalue.Marshal(
		struct {
			TestBench string
		}{TestBench: "Some longer string than usual"},
	)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		_, _ = av.AVtoGobBytes(av1)
	}
}

func BenchmarkGobBytesToAv(b *testing.B) {
	av1, err := attributevalue.Marshal(
		struct {
			TestBench string
		}{TestBench: "Some longer string than usual"},
	)
	if err != nil {
		panic(err)
	}

	gb, err := av.AVtoGobBytes(av1)
	for i := 0; i < b.N; i++ {
		_ = av.GobBytesToAv(gb, av1)
	}
}

func BenchmarkAVtoGobStream(b *testing.B) {
	av1, err := attributevalue.Marshal(
		struct {
			TestBench string
		}{TestBench: "Some longer string than usual"},
	)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		_, _ = av.AVtoGobStream(av1)
	}
}

func BenchmarkGobStreamToAv(b *testing.B) {
	av1, err := attributevalue.Marshal(
		struct {
			TestBench string
		}{TestBench: "Some longer string than usual"},
	)
	if err != nil {
		panic(err)
	}

	gs, err := av.AVtoGobStream(av1)

	for i := 0; i < b.N; i++ {
		_ = av.GobStreamToAv(gs, av1)
	}
}

func BenchmarkGobCombo(b *testing.B) {
	b.Run("BenchmarkAVtoGobBytes", BenchmarkAVtoGobBytes)
	b.Run("BenchmarkGobBytesToAv", BenchmarkGobBytesToAv)
	b.Run("BenchmarkAVtoGobStream", BenchmarkAVtoGobStream)
	b.Run("BenchmarkGobStreamToAv", BenchmarkGobStreamToAv)
}
