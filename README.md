# Why

You've crafted a pièce de résistance of a framework, a chef-d'œuvre of API's, a tour de force of chain'ed blocks.
You are about submit your final pull request.

### BUT WAIT!!!

Is the world ready for your masterpeice?
Do they really deserve this perfection of 1's and 0's? You check your bitcoin wallet, 899... 1872... 443... 210...
crap, it's fluctuating again, this is no time for a payout. GitHub's private repo's are too expensive, for you've spent this months
allowance on macchiatos and biscottis at the local coffee shop.

Wait, what is this? A state-of-the-art peer-to-peer file-transfer-protocal that hides your files in plain site?

The feather on your fedora just tingled.

# What

Redditfs is a command line interface to manage repositories that are stored, free of charge, on reddit.

<p align="center">
  <img src="https://github.com/maxchehab/redditfs/blob/master/images/demo.gif?raw=true" />
</p>

# How

Install [golang](https://golang.org/dl/), setup a [gopath](https://github.com/golang/go/wiki/SettingGOPATH), and `go get github.com/maxchehab/redditfs; go install github.com/maxchehab/redditfs`

### Creating an authorized application

The first time you use redditfs, you will be asked to create an [authorized reddit application](https://www.reddit.com/prefs/apps/).
To create and obtain the correct credentials, do the following:

1.  Create a subreddit. Preferably make this private.
2.  Go to http://www.reddit.com/prefs/apps/.
3.  Press the button on the bottom.

    ![Press the button on the bottom](https://github.com/maxchehab/redditfs/blob/master/images/image1.png?raw=true)

4.  Fill out the information

    ![Fill out the information](https://github.com/maxchehab/redditfs/blob/master/images/image2.png?raw=true)

5.  The following circled information is important.

    ![The following circled information is important](https://github.com/maxchehab/redditfs/blob/master/images/image3.png?raw=true)

6.  When prompted fill out the relevent information.

    ![When prompted fill out the relevent information](https://github.com/maxchehab/redditfs/blob/master/images/image4.png?raw=true)

7.  Voila!

    ![Voila](https://media.giphy.com/media/5heSxZbRPE2Ee6QaDi/giphy.gif)

## Features

#### Pull

`redditfs pull` will display, select, and download any repositories that are available into the current directory.

#### Push

`redditfs push` will upload current directory.

#### .redditfsignore

A `.redditfsignore` file specifies intentionally untracked files that Redditfs should ignore. This file should be located in the root of your repository

For example:

```git
dist/
node_modules/
.vscode/
```

This will ignore everything inside of a `dist`, `node_modules`, and `.vscode` directory.

## P.S.

I have only tested this on Debian Linux, therefore I can only confirm it works on my machine. For all y'all fancy programmers who can afford Macbook Pro's....

<p align="center">
  <img src="https://media.giphy.com/media/3ohA2ZD9EkeK2AyfdK/giphy.gif" width="450px" />
</p>

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE.md](LICENSE.md) file for details.

## Contact

Feel free to reachout to me at [twitter](https://twitter.com/maxchehab) or leave report an [issue](https://github.com/maxchehab/redditfs/issues) if there are any problems.


## P.S.S.
This is a joke.
