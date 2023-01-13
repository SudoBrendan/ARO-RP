package frontend

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/ugorji/go/codec"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/database/cosmosdb"
	"github.com/Azure/ARO-RP/pkg/frontend/middleware"
)

func (f *frontend) getAdminHiveClusterDeployment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ctx.Value(middleware.ContextKeyLog).(*logrus.Entry)

	b, err := f._getAdminHiveClusterDeployment(ctx, r, log)

	switch {
	case cosmosdb.IsErrorStatusCode(err, http.StatusNotFound):
		api.WriteError(w, http.StatusNotFound, api.CloudErrorCodeNotFound, "", "Cluster not found.")
		return
	case err != nil:
		api.WriteError(w, http.StatusInternalServerError, api.CloudErrorCodeInternalServerError, "", err.Error())
		return
	}

	adminReply(log, w, nil, b, err)
}

func (f *frontend) _getAdminHiveClusterDeployment(ctx context.Context, r *http.Request, log *logrus.Entry) ([]byte, error) {
	url := filepath.Dir(r.URL.Path)
	resourceID := strings.TrimPrefix(url, "/admin")
	doc, err := f.dbOpenShiftClusters.Get(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if doc.OpenShiftCluster.Properties.HiveProfile.Namespace == "" {
		return nil, errors.New("cluster is not managed by hive")
	}

	hr, err := f.hiveActionsFactory(log, f.env)
	if err != nil {
		return nil, err
	}

	cd, err := hr.GetClusterDeployment(ctx, doc)
	if err != nil {
		return nil, err
	}

	var b []byte
	err = codec.NewEncoderBytes(&b, &codec.JsonHandle{}).Encode(cd)
	if err != nil {
		return nil, err
	}

	return b, nil
}
