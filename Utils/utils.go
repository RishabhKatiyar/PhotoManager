package utils

import (	
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"errors"
	"sync"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rs/zerolog/log"
)

type Months int

const (
	January Months = iota + 1
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func (e Months) String() string {
	switch e {
	case January:
		return "January"
	case February:
		return "February"
	case March:
		return "March"
	case April:
		return "April"
	case May:
		return "May"
	case June:
		return "June"
	case July:
		return "July"
	case August:
		return "August"
	case September:
		return "September"
	case October:
		return "October"
	case November:
		return "November"
	case December:
		return "December"
	default:
		return ""
	}
}

type Utils struct {
	Destination_path	string
	Folder_tree        	map[int]map[int]map[int][]string	
	Fatal_files 		[]string   
	WaitGroupVar 		*sync.WaitGroup                   
}

func (u *Utils) get_date_taken_from_exif(file_path string) (int, int, int, error) {
	year, month, day := 0, 0, 0
	f, err := os.Open(file_path)
	if err != nil {
			return year, month, day, err
		}
	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
			return year, month, day, err
		}

	date_taken, err := x.DateTime()
	year = date_taken.Year()
	month = int(date_taken.Month())
	day = date_taken.Day()

	return year, month, day, err
}

func (u *Utils) get_date_taken_from_file_name(file_name string) (int, int, int, error) {
	year, month, day := 0, 0, 0
	var err error
	fileBaseName := filepath.Base(file_name)
	idx := strings.Index(fileBaseName, "-")
	if idx < 0 {
		idx = strings.Index(fileBaseName, "_")
	}
	if idx > 0 {
		year, _ = strconv.Atoi(fileBaseName[idx+1 : idx+5])
		month, _ = strconv.Atoi(fileBaseName[idx+5 : idx+7])
		day, _ = strconv.Atoi(fileBaseName[idx+7 : idx+9])
	} else {
		err = errors.New("Cannot get date taken from file name")
	}
	
	return year, month, day, err
}
	
/*
populate the data structure which holds folder tree metadata
*/
func (u *Utils) Create_folder_tree(list_of_files []string, calculate_date_taken_from_file_name bool, start_time time.Time) error {
	defer func() {
		log.Debug().Msgf("Time Elapsed in creating folder tree is %v", time.Since(start_time))
	}()

	folder_tree := make(map[int]map[int]map[int][]string)
	
	for _, file_name := range list_of_files {
		year, month, day := 0, 0, 0
		var err error
		if calculate_date_taken_from_file_name {
			year, month, day, err = u.get_date_taken_from_file_name(file_name)
			if err != nil {			
				log.Error().Stack().Err(err).Msgf("Date time not present in file name %s, queueing them to fatal files list", file_name)
				u.Fatal_files = append(u.Fatal_files, file_name)
				continue
			}		
		} else {
			year, month, day, err = u.get_date_taken_from_exif(filepath.Join(file_name))
			if err != nil {
				log.Error().Stack().Err(err).Msgf("Date taken not found in exif for file %s, trying with file name", file_name)
				year, month, day, err = u.get_date_taken_from_file_name(file_name)
				if err != nil {			
					log.Error().Stack().Err(err).Msgf("Date time not present in file name %s, queueing them to fatal files list", file_name)
					u.Fatal_files = append(u.Fatal_files, file_name)
					continue
				}
			}
		}

		if _, found := folder_tree[year][month][day]; found {
			folder_tree[year][month][day] = append(folder_tree[year][month][day], file_name)
		} else {
			if _, found := folder_tree[year]; !found {
				folder_tree[year] = make(map[int]map[int][]string)
			}
			if _, found := folder_tree[year][month]; !found {
				folder_tree[year][month] = make(map[int][]string)
			}
			folder_tree[year][month][day] = make([]string, 0)
			folder_tree[year][month][day] = append(folder_tree[year][month][day], file_name)
		}
	}
	u.Folder_tree = folder_tree
	log.Debug().Msg("Folder Tree Generated")
	return nil
}

func (u *Utils) Copy_files(file_list []string, file_path string) {
	defer u.WaitGroupVar.Done()

	if _, err := os.Stat(file_path); os.IsNotExist(err) {
		os.MkdirAll(file_path, os.ModeDir)
	}

	for _, fileName := range file_list {
		// copy the file to created or existing folder
		dest := filepath.Join(file_path, filepath.Base(fileName))
		log.Debug().Msgf("Copying file .. %s to %s", fileName, dest)
		_, err := u.copy(fileName, dest)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
		}
	}
}

/*
create year, month and dates folder tree as necessary ands
then paste the files in their respective folders
*/
func (u *Utils) Create_folders_and_copy_files(videos bool, start_time time.Time) error {
	defer func() {
		log.Debug().Msgf("Time Elapsed in copying files is %v", time.Since(start_time))
	}()
	
	destination_path := u.Destination_path

	for year_key, months := range u.Folder_tree {
		for month_key, days := range months {
			month_str := fmt.Sprintf("%v", Months(month_key))
			month_str_key := ""
			if month_key >= 10 {
				month_str_key = fmt.Sprintf("%d", month_key) + " (" + month_str + ")"
			} else {
				month_str_key = "0" + fmt.Sprintf("%d", month_key) + " (" + month_str + ")"
			}
			for day_key, file_list := range days {
				day_str_key := ""
				if day_key >= 10 {
					day_str_key = fmt.Sprintf("%d", day_key) + " " + month_str
				} else {
					day_str_key = "0" + fmt.Sprintf("%d", day_key) + " " + month_str
				}
				
				file_path := filepath.Join(destination_path, fmt.Sprintf("%d", year_key), month_str_key, day_str_key)
				
				if videos {
					file_path = filepath.Join(file_path, "Videos")
				}
				
				u.WaitGroupVar.Add(1)
				go u.Copy_files(file_list, file_path)
			}
		}
	}

	return nil
}

func (u *Utils) copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
