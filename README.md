# MyCloud (name will change)

MyCloud is about making your homeserver into a cloud server where you can upload your own files and such on it like Dropbox or OneDrive (but free!)
I'm mainly doing this to further strengthen my programming skills, learn new things on the way, and I also don't want to constantly have to open a web browser to access my server files for uploading things.

## To run the server locally

cd into the server directory in your terminal

    `cd server`

and run main.go

    `go run .`

## Uploading a file

    `curl -X POST -F "file=@C:\path\to\your\file.txt" "http://localhost:8080/upload?path=/test"`
