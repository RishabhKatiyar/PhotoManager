package utils

import (	
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"errors"
	"sync"

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
	Failed_files 		[]string
	Fatal_files 		[]string   
	WaitGroupVar 		*sync.WaitGroup                   
}

func (u *Utils) get_date_taken(file_path string) (time.Time, error) {
	f, err := os.Open(file_path)
	if err != nil {
			return time.Now(), err
		}
	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
			return time.Now(), err
		}

	tm, err := x.DateTime()

	return tm, err
	}
	
/*
populating the data structure which holds folder structure metadata
*/
func (u *Utils) Create_folder_tree(list_of_files []string) error {
	folder_tree := make(map[int]map[int]map[int][]string)
	
	for _, file_name := range list_of_files {
		date_taken, err := u.get_date_taken(filepath.Join(file_name))
		if err != nil {
			log.Error().Stack().Err(err).Msgf("Cannot Process file %s because date taken not found in exif, queueing them to process by name", file_name)
			u.Failed_files = append(u.Failed_files, file_name)
			continue
		}
		year := date_taken.Year()
		month := date_taken.Month()
		day := date_taken.Day()

		if _, found := folder_tree[year][int(month)][day]; found {
			folder_tree[year][int(month)][day] = append(folder_tree[year][int(month)][day], file_name)
		} else {
			if _, found := folder_tree[year]; !found {
				folder_tree[year] = make(map[int]map[int][]string)
			}
			if _, found := folder_tree[year][int(month)]; !found {
				folder_tree[year][int(month)] = make(map[int][]string)
			}
			folder_tree[year][int(month)][day] = make([]string, 0)
			folder_tree[year][int(month)][day] = append(folder_tree[year][int(month)][day], file_name)
		}
	}
	u.Folder_tree = folder_tree
	return nil
}

func (u *Utils) Create_folder_tree_with_name(list_of_files []string) error {
	folder_tree := make(map[int]map[int]map[int][]string)

	for _, file_name := range list_of_files {
		fileBaseName := filepath.Base(file_name)
	
		idx := strings.Index(fileBaseName, "-")
		if idx < 0 {
			idx = strings.Index(fileBaseName, "_")
		}
		if idx > 0 {
			year, _ := strconv.Atoi(fileBaseName[idx+1 : idx+5])
			month, _ := strconv.Atoi(fileBaseName[idx+5 : idx+7])
			day, _ := strconv.Atoi(fileBaseName[idx+7 : idx+9])
						
		if _, found := folder_tree[year][int(month)][day]; found {
				folder_tree[year][int(month)][day] = append(folder_tree[year][int(month)][day], file_name)
		} else {
				if _, found := folder_tree[year]; !found {
					folder_tree[year] = make(map[int]map[int][]string)
				}
				if _, found := folder_tree[year][int(month)]; !found {
					folder_tree[year][int(month)] = make(map[int][]string)
				}
				folder_tree[year][int(month)][day] = make([]string, 0)
				folder_tree[year][int(month)][day] = append(folder_tree[year][int(month)][day], file_name)
			}
		} else {
			err := errors.New("Cannot Process file")
			log.Error().Stack().Err(err).Msgf("Date time not present in file name %s, queueing them to fatal files list", file_name)
			u.Fatal_files = append(u.Fatal_files, file_name)
		}
	}
	u.Folder_tree = folder_tree
	return nil
}

func (u *Utils) Copy_files(file_list []string, file_path string) {
	defer u.WaitGroupVar.Done()
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
create year, month and dates folders as necessary
then paste the files in respective folders
*/
func (u *Utils) Create_folders_and_copy_files(videos bool) error {
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
				if _, err := os.Stat(file_path); os.IsNotExist(err) {
					os.MkdirAll(file_path, os.ModeDir)
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
