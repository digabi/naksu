# Testing and Building with Travis-CI

To build Naksu with Travis-CI.org:

  1. Create GitHub Personal Access Token:
    * Your Profile > Settings > Developer Settings > Personal access tokens
    * Click "Generate new token"
    * Note, e.g. "Automatic releases for mplattu/naksu"
    * Scope: `public_repo` (check this and only this)
    * Click "Generate token"
    * Get the token (e.g. `d982bdb37a4a1c215c09528da5a615cfac199924`)
  1. Enter the token as a Travis CI environment variable
    * Go to Naksu repo > More options > Settings
    * Add environment variable `GITHUB_TOKEN`, which has the token as a value,
      "all branches" and do not display the value in the build log
    * Add another variable `GITHUB_REPO` which contains the slug to your repo
      (e.g. `yourusername/naksu`). For debugging purposes you might want to
      display this value in the log.
  1. After a successful build Travis creates a Draft release with built binaries (`naksu_linux_amd64.zip` and `naksu_windows_amd64.zip`).
     The existing unnamed draft releases are deleted by `tools/github_purge_draft_releases.py` executed by Travis. If you
     wish to save a draft release after a build you can give a name to it.
