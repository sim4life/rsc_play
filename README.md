# Files resources Webservice

API to manage files as resources

Write a rest webservice in Go that allows to operate over text files as resources. The minimum requirements for the service are:  
  1. Create a text file with some contents stored in a given path.  
  2. Retrieve the contents of a text file under the given path.  
  3. Replace the contents of a text file.  
  4. Delete the resource that is stored under a given path.  

We would also need to get some statistics per folder basis and retrieve them through another entry point. These statistics are:  
  1. Total number of files in that folder.  
  2. Average number of alphanumeric characters per text file (and standard deviation) in that folder.  
  3. Average word length (and standard deviation) in that folder.  
  4. Total number of bytes stored in that folder.  
  5. Note: All these computations must be calculated recursively from the provided path to the entry point.  

Use all necessary libraries, including third-party libraries but make sure that they are easily fetched. Keep in mind that json is the transport format to be used.
