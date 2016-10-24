package Sync

// Syncer manage sync file
type Syncer struct {
}

// SyncForce force to write file to dst file
func (s *Syncer) SyncForce() error {
	return nil
}

// Sync write to new file to old file if md5 changes
func (s *Syncer) Sync() error {
	return nil
}
