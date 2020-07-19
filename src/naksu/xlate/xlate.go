package xlate

var currentLanguage string
var xlateStrings map[string]string

// SetLanguage sets active language
func SetLanguage(newLanguage string) {
	currentLanguage = newLanguage

	if currentLanguage == "fi" {
		xlateStrings = map[string]string{
			// Main window, buttons
			"Start Exam Server":                 "Käynnistä koetilan palvelin",
			"Start %s":                          "Käynnistä %s",
			"Abitti Exam":                       "Abitti-koe",
			"Abitti Exam (%s)":                  "Abitti-koe (%s)",
			"Matriculation Exam":                "Yo-koe",
			"Remove Exams":                      "Poista kokeet",
			"Remove Server":                     "Poista palvelin",
			"Make Exam Server Backup":           "Tee palvelimesta varmuuskopio",
			"Send logs to Abitti support":       "Lähetä lokitiedot Abitti-tukeen",
			"Open virtual USB stick (ktp-jako)": "Avaa virtuaalinen siirtotikku (ktp-jako)",

			// Main window, checkboxes
			"Show management features": "Näytä hallintaominaisuudet",

			// Main window, network status
			"Network status: ":                     "Verkon tila: ",
			"No network connection":                "Ei verkkoyhteyttä",
			"Network speed is too low (%d Mbit/s)": "Verkon nopeus ei riitä (%d Mbit/s)",
			"OK":                                   "OK",
			"Wireless connection":                  "Langaton yhteys",

			// Main window, log delivery statys
			"Sending logs":           "Lokitietoja lähetetään",
			"Sending logs: %d %%":    "Lokitietoja lähetetään: %d %%",
			"Error sending logs: %s": "Virhe lokitietojen lähetyksessä: %s",
			"Logs sent!":             "Lokitiedot lähetetty!",
			"Cannot send logs because there is no Internet connection. Logs are in a zip archive in the ktp folder.": "Lokitietoja ei voi lähettää, koska yhteys Internetiin ei toimi. Lokitiedot löytyvät zip-pakettina ktp-kansiosta.",
			"Copying logs: %s":       "Lokitietoja kopioidaan: %s",
			"0 % (this can take a while...)": "0 % (tässä voi mennä hetki...)",
			"Done copying":           "Lokitiedot kopioitu",
			"Zipping logs: %d %%":    "Lokitietoja pakataan: %d %%",
			"Done zipping":           "Lokitiedot pakattu",
			"Error zipping logs: %s": "Virhe lokitietojen pakkaamisessa: %s",
			"Copy to clipboard":      "Kopioi leikepöydälle",

			// Main window, other
			"Current version: %s":             "Asennettu versio: %s",
			"Update available: %s":            "Päivitys saatavilla: %s",
			"Network device:":                 "Verkkolaite:",
			"Select in terminal":              "Valitse pääteikkunassa",
			"Server networking hardware:":     "Palvelimen verkkolaite:",
			"Install/update server for:":      "Asenna tai päivitä palvelin:",
			"DANGER! Annihilate your server:": "VAARA! Palvelimen tuhoaminen:",
			"Naksu could not check for updates as there is no network connection.":                       "Naksu ei voinut etsiä uusia versioita, koska verkkoyhteys puuttuu.",
			"Naksu update failed. Maybe you don't have network connection?\n\nError: %s":                 "Naksun päivitys epäonnistui. Ehkä sinulla ei ole juuri nyt verkkoyhteyttä?\n\nVirhe: %s",
			"Did not get a path for a new Vagrantfile":                                                   "Uuden Vagrantfile-tiedoston sijainti on annettava",
			"Could not execute vagrant. Are you sure you have installed HashiCorp Vagrant?":              "Ohjelman Vagrant käynnistys epäonnistui. Oletko varma, että koneeseen on asennettu HashiCorp Vagrant?",
			"Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?":           "Ohjelman VBoxManage käynnistys epäonnistui. Oletko varma, että koneeseen on asennettu Oracle VirtualBox?",
			"Please turn Windows Hypervisor off as it may cause problems.":                               "Ole hyvä ja kytke Windows Hypervisor pois päältä, koska se voi aiheuttaa ongelmia.",
			"It appears your CPU does not support hardware virtualisation (VT-x or AMD-V).":              "Näyttää siltä, että prosessorisi ei tue laitetason virtualisointia (VT-x tai AMD-V).",
			"Hardware virtualisation (VT-x or AMD-V) is disabled. Please enable it before continuing.":   "Laitetason virtualisointi (VT-x tai AMD-V) on kytketty pois päältä. Kytke se päälle koneen asetuksista.",
			"Your home directory path (%s) contains characters which may cause problems to Vagrant.":     "Kotihakemistosi (%s) polku sisältää merkkejä, jotka voivat aiheuttaa ongelmia Vagrantille.",
			"Please place the new Exam Vagrantfile to another location (e.g. desktop or home directory)": "Sijoita yo-kokeen Vagrantfile johonkin toiseen paikkaan (esim. työpöydälle tai kotihakemistoon)",
			"Your free disk size is getting low (%s).":                                                   "Levytilasi on loppumassa (jäljellä %s).",
			"Start by installing a server: Show management features":                                     "Aloita asentamalla palvelin: Näytä hallintaominaisuudet",
			"You are starting Matriculation Examination server with an Internet connection.":             "Olet käynnistämässä yo-palvelinta, vaikka palvelimella on internet-yhteys.",

			"You have not set network device. Follow terminal for device selection menu.":                              "Et ole valinnut verkkolaitetta. Seuraa pääteikkunaa mahdollista verkkolaitteen valintaa varten.",
			"You have selected network device '%s' which is not available. Follow terminal for device selection menu.": "Olet valinnut verkkolaitteen '%s', joka ei ole käytettävissä. Seuraa pääteikkunaa mahdollista verkkolaitteen valintaa varten.",

			// Backup dialog
			"naksu: SaveTo":             "naksu: Tallennuspaikka",
			"Please select target path": "Valitse tallennuspaikka",
			"Save":                      "Tallenna",
			"Cancel":                    "Peruuta",

			// Log delivery dialog
			"naksu: Send Logs":             "naksu: Lähetä lokitiedot",
			"Filename for Abitti support:": "Tiedostonimi Abitti-tuelle:",
			"Wait...":						"Odota...",
			"Close":                        "Sulje",

			// Destroy dialog
			"naksu: Remove Exams": "naksu: Poista kokeet",
			"Remove Exams restores server to its initial status.":                   "Kokeiden poistaminen palauttaa palvelimen alkutilaan.",
			"Exams, responses and logs in the server will be irreversibly deleted.": "Kokeet, suoritukset ja lokitiedot poistetaan peruuttamattomasti.",
			"It is recommended to back up your server before removing exams.":       "On suositeltavaa ottaa palvelimesta varmuuskopio ennen kokeiden poistamista.",
			"Do you wish to remove all exams?":                                      "Haluatko poistaa kaikki kokeet?",
			"Yes, Remove":                                                           "Kyllä, poista",
			// Already defined: "Cancel"
			"Removing exams. This takes a while.": "Kokeita poistetaan. Odota hetki.",
			"Exams were removed successfully.":    "Kokeiden poistaminen onnistui.",
			"Failed to remove exams.":             "Kokeiden poistaminen epäonnistui.",

			// Remove dialog
			"naksu: Remove Server": "naksu: Poista palvelin",
			"Removing server destroys it and all downloaded disk images.": "Palvelimen poistaminen tuhoaa sen ja kaikki ladatut levynkuvat.",
			// Already defined: "Exams, responses and logs in the server will be irreversibly deleted."
			"It is recommended to back up your server before removing server.": "On suositeltavaa ottaa palvelimesta varmuuskopio ennen poistamista.",
			"Do you wish to remove the server?":                                "Halutko poistaa palvelimen?",
			// Already defined: "Yes, Remove"
			// Already defined: "Cancel"
			"Error while removing server: %v": "Palvelimen poistaminen epäonnistui: %v",
			"Server was removed succesfully.": "Palvelimen poistaminen onnistui.",

			// mebroutines
			"command failed: %s":       "komento epäonnistui: %s",
			"Failed to execute %s":     "Komennon suorittaminen epäonnistui: %s",
			"Failed to execute %s: %v": "Komennon suorittaminen epäonnistui: %s (%v)",
			"Could not chdir to %s":    "Hakemistoon %s siirtyminen epäonnistui",
			"Server failed to start. This is typical in Windows after an update. Please try again to start the server.": "Palvelimen käynnistys epäonnistui. Tämä on tyypillista Windows-koneissa päivityksen jälkeen. Yritä käynnistää palvelin uudelleen.",
			"Error":   "Virhe",
			"Warning": "Varoitus",
			"Info":    "Tiedoksi",

			// backup
			"File %s already exists":                                "Tiedosto %s on jo olemassa",
			"Backup has been made to %s":                            "Varmuuskopio on talletettu tiedostoon %s",
			"Could not get vagrantbox ID: %d":                       "Vagrantboxin ID:tä ei voitu lukea: %d",
			"Could not back up disk %s to %s":                       "Varmuuskopion ottaminen levystä %s tiedostoon %s epäonnistui",
			"Could not write backup file %s. Try another location.": "Varmuuskopion kirjoittaminen tiedostoon %s epäonnistui. Kokeile toista tallennuspaikkaa.",
			"Backup failed: %v":                                     "Varmuuskopiointi epäonnistui: %v",
			"The backup file is too large for a FAT32 filesystem. Please reformat the backup disk as exFAT.": "Varmuuskopio on liian suuri talletettavaksi FAT32-tiedostojärjestelmäään. Alusta varmuuskopiolevy uudelleen exFAT-tiedostojärjestelmällä.",

			"Checking existing file...":      "Etsin olemassaolevaa tiedostoa...",
			"Checking backup path...":        "Tarkistan varmistushakemistoa...",
			"Getting vagrantbox ID...":       "Haen vagrantbox ID:tä...",
			"Getting disk UUID...":           "Haen levyn UUID:tä...",
			"Please wait, writing backup...": "Odota, kirjoitan varmuuskopiota...",
			"Detaching backup disk image...": "Irrotan varmuuskopion levynkuvan...",
			"Backup done: %s":                "Varmuuskopio valmis: %s",

			// backup, getmediapath
			"Home directory":    "Kotihakemisto",
			"Temporary files":   "Tilapäishakemisto",
			"Profile directory": "Profiilihakemisto",
			"Desktop":           "Työpöytä",

			// install
			"Could not update Abitti stickless server. Please check your network connection.": "virtuaalisen Abitti-palvelimen päivitys epäonnistui. Tarkista verkkoyhteytesi.",
			"Could not change to vagrant directory ~/ktp":                                     "Vagrant-hakemistoon ~/ktp siirtyminen epäonnistui",
			"Error while copying new Vagrantfile: %d":                                         "Uuden Vagrantfile-tiedoston kopiointi epäonnistui: %d",
			"Could not create directory: %v":                                                  "Hakemiston luominen epäonnistui: %v",
			"Failed to delete %s":                                                             "Tiedoston %s poistaminen epäonnistui",
			"Failed to rename %s to %s":                                                       "Tiedoston %s nimeäminen tiedostoksi %s epäonnistui",
			"Failed to create file %s":                                                        "Tiedoston %s luominen epäonnistui",
			"Failed to retrieve %s":                                                           "Sijainnista %s lataaminen epäonnistui",
			"Could not copy body from %s to %s":                                               "Sisällön %s kopioint sijaintiin %s epäonnistui",

			// start
			// Already defined: "Could not change to vagrant directory ~/ktp"

			// boxversion
			"Abitti server":      "Abitti-palvelin",
			"Matric Exam server": "Yo-palvelin",
		}
	} else if currentLanguage == "sv" {
		xlateStrings = map[string]string{
			// Main window, buttons
			"Start Exam Server":                 "Starta provlokalsserver",
			"Start %s":                          "Starta %s",
			"Abitti Exam":                       "Abitti-prov",
			"Abitti Exam (%s)":                  "Abitti-prov (%s)",
			"Matriculation Exam":                "Studentprovet",
			"Remove Exams":                      "Avlägsna proven",
			"Remove Server":                     "Avlägsna servern",
			"Make Exam Server Backup":           "Säkerhetskopiera servern",
			"Send logs to Abitti support":       "Skicka logguppgifterna till Abitti-stödet",
			"Open virtual USB stick (ktp-jako)": "Öppna den virtuellaöverföringsstickan (ktp-jako)",

			// Main window, checkboxes
			"Show management features": "Visa hanteringsegenskaper",

			// Main window, network status
			"Network status :":                     "Nätverksstatus: ",
			"No network connection":                "Inget nätverk",
			"Network speed is too low (%d Mbit/s)": "Hastigheten är för låg (%d Mbit/s)",
			"OK":                                   "OK",
			"Wireless connection":                  "Trådlös anslutning",

			// Main window, log delivery statys
			"Sending logs":           "Skickar logguppgifter",
			"Sending logs: %d %%":    "Skickar logguppgifter: %d %%",
			"Error sending logs: %s": "Fel i skickande av logguppgifter: %s",
			"Logs sent!":             "Logguppgifterna har skickats!",
			"Cannot send logs because there is no Internet connection. Logs are in a zip archive in the ktp folder.": "Logguppgifterna kan inte skickas, eftersom anslutningen till Internet inte fungerar. Logguppgifterna finns sparade som en zip-fil i ktp-mappen.",
			"Copying logs: %s":       "Kopierar logguppgifter: %s",
			"0 % (this can take a while...)": "0 % (kan ta ett tag...)",
			"Done copying":           "Logguppgifterna är kopierade",
			"Zipping logs: %d %%":    "Komprimerar logguppgifter: %d %%",
			"Done zipping":           "Logguppgifterna är komprimerade",
			"Error zipping logs: %s": "Fel i komprimering av logguppgifter: %s",
			"Copy to clipboard":      "Kopiera till urklipp",

			// Main window, other
			"Current version: %s":             "Installerad version: %s",
			"Update available: %s":            "Uppdatering tillgänglig: %s",
			"Network device:":                 "Nätverksanordning:",
			"Select in terminal":              "Välj i terminalen",
			"Server networking hardware:":     "Servernätverkshårdvara:",
			"Install/update server for:":      "Installera eller uppdatera server för:",
			"DANGER! Annihilate your server:": "FARA! Förstörning av servern:",
			"Naksu could not check for updates as there is no network connection.":                       "Naksu kunde inte söka nya uppdateringar för att nätförbindelse saknades.",
			"Naksu update failed. Maybe you don't have network connection?\n\nError: %s":                 "Uppdateringen av Naksu misslyckades. Du saknar möjligtvis nätförbindelse för tillfället?\n\nFel: %s",
			"Did not get a path for a new Vagrantfile":                                                   "Ge sökvägen för den nya Vagrantfile-filen",
			"Could not execute vagrant. Are you sure you have installed HashiCorp Vagrant?":              "Utförandet av programmet Vagrant misslyckades. Är du säker, att HashiCorp Vagrant har installerats på datorn?",
			"Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?":           "Utförandet av programmet VBoxManage misslyckades. Är du säker, att Oracle VirtualBox har installerats på datorn?",
			"Please turn Windows Hypervisor off as it may cause problems.":                               "Vänligen stäng Windows Hypervisor eftersom den kan orsaka problem.",
			"It appears your CPU does not support hardware virtualisation (VT-x or AMD-V).":              "Det verkar som om din processor inte stöder virtualisering av hårdvara (VT-x eller AMD-V).",
			"Hardware virtualisation (VT-x or AMD-V) is disabled. Please enable it before continuing.":   "Virtualiseringen av hårdvaran (VT-x eller AMD-V) är avstängd. Vänligen aktivera den innan du fortsätter.",
			"Your home directory path (%s) contains characters which may cause problems to Vagrant.":     "Sökvägen till din hemkatalog (%s) innehåller tecken, som orsakar problem för Vagrant.",
			"Please place the new Exam Vagrantfile to another location (e.g. desktop or home directory)": "Placera Vagrantfile-filen för studentexamen på ett annat ställe (t.ex. på skrivbordet eller i hemkatalogen).",
			"Your free disk size is getting low (%s).":                                                   "Ditt diskutrymme börjar ta slut (kvar %s).",
			"Start by installing a server: Show management features":                                     "Börja med att installera servern: Visa hanteringsegenskaper",
			"You are starting Matriculation Examination server with an Internet connection.":             "Då håller på att starta provlokalsservern för studentexamen, fast servern har internetförbindelse.",

			"You have not set network device. Follow terminal for device selection menu.":                              "Du har inte valt nätverksanordning. Följ terminalfönstret för eventuellt val av nätverksanordning.",
			"You have selected network device '%s' which is not available. Follow terminal for device selection menu.": "Du har valt nätverksanordningen '%s', som inte är tillgänglig. Följ terminalfönstret för eventuellt val av nätverksanordning.",

			// Backup dialog
			"naksu: SaveTo":             "naksu: Spara till ",
			"Please select target path": "Välj sökväg",
			"Save":                      "Spara",
			"Cancel":                    "Avbryt",

			// Log delivery dialog
			"naksu: Send Logs":             "naksu: Skicka logguppgifterna",
			"Filename for Abitti support:": "Filnamn för Abitti-stödet",
			"Wait...":						"Vänta...",
			"Close":                        "Stäng",

			// Destroy dialog
			"naksu: Remove Exams": "naksu: Avlägsna proven",
			"Remove Exams restores server to its initial status.":                   "Avlägsnandet av proven återställer servern till sitt ursprungsläge.",
			"Exams, responses and logs in the server will be irreversibly deleted.": "Alla prov, loggfiler och svar på servern avlägsnas oåterkalleligt.",
			"It is recommended to back up your server before removing exams.":       "Det är rekommenderat att ta en säkerhetskopia av servern före proven avlägsnas.",
			"Do you wish to remove all exams?":                                      "Vill du avlägsna alla prov?",
			"Yes, Remove":                                                           "Ja, avlägsna",
			// Already defined: "Cancel"
			"Removing exams. This takes a while.": "Proven avlägsnas. Vänta en stund. ",
			"Exams were removed successfully.":    "Avlägsning av proven lyckades.",
			"Failed to remove exams.":             "Avlägsning av proven misslyckades.",

			// Remove dialog
			"naksu: Remove Server": "naksu: Avlägsna servern",
			"Removing server destroys it and all downloaded disk images.": "Avlägsnandet av servern förstör den och alla nerladdade skivavbilder.",
			// Already defined: "Exams, responses and logs in the server will be irreversibly deleted."
			"It is recommended to back up your server before removing server.": "Det är rekommenderat att ta en säkerhetskopia av servern före den avlägsnas.",
			"Do you wish to remove the server?":                                "Vill du avlägsna servern?",
			// Already defined: "Yes, Remove"
			// Already defined: "Cancel"
			"Error while removing server: %v": "Avlägsning av servern misslyckades: %v",
			"Server was removed succesfully.": "Avlägsning av servern lyckades.",

			// mebroutines
			"command failed: %s":       "Komandot misslyckades: %s",
			"Failed to execute %s":     "Utförning av komandot misslyckades: %s",
			"Failed to execute %s: %v": "Utförning av komandot misslyckades: %s (%v)",
			"Could not chdir to %s":    "Förflyttning till katalogen %s misslyckades",
			"Server failed to start. This is typical in Windows after an update. Please try again to start the server.": "Startandet av servern misslyckades. Detta är typiskt i Windows efter en uppdatering. Vänligen försök igen för att starta servern.",
			"Error":   "Fel",
			"Warning": "Varning",
			"Info":    "För information",

			// backup
			"File %s already exists":                                "Filen %s existerar redan",
			"Backup has been made to %s":                            "Säkerhetskopian har sparats i filen %s",
			"Could not get vagrantbox ID: %d":                       "Det gick inte att läsa ID:n på Vagrantboxen: %d",
			"Could not back up disk %s to %s":                       "Säkerhetskopieringen av skivan %s i filen %s misslyckades",
			"Could not write backup file %s. Try another location.": "Det gick inte att säkerhetskopiera till filen %s. Pröva att spara filen på ett annat ställe.",
			"Backup failed: %v":                                     "Säkerhetskopieringen misslyckades: %v",
			"The backup file is too large for a FAT32 filesystem. Please reformat the backup disk as exFAT.": "Säkerhetskopian är för stor för ett FAT32-filsystem. Vänligen formatera minnepinnen eller skivan som exFAT.",

			"Checking existing file...":      "Granskar existerande fil...",
			"Checking backup path...":        "Granskar säkerhetskopians katalog...",
			"Getting vagrantbox ID...":       "Söker vagrantbox ID...",
			"Getting disk UUID...":           "Söker skivans UUID...",
			"Please wait, writing backup...": "Vänligen vänta, skriver säkerhetskopian...",
			"Detaching backup disk image...": "Lösgör säkerhetskopians skivavbild...",
			"Backup done: %s":                "Säkerhetskopian färdig: %s",

			// backup, getmediapath
			"Home directory":    "Hemkatalog",
			"Temporary files":   "Tillfällig katalog",
			"Profile directory": "Profilkatalog",
			"Desktop":           "Skrivbord",

			// install
			"Could not update Abitti stickless server. Please check your network connection.": "Det gick inte att uppdatera den virtuella Abitti-servern. Kontrollera din nätförbindelse.",
			"Could not change to vagrant directory ~/ktp":                                     "Förflyttningen till Vagrant-katalogen ~/ktp misslyckades",
			"Error while copying new Vagrantfile: %d":                                         "Kopieringen av en ny Vagrantfile-fil misslyckades: %d",
			"Could not create directory: %v":                                                  "Det gick inte att skapa katalogen: %v",
			"Failed to delete %s":                                                             "Det gick inte att avlägsna filen %s",
			"Failed to rename %s to %s":                                                       "Det gick inte att namnge filen %s som %s",
			"Failed to create file %s":                                                        "Det gick inte att skapa filen %s",
			"Failed to retrieve %s":                                                           "Det gick inte att ladda ner från sökvägen %s",
			"Could not copy body from %s to %s":                                               "Det gick inte att kopiera från sökvägen %s till %s",

			// start
			// Already defined: "Could not change to vagrant directory ~/ktp"

			// boxversion
			"Abitti server":      "Abitti-server",
			"Matric Exam server": "Examensserver",
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
