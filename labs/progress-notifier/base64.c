 
#include <stdio.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <signal.h>
#include <fcntl.h>
#include <unistd.h>
#include "logger.h"


// Algorithm from: https://en.wikibooks.org/wiki/Algorithm_Implementation/Miscellaneous/Base64


#include <inttypes.h>
#include <string.h>
#define WHITESPACE 64
#define EQUALS     65
#define INVALID    66
#define ENCODED_SIZE 60000000

off_t fs;
int progress;

int base64encode(const void* data_buf, size_t dataLength, char* result, size_t resultSize);
int base64decode (char *in, size_t inLen, unsigned char *out, size_t *outLen);
void signalHandler(int sig);

static const char d[] = {
    66,66,66,66,66,66,66,66,66,66,64,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,62,66,66,66,63,52,53,
    54,55,56,57,58,59,60,61,66,66,66,65,66,66,66, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
    10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,66,66,66,66,66,66,26,27,28,
    29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66
};

int main(int argc, char *argv[]){

    signal(SIGINT, signalHandler);
    signal(29, signalHandler);

    int fd = open(argv[2], O_RDONLY);
    fs = lseek(fd,(size_t)0,SEEK_END);
    lseek(fd,(size_t)0,SEEK_SET);
    char* buffer = malloc(sizeof(char)*fs);
    if (read(fd, buffer, fs*sizeof(char))==-1){
        errorf("Something went wrong while reading the file");
        return -1;
    }

    if(strcmp(argv[1],"--encode")==0){
        char* encodedBuff = (char*)malloc(sizeof(char)*ENCODED_SIZE);
        base64encode(buffer, fs, encodedBuff, ENCODED_SIZE);
        int fen;
        fen = open("encoded.txt", O_WRONLY | O_CREAT, 0644);
        write(fen, encodedBuff, ENCODED_SIZE);
        close(fen);
        free(encodedBuff);
        infof("File encoded\n");    
    } else if(strcmp(argv[1],"--decode")==0){
        unsigned char* decodedBuff = (unsigned char*)malloc(sizeof(char)*fs);
        size_t dSize = sizeof(char)*fs;
        base64decode(buffer, fs, decodedBuff, &dSize);
        int fde = open("decoded.txt", O_WRONLY | O_CREAT,0644);
        write(fde, decodedBuff, dSize);
        close(fde);
        free(decodedBuff);
        infof("File decoded\n");
    }

    return 0;
}


int base64encode(const void* data_buf, size_t dataLength, char* result, size_t resultSize){
   const char base64chars[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
   const uint8_t *data = (const uint8_t *)data_buf;
   size_t resultIndex = 0;
   size_t x;
   uint32_t n = 0;
   int padCount = dataLength % 3;
   progress = 0;
   uint8_t n0, n1, n2, n3;

   for (x = 0; x < dataLength; x += 3) 
   {
      progress++;
      n = ((uint32_t)data[x]) << 16;

      if((x+1) < dataLength)
         n += ((uint32_t)data[x+1]) << 8;
      
      if((x+2) < dataLength)
         n += data[x+2];

      n0 = (uint8_t)(n >> 18) & 63;
      n1 = (uint8_t)(n >> 12) & 63;
      n2 = (uint8_t)(n >> 6) & 63;
      n3 = (uint8_t)n & 63;

      if(resultIndex >= resultSize) return 1;   
      result[resultIndex++] = base64chars[n0];
      
      if(resultIndex >= resultSize) return 1;   
      result[resultIndex++] = base64chars[n1];

      if((x+1) < dataLength)
      {
         if(resultIndex >= resultSize) return 1;   
         result[resultIndex++] = base64chars[n2];
      }

      if((x+2) < dataLength)
      {
         if(resultIndex >= resultSize) return 1;   
         result[resultIndex++] = base64chars[n3];
      }
   }

   if (padCount > 0) 
   { 
      for (; padCount < 3; padCount++) 
      { 
         if(resultIndex >= resultSize) return 1;  
         result[resultIndex++] = '=';
      } 
   }
   
   if(resultIndex >= resultSize) return 1;   
   result[resultIndex] = 0;
   return 0;   
}

int base64decode (char *in, size_t inLen, unsigned char *out, size_t *outLen) { 
    char *end = in + inLen;
    char iter = 0;
    uint32_t buf = 0;
    size_t len = 0;
    progress = 0;
    
    while (in < end) {
        char c = d[*in++];
        
        switch (c) {
        case WHITESPACE: continue;   
        case INVALID:    return 1;   
        case EQUALS:                 
            in = end;
            continue;
        default:
            buf = buf << 6 | c;
            iter++; 
            if (iter == 4) {
                if ((len += 3) > *outLen) return 1; 
                *(out++) = (buf >> 16) & 255;
                *(out++) = (buf >> 8) & 255;
                *(out++) = buf & 255;
                buf = 0; iter = 0;

            }   
        }
        progress++;
    }
   
    if (iter == 3) {
        if ((len += 2) > *outLen) return 1; 
        *(out++) = (buf >> 10) & 255;
        *(out++) = (buf >> 2) & 255;
    }
    else if (iter == 2) {
        if (++len > *outLen) return 1; 
        *(out++) = (buf >> 4) & 255;
    }

    *outLen = len;
    return 0;
}

void signalHandler(int sig){
    float fProgress = 100*progress/fs;
    infof("File progress: %lf\n", fProgress);
}