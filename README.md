# Final Fantasy IX Save Converter ![](.github/logo.jpg)
![build](https://github.com/EuanFH/ff9SaveLib/workflows/build/badge.svg?branch=master)
[![codecov](https://codecov.io/gh/EuanFH/ff9SaveLib/branch/master/graph/badge.svg?token=YRKO4WLF3T)](https://codecov.io/gh/EuanFH/ff9SaveLib/)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=EuanFH_ff9SaveLib&metric=bugs)](https://sonarcloud.io/dashboard?id=EuanFH_ff9SaveLib)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=EuanFH_ff9SaveLib&metric=code_smells)](https://sonarcloud.io/dashboard?id=EuanFH_ff9SaveLib)
[![GoDoc](https://godoc.org/github.com/nitishm/go-rejson?status.svg)](https://pkg.go.dev/github.com/euanfh/ff9SaveLib)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
## Description 📖
Full implementation of the 2016 FF9 save system. 
Allows you to read, write, and edit any 2016 save.

This library also has a command line utility if you just want to convert your saves between json and binary format.

### Supported Formats
 - Encrypted Save file used in the PC version
 - Json save file used in the Switch version
 
### Unsupported Formats
 - PS1 save file
 - original PC save file

## Warning ⚠️
Back up your save files before using this converter. I can't be 100% certain all save files will convert properly.

**Converting to switch saves is lossy.** Some information is not carried over to the switch save files since it doesn't
include those fields.

Here are the list of fields not carried over:
 - Auto Login
 - System Achievement Statuses
 - Screen Rotation
 
The command line utility will initialize these fields with default values.

## Bug Reporting 🐛
If this tool generates a corrupted save file please report it.
Please submit an issue on github including:
 - What is corrupted in the game
 - Attaching both save files original and generated
 - What direction you where converting switch to pc or pc to switch 

## Command line Utility ️🖥️
After downloading the application you will need to run it from the command prompt or terminal depending on your OS.

To use this application will ether need the PC save file located at 

AppData\LocalLow\SquareEnix\FINAL FANTASY IX\Steam\EncryptedSavedData\SavedData_ww.dat

or a switch save which can be dumped using homebrew like checkpoint.

```
Final Fantasy 9 Save Converter
Usage:
        ff9sc FILE DIRECTORY
        ff9sc DIRECTORY FILE
Examples:
        $ ff9sc SavedData_ww.dat switchSavesFolder
        $ ff9sc switchSavesFolder SavedData_ww.dat
```

The first file path in the command is what you're converting from the second is what your converting too

## Download 💽
Click the download link and select the version for your operating system 

[Download](https://github.com/EuanFH/ff9SaveLib/releases)
