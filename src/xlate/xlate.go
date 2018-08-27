package xlate


var current_language string
var xlate_strings map[string]string

func SetLanguage(new_language string) {
  current_language = new_language

  if current_language == "fi" {
    xlate_strings = map[string]string {
      // Main window, groups
      "Language": "Kieli",
      "Basic Functions": "Perustoiminnot",
      "Abitti": "Abitti-koe",
      "Matriculation Exam": "Ylioppilaskoe",

      // Main window, buttons
      "Start Stickless Exam Server": "Käynnistä tikuton koetilan palvelin",
      "Install or update Abitti Stickless Exam Server": "Asenna tai päivitä Abitin tikuton palvelin",
      "Install or update Stickless Matriculation Exam Server": "Asenna tai päivitä yo-kokeen tikuton palvelin",
      "Make Stickless Exam Server Backup": "Tee tikuttomasta palvelimesta varmuuskopio",
      "Exit": "Poistu",

      // Main window, other
      "Did not get a path for a new Vagrantfile": "Uuden Vagrantfile-tiedoston sijainti on annettava",
      "Could not execute vagrant. Are you sure you have installed HashiCorp Vagrant?": "Ohjelman Vagrant käynnistys epäonnistui. Oletko varma, että koneeseen on asennettu HashiCorp Vagrant?",
      "Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?": "Ohjelman VBoxManage käynnistys epäonnistui. Oletko varma, että koneeseen on asennettu Oracle VirtualBox?",

      // Backup dialog
      "naksu: SaveTo": "naksu: Tallennuspaikka",
      "Please select target path": "Valitse tallennuspaikka",
      "Save": "Tallenna",
      "Cancel": "Peruuta",

      // mebroutines
      "command failed: %s": "komento epäonnistui: %s",
      "Failed to execute %s": "Komennon suorittaminen epäonnistui: %s",
      "Could not chdir to %s": "Hakemistoon %s siirtyminen epäonnistui",
      "Error": "Virhe",
      "Warning": "Varoitus",
      "Info": "Tiedoksi",

      // backup
      "File %s already exists": "Tiedosto %s on jo olemassa",
      "Backup has been made to %s": "Varmuuskopio on talletettu tiedostoon %s",
      "Could not get vagrantbox ID: %d": "Vagrantboxin ID:tä ei voitu lukea: %d",
      "Could not make backup: failed to get disk UUID": "Varmuuskopion ottaminen epäonnistui: levyn UUID:tä ei löytynyt",
      "Could not back up disk %s to %s": "Varmuuskopion ottaminen levystä %s tiedostoon %s epäonnistui",

      // backup, getmediapath
      "Home directory": "Kotihakemisto",
      "Temporary files": "Tilapäishakemisto",
      "Profile directory": "Profiilihakemisto",

      // install
      "Could not change to vagrant directory ~/ktp": "Vagrant-hakemistoon ~/ktp siirtyminen epäonnistui",
      "Error while copying new Vagrantfile: %d": "Uuden Vagrantfile-tiedoston kopiointi epäonnistui: %d",
      "Could not create ~/ktp to %s": "Hakemiston ~/ktp luominen sijaintiin %s epäonnistui",
      "Could not create ~/ktp-jako to %s": "Hakemiston ~/ktp-jako luominen sijaintiin %s epäonnistui",
      "Failed to delete %s": "Tiedoston %s poistaminen epäonnistui",
      "Failed to rename %s to %s": "Tiedoston %s nimeäminen tiedostoksi %s epäonnistui",
      "Failed to create file %s": "Tiedoston %s luominen epäonnistui",
      "Failed to retrieve %s": "Sijainnista %s lataaminen epäonnistui",
      "Could not copy body from %s to %s": "Sisällön %s kopioint sijaintiin %s epäonnistui",

      // start
      // Already defined: "Could not change to vagrant directory ~/ktp"
    }
  } else if current_language == "sv" {
    xlate_strings = map[string]string {
      // Main window, groups
      "Language": "Språk",
      "Basic Functions": "Grundfunktionaliteter",
      "Abitti": "Abitti-prov",
      "Matriculation Exam": "Studentprov",

      // Main window, buttons
      "Start Stickless Exam Server": "Starta sticklös provlokalsserver",
      "Install or update Abitti Stickless Exam Server": "Installera eller uppdatera sticklös server för Abitti",
      "Install or update Stickless Matriculation Exam Server": "Installera eller uppdatera sticklös server för studentexamen",
      "Make Stickless Exam Server Backup": "Säkerhetskopiera den sticklösa servern",
      "Exit": "SV: Poistu",

      // Main window, other
      "Did not get a path for a new Vagrantfile": "Ge positionen för den nya Vagrantfile-filen",
      "Could not execute vagrant. Are you sure you have installed HashiCorp Vagrant?": "Startning av programmet Vagrant lyckades inte. Är du säker, att HashiCorp Vagrant har installerats på datorn?",
      "Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?": "Startning av programmet VBoxManage lyckades inte. Är du säker, att Oracle VirtualBox har installerats på datorn?",

      // Backup dialog
      "naksu: SaveTo": "naksu: Spara till ",
      "Please select target path": "Välj sparningsställe",
      "Save": "Spara",
      "Cancel": "Avbryt",

      // mebroutines
      "command failed: %s": "Komandot misslyckades: %s",
      "Failed to execute %s": "Utförning av komandot misslyckades: %s",
      "Could not chdir to %s": "Förflyttning till katalogen %s misslyckades",
      "Error": "Fel",
      "Warning": "Varning",
      "Info": "För information",

      // backup
      "File %s already exists": "SV: Tiedosto %s on jo olemassa",
      "Backup has been made to %s": "SV: Varmuuskopio on talletettu tiedostoon %s",
      "Could not get vagrantbox ID: %d": "SV: Vagrantboxin ID:tä ei voitu lukea: %d",
      "Could not make backup: failed to get disk UUID": "SV: Varmuuskopion ottaminen epäonnistui: levyn UUID:tä ei löytynyt",
      "Could not back up disk %s to %s": "SV: Varmuuskopion ottaminen levystä %s tiedostoon %s epäonnistui",

      // backup, getmediapath
      "Home directory": "SV: Kotihakemisto",
      "Temporary files": "SV: Tilapäishakemisto",
      "Profile directory": "SV: Profiilihakemisto",

      // install
      "Could not change to vagrant directory ~/ktp": "SV: Vagrant-hakemistoon ~/ktp siirtyminen epäonnistui",
      "Error while copying new Vagrantfile: %d": "SV: Uuden Vagrantfile-tiedoston kopiointi epäonnistui: %d",
      "Could not create ~/ktp to %s": "SV: Hakemiston ~/ktp luominen sijaintiin %s epäonnistui",
      "Could not create ~/ktp-jako to %s": "SV: Hakemiston ~/ktp-jako luominen sijaintiin %s epäonnistui",
      "Failed to delete %s": "SV: Tiedoston %s poistaminen epäonnistui",
      "Failed to rename %s to %s": "SV: Tiedoston %s nimeäminen tiedostoksi %s epäonnistui",
      "Failed to create file %s": "SV: Tiedoston %s luominen epäonnistui",
      "Failed to retrieve %s": "SV: Sijainnista %s lataaminen epäonnistui",
      "Could not copy body from %s to %s": "SV: Sisällön %s kopioint sijaintiin %s epäonnistui",

      // start
      // Already defined: "Could not change to vagrant directory ~/ktp"
    }
  } else {
    xlate_strings = nil
  }
}

func Get(key string) string {
  if xlate_strings == nil {
    return key
  }

  new_string := xlate_strings[key]

  if new_string == "" {
    return key
  }

  return new_string
}
