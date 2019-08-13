package content_files

import (
	"os"

	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/components/basis"
	"github.com/pavlo67/constructor/components/basis/config"
	"github.com/pavlo67/constructor/components/basis/joiner"
	"github.com/pavlo67/constructor/components/basis/logger"
	"github.com/pavlo67/constructor/components/basis/starter"
	"github.com/pavlo67/constructor/components/structura"
)

var l logger.Operator
var _ starter.Operator = &contentFilesStarter{}

func Starter() starter.Operator {
	return &contentFilesStarter{}
}

type contentFilesStarter struct {
	interfaceKey joiner.InterfaceKey
	path         string
	marshaler    basis.Marshaler
}

func (nms *contentFilesStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (nms *contentFilesStarter) Init(conf *config.Config, options basis.Info) (info []basis.Info, err error) {
	l = logger.Get()

	var ok bool
	nms.path, ok = options.String("path")
	if !ok {
		return nil, errors.New("no path for contentFilesStarter.Init()")
	}
	fileinfo, err := os.Stat(nms.path)
	if err != nil {
		return nil, errors.Wrapf(err, "on check directory '%s'", nms.path)
	}

	if !fileinfo.IsDir() {
		return nil, errors.Errorf("'%s' isn't a directory", nms.path)
	}

	nms.marshaler, ok = options["marshaler"].(basis.Marshaler)
	if !ok || nms.marshaler == nil {
		return nil, errors.New("no marshaler for contentFilesStarter.Init()")
	}

	nms.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(structura.InterfaceKey)))

	return nil, nil
}

func (nms *contentFilesStarter) Setup() error {
	err := os.MkdirAll(nms.path, 0755)
	if err != nil {
		return errors.Wrapf(err, "on create directory '%s'", nms.path)
	}

	return nil
}

func (nms *contentFilesStarter) Run(joiner joiner.Operator) error {

	//var err error
	//notesOp, err := New(
	//	nms.mysqlConfig,
	//	nms.index.Params["table"],
	//	nms.jointLinks,
	//	nil, // grpsOp,
	//	nil, // linksOp,
	//	nil,
	//)
	//if err != nil {
	//	return errors.Wrap(err, "can't init notes_mysql")
	//}
	//
	//err = joiner.Join(notesOp, nms.interfaceKey)
	//if err != nil {
	//	return errors.Wrap(err, "can't join notes_mysql as notes.Operator interface")
	//}

	//err = joiner.Join(dataOp.Clean, ds.cleanerInterfaceKey)
	//if err != nil {
	//	return errors.Wrapf(err, "can't join datastoremysql.Operator.Clean as %s", ds.cleanerInterfaceKey)
	//}

	return nil
}