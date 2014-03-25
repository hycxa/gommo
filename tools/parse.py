#!/usr/bin/env python
#coding=utf8

"""
协议生成工具
"""

import sys, os
from random import random

def random_shuffle(beginNum, endNum):
    retTab = []
    for i in range(beginNum, endNum):
        retTab.append(i)

    length = len(retTab)
    for i in range(0,length):
        randomi = int(random()*length)
        retTab[0],retTab[randomi] = retTab[randomi],retTab[0]
    return retTab

genPacketIDTab = random_shuffle(1000, 100000)
noUsePacketIDIndex = 0

def getOnePacketID():
    global genPacketIDTab
    global noUsePacketIDIndex

    retID = genPacketIDTab[noUsePacketIDIndex]
    noUsePacketIDIndex = noUsePacketIDIndex + 1
    length = len(genPacketIDTab)
    if noUsePacketIDIndex >= length:
        maxValue = genPacketIDTab[length - 1]
        genPacketIDTab = random_shuffle(maxValue + 1, maxValue + 100000)
        noUsePacketIDIndex = 0
    return retID

class Packet(object):
    # name 自动全部大写
    #结构体自动首字母大写
    def __init__(self, n, d, p):
        self.typeid = getOnePacketID()
        self.name = n.upper()
        self.desc = d
        self.payload = p[:1].upper() + p[1:]
        self.check()
    def __repr__(self):
        return "\nAPI typeid: %s \nname: %s desc: %s payload: %s\n" %\
            (self.typeid, self.name, self.desc, self.payload)
    def check(self):
        for i in [self.name, self.typeid, self.payload]:
            if (i is None) or (i == ''):
                print "API error:", self.typeid, self.name, self.payload
                sys.exit(-1)
        if self.desc is None:
            self.desc = ''

def new_packet(packet_list, packet_name, packet_desc, packet_payload):
    packet = Packet(packet_name, packet_desc, packet_payload)
    packet_list.append(packet)

def parse_packet(packet_buf):
    L = [line.strip() for line in packet_buf.split('\n')]
    L = [line for line in L if line and line[0] != '#']
    packet_list = []

    start_packet = False
    packet_name = None
    packet_desc = None
    packet_payload = None
    for line in L:
        idx = line.find(':')
        if idx < 0: continue
        if line[:idx] == 'name':
            if start_packet:
                new_packet(packet_list, packet_name, packet_desc, packet_payload)
                packet_name = packet_desc = packet_payload = None
            start_packet = True
            packet_name = line[idx+1:]
        elif line[:idx] == 'payload':
            packet_payload = line[idx+1:]
        elif line[:idx] == 'desc':
            packet_desc = line[idx+1:]

    if packet_name is not None:
        new_packet(packet_list, packet_name, packet_desc, packet_payload)
    return packet_list

def gen_go_packet(packet_list):
    f = open(os.path.join('./', 'packet_name.go'), 'w')
    f.write("""package proto\n\n
const (
    """)
    for item in packet_list:
        f.write("\t%s = %d\n" %(item.name, item.typeid))
    f.write(")\n")
    f.close()
    

def parse(packet_buf):
    packet_list = parse_packet(packet_buf)

    gen_go_packet(packet_list)

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print 'usage: ./parse.py proto_dir [gen_dir]'
        sys.exit(0)

    path_pre = sys.argv[1]
    try:
        packet_buf = open(os.path.join(path_pre, 'packet.txt'), 'r').read()
    except IOError, e:
        print 'Open proto file failed:', e
        sys.exit(0)
    parse(packet_buf)
