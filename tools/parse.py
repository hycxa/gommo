#!/usr/bin/env python
#coding=utf8

"""
协议生成工具
"""

import sys, os, re
from random import random

def random_shuffle(beginNum, endNum):
    retTab = []
    for i in range(beginNum, endNum):
        retTab.append(i)

    length = len(retTab)
    for i in range(0,length):
        randomi = int(random()*length)
        retTab[0],retTab[randomi] = retTab[randomi],retTab[0]
    retTab.insert(0, 1)
    return retTab

userPacketIDMin = 10000
userPacketIDMax = 100000
genPacketIDTab = random_shuffle(userPacketIDMin, userPacketIDMax)

sysPacketIDMin = 0
sysPacketIDMax = 9999
genSysTab = random_shuffle(0, 9999)

def getOnePacketID(tab):
    length = len(tab)
    if tab[0] >= length:
        return None
    retID = tab[tab[0]]
    tab[0] = tab[0] + 1
    return retID

class Packet(object):
    def __init__(self, n, s, d, p):
        if s == "PACKAGE_SYSTEM":
            self.typeid = getOnePacketID(genSysTab)
        elif s == "PACKAGE_USER":
            self.typeid = getOnePacketID(genPacketIDTab)
        self.name = n
        self.desc = d
        self.payload = p
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

def new_packet(packet_list, packet_name, packet_scope, packet_desc, packet_payload):
    if packet_name is not None and packet_payload is not None and packet_scope is not None:
        packet = Packet(packet_name, packet_scope, packet_desc, packet_payload)
        packet_list.append(packet)

def getOneLineStructName(line):
    nameTab = re.findall(r"type\s+(\S+)\s+struct\s+{", line)
    if len(nameTab) >= 1:
        return nameTab[0]
    else:
        return None

def parse_packet(packet_buf):
    L = [line.strip() for line in packet_buf.split('\n')]
    L = [line for line in L if line and line[0] != '#']
    packet_list = []

    start_packet = False
    curStructName = None
    packet_name = None
    packet_scope = None
    packet_desc = None
    for line in L:
        idx = line.find(':')
        if idx < 0:
            structName = getOneLineStructName(line)
            if structName is not None:
                new_packet(packet_list, packet_name, packet_scope, packet_desc, curStructName)
                packet_name = packet_desc = None
                curStructName = structName
            continue

        if line[:idx] == 'name':
            if start_packet:
                new_packet(packet_list, packet_name, packet_scope, packet_desc, curStructName)
                packet_name = packet_desc = None
            start_packet = True
            packet_name = line[idx+1:]
        elif line[:idx] == 'desc':
            packet_desc = line[idx+1:]
        elif line[:idx] == 'packageScope':
            packet_scope = line[idx+1:]

    new_packet(packet_list, packet_name, packet_scope, packet_desc, curStructName)
    packet_name = packet_desc = packet_scope = None
    return packet_list

def gen_go_packet(packet_list):
    #协议名定义生成
    f = open(os.path.join('./', 'name.go'), 'w')
    f.write("""package proto\n\n
        const (
    """)
    for item in packet_list:
        f.write("\t%s = %d\n" %(item.name, item.typeid))
    f.write(")\n")
    f.close()

    #协议解析及编码生成
    decodef = open(os.path.join('./', 'decode.go'), 'w')
    decodef.write("""package proto \n\n
        import (
        "bytes"
        "encoding/gob"
        )\n
        func GetPacketScope(id PacketID) int{""")

    decodef.write("if id >= %d && id <= %d {\nreturn PACKAGE_SYSTEM\n" %(sysPacketIDMin, sysPacketIDMax))
    decodef.write("}else if id >= %d && id <= %d{\nreturn PACKAGE_USER\n" %(userPacketIDMin, userPacketIDMax))

    decodef.write("""}else{
    return -1
        }
    }\n
    func EncodeMsg(buff * bytes.Buffer, msg *Message) bool {
	        enc := gob.NewEncoder(buff)
        	err := enc.Encode(msg.Sender)
	        if err != nil {
	            checkErr(err)	
                return false
	        }
            err = enc.Encode(msg.PacketID)
            if err != nil {
                checkErr(err)
                return false
            }
            switch msg.PacketID {\n""")

    for item in packet_list:
        decodef.write("case %s:\nerr = enc.Encode(msg.Data.(%s))\n" %(item.name, item.payload))

    decodef.write("""default:
                return false
            }
            if err != nil {
                checkErr(err)
                return false
            }
            return true
        }\n""")

    #协议解码生成
    decodef.write("""
        func DecodeMsg(buff *bytes.Buffer) (bool, *Message) {
           msg := Message{}
           dec := gob.NewDecoder(buff)
           err := dec.Decode(&(msg.Sender))
           if err != nil {
                checkErr(err)
                return false, nil
           }
           err = dec.Decode(&(msg.PacketID))
           if err != nil {
                checkErr(err)
                return false, nil
           }
           switch msg.PacketID {
                  """)
    for item in packet_list:
        decodef.write("case %s:\nvar data %s\nerr = dec.Decode(&data)\nmsg.Data = data\n" %(item.name, item.payload))
    decodef.write("""default:
                  return false, nil
              }
              if err != nil {
                    checkErr(err)
                    return false, nil
              }
                return true, &msg
            }\n\n""")

def parse(packet_buf):
    packet_list = parse_packet(packet_buf)

    gen_go_packet(packet_list)

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print 'usage: ./parse.py proto_dir [gen_dir]'
        sys.exit(0)

    path_pre = sys.argv[1]
    try:
        packet_buf = open(os.path.join(path_pre, 'struct.go'), 'r').read()
    except IOError, e:
        print 'Open proto file failed:', e
        sys.exit(0)
    parse(packet_buf)
