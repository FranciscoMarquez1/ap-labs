#include <stdio.h>

#define REPORT_FILE "packages_report.txt"
#define HASH_SIZE 1024

void analizeLog(char *logFile, char *report);

struct Package {
    int key;
    char* name; 
    char* install_date;
    char* last_update;
    char* removal_date;
    int updates;
};

struct Report {
    int installed;
    int removed;
    int upgraded;
};

struct Package* hashArray[HASH_SIZE]; 
struct Package* dummyItem;
struct Package* item;

int hashCode(int key) {
   return key % HASH_SIZE;
}

struct Package *search(int key) {
   int hashIndex = hashCode(key);  
	
   while(hashArray[hashIndex] != NULL) {
	
      if(hashArray[hashIndex]->key == key)
         return hashArray[hashIndex]; 
			
      ++hashIndex;
		
      hashIndex %= HASH_SIZE;
   }        
	
   return NULL;        
}

void insert(int key, char* name, char* install_date, char* last_update, char* removal_date) {

   struct Package *item = (struct Package*) malloc(sizeof(struct Package));
   item->name = name;
   item->install_date = install_date;
   item->last_update = last_update;
   item->removal_date = removal_date;  
   item->updates = 0;  
   item->key = key;

   int hashIndex = hashCode(key);

   while(hashArray[hashIndex] != NULL && hashArray[hashIndex]->key != -1) {
      ++hashIndex;
		
      hashIndex %= HASH_SIZE;
   }
	
   hashArray[hashIndex] = item;
}

struct Package* delete(struct Package* item) {
   int key = item->key;

   int hashIndex = hashCode(key);

   while(hashArray[hashIndex] != NULL) {
	
      if(hashArray[hashIndex]->key == key) {
         struct Package* temp = hashArray[hashIndex]; 
			
         hashArray[hashIndex] = dummyItem; 
         return temp;
      }
		
      ++hashIndex;
		
      hashIndex %= HASH_SIZE;
   }      
	
   return NULL;        
}

int main(int argc, char **argv) {

    if (argc < 2) {
	printf("Usage:./pacman-analizer.o pacman.log\n");
	return 1;
    }

    analizeLog(argv[1], REPORT_FILE);

    return 0;
}

void analizeLog(char *logFile, char *report) {
    printf("Generating Report from: [%s] log file\n", logFile);

    // Implement your solution here.

    printf("Report is generated at: [%s]\n", report);
}
