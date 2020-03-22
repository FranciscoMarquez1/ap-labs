

int mystrlen(char *str){
    int len = 0;
    for(int i=0;str[i]!='\0';i++){
        len++;
    }
    return len;
}

char *mystradd(char *origin, char *addition){
    char* concatenate;
    int lenOr = mystrlen(origin);
    int lenAd = mystrlen(addition);
    int lenConc = lenOr + lenAd;
    concatenate = malloc(sizeof(char)*(lenConc+1));
    int counter = 0;
    for(int i=0;i<lenConc;i++){
        if (i<lenOr){
            concatenate[i]=origin[i];
        } else {
            concatenate[i]=addition[counter];
            counter++;
        }
    }
    concatenate[lenConc] = '\0';
    return concatenate;
}

int mystrfind(char *origin, char *substr){
    int lenOr = mystrlen(origin);
    int lenSu = mystrlen(substr);
    if (lenSu > lenOr) {
        return 0;
    }
    for(int i=0;i<(lenOr);i++){
        if (origin[i]==substr[0]){
            if(lenSu==1){
                return 1;
            }
            for(int j=1;j<lenSu;j++){
                if(origin[i+j]!=substr[j]){
                    break;
                }
            }
            return 1;
        }
    }
    return 0;
}
