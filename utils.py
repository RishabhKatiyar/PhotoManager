import os
from collections import namedtuple
import PIL
from PIL import Image
from enum import IntEnum


class Months(IntEnum):
    January = 1
    February = 2
    March = 3
    April = 4
    May = 5
    June = 6
    July = 7
    August = 8
    September = 9
    October = 10
    November = 11
    December = 12


def get_date_taken(file_path):
    try:
        image_data = PIL.Image.open(file_path)
        data = image_data._getexif()
        return data[36867]
    except Exception:
        print(f"Cannot Process file : {file_path}")


def get_folder_tree(path, list_of_files):
    """populating the data structure which holds folder structure metadata"""
    folder_tree = {}
    DayFilePair = namedtuple('DayFilePair', ['day', 'filename'])

    for file_name in list_of_files:
        date_taken = get_date_taken(os.path.join(path, file_name))
        if date_taken is None:
            print(f"Cannot Process file : {file_name}")
            continue
        full_date = date_taken.split(sep=" ")[0]
        _year = int(full_date.split(sep=":")[0])
        _month = int(full_date.split(sep=":")[1])
        _day = int(full_date.split(sep=":")[2])

        day_file_pair = DayFilePair(_day, file_name)

        if _year in folder_tree:
            if _month in folder_tree[_year]:
                folder_tree[_year][_month].add(day_file_pair)
            else:
                folder_tree[_year][_month] = set([day_file_pair])
        else:
            folder_tree[_year] = {}
            folder_tree[_year][_month] = set([day_file_pair])

    return folder_tree


def create_folders_and_copy_files(folder_tree, source_path, destination_path):
    """create year, month and dates folders as necessary
        then paste the files in respective folders"""
    for year_key in folder_tree:
        file_path = os.path.join(destination_path, str(year_key))
        if os.path.exists(file_path):
            pass
        else:
            print(f"Creating path {file_path}")
            os.mkdir(file_path)

        for month_key in folder_tree[year_key]:
            for month_enum_val in Months:
                if month_enum_val == month_key:
                    month_str = str(month_enum_val).split(sep=".")[1]

            if month_key >= 10:
                month_str_key = str(month_key) + " (" + month_str + ")"
            else:
                month_str_key = "0" + str(month_key) + " (" + month_str + ")"

            file_path = os.path.join(destination_path, str(year_key), str(month_str_key))
            if os.path.exists(file_path):
                pass
            else:
                print(f"Creating path {file_path}")
                os.mkdir(file_path)

            for day_file_pair in folder_tree[year_key][month_key]:
                if day_file_pair.day >= 10:
                    day_str_key = str(day_file_pair.day) + " " + month_str
                else:
                    day_str_key = "0" + str(day_file_pair.day) + " " + month_str
                file_path = os.path.join(destination_path, str(year_key), str(month_str_key), str(day_str_key))
                if os.path.exists(file_path):
                    pass
                else:
                    print(f"Creating path {file_path}")
                    os.mkdir(file_path)

                # copy the file to created or existing folder
                import shutil
                print(f"Copying file .. {os.path.join(source_path, day_file_pair.filename)} to {file_path}")
                shutil.copy(os.path.join(source_path, day_file_pair.filename), file_path)


def get_folder_tree_for_videos(path, list_of_files):
    """populating the data structure which holds folder structure metadata"""
    folder_tree = {}
    DayFilePair = namedtuple('DayFilePair', ['day', 'filename'])

    for file_name in list_of_files:
        file_name_date = file_name.split(sep="_")[1]
        _year = int(file_name_date[0:4])
        _month = int(file_name_date[4:6])
        _day = int(file_name_date[6:8])

        day_file_pair = DayFilePair(_day, file_name)

        if _year in folder_tree:
            if _month in folder_tree[_year]:
                folder_tree[_year][_month].add(day_file_pair)
            else:
                folder_tree[_year][_month] = set([day_file_pair])
        else:
            folder_tree[_year] = {}
            folder_tree[_year][_month] = set([day_file_pair])

    return folder_tree


def create_folders_and_copy_files_for_videos(folder_tree, source_path, destination_path):
    """create year, month and dates folders as necessary
        then paste the files in respective folders"""
    for year_key in folder_tree:
        file_path = os.path.join(destination_path, str(year_key))
        if os.path.exists(file_path):
            pass
        else:
            print(f"Creating path {file_path}")
            os.mkdir(file_path)

        for month_key in folder_tree[year_key]:
            for month_enum_val in Months:
                if month_enum_val == month_key:
                    month_str = str(month_enum_val).split(sep=".")[1]

            if month_key >= 10:
                month_str_key = str(month_key) + " (" + month_str + ")"
            else:
                month_str_key = "0" + str(month_key) + " (" + month_str + ")"

            file_path = os.path.join(destination_path, str(year_key), str(month_str_key))
            if os.path.exists(file_path):
                pass
            else:
                print(f"Creating path {file_path}")
                os.mkdir(file_path)

            for day_file_pair in folder_tree[year_key][month_key]:
                if day_file_pair.day >= 10:
                    day_str_key = str(day_file_pair.day) + " " + month_str
                else:
                    day_str_key = "0" + str(day_file_pair.day) + " " + month_str
                file_path = os.path.join(destination_path, str(year_key), str(month_str_key), str(day_str_key))
                if os.path.exists(file_path):
                    pass
                else:
                    print(f"Creating path {file_path}")
                    os.mkdir(file_path)

                file_path = os.path.join(destination_path, str(year_key), str(month_str_key), str(day_str_key), "Videos")
                if os.path.exists(file_path):
                    pass
                else:
                    print(f"Creating path {file_path}")
                    os.mkdir(file_path)

                # copy the file to created or existing folder
                import shutil
                print(f"Copying file .. {os.path.join(source_path, day_file_pair.filename)} to {file_path}")
                shutil.copy(os.path.join(source_path, day_file_pair.filename), file_path)
