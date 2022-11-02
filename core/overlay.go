package core

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/otiai10/copy"
	"golang.org/x/sys/unix"
)

var (
	overlaysPath   = "/etc/almost/overlays"
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

	removeOrphanOverlays()
}

func initDb(db *sql.DB) {
	sqlStmt := `
	CREATE TABLE overlays (original TEXT NOT NULL PRIMARY KEY, workdir TEXT, timestamp TEXT, persist INTEGER DEFAULT 0);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		fmt.Println("failed to initialize the overlays db: ", err)
		os.Exit(1)
	}
}

func OverlayAdd(path string, force bool, verbose bool) error {
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

	// now we need to create a temporary directory where we will store the
	// overlay structure
	workDir := fmt.Sprintf("%s/%s", overlaysPath, uuid.New().String())
	if err := os.Mkdir(workDir, 0755); err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return err
	}

	for _, dir := range []string{"upper", "lower", "work"} {
		if err := os.MkdirAll(fmt.Sprintf("%s/%s", workDir, dir), 0755); err != nil {
			fmt.Println("Error creating", dir, "directory:", err)
			return err
		}
	}

	// here is where the overlay magic happens, we are going to mount the
	// temporary directory to the original one
	if err := unix.Mount("overlay", path, "overlay", 0, fmt.Sprintf("lowerdir=%s,upperdir=%s/upper,workdir=%s/work", path, workDir, workDir)); err != nil {
		fmt.Println("Error mounting overlay:", err)
		return err
	}

	// now we need to add the overlay information to the database so that we
	// can remove it later
	registerOverlay(path, workDir, verbose)

	fmt.Printf("Your new overlay is ready at %s\n", path)

	return nil
}

func OverlayRemove(path string, keep bool, verbose bool) error {
	// first we need to check if the given has an overlay
	if !overlayCheck(path, verbose) {
		return fmt.Errorf("path %s is not overlayed", path)
	}

	original, workDir := getOverlay(path, verbose)

	// then unmount the overlay
	if err := unix.Unmount(path, 0); err != nil {
		fmt.Println("The resource is busy, re-trying killing all processes using it..")
		if err := unix.Unmount(path, unix.MNT_DETACH); err != nil {
			fmt.Println("Error unmounting overlay:", err)
			return err
		}
	}

	// remove it from the internal database
	removeOverlay(path, verbose)

	// if the keep flag is set, we need to merge the upper and lower
	// directories and copy them to the original path
	if keep {
		if err := copy.Copy(fmt.Sprintf("%s/upper", workDir), original); err != nil {
			fmt.Println("Error copying overlay to original path:", err)
			return err
		}
	}

	// finally we need to remove the temporary directory
	if err := os.RemoveAll(workDir); err != nil {
		fmt.Println("Error removing temporary directory:", err)
		return err
	}

	fmt.Printf("Overlay at %s removed\n", path)

	return nil
}

func OverlayList() map[string]string {
	overlays := make(map[string]string)

	db, err := sql.Open("sqlite3", overlaysDbPath)
	if err != nil {
		fmt.Println(err)
		return overlays
	}
	defer db.Close()

	rows, err := db.Query("SELECT original, workdir FROM overlays")
	if err != nil {
		fmt.Println(err)
		return overlays
	}
	defer rows.Close()
	for rows.Next() {
		var original, workdir string
		if err := rows.Scan(&original, &workdir); err != nil {
			fmt.Println(err)
			return overlays
		}
		overlays[original] = workdir
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return overlays
	}

	return overlays
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
	var persist int
	err = stmt.QueryRow(path).Scan(&original, &workdir, &timestamp, &persist)
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

	if !srcInfo.IsDir() {
		return fmt.Errorf("Source is not a directory")
	}

	err = copy.Copy(src, dst)
	if err != nil {
		return err
	}

	return nil
}

func removeOrphanOverlays() error {
	overlays := OverlayList()
	for original, workdir := range overlays {
		if _, err := os.Stat(workdir); os.IsNotExist(err) {
			fmt.Println("Removing orphan overlay for:", original)
			removeOverlay(original, false)
		}
	}

	return nil
}
