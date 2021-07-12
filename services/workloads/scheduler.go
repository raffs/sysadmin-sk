/**
* Follows a simple producer-consumer model at first. Firstly need to create a stable set of
* communication.
 */
package workloads


func ScheduleTask(spec []*InitialTaskSpec, jobs chan<- InitialTaskSpec) {
	defer close(jobs)
	for _, task := range spec {
		jobs <- *task
	}
}


