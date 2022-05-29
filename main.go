package main

import (
	"log"
	"os"
	"path"

	"s3copy/locations"
)

var programName = path.Base(os.Args[0])

func usage() {
	log.Fatalf("%s: Required Arguments <source> <destination>", programName)
}

func main() {
	// we require a source and a destination
	if len(os.Args) < 3 {
		usage()
	}

	log.Println("Starting", programName)

	// capture source and destination from the command line arguments
	source := os.Args[1]
	dest := os.Args[2]

	log.Println("Copying From", source, "to", dest)

	// resolve source
	inputSource, err := locations.ResolveSource(source)
	if err != nil {
		log.Fatalln("Error Resolving Input Source", err)
	}

	// resolve destination
	outputDestination, err := locations.ResolveDestination(dest)
	if err != nil {
		log.Fatalln("Error Resolving Output Destination", err)
	}

	log.Println("Source:", inputSource.Name())
	log.Println("Destination:", outputDestination.Name())

	written, err := outputDestination.Write(inputSource)
	if err != nil {
		log.Fatalln("Error Writing to", outputDestination.Name(), "from", inputSource.Name(), "-", err)
	}

	log.Println("Finished Writing Bytes:", written)
}
