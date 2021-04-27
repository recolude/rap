# File Format V1

```bash
protoc -I=v1 --go_out=./ v1/recording.proto
```

```proto
option go_package = "github.com/recolude/rap";

syntax = "proto3";

message Recording {
    string name = 1;
    map<string, string> metadata = 2;
    repeated SubjectRecording subjects = 3;
    repeated CustomEventCapture customEvents = 4;
}

message SubjectRecording {
    int32 id = 1;
    string name = 2;
    map<string, string> metadata = 3;
    repeated CustomEventCapture customEvents = 4;
    repeated LifeCycleEventCapture lifecycleEvents = 5;
    repeated VectorCapture capturedPositions = 6;
    repeated VectorCapture capturedRotations = 7;
}

message CustomEventCapture {
    float time = 1;
    string name = 2;
    string contents = 3;
    map<string, string> data = 4;
}

message LifeCycleEventCapture {
    float time = 1;
    EnumLifeType type = 2;
}

enum EnumLifeType {
    START = 0;
    ENABLE = 1;
    DISABLE = 2;
    DESTROY = 3;
}

message VectorCapture {
    float time = 1;
    float x = 2;
    float y = 3;
    float z = 4;
}

```