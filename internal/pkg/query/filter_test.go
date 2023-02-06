package query

import (
	"testing"
)

func TestVersionFilter_getWhere(t *testing.T) {
	type fields struct {
		Application string
		Location    string
		Environment string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		args   int
	}{
		{"empty", fields{Application: "", Location: "", Environment: ""}, "", 0},
		{"loc", fields{Application: "", Location: "loc", Environment: ""}, "(location LIKE ?)", 1},
		{"app", fields{Application: "app", Location: "", Environment: ""}, "(application LIKE ?)", 1},
		{"env", fields{Application: "", Location: "", Environment: "env"}, "(environment LIKE ?)", 1},
		{"two", fields{Application: "app", Location: "", Environment: "env"}, "(application LIKE ? AND environment LIKE ?)", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := VersionFilter{
				Application: tt.fields.Application,
				Location:    tt.fields.Location,
				Environment: tt.fields.Environment,
			}
			got, got1 := f.getWhere()
			if got != tt.want {
				t.Errorf("VersionFilter.getWhere() got = %v, want %v", got, tt.want)
			}
			if len(got1) != tt.args {
				t.Errorf("VersionFilter.getWhere() got1 = %v, want %v", got1, tt.args)
			}
		})
	}
}
