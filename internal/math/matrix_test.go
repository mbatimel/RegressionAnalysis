package math

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SumMatrix(t *testing.T) {
	t.Run("ok, sum matrix", func(t *testing.T) {
		value := [][]float64{
			{1, 2, 3},
			{1, 2, 3},
			{1, 2, 3}}
		expected := [][]float64{
			{2, 4, 6},
			{2, 4, 6},
			{2, 4, 6}}
		matrix1, err := CreateMatrix(value)
		assert.NoError(t, err)
		matrix2, err := CreateMatrix(value)
		assert.NoError(t, err)
		err = Sum_matrix(matrix1, matrix2)
		assert.NoError(t, err)
		assert.Equal(t, matrix1.Values, expected)

	})
	t.Run("err, error in column", func(t *testing.T) {
		value1 := [][]float64{
			{1, 2},
			{1, 2},
			{1, 2}}
		value2 := [][]float64{
			{2, 4, 6},
			{2, 4, 6},
			{2, 4, 6}}
		matrix1, err := CreateMatrix(value1)
		assert.NoError(t, err)
		matrix2, err := CreateMatrix(value2)
		assert.NoError(t, err)
		err = Sum_matrix(matrix1, matrix2)
		assert.Error(t, err, errors.New("The number of columns does not match"))
	})
	t.Run("err, error in rows", func(t *testing.T) {
		value1 := [][]float64{
			{1, 2, 1},
			{1, 2, 1}}
		value2 := [][]float64{
			{2, 4, 6},
			{2, 4, 6},
			{2, 4, 6}}
		expextErr :=errors.New("The number of rows does not match")
		matrix1, err := CreateMatrix(value1)
		assert.NoError(t, err)
		matrix2, err := CreateMatrix(value2)
		assert.NoError(t, err)
		err = Sum_matrix(matrix1, matrix2)
		assert.Error(t, expextErr, err)
	})
}
