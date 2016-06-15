#!/usr/bin/python
# -*- coding: UTF-8 -*-

from multiprocessing import cpu_count
from multiprocessing import Pool

import MySQLdb as db
import demjson as json
import getopt
import time
import sys
import logging

connConf = {
    "db":       "go_test",
    "port":     3306,
    "host":     "localhost",
    "user":     "root",
    "passwd":   "rocky",
    "charset":  "utf8",
}

conf = {
    "back":     False,
    "store":    False,
    "table":    "tel_info",
    "proNum":   cpu_count(),
    "logLevel": "INFO",
}

selectSql = "select * from %s limit %d offset %d"
insertSql = "insert into %s values ('%s', '%s', '%s')"
sqlCount = 0

# Help information
def Help():
    print """\
Usage:
    python test.py command [args]

Commands:
    -?, --help:              Display this help and exit.
    -h, --host=name:         Connect to database.
    -u, --user=name:         User for login database.
    -p, --passwd=name:       Passwd for login database.
    -d, --database=name:     Database name.
    -t, --table=name:        Table name.
    -P, --port=name:         Port number to use for connection.
    -c, --charset=name:      Set character.
    -n, --num=name:          Child process number.
    -l, --logLevel=name:     Set log level.
    """


# Initialization log
def logInit():
    logger = logging.getLogger()
    handler = logging.FileHandler("log.log")
    logger.addHandler(handler)
    level = logging.getLevelName(conf["logLevel"])
    logger.setLevel(level)


# Parse command options
def getOption():
    for k, v in ops:
        if k in ["-?", "--help"]:
            help()

        elif k in ["-h", "--host"]:
            connConf["host"] = v

        elif k in ["-u", "--user"]:
            connConf["user"] = v

        elif k in ["-p", "--passwd"]:
            connConf["passwd"] = v

        elif k in ["-d", "--database"]:
            connConf["db"] = v

        elif k in ["-P", "--port"]:
            connConf["port"] = int(v)

        elif k in ["-c", "--charset"]:
            connConf["charset"] = v

        elif k in ["-t", "--table"]:
            conf["table"] = v

        elif k in ["-n", "--num"]:
            conf["proNum"] = int(v)

        elif k in ["-l", "--logLevel"]:
            conf["logLevel"] = v

    for k in args:
        conf[k] = True

# Backup Data
def backup(cursor, n, i):
    if i == proNum-1:
        c = sqlCount - n * i
    else:
        c = n

    # Select
    sql = selectSql % (conf["table"], c, i*n)
    logging.log(logging.DEBUG, "SQL: %s", sql)

    cursor.execute(sql)
    data = cursor.fetchall()

    # Encode
    logging.log(logging.INFO, "Start Child Process %d, have %d data, save data to file[bak%d.txt]...",
                i, len(data), i)
    json.encode_to_file("bak"+str(i)+".txt", data, overwrite=True)

    logging.log(logging.INFO, "Child Process %d Exit...", i)


# Restore Data
def restore(cursor, i):
    # Decode
    logging.log(logging.INFO, "Child Process %d begin decode file[bak%d.txt]...", i, i)
    data = json.decode_file("bak"+str(i)+".txt")

    logging.log(logging.INFO, "Child Process %d decode file end, begin insert data to database...", i)

    # Insert
    for v in data:
        sql = insertSql % (conf["table"], v[0], v[1], v[2])

        try:
            # logging.log(logging.DEBUG, i[0], i[1], i[2])
            cursor.execute(sql)
        except:
            logging.log(logging.ERROR, "InsertError: %s", sql)
            raise

    logging.log(logging.INFO, "Child Process %d exit...", i)


# Child process handle function
def handle(n, i):
    conn = db.connect(**connConf)
    cursor = conn.cursor()

    if conf["back"]:
        backup(cursor, n, i)

    if conf["store"]:
        restore(cursor, i)

    cursor.close()
    conn.commit()
    conn.close()


# Child process handle sql count
def getSqlCount():
    conn = db.connect(**connConf)

    cursor = conn.cursor()
    cursor.execute("select count(*) from %s" % conf["table"])
    count = cursor.fetchone()[0]

    n = count / proNum
    cursor.close()
    conn.close()

    return count, n


if __name__ == '__main__':
    try:
        ops, args = getopt.getopt(sys.argv[1:], "?h:u:p:d:t:P:c:n:l:",
                                 ["help", "host=", "user=", "passwd=", "database=", "table="])
    except:
        Help()
        sys.exit(1)

    start = time.time()
    getOption()

    logInit()
    logging.log(logging.DEBUG, "Config Information: \n conf: %s\n connConf: %s",
                str(conf), str(connConf))

    proNum = conf["proNum"]
    sqlCount, num = getSqlCount()
    logging.log(logging.DEBUG, "SQL total number: %d, child process handle SQL number: %d",
                sqlCount, num)

    # Create process pool object
    p = Pool(proNum)
    for i in range(proNum):
        p.apply_async(handle, args=(num, i))

    p.close()
    p.join()

    end = time.time()
    logging.log(logging.INFO, "Running Time %fs", (end - start))






