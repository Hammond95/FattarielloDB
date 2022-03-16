package v1

type RaftStats struct {
	State              string `json:"state" binding:"required"`
	Term               string `json:"term"`
	LastLogIndex       string `json:"last_log_index"`
	LastLogTerm        string `json:"last_log_term"`
	CommitIndex        string `json:"commit_index"`
	AppliedIndex       string `json:"applied_index"`
	FSMPending         string `json:"fsm_pending"`
	LastSnapshotIndex  string `json:"last_snapshot_index"`
	LastSnapshotTerm   string `json:"last_snapshot_term"`
	ProtocolVersion    string `json:"protocol_version"`
	ProtocolVersionMin string `json:"protocol_version_min"`
	ProtocolVersionMax string `json:"protocol_version_max"`
	PnapshotVersionMin string `json:"snapshot_version_min"`
	PnapshotVersionMax string `json:"snapshot_version_max"`
}
