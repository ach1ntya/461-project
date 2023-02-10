#! /usr/bin/env python3
import sys
import git
from git import Repo

def clone_repo(url):
    f = open(url, "r")
    url_string = f.readline().rstrip('\n')
    print(url_string) 
    Repo.clone_from(url_string, "/Users/Ben_Brown19/Desktop/school/ECE_461/cloneage/zlone")
    #repo = Repo.clone_from("https://github.com/brow1770/goprac.git", "/Users/Ben_Brown19/Desktop/school/ECE_461/cloneage/blone")
    print("success")
    f.close()

def main():
    repo_url = sys.argv[1]
    clone_repo(repo_url)    

if __name__ == "__main__":
    main()