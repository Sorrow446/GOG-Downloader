# GOG-Downloader
CLI GOG downloader written in Go for Windows, Linux, macOS and Android.
![](https://i.imgur.com/aelWCRa.png)
![](https://i.imgur.com/8zQrXYX.png)
![](https://i.imgur.com/cxun5l0.png)
![](https://i.imgur.com/75KKwSG.png)    
[Windows, Linux, macOS and Android binaries](https://github.com/Sorrow446/GOG-Downloader/releases)

## Features
- Interactive CLI
- Filter & template system
- Resumable downloads of incomplete downloads

## Setup
Dump cookies to `cookies.json`. EditThisCookie Chrome extension's recommended. Netscape will also be supported soon. 

|Option|Info|
| --- | --- |
|platform|Item platform. windows/win, linux, mac/osx.
|language|Item language. en, cz, de, es, fr, it, hu, nl, pl, pt, br, sv, tr, uk, ru, ar, ko, cn, jp, all.
|folderTemplate|Game folder naming template. title, titlePeriods. Ex: {{.title}} [GOG], {{.titlePeriods}}.GOG
|goodies|Include goodies.
|outPath|Where to download to. Path will be made if it doesn't already exist.

# Usage
Args take priority over the config file.

Download by search:   
`gog_dl_x64.exe "destroy all humans"`
If more than one result is yielded, you'll be asked to choose.

Download from all owned Windows games:   
`gog_dl_x64 -p windows`

```
 _____ _____ _____    ____                _           _
|   __|     |   __|  |    \ ___ _ _ _ ___| |___ ___ _| |___ ___
|  |  |  |  |  |  |  |  |  | . | | | |   | | . | .'| . | -_|  _|
|_____|_____|_____|  |____/|___|_____|_|_|_|___|__,|___|___|_|

Usage: gog_dl_x64.exe [--platform PLATFORM] [--language LANGUAGE] [--template TEMPLATE] [--goodies] [--out-path OUT-PATH] [QUERY]

Positional arguments:
  QUERY

Options:
  --platform PLATFORM, -p PLATFORM
                         Item platform. windows/win, linux, mac/osx.
  --language LANGUAGE, -l LANGUAGE
                         Item language.
                         en, cz, de, es, fr, it, hu, nl, pl, pt, br, sv, tr, uk, ru, ar, ko, cn, jp, all.
  --template TEMPLATE, -t TEMPLATE
                         Game folder naming template. title, titlePeriods.
                         Ex: {{.title}} [GOG], {{.titlePeriods}}.GOG
  --goodies, -g          Include goodies.
  --out-path OUT-PATH, -o OUT-PATH
                         Where to download to. Path will be made if it doesn't already exist.
  --help, -h             display this help and exit
```

# Disclaimer  
- GOG Downloader has no partnership, sponsorship or endorsement with GOG or CD PROJEKT.
