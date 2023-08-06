package main

import (
	"errors"
    "fmt"
	"io"
    "io/ioutil"
	"net/http"
	"os"
    "strings"
)

const baseURL = "https://raw.githubusercontent.com/deosjr/PokeCiv/master/Data/Graphics/"
const from = 0
const to = 494

func main() {
    downloadImage("Battlebacks", "battlebgIndoorA.png")

    // test: crawl all of bulbapedia for images
    pokedex := "https://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_National_Pok%C3%A9dex_number"
	response, err := http.Get(pokedex)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
    body, _ := ioutil.ReadAll(response.Body)
    ss := string(body)
    split := strings.Split(ss, "<tr")[3+from:3+to]
    lastnum := from
    for _, str := range split {
        newlinesplit := strings.Split(str, "\n")
        var num, ignore int
        fmt.Sscanf(newlinesplit[1], `<td rowspan="%d" style="font-family:monospace,monospace">#%d`, &ignore, &num)
        if num == 0 || lastnum > num {
            continue
        }
        if num > to {
            break
        }
        lastnum = num
        //<td><a href="/wiki/Sigilyph_(Pok%C3%A9mon)" title="Sigilyph"><img alt="Sigilyph" src="//archives.bulbagarden.net/media/upload/thumb/a/a5/0561Sigilyph.png/70px-0561Sigilyph.png" decoding="async" loading="lazy" width="70" height="70
        var href, name, altname, imgsrc string
        fmt.Sscanf(newlinesplit[2], `<td><a href=%q title=%q><img alt=%q src=%q decoding=`, &href, &name, &altname, &imgsrc)
        if name == "" {
            continue
        }
        fmt.Printf("%d: %s\n", num, name)
	    resp, err := http.Get("https://bulbapedia.bulbagarden.net" + href)
	    if err != nil {
		    panic(err)
	    }
	    defer resp.Body.Close()
        b, _ := ioutil.ReadAll(resp.Body)
        s := string(b)
        spritesplit := strings.Split(s, `>Sprites</span></h3>`)
        s = spritesplit[1]
        bottom := strings.Split(s, `Pokémon HeartGold and SoulSilver Versions"><span style="color:#000;">SoulSilver</span></a>`)[1]
        lines := strings.Split(bottom, "/th>")[1:3]
        front := true
        name = strings.ToLower(name)
        for _, line := range lines {
            var href2 string
            fmt.Sscanf(strings.Split(line, "src")[1], "=%q", &href2)
            if !front {
                name += "_back"
            }
            downloadFile("https:"+href2, "img/"+name+".png")
            front = false
        }
    }
}

func downloadImage(category, name string) error {
    return downloadFile(baseURL + category + "/" + name, "img/" + name)
}

func fileExists(name string) (bool, error) {
    _, err := os.Stat(name)
    if err == nil {
        return true, nil
    }
    if errors.Is(err, os.ErrNotExist) {
        return false, nil
    }
    return false, err
}

func downloadFile(URL, fileName string) error {
    ok, err := fileExists(fileName)
    if err != nil {
        return err
    }
    if ok {
        return nil
    }
	//Get the response bytes from the url
    fmt.Printf("GET: %s \n", URL)
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
