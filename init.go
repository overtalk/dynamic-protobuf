package dynamic

import (
	"errors"
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

var (
	RepetitivePackageNameErr = errors.New("repetitive package name")
)

type DynamicParser struct {
	packageName    string
	parser         *protoparse.Parser
	fileDesc       []*desc.FileDescriptor
	enumDescMap    map[string]*desc.EnumDescriptor
	messageDescMap map[string]*desc.MessageDescriptor
}

func NewDynamicParser() *DynamicParser {
	return &DynamicParser{
		parser:         &protoparse.Parser{},
		fileDesc:       nil,
		enumDescMap:    make(map[string]*desc.EnumDescriptor),
		messageDescMap: make(map[string]*desc.MessageDescriptor),
	}
}

func (dynamicParser *DynamicParser) SetImportPath(importPaths ...string) {
	dynamicParser.parser.ImportPaths = importPaths
}

func (dynamicParser *DynamicParser) ParseFiles(filenames ...string) error {
	fileDesc, err := dynamicParser.parser.ParseFiles(filenames...)
	if err != nil {
		return err
	}

	dynamicParser.fileDesc = fileDesc

	// get all desc
	for _, _desc := range fileDesc {
		if dynamicParser.packageName != "" && dynamicParser.packageName != _desc.GetPackage() {
			return RepetitivePackageNameErr
		}
		dynamicParser.packageName = _desc.GetPackage()

		for _, v := range _desc.GetMessageTypes() {
			dynamicParser.messageDescMap[v.GetName()] = v
		}

		for _, v := range _desc.GetEnumTypes() {
			dynamicParser.enumDescMap[v.GetName()] = v
		}
	}

	return nil
}

func (dynamicParser *DynamicParser) ParseToMap(messageName string, protocolID int32, data []byte) (string, []interface{}, error) {
	msgDesc, isExist := dynamicParser.messageDescMap[messageName]
	if !isExist {
		return "", nil, errors.New("message `" + messageName + "` is absent")
	}

	msg := dynamic.NewMessage(msgDesc)
	if err := proto.Unmarshal(data, msg); err != nil {
		return "", nil, err
	}

	return msgDesc.FindFieldByNumber(protocolID).GetMessageType().GetName(), dynamicParser.toMap(msgDesc, msg), nil
}

func (dynamicParser *DynamicParser) toMap(msgDesc *desc.MessageDescriptor, msg *dynamic.Message) []interface{} {
	var ret []interface{}

	for _, fieldDesc := range msgDesc.GetFields() {
		fieldValue := msg.GetField(fieldDesc)
		switch fieldDesc.GetType() {
		case dpb.FieldDescriptorProto_TYPE_MESSAGE:
			subMsgDesc, isExist := dynamicParser.messageDescMap[fieldDesc.GetMessageType().GetName()]
			if isExist {
				if t, flag := fieldValue.(*dynamic.Message); !flag {
					continue
				} else {
					ret = append(ret, dynamicParser.toMap(subMsgDesc, t)...)
				}
			}

		default:
			ret = append(ret, fieldValue)
		}
	}

	return ret
}
