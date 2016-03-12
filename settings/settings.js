

var static = {
    virtual: "/static/",
    localroot: "./public",
    compress: "./settings/compress.json"
}
var database = {
    type: "sqlite3",
    uri: "./resource/db/lottery.2.sqlite3"
}

var server = {
    port: ":9000",
    host: ":9000"
}
var template = {
    home: "./templates",
    delimes: {
        left: "{%",
        right: "%}"
    },
    charset: "UTF-8",
    reload: true
}
var defaultvars = {
    appname: ""
}
var admin = {
    passwd: "micro1867321",
    salt:"kjH*(Y(*Y9uIHf34f34dON3JK*(y*T&T)))"
}
var log = {
    path: "./.log",
    format: "060102.log",
}
var headers = {
    Server: "Go/1.6.0",
    "X-Tango-Version":"0.4.6"
}

var time = {
    zone: "Hongkong"
}

var tls = {
    use: false,
    cert: "./resource/tls/cert.pem",
    key: "./resource/tls/key.pem"
}