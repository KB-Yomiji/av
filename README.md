# AWS Transmit AV

### Purpose
Convert aws AttributeValues to a wire format for transmission, save, distribution

### Use
To JSON:
```go
av, _ := attributevalue.Marshal(
    Mytest{
        Matrix: "testValue",
        Test:   []byte{1, 2, 4},
    },
)

js, _ := av.ToJSON(av)
// output:
//    {
//        "Matrix": "testValue",
//        "Test": "AQIE"
//    }
//
```
From JSON:
```go
av1, err = av.FromJSON(js, &Mytest{})
if err != nil {
    t.Fatalf("error occurred: %s", err)
}
```
