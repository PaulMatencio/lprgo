This directory contains go command line
### Available Commands:
- completion
  - Generate the autocompletion script for the specified shell


- gs 
  - Execute ghostscript to convert Pdf or ps to PCL4/5


- gsLpr       
  - Execute ghostscript command to convert a Pdf or PS file to a PCL format  then send the result to a remote printer server 
   
     
- help        
  - Help about any command
  
      
- lpd         
  - Start a lpd server

        
- lpr         
  - Send a file to an lpd server


### Examples

- Start an LPD server  
  - lprgo  lpd  -p 1515 -S /home/paul/lpd/spool/
        

- Send an input PDF  file to  LPD port 1515
  - lprgo  lpr -p 1515 -q p0001 -U paul -f Oct2021.pdf 
        

- transform an input PDF  file to PCL file and send the result to LPD 
  - lprgo gsLpr  -p 1515 -q p0005  -U Michel  -f Oct2021.pdf 
        


