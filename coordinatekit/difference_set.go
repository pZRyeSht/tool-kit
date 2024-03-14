package coordinatekit

// 坐标范围由左下角到右上角组成的规则矩形，大坐标范围包含小坐标范围，得到除开小坐标范围的剩余坐标坐标范围的规则矩形坐标范围（即两个大小坐标范围的差集）

// 计算两个坐标范围的差集
func difference(a, b Coordinate) []Coordinate {
	result := make([]Coordinate, 0)
	
	// 计算左侧区域
	left := Coordinate{
		BottomLeftX: b.BottomLeftX,
		BottomLeftY: b.BottomLeftY,
		TopRightX:   a.BottomLeftX,
		TopRightY:   b.TopRightY,
	}
	if left.BottomLeftX < left.TopRightX {
		result = append(result, left)
	}
	
	// 计算右侧区域
	right := Coordinate{
		BottomLeftX: a.TopRightX,
		BottomLeftY: b.BottomLeftY,
		TopRightX:   b.TopRightX,
		TopRightY:   b.TopRightY,
	}
	if right.BottomLeftX < right.TopRightX {
		result = append(result, right)
	}
	
	// 计算上侧区域
	upper := Coordinate{
		BottomLeftX: a.BottomLeftX,
		BottomLeftY: a.TopRightY,
		TopRightX:   a.TopRightX,
		TopRightY:   b.TopRightY,
	}
	if upper.BottomLeftY < upper.TopRightY {
		result = append(result, upper)
	}
	
	// 计算下侧区域
	lower := Coordinate{
		BottomLeftX: a.BottomLeftX,
		BottomLeftY: b.BottomLeftY,
		TopRightX:   a.TopRightX,
		TopRightY:   a.BottomLeftY,
	}
	if lower.BottomLeftY < lower.TopRightY {
		result = append(result, lower)
	}
	
	return result
}
