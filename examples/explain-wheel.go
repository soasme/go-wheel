package main

import (
	"fmt"
	"github.com/soasme/go-wheel"
	"os"
)

func main() {
	filename := os.Args[1]
	wheelFile, err := wheel.ParseFilename(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(wheelFile)

	r, err := wheel.Open(filename)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	fWheel, err := wheel.FindFile(r, wheelFile.PathToWheel())
	if err != nil {
		panic(err)
	}

	wheelData, err := wheel.ReadMetadata(fWheel)
	if err != nil {
		panic(err)
	}

	for _, k := range []string{"Wheel-Version", "Root-Is-Purelib", "Generator", "Build"} {
		if v, ok := wheelData.FetchOne(k); ok {
			fmt.Println(k, v)
		}
	}
	tags, _ := wheelData.FetchAll("Tag")
	fmt.Println("Tag", tags)

	fMetadata, err := wheel.FindFile(r, wheelFile.PathToMetadata())
	if err != nil {
		panic(err)
	}
	metadata, err := wheel.ReadMetadata(fMetadata)
	if err != nil {
		panic(err)
	}
	fmt.Println(metadata.Header)
	fmt.Println(metadata.Body)
}
