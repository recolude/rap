package position_test

import (
	"testing"

	"github.com/EliCDavis/vector"
	"github.com/recolude/rap/pkg/encoding/position"
	"github.com/stretchr/testify/assert"
)

func Test_Vec3ToOctCells_TopRightForward(t *testing.T) {
	cells := make([]position.OctCell, 1)

	position.Vec3ToOctCells(
		vector.NewVector3(0.75, 0.75, 0.75),
		vector.Vector3Zero(),
		vector.Vector3One(),
		cells,
	)

	pos := position.OctCellsToVec3(
		vector.Vector3Zero(),
		vector.Vector3One(),
		cells,
	)

	assert.InDelta(t, 0.75, pos.X(), 0.001)
	assert.InDelta(t, 0.75, pos.Y(), 0.001)
	assert.InDelta(t, 0.75, pos.Z(), 0.001)
}

func Test_Vec3ToOctCells_TopLeftForward(t *testing.T) {
	cells := make([]position.OctCell, 1)

	start := vector.Vector3Zero()
	end := vector.Vector3One()

	tests := map[string]struct {
		x float64
		y float64
		z float64
	}{
		"a": {x: 0.75, y: 0.75, z: 0.75},
		"b": {x: 0.75, y: 0.75, z: 0.25},
		"c": {x: 0.75, y: 0.25, z: 0.75},
		"d": {x: 0.75, y: 0.25, z: 0.25},

		"e": {x: 0.25, y: 0.75, z: 0.75},
		"f": {x: 0.25, y: 0.75, z: 0.25},
		"g": {x: 0.25, y: 0.25, z: 0.75},
		"h": {x: 0.25, y: 0.25, z: 0.25},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			position.Vec3ToOctCells(vector.NewVector3(tc.x, tc.y, tc.z), start, end, cells)

			pos := position.OctCellsToVec3(start, end, cells)

			assert.InDelta(t, tc.x, pos.X(), 0.001)
			assert.InDelta(t, tc.y, pos.Y(), 0.001)
			assert.InDelta(t, tc.z, pos.Z(), 0.001)
		})
	}
}

func Test_Vec3ToOctCells_Weird(t *testing.T) {
	cells := make([]position.OctCell, 8)

	start := vector.NewVector3(2, 3, 4)
	end := vector.NewVector3(7, 8, 9)

	x := 4.
	y := 5.
	z := 6.

	position.Vec3ToOctCells(vector.NewVector3(x, y, z), start, end, cells)

	pos := position.OctCellsToVec3(start, end, cells)

	assert.InDelta(t, x, pos.X(), 0.01)
	assert.InDelta(t, y, pos.Y(), 0.01)
	assert.InDelta(t, z, pos.Z(), 0.01)
}
