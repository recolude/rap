package euler

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/EliCDavis/vector"
	"github.com/recolude/rap/format/streams/euler"
)

type matrix [][]float64

func Multiply3x3(m1, m2 matrix) matrix {
	m3 := [][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			m3[i][j] = 0
			for k := 0; k < 3; k++ {
				m3[i][j] = m3[i][j] + (m1[i][k] * m2[k][j])
			}
		}
	}
	return m3
}

func Multiply3x3by3x1(m matrix, v vector.Vector3) vector.Vector3 {
	return vector.NewVector3(
		(m[0][0]*v.X())+(m[0][1]*v.Y())+(m[0][2]*v.Z()),
		(m[1][0]*v.X())+(m[1][1]*v.Y())+(m[1][2]*v.Z()),
		(m[2][0]*v.X())+(m[2][1]*v.Y())+(m[2][2]*v.Z()),
	)
}

func Transpose(in matrix) matrix {
	return [][]float64{
		{in[0][0], in[1][0], in[2][0]},
		{in[0][1], in[1][1], in[2][1]},
		{in[0][2], in[1][2], in[2][2]},
	}
}

func Rx(theta float64) matrix {
	return [][]float64{
		{1, 0, 0},
		{0, math.Cos(theta), -math.Sin(theta)},
		{0, math.Sin(theta), math.Cos(theta)},
	}
}

func RxT(theta float64) matrix {
	return Transpose(Rx(theta))
}

func Ry(theta float64) matrix {
	return [][]float64{
		{math.Cos(theta), 0, math.Sin(theta)},
		{0, 1, 0},
		{-math.Sin(theta), 0, math.Cos(theta)},
	}
}

func RyT(theta float64) matrix {
	return Transpose(Ry(theta))
}

func Rz(theta float64) matrix {
	return [][]float64{
		{math.Cos(theta), -math.Sin(theta), 0},
		{math.Sin(theta), math.Cos(theta), 0},
		{0, 0, 1},
	}
}

func RzT(theta float64) matrix {
	return Transpose(Rz(theta))
}

func ToNormal(inEulerAngle vector.Vector3) vector.Vector3 {
	vectorToTransform := vector.Vector3Forward()

	Sx := math.Sin(inEulerAngle.X() * (math.Pi / 180.))
	Sy := math.Sin(inEulerAngle.Y() * (math.Pi / 180.))
	Sz := math.Sin(inEulerAngle.Z() * (math.Pi / 180.))
	Cx := math.Cos(inEulerAngle.X() * (math.Pi / 180.))
	Cy := math.Cos(inEulerAngle.Y() * (math.Pi / 180.))
	Cz := math.Cos(inEulerAngle.Z() * (math.Pi / 180.))

	var Mx matrix = make([][]float64, 3)
	Mx[0] = []float64{0, 0, 0}
	Mx[1] = []float64{0, 0, 0}
	Mx[2] = []float64{0, 0, 0}

	Mx[0][0] = Cy*Cz - Sx*Sy*Sz
	Mx[0][1] = -Cx * Sz
	Mx[0][2] = Cz*Sy + Cy*Sx*Sz
	Mx[1][0] = Cz*Sx*Sy + Cy*Sz
	Mx[1][1] = Cx * Cz
	Mx[1][2] = -Cy*Cz*Sx + Sy*Sz
	Mx[2][0] = -Cx * Sy
	Mx[2][1] = Sx
	Mx[2][2] = Cx * Cy

	return vector.NewVector3(
		(Mx[0][0]*vectorToTransform.X())+(Mx[0][1]*vectorToTransform.Y())+(Mx[0][2]*vectorToTransform.Z()),
		(Mx[1][0]*vectorToTransform.X())+(Mx[1][1]*vectorToTransform.Y())+(Mx[1][2]*vectorToTransform.Z()),
		(Mx[2][0]*vectorToTransform.X())+(Mx[2][1]*vectorToTransform.Y())+(Mx[2][2]*vectorToTransform.Z()),
	)
}

func encodeOct24(captures []euler.Capture) []byte {
	// https://stackoverflow.com/questions/1568568/how-to-convert-euler-angles-to-directional-vector

	streamData := new(bytes.Buffer)

	buf := make([]byte, 8)
	size := binary.PutUvarint(buf, uint64(len(captures)))
	streamData.Write(buf[:size])

	for _, capture := range captures {
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.Time()))
		streamData.Write(buf)
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.EulerZXY().X()))
		streamData.Write(buf)
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.EulerZXY().Y()))
		streamData.Write(buf)
		binary.LittleEndian.PutUint64(buf, math.Float64bits(capture.EulerZXY().Z()))
		streamData.Write(buf)
	}
	return streamData.Bytes()
}

func decodeOct24(streamData *bytes.Reader) ([]euler.Capture, error) {
	numCaptures, err := binary.ReadUvarint(streamData)
	if err != nil {
		return nil, err
	}

	captures := make([]euler.Capture, numCaptures)
	for i := 0; i < int(numCaptures); i++ {
		var time float32
		var x float32
		var y float32
		var z float32

		binary.Read(streamData, binary.LittleEndian, &time)
		binary.Read(streamData, binary.LittleEndian, &x)
		binary.Read(streamData, binary.LittleEndian, &y)
		binary.Read(streamData, binary.LittleEndian, &z)
		captures[i] = euler.NewEulerZXYCapture(float64(time), float64(x), float64(y), float64(z))
	}

	return captures, nil
}
