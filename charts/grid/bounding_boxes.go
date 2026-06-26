package grid

// AreBoundingBoxTouching reports whether two axis-aligned bounding boxes
// touch or overlap. Mirrors @nivo/grid areBoundingBoxTouching.
func AreBoundingBoxTouching(boxA, boxB BoundingBox) bool {
	touchX := boxA.Left <= boxB.Right && boxB.Left <= boxA.Right
	touchY := boxA.Top <= boxB.Bottom && boxB.Top <= boxA.Bottom
	return touchX && touchY
}
