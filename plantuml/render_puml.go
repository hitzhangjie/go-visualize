package plantuml

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hitzhangjie/go-visualize/http"
)

func RenderPlantUML(pumlFile string) error {

	// download plantuml.jar if not found
	home, _ := os.UserHomeDir()
	jar := filepath.Join(home, "plantuml.jar")
	_, err := os.Lstat(jar)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err := http.DownloadFile(jar, "https://nchc.dl.sourceforge.net/project/plantuml/plantuml.jar")
		if err != nil {
			return fmt.Errorf("download plantuml.jar error: %v", err)
		}
	}

	cmd := exec.Command("java", "-jar", jar, "-progress", pumlFile)
	if msg, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("run plantuml error: %v,\nmsg:%s", err, msg)
	}
	return nil
}
