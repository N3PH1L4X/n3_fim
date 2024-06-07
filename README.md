# N3 FILE INTEGRITY MONITOR
N3 FIM is a proof of concept and a basic implementation of a File Integrity Monitor (FIM) coded in Go designed to help users know when their filesystem is being manipulated, either by modifying, deleting or creating new files.

## Objective
This is a small project I created to learn and practice the programming language Go along with other programming concepts. Also i wanted to create something cybersecurity related so a FIM was a sensible choice.
The project was inspired by this [Josh Madakor's video](https://www.youtube.com/watch?v=WJODYmk4ys8) and created using **Go version go1.22.0**. It's a similar implementation but I added some things on my own.

## Features
1. Analyze recursively the current directory to create a new baseline, gathering files, calculating their hashes and storing them in a CSV file
2. Monitor the specified routes using the information stored in the baseline file and checking if a file was either modified or created (It currently does not check if a file was deleted)
4. Configuration based functionalities. You can choose what directories in your filesystem you want to monitor by adding them to the config.yaml file. If there is no such file the program will create a default one
5. Specify the time between each iteration of file checking. This is 1 second by default
6. Multiple hashing algorithms supported. You can choose MD5, SHA256 and CRC32
7. Stores any flagged file in a .log file indicating the timestamp of detection, the file flagged and the reason why it was flagged
8. An example directory with multiple files to test the functionality

## Heads up
As stated earlier this was built as a learning project so keep in mind it may have bugs and may break sometimes. If you find a bug please consider letting me know and I would gladly fix it!

## Test the project
If you want to try it out then you can clone the repository and run the prebuilt .exe file or you can build the binaries yourself. I strongly suggest you to inspect the code and build it yourself for your specific system.
If you choose to built the binary yourself then you will need to install [Git](https://www.git-scm.com/downloads) and [Go](https://go.dev/doc/install) for your local system.
You can use these commands:

```
$ git clone https://github.com/N3PH1L4X/n3_fim.git
$ cd ./n3_fim
$ go build
```
After you build the project you will need to edit the config.yaml file and add the directories you want to monitor. After doing that you can simply run the binary with the **baseline** argument to create the new baseline (this overwrites any previous one) and then run it again with the **monitor** argument to begin monitoring.
```
$ ./n3_fim.exe baseline
$ ./n3_fim.exe monitor
```

## Collaboration
If you find it useful then feel free to use the code as you want. I'm also open to suggestions on how to make this project more stable, cleaner and better in general. Thank you!

## Screenshots
![image](https://github.com/N3PH1L4X/n3_fim/assets/71483185/b5b445d2-0d67-4cd0-85e9-042caae4d543)

![image](https://github.com/N3PH1L4X/n3_fim/assets/71483185/1376414d-7830-4594-8cc4-cbbf1a0ed6d4)

![image](https://github.com/N3PH1L4X/n3_fim/assets/71483185/8c553a47-9095-4820-8d6c-3597b4a7d082)

