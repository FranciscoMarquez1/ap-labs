 
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <fcntl.h>
#include <unistd.h>

void filecopy(int fd);

/* cat:  concatenate files, version 2 */
int main(int argc, char *argv[])
{
    int fd;
    void filecopy(int fd);
    char *prog = argv[0];   /* program name for errors */

    if (argc == 1)  /* no args; copy standard input */
        filecopy(0);
    else
        while (--argc > 0)
            if ((fd = open(*++argv, O_RDONLY)) == -1) {
                fprintf(stderr, "%s: can′t open %s\n",
			prog, *argv);
                return 1;
            } else {
                filecopy(fd);
                close(fd);
            }

    if (ferror(stdout)) {
        fprintf(stderr, "%s: error writing stdout\n", prog);
        return 2;
    }

    return 0;
}

void filecopy(int fd)
{

    if(fd == 0){
        int size = 512;
        char buffer[size+1];
        while(read(fd,buffer,size)>0){
            buffer[size] = '\0';
            printf("%s", buffer);
            memset(buffer,0,size);
        }
        return;
    }

    off_t fs = lseek(fd,(size_t)0,SEEK_END);
    lseek(fd,(size_t)0,SEEK_SET);
    char buff[fs+1];
    read(fd, buff, fs);
    write(1, buff, fs);
}
