#include <stdio.h>


int mystrlen(char *str);
char *mystradd(char *origin, char *addition);
int mystrfind(char *origin, char *substr);

int main(int argc, char *argv[]) {
    char* originalStr = argv[1];
    char* addStr = argv[2];
    char* subStr = argv[3];
    char* concStr = mystradd(originalStr, addStr);

    printf("Initial Lenght     : %d\n", mystrlen(originalStr));
    printf("New String         : %s\n", concStr);

    if (mystrfind(concStr, subStr)==1){
        printf("SubString was found: yes\n");
    } else {
        printf("SubString was found: no\n");
    }
    
    return 0;
}
