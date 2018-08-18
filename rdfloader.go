package main

import (
    "io"
    "log"
    "net/http"
    "os"
    "path"

    rdf "github.com/knakk/rdf"
)

func main() {

    // Get the arguments: [source] [download-directory]
    if len(os.Args) < 2 {
        log.Print("Missing RDF source and/or download dir arguments")
        return
    }
    source := os.Args[1]
    download_dir := os.Args[2]

    void_url := ""
    src_directory := ""

    switch source {
      case "bcodmo":
        void_url = "https://www.bco-dmo.org/.well-known/void"
        src_directory = "bcodmo"
        break

      case "lter":
        void_url = "https://www.bco-dmo.org/lter/void"
        src_directory = "lter"
        break

      default:
        log.Print("Unknown Source")
        return
    }

    // Generate the full directory path.
    directory := download_dir + "/" + src_directory + "/"

    log.Print("VoID URL: "+ void_url + "\n")
    log.Print("Directory: " + directory + "\n")

    // Create the directory if it doesn't exist.
    if _, err := os.Stat(directory); os.IsNotExist(err) {
        err = os.MkdirAll(directory, 0750)
        if err != nil {
            log.Print("\nCould not create directories: " + directory + "\n ERROR: " + err.Error())
            return
        }
    } else if err != nil {
        log.Print("\nCould not check directories: " + directory + "\n ERROR: " + err.Error())
        return
    }

    // Read the contents of the VoID document.
    response, err := http.Get(void_url)
    if err != nil {
            // Send the response.
            log.Print("Could not read: " + void_url + "\n ERROR: " + err.Error())
            return
    }
    defer response.Body.Close()

    dec := rdf.NewTripleDecoder(response.Body, rdf.RDFXML)
    for triple, err := dec.Decode(); err != io.EOF; triple, err = dec.Decode() {
        // For each void:dataDump triple,
        if ("http://rdfs.org/ns/void#dataDump" == triple.Pred.String()) {
            log.Print("void:dataDump: " + triple.Obj.String())
            // Download the file to the directory.
            filepath := directory + path.Base(triple.Obj.String())
            log.Print("-- write: " + filepath)

            err := DownloadFile(filepath, triple.Obj.String())
            if err != nil {
                // Send the response.
                log.Print("Could not write: " + triple.Obj.String() + " to " + filepath + "\n ERROR: " + err.Error())
                return
            }
        }
    }
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}
