# This is a Python script to arrange photos according to date taken,
# user needs to set
# 1. source_path where photos are dumped at from an external storage device
# 2. destination_path where photos will be copied to

import glob
import os
from utils import *


# to be set by user
#source_path = "D:\\dump"
source_path = "D:\\dumpSource"
#destination_path = "D:\\Pictures\\Pictures\\1. Home"
destination_path = "D:\\dumpDest"

process_photos = True
process_videos = True


if __name__ == '__main__':
    # os.chdir is necessary to run glob.glbob method
    os.chdir(source_path)

    util_object = Util()
    util_object.source_path = source_path
    util_object.destination_path = destination_path

    # Photos ...
    if process_photos:
        list_of_files = glob.glob("*.jpg")
        print(f'Processing photos with date metadata')
        util_object.get_folder_tree(list_of_files)
        util_object.create_folders_and_copy_files()

        if len(util_object.failed_files):
            print(f'Processing photos with file name')
            util_object.get_folder_tree_with_name(util_object.failed_files)
            util_object.create_folders_and_copy_files()

    # Videos ...
    if process_videos:
        list_of_files = glob.glob("*.mp4")
        
        print(f'Processing videos with file name')
        util_object.get_folder_tree_with_name(list_of_files)
        util_object.create_folders_and_copy_files(videos=True)

    if len(util_object.fatal_files):
        print(f'Could not process - ')
        print(util_object.fatal_files)
    else:
        print(f'Success..')
