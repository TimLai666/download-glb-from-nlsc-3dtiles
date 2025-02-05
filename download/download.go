package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/HazelnutParadise/Go-Utils/asyncutil"
)

func DownloadGLBsFromJsonUrl(jsonUrl string, outputDir string, numGoroutines ...int) {
	// 建立輸出目錄
	makeDirIfNotExist(outputDir)
	// 取得所有 GLB 檔案的 URL
	glbUrls := getGlbUrlsFromJsonUrl(jsonUrl)

	var maxParallelDownloads int
	if len(numGoroutines) > 1 {
		log.Fatalf("Too many arguments for DownloadGLBsFromJsonUrl function")
	} else if len(numGoroutines) == 0 {
		// 一次下載10個檔案
		maxParallelDownloads = 60
	} else {
		maxParallelDownloads = numGoroutines[0]
	}

	fmt.Printf("Downloading %d GLB files to %s\n", len(glbUrls), outputDir)
	asyncutil.ParallelForEach(glbUrls, func(_ int, url string) uint8 {
		DownloadGLB(url, outputDir, DownloadOptions{DontCheckDirExist: true})
		return 0
	}, maxParallelDownloads)
	fmt.Println("All downloads completed. Saved to: ", outputDir)
}

type DownloadOptions struct {
	DontCheckDirExist bool
}

func DownloadGLB(url string, outputDir string, options ...DownloadOptions) {
	if len(options) > 1 {
		log.Fatalf("Too many arguments for DownloadGLB function")
	}
	if len(options) == 0 || !options[0].DontCheckDirExist {
		makeDirIfNotExist(outputDir)
	}
	// 下載 GLB 檔案
	var resp *http.Response
	for {
		var err error
		resp, err = http.Get(url)
		if err != nil {
			log.Printf("Failed to download GLB file: %v, retrying\n", err)
			continue
		}
		break
	}

	defer resp.Body.Close()

	// 建立輸出檔案
	fileName := filepath.Base(url)
	outputPath := filepath.Join(outputDir, fileName)
	output, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create output GLB file: %v\n", err)
	}
	defer output.Close()

	// 將下載的內容寫入輸出檔案
	_, err = io.Copy(output, resp.Body)
	if err != nil {
		log.Fatalf("Failed to save GLB file: %v\n", err)
	}

	fmt.Printf("GLB file downloaded successfully: %s\n", outputPath)
}

func makeDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Failed to create output directory: %v\n", err)
		}
	}
}
