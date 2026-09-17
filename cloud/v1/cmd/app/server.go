package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func RunWebServer() {
	port := getEnv("PORT", defPort)
	host := getEnv("HOST", defHost)
	dir := getEnv("APPDIR", defDir)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "OK", "service": "backend-api"})
	})
	mux.HandleFunc("POST /api/data", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Data received successfully"}`))
	})
	RegisterStaticApp(mux, "/auth", filepath.Join(dir, "auth"))
	RegisterStaticApp(mux, "/home", filepath.Join(dir, "home"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}
		http.NotFound(w, r)
	})
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("🚀 Unified Server running on http://%s:%s\n", host, port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}

func RegisterStaticApp(mux *http.ServeMux, prefix string, appDir string) {
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		fmt.Printf("⚠️ Warning: Frontend directory %s does not exist yet.\n", appDir)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		// Strip the prefix (e.g., /auth/signin -> /signin)
		cleanPath := strings.TrimPrefix(r.URL.Path, prefix)
		cleanPath = strings.TrimPrefix(filepath.Clean(cleanPath), "/")

		if cleanPath == "." || cleanPath == "" {
			cleanPath = "index.html"
		}

		targetPath := filepath.Join(appDir, cleanPath)

		// 1. Serve direct static assets (chunks, images, CSS, JS inside _next folder, etc.)
		if info, err := os.Stat(targetPath); err == nil && !info.IsDir() {
			http.StripPrefix(prefix, http.FileServer(http.Dir(appDir))).ServeHTTP(w, r)
			return
		}

		// 2. Check if it's a directory containing an index.html (Next.js route export structure)
		if info, err := os.Stat(targetPath); err == nil && info.IsDir() {
			indexPath := filepath.Join(targetPath, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				http.ServeFile(w, r, indexPath)
				return
			}
		}

		// 3. Check if appending .html matches a file export variant
		htmlPath := targetPath + ".html"
		if _, err := os.Stat(htmlPath); err == nil {
			http.ServeFile(w, r, htmlPath)
			return
		}

		// 4. Fallback for Next.js SPA client-side routing (try 404.html, then main index.html)
		fallbackPath := filepath.Join(appDir, "404.html")
		if _, err := os.Stat(fallbackPath); err != nil {
			fallbackPath = filepath.Join(appDir, "index.html")
		}

		if _, err := os.Stat(fallbackPath); err == nil {
			http.ServeFile(w, r, fallbackPath)
			return
		}

		http.NotFound(w, r)
	}

	// Register both base route and trailing path handlers
	mux.HandleFunc(prefix, handler)
	mux.HandleFunc(prefix+"/", handler)
}
