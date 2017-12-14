#!/usr/bin/env python
# -*- coding: UTF-8 -*-

from multiprocessing.dummy import Pool as ThreadPool
from multiprocessing import Pool 

## from http://segmentfault.com/a/1190000000414339

## IO密集 多线程
def RunMultiIoJob(func, params, num=5):
    # Make the Pool of workers
    pool = ThreadPool(num) 
    # Open the urls in their own threads
    # and return the results
    results = pool.map(func, params)
    #close the pool and wait for the work to finish 
    pool.close() 
    pool.join() 
    return results

## CPU密集 多进程
def RunMultiCpuJob(func, params, num=5):
    pool = Pool()
    results = pool.map(func, params)
    pool.close()
    pool.join()
    return results