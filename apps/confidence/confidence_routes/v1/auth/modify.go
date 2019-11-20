package v1_auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/apps/confidence/confidence_routes"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/components/auth"
	"github.com/pavlo67/workshop/libraries/filelib"
)

var _ = server_http.InitEndpoint(&confidence_routes.Endpoints, "POST", filelib.RelativePath(filelib.CurrentFile(true), confidence_routes.BasePath, confidence_routes.Prefix),
	nil, workerModify, "")
var _ server_http.WorkerHTTP = workerModify

func workerModify(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	if user == nil {
		return server.ResponseRESTError(http.StatusForbidden, errors.New("no user authorized"))
	}

	credsJSON, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Wrap(err, "can't read body"))
	}

	var toReplace auth.Creds
	err = json.Unmarshal(credsJSON, &toReplace)
	if err != nil {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Wrapf(err, "can't unmarshal body: %s", credsJSON))
	}

	// toReplace = append(user.Creds, toReplace...)

	// !!! previous user.Creds are ignored here
	toReplaceModified, err := confidence_routes.AuthOpToSetToken.SetCreds(*user, toReplace) // TODO!!! add custom toReplace
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Wrap(err, "can't create JWT"))
	}

	if toReplaceModified != nil {
		for t, c := range toReplaceModified.Values {
			user.Creds.Values[t] = c
		}
	}
	return server.ResponseRESTOk(map[string]interface{}{"user": user})
}
