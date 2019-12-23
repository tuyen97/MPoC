import os
import sys
import time
import json

numnode = int(sys.argv[1])
numbp = int(sys.argv[2])

if __name__ == "__main__":
    balance = {}
    # init balance
    for i in range(numnode):
        balance[str(i)] = 1000000

    BPs = []
    for i in range(numbp):
        BPs.append(str(i))

    genesis = {}
    genesis["Balance"] = balance
    genesis["BPs"] = BPs
    genesis["Timestamp"] = int(round(time.time() * 1000000000))
    w = open("genesis.json", "w")
    w.write(json.dumps(genesis))
    w.close
    print(json.dumps(genesis))
