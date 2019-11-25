package data_sqlite

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libraries/encodelib"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/components/data"
)

func TestCRUD(t *testing.T) {
	env := "test"
	err := os.Setenv("ENV", env)
	require.NoError(t, err)

	l, err = logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	configPath := filelib.CurrentPath() + "../../../environments/" + env + ".yaml"
	cfg, err := config.Get(configPath, encodelib.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	cfgSQLite := config.Access{}
	err = cfg.Value("sqlite", &cfgSQLite)
	require.NoError(t, err)

	l.Debugf("%#v", cfgSQLite)

	dataOp, cleanerOp, err := NewData(cfgSQLite, "", 0)
	require.NoError(t, err)

	l.Debugf("%#v", dataOp)

	testCases := data.TestCases(dataOp, cleanerOp)

	data.OperatorTestScenario(t, testCases, l)
}