package coordinatekit

// Coordinate 坐标结构体，包含左下角和右上角坐标
type Coordinate struct {
	BottomLeftX, BottomLeftY int
	TopRightX, TopRightY     int
}

// isOverlap 范围重叠
func isOverlap(rect1, rect2 Coordinate) bool {
	// Check if rect1 is to the right of rect2
	if rect1.BottomLeftX > rect2.TopRightX {
		return false
	}
	// Check if rect1 is to the left of rect2
	if rect1.TopRightX < rect2.BottomLeftX {
		return false
	}
	// Check if rect1 is above rect2
	if rect1.BottomLeftY > rect2.TopRightY {
		return false
	}
	// Check if rect1 is below rect2
	if rect1.TopRightY < rect2.BottomLeftY {
		return false
	}
	// If none of the above conditions are met, the rectangles must overlap
	return true
}

// isRangeContained 坐标覆盖
func isRangeContained(range1, range2 Coordinate) bool {
	if range1.BottomLeftX <= range2.BottomLeftX && range1.BottomLeftY <= range2.BottomLeftY &&
		range1.TopRightX >= range2.TopRightX && range1.TopRightY >= range2.TopRightY {
		return true
	}
	
	if range2.BottomLeftX <= range1.BottomLeftX && range2.BottomLeftY <= range1.BottomLeftY &&
		range2.TopRightX >= range1.TopRightX && range2.TopRightY >= range1.TopRightY {
		return true
	}
	return false
}
