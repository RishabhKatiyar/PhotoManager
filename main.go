package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"errors"
	"sync"

	"github.com/RishabhKatiyar/PhotoManager/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	//source_path      = "D:\\dumpSource"
	source_path      = ""
	//destination_path = "D:\\dumpDest"
	destination_path = ""

	process_photos = true
	process_videos = true
)

func main() {
	
	fmt.Println("Enter Source Path")
    fmt.Scanln(&source_path)

	fmt.Println("Enter Destination Path")
    fmt.Scanln(&destination_path)
	
	start := time.Now()

	// UNIX Time is faster and smaller than most timestamps
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	//Set log level
	//log_level := os.Getenv("LOG_LEVEL")
	log_level := "debug"
	switch {
	case log_level == "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case log_level == "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case log_level == "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case log_level == "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case log_level == "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case log_level == "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case log_level == "no":
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case log_level == "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case log_level == "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
	// Log a human-friendly, colorized output
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// wait for all go routines to complete
	var wg sync.WaitGroup

	
	wg.Add(1)
	// Photos
	// Passing same wg pointer to wait for copy go routines to finish
	util_object_photos := utils.Utils{Destination_path: destination_path, WaitGroupVar : &wg}
	go func () {
		defer wg.Done()
		
		if process_photos {
			var list_of_files []string
			log.Debug().Msg("Reading Photos..")
			err := filepath.Walk(source_path, func(path string, info os.FileInfo, err error) error {
				if filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".JPG"  {
					list_of_files = append(list_of_files, path)
				}
				return nil
			})

			if err != nil {
				log.Error().Stack().Err(err).Msg("")
			}

			util_object_photos.Create_folder_tree(list_of_files, false, time.Now())
			util_object_photos.Create_folders_and_copy_files(false, time.Now())
		}
	}()

	wg.Add(1)
	// Videos
	// Passing same wg pointer to wait for copy go routines to finish
	util_object_videos := utils.Utils{Destination_path: destination_path, WaitGroupVar : &wg}
	go func () {
		defer wg.Done()
		
		if process_videos {
			var list_of_files []string
			log.Debug().Msg("Reading Videos..")
			err := filepath.Walk(source_path, func(path string, info os.FileInfo, err error) error {
				if filepath.Ext(path) == ".mp4" || filepath.Ext(path) == ".MP4" {
					list_of_files = append(list_of_files, path)
				}
				return nil
			})
			
			if err != nil {
				log.Error().Stack().Err(err).Msg("")
			}

			util_object_videos.Create_folder_tree(list_of_files, true, time.Now())
			util_object_videos.Create_folders_and_copy_files(true, time.Now())
		}
	}()
	
	//
	// wait here till all copying go routines have completed
	//
	fmt.Printf("\n\nWait till copying of files is in progress..\n\n")
	wg.Wait()

	// Disply Fatal files list
	if len(util_object_photos.Fatal_files) > 0 {
		err := errors.New("Could not process files")
		log.Error().Stack().Err(err).Msgf("%v", util_object_photos.Fatal_files)
	}

	// Disply Fatal files list
	if len(util_object_videos.Fatal_files) > 0 {
		err := errors.New("Could not process files")
		log.Error().Stack().Err(err).Msgf("%v", util_object_videos.Fatal_files)
	} else {
		log.Debug().Msg("Success!")
	}

	log.Debug().Msgf("Time Taken = %v", time.Since(start))

	fmt.Println("You can exit the application now..")
	exitInput := ""
	fmt.Scanln(&exitInput)
}
