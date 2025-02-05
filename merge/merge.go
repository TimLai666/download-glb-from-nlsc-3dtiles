package merge

import (
	"fmt"
	"log"
	"os"

	"github.com/qmuntal/gltf"
)

func MergeAllGLBs(outputFile string, glbDir string) {
	files, err := os.ReadDir(glbDir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v\n", err)
	}
	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	MergeGLB(outputFile, fileNames)
}

func MergeGLB(outputFile string, inputFiles []string) {
	// 建立一個新的 GLTF 結構作為合併的結果
	mergedScene := &gltf.Document{
		Asset: gltf.Asset{Version: "2.0"},
	}

	// 讀取每個 GLB 檔案並合併
	for _, file := range inputFiles {
		doc, err := gltf.Open(file)
		if err != nil {
			log.Fatalf("Failed to open GLB file: %v\n", err)
		}

		// 合併節點與網格
		mergedScene.Nodes = append(mergedScene.Nodes, doc.Nodes...)
		mergedScene.Meshes = append(mergedScene.Meshes, doc.Meshes...)
		mergedScene.Materials = append(mergedScene.Materials, doc.Materials...)
	}

	// 輸出合併後的 GLB 文件
	output, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output GLB file: %v\n", err)
	}
	defer output.Close()

	err = gltf.Save(mergedScene, outputFile)
	if err != nil {
		log.Fatalf("Failed to save merged GLB file: %v\n", err)
	}

	fmt.Printf("GLB files merged successfully into %s\n", outputFile)
}

// func main() {
// 	// 呼叫合併函式，將多個 GLB 合併成一個
// 	mergeGLB("merged.glb", []string{"model1.glb", "model2.glb", "model3.glb"})
// }
