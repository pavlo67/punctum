package apps

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/serializer"
)

func Prepare(buildDate, buildTag, buildCommit, serviceName, appsSubpathDefault string) (versionOnly bool, envPath string, cfgService *config.Config, l logger.Operator) {

	rand.Seed(time.Now().UnixNano())

	var appsSubpath string
	flag.BoolVar(&versionOnly, "v", false, "show build vars only")
	flag.StringVar(&appsSubpath, "apps_subpath", appsSubpathDefault, "subpath to /apps directory")
	flag.Parse()

	log.Printf("builded: %s, tag: %s, commit: %s\n", buildDate, buildTag, buildCommit)

	if versionOnly {
		return versionOnly, "", nil, nil
	}

	// logger

	l, err := logger.Init(logger.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// getting config environments

	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}

	cwd, err := os.Getwd()
	if err != nil {
		l.Fatal("can't os.Getwd(): ", err)
	}
	cwd += "/"
	l.Info("CWD: ", cwd)

	// get config

	envPath = cwd + appsSubpath + "_environments/"
	cfgServicePath := envPath + configEnv + ".yaml"
	cfgService, err = config.Get(cfgServicePath, serviceName, serializer.MarshalerYAML)
	if err != nil || cfgService == nil {
		l.Fatalf("on config.Get(%s, %s, serializer.MarshalerYAML)", cfgServicePath, serviceName, cfgService, err)
	}
	return versionOnly, envPath, cfgService, l
}

func PrepareTests(t *testing.T, serviceName, appsSubpath, configEnv string) (envPath string, cfgService *config.Config) {
	os.Setenv("ENV", configEnv)

	l, err := logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	cwd, err := os.Getwd()
	require.NoError(t, err)

	cwd += "/"
	// t.Log("CWD: ", cwd)

	// get config

	envPath = cwd + appsSubpath + "_environments/"
	cfgServicePath := envPath + configEnv + ".yaml"
	cfgService, err = config.Get(cfgServicePath, serviceName, serializer.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfgService)

	return envPath, cfgService

}