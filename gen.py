import sys
import requests
import json

numnode = sys.argv[1]


def genDockerCompose():
    w = open("docker-compose.yml", "w")
    w.write("version: '3'\n")
    w.write("services:\n")
    for i in range(int(numnode)):
        w.write("  node_"+str(i)+":\n")
        w.write("    build:\n")
        w.write("      context: .\n")
        w.write("      dockerfile: ./Dockerfile\n")
        w.write("    ports:\n")
        w.write("      - "+str(8000+i)+":8000\n")
        w.write("    command: ./dpos 8000 8080 {}\n".format(i))
        w.write("    hostname: node_"+str(i)+"\n")
        # w.write("\t\tnetwork_mode: \"host\"\n")
    w.close()


if __name__ == "__main__":
    genDockerCompose()
    # startcluster()
