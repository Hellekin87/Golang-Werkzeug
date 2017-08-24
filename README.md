# Golang-Werkzeug
Automatic Build Tool for Golang 

This tool checks for changes on your go-files in a repository and executes or builds your Go-Projekt constantly. The output of your Go-Application is redirect to the console. For a better usage add the executable to your environment variables .

<b>Examples:</b>

Builds or runs the target_file when it changes: <br>
<pre> werkzeug [run|build] -f [target_file] </pre>

Builds or runs the target_file when one of the files within the directory are changing:<br>
<pre> werkzeug [run|build] -f [target_file] -all </pre>

Builds or runs the target_file with input arguments (don't forget to surround your args with those -> ""):<br>
<pre> werkzeug [run|build] -f [target_file] -arg " -t 'C:/Users/...' " </pre>

Builds or runs the file which was changed in directory <br> 
<pre> werkzeug [run|build] -d [target_dir] </pre>





