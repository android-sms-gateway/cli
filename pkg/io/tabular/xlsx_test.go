package tabular_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/android-sms-gateway/cli/pkg/io/tabular"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func TestXLSXReader_Read_WithHeader(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "input.xlsx")

	f := excelize.NewFile()
	defer f.Close()
	index, err := f.NewSheet("Data")
	require.NoError(t, err)
	f.SetActiveSheet(index)

	require.NoError(t, f.SetCellValue("Data", "A1", "Phone"))
	require.NoError(t, f.SetCellValue("Data", "B1", "Message"))
	require.NoError(t, f.SetCellValue("Data", "A2", "+12025550123"))
	require.NoError(t, f.SetCellValue("Data", "B2", "Hello"))
	require.NoError(t, f.SaveAs(path))

	reader := tabular.NewXLSXReader(tabular.XLSXConfig{Path: path, Sheet: "Data", HasHeader: true})
	records, err := reader.Read(context.Background())
	require.NoError(t, err)
	require.Len(t, records, 1)

	assert.Equal(t, 2, records[0].RowNumber)
	assert.Equal(t, "+12025550123", records[0].Values["Phone"])
	assert.Equal(t, "Hello", records[0].Values["Message"])
}

func TestXLSXReader_Read_WithoutHeader(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "input.xlsx")

	f := excelize.NewFile()
	defer f.Close()
	index, err := f.NewSheet("Data")
	require.NoError(t, err)
	f.SetActiveSheet(index)

	require.NoError(t, f.SetCellValue("Data", "A1", "+12025550123"))
	require.NoError(t, f.SetCellValue("Data", "B1", "Hello"))
	require.NoError(t, f.SaveAs(path))

	reader := tabular.NewXLSXReader(tabular.XLSXConfig{Path: path, Sheet: "Data", HasHeader: false})
	records, err := reader.Read(context.Background())
	require.NoError(t, err)
	require.Len(t, records, 1)

	assert.Equal(t, 1, records[0].RowNumber)
	assert.Equal(t, "+12025550123", records[0].Values["col_1"])
	assert.Equal(t, "Hello", records[0].Values["col_2"])
}
