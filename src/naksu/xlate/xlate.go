package xlate

var currentLanguage string
var xlateStrings map[string]string

// SetLanguage sets active language
func SetLanguage(newLanguage string) {
	currentLanguage = newLanguage

	if currentLanguage == "fi" {
		xlateStrings = map[string]string{
			// Main window, buttons
			"Start Stickless Exam Server":                           "Käynnistä virtuaalinen koetilan palvelin",
			"Install or update Abitti Stickless Exam Server":        "Asenna tai päivitä Abitin virtuaalinen palvelin",
			"Install or update Stickless Matriculation Exam Server": "Asenna tai päivitä yo-kokeen virtuaalinen palvelin",
			"Make Stickless Exam Server Backup":                     "Tee virtuaalisesta palvelimesta varmuuskopio",
			"Open virtual USB stick (ktp-jako)":                     "Avaa virtuaalinen siirtotikku (ktp-jako)",

			// Main window, checkboxes
			"Show management features": "Näytä hallintaominaisuudet",

			// Main window, other
			"Current version: %s": "Asennettu versio: %s",
			"Naksu update failed. Maybe you don't have network connection?\n\nError: %s":                 "Naksun päivitys epäonnistui. Ehkä sinulla ei ole juuri nyt verkkoyhteyttä?\n\nVirhe: %s",
			"Did not get a path for a new Vagrantfile":                                                   "Uuden Vagrantfile-tiedoston sijainti on annettava",
			"Could not execute vagrant. Are you sure you have installed HashiCorp Vagrant?":              "Ohjelman Vagrant käynnistys epäonnistui. Oletko varma, että koneeseen on asennettu HashiCorp Vagrant?",
			"Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?":           "Ohjelman VBoxManage käynnistys epäonnistui. Oletko varma, että koneeseen on asennettu Oracle VirtualBox?",
			"Your home directory path (%s) contains characters which may cause problems to Vagrant.":     "Kotihakemistosi (%s) polku sisältää merkkejä, jotka voivat aiheuttaa ongelmia Vagrantille.",
			"Please place the new Exam Vagrantfile to another location (e.g. desktop or home directory)": "Sijoita yo-kokeen Vagrantfile johonkin toiseen paikkaan (esim. työpöydälle tai kotihakemistoon)",

			// Backup dialog
			"naksu: SaveTo":             "naksu: Tallennuspaikka",
			"Please select target path": "Valitse tallennuspaikka",
			"Save":                      "Tallenna",
			"Cancel":                    "Peruuta",

			// mebroutines
			"Abitti server":         "Abitti-palvelin",
			"Matric Exam server":    "Yo-palvelin",
			"command failed: %s":    "komento epäonnistui: %s",
			"Failed to execute %s":  "Komennon suorittaminen epäonnistui: %s",
			"Could not chdir to %s": "Hakemistoon %s siirtyminen epäonnistui",
			"Server failed to start. This is typical in Windows after an update. Please try again to start the server.": "Palvelimen käynnistys epäonnistui. Tämä on tyypillista Windows-koneissa päivityksen jälkeen. Yritä käynnistää palvelin uudelleen.",
			"Error":   "Virhe",
			"Warning": "Varoitus",
			"Info":    "Tiedoksi",

			// backup
			"File %s already exists":                                "Tiedosto %s on jo olemassa",
			"Backup has been made to %s":                            "Varmuuskopio on talletettu tiedostoon %s",
			"Could not get vagrantbox ID: %d":                       "Vagrantboxin ID:tä ei voitu lukea: %d",
			"Could not make backup: failed to get disk UUID":        "Varmuuskopion ottaminen epäonnistui: levyn UUID:tä ei löytynyt",
			"Could not back up disk %s to %s":                       "Varmuuskopion ottaminen levystä %s tiedostoon %s epäonnistui",
			"Could not write backup file %s. Try another location.": "Varmuuskopion kirjoittaminen tiedostoon %s epäonnistui. Kokeile toista tallennuspaikkaa.",
			"Backup failed.":                                        "Varmuuskopiointi epäonnistui.",

			// backup, getmediapath
			"Home directory":    "Kotihakemisto",
			"Temporary files":   "Tilapäishakemisto",
			"Profile directory": "Profiilihakemisto",
			"Desktop":           "Työpöytä",

			// install
			"Could not update Abitti stickless server. Please check your network connection.": "virtuaalisen Abitti-palvelimen päivitys epäonnistui. Tarkista verkkoyhteytesi.",
			"Could not change to vagrant directory ~/ktp":                                     "Vagrant-hakemistoon ~/ktp siirtyminen epäonnistui",
			"Error while copying new Vagrantfile: %d":                                         "Uuden Vagrantfile-tiedoston kopiointi epäonnistui: %d",
			"Could not create ~/ktp to %s":                                                    "Hakemiston ~/ktp luominen sijaintiin %s epäonnistui",
			"Could not create ~/ktp-jako to %s":                                               "Hakemiston ~/ktp-jako luominen sijaintiin %s epäonnistui",
			"Failed to delete %s":                                                             "Tiedoston %s poistaminen epäonnistui",
			"Failed to rename %s to %s":                                                       "Tiedoston %s nimeäminen tiedostoksi %s epäonnistui",
			"Failed to create file %s":                                                        "Tiedoston %s luominen epäonnistui",
			"Failed to retrieve %s":                                                           "Sijainnista %s lataaminen epäonnistui",
			"Could not copy body from %s to %s":                                               "Sisällön %s kopioint sijaintiin %s epäonnistui",

			// start
			// Already defined: "Could not change to vagrant directory ~/ktp"
		}
	} else if currentLanguage == "sv" {
		xlateStrings = map[string]string{
			// Main window, buttons
			"Start Stickless Exam Server":                           "Starta virtuell provlokalsserver",
			"Install or update Abitti Stickless Exam Server":        "Installera eller uppdatera virtuell server för Abitti",
			"Install or update Stickless Matriculation Exam Server": "Installera eller uppdatera virtuell server för studentexamen",
			"Make Stickless Exam Server Backup":                     "Säkerhetskopiera den virtuella servern",
			"Open virtual USB stick (ktp-jako)":                     "Öppna den virtuellaöverföringsstickan (ktp-jako)",

			// Main window, checkboxes
			"Show management features": "Visa hanteringsegenskaper",

			// Main window, other
			"Current version: %s": "Installerad version: %s",
			"Naksu update failed. Maybe you don't have network connection?\n\nError: %s":                 "Uppdateringen av Naksu misslyckades. Du saknar möjligtvis nätförbindelse för tillfället?\n\nFel: %s",
			"Did not get a path for a new Vagrantfile":                                                   "Ge sökvägen för den nya Vagrantfile-filen",
			"Could not execute vagrant. Are you sure you have installed HashiCorp Vagrant?":              "Utförandet av programmet Vagrant misslyckades. Är du säker, att HashiCorp Vagrant har installerats på datorn?",
			"Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?":           "Utförandet av programmet VBoxManage misslyckades. Är du säker, att Oracle VirtualBox har installerats på datorn?",
			"Your home directory path (%s) contains characters which may cause problems to Vagrant.":     "Sökvägen till din hemkatalog (%s) innehåller tecken, som orsakar problem för Vagrant.",
			"Please place the new Exam Vagrantfile to another location (e.g. desktop or home directory)": "Placera Vagrantfile-filen för studentexamen på ett annat ställe (t.ex. på skrivbordet eller i hemkatalogen).",

			// Backup dialog
			"naksu: SaveTo":             "naksu: Spara till ",
			"Please select target path": "Välj sökväg",
			"Save":                      "Spara",
			"Cancel":                    "Avbryt",

			// mebroutines
			"Abitti server":         "Abitti-server",
			"Matric Exam server":    "Examensserver",
			"command failed: %s":    "Komandot misslyckades: %s",
			"Failed to execute %s":  "Utförning av komandot misslyckades: %s",
			"Could not chdir to %s": "Förflyttning till katalogen %s misslyckades",
			"Server failed to start. This is typical in Windows after an update. Please try again to start the server.": "Startandet av servern misslyckades. Detta är typiskt i Windows efter en uppdatering. Vänligen försök igen för att starta servern.",
			"Error":   "Fel",
			"Warning": "Varning",
			"Info":    "För information",

			// backup
			"File %s already exists":                                "Filen %s existerar redan",
			"Backup has been made to %s":                            "Säkerhetskopian har sparats i filen %s",
			"Could not get vagrantbox ID: %d":                       "Det gick inte att läsa ID:n på Vagrantboxen: %d",
			"Could not make backup: failed to get disk UUID":        "Säkerhetskopieringen misslyckades: skivans UUID hittades inte",
			"Could not back up disk %s to %s":                       "Säkerhetskopieringen av skivan %s i filen %s misslyckades",
			"Could not write backup file %s. Try another location.": "Det gick inte att säkerhetskopiera till filen %s. Pröva att spara filen på ett annat ställe.",
			"Backup failed.":                                        "Säkerhetskopieringen misslyckades.",

			// backup, getmediapath
			"Home directory":    "Hemkatalog",
			"Temporary files":   "Tillfällig katalog",
			"Profile directory": "Profilkatalog",
			"Desktop":           "Skrivbord",

			// install
			"Could not update Abitti stickless server. Please check your network connection.": "Det gick inte att uppdatera den virtuella Abitti-servern. Kontrollera din nätförbindelse.",
			"Could not change to vagrant directory ~/ktp":                                     "Förflyttningen till Vagrant-katalogen ~/ktp misslyckades",
			"Error while copying new Vagrantfile: %d":                                         "Kopieringen av en ny Vagrantfile-fil misslyckades: %d",
			"Could not create ~/ktp to %s":                                                    "Det gick inte att skapa katalogen ~/ktp i sökvägen %s",
			"Could not create ~/ktp-jako to %s":                                               "Det gick inte att skapa katalogen ~/ktp-jako i sökvägen %s",
			"Failed to delete %s":                                                             "Det gick inte att avlägsna filen %s",
			"Failed to rename %s to %s":                                                       "Det gick inte att namnge filen %s som %s",
			"Failed to create file %s":                                                        "Det gick inte att skapa filen %s",
			"Failed to retrieve %s":                                                           "Det gick inte att ladda ner från sökvägen %s",
			"Could not copy body from %s to %s":                                               "Det gick inte att kopiera från sökvägen %s till %s",

			// start
			// Already defined: "Could not change to vagrant directory ~/ktp"
		}
	} else {
		xlateStrings = nil
	}
}

// Get returns translated string for given key
func Get(key string) string {
	if xlateStrings == nil {
		return key
	}

	newString := xlateStrings[key]

	if newString == "" {
		return key
	}

	return newString
}
