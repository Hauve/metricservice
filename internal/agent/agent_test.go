package agent

//func TestNew(t *testing.T) {
//	tt := struct {
//		name string
//		want *MyAgent
//	}{
//		name: "Single test New",
//		want: &MyAgent{
//			storage:        storage.NewMemStorage(),
//			address:        "localhost:8080",
//			reportInterval: 10,
//			pollInterval:   2,
//		},
//	}
//	t.Run(tt.name, func(t *testing.T) {
//		if got := New(); !reflect.DeepEqual(got, tt.want) {
//			t.Errorf("New() = %v, want %v", got, tt.want)
//		}
//	})
//}
//
//// BAD TEST
//func TestMyAgent_collectMetrics(t *testing.T) {
//	type fields struct {
//		storage handlers.Storage
//	}
//	tt := struct {
//		name   string
//		fields fields
//	}{
//		name: "Single test",
//		fields: fields{
//			storage: storage.NewMemStorage(),
//		},
//	}
//
//	t.Run(tt.name, func(t *testing.T) {
//		ag := &MyAgent{
//			storage: tt.fields.storage,
//		}
//		ag.collectMetrics()
//
//		gotGaugeKeys := ag.storage.GetGaugeKeys()
//		temp := "HeapAlloc HeapSys MSpanInuse GCCPUFraction StackInuse GCSys MSpanSys Mallocs OtherSys HeapIdle HeapInuse HeapObjects NumForcedGC PauseTotalNs Sys Alloc Frees LastGC NumGC StackSys RandomValue Lookups MCacheInuse MCacheSys NextGC BuckHashSys HeapReleased"
//		mustGaugeKeys := strings.Split(temp, " ")
//
//		for _, must := range mustGaugeKeys {
//			assert.Contains(t, *gotGaugeKeys, must, "Value that must be in storage is missed")
//		}
//	})
//
//}
