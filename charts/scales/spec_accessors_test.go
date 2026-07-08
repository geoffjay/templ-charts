package scales

import "testing"

func TestSpecReverse(t *testing.T) {
	cases := []struct {
		name string
		spec ScaleSpec
		want bool
	}{
		{"linear-forward", ScaleLinearSpec{}, false},
		{"linear-reverse", ScaleLinearSpec{Reverse: true}, true},
		{"log-reverse", ScaleLogSpec{Reverse: true}, true},
		{"symlog-reverse", ScaleSymlogSpec{Reverse: true}, true},
		{"band-none", ScaleBandSpec{}, false},
		{"time-none", ScaleTimeSpec{}, false},
	}
	for _, c := range cases {
		if got := SpecReverse(c.spec); got != c.want {
			t.Errorf("%s: SpecReverse = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestSpecMinIsAuto(t *testing.T) {
	cases := []struct {
		name string
		spec ScaleSpec
		want bool
	}{
		{"linear-fixed", ScaleLinearSpec{Min: FloatVal(0)}, false},
		{"linear-auto", ScaleLinearSpec{Min: AutoFloat()}, true},
		{"log-fixed", ScaleLogSpec{Min: FloatVal(1)}, false},
		{"log-auto", ScaleLogSpec{Min: AutoFloat()}, true},
		{"symlog-auto", ScaleSymlogSpec{Min: AutoFloat()}, true},
		// Discrete specs carry no floor; report auto.
		{"band", ScaleBandSpec{}, true},
		{"point", ScalePointSpec{}, true},
	}
	for _, c := range cases {
		if got := SpecMinIsAuto(c.spec); got != c.want {
			t.Errorf("%s: SpecMinIsAuto = %v, want %v", c.name, got, c.want)
		}
	}
}
