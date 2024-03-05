#include <stdio.h>
#include <stdlib.h>
#include <hiredis/hiredis.h>
#include "redistest.h"

void pingRedis() {
    redisReply *reply;
    redisContext *c;

    c = redisConnect("127.0.0.1", 6379);
    if (c->err) {
        printf("error: %s\n", c->errstr);
        return;
    }

    reply = redisCommand(c,"PING %s", "Hello World");
    printf("PING Response: %s\n", reply->str);
    freeReplyObject(reply);

    redisFree(c);
}
