import os
import sys
import time
import random
import queue
import multiprocessing
from multiprocessing import Pool
from queue import Queue
from sklearn.model_selection import ParameterGrid


def benchmark(item):
    # print(item)
    os.system(item['command'])


if __name__ == "__main__":
    wrk_folder = sys.argv[1]
    # command=['pwd', 'pwd', 'pwd']
    # print(command)
    path_wrk = '{}/wrk'.format(wrk_folder)
    destination_ip = 'localhost'

    start_time = time.time()
    print(start_time)
    command_round = []
    for i in range(10):
        _command_round = []
        for j in range(4):
            destimation_port = 8000 + j
            time_benchwrk = random.choice(range(20))
            _command = '{} -t12 -c400 -d{}s -s {}/scripts/post_lua/post_{}.lua  http://{}:{}'\
                .format(path_wrk, time_benchwrk, wrk_folder, j, destination_ip, destimation_port)
            _command_round.append(_command)
        command_round.append(_command_round)

    for i in range(10):
        param_grid = {
            'command': command_round[i]
        }
        print(param_grid)
        queue = Queue()
        for item in list(ParameterGrid(param_grid)):
            queue.put_nowait(item)
        pool = Pool(10)
        pool.map(benchmark, list(queue.queue))
        pool.close()
        pool.join()
        pool.terminate()
        print("Done")
