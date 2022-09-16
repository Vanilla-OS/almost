package core

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sys/unix"
)

var (
	overlaysPath = "/etc/almost/overlays"
	overlaysDbPath = "/etc/almost/overlays.db"
)

func init() {
	if !RootCheck(false) {
		return
	}

	if _, err := os.Stat(overlaysPath); os.IsNotExist(err) {
		os.Mkdir(overlaysPath, 0755)
	}

	if _, err := os.Stat(overlaysDbPath); os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", overlaysDbPath)
		if err != nil {
			panic(err)
		}
		initDb(db)
		defer db.Close()
	}
}

func initDb(db *sql.DB) {
	sqlStmt := `
	CREATE TABLE overlays (original TEXT NOT NULL PRIMARY KEY, workdir TEXT, timestamp TEXT);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		fmt.Println("failed to initialize the overlays db: ", err)
		os.Exit(1)
	}
}

func OverlayAdd(path string, force bool, verbose bool) error {
	if verbose {
		fmt.Println("Preparing a new overlay for:", path)
	}
	// first we need to check if the given path exists, to avoid overlaying
	// non-existing directories which is not the desired behaviour
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path %s does not exist", path)
	}

	// now we need to check if the path is already overlayed if it is, we
	// need to check if the force flag is set to allow overwriting the
	// existing overlay
	if overlayCheck(path, verbose) && !force {
		return fmt.Errorf("path %s is already overlayed, remove it first", path)
	}

	// now we need to create a temporary directory where we will store a copy
	// of the original directory contents
	workDir := fmt.Sprintf("%s/%s", overlaysPath, uuid.New().String())
	if err := os.Mkdir(workDir, 0755); err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return err
	}
	if err := unix.Mount("tmpfs", workDir, "tmpfs", 0, ""); err != nil {
		fmt.Println("Error mounting a tmpfs to the temporary directory:", err)
		return err
	}
	
	if err := copyDir(path, workDir, verbose); err != nil {
		fmt.Println("Error copying directory:", err)
		return err
	}

	// here is where the magic happens, we bind mount the temporary directory
	// to the original directory so that any changes made to the original
	// directory will be stored in the temporary directory
	if err := unix.Mount(workDir, path, "", unix.MS_BIND, ""); err != nil {
		fmt.Println("Error binding mount:", err)
		return err
	}
	
	// now we need to add the overlay information to the database so that we
	// can remove it later
	registerOverlay(path, workDir, verbose)

	return nil
}

func OverlayRemove(path string, keep bool, verbose bool) error {
	if verbose {
		fmt.Println("Removing overlay for:", path)
	}

	// first we need to check if the given has an overlay
	if !overlayCheck(path, verbose) {
		return fmt.Errorf("path %s is not overlayed", path)
	}

	original, workDir := getOverlay(path, verbose)

	// we are going to unmount the temporary directory from the original one
	if err := unix.Unmount(original, 0); err != nil {
		fmt.Println("Error unmounting the overlay:", err)
		return err
	}

	// here we check if the user wants to keep the temporary directory
	// or trash it, in the first case we copy the contents of the temporary
	// directory to the original one
	if keep {
		if err := copyDir(workDir, original, verbose); err != nil {
			fmt.Println("Error copying directory:", err)
			return err
		}
	}

	// now we need to remove the overlay information from the database and
	// free the path for future overlays
	removeOverlay(path, verbose)

	return nil
}

func OverlayList() ([]string, error) {
	db, err := sql.Open("sqlite3", overlaysDbPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	rows, err := db.Query("SELECT original FROM overlays")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rows.Close()

	var overlays []string
	for rows.Next() {
		var original string
		err = rows.Scan(&original)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		overlays = append(overlays, original)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return overlays, nil
}

func overlayCheck(path string, verbose bool) bool {
	if verbose {
		fmt.Println("Checking if", path, "is overlayed")
	}

	db, err := sql.Open("sqlite3", overlaysDbPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM overlays WHERE original = ?")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	var original string
	var workdir string
	var timestamp string
	err = stmt.QueryRow(path).Scan(&original, &workdir, &timestamp)
	return err == nil
}

func getOverlay(path string, verbose bool) (string, string) {
	if verbose {
		fmt.Println("Getting overlay information for:", path)
	}

	db, err := sql.Open("sqlite3", overlaysDbPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT original, workdir FROM overlays WHERE original = ?")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer stmt.Close()

	var original, workDir string
	err = stmt.QueryRow(path).Scan(&original, &workDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return original, workDir
}

func registerOverlay(path, workDir string, verbose bool) {
	if verbose {
		fmt.Println("Registering overlay for:", path)
	}

	db, err := sql.Open("sqlite3", overlaysDbPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	
	stmt, err := db.Prepare("INSERT INTO overlays(original, workdir, timestamp) VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = stmt.Exec(path, workDir, time.Now().Format(time.RFC3339))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer stmt.Close()
}

func removeOverlay(path string, verbose bool) {
	if verbose {
		fmt.Println("Removing overlay for:", path)
	}

	db, err := sql.Open("sqlite3", overlaysDbPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	
	stmt, err := db.Prepare("DELETE FROM overlays WHERE original = ?")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = stmt.Exec(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer stmt.Close()
}

func copyDir(src, dst string, verbose bool) error {
	if verbose {
		fmt.Println("Copying", src, "to", dst)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath, verbose); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath, verbose); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string, verbose bool) error {
	if verbose {
		fmt.Println("Copying file", src, "to", dst)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}
