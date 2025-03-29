#include <stdio.h>
#include <stdlib.h>
int sexec(const char *🤠){
    if(system(🤠)!=0){
        printf("Error executing part: %s\n",🤠);
        return(1);}
    return(0);
}
int main(){
    const char *😸[]={
        "rnodeconf -T /dev/ttyUSB0",
        "rdvce -k && kishell",
        "cXXXXX",
        "sc0",
        "S &<< mcXXXXX ms1 #success!"
    };for(int i=0;i<sizeof(😸)/sizeof(😸[0]);i++){
        if(sexec(😸[i]) != 0){
            printf("Hardware failure.\n");
            return(1);}}
    return(0);}
