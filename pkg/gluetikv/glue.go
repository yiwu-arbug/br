// Copyright 2020 PingCAP, Inc. Licensed under Apache-2.0.

package gluetikv

import (
	"context"

	pd "github.com/pingcap/pd/v4/client"
	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/domain"
	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tidb/store/tikv"

	"github.com/pingcap/br/pkg/glue"
	"github.com/pingcap/br/pkg/utils"
)

// Glue is an implementation of glue.Glue that accesses only TiKV without TiDB.
type Glue struct{}

// GetDomain implements glue.Glue
func (Glue) GetDomain(store kv.Storage) (*domain.Domain, error) {
	return nil, nil
}

// CreateSession implements glue.Glue
func (Glue) CreateSession(store kv.Storage) (glue.Session, error) {
	return nil, nil
}

// Open implements glue.Glue
func (Glue) Open(path string, option pd.SecurityOption) (kv.Storage, error) {
	if option.CAPath != "" {
		conf := config.GetGlobalConfig()
		conf.Security.ClusterSSLCA = option.CAPath
		conf.Security.ClusterSSLCert = option.CertPath
		conf.Security.ClusterSSLKey = option.KeyPath
		config.StoreGlobalConfig(conf)
	}
	return tikv.Driver{}.Open(path)
}

// OwnsStorage implements glue.Glue
func (Glue) OwnsStorage() bool {
	return true
}

// StartProgress implements glue.Glue
func (Glue) StartProgress(ctx context.Context, cmdName string, total int64, redirectLog bool) glue.Progress {
	return progress{ch: utils.StartProgress(ctx, cmdName, total, redirectLog)}
}

type progress struct {
	ch chan<- struct{}
}

// Inc implements glue.Progress
func (p progress) Inc() {
	p.ch <- struct{}{}
}

// Close implements glue.Progress
func (p progress) Close() {
	close(p.ch)
}
