#!/usr/bin/python

#
# Copyright (C) 2019 Nethesis S.r.l.
# http://www.nethesis.it - info@nethesis.it
#
# This script is part of Dante.
#
# Dante is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License,
# or any later version.
#
# Dante is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with Dante.  If not, see COPYING.
#

# Library functions for squidguard miners

import os
import gzip
import re

from datetime import datetime

def find_most_recent(directory, partial_file_name):
    # list all the files in the directory
    files = os.listdir(directory)

    # remove all file names that don't match partial_file_name string
    files = filter(lambda x: x.find(partial_file_name) > -1, files)

    # create a dict that contains list of files and their modification timestamps
    name_n_timestamp = dict([(x, os.stat(directory+x).st_mtime) for x in files])

    # return the file with the latest timestamp
    return max(name_n_timestamp, key=lambda k: name_n_timestamp.get(k))

def grep_pattern(pattern, file_or_list):
    lines = []

    for line in file_or_list:
        if re.search(pattern, line):
            lines.append(line)
    return lines

def grep_blocked_lines():
    log_directory = "/var/log/ufdbguard/"
    log_file = log_directory + "ufdbguardd.log"
    today = datetime.today().strftime('%Y-%m-%d')

    # Exit if ufdb is not installed or log is empty
    if (not os.path.isfile(log_file)) or (not os.path.isfile('/etc/e-smith/db/configuration/defaults/ufdb/type')):
        exit(1)

    today_lines = None

    # scan current log file for today logs
    with open(log_file, 'r') as log_file_fd:
        today_lines = grep_pattern(r"{0}".format(today), log_file_fd)

    previous_log_file = find_most_recent(log_directory, "ufdbguardd.log-")
    previous_log_file = log_directory + previous_log_file

    # scan previous log file for today logs
    if previous_log_file.endswith(".gz"):
        with gzip.open(previous_log_file, 'rb') as previous_log_file_fd:
            today_lines_previous_log = grep_pattern(r"{0}".format(today), previous_log_file_fd)
            today_lines = today_lines + today_lines_previous_log
    else:
        with open(previous_log_file) as previous_log_file_fd:
            today_lines_previous_log = grep_pattern(r"{0}".format(today), previous_log_file_fd)
            today_lines = today_lines + today_lines_previous_log

    blocked_lines = grep_pattern(r' BLOCK ', today_lines)
    return blocked_lines
