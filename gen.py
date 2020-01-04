import sys
import json

numnode = int(sys.argv[1])


def genDockerCompose():
    perFile = 50
    numFile = int(numnode/50)
    for i in range(numFile):
        w = open("docker-compose-{}.yml".format(i), "w")
        w.write("version: '3'\n")
        w.write("services:\n")
        if i == 0:
            w.write("  dns:\n")
            w.write("    image: mdpos\n")
            w.write("    network_mode: host\n")
            w.write("    command: ./mdpos 1\n")
        for j in range(i*50, min((i+1)*50, numnode)):
            # print(j)
            w.write("  node_"+str(j)+":\n")
            # w.write("    image: mdpos\n")
            # w.write("    build:\n")
            # w.write("      context: .\n")
            # w.write("      dockerfile: ./Dockerfile\n")
            w.write("    image: mdpos\n")
            w.write("    network_mode: host\n")
            w.write("    ports:\n")
            w.write("      - "+str(8000+j)+":8000\n")
            w.write("    command: ./mdpos 0 {} {} {}\n".format(8000+j, 9000+j, j))
            w.write("    hostname: node_"+str(j)+"\n")
            # w.write("\t\tnetwork_mode: \"host\"\n")
        w.close()


if __name__ == "__main__":
    genDockerCompose()
    # startcluster()
