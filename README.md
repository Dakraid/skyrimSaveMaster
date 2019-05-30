# SkyrimSaveMaster
SkyrimSaveMaster is a currently work in progress project that will consist of a library and an utility using aforementioned library.

## Purpose
This project is the results of a desire to be able to fix a renamed esp or assign old references to the new mods (but that still retains the same structure). 

Example would be a merge of two already in use plugins into one. The contents of the plugin are pretty much the same as with the individual plugins, but due to the changed esp the game could fail to realize this.

The tool is supposed to help with that alongside a set of other functionality like save editing and cleaning.

## Limitations
There is no telling if my concept will actually work in practice or what challenges I will face while continuing work on this project.

Right now the tool is also mainly focused on TES V - Skyrim Classic, but support for the remaster TES V - Skyrim: Special Edition is planned.

## Library
The internal core of this utility will be based on a C-shared DLL written in GoLang, which will allow other programmers to possibly use this library to write their own set of tools.

The library mainly handles reading of and writing to the save file, with a few generic functions that I deem generally useful. After the main functionality is implemented and tested, the library will be isolated from the tool. At this point I will also focus on adding a variety of functions as part of an API.

## Utility
The utility will be based on the library, but feature additional functionality which will be contained in its own code. This will include any algorithms and procedures that process the data from the save, as well as any other parts that are required.

## License and User Agreement
This library and the utility will be freely available under the MIT license terms.

I am not responsible for any broken data or files that has resulted as part of using either the library or the utility. Use of either is at your own risk.

## Disclaimer
This library and utility is not made, guaranteed, or supported by Zenimax, Bethesda Game Studios, or any of their affiliates.