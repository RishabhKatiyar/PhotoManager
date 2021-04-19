# This is a Python script to arrange photos according to date taken,
# user needs to set
# 1. destination_path where photos will be copied to
# 2. source_path where photos are dumped at from an external storage device

import glob
import os
from utils import get_folder_tree, create_folders_and_copy_files, get_folder_tree_for_videos, create_folders_and_copy_files_for_videos


# Press the green button in the gutter to run the script.
if __name__ == '__main__':
    '''
    print("Enter source_path")
    source_path = input()
    '''
    # to be set by user
    source_path = "D:\\One plus data\\Camera"
    destination_path = "D:\\Pictures\\Pictures\\1. Home"

    # getting list of files non recursive in given source_path
    os.chdir(source_path)
    list_of_files = glob.glob("*.jpg")
    # done

    folder_tree = get_folder_tree(source_path, list_of_files)
    create_folders_and_copy_files(folder_tree, source_path, destination_path)

    print("Processing Videos ...")
    list_of_files = glob.glob("*.mp4")

    folder_tree = get_folder_tree_for_videos(source_path, list_of_files)
    create_folders_and_copy_files_for_videos(folder_tree, source_path, destination_path)
