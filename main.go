package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/3d0c/gmf"
)

var (
	section     *io.SectionReader
	srcfileName string
	seq         int
)

func customReader() ([]byte, int) {
	var file *os.File
	var err error

	if section == nil {
		file, err = os.Open(srcfileName)
		if err != nil {
			panic(err)
		}

		fi, err := file.Stat()
		if err != nil {
			panic(err)
		}

		section = io.NewSectionReader(file, 0, fi.Size())
	}

	b := make([]byte, gmf.IO_BUFFER_SIZE)

	n, err := section.Read(b)
	if err != nil && err == io.EOF {
		file.Close()
		return b, n
	}
	if err != nil {
		return nil, n
	}

	return b, n
}

func writeFile(b []byte) {
	name := "tmp-img/" + strconv.Itoa(seq) + ".jpg"

	fp, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := fp.Close(); err != nil {
			log.Fatal(err)
		}
		seq++
	}()

	if n, err := fp.Write(b); err != nil {
		log.Fatal(err)
	} else {
		log.Println(n, "bytes written to", name)
	}
}

func main() {
	if len(os.Args) > 1 {
		srcfileName = os.Args[1]
	} else {
		srcfileName = "SampleVideo_1280x720_2mb.mp4"
	}

	ictx := gmf.NewCtx()
	defer ictx.CloseInput()
	defer ictx.Free()

	// if err := ictx.SetInputFormat("mp4"); err != nil {
	// 	log.Fatal(err)
	// }

	avioCtx, err := gmf.NewAVIOContext(ictx, &gmf.AVIOHandlers{ReadPacket: customReader})
	defer gmf.Release(avioCtx)
	if err != nil {
		log.Fatal(err)
	}

	ictx.SetPb(avioCtx).OpenInput("")

	ist, err := ictx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		log.Println("No video stream found in", srcfileName)
	}

	fmt.Println("ictx.Duration:", ictx.Duration())
	fmt.Printf("bitrate: %d/sec\n", ictx.BitRate())
	fmt.Println(ist.CodecCtx().Width(), ist.CodecCtx().Height())
}
