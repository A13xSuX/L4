package metrics

import (
	"fmt"
	"net/http"
	"runtime"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// доступность сервера
	fmt.Fprintln(w, "# HELP app_up Shows that the application is running")
	fmt.Fprintln(w, "# TYPE app_up gauge")
	fmt.Fprintln(w, "app_up 1")

	// кол-во байт занятых живыми объектами
	fmt.Fprintln(w, "# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use")
	fmt.Fprintln(w, "# TYPE go_memstats_alloc_bytes gauge")
	fmt.Fprintf(w, "go_memstats_alloc_bytes %d\n", mem.Alloc)

	// кол-во байт выделенных за все время работы
	fmt.Fprintln(w, "# HELP go_memstats_allocated_bytes_total Total number of bytes allocated even if freed")
	fmt.Fprintln(w, "# TYPE go_memstats_allocated_bytes_total counter")
	fmt.Fprintf(w, "go_memstats_allocated_bytes_total %d\n", mem.TotalAlloc)

	// кол-во аллокаций
	fmt.Fprintln(w, "# HELP go_memstats_mallocs_total Total number of mallocs")
	fmt.Fprintln(w, "# TYPE go_memstats_mallocs_total counter")
	fmt.Fprintf(w, "go_memstats_mallocs_total %d\n", mem.Mallocs)

	// кол-во памяти полученной от OC
	fmt.Fprintln(w, "# HELP go_memstats_sys_bytes Number of bytes obtained from system")
	fmt.Fprintln(w, "# TYPE go_memstats_sys_bytes gauge")
	fmt.Fprintf(w, "go_memstats_sys_bytes %d\n", mem.Sys)

	// текущий объем памяти heap
	fmt.Fprintln(w, "# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use")
	fmt.Fprintln(w, "# TYPE go_memstats_heap_alloc_bytes gauge")
	fmt.Fprintf(w, "go_memstats_heap_alloc_bytes %d\n", mem.HeapAlloc)

	// текущее количество объектов в heap
	fmt.Fprintln(w, "# HELP go_memstats_heap_objects Number of allocated heap objects")
	fmt.Fprintln(w, "# TYPE go_memstats_heap_objects gauge")
	fmt.Fprintf(w, "go_memstats_heap_objects %d\n", mem.HeapObjects)

	// кол-во сборок мусора
	fmt.Fprintln(w, "# HELP go_gc_cycles_total Number of completed gc cycles")
	fmt.Fprintln(w, "# TYPE go_gc_cycles_total counter")
	fmt.Fprintf(w, "go_gc_cycles_total %d\n", mem.NumGC)

	// последний вызов gc
	fmt.Fprintln(w, "# HELP go_gc_last_time_seconds Unix time of the last gc in seconds")
	fmt.Fprintln(w, "# TYPE go_gc_last_time_seconds gauge")
	fmt.Fprintf(w, "go_gc_last_time_seconds %.9f\n", float64(mem.LastGC)/1e9)
}
