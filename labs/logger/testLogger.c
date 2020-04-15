#include <stdio.h>
#include "logger.h"

int main() {
    infof("INFOF Test %d\n", 10);
    warnf("WARNF Test %s\n", "perro");
    errorf("ERRORF Test %s\n", "gato");
    panicf("PANICF Test %d\n", 1);
    return 0;
}