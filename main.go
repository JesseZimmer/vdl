package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bodgit/sevenzip"
	"github.com/eiannone/keyboard"
	"github.com/inancgumus/screen"
)

func gameList(systemUrl string) map[string]string {

	gamesList := make(map[string]string)
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	res, err := c.Get(systemUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Find the table with class "vault_table"
	doc.Find("table.rounded").Each(func(b int, tableHtml *goquery.Selection) {
		// Find all rows in the table
		tableHtml.Find("tr").Each(func(j int, rowHtml *goquery.Selection) {
			// Find all cells in the row
			gameTableRow := rowHtml.Find("td").First()
			// Find all <a> elements in the cell
			gameAtrib := gameTableRow.Find("a").First()
			// Get the text and href attribute of the <a> element

			gameName := gameAtrib.Text()
			gameHref, _ := gameAtrib.Attr("href")
			vaultSplit := strings.Split(gameHref, "/")
			vaultID := vaultSplit[len(vaultSplit)-1]

			gamesList[gameName] = vaultID

		})
	})
	return gamesList
}

func extractArchive(archivePath string) error {

	r, err := sevenzip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".bin") || strings.HasSuffix(f.Name, ".cue") {

			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			outFile, err := os.Create("/mnt/SDCARD/Roms/PS/" + strings.Split(f.Name, "/")[1])
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}

			fmt.Printf("Extracted: %s\n", "/mnt/SDCARD/Roms/PS/"+f.Name)
		}
	}

	return nil
}

func parseRom(vaultId string) (mediaId string, romFolder string) {

	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := c.Get("https://vimm.net/vault/" + vaultId)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	mediaId = doc.Find("input[name='mediaId']").AttrOr("value", "")

	img := doc.Find("input[name='system']")
	if img.Length() > 0 {
		console, _ := img.Attr("value")

		if console == "GB" {
			romFolder = "GB/"
		} else if console == "GBC" {
			romFolder = "GBC/"
		} else if console == "GBA" {
			romFolder = "GBA/"
		} else if console == "DS" {
			romFolder = "NDS/"
		} else if console == "Atari2600" {
			romFolder = "ATARI/"
		} else if console == "Atari5200" {
			romFolder = "FIFTYTWOHUNDRED/"
		} else if console == "NES" {
			romFolder = "FC/"
		} else if console == "SMS" {
			romFolder = "MS/"
		} else if console == "Atari7800" {
			romFolder = "SEVENTYEIGHTHUNDRED/"
		} else if console == "Genesis" {
			romFolder = "MD/"
		} else if console == "SNES" {
			romFolder = "SFC/"
		} else if console == "32X" {
			romFolder = "THIRTYTWOX/"
		} else if console == "PS1" {
			romFolder = "PS/"
		} else if console == "Lynx" {
			romFolder = "LYNX/"
		} else if console == "GG" {
			romFolder = "GG/"
		} else if console == "VB" {
			romFolder = "VB/"
		} else {
			fmt.Println("No console match - saving ROM to /mnt/SDCARD/Roms/.")
		}

	} else {
		fmt.Println("No console found - saving ROM to /mnt/SDCARD/Roms/.")
	}

	return mediaId, romFolder
}

func downloadRom(filepath string, romUrl string, downloadUrl string) (err error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	fmt.Println("Downloading... this may take some time for larger game files... (PS1, NDS)")

	req, err := http.NewRequest("GET", downloadUrl, nil)
	req.Header.Add("Referer", romUrl)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0")
	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Error: %s", resp.Status)
	}

	gameName := strings.Split(resp.Header.Get("Content-Disposition"), "\"")[1]

	filepath = filepath + gameName

	rom, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer rom.Close()

	_, err = io.Copy(rom, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Download complete -", filepath)

	if strings.Contains(filepath, "/PS/") {
		fmt.Println("Extracting Archive... Please wait...")
		extractArchive(filepath)
		fmt.Println("Extraction Complete")
		defer os.Remove(filepath) // remove a single file
	}

	return nil
}

func main() {
	//SystemList
	sysList := map[int]map[string]string{
		1:  {"Atari 2600": "https://vimm.net/vault/Atari2600"},
		2:  {"Atari 5200": "https://vimm.net/vault/Atari5200"},
		3:  {"Nintendo": "https://vimm.net/vault/NES"},
		4:  {"Master System": "https://vimm.net/vault/SMS"},
		5:  {"Atari 7800": "https://vimm.net/vault/Atari7800"},
		6:  {"Genesis": "https://vimm.net/vault/Genesis"},
		7:  {"Super Nintendo": "https://vimm.net/vault/SNES"},
		8:  {"Sega 32X": "https://vimm.net/vault/32X"},
		9:  {"Saturn": "https://vimm.net/vault/Saturn"},
		10: {"PlayStation": "https://vimm.net/vault/PS1"},
		11: {"Game Boy": "https://vimm.net/vault/GB"},
		12: {"Lynx": "https://vimm.net/vault/Lynx"},
		13: {"Game Gear": "https://vimm.net/vault/GG"},
		14: {"Virtual Boy": "https://vimm.net/vault/VB"},
		15: {"Game Boy Color": "https://vimm.net/vault/GBC"},
		16: {"Game Boy Advance": "https://vimm.net/vault/GBA"},
		17: {"Nintendo DS": "https://vimm.net/vault/DS"},
	}
	screen.Clear()
	fmt.Printf("Unoffical Vimm.net Game Downloader\n\n")

	//sort list
	keys := make([]int, 0, len(sysList))
	for key := range sysList {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	i := 1
	for _, key := range keys {
		sysMap := sysList[key]
		for sysName := range sysMap {
			fmt.Printf("%d. %s\n", i, sysName)
			i += 1
		}
	}

	var sysSelection int
	var letterSelection string
	fmt.Print("\nSelect the system # : ")
	fmt.Scanln(&sysSelection)

	selectedSystem, found := sysList[sysSelection]
	if !found {
		fmt.Println("Invalid system selection")
		return
	}

	var sysUrl string

	for _, value := range selectedSystem {
		sysUrl = value
		break
	}

	fmt.Print("Enter A-Z to see a list of games: ")
	fmt.Scanln(&letterSelection)

	//fmt.Println(sysUrl + "/" + letterSelection)

	games := gameList(sysUrl + "/" + letterSelection)

	gameKeys := make([]string, 0, len(games))
	for key := range games {
		gameKeys = append(gameKeys, key)
	}
	sort.Strings(gameKeys)

	pageSize := 10
	currentPage := 0
	totalPages := len(gameKeys) / pageSize
	if len(keys)%pageSize != 0 {
		totalPages++
	}

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		screen.Clear()
		printPage(games, gameKeys, currentPage, pageSize)
		fmt.Printf("Page %d/%d - Use L2/R2 to navigate, Q/q to download game\n", currentPage+1, totalPages)

		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc || char == 'q' || char == 'Q' {
			var vaultID string
			fmt.Println()
			fmt.Printf("Enter the vault ID from the game manu above: ")
			fmt.Scanln(&vaultID)

			mediaId, romFolder := parseRom(vaultID)
			romFilepath := "/mnt/SDCARD/Roms/" + romFolder
			downloadRom(romFilepath, "https://vimm.net/vault/"+vaultID, "https://download3.vimm.net/download/?mediaId="+mediaId)
			os.Exit(1)

		}

		if key == keyboard.KeyArrowLeft {
			if currentPage > 0 {
				currentPage--
			}
		} else if key == keyboard.KeyArrowRight {
			if currentPage < totalPages-1 {
				currentPage++
			}
		}

	}

}

func printPage(items map[string]string, keys []string, currentPage, pageSize int) {
	start := currentPage * pageSize
	end := start + pageSize
	if end > len(keys) {
		end = len(keys)
	}

	for _, key := range keys[start:end] {
		fmt.Printf("%s: %s\n", key, items[key])
	}
	fmt.Println()
}
