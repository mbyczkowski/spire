package sql

import (
	"time"
)

// Model is used as a base for other models. Similar to gorm.Model without `DeletedAt`.
// We don't want soft-delete support.
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Bundle holds a trust bundle.
type Bundle struct {
	Model

	TrustDomain string `gorm:"not null;unique_index"`
	Data        []byte `gorm:"size:65535"`

	FederatedEntries []RegisteredEntry `gorm:"many2many:federated_registration_entries;"`
}

// MysqlBundle holds a trust bundle. It is equivalent to Bundle, but contains
// MySQL-specific gorm tags.
type MysqlBundle struct {
	Model

	TrustDomain string `gorm:"size:191;not null;unique_index"`
	Data        []byte `gorm:"size:65535"`

	FederatedEntries []RegisteredEntry `gorm:"many2many:federated_registration_entries;jointable_foreignkey:bundle_id"`
}

var _ = MysqlBundle(Bundle{})

// TableName gets table name of MysqlJoinToken
func (MysqlBundle) TableName() string {
	return "bundles"
}

// AttestedNode holds an attested node (agent)
type AttestedNode struct {
	Model

	SpiffeID     string `gorm:"unique_index"`
	DataType     string
	SerialNumber string
	ExpiresAt    time.Time
}

// TableName gets table name of AttestedNode
func (AttestedNode) TableName() string {
	return "attested_node_entries"
}

// MysqlAttestedNode holds an attested node (agent). It is equivalent to Bundle,
// but contains MySQL-specific gorm tags.
type MysqlAttestedNode struct {
	Model

	SpiffeID     string `gorm:"size:191;unique_index"`
	DataType     string
	SerialNumber string
	ExpiresAt    time.Time
}

var _ = MysqlAttestedNode(AttestedNode{})

// TableName gets table name of MysqlAttestedNode
func (MysqlAttestedNode) TableName() string {
	return "attested_node_entries"
}

// NodeSelector holds a node selector by spiffe ID
type NodeSelector struct {
	Model

	SpiffeID string `gorm:"unique_index:idx_node_resolver_map"`
	Type     string `gorm:"unique_index:idx_node_resolver_map"`
	Value    string `gorm:"unique_index:idx_node_resolver_map"`
}

// TableName gets table name of NodeSelector
func (NodeSelector) TableName() string {
	return "node_resolver_map_entries"
}

// MysqlNodeSelector holds a node selector by spiffe ID. It is equivalent to NodeSelector,
// but contains MySQL-specific gorm tags.
type MysqlNodeSelector struct {
	Model

	SpiffeID string `gorm:"size:191;unique_index:idx_node_resolver_map"`
	Type     string `gorm:"size:191;unique_index:idx_node_resolver_map"`
	Value    string `gorm:"size:191;unique_index:idx_node_resolver_map"`
}

// TableName gets table name of MysqlNodeSelector
func (MysqlNodeSelector) TableName() string {
	return "node_resolver_map_entries"
}

// RegisteredEntry holds a registered entity entry
type RegisteredEntry struct {
	Model

	EntryID       string `gorm:"unique_index"`
	SpiffeID      string
	ParentID      string
	TTL           int32
	Selectors     []Selector
	FederatesWith []Bundle `gorm:"many2many:federated_registration_entries;"`
	Admin         bool
	Downstream    bool
}

// MysqlRegisteredEntry holds a registered entity entry
type MysqlRegisteredEntry struct {
	Model

	EntryID       string `gorm:"size:191;unique_index"`
	SpiffeID      string
	ParentID      string
	TTL           int32
	Selectors     []Selector
	FederatesWith []Bundle `gorm:"many2many:federated_registration_entries;jointable_foreignkey:registered_entry_id"`
	Admin         bool
	Downstream    bool
}

var _ = MysqlRegisteredEntry(RegisteredEntry{})

// TableName gets table name of MysqlRegisteredEntry
func (MysqlRegisteredEntry) TableName() string {
	return "registered_entries"
}

// JoinToken holds a join token
type JoinToken struct {
	Model

	Token  string `gorm:"unique_index"`
	Expiry int64
}

// MysqlJoinToken holds a join token. It is equivalent to JoinToken, but contains
// MySQL-specific gorm tags.
type MysqlJoinToken struct {
	Model

	Token  string `gorm:"size:191;unique_index"` // limit varchar for DBs that don't have `innodb_large_prefix` set
	Expiry int64
}

var _ = MysqlJoinToken(JoinToken{})

// TableName gets table name of MysqlJoinToken
func (MysqlJoinToken) TableName() string {
	return "join_tokens"
}

type Selector struct {
	Model

	RegisteredEntryID uint   `gorm:"unique_index:idx_selector_entry"`
	Type              string `gorm:"unique_index:idx_selector_entry"`
	Value             string `gorm:"unique_index:idx_selector_entry"`
}

type MysqlSelector struct {
	Model

	RegisteredEntryID uint   `gorm:"size:191;unique_index:idx_selector_entry"`
	Type              string `gorm:"size:191;unique_index:idx_selector_entry"`
	Value             string `gorm:"size:191;unique_index:idx_selector_entry"`
}

var _ = MysqlSelector(Selector{})

// TableName gets table name of MysqlSelector
func (MysqlSelector) TableName() string {
	return "selectors"
}

// Migration holds version information
type Migration struct {
	Model

	// Database version
	Version int
}

// modelForDialect returns database-specific model structs
// This function uses the language change introduced in Go 1.8 where the tags
// are ignored when explicitly converting a value from one struct type to another.
// (https://golang.org/doc/go1.8#language)
// With that we can maintain separate structs for certain databases.
func modelForDialect(model interface{}, dbType string) interface{} {
	if dbType == "mysql" {
		switch v := model.(type) {
		case JoinToken:
			return MysqlJoinToken(v)
		case Bundle:
			return MysqlBundle(v)
		case AttestedNode:
			return MysqlAttestedNode(v)
		case NodeSelector:
			return MysqlNodeSelector(v)
		case RegisteredEntry:
			return MysqlRegisteredEntry(v)
		case Selector:
			return MysqlSelector(v)
		default:
			return model
		}
	}
	return model
}
