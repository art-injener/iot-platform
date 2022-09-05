package params

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/art-injener/iot-platform/pkg/models/device"
)

func TestParameters_Serialize(t *testing.T) { //nolint:paralleltest
	type fields struct {
		ConstParam   *device.SettingsModel
		MutableParam *device.ParamsModel
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "test valid serialize",
			fields: fields{
				ConstParam:   device.NewDeviceSettings("79262279262", 1),
				MutableParam: device.NewDeviceParams("79262279262"),
			},
			want: "ID=79262279262&" +
				"VER=0.0.1&" +
				"TZ=220&" +
				"WUI=1&" +
				"GPST=3&" +
				"STIME=210726102335&" +
				"BAL=20.00&" +
				"TE=25.0&" +
				"VB=100%(5.95V)&" +
				"IC=1&" +
				"SQ=69&" +
				"LA=55.758436&" +
				"LAD=N&" +
				"LO=37.551049&" +
				"LOD=E&" +
				"SPD=0&" +
				"DM=0.0&" +
				"GT=210726072255&&",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parameters{
				ConstParam:   tt.fields.ConstParam,
				MutableParam: tt.fields.MutableParam,
			}
			got := p.Serialize()

			// обновляем системное и utc время
			tt.want = strings.ReplaceAll(tt.want, "STIME=210726102335", fmt.Sprintf("STIME=%s",
				time.Unix(p.MutableParam.SystemTime, 0).Local().Format("060102150405")))
			tt.want = strings.ReplaceAll(tt.want, "GT=210726072255", fmt.Sprintf("GT=%s",
				time.Unix(p.MutableParam.SystemTime, 0).UTC().Format("060102150405")))

			if got != tt.want {
				t.Errorf("Serialize() \n\t got = %v, \n\t want= %v", got, tt.want)
			}
		})
	}
}
