package workloads

type InitialTaskSpec struct {
	Name 	int 	`json:name"`
	Version string 	`json:version`
	Path    string  `json:path`
}
