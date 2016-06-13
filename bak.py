#!/usr/bin/python2.6
# -*- coding: UTF-8 -*-

import MySQLdb as db
import demjson as json
import sys
import os

HOST = "localhost"
USER = "root"
PASSWD = "rocky"
DBNAME = "go_test"

selectSql = "select * from %s order by mts"
insertSql = "insert into %s values ('%s', '%s', '%s')"

# 备份数据
def backup(cursor, param):
    sql = selectSql % param

    cursor.execute(sql)
    data = cursor.fetchall()

    # 转换成json格式
    if os.path.exists("bak.txt"):
        os.remove("bak.txt")

    print "有%d条数据, 将数据存入文件" % len(data)
    json.encode_to_file("bak.txt", data)
    print "完成"

# 还原数据
def restore(cursor, param):
    print "开始解析文件"
    data = json.decode_file("bak.txt")
    print "解析文件结束"

    print "开始插入数据"
    for i in data:
        print i[0], i[1], i[2]
        sql = insertSql % (param, i[0], i[1], i[2])

        try:
            cursor.execute(sql)
        except:
            print sql
    print "完成"


if __name__ == '__main__':
    if len(sys.argv) < 3:
        print "Usage: python bak.py [ back | store ] TableName"
        os._exit(1)

    # 连接数据库
    conn = db.connect(HOST, USER, PASSWD, DBNAME, charset='utf8')

    # 获取游标执行sql
    cursor = conn.cursor()

    if sys.argv[1] == "back":
        backup(cursor, sys.argv[2])

    if sys.argv[1] == "store":
        restore(cursor, sys.argv[2])

    # 关闭游标
    cursor.close()
    conn.close()
