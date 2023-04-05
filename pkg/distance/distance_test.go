package distance_test

import (
	"testing"

	"github.com/alaleks/geospace/pkg/distance"
)

func TestCalcDistance(t *testing.T) {
	krd := struct {
		Lat float64
		Lon float64
	}{
		Lat: 45.04484,
		Lon: 38.97603,
	}

	msk := struct {
		Lat float64
		Lon float64
	}{
		Lat: 55.75222,
		Lon: 37.61556,
	}

	tests := []struct {
		function func(lat1, lon1, lat2, lon2 float64) float64
		name     string
		city1    struct {
			Lat float64
			Lon float64
		}
		city2 struct {
			Lat float64
			Lon float64
		}
		dist int
	}{
		{
			name:     "Calculation of distance to CalcGreatCirlcle",
			function: distance.CalcGreatCircle,
			city1:    krd,
			city2:    msk,
			dist:     1194,
		},
		{
			name:     "Calculation of distance to CalcHaversine",
			function: distance.CalcHaversine,
			city1:    krd,
			city2:    msk,
			dist:     1194,
		},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.function(tt.city1.Lat, tt.city1.Lon, tt.city2.Lat, tt.city2.Lon)
			if int(result) != tt.dist {
				t.Errorf(`the distance by coordinates was calculated incorrectly, 
				it must be %d:, and the calculation returns: %d`, tt.dist, int(result))
			}
		})
	}
}

func BenchmarkCalcDistance(b *testing.B) {
	krd := struct {
		Lat float64
		Lon float64
	}{
		Lat: 45.04484,
		Lon: 38.97603,
	}

	msk := struct {
		Lat float64
		Lon float64
	}{
		Lat: 55.75222,
		Lon: 37.61556,
	}

	b.ResetTimer()

	b.Run("Calculation of distance to CalcGreatCirlcle", func(b *testing.B) {
		_ = distance.CalcGreatCircle(krd.Lat, krd.Lon, msk.Lat, msk.Lon)
	})

	b.ResetTimer()

	b.Run("Calculation of distance to CalcHaversine", func(b *testing.B) {
		_ = distance.CalcHaversine(krd.Lat, krd.Lon, msk.Lat, msk.Lon)
	})
}
