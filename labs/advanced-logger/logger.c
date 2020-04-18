#include <stdio.h>
#include <stdarg.h>
#include <string.h>
#include <syslog.h>
#include "logger.h"

char lType = 0;

int initLogger(char *logType){
    printf("Initializing Logger on: %s\n", logType);
    if (strcmp(logType, "stdout")==0){
        lType= 0;
        return 0;
    }
    if (strcmp(logType, "syslog")==0){
        lType= 1;
        return 0;
    }
    return -1;
}

int infof(const char *format, ...){
    va_list text;
    va_start(text, format);
    if(lType==1){
        vsyslog(LOG_INFO, format, text);
        va_end(text);
        return 1;
    }
    textcolor(BRIGHT, GREEN, BLACK);
    int ok = vfprintf(stdout, format, text);
    textcolor(RESET, WHITE, BLACK);
    va_end(text);
    return ok;
}

int warnf(const char *format, ...){
    va_list text;
    va_start(text, format);
    if(lType==1){
        vsyslog(LOG_WARNING, format, text);
        va_end(text);
        return 1;
    }
    textcolor(BRIGHT, YELLOW, BLACK);
    int ok = vfprintf(stdout, format, text);
    textcolor(RESET, WHITE, BLACK);
    va_end(text);
    return ok;
}

int errorf(const char *format, ...){
    va_list text;
    va_start(text, format);
    if(lType==1){
        vsyslog(LOG_ERR, format, text);
        va_end(text);
        return 1;
    }
    textcolor(BRIGHT, MAGENTA, BLACK);
    int ok = vfprintf(stdout, format, text);
    textcolor(RESET, WHITE, BLACK);
    va_end(text);
    return ok;
}

int panicf(const char *format, ...){
    va_list text;
    va_start(text, format);
    if(lType==1){
        vsyslog(LOG_EMERG, format, text);
        va_end(text);
        return 1;
    }
    textcolor(BRIGHT, RED, WHITE);
    int ok = vfprintf(stdout, format, text);
    textcolor(RESET, WHITE, BLACK);
    va_end(text);
    return ok;
}

void textcolor(int attr, int fg, int bg){
        char command[13];
	    sprintf(command, "%c[%d;%d;%dm", 0x1B, attr, fg + 30, bg + 40);
	    printf("%s", command);
}
