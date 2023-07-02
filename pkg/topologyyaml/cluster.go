package topologyyaml

type ClusterStatus struct {
	FeMasterHealthy bool
	Fes             []FeStatus
	Bes             []BeStatus
}

type FeStatus struct {
	Host      string `db:"Host"`
	IsMaster  bool   `db:"IsMaster"`
	QueryPort int    `db:"QueryPort"`
	Alive     bool   `db:"Alive"`
	Version   string `db:"Version"`
}

type BeStatus struct {
	Host          string `db:"Host"`
	Alive         bool   `db:"Alive"`
	AvailCapacity string `db:"AvailCapacity"`
	TotalCapacity string `db:"TotalCapacity"`
	// UsedUsedPct is the percentage of used capacity
	UsedPct string `db:"UsedPct"`
	Version string `db:"Version"`
}
