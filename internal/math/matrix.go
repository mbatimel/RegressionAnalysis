package math

import (
	"errors"
	"fmt"
	"github.com/mbatimel/RegressionAnalysis/internal/models"
	"math"
)

type Matrix interface {
	Sum_matrix(matrix_1, matrix_2 *models.Matrix) (err error)
	Sub_matrix(matrix_1, matrix_2 *models.Matrix) (err error)
	Mul_matrix(matrix_1, matrix_2 *models.Matrix) (err error)
	MulNumber_matrix(m *models.Matrix, scalar float64) *models.Matrix
	Div_matrix(m *models.Matrix, scalar float64) (*models.Matrix, error)
	Transpose_matrix(m *models.Matrix) *models.Matrix
	Inverse_matrix(m *models.Matrix) (*models.Matrix, error)
	Exponential_matrix(m *models.Matrix) *models.Matrix
}

// CreateMatrix creates a new models.Matrix with the given values.
func CreateMatrix(values [][]float64) (*models.Matrix, error) {
	rows := int64(len(values))
	if rows == 0 {
		return nil, errors.New("models.Matrix cannot have zero rows")
	}
	cols := int64(len(values[0]))
	for _, row := range values {
		if int64(len(row)) != cols {
			return nil, errors.New("all rows must have the same number of columns")
		}
	}
	return &models.Matrix{Values: values, Rows: rows, Cols: cols}, nil
}

// Mul_matrix multiplies two matrices.
func Mul_matrix(matrix1, matrix2 *models.Matrix) (*models.Matrix, error) {
	if matrix1.Cols != matrix2.Rows {
		return nil, errors.New("number of columns in matrix1 must equal the number of rows in matrix2")
	}

	result := make([][]float64, matrix1.Rows)
	for i := range result {
		result[i] = make([]float64, matrix2.Cols)
	}

	for i := int64(0); i < matrix1.Rows; i++ {
		for j := int64(0); j < matrix2.Cols; j++ {
			for k := int64(0); k < matrix1.Cols; k++ {
				result[i][j] += matrix1.Values[i][k] * matrix2.Values[k][j]
			}
		}
	}

	return CreateMatrix(result)
}

// MulNumber_matrix multiplies a models.Matrix by a scalar.
func MulNumber_matrix(m *models.Matrix, scalar float64) *models.Matrix {
	result := make([][]float64, m.Rows)
	for i := range result {
		result[i] = make([]float64, m.Cols)
		for j := range result[i] {
			result[i][j] = m.Values[i][j] * scalar
		}
	}
	return &models.Matrix{Values: result, Rows: m.Rows, Cols: m.Cols}
}

// Div_matrix divides a models.Matrix by a scalar.
func Div_matrix(m *models.Matrix, scalar float64) (*models.Matrix, error) {
	if scalar == 0 {
		return nil, errors.New("division by zero is not allowed")
	}
	return MulNumber_matrix(m, 1/scalar), nil
}

// Transpose_matrix returns the transpose of a models.Matrix.
func Transpose_matrix(m *models.Matrix) *models.Matrix {
	result := make([][]float64, m.Cols)
	for i := range result {
		result[i] = make([]float64, m.Rows)
		for j := range result[i] {
			result[i][j] = m.Values[j][i]
		}
	}
	return &models.Matrix{Values: result, Rows: m.Cols, Cols: m.Rows}
}

// Inverse_matrix calculates the inverse of a 2x2 models.Matrix.
func Inverse_matrix(m *models.Matrix) (*models.Matrix, error) {
	if m.Rows != 2 || m.Cols != 2 {
		return nil, errors.New("only 2x2 matrices are supported for inversion")
	}

	a, b := m.Values[0][0], m.Values[0][1]
	c, d := m.Values[1][0], m.Values[1][1]
	det := a*d - b*c

	if det == 0 {
		return nil, errors.New("models.Matrix is singular and cannot be inverted")
	}

	inverseValues := [][]float64{
		{d / det, -b / det},
		{-c / det, a / det},
	}
	return CreateMatrix(inverseValues)
}

// Exponential_matrix raises each element of the models.Matrix to the power of e.
func Exponential_matrix(m *models.Matrix) *models.Matrix {
	result := make([][]float64, m.Rows)
	for i := range result {
		result[i] = make([]float64, m.Cols)
		for j := range result[i] {
			result[i][j] = math.Exp(m.Values[i][j])
		}
	}
	return &models.Matrix{Values: result, Rows: m.Rows, Cols: m.Cols}
}

func Sum_matrix(matrix_1, matrix_2 *models.Matrix) (err error) {
	if matrix_1 == nil || matrix_2 == nil {
		return fmt.Errorf("someone models.Matrix is nil")
	}

	if matrix_1.Cols != matrix_2.Cols {
		return fmt.Errorf("The number of columns does not match")
	}
	if matrix_1.Rows != matrix_2.Rows {
		return fmt.Errorf("The number of rows does not match")
	}
	for i := 0; i < int(matrix_1.Rows); i++ {
		for j := 0; j < int(matrix_2.Rows); j++ {
			matrix_1.Values[i][j] += matrix_2.Values[i][j]
		}
	}
	return nil
}

func Sub_matrix(matrix_1, matrix_2 *models.Matrix) (err error) {
	if matrix_1 == nil || matrix_2 == nil {
		return fmt.Errorf("someone models.Matrix is nil")
	}
	if len(matrix_1.Values) != len(matrix_2.Values) {
		if matrix_1.Cols != matrix_2.Cols {
			return fmt.Errorf("The number of columns does not match")
		}
		if matrix_1.Rows != matrix_2.Rows {
			return fmt.Errorf("The number of rows does not match")
		}
	}
	for i := 0; i < int(matrix_1.Rows); i++ {
		for j := 0; j < int(matrix_2.Rows); j++ {
			matrix_1.Values[i][j] += matrix_2.Values[i][j]
		}
	}
	return nil
}
