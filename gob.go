package av

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// AtributeValue types must be registered in gob
func init() {
	gob.Register(new(types.AttributeValueMemberS))
	gob.Register(new(types.AttributeValueMemberBS))
	gob.Register(new(types.AttributeValueMemberB))
	gob.Register(new(types.AttributeValueMemberL))
	gob.Register(new(types.AttributeValueMemberBOOL))
	gob.Register(new(types.AttributeValueMemberM))
	gob.Register(new(types.AttributeValueMemberN))
	gob.Register(new(types.AttributeValueMemberNS))
	gob.Register(new(types.AttributeValueMemberNULL))
	gob.Register(new(types.AttributeValueMemberSS))
}

func AVtoGobBytes(av types.AttributeValue) ([]byte, error) {
	bytesBuffer := new(bytes.Buffer)
	enc := gob.NewEncoder(bytesBuffer)

	err := enc.Encode(av)
	if err != nil {
		return nil, err
	}

	return bytesBuffer.Bytes(), nil
}

func AVtoGobStream(av types.AttributeValue) (io.Reader, error) {
	bytesBuffer := new(bytes.Buffer)
	enc := gob.NewEncoder(bytesBuffer)

	err := enc.Encode(av)
	if err != nil {
		return nil, err
	}

	return bytesBuffer, nil
}

func GobBytesToAv(b []byte, av types.AttributeValue) error {
	bytesBuffer := bytes.NewBuffer(b)
	dec := gob.NewDecoder(bytesBuffer)
	err := dec.Decode(av)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func GobStreamToAv(reader io.Reader, av types.AttributeValue) error {
	dec := gob.NewDecoder(reader)
	err := dec.Decode(av)
	if err != nil {
		return err
	}

	return nil
}
