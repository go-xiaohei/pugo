```toml
title = "Git"
date = "2016-02-04 15:00:00"
slug = "en/docs/deploy/git"
hover = "docs"
lang = "en"
template = "docs.html"
```

`PuGo` can copy local static files to your cloned repository, and push a new commit to remote directory.

```bash
pugo deploy git --local="dest" --repo="repo" --message="commit {time}" --branch="master"
```

`--local` set local directory that saving static files after building.

`--repo` set cloned git repository directory.

`--message` set commit message, `{time}` is a placeholder of the commit created time.

`--branch` set remote branch that new commit push to.

#### Steps

- Git clone your repository to save static files. Because `PuGo` just run `git push` wihout username and password options, so you need clone the repository via `git://` or `user:password`.

```bash
    git clone git://your-repo.git
    // or
    git clone https://user:password@your-repo.git
```
    
- Make sure the branch is proper to push to. Such as in Github, you need `gh-pages`.

```bash
    git checkout gh-pages
```
    
- Then run `build` command to build static files in [dest] directory.

```bash
    pugo build --dest="[dest]"
```

- Run `deploy` command with correct flags from above steps:

```bash
    pugo deploy --local="dest" --repo="your-repo" --branch="gh-pages"
```
    
#### Warning

The strategy of `PuGo` git deployment is copying compiled files to your repository directory. It will overwrite all same-named files, and will not affect other files if not same. So sometimes you need check whether the git changes following your mind. 