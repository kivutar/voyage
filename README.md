# voyage

A gemini browser with lua scripting support and gaming / UI capabilies.

Instead of browsing gmi files, you browse lua programs. These programs can be games, applications, websites. They execute locally.

This is a work in progress, I did nothing to protect the user with proper sandboxing. People can use this to steal your personal informations, so don't run it unless you know what you are doing.

## Usage

    go build
    go build && ./voyage gemini://blog.kivutar.me/test.gmi
