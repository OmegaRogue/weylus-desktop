build:
	go build weylus-desktop
run:
	go run weylus-desktop
lint:
	golangci-lint run --color always
test-ffmpeg-rtmp:
	go run weylus-desktop | ffmpeg -f mp4 -i - -c copy -f flv -listen 1 rtmp://localhost:1935/live/app
test-ffmpeg:
	go run weylus-desktop | ffmpeg -re -i - -listen 1 -f mp4 -fflags nobuffer -movflags frag_keyframe+empty_moov http://localhost:8080/test.mp4
