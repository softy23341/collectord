package dalpg

import (
	"bytes"
	"fmt"

	"git.softndit.com/collector/backend/dal"
	"github.com/jackc/pgx"
	"gopkg.in/inconshreveable/log15.v2"
)

// Querier TBD
type Querier interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}

// Transactionable TBD
type Transactionable interface {
	Begin() (*pgx.Tx, error)
}

// Manager TBD
type Manager struct {
	p   Querier
	ptx Transactionable
	log log15.Logger
}

// ManagerContext TBD
type ManagerContext struct {
	Log        log15.Logger
	PoolConfig *pgx.ConnPoolConfig
}

// NewManager TBD
func NewManager(mc *ManagerContext) (*Manager, error) {
	pool, err := pgx.NewConnPool(*mc.PoolConfig)
	if err != nil {
		return nil, err
	}

	return &Manager{p: pool, ptx: pool, log: mc.Log}, nil
}

// BeginTx TBD
func (m *Manager) BeginTx() (dal.TxManager, error) {
	tx, err := m.ptx.Begin()
	if err != nil {
		return nil, err
	}

	return &TxManager{
		p:       tx,
		log:     m.log,
		Manager: &Manager{p: tx, ptx: m.ptx, log: m.log},
	}, nil
}

var _ dal.TrManager = (*Manager)(nil)

// TxManager TBD
type TxManager struct {
	*Manager
	p   *pgx.Tx
	log log15.Logger
}

// Commit TBD
func (tm *TxManager) Commit() error {
	return tm.p.Commit()
}

// Rollback TBD
func (tm *TxManager) Rollback() error {
	return tm.p.Rollback()
}

var _ dal.TxManager = (*TxManager)(nil)

// MakePlaceholders TBD
func MakePlaceholders(n, begin int) string {
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprint("$", i+begin))
	}
	return buf.String()
}
