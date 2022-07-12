package parsing

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/EliCDavis/vector"
	"github.com/Jeffail/gabs"
	"github.com/recolude/rap/format"
	"github.com/recolude/rap/format/collection/enum"
	"github.com/recolude/rap/format/collection/euler"
	"github.com/recolude/rap/format/collection/event"
	"github.com/recolude/rap/format/collection/position"
	"github.com/recolude/rap/format/io"
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

func parseReferenceFromJSON(jsonObj *gabs.Container) (format.BinaryReference, error) {
	name, err := parseRequiredStringKey(jsonObj, "reference", "name")
	if err != nil {
		return nil, err
	}

	uri, err := parseRequiredStringKey(jsonObj, "reference", "uri")
	if err != nil {
		return nil, err
	}

	parsedMetadata, err := parseMetadata(jsonObj)
	if err != nil {
		return nil, err
	}

	sizeNode := jsonObj.Path("size")
	if sizeNode == nil {
		return nil, fmt.Errorf("reference requires a size property")
	}

	size, err := strconv.Atoi(sizeNode.String())
	if err != nil {
		return nil, fmt.Errorf("reference size must be int")
	}

	if size < 0 {
		return nil, fmt.Errorf("reference size must be non-negative")
	}

	return io.NewBinaryReference(name, uri, uint64(size), parsedMetadata), nil
}

func parseReferences(jsonObj *gabs.Container) ([]format.BinaryReference, error) {
	referencesNode := jsonObj.Path("references")
	if referencesNode == nil {
		return []format.BinaryReference{}, nil
	}

	referenceNodes, err := referencesNode.Children()
	if err != nil {
		return nil, errors.New("references property should be an array")
	}

	_, properInternal := referencesNode.Data().([]interface{})
	if !properInternal {
		return nil, errors.New("references property should be an array")
	}

	references := make([]format.BinaryReference, len(referenceNodes))
	for i, child := range referenceNodes {
		if isNotJSONObj(child) {
			return nil, errors.New("references property should be array of reference objects")
		}

		parsed, err := parseReferenceFromJSON(child)
		if err != nil {
			return nil, err
		}

		references[i] = parsed
	}
	return references, nil
}

func parseRequiredStringKey(jsonObj *gabs.Container, thing string, key string) (string, error) {
	node := jsonObj.Path(key)
	if node == nil {
		return "", fmt.Errorf("%s requires %s", thing, key)
	}

	id, idParsed := node.Data().(string)
	if !idParsed {
		return "", fmt.Errorf("%s %s must be string", thing, key)
	}

	return id, nil
}

// func parseRequiredIntKey(jsonObj *gabs.Container, thing, key string) (int, error) {
// 	intNode := jsonObj.Path(key)
// 	if intNode == nil {
// 		return 0, fmt.Errorf("%s requires %s", thing, key)
// 	}

// 	parsedInt, err := strconv.Atoi(intNode.String())
// 	if err != nil {
// 		return 0, fmt.Errorf("%s %s must be int", thing, key)
// 	}

// 	return parsedInt, nil
// }

func parseRequiredFloatKey(jsonObj *gabs.Container, thing string, key string) (float64, error) {
	node := jsonObj.Path(key)
	if node == nil {
		return 0, fmt.Errorf("%s requires %s property", thing, key)
	}

	id, idParsed := node.Data().(float64)
	if !idParsed {
		return 0, fmt.Errorf("%s %s must be number", thing, key)
	}

	return id, nil
}

func parseCaptureTime(jsonCapture *gabs.Container) (float64, error) {
	timeNode := jsonCapture.Path("time")
	if timeNode == nil {
		return -1, fmt.Errorf("capture object must contain time property")
	}

	time, correct := timeNode.Data().(float64)
	if !correct {
		return -1, fmt.Errorf("capture object's time property must be a number")
	}

	return time, nil
}

func parseVector3(jsonObj *gabs.Container) (vector.Vector3, error) {
	x, err := parseRequiredFloatKey(jsonObj, "position capture", "x")
	if err != nil {
		return vector.Vector3Zero(), err
	}

	y, err := parseRequiredFloatKey(jsonObj, "position capture", "y")
	if err != nil {
		return vector.Vector3Zero(), err
	}

	z, err := parseRequiredFloatKey(jsonObj, "position capture", "z")
	if err != nil {
		return vector.Vector3Zero(), err
	}

	return vector.NewVector3(x, y, z), err
}

func parsePositionCollection(name string, jsonCaptures []*gabs.Container) (format.CaptureCollection, error) {
	captures := make([]position.Capture, len(jsonCaptures))

	for i, jsonCapture := range jsonCaptures {
		time, err := parseCaptureTime(jsonCapture)
		if err != nil {
			return nil, err
		}

		pos, err := parseVector3(jsonCapture.Path("data"))
		if err != nil {
			return nil, err
		}

		captures[i] = position.NewCapture(time, pos.X(), pos.Y(), pos.Z())
	}

	return position.NewCollection(name, captures), nil
}

func parseEulerCollection(name string, jsonCaptures []*gabs.Container) (format.CaptureCollection, error) {
	captures := make([]euler.Capture, len(jsonCaptures))

	for i, jsonCapture := range jsonCaptures {
		time, err := parseCaptureTime(jsonCapture)
		if err != nil {
			return nil, err
		}

		pos, err := parseVector3(jsonCapture.Path("data"))
		if err != nil {
			return nil, err
		}

		captures[i] = euler.NewEulerZXYCapture(time, pos.X(), pos.Y(), pos.Z())
	}

	return euler.NewCollection(name, captures), nil
}

func parseEnumCollection(name string, jsonCaptures []*gabs.Container) (format.CaptureCollection, error) {
	captures := make([]enum.Capture, len(jsonCaptures))

	allEnumEntriesMapping := make(map[string]int)
	allEnumEntries := make([]string, 0)

	for i, jsonCapture := range jsonCaptures {
		time, err := parseCaptureTime(jsonCapture)
		if err != nil {
			return nil, err
		}

		enumEntry, err := parseRequiredStringKey(jsonCapture, "enum capture", "data")
		if err != nil {
			return nil, err
		}

		enumIndex := len(allEnumEntriesMapping)
		if val, ok := allEnumEntriesMapping[enumEntry]; ok {
			enumIndex = val
		} else {
			allEnumEntries = append(allEnumEntries, enumEntry)
		}

		captures[i] = enum.NewCapture(time, enumIndex)
	}

	return enum.NewCollection(name, allEnumEntries, captures), nil
}

func parseEventCollection(name string, jsonCaptures []*gabs.Container) (format.CaptureCollection, error) {
	captures := make([]event.Capture, len(jsonCaptures))

	for i, jsonCapture := range jsonCaptures {
		time, err := parseCaptureTime(jsonCapture)
		if err != nil {
			return nil, err
		}

		dataNode := jsonCapture.Path("data")
		if dataNode == nil {
			return nil, errors.New("event capture requires data property object")
		}

		name, err := parseRequiredStringKey(dataNode, "event capture", "name")
		if err != nil {
			return nil, err
		}

		parsedMetadata, err := parseMetadata(dataNode)
		if err != nil {
			return nil, err
		}

		captures[i] = event.NewCapture(time, name, parsedMetadata)
	}

	return event.NewCollection(name, captures), nil
}

func parseCollectionFromJSON(jsonObj *gabs.Container) (format.CaptureCollection, error) {
	name, err := parseRequiredStringKey(jsonObj, "collection", "name")
	if err != nil {
		return nil, err
	}

	collectionType, err := parseRequiredStringKey(jsonObj, "collection", "type")
	if err != nil {
		return nil, err
	}

	capturesNode := jsonObj.Path("captures")
	if capturesNode == nil {
		return nil, errors.New("collection object requires captures property")
	}

	childCaptures, err := capturesNode.Children()
	if err != nil {
		return nil, errors.New("collection's captures property must be an array")
	}

	_, properInternal := capturesNode.Data().([]interface{})
	if !properInternal {
		return nil, errors.New("collection's captures property must be an array")
	}

	switch collectionType {
	case "recolude.position":
		return parsePositionCollection(name, childCaptures)

	case "recolude.euler":
		return parseEulerCollection(name, childCaptures)

	case "recolude.event":
		return parseEventCollection(name, childCaptures)

	case "recolude.enum":
		return parseEnumCollection(name, childCaptures)
	}
	return nil, fmt.Errorf("unrecognized collection type: '%s'", collectionType)
}

func parseCollectionsFromJSON(jsonObj *gabs.Container) ([]format.CaptureCollection, error) {
	collectionsNode := jsonObj.Path("collections")
	if collectionsNode == nil {
		return []format.CaptureCollection{}, nil
	}

	referenceNodes, err := collectionsNode.Children()
	if err != nil {
		return nil, errors.New("collections property should be an array")
	}

	_, properInternal := collectionsNode.Data().([]interface{})
	if !properInternal {
		return nil, errors.New("collections property should be an array")
	}

	references := make([]format.CaptureCollection, len(referenceNodes))
	for i, child := range referenceNodes {
		if isNotJSONObj(child) {
			return nil, errors.New("collections property should be array of collection objects")
		}

		parsed, err := parseCollectionFromJSON(child)
		if err != nil {
			return nil, err
		}

		references[i] = parsed
	}
	return references, nil
}

func parseRecordingFromJSON(jsonObj *gabs.Container) (format.Recording, error) {
	children, err := jsonObj.Children()
	if err != nil {
		return nil, err
	}

	if len(children) == 0 {
		return nil, errors.New("recording object can not be empty")
	}

	id, err := parseRequiredStringKey(jsonObj, "recording", "id")
	if err != nil {
		return nil, err
	}

	name, err := parseRequiredStringKey(jsonObj, "recording", "name")
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

	references, err := parseReferences(jsonObj)
	if err != nil {
		return nil, err
	}

	collections, err := parseCollectionsFromJSON(jsonObj)
	if err != nil {
		return nil, err
	}

	return format.NewRecording(
		id,
		name,
		collections,
		subRecordings,
		metadata,
		nil,
		references,
	), nil
}

func FromJSON(jsonData []byte) (format.Recording, error) {
	rootObj, err := gabs.ParseJSON(jsonData)
	if err != nil {
		return nil, err
	}
	return parseRecordingFromJSON(rootObj)
}
