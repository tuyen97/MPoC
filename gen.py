import sys
import json

numnode = sys.argv[1]


def genDockerCompose(start_from):
    w = open("docker-compose-{}.yml".format(start_from), "w")
    w.write("version: '3'\n")
    w.write("services:\n")
    for i in range(int(numnode)):
        w.write("  node_"+str(i+int(start_from))+":\n")
        # w.write("    image: mdpos\n")
        w.write("    build:\n")
        w.write("      context: .\n")
        w.write("      dockerfile: ./Dockerfile\n")
        w.write("    network_mode: host\n")
        w.write("    ports:\n")
        w.write("      - "+str(8000+i+int(start_from))+":8000\n")
        w.write("    command: ./mdpos {} {} {}\n".format(8000+i+int(start_from),9000+i+int(start_from),i+int(start_from)))
        w.write("    hostname: node_"+str(i+int(start_from))+"\n")
        # w.write("\t\tnetwork_mode: \"host\"\n")
    w.close()


if __name__ == "__main__":
    genDockerCompose(sys.argv[2])
    # startcluster()
