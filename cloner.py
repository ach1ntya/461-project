#! /usr/bin/env python3
import sys
import git
import subprocess
from git import Repo

def clone_repo(url, dir_path):

    Repo.clone_from(str(url), "/home/shay/a/dkholode/Desktop/461-project/clonedir" + str(dir_path))
    subp = subprocess.run(["git rev-list --all --count"], stdout=subprocess.PIPE, text=True, shell=True, cwd= "/home/shay/a/dkholode/Desktop/461-project/clonedir" + str(dir_path))
    print(subp.stdout)
    

def main():
    repo_url = sys.argv[1]
    dir_path = sys.argv[2]
    clone_repo(repo_url, dir_path)    

if __name__ == "__main__":
    main()