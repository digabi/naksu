package install

import (
	"naksu/xlate"
	"naksu/progress"
	"naksu/mebroutines"

	"fmt"
	"os"
	"io"
	"net/http"
	"errors"
)

func Do_get_server(path_new_vagrantfile string) {
	const URL_VAGRANT = "http://static.abitti.fi/usbimg/qa/vagrant/Vagrantfile"

	// Create ~/ktp if missing
	progress.Set_message_xlate("Creating ~/ktp")
	var path_ktp = create_dir_ktp()
	mebroutines.Message_debug(fmt.Sprintf("path_ktp is %s", path_ktp))

  // Create ~/ktp-jako if missing
	progress.Set_message_xlate("Creating ~/ktp-jako")
	var path_ktpjako = create_dir_ktpjako()
	mebroutines.Message_debug(fmt.Sprintf("path_ktpjako is %s", path_ktpjako))

	var path_vagrantfile = path_ktp + string(os.PathSeparator) + "Vagrantfile"

	// If no given Vagrantfile try to download one
	if path_new_vagrantfile == "" {
		// Download Vagrantfile (Abitti)
		path_abitti_vagrantfile := path_vagrantfile + ".abitti"

		progress.Set_message_xlate("Downloading Abitti Vagrantfile")
		err_download := download_file(URL_VAGRANT, path_abitti_vagrantfile)
		if err_download != nil {
			mebroutines.Message_debug(fmt.Sprintf("Download failed: %v", err_download))
			mebroutines.Message_warning(xlate.Get("Could not update Abitti stickless server. Please check your network connection."))
			return
		}

		path_new_vagrantfile = path_abitti_vagrantfile
	}

	// chdir ~/ktp
	if (! mebroutines.Chdir_vagrant_directory()) {
		mebroutines.Message_error("Could not change to vagrant directory ~/ktp")
	}

	// If there is ~/ktp/Vagrantfile
	if (mebroutines.ExistsFile(path_vagrantfile)) {
		// Destroy current VM
		progress.Set_message_xlate("Destroying existing server")
		run_params_destroy := []string{"destroy","-f"}
		mebroutines.Run_vagrant(run_params_destroy)

		remove_vagrantfile(path_vagrantfile)
	}

	progress.Set_message_xlate("Copying Vagrantfile")
	err := mebroutines.CopyFile(path_new_vagrantfile, path_vagrantfile)

	if (err != nil) {
		mebroutines.Message_debug(fmt.Sprintf("Copying failed, error: %v", err))
		mebroutines.Message_error(fmt.Sprintf(xlate.Get("Error while copying new Vagrantfile: %d"), err))
	}

	progress.Set_message_xlate("Installing/updating VM: box update")
	run_params_update := []string{"box","update"}
	mebroutines.Run_vagrant(run_params_update)

	progress.Set_message_xlate("Installing/updating VM: box prune")
	run_params_prune := []string{"box","prune"}
	mebroutines.Run_vagrant(run_params_prune)

	progress.Set_message_xlate("Downloading stickless server and starting it for the first time. This takes a long time...\n\nIf the server fails to start please try to start it again from the Naksu main menu.")
	run_params_up := []string{"up"}
	mebroutines.Run_vagrant(run_params_up)
}


func create_dir_ktp () string {
	var path_ktp = mebroutines.Get_vagrant_directory()

  if (! mebroutines.ExistsDir(path_ktp)) {
    if (mebroutines.CreateDir(path_ktp) != nil) {
			mebroutines.Message_error(fmt.Sprintf(xlate.Get("Could not create ~/ktp to %s"), path_ktp))
    }
  }

	return path_ktp
}

func create_dir_ktpjako () string {
	var path_ktpjako = mebroutines.Get_mebshare_directory()

  if (! mebroutines.ExistsDir(path_ktpjako)) {
    if (mebroutines.CreateDir(path_ktpjako) != nil) {
			mebroutines.Message_error(fmt.Sprintf(xlate.Get("Could not create ~/ktp-jako to %s"), path_ktpjako))
    }
  }

	return path_ktpjako
}

func remove_vagrantfile (path_vagrantfile string) {
	// Delete Vagrantfile.bak
	if (mebroutines.ExistsFile(path_vagrantfile+".bak")) {
		err := os.Remove(path_vagrantfile+".bak")
		if (err != nil) {
			mebroutines.Message_warning(fmt.Sprintf(xlate.Get("Failed to delete %s"), path_vagrantfile+".bak"))
		}
	}

	// Rename Vagrantfile to Vagrantfile.bak
	err := os.Rename(path_vagrantfile, path_vagrantfile+".bak")
	if (err != nil) {
		mebroutines.Message_warning(fmt.Sprintf(xlate.Get("Failed to rename %s to %s"), path_vagrantfile, path_vagrantfile+".bak"))
	}
}

func download_file (url string, filepath string) error {
	mebroutines.Message_debug(fmt.Sprintf("Starting download from URL %s to file %s", url, filepath))

	out, err1 := os.Create(filepath)
	defer out.Close()
	if (err1 != nil) {
		return errors.New("Failed to create file")
	}

	resp, err2 := http.Get(url)
	if (err2 != nil) {
		return errors.New("Failed to retrieve file")
	}
	defer resp.Body.Close()

	_, err3 := io.Copy(out, resp.Body)
	if (err3 != nil) {
		return errors.New("Failed to copy body")
	}

	mebroutines.Message_debug(fmt.Sprintf("Finished download from URL %s to file %s", url, filepath))
	return nil
}

func If_http_get (url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		mebroutines.Message_debug(fmt.Sprintf("Testing HTTP GET %s and got error %v", url, err.Error()))
		return false
	}
	defer resp.Body.Close()

	mebroutines.Message_debug(fmt.Sprintf("Testing HTTP GET %s succeeded", url))

	return true
}
