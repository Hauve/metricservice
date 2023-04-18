package storage

import (
	"reflect"
	"testing"
)

func TestMemStorage_GetCounter(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue int64
		wantOk    bool
	}{
		{
			name: "positive test with created counter",
			fields: fields{
				gauge: make(map[string]float64),
				counter: map[string]int64{
					"name": 5,
				},
			},
			args:      args{"name"},
			wantValue: 5,
			wantOk:    true,
		}, {
			name: "negative test with not created counter",
			fields: fields{
				gauge: make(map[string]float64),
				counter: map[string]int64{
					"name": 5,
				},
			},
			args:      args{"name2"},
			wantValue: 0,
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			gotValue, gotOk := st.GetCounter(tt.args.key)
			if gotValue != tt.wantValue {
				t.Errorf("GetCounter() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("GetCounter() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMemStorage_GetCounterKeys(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	tests := []struct {
		name   string
		fields fields
		want   *[]string
	}{
		{
			name: "positive test",
			fields: fields{
				gauge: map[string]float64{
					"name": 5,
				},
				counter: map[string]int64{
					"name":  5,
					"name1": 7,
					"name2": 4,
					"name3": 6,
				},
			},
			want: &[]string{"name", "name1", "name2", "name3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			if got := st.GetCounterKeys(); !reflect.DeepEqual(got, tt.want) {
				//GetGaugeKeys не гарантирует определённый порядок элементов
				equal := false
				for gotKey := range got {
					for wantKey := range *tt.want {
						if gotKey == wantKey {
							equal = true
							break
						}
					}
					if !equal {
						t.Errorf("GetGaugeKeys() = %v, want %v", got, tt.want)
						return
					}
					equal = false
				}
			}
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue float64
		wantOk    bool
	}{
		{
			name: "positive test with created gauge",
			fields: fields{
				gauge: map[string]float64{
					"name": 5.1,
				},
				counter: make(map[string]int64),
			},
			args:      args{"name"},
			wantValue: 5.1,
			wantOk:    true,
		}, {
			name: "negative test with not created gauge",
			fields: fields{
				gauge: map[string]float64{
					"name": 5.1,
				},
				counter: make(map[string]int64),
			},
			args:      args{"name2"},
			wantValue: float64(0),
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			gotValue, gotOk := st.GetGauge(tt.args.key)
			if gotValue != tt.wantValue {
				t.Errorf("GetGauge() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("GetGauge() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestMemStorage_GetGaugeKeys(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	tests := []struct {
		name   string
		fields fields
		want   *[]string
	}{
		{
			name: "positive test",
			fields: fields{
				gauge: map[string]float64{
					"name":  5.11,
					"name1": 7.52,
					"name2": 4.7743,
					"name3": 6.324,
				},
				counter: map[string]int64{
					"name": 5,
				},
			},
			want: &[]string{"name", "name1", "name2", "name3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			if got := st.GetGaugeKeys(); !reflect.DeepEqual(got, tt.want) {
				//GetGaugeKeys не гарантирует определённый порядок элементов
				equal := false
				for gotKey := range got {
					for wantKey := range *tt.want {
						if gotKey == wantKey {
							equal = true
							break
						}
					}
					if !equal {
						t.Errorf("GetGaugeKeys() = %v, want %v", got, tt.want)
						return
					}
					equal = false
				}
			}
		})
	}
}

func TestMemStorage_SetCounter(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		key string
		val int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive test",
			fields: fields{
				gauge: map[string]float64{
					"name": 5.1,
				},
				counter: map[string]int64{
					"name": 55,
				},
			},
			args: args{"name", 55},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			st.AddCounter(tt.args.key, tt.args.val)
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		key string
		val float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive test",
			fields: fields{
				gauge: map[string]float64{
					"name": 5.1,
				},
				counter: map[string]int64{
					"name": 55,
				},
			},
			args: args{"name", 55.25},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			st.SetGauge(tt.args.key, tt.args.val)
		})
	}
}
