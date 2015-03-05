#Lib Folder
This folder contains the library files used in this project.
A .gitignore rule is used to prevent them from being uploaded to GitHub. So you'll
need to download the required libraries yourself and put them in this folder.
This is a temporary solution until I've written a ant or maven script that downloads the
required libraries automatically.

##Temporary List of Libraries used:
 + **JNetPcap** http://jnetpcap.com/download

##Notes:
The jnetpcap.jar must support the installed native .so files on Linux to work properly
otherwise you may get error: "size of array must be max_id_count size".
**TODO:** check versions, add link for correct version.