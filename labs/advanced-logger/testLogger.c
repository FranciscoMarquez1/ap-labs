#include <stdio.h>
#include "logger.h"

int main(){

    // default logging
    infof("INFOF Test %d\n", 10);
    warnf("WARNF Test %s\n", "perro");
    errorf("ERRORF Test %s\n", "gato");

    // stdout logging
    initLogger("stdout");
    infof("INFOF Test %d\n", 10);
    warnf("WARNF Test %s\n", "perro");
    errorf("ERRORF Test %s\n", "gato");

    // syslog logging
    initLogger("syslog");
    infof("INFOF Test %d\n", 10);
    warnf("WARNF Test %s\n", "perro");
    errorf("ERRORF Test %s\n", "gato");
    

    return 0;
}
