package adot

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

var tmp = "tmp"
var homeA = tmp + "/home"
var homeB = tmp + "/home"
var remote = tmp + "/remote"
var fileZ = "testing"

func TestAll(t *testing.T) {
	require.NoError(t, os.RemoveAll(homeA))
	require.NoError(t, os.Mkdir(homeA, 700))
	require.NoError(t, os.RemoveAll(homeB))
	require.NoError(t, os.Mkdir(homeB, 700))
	require.NoError(t, os.RemoveAll(remote))
	require.NoError(t, os.Mkdir(remote, 700))

	a := ADot{}
	require.NoError(t, os.Chdir(homeA))
	newFile(t, fileZ, "123\n466\n789")
	require.NoError(t, a.InitNew(remote))
	require.NoError(t, a.Add(fileZ))

	b := ADot{}
	require.NoError(t, os.Chdir(homeB))
	require.NoError(t, b.InitExisting(remote))
	checkFile(t, "testing", "123\n466\n789")
}

func newFile(t *testing.T, path, content string) {
	t.Helper()

	fp, err := os.Create(path)
	require.NoError(t, err)
	_, err = fp.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, fp.Close())
}

func checkFile(t *testing.T, path, expected string) {
	t.Helper()

	fp, err := os.Open(path)
	require.NoError(t, err)

	var content []byte
	content, err = ioutil.ReadAll(fp)
	require.NoError(t, err)
	require.NoError(t, fp.Close())

	require.Equal(t, expected, string(content))
}
