package binary_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/recolude/rap/internal/io/binary"
	"github.com/recolude/rap/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_Reader(t *testing.T) {
	// ARRANGE ================================================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockReader(ctrl)
	m.EXPECT().Read(gomock.Any()).Return(1, nil)
	m.EXPECT().Read(gomock.Any()).Return(1, nil)
	m.EXPECT().Read(gomock.Any()).Return(1, nil)
	errReader := binary.NewErrReader(m)

	// ACT ====================================================================
	buf := []byte{0}
	r1, err1 := errReader.Read(buf)
	r2, err2 := errReader.Read(buf)
	r3, err3 := errReader.Read(buf)

	// ASSERT =================================================================
	assert.NoError(t, err1)
	assert.Equal(t, 1, r1)
	assert.NoError(t, err2)
	assert.Equal(t, 1, r2)
	assert.NoError(t, err3)
	assert.Equal(t, 1, r3)
	assert.Equal(t, 3, errReader.TotalRead())
	assert.NoError(t, errReader.Error())
}

func Test_Read_Byte(t *testing.T) {
	// ARRANGE ================================================================
	errReader := binary.NewErrReader(bytes.NewBuffer([]byte{44, 55, 66}))

	// ACT ====================================================================
	r1, err1 := errReader.ReadByte()
	r2, err2 := errReader.ReadByte()
	r3, err3 := errReader.ReadByte()

	// ASSERT =================================================================
	assert.NoError(t, err1)
	assert.Equal(t, byte(44), r1)
	assert.NoError(t, err2)
	assert.Equal(t, byte(55), r2)
	assert.NoError(t, err3)
	assert.Equal(t, byte(66), r3)
	assert.Equal(t, 3, errReader.TotalRead())
	assert.NoError(t, errReader.Error())
}

func Test_ReaderContinuesAfterError(t *testing.T) {
	// ARRANGE ================================================================
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockReader(ctrl)
	m.EXPECT().
		Read(gomock.Any()).
		Return(0, errors.New("test3")).
		Return(0, errors.New("test2")).
		Return(0, errors.New("test1"))

	errReader := binary.NewErrReader(m)

	// ACT ====================================================================
	buf := []byte{0}
	r1, err1 := errReader.Read(buf)
	r2, err2 := errReader.Read(buf)
	r3, err3 := errReader.Read(buf)

	// ASSERT =================================================================
	assert.EqualError(t, err1, "test1")
	assert.Equal(t, 0, r1)
	assert.EqualError(t, err2, "test1")
	assert.Equal(t, 0, r2)
	assert.EqualError(t, err3, "test1")
	assert.Equal(t, 0, r3)
	assert.Equal(t, 0, errReader.TotalRead())
	assert.EqualError(t, errReader.Error(), "test1")
}
