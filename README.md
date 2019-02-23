# ADot

ADot is short for automatic dotfiles.

This project is not complete!

It automatically synchronises selected dotfiles to other machines using the power of git.

## Planned functionality

adot new \<URL>
 - git init
 - git add empty .adot file with URL
 - git push
 - Print out some instructions like running `adot add`
 - Print out warning about making sure the repo is private and safe! It could lead to system pwnage otherwise.

adot init
 - clone
 - copy down files, if a file exists, rename

adot link
 - adds a new file or directory to be synced.
 - copy up the file/dir.
 - git commit and push

adot unlink
 - unlinks a file or directory
 - file won't be deleted from your system

adot rm
 - unlinks and removes a file or directory

adot push
 - check repo state to be clean
 - iterate
   - if file is different to repo, copy to repo
 - if there are changes, commit, push
 - if there's a conflict do a pull rebase
 - otherwise let the user handle the conflict

adot pull
 - check repo state to be clean
 - git pull
 - iterate
   - if any files are different, backup/copy.

adot service
 - monitors filesystem for changes and can push automatically
 - adot pull runs on a regular basis

