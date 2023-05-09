package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func validateError(t *testing.T, err error, errCode codes.Code) {
	if err != nil {
		if errCode == codes.OK {
			st, ok := status.FromError(err)
			if ok {
				t.Error("not expected error: ", err.Error(), st.Code())
			} else {
				t.Error("not expected error: ", err.Error())
			}
		} else {
			st, ok := status.FromError(err)
			if !ok {
				t.Error("invalid error, must be code:", errCode)
				return
			}

			if st.Code() != errCode {
				t.Errorf("error code: got %d, waited %d", st.Code(), errCode)
			}
		}
	}
}
func validateEmptyResponse(t *testing.T, got, wait *emptypb.Empty) {
	var err error
	if wait == nil {
		if got == nil {
			return
		}
		t.Error("response must be empty")
		var data []byte
		data, err = json.MarshalIndent(got, "", "	")
		if err != nil {
			t.Error("unable marshal response")
		} else {
			t.Log(string(data))
		}
		return
	}

	if got == nil {
		t.Error("response must be not nil")
		return
	}

	if !reflect.DeepEqual(got, wait) {
		t.Error("response: expected", wait, "received", got)
	}
}
func validateResponse[R any](t *testing.T, got, wait *R) {
	var err error
	if wait == nil {
		if got == nil {
			return
		}
		t.Error("response must be empty")
		var data []byte
		data, err = json.MarshalIndent(got, "", "	")
		if err != nil {
			t.Error("unable marshal response")
		} else {
			t.Log(string(data))
		}
		return
	}

	if got == nil {
		t.Error("response must be not nil")
		return
	}

	if !reflect.DeepEqual(got, wait) {
		t.Error("response: expected", wait, "received", got)
	}
}
