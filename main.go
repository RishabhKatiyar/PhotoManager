package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RishabhKatiyar/PhotoManager/utils"
)

var (
	source_path      = "D:\\dumpSource"
	destination_path = "D:\\dumpDest"

	process_photos = true
	process_videos = true
)

func main() {
	util_object := utils.Utils{Destination_path: destination_path}

	// Photos
	if process_photos {
		var list_of_files []string
		fmt.Println("Processing photos with date metadata")
		err := filepath.Walk(source_path, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".jpg" {
				list_of_files = append(list_of_files, path)
			}
			return nil
		})

		if err != nil {
			panic(err)
		}

		util_object.Get_folder_tree(list_of_files)
		util_object.Create_folders_and_copy_files(false)

		if len(util_object.Failed_files) > 0 {
			fmt.Println("Processing photos with file name")
			util_object.Get_folder_tree_with_name(util_object.Failed_files)
			util_object.Create_folders_and_copy_files(false)
		}
	}

	// Videos
	if process_videos {
		var list_of_files []string
		fmt.Println("Processing photos with date metadata")
		err := filepath.Walk(source_path, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".mp4" {
				list_of_files = append(list_of_files, path)
			}
			return nil
		})
		
		if err != nil {
			panic(err)
		}

		fmt.Println("Processing videos with file name")
		util_object.Get_folder_tree_with_name(list_of_files)
		util_object.Create_folders_and_copy_files(true)
	}

	// Disply Fatal files list
	if len(util_object.Fatal_files) > 0 {
		fmt.Println("Could not process")
		fmt.Printf("%v", util_object.Fatal_files)
	} else {
		fmt.Println("Success!")
	}
}
