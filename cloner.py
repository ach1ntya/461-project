#! /usr/bin/env python3
import sys
import git
import subprocess
from git import Repo

def clone_repo(url, dir_path):
    #f = open(url, "r")
    #url_string = f.readline().rstrip('\n')
    #print(url_string)
    Repo.clone_from(str(url), "/Users/Ben_Brown19/Desktop/school/ECE_461/461-Project/cloneDir/g" + str(dir_path))
    #subp = subprocess.Popen([])
    #subp = subprocess.call(["cd", "/Users/Ben_Brown19/Desktop/school/ECE_461/461-Project/cloneDir/g" + str(dir_path), "&&", "git", "rev-list", "--all", "--count"], shell= True, cwd= "/Users/Ben_Brown19/Desktop/school/ECE_461/461-Project")
    subp = subprocess.run(["git rev-list --all --count"], stdout=subprocess.PIPE, text=True, shell=True, cwd= "/Users/Ben_Brown19/Desktop/school/ECE_461/461-Project/cloneDir/g" + str(dir_path))
    print(subp.stdout)
    

def main():
    repo_url = sys.argv[1]
    dir_path = sys.argv[2]
    clone_repo(repo_url, dir_path)    

if __name__ == "__main__":
    main()