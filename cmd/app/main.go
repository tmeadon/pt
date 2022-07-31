package main

import (
	"log"
	"os"
	"time"

	"github.com/tmeadon/pt/pkg/backupmgr"
	"github.com/tmeadon/pt/pkg/blobstorage"
	"github.com/tmeadon/pt/pkg/webapp"
)

var dbPath string = "pt-gorm.db"
var backupContainer *blobstorage.BackupContainer

func main() {
	getBackupContainer()
	restoreDbIfNeeded()
	go backupRoutine()

	s := webapp.NewServer(dbPath)
	s.Start()
}

func restoreDbIfNeeded() {
	if restoreNeeded() {
		log.Print("database file missing, restoring from backup")
		err := backupmgr.RestoreFromLatest(dbPath, backupContainer)

		if err != nil {
			log.Fatalf("failed to restore backup: %v", err)
		}
	}
}

func restoreNeeded() bool {
	if _, err := os.Stat(dbPath); err == nil {
		return false
	}

	return true
}

func getBackupContainer() {
	sas, ok := os.LookupEnv("PT_BACKUP_SAS")
	if !ok {
		log.Fatal("PT_BACKUP_SAS environment variable not set")
	}

	container, err := blobstorage.NewBackupContainer(sas)
	if err != nil {
		log.Fatalf("unable to connect to backup container: %v", err)
	}

	backupContainer = container
}

func backupRoutine() {
	for {
		time.Sleep(8 * time.Hour)

		err := backupmgr.BackupAndShip(dbPath, backupContainer)

		if err != nil {
			log.Printf("failed to execute database backup: %v", err)
		} else {
			log.Printf("backup completed")
		}
	}
}
