import os
import queue
import multiprocessing
from multiprocessing import Pool
from queue import Queue
from sklearn.model_selection import ParameterGrid


def benchmark(item):
    # print(item)
    os.system(item['command'])


if __name__ == "__main__":
    # command=['pwd', 'pwd', 'pwd']
    # print(command)
    command = ['/home/bkc_3/Desktop/wrk/wrk -t12 -c400 -d30s -s /home/bkc_3/Desktop/wrk/scripts/post.lua  http://127.0.0.1:8000',
               '/home/bkc_3/Desktop/wrk/wrk -t12 -c400 -d30s -s /home/bkc_3/Desktop/wrk/scripts/post.lua  http://127.0.0.1:8001',
               '/home/bkc_3/Desktop/wrk/wrk -t12 -c400 -d30s -s /home/bkc_3/Desktop/wrk/scripts/post.lua  http://127.0.0.1:8002',
               '/home/bkc_3/Desktop/wrk/wrk -t12 -c400 -d30s -s /home/bkc_3/Desktop/wrk/scripts/post.lua  http://127.0.0.1:8003']
    param_grid = {
        'command': command
    }
    print(param_grid)
    queue = Queue()
    for item in list(ParameterGrid(param_grid)):
        queue.put_nowait(item)
    pool = Pool(4)
    pool.map(benchmark, list(queue.queue))
    pool.close()
    pool.join()
    pool.terminate()
    print("Done")

# import subprocess

# for i in range(5):
#     subprocess.call("/home/bkc_3/Desktop/wrk/wrk -t12 -c400 -d30s -s /home/bkc_3/Desktop/wrk/scripts/post.lua  http://127.0.0.1:{}".format(str(8000+i)),shell=True)
