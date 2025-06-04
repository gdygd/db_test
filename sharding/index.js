const app = require("express")();
const {Client} = require("pg");
const crypto = require("crypto")
const ConsistentHash = require("consistent-hash");
const hr = new ConsistentHash();
hr.add("5433")
hr.add("5434")
hr.add("5435")

const clients = {
    "5433": new Client ({
        "host": "10.1.0.119",
        "port": "5433",
        "user": "postgres",
        "password": "secret",
        "database": "postgres"
    }),
    "5434": new Client ({
        "host": "10.1.0.119",
        "port": "5434",
        "user": "postgres",
        "password": "secret",
        "database": "postgres"
    }),
    "5435": new Client ({
        "host": "10.1.0.119",
        "port": "5435",
        "user": "postgres",
        "password": "secret",
        "database": "postgres"
    })
}

connect();
async function connect() {
     await clients["5433"].connect();
     await clients["5434"].connect();
     await clients["5435"].connect();
}



app.get("/:urlId",  (req, res) => {
    //https://localhost:8081/fhy2h
    const urlId = req.params.urlId;
    const server = hr.get(urlId)
    await clients[server].query("select * from url_table where URL_ID = $1", [urlId]);


})

app.post("/", async (req, res) => {
    const url = req.query.url;
    // consistently hash this to get a port!
    const hash = crypto.createHash("sha256").update(url).digest("base64")
    const urlId = hash.substr(0,5);
    const server = hr.get(urlId)

    const result = await clients[server].query("insert into url_table (url, url_id) valudes ($1, $2)" [url, urlId]);
    if (result.rowsCount > 0) {
        res.send({
            "urlId": urlId,
            "url":url,
            "server": server
        })
    }
    else {
        res.sendStatus(404)
    }

    res.send({
        "urlId": urlId,
        "url":url,
        "server": server
    })
})

app.listen(8081, () => console.log("Listening 8081"))