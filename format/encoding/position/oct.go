package position

import "github.com/EliCDavis/vector"

type OctCell int

// ORDER MATTERS: topBit << 2 | rightBit << 1 | forwardBit
const (
	TopRightForward OctCell = iota
	TopRightBackward
	TopLeftForward
	TopLeftBackward
	BottomRightForward
	BottomRightBackward
	BottomLeftForward
	BottomLeftBackward
)

func Vec3ToOctCells(v, min, max vector.Vector3, cells []OctCell) {
	center := min.Add(max).DivByConstant(2.0)
	crossSection := max.Sub(min)
	incrementX := crossSection.X() / 4.0
	incrementY := crossSection.Y() / 4.0
	incrementZ := crossSection.Z() / 4.0

	for cellIndex := 0; cellIndex < len(cells); cellIndex++ {
		topBit := 0
		newY := center.Y() + incrementY
		if v.Y() < center.Y() {
			topBit = 1
			newY = center.Y() - incrementY
		}

		rightBit := 0
		newX := center.X() + incrementX
		if v.X() < center.X() {
			rightBit = 1
			newX = center.X() - incrementX
		}

		forwardBit := 0
		newZ := center.Z() + incrementZ
		if v.Z() < center.Z() {
			forwardBit = 1
			newZ = center.Z() - incrementZ
		}
		cells[cellIndex] = OctCell(topBit<<2 | rightBit<<1 | forwardBit)
		center = vector.NewVector3(newX, newY, newZ)
		incrementX /= 2.0
		incrementY /= 2.0
		incrementZ /= 2.0
	}
}

func OctCellsToVec3(min, max vector.Vector3, cells []OctCell) vector.Vector3 {
	center := min.Add(max).DivByConstant(2.0)
	crossSection := max.Sub(min)
	incrementX := crossSection.X() / 4.0
	incrementY := crossSection.Y() / 4.0
	incrementZ := crossSection.Z() / 4.0

	for cellIndex := 0; cellIndex < len(cells); cellIndex++ {
		newY := center.Y() - incrementY
		if cells[cellIndex]&0b100 == 0 {
			newY = center.Y() + incrementY
		}

		newX := center.X() - incrementX
		if cells[cellIndex]&0b010 == 0 {
			newX = center.X() + incrementX
		}

		newZ := center.Z() - incrementZ
		if cells[cellIndex]&0b001 == 0 {
			newZ = center.Z() + incrementZ
		}
		center = vector.NewVector3(newX, newY, newZ)
		incrementX /= 2.0
		incrementY /= 2.0
		incrementZ /= 2.0
	}
	return center
}
