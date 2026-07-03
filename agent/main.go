package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

const SECRET_KEY = "bngwiro" 

func handleCreateWebsite(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token != "Bearer "+SECRET_KEY {
		http.Error(w, "Akses Ditolak!", http.StatusUnauthorized)
		return
	}

	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "Error: Nama domain kosong!", http.StatusBadRequest)
		return
	}

	folderPath := "/tmp/" + domain
	exec.Command("mkdir", "-p", folderPath).Run()

	nginxConfigDir := "/tmp/nginx-configs"
	exec.Command("mkdir", "-p", nginxConfigDir).Run()

	configContent := fmt.Sprintf(`server {
        listen 80;
        server_name %s www.%s;
        root %s;
    }`, domain, domain, folderPath)

	filePath := nginxConfigDir + "/" + domain + ".conf"
	os.WriteFile(filePath, []byte(configContent), 0644)

	fmt.Fprintf(w, "Sukses mendeploy %s", domain)
}

func main() {
	http.HandleFunc("/api/panel", handleCreateWebsite)
	fmt.Println("Xeroint panel host web")
	http.ListenAndServe(":8080", nil)
}