package beacon

import (
	"testing"

	params2 "github.com/art-injener/iot-platform/internal/imitation/virtualdevice/beacon/params"
)

func TestVirtualBeaconImpl_checkResponseOnSuccessState(t *testing.T) {

	type args struct {
		response   string
		requestCRC uint8
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Check valid response",
			args: args{
				response:   "RE=0&CRC=20&STIME=210729140544&P=0000&&*180",
				requestCRC: 20,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Check error CRC response",
			args: args{
				response:   "RE=0&CRC=33&STIME=210728174951&P=1997&&*220",
				requestCRC: 33,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Check error response",
			args: args{
				response:   "ERR&CRC=196&&*19",
				requestCRC: 196,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VirtualBeaconImpl{}
			got, err := v.checkResponseOnSuccessState(tt.args.response, tt.args.requestCRC)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkResponseOnSuccessState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkResponseOnSuccessState() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVirtualBeaconImpl_wakeUpIntervalParser(t *testing.T) {

	tests := []struct {
		name  string
		data  string
		value int
	}{
		{
			name:  "Test parsing WUI",
			data:  "RE=0&CRC=228&STIME=210829224406&P=0000&WUI=10&&*170",
			value: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VirtualBeaconImpl{
				DevParam: params2.NewParameters("7911111111", 1),
			}
			if err := v.wakeUpIntervalParser(tt.data); err != nil {
				t.Errorf("wakeUpIntervalParser() error = %v", err)
			}
			if err := v.wakeUpIntervalParser(tt.data); err == nil {
				if v.DevParam.ConstParam.WUI != tt.value {
					t.Errorf("wakeUpIntervalParser() value = %v want = %v", v.DevParam.ConstParam.WUI, tt.value)
				}
			}
		})
	}
}
