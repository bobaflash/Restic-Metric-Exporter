package internal

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type ResticRequester struct {
	Name              string
	Path              string
	Secret            string
	MainWg            *sync.WaitGroup
	resticRequesterWg *sync.WaitGroup
	metrics           Metrics
}

func NewResticRequester(repo Repository, wg *sync.WaitGroup, metrics Metrics) ResticRequester {

	return ResticRequester{
		Name:              repo.Name,
		Path:              repo.Path,
		Secret:            repo.Secret,
		MainWg:            wg,
		resticRequesterWg: &sync.WaitGroup{},
		metrics:           metrics,
	}
}

func (rr *ResticRequester) Run() {
	rr.resticRequesterWg.Add(2)
	go rr.getStats(rr.resticRequesterWg)
	go rr.getLastSnapshotInfo(rr.resticRequesterWg)
	rr.resticRequesterWg.Wait()
	rr.MainWg.Done()
}

type ResticStats struct {
	TotalSize      int `json:"total_size"`
	TotalFileCount int `json:"total_file_count"`
	SnapshotsCount int `json:"snapshots_count"`
}

func (rr *ResticRequester) getStats(wg *sync.WaitGroup) {

	var envs []string

	envs = append(envs, fmt.Sprintf("RESTIC_REPOSITORY=%s", rr.Path))
	envs = append(envs, fmt.Sprintf("RESTIC_PASSWORD=%s", rr.Secret))

	out, err := RunCmd("restic", []string{"stats"}, envs)

	if err != nil {
		slog.Error("Could not run 'restic stats' command", "error", err)
		wg.Done()
		return
	}

	var result ResticStats

	json.Unmarshal(out, &result)

	rr.metrics.SnapshotCount.WithLabelValues(rr.Name).Set(float64(result.SnapshotsCount))
	rr.metrics.SnapshotTotalFileCount.WithLabelValues(rr.Name).Set(float64(result.TotalFileCount))
	rr.metrics.SnapshotTotalSize.WithLabelValues(rr.Name).Set(float64(result.TotalSize))
	wg.Done()
}

type LastSnapshotInfo struct {
	Time           time.Time `json:"time"`
	Parent         string    `json:"parent"`
	Tree           string    `json:"tree"`
	Paths          []string  `json:"paths"`
	Hostname       string    `json:"hostname"`
	Username       string    `json:"username"`
	ProgramVersion string    `json:"program_version"`
	Summary        struct {
		BackupStart         time.Time `json:"backup_start"`
		BackupEnd           time.Time `json:"backup_end"`
		FilesNew            int       `json:"files_new"`
		FilesChanged        int       `json:"files_changed"`
		FilesUnmodified     int       `json:"files_unmodified"`
		DirsNew             int       `json:"dirs_new"`
		DirsChanged         int       `json:"dirs_changed"`
		DirsUnmodified      int       `json:"dirs_unmodified"`
		DataBlobs           int       `json:"data_blobs"`
		TreeBlobs           int       `json:"tree_blobs"`
		DataAdded           int       `json:"data_added"`
		DataAddedPacked     int       `json:"data_added_packed"`
		TotalFilesProcessed int       `json:"total_files_processed"`
		TotalBytesProcessed int       `json:"total_bytes_processed"`
	} `json:"summary"`
	ID      string `json:"id"`
	ShortID string `json:"short_id"`
}

func (rr *ResticRequester) getLastSnapshotInfo(wg *sync.WaitGroup) {

	var envs []string

	envs = append(envs, fmt.Sprintf("RESTIC_REPOSITORY=%s", rr.Path))
	envs = append(envs, fmt.Sprintf("RESTIC_PASSWORD=%s", rr.Secret))

	out, err := RunCmd("restic", []string{"snapshots", "latest"}, envs)

	if err != nil {
		slog.Error("Could not run 'restic snapshot latest' command", "error", err)
		wg.Done()
		return
	}

	var result []LastSnapshotInfo

	json.Unmarshal(out, &result)

	rr.metrics.SnapshotTimestamp.WithLabelValues(rr.Name, result[0].Hostname).Set(float64(result[0].Time.Unix()))
	wg.Done()
}
