package tabular_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/android-sms-gateway/cli/pkg/io/tabular"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSVReader_Read_WithHeader(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "input.csv")
	require.NoError(t, os.WriteFile(path, []byte("Phone,Message\n+12025550123, Hello "), 0o600))

	reader := tabular.NewCSVReader(tabular.CSVConfig{Path: path, HasHeader: true})
	records, err := reader.Read(context.Background())
	require.NoError(t, err)
	require.Len(t, records, 1)

	assert.Equal(t, 2, records[0].RowNumber)
	assert.Equal(t, "+12025550123", records[0].Values["Phone"])
	assert.Equal(t, "Hello", records[0].Values["Message"])
}

func TestCSVReader_Read_WithoutHeader(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "input.csv")
	require.NoError(t, os.WriteFile(path, []byte("+12025550123,Hello"), 0o600))

	reader := tabular.NewCSVReader(tabular.CSVConfig{Path: path, HasHeader: false})
	records, err := reader.Read(context.Background())
	require.NoError(t, err)
	require.Len(t, records, 1)

	assert.Equal(t, "+12025550123", records[0].Values["col_1"])
	assert.Equal(t, "Hello", records[0].Values["col_2"])
}
