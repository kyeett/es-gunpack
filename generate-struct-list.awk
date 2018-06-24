BEGIN {
   print "Finding structs!\n"

   print "switch signalType {"
}
/^package [a-z]+/ { PACKAGE_NAME=substr($2, 1, length($2));}
PACKAGE_NAME != "" && ( $1 == "type" ) && ( $3 == "struct" ) {print "case \"" PACKAGE_NAME "." $2 "\":\n   newTest := &" PACKAGE_NAME "." $2 "{}\n"}
END {
    print "default:\n   fmt.Println(\"Linux.\")\n}"

}