 package main

 import (
         "fmt"
         "os"
         "path/filepath"
         "time"
         "strings"
         "os/exec"
         "bytes"        
 )

type GO_FILE struct {
    name string
    size int64
}

 var (
         targetFile string
         targetFolder string
         build_option string
         last_go_files []GO_FILE
         dir_mode bool = false
         all_files bool = false
         arg_string string
 )

func contains(gf []GO_FILE, fn string) bool {
    for _, a := range gf {  
        if a.name == fn {
            return true
        }
    }
    return false
}

func sizeChanged(gf []GO_FILE, size int64) bool {
    for _, a := range gf {     
        if a.size == size {
            return false
        }
    }
    return true
}

func setSizeOnIndex(gf []GO_FILE, newsize int64, where string){
   for i, a := range gf {     
         if a.name == where{
            gf[i].size = newsize
         }
    }

}

func trim(s string) string{
	ret := strings.Replace(s, "'", "", -1)
	return ret
}

 func FileWalkerAll(path string, fileInfo os.FileInfo, err error) error {

         if err != nil {
                 fmt.Println(err)
                 return nil
         }

         absolute, err := filepath.Abs(path)

         if err != nil {
                 fmt.Println(err)
                 return nil
         }

         if fileInfo.IsDir() {
                 testDir, err := os.Open(absolute)

                 if err != nil {
                         if os.IsPermission(err) {
                                 fmt.Println("No permission to scan ... ", absolute)
                                 fmt.Println(err)
                         }
                 }
                 testDir.Close()
                 return nil
         } else {
                if filepath.Ext(path) == ".go"{
                    if contains(last_go_files,fileInfo.Name()) {
                         if sizeChanged(last_go_files,fileInfo.Size()){

                            fmt.Println("File Size of",fileInfo.Name(),"changed to",fileInfo.Size())
                            setSizeOnIndex(last_go_files, fileInfo.Size(), fileInfo.Name())

                            if(dir_mode){
                            	targetFile = path
                                RUN_CMD(build_option,path)
                            }else{
                                if all_files{
                                         RUN_CMD(build_option,targetFile)
                                    }else if targetFile == path {
                                         RUN_CMD(build_option,targetFile)
                                    }                               
                            }
                         }
                    }else{
                         fmt.Println("File",fileInfo.Name(),"added..")
                         last_go_files = append(last_go_files, GO_FILE{ fileInfo.Name(), fileInfo.Size() })
                    }               
                }
         }
         return nil
 }

 func main() {
  
         if len(os.Args) < 2 {
         		 fmt.Printf("--------------------------------------------------\n")
                 fmt.Printf(" Usage of Werkzeug:\n")
                 fmt.Printf("--------------------------------------------------\n")
                 fmt.Printf(" [run|build] -f <target_file> (-all, -arg) \n")
                 fmt.Printf(" [run|build] -d <target_dir> \n")
                 fmt.Printf("--------------------------------------------------\n")
                 os.Exit(0)
         }
         
        fmt.Println("\nWerkzeug started...")

        for i, arg := range os.Args {
            if arg == "-f"{
                targetFile = os.Args[i+1]
            }
            if arg == "-d"{
                targetFolder = os.Args[i+1]
                dir_mode = true
            }
            if arg == "-all"{
                all_files = true
            }
            if strings.ToLower(arg) == "run"{
                build_option = strings.ToLower(arg)
            }
            if strings.ToLower(arg) == "build"{
                build_option = strings.ToLower(arg)
            }
            if strings.ToLower(arg) == "-arg"{
                arg_string = os.Args[i+1]
            }

         }
       
         if build_option == ""{
            fmt.Println("[ERROR] Missing build method [run|build]")
            os.Exit(-1)
         }

         if build_option == "build" && arg_string != ""{
         	fmt.Println("[ERROR] Can't use method 'build' with flag -arg")
            os.Exit(-1)
         }

        if targetFolder == ""{
                  if info, err := os.Stat(targetFile); err == nil && info.IsDir() {
                    fmt.Println("[ERROR] Used -f option but given a directory")
                    os.Exit(-1)
                  }else{
                        targetFolder = filepath.Dir(targetFile)
                  }  
         }

         if !dir_mode {
         	 if targetFile ==""{
         	 		fmt.Println("[ERROR] No TARGETFILE was set")
                	os.Exit(-1)
         	 	}else{
         	 		fmt.Println("TARGET_FILE: [", targetFile,"]")
         	 	}             
         }else{
            if targetFile != ""{
                fmt.Println("[ERROR] Flags -f and -d combined are not allowed")
                os.Exit(-1)
            }
         }

         fmt.Println("TARGET_FOLDER: [", targetFolder,"]")

         fmt.Println("Starting from directory [", targetFolder, "]")
         
         testFile, err := os.Open(targetFolder)
         if err != nil {
                 fmt.Println(err)
                 os.Exit(-1)
         }
         defer testFile.Close()

         testFileInfo, _ := testFile.Stat()
         if !testFileInfo.IsDir() {
                 fmt.Println("[ERROR]",targetFolder, " is not a directory!")
                 os.Exit(-1)
         }
    
        for{
            err = filepath.Walk(targetFolder, FileWalkerAll)

            if err != nil {
                     fmt.Println(err)
                     os.Exit(-1)
            }
            time.Sleep(1000 * time.Millisecond)
        }
 } 


func RUN_CMD(buildoption string,filename string){
    
    app := "go"
    fmt.Println("Build option:",buildoption)
    final_args := []string{ buildoption, filename}

    if arg_string != ""{
         fmt.Println("Running with args:",arg_string)
    
    		current := strings.Split(arg_string, " ")
         	
         	wait := false
         	s_builder := ""
         	for _, elem := range current {       		
         		if strings.Contains(elem, "'"){
         			num := strings.Count(elem, "'")
         			if num > 1{
         				final_args = append(final_args, trim(elem))
         			}else{
         				s_builder += elem+" "
         				if !wait{
         					wait = true
         				}else{
         					final_args = append(final_args,trim(s_builder))
         					s_builder = ""
        					wait = false
         				}
         			}
         		}else{
         				if wait {
         					if elem == ""{
         						s_builder += " "
         					}else{
         						s_builder += elem+" "
         					}
         				}else{
         					elem = strings.Replace(elem, " ", "", -1)         		
	         				if elem != "" {
	         					final_args = append(final_args, elem)	
	         				}
         				}        			
        			}  	
         	}
    }

 	if buildoption == "run"{
        fmt.Println("\n****************************************************************")
        fmt.Println(targetFile)
        fmt.Println("****************************************************************\n")
    }
    
    cmd := exec.Command(app, final_args...)
   
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err := cmd.Run()
    
  	if err != nil {
    	fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
    	fmt.Println("Searching for changes...")
    	return
  	}
  	if buildoption == "run"{
        fmt.Println(out.String())

        fmt.Println("\n****************************************************************")
        fmt.Println("***************************    END   ***************************")
        fmt.Println("****************************************************************\n")
  	}else{
       fmt.Println("Creating Executable!....")
  	}
    fmt.Println("Searching for changes...")
}
