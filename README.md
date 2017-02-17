SOUP ARCHIVE
------------

[![Build Status](https://travis-ci.org/bestform/souparchive.svg?branch=master)](https://travis-ci.org/bestform/souparchive)

Usage:

    ./souparchive -user YOURUSERNAME
    
This will save all the entries in your soup.io rss feed in the archive folder.

Subsequent calls will remember already saved items so you can run this script as a cron job to continiously archive your soup.io feed.

Keep in mind that at its current state the script will not keep track of the ordering of the files. Nor will it save two different files with the same name as different entries. (PRs welcome)