#include <stdio.h>
#include <stdlib.h>

static char daytab[2][13] = {
    {0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31},
    {0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
};

const char * names[] = {
       "Illegal month",
       "January", "February", "March",
       "April", "May", "June",
       "July", "August", "September",
       "October", "November", "December"
   };
/* month_day function's prototype*/
void month_day(int year, int yearday, int *pmonth, int *pday);

int main(int argc, char **argv) {
    int month;
    int day;
    month_day(atoi(argv[1]), atoi(argv[2]), &month, &day);
    printf("%s %d, %d\n", names[month], day, atoi(argv[1]));
    return 0;
}

void month_day(int year, int yearday, int *pmonth, int *pday){
	int i, leap, counter;
    i = 1;
    counter = 0;
    leap = year%4 == 0 && year%100 != 0 || year%400 == 0;
    for (i = 1; i <= 12; i++)
    {
        if (yearday < (counter + daytab[leap][i])){
            *pmonth = i;
            *pday = yearday - counter;
            break;
        }
        else{
            counter += daytab[leap][i];
        }
    }
    
}
