package parsing

import (
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/metadata"
)

func isJSONObj(container *gabs.Container) bool {
	_, err := container.ChildrenMap()
	return err == nil
}

func isNotJSONObj(container *gabs.Container) bool {
	return !isJSONObj(container)
}

func parseMetadata(jsonObj *gabs.Container) (metadata.Block, error) {
	metadataNode := jsonObj.Path("metadata")
	if metadataNode == nil {
		return metadata.EmptyBlock(), nil
	}

	_, err := metadataNode.ChildrenMap()
	if err != nil {
		return metadata.EmptyBlock(), errors.New("metadata should be object")
	}

	resultingMetadata := metadata.NewMetadataProperty(metadata.EmptyBlock())
	resultingMetadata.UnmarshalJSON([]byte(metadataNode.String()))

	return resultingMetadata.Block(), nil
}

func parseSubRecordings(jsonObj *gabs.Container) ([]format.Recording, error) {
	recordingsNode := jsonObj.Path("recordings")
	if recordingsNode == nil {
		return []format.Recording{}, nil
	}

	subRecordingNodes, err := recordingsNode.Children()
	if err != nil {
		return nil, errors.New("recordings property should be array")
	}

	_, properInternal := recordingsNode.Data().([]interface{})
	if !properInternal {
		return nil, errors.New("recordings property should be array")
	}

	subRecordings := make([]format.Recording, len(subRecordingNodes))
	for i, child := range subRecordingNodes {
		if isNotJSONObj(child) {
			return nil, errors.New("recordings property should be array of recording objects")
		}

		parsed, err := parseRecordingFromJSON(child)
		if err != nil {
			return nil, err
		}

		subRecordings[i] = parsed
	}
	return subRecordings, nil
}

func parseRequiredStringKey(jsonObj *gabs.Container, key string) (string, error) {
	node := jsonObj.Path(key)
	if node == nil {
		return "", errors.New("recording requires " + key)
	}

	id, idParsed := node.Data().(string)
	if !idParsed {
		return "", errors.New("recording " + key + " must be string")
	}

	return id, nil
}

func parseRecordingFromJSON(jsonObj *gabs.Container) (format.Recording, error) {
	children, err := jsonObj.Children()
	if err != nil {
		return nil, err
	}

	if len(children) == 0 {
		return nil, errors.New("recording object can not be empty")
	}

	id, err := parseRequiredStringKey(jsonObj, "id")
	if err != nil {
		return nil, err
	}

	name, err := parseRequiredStringKey(jsonObj, "name")
	if err != nil {
		return nil, err
	}

	subRecordings, err := parseSubRecordings(jsonObj)
	if err != nil {
		return nil, err
	}

	metadata, err := parseMetadata(jsonObj)
	if err != nil {
		return nil, err
	}

	return format.NewRecording(
		id,
		name,
		nil,
		subRecordings,
		metadata,
		nil,
		nil,
	), nil
}

func FromJSON(jsonData []byte) (format.Recording, error) {
	rootObj, err := gabs.ParseJSON(jsonData)
	if err != nil {
		return nil, err
	}
	return parseRecordingFromJSON(rootObj)
}
