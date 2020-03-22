#include <stdio.h>
#include <string.h>
#include <stdlib.h>

void mergesort(void *lineptr, int l, int r, size_t size, int (*comp)(void *, void *));
void merge(void* lineptr, int l, int m, int r, size_t size, int(*comp)(void *, void *));
int numcmp(void* x, void* y);
int chrcmp(void* x, void* y);

void mergesort(void *lineptr, int l, int r, size_t size, int (*comp)(void *, void *)){
    if (l >= r) { /* do nothing if array contains */
        return;   /* fewer than two elements */
    }
    int m = l+(r-l)/2;
    mergesort(lineptr, l, m, size, comp);
    mergesort(lineptr, m+1, r, size, comp);
    merge(lineptr, l, m, r, size, comp);
}

void merge(void* lineptr, int l, int m, int r, size_t size, int(*comp)(void *, void *)){
    int i, j, k;
    int s1 = m - l + 1;
    int s2 = r -m;
    void* L = malloc(size*s1);
    void* R = malloc(size*s2);

    for(i=0; i < s1; i++){
        memcpy(L+size*i, lineptr+size*(l+i), size);
    }
    for(j=0; j < s2; j++){
        memcpy(R+size*j, lineptr+size*(m+1+j), size);
    }
    i = 0;
    j = 0;
    k = l;

    while(i<s1 && j < s2){
        if(comp(L+size*i, R+size*j)){
            memcpy(lineptr + size*k, L+size*i, size);
            i++;
        } else {
            memcpy(lineptr + size*k, R+size*j, size);
            j++;
        }
        k++;
    }
    while(i < s1) {
        memcpy(lineptr+size*k, L+size*i, size);
        i++;
        k++;
    }
    while(j < s2){
        memcpy(lineptr+size*k, R+size*j, size);
        j++;
        k++;
    }
    free(L);
    free(R);
}

/* numcmp: compare x and y numerically */
int numcmp(void* x, void* y){
    int * a = (int*)x;
    int * b = (int*)y;
    if (*a <= *b)
        return 1;
    else
        return 0;
}

int chrcmp(void* x, void* y){
    char * a = (char*)x;
    char * b = (char*)y;
    if (*a <= *b)
        return 1;
    else
        return 0;
}

/* sort input lines */
int main(int argc, char *argv) {
    printf("Ints\n");
    int a[] = {9,6,7,6,2,3};
    for(int i = 0; i<6;i++){
        printf("%d ", a[i]);
    }
    printf("\n");
    mergesort(a, 0, 5, sizeof(int), numcmp);
    for(int i = 0; i<6;i++){
        printf("%d ", a[i]);
    }
    printf("\n");

    printf("Chars\n");
    char b[] = {'o', 'b', 'e', 'd'};
    for(int i = 0; i<4;i++){
        printf("%c ", b[i]);
    }
    printf("\n");
    mergesort(b, 0, 3, sizeof(char), chrcmp);
    for(int i = 0; i<4;i++){
        printf("%c ", b[i]);
    }
    printf("\n");

    return 0;
}