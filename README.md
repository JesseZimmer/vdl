# vdl
 A direct download utility for the Miyoo Mini Plus. Note: this utility will only work with roms hosted on Vimm's Lair, and is an unofficial tool made by me for me. This is the first project I've made in Go, and was mostly done for fun and to learn, so please take it easy on me when reading the code, and I'd be happy to take any feedback on improvements that could be made. :)

 This was forked from [DanCousins/vdl](https://github.com/DanCousins/vdl) to build out a ROM browser menu

## Installation and Usage


Grab the vdl file from the releases and store it on your SD card somewhere. I put mine in the root directory, as that's what loads up first when you enter the Terminal on the console. This gives you easy access to the download utility with minimal button presses.

Navigate to the Terminal on your Miyoo (in the Apps folder, and can be installed using the Package Manager if you don't have it already). Run the download utility by typing: 
```
./vdl
```
This assumes you've stored it in the root of your SD card. You can use the SELECT button to tab, and the START button as enter. So you can type "./v", then press SELECT to autocomplete, and then press START to run the application. 

## Console Compatibility
| Console  | Compatible | Notes  |
| ------------- | ------------- | ------------- |
| Atari 2600  | ✔️  | -  |
| Atari 5600  | ✔️  | -  |
| Nintendo  | ✔️  | -  |
| Master System  | ✔️  | -  |
| Atari 7800  | ✔️  | -  |
| Gensis  | ✔️  | -  |
| Super Nintendo  | ✔️  | -  |
| Sega 32X  | ✔️  | -  |
| Playstation  | ✔️  | Only single disc games currently supported, and they can be very slow to extract. ~28 seconds for a 57MB game uncompressed. |
| Game Boy  | ✔️  | -  |
| Lynx  | ✔️  | -  |
| Game Gear  | ✔️  | -  |
| Virtual Boy  | ✔️  | -  |
| Game Boy Colour  | ✔️  | -  |
| Game Boy Advance  | ✔️  | -  |
| Nintendo DS  | ✔️  | -  |

## Building From Source
Install Go, clone the repo, modify your environment variables to compile for ARM and Linux, this is how I've done it in Powershell:
```
$Env:GOOS = "linux"; $Env:GOARCH = "arm"
```
In the folder you cloned run:
```
go build
```
It should spit you out an executable which you can then transfer to your Miyoo and run as above. 

## Planned Features
- Download progress bar
- Better error handling around failed downloads
- Rom browser to make solution self-contained - Done
- Support for multi-disc downloads
- Investigate faster unarchiving of PS1 games
