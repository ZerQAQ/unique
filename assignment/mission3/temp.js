eid = 10
photoNum = 3
imgurls = []
skey = "skey"
root = "localhost:8080/kuro"

for (let i = 1; i <= 3; i++){
    imgurls.push(root + "/src/photo/" + eid.toString() + "/" + i.toString() + "?skey=" + skey)
}