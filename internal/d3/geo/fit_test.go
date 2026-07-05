package geo

import "testing"

// Projection fitting tests. Expected scale/translate values were produced by
// d3-geo@3's projection.fitExtent/fitSize/fitWidth/fitHeight.

// cwWorld is a clockwise 60°×40° box (small, unambiguous interior).
func cwWorld() Geometry {
	return Geometry{Type: TypePolygon, Coordinates: [][][2]float64{{
		{-30, -20}, {-30, 20}, {30, 20}, {30, -20}, {-30, -20},
	}}}
}

func TestFitExtentAndSize(t *testing.T) {
	fe := GeoEquirectangular().FitExtent(0, 0, 300, 200, cwWorld())
	approx(t, "fitExtent.scale", fe.ScaleValue(), 251.3427)
	fx, fy := fe.TranslateValue()
	approx(t, "fitExtent.tx", fx, 150)
	approx(t, "fitExtent.ty", fy, 100)

	fs := GeoEquirectangular().FitSize(400, 300, cwWorld())
	approx(t, "fitSize.scale", fs.ScaleValue(), 377.0141)
	sx, sy := fs.TranslateValue()
	approx(t, "fitSize.tx", sx, 200)
	approx(t, "fitSize.ty", sy, 150)

	fw := GeoEquirectangular().FitWidth(400, cwWorld())
	approx(t, "fitWidth.scale", fw.ScaleValue(), 381.9719)
	wx, wy := fw.TranslateValue()
	approx(t, "fitWidth.tx", wx, 200)
	approx(t, "fitWidth.ty", wy, 151.9725)

	fh := GeoEquirectangular().FitHeight(300, cwWorld())
	approx(t, "fitHeight.scale", fh.ScaleValue(), 377.0141)
	hx, hy := fh.TranslateValue()
	approx(t, "fitHeight.tx", hx, 197.4041)
	approx(t, "fitHeight.ty", hy, 150)
}

// TestFit_PreservesClipExtent verifies fit restores a rectangular postclip that
// was installed before fitting (the fitted world still renders non-empty).
func TestFit_PreservesClipExtent(t *testing.T) {
	p := GeoEquirectangular().ClipExtent(10, 20, 300, 400)
	p.FitSize(400, 300, cwWorld())
	if NewPath(p).Geometry(cwWorld()) == "" {
		t.Error("fitted world should render non-empty")
	}
}
